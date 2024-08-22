package service

import (
	"context"
	"time"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

type UserService interface {
	Create(ctx context.Context, info *model.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
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

type AuthService interface {
	Login(ctx context.Context, info *model.LoginUserInfo) (string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	AccessToken(ctx context.Context, refreshToken string) (string, error)
}

type AccessService interface {
	Check(ctx context.Context, accessToken, endpointAddress string) error
}

type TokenService interface {
	Generate(info model.UserInfo, secretKey []byte, duration time.Duration) (string, error)
	Verify(tokenStr string, secretKey []byte) (*model.UserClaims, error)
}

type PasswordService interface {
	Verify(hashedPassword string, candidatePassword string) bool
	Hash(password string) (string, error)
}
