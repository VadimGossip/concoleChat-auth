package user

import (
	"context"
	"github.com/VadimGossip/concoleChat-auth/internal/model"
	"github.com/VadimGossip/concoleChat-auth/internal/repository"
	def "github.com/VadimGossip/concoleChat-auth/internal/service"
	"github.com/VadimGossip/concoleChat-auth/internal/service/user/validator"
)

var _ def.UserService = (*service)(nil)

type service struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) *service {
	return &service{
		userRepository: userRepository,
	}
}

func (s *service) Create(ctx context.Context, info *model.UserInfo) (int64, error) {
	if err := validator.CreateValidation(info); err != nil {
		return 0, err
	}
	return s.userRepository.Create(ctx, info)
}

func (s *service) Get(ctx context.Context, id int64) (*model.User, error) {
	return s.userRepository.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, id int64, updateInfo *model.UpdateUserInfo) error {
	if err := validator.UpdateValidation(updateInfo); err != nil {
		return err
	}
	return s.userRepository.Update(ctx, id, updateInfo)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.userRepository.Delete(ctx, id)
}
