package initializers

import (
	"fmt"
	"github.com/SiriusServiceDesk/auth-service/internal/app/dependencies"
	"github.com/SiriusServiceDesk/auth-service/internal/config"
	"github.com/SiriusServiceDesk/auth-service/internal/grpc/server"
	"github.com/SiriusServiceDesk/auth-service/internal/grpc/server/handlers"
	"github.com/SiriusServiceDesk/auth-service/pkg/logger"
	"github.com/SiriusServiceDesk/gateway-service/pkg/auth_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func InitializeGRPCListener() (net.Listener, error) {
	cfg := config.GetConfig().GrpcServer
	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	logger.Info("Initializing GRPCListener", zap.String("address", address))
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatal("Failed to listen grpc server", zap.Error(err))
	}
	return listener, nil
}

func InitializeGRPCServer(listener net.Listener, container *dependencies.Container) *server.GRPCServer {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	auth_v1.RegisterAuthV1Server(grpcServer, handlers.NewHandler(container.UserService, container.Redis))

	logger.Info("Complete register grpc handlers")
	return server.NewGRPCServer(listener, grpcServer)
}

func StartGRPCServer(grpcServer *server.GRPCServer) {
	go func() {
		err := grpcServer.Start()
		if err != nil {
			logger.Fatal("failed starting GRPC server", "error", err)
		}
	}()
}
