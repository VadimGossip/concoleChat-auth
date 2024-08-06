package user

import (
	"context"
	"fmt"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

func (s *service) Delete(ctx context.Context, ID int64) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error
		if txErr = s.userRepository.Delete(ctx, ID); txErr != nil {
			return txErr
		}

		txErr = s.auditService.Create(ctx, &model.Audit{
			Action:     "delete user",
			CallParams: fmt.Sprintf("id %d", ID),
		})
		if txErr != nil {
			return txErr
		}

		return nil
	})
}
