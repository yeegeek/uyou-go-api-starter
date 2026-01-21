package redis

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeegeek/uyou-go-api-starter/internal/config"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
		errMsg  string
	}{
		{
			// 注意：本测试在“Redis 未运行”的机器上用于验证失败分支；
			// 如果本机恰好有 Redis 运行，则会返回成功，这种情况下我们跳过该用例。
			name: "valid config but no redis server (expected failure if redis is down)",
			cfg: &config.Config{
				Redis: config.RedisConfig{
					Host:         "localhost",
					Port:         6379,
					Password:     "",
					DB:           0,
					DialTimeout:  5,
					ReadTimeout:  3,
					WriteTimeout: 3,
					PoolSize:     10,
					MinIdleConns: 5,
				},
			},
			wantErr: true,
			errMsg:  "failed to connect to redis",
		},
		{
			name: "invalid host",
			cfg: &config.Config{
				Redis: config.RedisConfig{
					Host:         "invalid-host-that-does-not-exist",
					Port:         6379,
					Password:     "",
					DB:           0,
					DialTimeout:  5,
					ReadTimeout:  3,
					WriteTimeout: 3,
					PoolSize:     10,
					MinIdleConns: 5,
				},
			},
			wantErr: true,
			errMsg:  "failed to connect to redis",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.cfg)
			if tt.wantErr {
				if err == nil && client != nil {
					// 本机 Redis 可用时，失败分支不可复现，跳过即可
					_ = client.Close()
					t.Skip("Redis is available on localhost; skipping failure-path test case")
				}
				assert.Error(t, err)
				assert.Nil(t, client)
				if err != nil {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				if client != nil {
					defer client.Close()
				}
			}
		})
	}
}

func TestClient_HealthCheck(t *testing.T) {
	// 使用内存 Redis 客户端进行测试（如果可用）
	// 否则测试会失败，这是预期的
	cfg := &config.Config{
		Redis: config.RedisConfig{
			Host:         "localhost",
			Port:         6379,
			Password:     "",
			DB:           0,
			DialTimeout:  5,
			ReadTimeout:  3,
			WriteTimeout: 3,
			PoolSize:     10,
			MinIdleConns: 5,
		},
	}

	client, err := NewClient(cfg)
	if err != nil {
		// Redis 不可用时跳过测试
		t.Skipf("Redis not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	err = client.HealthCheck(ctx)
	assert.NoError(t, err)
}

func TestClient_SetWithExpiration(t *testing.T) {
	cfg := &config.Config{
		Redis: config.RedisConfig{
			Host:         "localhost",
			Port:         6379,
			Password:     "",
			DB:           0,
			DialTimeout:  5,
			ReadTimeout:  3,
			WriteTimeout: 3,
			PoolSize:     10,
			MinIdleConns: 5,
		},
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	key := "test:set:key"
	value := "test-value"
	expiration := 1 * time.Minute

	err = client.SetWithExpiration(ctx, key, value, expiration)
	assert.NoError(t, err)

	// 验证值
	result, err := client.GetString(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, value, result)

	// 清理
	_ = client.Delete(ctx, key)
}

func TestClient_GetString(t *testing.T) {
	cfg := &config.Config{
		Redis: config.RedisConfig{
			Host:         "localhost",
			Port:         6379,
			Password:     "",
			DB:           0,
			DialTimeout:  5,
			ReadTimeout:  3,
			WriteTimeout: 3,
			PoolSize:     10,
			MinIdleConns: 5,
		},
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	key := "test:get:key"
	value := "test-value"

	// 设置值
	err = client.SetWithExpiration(ctx, key, value, 1*time.Minute)
	require.NoError(t, err)

	// 获取值
	result, err := client.GetString(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, value, result)

	// 获取不存在的键
	_, err = client.GetString(ctx, "test:nonexistent")
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)

	// 清理
	_ = client.Delete(ctx, key)
}

func TestClient_Delete(t *testing.T) {
	cfg := &config.Config{
		Redis: config.RedisConfig{
			Host:         "localhost",
			Port:         6379,
			Password:     "",
			DB:           0,
			DialTimeout:  5,
			ReadTimeout:  3,
			WriteTimeout: 3,
			PoolSize:     10,
			MinIdleConns: 5,
		},
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	key := "test:delete:key"
	value := "test-value"

	// 设置值
	err = client.SetWithExpiration(ctx, key, value, 1*time.Minute)
	require.NoError(t, err)

	// 验证存在
	exists, err := client.Exists(ctx, key)
	assert.NoError(t, err)
	assert.True(t, exists)

	// 删除
	err = client.Delete(ctx, key)
	assert.NoError(t, err)

	// 验证不存在
	exists, err = client.Exists(ctx, key)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestClient_Exists(t *testing.T) {
	cfg := &config.Config{
		Redis: config.RedisConfig{
			Host:         "localhost",
			Port:         6379,
			Password:     "",
			DB:           0,
			DialTimeout:  5,
			ReadTimeout:  3,
			WriteTimeout: 3,
			PoolSize:     10,
			MinIdleConns: 5,
		},
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	key := "test:exists:key"
	value := "test-value"

	// 不存在的键
	exists, err := client.Exists(ctx, key)
	assert.NoError(t, err)
	assert.False(t, exists)

	// 设置值
	err = client.SetWithExpiration(ctx, key, value, 1*time.Minute)
	require.NoError(t, err)

	// 存在的键
	exists, err = client.Exists(ctx, key)
	assert.NoError(t, err)
	assert.True(t, exists)

	// 清理
	_ = client.Delete(ctx, key)
}

func TestClient_Increment(t *testing.T) {
	cfg := &config.Config{
		Redis: config.RedisConfig{
			Host:         "localhost",
			Port:         6379,
			Password:     "",
			DB:           0,
			DialTimeout:  5,
			ReadTimeout:  3,
			WriteTimeout: 3,
			PoolSize:     10,
			MinIdleConns: 5,
		},
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	key := "test:increment:key"

	// 第一次递增
	val, err := client.Increment(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), val)

	// 第二次递增
	val, err = client.Increment(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), val)

	// 清理
	_ = client.Delete(ctx, key)
}

func TestClient_IncrementBy(t *testing.T) {
	cfg := &config.Config{
		Redis: config.RedisConfig{
			Host:         "localhost",
			Port:         6379,
			Password:     "",
			DB:           0,
			DialTimeout:  5,
			ReadTimeout:  3,
			WriteTimeout: 3,
			PoolSize:     10,
			MinIdleConns: 5,
		},
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	key := "test:incrementby:key"

	// 递增 5
	val, err := client.IncrementBy(ctx, key, 5)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), val)

	// 再递增 3
	val, err = client.IncrementBy(ctx, key, 3)
	assert.NoError(t, err)
	assert.Equal(t, int64(8), val)

	// 清理
	_ = client.Delete(ctx, key)
}

func TestClient_SetNX(t *testing.T) {
	cfg := &config.Config{
		Redis: config.RedisConfig{
			Host:         "localhost",
			Port:         6379,
			Password:     "",
			DB:           0,
			DialTimeout:  5,
			ReadTimeout:  3,
			WriteTimeout: 3,
			PoolSize:     10,
			MinIdleConns: 5,
		},
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	key := "test:setnx:key"
	value := "test-value"
	expiration := 1 * time.Minute

	// 第一次设置应该成功
	ok, err := client.SetNX(ctx, key, value, expiration)
	assert.NoError(t, err)
	assert.True(t, ok)

	// 第二次设置应该失败（键已存在）
	ok, err = client.SetNX(ctx, key, "another-value", expiration)
	assert.NoError(t, err)
	assert.False(t, ok)

	// 清理
	_ = client.Delete(ctx, key)
}

func TestClient_Expire(t *testing.T) {
	cfg := &config.Config{
		Redis: config.RedisConfig{
			Host:         "localhost",
			Port:         6379,
			Password:     "",
			DB:           0,
			DialTimeout:  5,
			ReadTimeout:  3,
			WriteTimeout: 3,
			PoolSize:     10,
			MinIdleConns: 5,
		},
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	key := "test:expire:key"
	value := "test-value"

	// 设置值（无过期时间）
	err = client.SetWithExpiration(ctx, key, value, 0)
	require.NoError(t, err)

	// 设置过期时间
	err = client.Expire(ctx, key, 1*time.Minute)
	assert.NoError(t, err)

	// 清理
	_ = client.Delete(ctx, key)
}

func TestClient_Close(t *testing.T) {
	cfg := &config.Config{
		Redis: config.RedisConfig{
			Host:         "localhost",
			Port:         6379,
			Password:     "",
			DB:           0,
			DialTimeout:  5,
			ReadTimeout:  3,
			WriteTimeout: 3,
			PoolSize:     10,
			MinIdleConns: 5,
		},
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	err = client.Close()
	assert.NoError(t, err)
}
