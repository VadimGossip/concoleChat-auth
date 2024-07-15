package app

import (
	"google.golang.org/grpc"

	"github.com/VadimGossip/concoleChat-auth/pkg/user_v1"
)

func initGrpcRouter(app *App) func(*grpc.Server) {
	return func(s *grpc.Server) {
		user_v1.RegisterUserV1Server(s, app.userImpl)
	}
}
