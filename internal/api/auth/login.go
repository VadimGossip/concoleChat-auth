package auth

import (
	"context"
	"github.com/VadimGossip/concoleChat-auth/internal/converter"
	desc "github.com/VadimGossip/concoleChat-auth/pkg/auth_v1"
)

func (s *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	refreshToken, err := s.authService.Login(ctx, converter.ToLoginUserInfoFromDesc(req))
	if err != nil {
		return nil, err
	}
	return &desc.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}
