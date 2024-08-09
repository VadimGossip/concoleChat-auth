package user

import (
	"context"
)

func (s *service) Delete(ctx context.Context, ID int64) error {
	return s.userCacheRepository.Delete(ctx, ID)
}
