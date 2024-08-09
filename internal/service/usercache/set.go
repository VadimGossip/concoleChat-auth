package user

import (
	"context"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

func (s *service) Set(ctx context.Context, user *model.User) error {
	return s.userCacheRepository.Set(ctx, user, s.userCacheCfg.Expire)
}
