package server

import (
	"google.golang.org/grpc"
	"net"
)

type GRPCServer struct {
	listener   net.Listener
	grpcServer *grpc.Server
}

func NewGRPCServer(listener net.Listener, server *grpc.Server) *GRPCServer {
	return &GRPCServer{
		listener:   listener,
		grpcServer: server,
	}
}

func (s *GRPCServer) Start() error {
	return s.grpcServer.Serve(s.listener)
}

func (s *GRPCServer) Shutdown() {
	s.grpcServer.GracefulStop()
}
