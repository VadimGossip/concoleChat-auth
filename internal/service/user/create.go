package user

import (
	"context"
	"fmt"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
	"github.com/VadimGossip/concoleChat-auth/internal/service/user/validator"
)

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

		if txErr = s.auditService.Create(ctx, &model.Audit{
			Action:     "create user",
			CallParams: fmt.Sprintf("info %+v", info),
		}); txErr != nil {
			return txErr
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
