package app

import (
	"context"
	"fmt"

	"log"
	"time"

	"github.com/VadimGossip/platform_common/pkg/closer"
	"github.com/VadimGossip/platform_common/pkg/db/postgres"
	"github.com/VadimGossip/platform_common/pkg/db/postgres/pg"
	"github.com/VadimGossip/platform_common/pkg/db/postgres/transaction"
	"github.com/VadimGossip/platform_common/pkg/db/redis"
	"github.com/VadimGossip/platform_common/pkg/db/redis/rdb"
	"github.com/sirupsen/logrus"

	"github.com/VadimGossip/concoleChat-auth/internal/api/user"
	"github.com/VadimGossip/concoleChat-auth/internal/model"
	"github.com/VadimGossip/concoleChat-auth/internal/repository"
	auditRepo "github.com/VadimGossip/concoleChat-auth/internal/repository/audit"
	userRepo "github.com/VadimGossip/concoleChat-auth/internal/repository/user/pg"
	userCacheRepo "github.com/VadimGossip/concoleChat-auth/internal/repository/user/redis"
	"github.com/VadimGossip/concoleChat-auth/internal/service"
	auditService "github.com/VadimGossip/concoleChat-auth/internal/service/audit"
	userService "github.com/VadimGossip/concoleChat-auth/internal/service/user"
	userCacheService "github.com/VadimGossip/concoleChat-auth/internal/service/usercache"
)

type serviceProvider struct {
	cfg *model.Config

	pgDbClient    postgres.Client
	txManager     postgres.TxManager
	redisDbClient redis.Client
	auditRepo     repository.AuditRepository
	userCacheRepo repository.UserCacheRepository
	userRepo      repository.UserRepository

	auditService     service.AuditService
	userCacheService service.UserCacheService
	userService      service.UserService

	userImpl *user.Implementation
}

func newServiceProvider(cfg *model.Config) *serviceProvider {
	return &serviceProvider{cfg: cfg}
}

func (s *serviceProvider) PgDbClient(ctx context.Context) postgres.Client {
	if s.pgDbClient == nil {
		dbDSN := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s", s.cfg.PgDb.Host, s.cfg.PgDb.Port, s.cfg.PgDb.Name, s.cfg.PgDb.Username, s.cfg.PgDb.Password, s.cfg.PgDb.SSLMode)
		cl, err := pg.New(ctx, dbDSN)
		if err != nil {
			logrus.Fatalf("failed to create db client: %s", err)
		}

		if err = cl.DB().Ping(ctx); err != nil {
			log.Fatalf("ping error: %s", err)
		}
		closer.Add(cl.Close)
		s.pgDbClient = cl
	}

	return s.pgDbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) postgres.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.PgDbClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) RedisDbClient(ctx context.Context) redis.Client {
	if s.redisDbClient == nil {
		cl := rdb.New(&rdb.ClientOptions{
			Addr:         fmt.Sprintf("%s:%d", s.cfg.RedisDb.Host, s.cfg.RedisDb.Port),
			Username:     s.cfg.RedisDb.Username,
			Password:     s.cfg.RedisDb.Password,
			DB:           s.cfg.RedisDb.Db,
			ReadTimeout:  time.Duration(s.cfg.RedisDb.ReadTimeoutSec) * time.Second,
			WriteTimeout: time.Duration(s.cfg.RedisDb.WriteTimeoutSec) * time.Second,
		})

		if err := cl.DB().Ping(ctx); err != nil {
			log.Fatalf("kdb ping error: %s", err)
		}

		closer.Add(cl.Close)
		s.redisDbClient = cl
	}

	return s.redisDbClient
}

func (s *serviceProvider) AuditRepository(ctx context.Context) repository.AuditRepository {
	if s.auditRepo == nil {
		s.auditRepo = auditRepo.NewRepository(s.PgDbClient(ctx))
	}
	return s.auditRepo
}

func (s *serviceProvider) AuditService(ctx context.Context) service.AuditService {
	if s.auditService == nil {
		s.auditService = auditService.NewService(s.AuditRepository(ctx))
	}
	return s.auditService
}

func (s *serviceProvider) UserCacheRepository(ctx context.Context) repository.UserCacheRepository {
	if s.userCacheRepo == nil {
		s.userCacheRepo = userCacheRepo.NewRepository(s.RedisDbClient(ctx))
	}
	return s.userCacheRepo
}

func (s *serviceProvider) UserCacheService(ctx context.Context) service.UserCacheService {
	if s.userCacheService == nil {
		s.userCacheService = userCacheService.NewService(s.cfg.UserCache, s.UserCacheRepository(ctx))
	}
	return s.userCacheService
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepo == nil {
		s.userRepo = userRepo.NewRepository(s.PgDbClient(ctx))
	}
	return s.userRepo
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository(ctx), s.UserCacheService(ctx), s.AuditService(ctx), s.TxManager(ctx))
	}
	return s.userService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}
