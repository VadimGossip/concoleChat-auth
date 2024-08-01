package repository

import (
	"context"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, info *model.UserInfo) (int64, error)
	Get(ctx context.Context, ID int64) (*model.User, error)
	Update(ctx context.Context, ID int64, updateInfo *model.UpdateUserInfo) error
	Delete(ctx context.Context, ID int64) error
}

type AuditRepository interface {
	Create(ctx context.Context, audit *model.Audit) error
}
