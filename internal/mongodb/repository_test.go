package mongodb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeegeek/uyou-go-api-starter/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestDocument 测试文档结构
type TestDocument struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Value int                `bson:"value"`
}

func TestNewBaseRepository(t *testing.T) {
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

	repo := NewBaseRepository(client, "test_collection")
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.collection)
}

func TestBaseRepository_Create(t *testing.T) {
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

	repo := NewBaseRepository(client, "test_documents")
	ctx := context.Background()

	doc := TestDocument{
		Name:  "test",
		Value: 123,
	}

	id, err := repo.Create(ctx, doc)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	// 清理
	objectID, _ := primitive.ObjectIDFromHex(id)
	_, _ = repo.collection.DeleteOne(ctx, bson.M{"_id": objectID})
}

func TestBaseRepository_FindByID(t *testing.T) {
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

	repo := NewBaseRepository(client, "test_documents")
	ctx := context.Background()

	// 创建文档
	doc := TestDocument{
		Name:  "test",
		Value: 123,
	}
	id, err := repo.Create(ctx, doc)
	require.NoError(t, err)

	// 查找文档
	var result TestDocument
	err = repo.FindByID(ctx, id, &result)
	assert.NoError(t, err)
	assert.Equal(t, doc.Name, result.Name)
	assert.Equal(t, doc.Value, result.Value)

	// 查找不存在的文档
	var emptyResult TestDocument
	err = repo.FindByID(ctx, primitive.NewObjectID().Hex(), &emptyResult)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "document not found")

	// 清理
	objectID, _ := primitive.ObjectIDFromHex(id)
	_, _ = repo.collection.DeleteOne(ctx, bson.M{"_id": objectID})
}

func TestBaseRepository_FindOne(t *testing.T) {
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

	repo := NewBaseRepository(client, "test_documents")
	ctx := context.Background()

	// 创建文档
	doc := TestDocument{
		Name:  "test",
		Value: 123,
	}
	id, err := repo.Create(ctx, doc)
	require.NoError(t, err)

	// 查找文档
	var result TestDocument
	err = repo.FindOne(ctx, bson.M{"name": "test"}, &result)
	assert.NoError(t, err)
	assert.Equal(t, doc.Name, result.Name)

	// 清理
	objectID, _ := primitive.ObjectIDFromHex(id)
	_, _ = repo.collection.DeleteOne(ctx, bson.M{"_id": objectID})
}

func TestBaseRepository_FindMany(t *testing.T) {
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

	repo := NewBaseRepository(client, "test_documents")
	ctx := context.Background()

	// 创建多个文档
	for i := 0; i < 3; i++ {
		doc := TestDocument{
			Name:  "test",
			Value: i,
		}
		_, _ = repo.Create(ctx, doc)
	}

	// 查找多个文档
	var results []TestDocument
	err = repo.FindMany(ctx, bson.M{"name": "test"}, &results)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 3)

	// 使用选项
	var limitedResults []TestDocument
	opts := options.Find().SetLimit(2)
	err = repo.FindMany(ctx, bson.M{"name": "test"}, &limitedResults, opts)
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(limitedResults), 2)

	// 清理
	_, _ = repo.collection.DeleteMany(ctx, bson.M{"name": "test"})
}

func TestBaseRepository_UpdateByID(t *testing.T) {
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

	repo := NewBaseRepository(client, "test_documents")
	ctx := context.Background()

	// 创建文档
	doc := TestDocument{
		Name:  "test",
		Value: 123,
	}
	id, err := repo.Create(ctx, doc)
	require.NoError(t, err)

	// 更新文档
	update := bson.M{"$set": bson.M{"value": 456}}
	err = repo.UpdateByID(ctx, id, update)
	assert.NoError(t, err)

	// 验证更新
	var result TestDocument
	err = repo.FindByID(ctx, id, &result)
	assert.NoError(t, err)
	assert.Equal(t, 456, result.Value)

	// 清理
	objectID, _ := primitive.ObjectIDFromHex(id)
	_, _ = repo.collection.DeleteOne(ctx, bson.M{"_id": objectID})
}

func TestBaseRepository_DeleteByID(t *testing.T) {
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

	repo := NewBaseRepository(client, "test_documents")
	ctx := context.Background()

	// 创建文档
	doc := TestDocument{
		Name:  "test",
		Value: 123,
	}
	id, err := repo.Create(ctx, doc)
	require.NoError(t, err)

	// 删除文档
	err = repo.DeleteByID(ctx, id)
	assert.NoError(t, err)

	// 验证已删除
	var result TestDocument
	err = repo.FindByID(ctx, id, &result)
	assert.Error(t, err)
}

func TestBaseRepository_Count(t *testing.T) {
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

	repo := NewBaseRepository(client, "test_documents")
	ctx := context.Background()

	// 创建文档
	doc := TestDocument{
		Name:  "count_test",
		Value: 123,
	}
	id, err := repo.Create(ctx, doc)
	require.NoError(t, err)

	// 统计文档
	count, err := repo.Count(ctx, bson.M{"name": "count_test"})
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(1))

	// 清理
	objectID, _ := primitive.ObjectIDFromHex(id)
	_, _ = repo.collection.DeleteOne(ctx, bson.M{"_id": objectID})
}

func TestBaseRepository_Exists(t *testing.T) {
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

	repo := NewBaseRepository(client, "test_documents")
	ctx := context.Background()

	// 不存在的文档
	exists, err := repo.Exists(ctx, bson.M{"name": "nonexistent"})
	assert.NoError(t, err)
	assert.False(t, exists)

	// 创建文档
	doc := TestDocument{
		Name:  "exists_test",
		Value: 123,
	}
	id, err := repo.Create(ctx, doc)
	require.NoError(t, err)

	// 存在的文档
	exists, err = repo.Exists(ctx, bson.M{"name": "exists_test"})
	assert.NoError(t, err)
	assert.True(t, exists)

	// 清理
	objectID, _ := primitive.ObjectIDFromHex(id)
	_, _ = repo.collection.DeleteOne(ctx, bson.M{"_id": objectID})
}
