package auth

//AuthServiceConfig

import (
	"github.com/VadimGossip/concoleChat-auth/internal/config"
	def "github.com/VadimGossip/concoleChat-auth/internal/service"
)

var _ def.AuthService = (*service)(nil)

type service struct {
	authServiceConfig config.AuthServiceConfig
	userService       def.UserService
	passwordService   def.PasswordService
	tokenService      def.TokenService
}

func NewService(authServiceConfig config.AuthServiceConfig,
	userService def.UserService,
	passwordService def.PasswordService,
	tokenService def.TokenService) *service {
	return &service{
		userService:     userService,
		passwordService: passwordService,
		tokenService:    tokenService,
	}
}
