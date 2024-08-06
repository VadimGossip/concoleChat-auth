package user

import (
	"context"
	"fmt"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

func (s *service) Get(ctx context.Context, ID int64) (*model.User, error) {
	user := &model.User{}
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error
		user, txErr = s.userRepository.Get(ctx, ID)
		if txErr != nil {
			return txErr
		}

		if txErr = s.auditService.Create(ctx, &model.Audit{
			Action:     "get user",
			CallParams: fmt.Sprintf("id %d", ID),
		}); txErr != nil {
			return txErr
		}

		return txErr
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}
