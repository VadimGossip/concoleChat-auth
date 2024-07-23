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

	tx, err := s.userRepository.BeginTxSerializable(ctx)
	if err != nil {
		return 0, err
	}

	id, err := s.userRepository.Create(ctx, tx, info)
	if err != nil {
		return 0, s.userRepository.StopTx(ctx, tx, err)
	}

	return id, s.userRepository.StopTx(ctx, tx, nil)
}

func (s *service) Get(ctx context.Context, ID int64) (*model.User, error) {
	tx, err := s.userRepository.BeginTxSerializable(ctx)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.Get(ctx, tx, ID)
	if err != nil {
		return nil, s.userRepository.StopTx(ctx, tx, err)
	}

	return user, s.userRepository.StopTx(ctx, tx, nil)
}

func (s *service) Update(ctx context.Context, ID int64, updateInfo *model.UpdateUserInfo) error {
	if err := validator.UpdateValidation(updateInfo); err != nil {
		return err
	}

	tx, err := s.userRepository.BeginTxSerializable(ctx)
	if err != nil {
		return err
	}
	return s.userRepository.StopTx(ctx, tx, s.userRepository.Update(ctx, tx, ID, updateInfo))
}

func (s *service) Delete(ctx context.Context, ID int64) error {
	tx, err := s.userRepository.BeginTxSerializable(ctx)
	if err != nil {
		return err
	}

	return s.userRepository.StopTx(ctx, tx, s.userRepository.Delete(ctx, tx, ID))
}
