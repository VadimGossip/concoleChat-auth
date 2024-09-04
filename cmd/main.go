package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/VadimGossip/concoleChat-auth/internal/app"
	"github.com/VadimGossip/concoleChat-auth/internal/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var appName = "Console Chat Auth"
var logLevel = flag.String("l", "info", "log level")

func main() {
	flag.Parse()
	logger.Init(getCore(getAtomicLevel(*logLevel)))

	ctx := context.Background()
	a, err := app.NewApp(ctx, appName, time.Now())
	if err != nil {
		logger.Fatalf("failed to init app%s: %s", appName, err)
	}

	if err = a.Run(ctx); err != nil {
		logger.Infof("app%s run process finished with error: %s", appName, err)
	}
}

func getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
	)
}

func getAtomicLevel(loglevel string) zap.AtomicLevel {
	var level zapcore.Level
	if err := level.Set(loglevel); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}

	return zap.NewAtomicLevelAt(level)
}
