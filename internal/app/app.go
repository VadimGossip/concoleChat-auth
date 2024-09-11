package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/VadimGossip/concoleChat-auth/internal/logger"
	"github.com/VadimGossip/concoleChat-auth/internal/metric"
	"github.com/VadimGossip/platform_common/pkg/closer"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	//import for init
	_ "github.com/VadimGossip/concoleChat-auth/statik"
)

type App struct {
	serviceProvider  *serviceProvider
	name             string
	appStartedAt     time.Time
	grpcServer       *grpc.Server
	httpServer       *http.Server
	swaggerServer    *http.Server
	prometheusServer *http.Server
}

func NewApp(ctx context.Context, name string, appStartedAt time.Time) (*App, error) {
	a := &App{
		name:         name,
		appStartedAt: appStartedAt,
	}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initConfig(_ context.Context) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		metric.Init,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
		a.initPrometheusServer,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()
	ctx, cancel := context.WithCancel(ctx)

	wg := &sync.WaitGroup{}
	wg.Add(5)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			logger.Fatalf("%s failed to run GRPC server: %v", a.name, err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			logger.Fatalf("%s failed to run HTTP server: %v", a.name, err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			logger.Fatalf("%s failed to run Swagger server: %v", a.name, err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runPrometheusServer()
		if err != nil {
			logger.Fatalf("%s failed to run Prometheus server: %v", a.name, err)
		}
	}()

	go func() {
		defer wg.Done()
		err := a.serviceProvider.UserConsumerService(ctx).RunConsumer(ctx)
		if err != nil {
			logger.Fatalf("%s failed to run consumer: %s", a.name, err)
		}
	}()

	gracefulShutdown(ctx, cancel, wg)
	return nil
}

func gracefulShutdown(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	select {
	case <-ctx.Done():
		logger.Info("terminating: context cancelled")
	case c := <-waitSignal():
		logger.Infof("terminating: got signal: %s", c)
	}

	cancel()
	if wg != nil {
		wg.Wait()
	}
}

func waitSignal() chan os.Signal {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	return sigterm
}
