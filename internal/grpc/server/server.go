// Package server 提供 gRPC 服务器管理
package server

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/yeegeek/uyou-go-api-starter/internal/config"
	"github.com/yeegeek/uyou-go-api-starter/internal/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server gRPC 服务器
type Server struct {
	grpcServer  *grpc.Server
	config      *config.Config
	userService user.Service
}

// NewServer 创建 gRPC 服务器
func NewServer(cfg *config.Config, userService user.Service) *Server {
	// 创建 gRPC 服务器选项
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(cfg.GRPC.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(cfg.GRPC.MaxSendMsgSize),
	}

	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer(opts...)

	// 注册用户服务
	// 注意：由于 protobuf 代码是占位符，RegisterUserServiceServer 可能不存在
	// 在实际使用前需要运行 protoc 生成真实的 protobuf 代码
	// TODO: 传入 userRepository 并取消注释以下代码
	// userServer := NewUserServiceServer(userService, userRepo)
	// pb.RegisterUserServiceServer(grpcServer, userServer)

	// 注册反射服务（用于 grpcurl 等工具）
	reflection.Register(grpcServer)

	return &Server{
		grpcServer:  grpcServer,
		config:      cfg,
		userService: userService,
	}
}

// Start 启动 gRPC 服务器
func (s *Server) Start() error {
	// 监听端口
	lis, err := net.Listen("tcp", ":"+s.config.GRPC.Port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", s.config.GRPC.Port, err)
	}

	slog.Info("gRPC server starting", "port", s.config.GRPC.Port)

	// 启动服务器
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

// Stop 停止 gRPC 服务器
func (s *Server) Stop() {
	slog.Info("Stopping gRPC server...")
	s.grpcServer.GracefulStop()
}
