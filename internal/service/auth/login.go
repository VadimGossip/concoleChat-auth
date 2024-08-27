package auth

import (
	"context"
	"fmt"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

func (s *service) Login(ctx context.Context, info *model.LoginUserInfo) (string, error) {
	user, err := s.userService.GetByUsername(ctx, info.Name)
	if err != nil {
		return "", fmt.Errorf("user with username %s not found", info.Name)
	}

	if !s.passwordService.Verify(user.Info.Password, info.Password) {
		return "", fmt.Errorf("password is incorrect")
	}

	refreshToken, err := s.tokenService.Generate(user.Info, []byte(s.tokenConfig.RefreshTokenSecretKey()), s.tokenConfig.RefreshTokenExpiration())
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token")
	}
	return refreshToken, nil
}
