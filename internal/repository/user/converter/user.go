package converter

import (
	"time"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
	repoModel "github.com/VadimGossip/concoleChat-auth/internal/repository/user/model"
)

func ToUserFromRepo(user *repoModel.User) *model.User {
	var updatedAt *time.Time
	if user.UpdatedAt.Valid {
		updatedAt = &user.UpdatedAt.Time
	}

	return &model.User{
		Id:        user.Id,
		Info:      ToUserInfoFromRepo(user.Info),
		CreatedAt: user.CreatedAt,
		UpdatedAt: updatedAt,
	}
}

func ToUserInfoFromRepo(info repoModel.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:            info.Name,
		Email:           info.Email,
		Password:        info.Password,
		PasswordConfirm: info.PasswordConfirm,
		Role:            info.Role,
	}
}

func ToRepoFromUserInfo(info *model.UserInfo) repoModel.UserInfo {
	return repoModel.UserInfo{
		Name:            info.Name,
		Email:           info.Email,
		Password:        info.Password,
		PasswordConfirm: info.PasswordConfirm,
		Role:            info.Role,
	}
}
