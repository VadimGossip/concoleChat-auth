package user

import (
	"context"
	"github.com/VadimGossip/concoleChat-auth/internal/client/db"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
	"github.com/VadimGossip/concoleChat-auth/internal/repository"
	def "github.com/VadimGossip/concoleChat-auth/internal/service"
	"github.com/VadimGossip/concoleChat-auth/internal/service/user/validator"
)

var _ def.UserService = (*service)(nil)

type service struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

func NewService(userRepository repository.UserRepository, txManager db.TxManager) *service {
	return &service{
		userRepository: userRepository,
		txManager:      txManager,
	}
}

func (s *service) Create(ctx context.Context, info *model.UserInfo) (int64, error) {
	if err := validator.CreateValidation(info); err != nil {
		return 0, err
	}
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error
		id, txErr = s.userRepository.Create(ctx, info)
		if txErr != nil {
			return txErr
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *service) Get(ctx context.Context, ID int64) (*model.User, error) {
	user := &model.User{}
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error
		user, txErr = s.userRepository.Get(ctx, ID)
		if txErr != nil {
			return txErr
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) Update(ctx context.Context, ID int64, updateInfo *model.UpdateUserInfo) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error
		if txErr = s.userRepository.Update(ctx, ID, updateInfo); txErr != nil {
			return txErr
		}
		return nil
	})
}

func (s *service) Delete(ctx context.Context, ID int64) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error
		if txErr = s.userRepository.Delete(ctx, ID); txErr != nil {
			return txErr
		}
		return nil
	})
}
