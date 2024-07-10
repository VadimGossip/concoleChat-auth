package main

import (
	"time"

	"github.com/VadimGossip/concoleChat-auth/internal/app"
)

var configDir = "config"

func main() {
	auth := app.NewApp("Console Chat Auth", configDir, time.Now())
	auth.Run()
}
