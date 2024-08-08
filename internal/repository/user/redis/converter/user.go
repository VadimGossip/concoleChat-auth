package converter

import (
	repoModel "github.com/VadimGossip/concoleChat-auth/internal/repository/user/redis/model"
	"time"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

func ToUserFromRepo(user *repoModel.User) *model.User {
	var updatedAt *time.Time
	if user.UpdatedAt != nil {
		tmpUpdatedAt := time.Unix(*user.UpdatedAt, 0)
		updatedAt = &tmpUpdatedAt
	}

	return &model.User{
		ID: user.ID,
		Info: model.UserInfo{
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
			Role:     user.Role,
		},
		CreatedAt: time.Unix(user.CreatedAt, 0),
		UpdatedAt: updatedAt,
	}
}

func ToRepoFromUser(user *model.User) *repoModel.User {
	var updatedAt *int64
	if user.UpdatedAt != nil {
		tmpUpdatedAt := user.UpdatedAt.Unix()
		updatedAt = &tmpUpdatedAt
	}

	return &repoModel.User{
		ID:        user.ID,
		Name:      user.Info.Name,
		Email:     user.Info.Email,
		Password:  user.Info.Password,
		Role:      user.Info.Role,
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: updatedAt,
	}
}
