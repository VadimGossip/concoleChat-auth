package service

import (
	"context"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

type UserService interface {
	Create(ctx context.Context, info *model.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, id int64, updateInfo *model.UpdateUserInfo) error
	Delete(ctx context.Context, id int64) error
}

type AuditService interface {
	Create(ctx context.Context, audit *model.Audit) error
}

type UserCacheService interface {
	Get(ctx context.Context, ID int64) (*model.User, error)
	Set(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, ID int64) error
}

type UserConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type UserProducerService interface {
	ProduceCreate(info *model.UserInfo) error
}
