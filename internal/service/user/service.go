package user

import (
	db "github.com/VadimGossip/platform_common/pkg/db/postgres"

	"github.com/VadimGossip/concoleChat-auth/internal/repository"
	def "github.com/VadimGossip/concoleChat-auth/internal/service"
)

var _ def.UserService = (*service)(nil)

type service struct {
	userRepository   repository.UserRepository
	passwordService  def.PasswordService
	userCacheService def.UserCacheService
	auditService     def.AuditService
	txManager        db.TxManager
}

func NewService(userRepository repository.UserRepository,
	passwordService def.PasswordService,
	userCacheService def.UserCacheService,
	auditService def.AuditService,
	txManager db.TxManager) *service {
	return &service{
		userRepository:   userRepository,
		passwordService:  passwordService,
		userCacheService: userCacheService,
		auditService:     auditService,
		txManager:        txManager,
	}
}
