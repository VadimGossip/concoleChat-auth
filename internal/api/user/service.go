package user

import (
	"github.com/VadimGossip/concoleChat-auth/internal/service"
	asyncService "github.com/VadimGossip/concoleChat-auth/internal/service/producer"
	desc "github.com/VadimGossip/concoleChat-auth/pkg/user_v1"
)

type Implementation struct {
	desc.UnimplementedUserV1Server
	userService      service.UserService
	userAsyncService asyncService.UserProducerService
}

func NewImplementation(userService service.UserService,
	userAsyncService asyncService.UserProducerService) *Implementation {
	return &Implementation{
		userService:      userService,
		userAsyncService: userAsyncService,
	}
}
