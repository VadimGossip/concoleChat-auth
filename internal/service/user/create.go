package user

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
	"github.com/VadimGossip/concoleChat-auth/internal/service/user/validator"
)

func (s *service) Create(ctx context.Context, info *model.UserInfo) (int64, error) {
	if err := validator.CreateValidation(info); err != nil {
		return 0, err
	}
	user := &model.User{}
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error
		user, txErr = s.userRepository.Create(ctx, info)
		if txErr != nil {
			return txErr
		}

		txErr = s.auditService.Create(ctx, &model.Audit{
			Action:     "create user",
			CallParams: fmt.Sprintf("info %+v", info),
		})
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	if err = s.userCacheService.Set(ctx, user); err != nil {
		logrus.Infof("User cache service err %s on set user = %+v", err, user)
	}

	return user.ID, nil
}
