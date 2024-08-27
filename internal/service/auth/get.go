package auth

import (
	"context"
	"fmt"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

func (s *service) verifyUserInfo(ctx context.Context, refreshToken string) (model.UserInfo, error) {
	claims, err := s.tokenService.Verify(refreshToken, []byte(s.tokenConfig.RefreshTokenSecretKey()))
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("invalid refresh token")
	}

	user, err := s.userService.GetByUsername(ctx, claims.Name)
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("user with username %s not found", claims.Name)
	}
	return user.Info, nil
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	info, err := s.verifyUserInfo(ctx, refreshToken)
	if err != nil {
		return "", err
	}

	rt, err := s.tokenService.Generate(info, []byte(s.tokenConfig.RefreshTokenSecretKey()), s.tokenConfig.RefreshTokenExpiration())
	if err != nil {
		return "", err
	}
	return rt, nil
}

func (s *service) AccessToken(ctx context.Context, refreshToken string) (string, error) {
	info, err := s.verifyUserInfo(ctx, refreshToken)
	if err != nil {
		return "", err
	}
	rt, err := s.tokenService.Generate(info, []byte(s.tokenConfig.AccessTokenSecretKey()), s.tokenConfig.AccessTokenExpiration())
	if err != nil {
		return "", err
	}
	return rt, nil
}
