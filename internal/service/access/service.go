package access

import (
	"github.com/VadimGossip/concoleChat-auth/internal/config"
	"github.com/VadimGossip/concoleChat-auth/internal/repository"

	def "github.com/VadimGossip/concoleChat-auth/internal/service"
)

var _ def.AccessService = (*service)(nil)

type service struct {
	tokenConfig      config.TokenConfig
	accessRepository repository.AccessRepository
	userService      def.UserService
	passwordService  def.PasswordService
	tokenService     def.TokenService
}

func NewService(tokenConfig config.TokenConfig,
	accessRepository repository.AccessRepository,
	userService def.UserService,
	passwordService def.PasswordService,
	tokenService def.TokenService) *service {
	return &service{
		tokenConfig:      tokenConfig,
		accessRepository: accessRepository,
		userService:      userService,
		passwordService:  passwordService,
		tokenService:     tokenService,
	}
}
