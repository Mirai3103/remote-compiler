package main

import (
	"fmt"
	"net"

	grpcIpml "github.com/Mirai3103/remote-compiler/internal/handler/grpc"
	"github.com/Mirai3103/remote-compiler/pkg/config"
	"github.com/Mirai3103/remote-compiler/pkg/logger"
	"github.com/Mirai3103/remote-compiler/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	log := logger.GetLogger()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.GRPC.Port))
	if err != nil {
		log.Fatal("Failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	executionHandler := grpcIpml.NewExecutionHandler(cfg)
	proto.RegisterExecutionServiceServer(grpcServer, executionHandler)
	reflection.Register(grpcServer)
	log.Info("Starting gRPC server", zap.Int("port", cfg.GRPC.Port))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to serve gRPC", zap.Error(err))
	}
}
