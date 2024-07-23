package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/VadimGossip/concoleChat-auth/internal/model"
	def "github.com/VadimGossip/concoleChat-auth/internal/repository"
	"github.com/VadimGossip/concoleChat-auth/internal/repository/user/converter"
	repoModel "github.com/VadimGossip/concoleChat-auth/internal/repository/user/model"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

var _ def.UserRepository = (*repository)(nil)

type repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) BeginTxSerializable(ctx context.Context) (pgx.Tx, error) {
	return r.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
}

func (r *repository) StopTx(ctx context.Context, tx pgx.Tx, err error) error {
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			logrus.Errorf("error while rollback transaction: %s", rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

// Create need to hash password
func (r *repository) Create(ctx context.Context, tx pgx.Tx, info *model.UserInfo) (int64, error) {
	repoInfo := converter.ToRepoFromUserInfo(info)
	userInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("username", "password", "email", "role", "created_at").
		Values(repoInfo.Name, repoInfo.Password, repoInfo.Email, repoInfo.Role, time.Now()).
		Suffix("RETURNING id")

	query, args, err := userInsert.ToSql()
	if err != nil {
		return 0, err
	}
	var id int64
	if err = tx.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repository) Get(ctx context.Context, tx pgx.Tx, ID int64) (*model.User, error) {
	userSelect := sq.Select("username", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": ID})

	query, args, err := userSelect.ToSql()
	if err != nil {
		return nil, err
	}

	repoUser := &repoModel.User{ID: ID}
	if err = tx.QueryRow(ctx, query, args...).Scan(&repoUser.Info.Name, &repoUser.Info.Email, &repoUser.Info.Role, &repoUser.CreatedAt, &repoUser.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return converter.ToUserFromRepo(repoUser), nil
}

func (r *repository) Update(ctx context.Context, tx pgx.Tx, ID int64, updateInfo *model.UpdateUserInfo) error {
	userUpdate := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": ID})

	if updateInfo.Name != nil {
		userUpdate = userUpdate.Set("username", *updateInfo.Name)
	}

	if updateInfo.Email != nil {
		userUpdate = userUpdate.Set("email", *updateInfo.Email)
	}
	if updateInfo.Role != model.UnknownRole {
		userUpdate = userUpdate.Set("role", updateInfo.Role)
	}

	query, args, err := userUpdate.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, tx pgx.Tx, ID int64) error {
	chatDelete := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": ID})

	query, args, err := chatDelete.ToSql()
	if err != nil {
		return err
	}

	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
