//go:build integration

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeegeek/uyou-go-api-starter/internal/config"
	"github.com/yeegeek/uyou-go-api-starter/internal/redis"
)

func TestRedisConnection(t *testing.T) {
	cfg, err := config.LoadConfig("../configs/config.yaml")
	require.NoError(t, err)

	cfg.Redis.Enabled = true
	cfg.Redis.Host = "localhost"

	client, err := redis.NewClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	err = client.HealthCheck(ctx)
	assert.NoError(t, err)
}

func TestRedisCache(t *testing.T) {
	cfg, err := config.LoadConfig("../configs/config.yaml")
	require.NoError(t, err)

	cfg.Redis.Enabled = true
	cfg.Redis.Host = "localhost"

	client, err := redis.NewClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	cache := redis.NewCache(client)
	ctx := context.Background()

	type testData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	key := "test:cache"
	data := testData{Name: "test", Value: 123}

	// Test Set
	err = cache.Set(ctx, key, data, 10*time.Second)
	require.NoError(t, err)

	// Test Get
	var result testData
	err = cache.Get(ctx, key, &result)
	require.NoError(t, err)
	assert.Equal(t, data, result)

	// Test Delete
	err = cache.Delete(ctx, key)
	require.NoError(t, err)

	// Test Get after delete
	err = cache.Get(ctx, key, &result)
	assert.Error(t, err)
}
