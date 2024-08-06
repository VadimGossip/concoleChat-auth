package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/VadimGossip/concoleChat-auth/internal/closer"
	"github.com/VadimGossip/concoleChat-auth/internal/config"
	"github.com/VadimGossip/concoleChat-auth/internal/model"
	desc "github.com/VadimGossip/concoleChat-auth/pkg/user_v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

type App struct {
	serviceProvider *serviceProvider
	name            string
	configDir       string
	appStartedAt    time.Time
	cfg             *model.Config
	grpcServer      *grpc.Server
}

func NewApp(ctx context.Context, name, configDir string, appStartedAt time.Time) (*App, error) {
	a := &App{
		name:         name,
		configDir:    configDir,
		appStartedAt: appStartedAt,
	}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	cfg, err := config.Init(a.configDir)
	if err != nil {
		return fmt.Errorf("[%s] config initialization error: %s", a.name, err)
	}
	a.cfg = cfg
	logrus.Infof("[%s] got config: [%+v]", a.name, *a.cfg)
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider(a.cfg)
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	desc.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserImpl(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	logrus.Infof("[grpc/server] Starting on port: %d", a.cfg.AppGrpcServer.Port)

	list, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.AppGrpcServer.Port))
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGRPCServer()
}
