package access

import (
	"context"
	"fmt"
	desc "github.com/VadimGossip/concoleChat-auth/pkg/access_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

const authPrefix = "Bearer "

func (s *Implementation) Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, fmt.Errorf("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	err := s.accessService.Check(ctx, strings.TrimPrefix(authHeader[0], authPrefix), req.EndpointAddress)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
