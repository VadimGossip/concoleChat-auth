package app

import (
	"context"
	"fmt"

	"github.com/VadimGossip/concoleChat-auth/internal/api/user"
	"github.com/VadimGossip/concoleChat-auth/internal/closer"
	"github.com/VadimGossip/concoleChat-auth/internal/model"
	"github.com/VadimGossip/concoleChat-auth/internal/repository"
	userRepo "github.com/VadimGossip/concoleChat-auth/internal/repository/user"
	"github.com/VadimGossip/concoleChat-auth/internal/service"
	userService "github.com/VadimGossip/concoleChat-auth/internal/service/user"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type serviceProvider struct {
	cfg *model.Config

	db       *pgx.Conn
	userRepo repository.UserRepository

	userService service.UserService

	userImpl *user.Implementation
}

func newServiceProvider(cfg *model.Config) *serviceProvider {
	return &serviceProvider{cfg: cfg}
}

func (s *serviceProvider) DBClient(ctx context.Context) *pgx.Conn {
	if s.db == nil {
		dbDSN := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s", s.cfg.Db.Host, s.cfg.Db.Port, s.cfg.Db.Name, s.cfg.Db.Username, s.cfg.Db.Password, s.cfg.Db.SSLMode)
		db, err := pgx.Connect(ctx, dbDSN)
		if err != nil {
			logrus.Fatalf("failed to connect to database: %v", err)
		}
		if err = db.Ping(ctx); err != nil {
			logrus.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(func() error {
			return db.Close(ctx)
		})
		s.db = db
	}

	return s.db
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepo == nil {
		s.userRepo = userRepo.NewRepository(s.DBClient(ctx))
	}
	return s.userRepo
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository(ctx))
	}
	return s.userService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}
