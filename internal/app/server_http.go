package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	desc "github.com/VadimGossip/concoleChat-auth/pkg/user_v1"
)

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := desc.RegisterUserV1HandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", a.cfg.AppGrpcServer.Port), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.cfg.AppHttpServer.Port),
		Handler: corsMiddleware.Handler(mux),
	}

	return nil
}

func (a *App) runHTTPServer() error {
	logrus.Infof("[%s] HTTP server is running on: %s", a.name, a.httpServer.Addr)
	return a.httpServer.ListenAndServe()
}
