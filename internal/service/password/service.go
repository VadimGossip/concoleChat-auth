package password

import (
	def "github.com/VadimGossip/concoleChat-auth/internal/service"
)

var _ def.PasswordService = (*service)(nil)

type service struct {
}

func NewService() *service {
	return &service{}
}
