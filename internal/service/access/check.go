package access

import (
	"context"
	"fmt"
)

func (s *service) Check(ctx context.Context, accessToken, endpointAddress string) error {
	claims, err := s.tokenService.Verify(accessToken, []byte(s.tokenConfig.AccessTokenSecretKey()))
	if err != nil {
		return fmt.Errorf("access token is invalid")
	}
	accessible, err := s.accessRepository.AccessibleByRole(ctx, claims.Role, endpointAddress)
	if err != nil {
		return err
	}
	if !accessible {
		return fmt.Errorf("access denied")
	}
	return nil
}
