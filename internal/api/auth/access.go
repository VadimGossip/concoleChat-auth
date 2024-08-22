package auth

import (
	"context"

	desc "github.com/VadimGossip/concoleChat-auth/pkg/auth_v1"
)

func (s *Implementation) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	accessToken, err := s.authService.AccessToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &desc.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}
