package user

import (
	"github.com/VadimGossip/concoleChat-auth/internal/config"
	"github.com/VadimGossip/concoleChat-auth/internal/repository"
	def "github.com/VadimGossip/concoleChat-auth/internal/service"
)

var _ def.UserCacheService = (*service)(nil)

type service struct {
	userCacheConfig     config.UserCacheConfig
	userCacheRepository repository.UserCacheRepository
}

func NewService(userCacheConfig config.UserCacheConfig, userCacheRepository repository.UserCacheRepository) *service {
	return &service{
		userCacheConfig:     userCacheConfig,
		userCacheRepository: userCacheRepository,
	}
}
