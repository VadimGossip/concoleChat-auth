package validator

import (
	"fmt"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

func passwordMatch(password, passwordConfirm string) error {
	if password != passwordConfirm {
		return fmt.Errorf("password and password confirm not match")
	}
	return nil
}

func CreateValidation(info *model.UserInfo) error {
	return passwordMatch(info.Password, info.PasswordConfirm)
}

func emptyUpdate(updateInfo *model.UpdateUserInfo) error {
	if updateInfo.Name == nil && updateInfo.Email == nil && updateInfo.Role == model.UnknownRole {
		return fmt.Errorf("update with no updated fields not allowed")
	}
	return nil
}

func UpdateValidation(updateInfo *model.UpdateUserInfo) error {
	return emptyUpdate(updateInfo)
}
