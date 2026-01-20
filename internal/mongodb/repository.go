// Package mongodb 提供 MongoDB 仓库基类
package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BaseRepository MongoDB 基础仓库，提供通用 CRUD 操作
type BaseRepository struct {
	collection *mongo.Collection
}

// NewBaseRepository 创建基础仓库实例
func NewBaseRepository(client *Client, collectionName string) *BaseRepository {
	return &BaseRepository{
		collection: client.Collection(collectionName),
	}
}

// Create 创建文档
func (r *BaseRepository) Create(ctx context.Context, document interface{}) (string, error) {
	result, err := r.collection.InsertOne(ctx, document)
	if err != nil {
		return "", fmt.Errorf("failed to create document: %w", err)
	}
	
	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert inserted id to ObjectID")
	}
	
	return id.Hex(), nil
}

// FindByID 根据 ID 查找文档
func (r *BaseRepository) FindByID(ctx context.Context, id string, result interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid object id: %w", err)
	}
	
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("document not found")
		}
		return fmt.Errorf("failed to find document: %w", err)
	}
	
	return nil
}

// FindOne 根据过滤条件查找单个文档
func (r *BaseRepository) FindOne(ctx context.Context, filter interface{}, result interface{}) error {
	err := r.collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("document not found")
		}
		return fmt.Errorf("failed to find document: %w", err)
	}
	return nil
}

// FindMany 根据过滤条件查找多个文档
func (r *BaseRepository) FindMany(ctx context.Context, filter interface{}, results interface{}, opts ...*options.FindOptions) error {
	cursor, err := r.collection.Find(ctx, filter, opts...)
	if err != nil {
		return fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)
	
	if err := cursor.All(ctx, results); err != nil {
		return fmt.Errorf("failed to decode documents: %w", err)
	}
	
	return nil
}

// UpdateByID 根据 ID 更新文档
func (r *BaseRepository) UpdateByID(ctx context.Context, id string, update interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid object id: %w", err)
	}
	
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	
	if result.MatchedCount == 0 {
		return fmt.Errorf("document not found")
	}
	
	return nil
}

// UpdateOne 根据过滤条件更新单个文档
func (r *BaseRepository) UpdateOne(ctx context.Context, filter interface{}, update interface{}) error {
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	
	if result.MatchedCount == 0 {
		return fmt.Errorf("document not found")
	}
	
	return nil
}

// UpdateMany 根据过滤条件更新多个文档
func (r *BaseRepository) UpdateMany(ctx context.Context, filter interface{}, update interface{}) (int64, error) {
	result, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, fmt.Errorf("failed to update documents: %w", err)
	}
	
	return result.ModifiedCount, nil
}

// DeleteByID 根据 ID 删除文档
func (r *BaseRepository) DeleteByID(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid object id: %w", err)
	}
	
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	
	if result.DeletedCount == 0 {
		return fmt.Errorf("document not found")
	}
	
	return nil
}

// DeleteOne 根据过滤条件删除单个文档
func (r *BaseRepository) DeleteOne(ctx context.Context, filter interface{}) error {
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	
	if result.DeletedCount == 0 {
		return fmt.Errorf("document not found")
	}
	
	return nil
}

// DeleteMany 根据过滤条件删除多个文档
func (r *BaseRepository) DeleteMany(ctx context.Context, filter interface{}) (int64, error) {
	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to delete documents: %w", err)
	}
	
	return result.DeletedCount, nil
}

// Count 统计文档数量
func (r *BaseRepository) Count(ctx context.Context, filter interface{}) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}
	return count, nil
}

// Exists 检查文档是否存在
func (r *BaseRepository) Exists(ctx context.Context, filter interface{}) (bool, error) {
	count, err := r.Count(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
