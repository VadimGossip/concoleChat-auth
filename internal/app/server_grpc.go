package app

import (
	"context"
	"net"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/VadimGossip/concoleChat-auth/internal/interceptor"
	"github.com/VadimGossip/concoleChat-auth/internal/logger"
	descAccess "github.com/VadimGossip/concoleChat-auth/pkg/access_v1"
	descAuth "github.com/VadimGossip/concoleChat-auth/pkg/auth_v1"
	descUser "github.com/VadimGossip/concoleChat-auth/pkg/user_v1"
)

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			interceptor.LogInterceptor,
			interceptor.ValidateInterceptor,
		)),
	)

	reflection.Register(a.grpcServer)
	descAccess.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AccessImpl(ctx))
	descAuth.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthImpl(ctx))
	descUser.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserImpl(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}
	logger.Infof("%s GRPC server is running on: %s", a.name, a.serviceProvider.GRPCConfig().Address())

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}
