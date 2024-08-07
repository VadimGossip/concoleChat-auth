package pg

import (
	"context"
	"strconv"

	db "github.com/VadimGossip/concoleChat-auth/internal/client/db/redis"
	"github.com/VadimGossip/concoleChat-auth/internal/model"
	def "github.com/VadimGossip/concoleChat-auth/internal/repository"
	"github.com/VadimGossip/concoleChat-auth/internal/repository/user/redis/converter"
	repoModel "github.com/VadimGossip/concoleChat-auth/internal/repository/user/redis/model"
)

var _ def.UserCacheRepository = (*repository)(nil)

type repository struct {
	db db.Client
}

func NewRepository(db db.Client) *repository {
	return &repository{db: db}
}

func (r *repository) Get(ctx context.Context, ID int64) (*model.User, error) {
	repoUser := &repoModel.User{}
	err := r.db.DB().HGetAll(ctx, strconv.FormatInt(ID, 10), repoUser)
	if err != nil {
		return nil, err
	}
	return converter.ToUserFromRepo(repoUser), nil
}

func (r *repository) Delete(ctx context.Context, ID int64) error {
	return r.db.DB().Del(ctx, strconv.FormatInt(ID, 10))
}
