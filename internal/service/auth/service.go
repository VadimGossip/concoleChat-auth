package auth

import (
	"github.com/VadimGossip/concoleChat-auth/internal/config"
	def "github.com/VadimGossip/concoleChat-auth/internal/service"
)

var _ def.AuthService = (*service)(nil)

type service struct {
	tokenConfig     config.TokenConfig
	userService     def.UserService
	passwordService def.PasswordService
	tokenService    def.TokenService
}

func NewService(tokenConfig config.TokenConfig,
	userService def.UserService,
	passwordService def.PasswordService,
	tokenService def.TokenService) *service {
	return &service{
		tokenConfig:     tokenConfig,
		userService:     userService,
		passwordService: passwordService,
		tokenService:    tokenService,
	}
}
