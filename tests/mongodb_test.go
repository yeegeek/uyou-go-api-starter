//go:build integration

package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeegeek/uyou-go-api-starter/internal/config"
	"github.com/yeegeek/uyou-go-api-starter/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMongoDBConnection(t *testing.T) {
	cfg, err := config.LoadConfig("../configs/config.yaml")
	require.NoError(t, err)

	cfg.MongoDB.Enabled = true
	cfg.MongoDB.URI = "mongodb://localhost:27017"

	client, err := mongodb.NewClient(cfg)
	require.NoError(t, err)
	defer client.Close(context.Background())

	ctx := context.Background()
	err = client.HealthCheck(ctx)
	assert.NoError(t, err)
}

func TestMongoDBRepository(t *testing.T) {
	cfg, err := config.LoadConfig("../configs/config.yaml")
	require.NoError(t, err)

	cfg.MongoDB.Enabled = true
	cfg.MongoDB.URI = "mongodb://localhost:27017"

	client, err := mongodb.NewClient(cfg)
	require.NoError(t, err)
	defer client.Close(context.Background())

	repo := mongodb.NewBaseRepository(client, "test_collection")
	ctx := context.Background()

	type testDoc struct {
		ID   string `bson:"_id,omitempty"`
		Name string `bson:"name"`
	}

	// Clean up
	client.Collection("test_collection").Drop(ctx)

	// Test Create
	doc := testDoc{Name: "test"}
	id, err := repo.Create(ctx, &doc)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	// Test FindByID
	var result testDoc
	err = repo.FindByID(ctx, id, &result)
	require.NoError(t, err)
	assert.Equal(t, "test", result.Name)

	// Test UpdateByID
	err = repo.UpdateByID(ctx, id, bson.M{"\$set": bson.M{"name": "updated"}})
	require.NoError(t, err)

	// Test FindOne after update
	var updatedResult testDoc
	err = repo.FindOne(ctx, bson.M{"_id": result.ID}, &updatedResult)
	require.NoError(t, err)
	assert.Equal(t, "updated", updatedResult.Name)

	// Test DeleteByID
	err = repo.DeleteByID(ctx, id)
	require.NoError(t, err)

	// Test FindByID after delete
	err = repo.FindByID(ctx, id, &result)
	assert.Error(t, err)
}
