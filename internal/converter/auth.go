package converter

import (
	"github.com/VadimGossip/concoleChat-auth/internal/model"
	desc "github.com/VadimGossip/concoleChat-auth/pkg/auth_v1"
)

func ToLoginUserInfoFromDesc(info *desc.LoginRequest) *model.LoginUserInfo {
	return &model.LoginUserInfo{
		Name:     info.Username,
		Password: info.Password,
	}
}
