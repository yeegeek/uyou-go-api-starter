package mongodb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yeegeek/uyou-go-api-starter/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
		errMsg  string
	}{
		{
			// 注意：本测试在“MongoDB 未运行”的机器上用于验证失败分支；
			// 如果本机恰好有 MongoDB 运行，则会返回成功，这种情况下我们跳过该用例。
			name: "valid config but no mongodb server (expected failure if mongodb is down)",
			cfg: &config.Config{
				MongoDB: config.MongoDBConfig{
					URI:             "mongodb://localhost:27017",
					Database:        "testdb",
					MaxPoolSize:     100,
					MinPoolSize:     10,
					ConnectTimeout:  10,
				},
			},
			wantErr: true,
			errMsg:  "failed to",
		},
		{
			name: "invalid URI",
			cfg: &config.Config{
				MongoDB: config.MongoDBConfig{
					URI:             "invalid-uri",
					Database:        "testdb",
					MaxPoolSize:     100,
					MinPoolSize:     10,
					ConnectTimeout:  10,
				},
			},
			wantErr: true,
			errMsg:  "failed to",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.cfg)
			if tt.wantErr {
				if err == nil && client != nil {
					// 本机 MongoDB 可用时，失败分支不可复现，跳过即可
					_ = client.Close(context.Background())
					t.Skip("MongoDB is available on localhost; skipping failure-path test case")
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
					defer client.Close(context.Background())
				}
			}
		})
	}
}

func TestClient_Database(t *testing.T) {
	cfg := &config.Config{
		MongoDB: config.MongoDBConfig{
			URI:             "mongodb://localhost:27017",
			Database:        "testdb",
			MaxPoolSize:     100,
			MinPoolSize:     10,
			ConnectTimeout:  10,
		},
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("MongoDB not available: %v", err)
	}
	defer client.Close(context.Background())

	database := client.Database()
	assert.NotNil(t, database)
	assert.Equal(t, "testdb", database.Name())
}

func TestClient_Collection(t *testing.T) {
	cfg := &config.Config{
		MongoDB: config.MongoDBConfig{
			URI:             "mongodb://localhost:27017",
			Database:        "testdb",
			MaxPoolSize:     100,
			MinPoolSize:     10,
			ConnectTimeout:  10,
		},
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("MongoDB not available: %v", err)
	}
	defer client.Close(context.Background())

	collection := client.Collection("test_collection")
	assert.NotNil(t, collection)
	assert.Equal(t, "test_collection", collection.Name())
}

func TestClient_HealthCheck(t *testing.T) {
	cfg := &config.Config{
		MongoDB: config.MongoDBConfig{
			URI:             "mongodb://localhost:27017",
			Database:        "testdb",
			MaxPoolSize:     100,
			MinPoolSize:     10,
			ConnectTimeout:  10,
		},
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("MongoDB not available: %v", err)
	}
	defer client.Close(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.HealthCheck(ctx)
	assert.NoError(t, err)
}

func TestClient_Close(t *testing.T) {
	cfg := &config.Config{
		MongoDB: config.MongoDBConfig{
			URI:             "mongodb://localhost:27017",
			Database:        "testdb",
			MaxPoolSize:     100,
			MinPoolSize:     10,
			ConnectTimeout:  10,
		},
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("MongoDB not available: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Close(ctx)
	assert.NoError(t, err)
}

func TestClient_WithTransaction(t *testing.T) {
	cfg := &config.Config{
		MongoDB: config.MongoDBConfig{
			URI:             "mongodb://localhost:27017",
			Database:        "testdb",
			MaxPoolSize:     100,
			MinPoolSize:     10,
			ConnectTimeout:  10,
		},
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skipf("MongoDB not available: %v", err)
	}
	defer client.Close(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试事务执行
	err = client.WithTransaction(ctx, func(sessCtx mongo.SessionContext) error {
		// 模拟事务操作
		return nil
	})
	assert.NoError(t, err)
}
