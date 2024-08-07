package app

import (
	"context"
	"fmt"
	userRepo "github.com/VadimGossip/concoleChat-auth/internal/repository/user/pg"
	"log"

	"github.com/VadimGossip/platform_common/pkg/closer"
	db "github.com/VadimGossip/platform_common/pkg/db/postgres"
	"github.com/VadimGossip/platform_common/pkg/db/postgres/pg"
	"github.com/VadimGossip/platform_common/pkg/db/postgres/transaction"
	"github.com/sirupsen/logrus"

	"github.com/VadimGossip/concoleChat-auth/internal/api/user"
	"github.com/VadimGossip/concoleChat-auth/internal/model"
	"github.com/VadimGossip/concoleChat-auth/internal/repository"
	auditRepo "github.com/VadimGossip/concoleChat-auth/internal/repository/audit"
	"github.com/VadimGossip/concoleChat-auth/internal/service"
	auditService "github.com/VadimGossip/concoleChat-auth/internal/service/audit"
	userService "github.com/VadimGossip/concoleChat-auth/internal/service/user"
)

type serviceProvider struct {
	cfg *model.Config

	dbClient  db.Client
	txManager db.TxManager
	auditRepo repository.AuditRepository
	userRepo  repository.UserRepository

	auditService service.AuditService
	userService  service.UserService

	userImpl *user.Implementation
}

func newServiceProvider(cfg *model.Config) *serviceProvider {
	return &serviceProvider{cfg: cfg}
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		dbDSN := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s", s.cfg.Db.Host, s.cfg.Db.Port, s.cfg.Db.Name, s.cfg.Db.Username, s.cfg.Db.Password, s.cfg.Db.SSLMode)
		cl, err := pg.New(ctx, dbDSN)
		if err != nil {
			logrus.Fatalf("failed to create db client: %s", err)
		}

		if err = cl.DB().Ping(ctx); err != nil {
			log.Fatalf("ping error: %s", err)
		}
		closer.Add(cl.Close)
		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) AuditRepository(ctx context.Context) repository.AuditRepository {
	if s.auditRepo == nil {
		s.auditRepo = auditRepo.NewRepository(s.DBClient(ctx))
	}
	return s.auditRepo
}

func (s *serviceProvider) AuditService(ctx context.Context) service.AuditService {
	if s.auditService == nil {
		s.auditService = auditService.NewService(s.AuditRepository(ctx))
	}
	return s.auditService
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepo == nil {
		s.userRepo = userRepo.NewRepository(s.DBClient(ctx))
	}
	return s.userRepo
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository(ctx), s.AuditService(ctx), s.TxManager(ctx))
	}
	return s.userService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}
