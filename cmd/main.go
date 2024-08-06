package main

import (
	"context"
	"time"

	"github.com/VadimGossip/concoleChat-auth/internal/app"
	"github.com/sirupsen/logrus"
)

var configDir = "config"
var appName = "Console Chat Auth"

func main() {
	ctx := context.Background()
	a, err := app.NewApp(ctx, appName, configDir, time.Now())
	if err != nil {
		logrus.Fatalf("failed to init app[%s]: %s", appName, err)
	}

	if err = a.Run(); err != nil {
		logrus.Infof("Application run process finished with error: %s", err)
	}
}
