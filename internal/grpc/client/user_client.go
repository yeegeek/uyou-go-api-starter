// Package client 提供 gRPC 客户端实现
package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/uyou/uyou-go-api-starter/api/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UserClient 用户服务 gRPC 客户端
type UserClient struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
}

// NewUserClient 创建用户服务 gRPC 客户端
func NewUserClient(address string, timeout int) (*UserClient, error) {
	// 创建连接选项
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Duration(timeout) * time.Second),
	}

	// 连接到 gRPC 服务器
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	// 创建客户端
	client := pb.NewUserServiceClient(conn)

	return &UserClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close 关闭连接
func (c *UserClient) Close() error {
	return c.conn.Close()
}

// GetUser 获取用户信息
func (c *UserClient) GetUser(ctx context.Context, userID uint32) (*pb.User, error) {
	resp, err := c.client.GetUser(ctx, &pb.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}

// GetUserByEmail 根据邮箱获取用户信息
func (c *UserClient) GetUserByEmail(ctx context.Context, email string) (*pb.User, error) {
	resp, err := c.client.GetUserByEmail(ctx, &pb.GetUserByEmailRequest{
		Email: email,
	})
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}

// ListUsers 获取用户列表
func (c *UserClient) ListUsers(ctx context.Context, page, pageSize int32) ([]*pb.User, int32, error) {
	resp, err := c.client.ListUsers(ctx, &pb.ListUsersRequest{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, 0, err
	}
	return resp.Users, resp.Total, nil
}

// UpdateUser 更新用户信息
func (c *UserClient) UpdateUser(ctx context.Context, userID uint32, name, email string) (*pb.User, error) {
	resp, err := c.client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:    userID,
		Name:  name,
		Email: email,
	})
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}

// DeleteUser 删除用户
func (c *UserClient) DeleteUser(ctx context.Context, userID uint32) (bool, string, error) {
	resp, err := c.client.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: userID,
	})
	if err != nil {
		return false, "", err
	}
	return resp.Success, resp.Message, nil
}
