package app

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcServer struct {
	port     int
	server   *grpc.Server
	listener net.Listener
}

func NewGrpcServer(port int) *GrpcServer {
	return &GrpcServer{
		port: port,
	}
}

func (s *GrpcServer) Listen(grpcRouter func(*grpc.Server)) error {
	s.server = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	grpcRouter(s.server)

	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	logrus.Infof("[grpc/server] Starting on port: %d", s.port)
	return s.server.Serve(s.listener)
}
