package app

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/VadimGossip/concoleChat-auth/internal/interceptor"
	desc "github.com/VadimGossip/concoleChat-auth/pkg/user_v1"
)

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)

	reflection.Register(a.grpcServer)

	desc.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserImpl(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	list, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.AppGrpcServer.Port))
	if err != nil {
		return err
	}
	logrus.Infof("[%s] GRPC server is running on: %d", a.name, a.cfg.AppGrpcServer.Port)

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}
