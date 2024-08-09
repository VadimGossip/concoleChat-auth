package user

import (
	"github.com/VadimGossip/concoleChat-auth/internal/model"
	"github.com/VadimGossip/concoleChat-auth/internal/repository"
	def "github.com/VadimGossip/concoleChat-auth/internal/service"
)

var _ def.UserCacheService = (*service)(nil)

type service struct {
	userCacheCfg        model.UserCacheCfg
	userCacheRepository repository.UserCacheRepository
}

func NewService(userCacheCfg model.UserCacheCfg, userCacheRepository repository.UserCacheRepository) *service {
	return &service{
		userCacheCfg:        userCacheCfg,
		userCacheRepository: userCacheRepository,
	}
}
