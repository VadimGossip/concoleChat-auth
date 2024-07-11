package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
	def "github.com/VadimGossip/concoleChat-auth/internal/repository"
	"github.com/VadimGossip/concoleChat-auth/internal/repository/user/converter"
	repoModel "github.com/VadimGossip/concoleChat-auth/internal/repository/user/model"
)

var _ def.UserRepository = (*repository)(nil)

type repository struct {
	m      sync.RWMutex
	data   map[int64]*repoModel.User
	lastID int64
}

func NewRepository() *repository {
	return &repository{
		data: make(map[int64]*repoModel.User),
	}
}

func (r *repository) Create(_ context.Context, info *model.UserInfo) (int64, error) {
	r.m.Lock()
	defer r.m.Unlock()

	r.lastID++
	r.data[r.lastID] = &repoModel.User{
		ID:        r.lastID,
		Info:      converter.ToRepoFromUserInfo(info),
		CreatedAt: time.Now(),
	}
	logrus.Infof("User created %+v", info)
	return r.lastID, nil
}

func (r *repository) Get(_ context.Context, id int64) (*model.User, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	if user, ok := r.data[id]; ok {
		return converter.ToUserFromRepo(user), nil
	}
	return nil, fmt.Errorf("user id=%d not found", id)
}

func (r *repository) Update(_ context.Context, id int64, updateInfo *model.UpdateUserInfo) error {
	r.m.Lock()
	defer r.m.Unlock()
	if user, ok := r.data[id]; ok {
		user.UpdatedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		user.Info.Name = updateInfo.Name
		user.Info.Email = updateInfo.Email

		if updateInfo.Role != model.UnknownRole {
			user.Info.Role = updateInfo.Role
		}
		logrus.Infof("User updated %+v", r.data[id])
		return nil
	}
	return fmt.Errorf("user id=%d not found", id)
}

func (r *repository) Delete(_ context.Context, id int64) error {
	r.m.Lock()
	defer r.m.Unlock()
	if _, ok := r.data[id]; ok {
		delete(r.data, id)
		logrus.Infof("User id=%d deleted", id)
		return nil
	}
	return fmt.Errorf("user id=%d not found", id)
}
