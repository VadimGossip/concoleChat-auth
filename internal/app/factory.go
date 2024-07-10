package app

import (
	"github.com/VadimGossip/concoleChat-auth/internal/api/user"
	"github.com/VadimGossip/concoleChat-auth/internal/service"
	userService "github.com/VadimGossip/concoleChat-auth/internal/service/user"
)

type Factory struct {
	dbAdapter *DBAdapter

	userService service.UserService

	userImpl *user.Implementation
}

var factory *Factory

func newFactory(dbAdapter *DBAdapter) *Factory {
	factory = &Factory{dbAdapter: dbAdapter}
	factory.userService = userService.NewService(dbAdapter.userRepo)
	factory.userImpl = user.NewImplementation(factory.userService)
	return factory
}
