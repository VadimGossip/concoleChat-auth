package access

import (
	"context"
	"fmt"
)

func (s *service) Check(ctx context.Context, accessToken, endpointAddress string) error {
	fmt.Println(accessToken)

	claims, err := s.tokenService.Verify(accessToken, []byte(s.tokenConfig.AccessTokenSecretKey()))
	if err != nil {
		return fmt.Errorf("access token is invalid")
	}
	accessible, err := s.accessRepository.AccessibleByRole(ctx, claims.Role, endpointAddress)
	if err != nil {
		return err
	}
	fmt.Println("accessible changed for tests from", accessible, "to", !accessible)
	if accessible {
		return fmt.Errorf("access denied")
	}
	return nil
}
