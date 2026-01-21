package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeegeek/uyou-go-api-starter/internal/config"
)

func TestNewCache(t *testing.T) {
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

	cache := NewCache(client)
	assert.NotNil(t, cache)
	assert.Equal(t, client, cache.client)
}

func TestCache_Set(t *testing.T) {
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

	cache := NewCache(client)
	ctx := context.Background()
	key := "test:cache:set"

	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	testData := TestStruct{
		Name:  "test",
		Value: 123,
	}

	err = cache.Set(ctx, key, testData, 1*time.Minute)
	assert.NoError(t, err)

	// 清理
	_ = cache.Delete(ctx, key)
}

func TestCache_Get(t *testing.T) {
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

	cache := NewCache(client)
	ctx := context.Background()
	key := "test:cache:get"

	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	testData := TestStruct{
		Name:  "test",
		Value: 123,
	}

	// 设置值
	err = cache.Set(ctx, key, testData, 1*time.Minute)
	require.NoError(t, err)

	// 获取值
	var result TestStruct
	err = cache.Get(ctx, key, &result)
	assert.NoError(t, err)
	assert.Equal(t, testData.Name, result.Name)
	assert.Equal(t, testData.Value, result.Value)

	// 获取不存在的键
	var emptyResult TestStruct
	err = cache.Get(ctx, "test:cache:nonexistent", &emptyResult)
	assert.Error(t, err)

	// 清理
	_ = cache.Delete(ctx, key)
}

func TestCache_GetOrSet(t *testing.T) {
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

	cache := NewCache(client)
	ctx := context.Background()
	key := "test:cache:getorset"

	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	callCount := 0
	fn := func() (interface{}, error) {
		callCount++
		return TestStruct{
			Name:  "generated",
			Value: callCount,
		}, nil
	}

	// 第一次调用，应该执行函数
	var result TestStruct
	err = cache.GetOrSet(ctx, key, &result, 1*time.Minute, fn)
	assert.NoError(t, err)
	assert.Equal(t, "generated", result.Name)
	assert.Equal(t, 1, result.Value)
	assert.Equal(t, 1, callCount)

	// 第二次调用，应该从缓存获取
	var cachedResult TestStruct
	err = cache.GetOrSet(ctx, key, &cachedResult, 1*time.Minute, fn)
	assert.NoError(t, err)
	assert.Equal(t, "generated", cachedResult.Name)
	assert.Equal(t, 1, cachedResult.Value)
	assert.Equal(t, 1, callCount) // 函数不应该再次被调用

	// 清理
	_ = cache.Delete(ctx, key)
}

func TestCache_Remember(t *testing.T) {
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

	cache := NewCache(client)
	ctx := context.Background()
	key := "test:cache:remember"

	callCount := 0
	fn := func() (interface{}, error) {
		callCount++
		return map[string]interface{}{
			"value": callCount,
		}, nil
	}

	// 第一次调用，应该执行函数
	result, err := cache.Remember(ctx, key, 1*time.Minute, fn)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, callCount)

	// 第二次调用，应该从缓存获取
	cachedResult, err := cache.Remember(ctx, key, 1*time.Minute, fn)
	assert.NoError(t, err)
	assert.NotNil(t, cachedResult)
	assert.Equal(t, 1, callCount) // 函数不应该再次被调用

	// 清理
	_ = cache.Delete(ctx, key)
}

func TestCache_Delete(t *testing.T) {
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

	cache := NewCache(client)
	ctx := context.Background()
	key1 := "test:cache:delete:1"
	key2 := "test:cache:delete:2"

	// 设置值
	err = cache.Set(ctx, key1, "value1", 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, key2, "value2", 1*time.Minute)
	require.NoError(t, err)

	// 删除多个键
	err = cache.Delete(ctx, key1, key2)
	assert.NoError(t, err)

	// 验证已删除
	var result string
	err = cache.Get(ctx, key1, &result)
	assert.Error(t, err)
}

func TestCache_Increment(t *testing.T) {
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

	cache := NewCache(client)
	ctx := context.Background()
	key := "test:cache:increment"

	// 递增 5
	val, err := cache.Increment(ctx, key, 5)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), val)

	// 再递增 3
	val, err = cache.Increment(ctx, key, 3)
	assert.NoError(t, err)
	assert.Equal(t, int64(8), val)

	// 清理
	_ = cache.Delete(ctx, key)
}

func TestCache_Expire(t *testing.T) {
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

	cache := NewCache(client)
	ctx := context.Background()
	key := "test:cache:expire"

	// 设置值
	err = cache.Set(ctx, key, "value", 0)
	require.NoError(t, err)

	// 设置过期时间
	err = cache.Expire(ctx, key, 1*time.Minute)
	assert.NoError(t, err)

	// 清理
	_ = cache.Delete(ctx, key)
}

func TestCache_TTL(t *testing.T) {
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

	cache := NewCache(client)
	ctx := context.Background()
	key := "test:cache:ttl"

	// 设置值（带过期时间）
	err = cache.Set(ctx, key, "value", 1*time.Minute)
	require.NoError(t, err)

	// 获取 TTL
	ttl, err := cache.TTL(ctx, key)
	assert.NoError(t, err)
	assert.Greater(t, ttl, time.Duration(0))
	assert.LessOrEqual(t, ttl, 1*time.Minute)

	// 清理
	_ = cache.Delete(ctx, key)
}
