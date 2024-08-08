package user

import (
	"context"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

func (s *service) Get(ctx context.Context, ID int64) (*model.User, error) {
	return s.userCacheRepository.Get(ctx, ID)
}
