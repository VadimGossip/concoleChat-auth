package main

import (
	"time"

	"github.com/VadimGossip/concoleChat-auth/internal/app"
	"github.com/sirupsen/logrus"
)

var configDir = "config"

func main() {
	auth := app.NewApp("Console Chat Auth", configDir, time.Now())
	if err := auth.Run(); err != nil {
		logrus.Infof("Application run process finished with error")
	}
}
