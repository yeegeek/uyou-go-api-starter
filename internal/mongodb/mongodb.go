// Package mongodb 提供 MongoDB 连接和操作的封装
package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"github.com/uyou/uyou-go-api-starter/internal/config"
)

// Client 封装 MongoDB 客户端
type Client struct {
	*mongo.Client
	database *mongo.Database
}

// NewClient 创建新的 MongoDB 客户端连接
func NewClient(cfg *config.Config) (*Client, error) {
	// 创建客户端选项
	clientOptions := options.Client().
		ApplyURI(cfg.MongoDB.URI).
		SetMaxPoolSize(uint64(cfg.MongoDB.MaxPoolSize)).
		SetMinPoolSize(uint64(cfg.MongoDB.MinPoolSize)).
		SetConnectTimeout(time.Duration(cfg.MongoDB.ConnectTimeout) * time.Second)

	// 连接到 MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.MongoDB.ConnectTimeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	// 测试连接
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	// 获取数据库实例
	database := client.Database(cfg.MongoDB.Database)

	return &Client{
		Client:   client,
		database: database,
	}, nil
}

// Database 返回数据库实例
func (c *Client) Database() *mongo.Database {
	return c.database
}

// Collection 返回指定集合
func (c *Client) Collection(name string) *mongo.Collection {
	return c.database.Collection(name)
}

// Close 关闭 MongoDB 连接
func (c *Client) Close(ctx context.Context) error {
	return c.Client.Disconnect(ctx)
}

// HealthCheck 检查 MongoDB 连接健康状态
func (c *Client) HealthCheck(ctx context.Context) error {
	return c.Ping(ctx, readpref.Primary())
}

// WithTransaction 执行事务
func (c *Client) WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) error) error {
	session, err := c.Client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx)
	})
	return err
}
