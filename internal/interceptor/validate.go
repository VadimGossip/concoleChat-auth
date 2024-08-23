package interceptor

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type validator interface {
	Validate() error
}

func BuildInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logrus.Infof("Intercepted request %+v method %+v", req, info.FullMethod)
		if val, ok := req.(validator); ok {
			if err := val.Validate(); err != nil {
				return nil, err
			}
		}
		return handler(ctx, req)
	}
}
