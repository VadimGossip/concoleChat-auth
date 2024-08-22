package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	db "github.com/VadimGossip/platform_common/pkg/db/postgres"
	"github.com/jackc/pgx/v4"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
	def "github.com/VadimGossip/concoleChat-auth/internal/repository"
	"github.com/VadimGossip/concoleChat-auth/internal/repository/user/pg/converter"
	repoModel "github.com/VadimGossip/concoleChat-auth/internal/repository/user/pg/model"
)

const (
	usersTableName  string = "users"
	idColumn        string = "id"
	usernameColumn  string = "username"
	passwordColumn  string = "password"
	emailColumn     string = "email"
	roleColumn      string = "role"
	createdAtColumn string = "created_at"
	updatedAtColumn string = "updated_at"
	repoName        string = "user_repository"
)

var _ def.UserRepository = (*repository)(nil)

type repository struct {
	db db.Client
}

func NewRepository(db db.Client) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, info *model.UserInfo) (*model.User, error) {
	repoInfo := converter.ToRepoFromUserInfo(info)
	userInsert := sq.Insert(usersTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(usernameColumn, passwordColumn, emailColumn, roleColumn, createdAtColumn).
		Values(repoInfo.Name, repoInfo.Password, repoInfo.Email, repoInfo.Role, time.Now()).
		Suffix("RETURNING " + idColumn + "," + usernameColumn + "," + passwordColumn + "," + emailColumn + "," + roleColumn + "," + createdAtColumn)

	query, args, err := userInsert.ToSql()
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	q := db.Query{
		Name:     repoName + ".Create",
		QueryRaw: query,
	}
	if err = r.db.DB().QueryRowContext(ctx, q, args...).
		Scan(&user.ID,
			&user.Info.Name,
			&user.Info.Password,
			&user.Info.Email,
			&user.Info.Role,
			&user.CreatedAt); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) Get(ctx context.Context, ID int64) (*model.User, error) {
	userSelect := sq.Select(usernameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(usersTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: ID})

	query, args, err := userSelect.ToSql()
	if err != nil {
		return nil, err
	}

	repoUser := &repoModel.User{ID: ID}
	q := db.Query{
		Name:     repoName + ".Get",
		QueryRaw: query,
	}
	if err = r.db.DB().ScanOneContext(ctx, repoUser, q, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return converter.ToUserFromRepo(repoUser), nil
}

func (r *repository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	userSelect := sq.Select(idColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(usersTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{usernameColumn: username})

	query, args, err := userSelect.ToSql()
	if err != nil {
		return nil, err
	}

	repoUser := &repoModel.User{Info: repoModel.UserInfo{Name: username}}
	q := db.Query{
		Name:     repoName + ".GetByUsername",
		QueryRaw: query,
	}

	if err = r.db.DB().ScanOneContext(ctx, repoUser, q, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return converter.ToUserFromRepo(repoUser), nil
}

func (r *repository) Update(ctx context.Context, ID int64, updateInfo *model.UpdateUserInfo) error {
	userUpdate := sq.Update(usersTableName).
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", time.Now()).
		Where(sq.Eq{idColumn: ID})

	if updateInfo.Name != nil {
		userUpdate = userUpdate.Set(usernameColumn, *updateInfo.Name)
	}

	if updateInfo.Email != nil {
		userUpdate = userUpdate.Set(emailColumn, *updateInfo.Email)
	}
	if updateInfo.Role != model.UnknownRole {
		userUpdate = userUpdate.Set(roleColumn, updateInfo.Role)
	}

	query, args, err := userUpdate.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     repoName + ".Update",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, ID int64) error {
	chatDelete := sq.Delete(usersTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: ID})

	query, args, err := chatDelete.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     repoName + ".Delete",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}
