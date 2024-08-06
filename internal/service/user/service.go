package user

import (
	"github.com/VadimGossip/concoleChat-auth/internal/client/db"
	"github.com/VadimGossip/concoleChat-auth/internal/repository"
	def "github.com/VadimGossip/concoleChat-auth/internal/service"
)

var _ def.UserService = (*service)(nil)

type service struct {
	userRepository repository.UserRepository
	auditService   def.AuditService
	txManager      db.TxManager
}

func NewService(userRepository repository.UserRepository, auditService def.AuditService, txManager db.TxManager) *service {
	return &service{
		userRepository: userRepository,
		auditService:   auditService,
		txManager:      txManager,
	}
}
