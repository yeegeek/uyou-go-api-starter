// Package redis 提供 Redis 连接和操作的封装
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/uyou/uyou-go-api-starter/internal/config"
)

// Client 封装 Redis 客户端
type Client struct {
	*redis.Client
}

// NewClient 创建新的 Redis 客户端连接
func NewClient(cfg *config.Config) (*Client, error) {
	// 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		DialTimeout:  time.Duration(cfg.Redis.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Redis.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Redis.WriteTimeout) * time.Second,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Client{Client: rdb}, nil
}

// Close 关闭 Redis 连接
func (c *Client) Close() error {
	return c.Client.Close()
}

// HealthCheck 检查 Redis 连接健康状态
func (c *Client) HealthCheck(ctx context.Context) error {
	return c.Ping(ctx).Err()
}

// SetWithExpiration 设置键值对并指定过期时间
func (c *Client) SetWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Set(ctx, key, value, expiration).Err()
}

// GetString 获取字符串值
func (c *Client) GetString(ctx context.Context, key string) (string, error) {
	return c.Get(ctx, key).Result()
}

// Delete 删除键
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	return c.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (c *Client) Exists(ctx context.Context, keys ...string) (bool, error) {
	result, err := c.Client.Exists(ctx, keys...).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Increment 原子递增
func (c *Client) Increment(ctx context.Context, key string) (int64, error) {
	return c.Incr(ctx, key).Result()
}

// IncrementBy 原子递增指定值
func (c *Client) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.IncrBy(ctx, key, value).Result()
}

// SetNX 仅当键不存在时设置值（分布式锁）
func (c *Client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return c.Client.SetNX(ctx, key, value, expiration).Result()
}

// Expire 设置键的过期时间
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.Client.Expire(ctx, key, expiration).Err()
}
