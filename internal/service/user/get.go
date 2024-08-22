package user

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

func (s *service) Get(ctx context.Context, ID int64) (*model.User, error) {
	user, err := s.userCacheService.Get(ctx, ID)
	if err != nil {
		logrus.Infof("User cache service err %s on get user id = %d", err, ID)
	}
	if user != nil {
		return user, nil
	}

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
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

	if err = s.userCacheService.Set(ctx, user); err != nil {
		logrus.Infof("User cache service err %s on set user = %+v", err, user)
	}

	return user, nil
}

func (s *service) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.userRepository.GetByUsername(ctx, username)
}
