// Package user 提供用户缓存服务
package user

import (
	"context"
	"fmt"
	"time"

	"github.com/yeegeek/uyou-go-api-starter/internal/redis"
)

const (
	// 用户信息缓存 TTL
	userCacheTTL = 5 * time.Minute
	// 用户角色缓存 TTL
	roleCacheTTL = 10 * time.Minute
)

// CachedService 带缓存的用户服务
type CachedService struct {
	service Service
	cache   *redis.Cache
}

// NewCachedService 创建带缓存的用户服务
func NewCachedService(service Service, cache *redis.Cache) Service {
	if cache == nil {
		// 如果没有 Redis，直接返回原始服务
		return service
	}

	return &CachedService{
		service: service,
		cache:   cache,
	}
}

// GetUserByID 获取用户信息（带缓存）
func (s *CachedService) GetUserByID(ctx context.Context, id uint) (*User, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("user:%d", id)
	var user User

	err := s.cache.Get(ctx, cacheKey, &user)
	if err == nil {
		// 缓存命中
		return &user, nil
	}

	// 缓存未命中，从数据库查询
	userPtr, err := s.service.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	_ = s.cache.Set(ctx, cacheKey, userPtr, userCacheTTL)

	return userPtr, nil
}

// UpdateUser 更新用户信息（清除缓存）
func (s *CachedService) UpdateUser(ctx context.Context, id uint, req UpdateUserRequest) (*User, error) {
	user, err := s.service.UpdateUser(ctx, id, req)
	if err != nil {
		return nil, err
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("user:%d", id)
	_ = s.cache.Delete(ctx, cacheKey)

	return user, nil
}

// DeleteUser 删除用户（清除缓存）
func (s *CachedService) DeleteUser(ctx context.Context, id uint) error {
	err := s.service.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("user:%d", id)
	_ = s.cache.Delete(ctx, cacheKey)

	return nil
}

// RegisterUser 注册用户（不缓存）
func (s *CachedService) RegisterUser(ctx context.Context, req RegisterRequest) (*User, error) {
	return s.service.RegisterUser(ctx, req)
}

// AuthenticateUser 认证用户（不缓存）
func (s *CachedService) AuthenticateUser(ctx context.Context, req LoginRequest) (*User, error) {
	return s.service.AuthenticateUser(ctx, req)
}

// ListUsers 列出用户（不缓存）
func (s *CachedService) ListUsers(ctx context.Context, filters UserFilterParams, page, perPage int) ([]User, int64, error) {
	return s.service.ListUsers(ctx, filters, page, perPage)
}

// PromoteToAdmin 提升为管理员（清除缓存）
func (s *CachedService) PromoteToAdmin(ctx context.Context, userID uint) error {
	err := s.service.PromoteToAdmin(ctx, userID)
	if err != nil {
		return err
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("user:%d", userID)
	_ = s.cache.Delete(ctx, cacheKey)

	return nil
}

// InvalidateUserCache 使用户缓存失效
func (s *CachedService) InvalidateUserCache(ctx context.Context, userID uint) error {
	cacheKey := fmt.Sprintf("user:%d", userID)
	return s.cache.Delete(ctx, cacheKey)
}
