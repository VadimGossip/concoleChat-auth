package converter

import (
	"time"

	repoModel "github.com/VadimGossip/concoleChat-auth/internal/repository/user/pg/model"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

func ToUserFromRepo(user *repoModel.User) *model.User {
	var updatedAt *time.Time
	if user.UpdatedAt.Valid {
		updatedAt = &user.UpdatedAt.Time
	}

	return &model.User{
		ID:        user.ID,
		Info:      ToUserInfoFromRepo(user.Info),
		CreatedAt: user.CreatedAt,
		UpdatedAt: updatedAt,
	}
}

func ToUserInfoFromRepo(info repoModel.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  info.Role,
	}
}

func ToRepoFromUserInfo(info *model.UserInfo) repoModel.UserInfo {
	return repoModel.UserInfo{
		Name:     info.Name,
		Email:    info.Email,
		Password: info.Password,
		Role:     info.Role,
	}
}
