package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/VadimGossip/platform_common/pkg/closer"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/VadimGossip/concoleChat-auth/internal/config"
	"github.com/VadimGossip/concoleChat-auth/internal/model"
	//import for init
	_ "github.com/VadimGossip/concoleChat-auth/statik"
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
	httpServer      *http.Server
	swaggerServer   *http.Server
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
		a.initHTTPServer,
		a.initSwaggerServer,
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

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			logrus.Fatalf("failed to run GRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			logrus.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			logrus.Fatalf("failed to run Swagger server: %v", err)
		}
	}()

	wg.Wait()

	return nil
}
