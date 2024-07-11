package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
	desc "github.com/VadimGossip/concoleChat-auth/pkg/user_v1"
)

func ToUserFromService(user *model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt != nil {
		updatedAt = timestamppb.New(*user.UpdatedAt)
	}

	return &desc.User{
		Id:        user.ID,
		Info:      ToUserInfoFromService(user.Info),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToUserInfoFromService(info model.UserInfo) *desc.UserInfo {
	return &desc.UserInfo{
		Name:            info.Name,
		Email:           info.Email,
		Password:        info.Password,
		PasswordConfirm: info.PasswordConfirm,
		Role:            desc.Role(desc.Role_value[info.Role]),
	}
}

func ToUserInfoFromDesc(info *desc.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:            info.Name,
		Email:           info.Email,
		Password:        info.Password,
		PasswordConfirm: info.PasswordConfirm,
		Role:            desc.Role_name[int32(info.Role)],
	}
}

func ToUpdateUserInfoFromDesc(info *desc.UpdateUserInfo) *model.UpdateUserInfo {
	m := &model.UpdateUserInfo{
		Name:  info.GetName().Value,
		Email: info.GetEmail().Value,
		Role:  model.UnknownRole,
	}
	if info.Role.Enum() != nil {
		m.Role = desc.Role_name[int32(*info.Role.Enum())]
	}
	return m
}
