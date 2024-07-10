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

func emptyName(name string) error {
	if name == "" {
		return fmt.Errorf("empty name not allowed")
	}
	return nil
}

func emptyEmail(email string) error {
	if email == "" {
		return fmt.Errorf("empty email not allowed")
	}
	return nil
}

func CreateValidation(info *model.UserInfo) error {
	if err := emptyName(info.Name); err != nil {
		return err
	}
	if err := emptyEmail(info.Email); err != nil {
		return err
	}
	if err := passwordMatch(info.Password, info.PasswordConfirm); err != nil {
		return err
	}
	return nil
}

func UpdateValidation(updateInfo *model.UpdateUserInfo) error {
	if err := emptyName(updateInfo.Name); err != nil {
		return err
	}
	if err := emptyEmail(updateInfo.Email); err != nil {
		return err
	}
	return nil
}
