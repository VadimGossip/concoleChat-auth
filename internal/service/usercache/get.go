package user

import (
	"context"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

func (s *service) Get(ctx context.Context, ID int64) (*model.User, error) {
	user, err := s.userCacheRepository.Get(ctx, ID)
	if err != nil {
		return nil, err
	}

	if user.ID != ID {
		return nil, err
	}
	return user, nil
}
