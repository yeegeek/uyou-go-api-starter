// Package redis 提供缓存相关功能
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Cache 提供高级缓存操作
type Cache struct {
	client *Client
}

// NewCache 创建缓存实例
func NewCache(client *Client) *Cache {
	return &Cache{client: client}
}

// Set 设置缓存（自动序列化为 JSON）
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	return c.client.SetWithExpiration(ctx, key, data, expiration)
}

// Get 获取缓存（自动反序列化 JSON）
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.GetString(ctx, key)
	if err != nil {
		return err
	}
	
	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}
	return nil
}

// GetOrSet 获取缓存，如果不存在则执行函数并缓存结果
func (c *Cache) GetOrSet(ctx context.Context, key string, dest interface{}, expiration time.Duration, fn func() (interface{}, error)) error {
	// 尝试从缓存获取
	err := c.Get(ctx, key, dest)
	if err == nil {
		return nil
	}
	
	// 缓存不存在，执行函数
	value, err := fn()
	if err != nil {
		return err
	}
	
	// 设置缓存
	if err := c.Set(ctx, key, value, expiration); err != nil {
		return err
	}
	
	// 将结果复制到 dest
	data, _ := json.Marshal(value)
	return json.Unmarshal(data, dest)
}

// Delete 删除缓存
func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	return c.client.Delete(ctx, keys...)
}

// DeletePattern 删除匹配模式的所有键
func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	var cursor uint64
	for {
		keys, nextCursor, err := c.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		
		if len(keys) > 0 {
			if err := c.client.Delete(ctx, keys...); err != nil {
				return err
			}
		}
		
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

// Remember 缓存记忆模式：先查缓存，不存在则执行函数并缓存
func (c *Cache) Remember(ctx context.Context, key string, expiration time.Duration, fn func() (interface{}, error)) (interface{}, error) {
	// 尝试从缓存获取
	var result interface{}
	err := c.Get(ctx, key, &result)
	if err == nil {
		return result, nil
	}
	
	// 执行函数
	result, err = fn()
	if err != nil {
		return nil, err
	}
	
	// 设置缓存
	if err := c.Set(ctx, key, result, expiration); err != nil {
		return result, err
	}
	
	return result, nil
}
