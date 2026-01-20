// Package server 提供 gRPC 服务器实现
package server

import (
	"context"
	"fmt"

	pb "github.com/uyou/uyou-go-api-starter/api/proto/user"
	"github.com/uyou/uyou-go-api-starter/internal/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServiceServer 用户服务 gRPC 服务器实现
type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	userService user.Service
}

// NewUserServiceServer 创建用户服务 gRPC 服务器
func NewUserServiceServer(userService user.Service) *UserServiceServer {
	return &UserServiceServer{
		userService: userService,
	}
}

// GetUser 获取用户信息
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// 调用用户服务
	usr, err := s.userService.GetUserByID(ctx, uint(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	// 转换为 protobuf 消息
	return &pb.GetUserResponse{
		User: convertUserToProto(usr),
	}, nil
}

// GetUserByEmail 根据邮箱获取用户信息
func (s *UserServiceServer) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserResponse, error) {
	// 调用用户服务
	usr, err := s.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	// 转换为 protobuf 消息
	return &pb.GetUserResponse{
		User: convertUserToProto(usr),
	}, nil
}

// ListUsers 获取用户列表
func (s *UserServiceServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// 设置默认分页参数
	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 调用用户服务
	users, total, err := s.userService.ListUsers(ctx, page, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	// 转换为 protobuf 消息
	pbUsers := make([]*pb.User, len(users))
	for i, usr := range users {
		pbUsers[i] = convertUserToProto(&usr)
	}

	return &pb.ListUsersResponse{
		Users:    pbUsers,
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}

// UpdateUser 更新用户信息
func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// 构建更新请求
	updateReq := user.UpdateUserRequest{
		Name:  req.Name,
		Email: req.Email,
	}

	// 调用用户服务
	usr, err := s.userService.UpdateUser(ctx, uint(req.Id), updateReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	// 转换为 protobuf 消息
	return &pb.UpdateUserResponse{
		User: convertUserToProto(usr),
	}, nil
}

// DeleteUser 删除用户
func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// 调用用户服务
	err := s.userService.DeleteUser(ctx, uint(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &pb.DeleteUserResponse{
		Success: true,
		Message: fmt.Sprintf("User %d deleted successfully", req.Id),
	}, nil
}

// convertUserToProto 将用户模型转换为 protobuf 消息
func convertUserToProto(usr *user.User) *pb.User {
	return &pb.User{
		Id:        uint32(usr.ID),
		Name:      usr.Name,
		Email:     usr.Email,
		Roles:     usr.GetRoleNames(),
		CreatedAt: usr.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: usr.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
