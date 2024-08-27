package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

type GRPCInterceptor interface {
	Hook() grpc.UnaryServerInterceptor
}

type interceptor struct{}

func NewInterceptor() *interceptor {
	return &interceptor{}
}

type validator interface {
	Validate() error
}

func (*interceptor) Hook() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if val, ok := req.(validator); ok {
			if err := val.Validate(); err != nil {
				return nil, err
			}
		}
		return handler(ctx, req)
	}
}
