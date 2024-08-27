package token

import (
	def "github.com/VadimGossip/concoleChat-auth/internal/service"
)

var _ def.TokenService = (*service)(nil)

type service struct {
}

func NewService() *service {
	return &service{}
}
