package user

import (
	"github.com/VadimGossip/concoleChat-auth/internal/service"
	desc "github.com/VadimGossip/concoleChat-auth/pkg/user_v1"
)

type Implementation struct {
	desc.UnimplementedUserV1Server
	userService service.UserService
}

func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
