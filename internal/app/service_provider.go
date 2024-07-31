package app

import (
	"context"
	"fmt"
	"github.com/VadimGossip/concoleChat-auth/internal/client/db/transaction"
	"github.com/VadimGossip/concoleChat-auth/internal/closer"
	"log"

	"github.com/VadimGossip/concoleChat-auth/internal/api/user"
	"github.com/VadimGossip/concoleChat-auth/internal/client/db"
	"github.com/VadimGossip/concoleChat-auth/internal/client/db/pg"
	"github.com/VadimGossip/concoleChat-auth/internal/model"
	"github.com/VadimGossip/concoleChat-auth/internal/repository"
	userRepo "github.com/VadimGossip/concoleChat-auth/internal/repository/user"
	"github.com/VadimGossip/concoleChat-auth/internal/service"
	userService "github.com/VadimGossip/concoleChat-auth/internal/service/user"
	"github.com/sirupsen/logrus"
)

type serviceProvider struct {
	cfg *model.Config

	dbClient  db.Client
	txManager db.TxManager
	userRepo  repository.UserRepository

	userService service.UserService

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

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepo == nil {
		s.userRepo = userRepo.NewRepository(s.DBClient(ctx))
	}
	return s.userRepo
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository(ctx), s.TxManager(ctx))
	}
	return s.userService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}
