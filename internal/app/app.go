package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VadimGossip/concoleChat-auth/internal/config"
	"github.com/VadimGossip/concoleChat-auth/internal/model"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

type App struct {
	*Factory
	name         string
	configDir    string
	appStartedAt time.Time
	cfg          *model.Config
	grpcServer   *GrpcServer
}

func NewApp(name, configDir string, appStartedAt time.Time) *App {
	return &App{
		name:         name,
		configDir:    configDir,
		appStartedAt: appStartedAt,
	}
}

func (app *App) Run() {
	ctx := context.Background()
	cfg, err := config.Init(app.configDir)
	if err != nil {
		logrus.Fatalf("Config initialization error: %s", err)
	}
	app.cfg = cfg
	logrus.Infof("[%s] got config: [%+v]", app.name, *app.cfg)

	dbAdapter := NewDBAdapter(cfg.Db)
	if err = dbAdapter.Connect(ctx); err != nil {
		logrus.Fatalf("Fail to connect db %s", err)
	}
	app.Factory = newFactory(dbAdapter)

	go func() {
		app.grpcServer = NewGrpcServer(cfg.AppGrpcServer.Port)
		grpcRouter := initGrpcRouter(app)
		if err = app.grpcServer.Listen(grpcRouter); err != nil {
			logrus.Fatalf("Failed to start GRPC server %s", err)
		}
	}()

	logrus.Infof("[%s] started", app.name)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	logrus.Infof("[%s] got signal: [%s]", app.name, <-c)

	if err = dbAdapter.Disconnect(ctx); err != nil {
		logrus.Errorf("Fail to diconnect db %s", err)
	}

	logrus.Infof("[%s] stopped", app.name)
}
