package user

import (
	"context"
	"fmt"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

func (s *service) Update(ctx context.Context, ID int64, updateInfo *model.UpdateUserInfo) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error
		if txErr = s.userRepository.Update(ctx, ID, updateInfo); txErr != nil {
			return txErr
		}
		return s.auditService.Create(ctx, &model.Audit{
			Action:     "update user",
			CallParams: fmt.Sprintf("id %d updateInfo %+v", ID, updateInfo),
		})
	})
}
