package repository

import (
	"context"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
	"github.com/jackc/pgx/v4"
)

type UserRepository interface {
	BeginTxSerializable(ctx context.Context) (pgx.Tx, error)
	StopTx(ctx context.Context, tx pgx.Tx, err error) error
	Create(ctx context.Context, tx pgx.Tx, info *model.UserInfo) (int64, error)
	Get(ctx context.Context, tx pgx.Tx, ID int64) (*model.User, error)
	Update(ctx context.Context, tx pgx.Tx, ID int64, updateInfo *model.UpdateUserInfo) error
	Delete(ctx context.Context, tx pgx.Tx, ID int64) error
}
