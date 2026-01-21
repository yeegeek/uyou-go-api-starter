// Package user 提供账户锁定功能
package user

import (
	"context"
	"fmt"
	"time"

	"github.com/yeegeek/uyou-go-api-starter/internal/config"
	"github.com/yeegeek/uyou-go-api-starter/internal/redis"
)

// AccountLockoutService 账户锁定服务
type AccountLockoutService struct {
	cache              *redis.Cache
	maxAttempts        int
	lockoutDuration    time.Duration
	enabled            bool
}

// NewAccountLockoutService 创建账户锁定服务
func NewAccountLockoutService(cache *redis.Cache, cfg *config.SecurityConfig) *AccountLockoutService {
	lockoutDuration := time.Duration(cfg.LockoutDuration) * time.Minute
	if lockoutDuration == 0 {
		lockoutDuration = 15 * time.Minute
	}

	maxAttempts := cfg.MaxLoginAttempts
	if maxAttempts == 0 {
		maxAttempts = 5
	}

	return &AccountLockoutService{
		cache:           cache,
		maxAttempts:     maxAttempts,
		lockoutDuration: lockoutDuration,
		enabled:         cache != nil,
	}
}

// RecordFailedAttempt 记录登录失败尝试
func (s *AccountLockoutService) RecordFailedAttempt(ctx context.Context, email string) error {
	if !s.enabled {
		return nil
	}

	key := s.getAttemptsKey(email)
	
	// 获取当前失败次数
	attempts, err := s.cache.Increment(ctx, key, 1)
	if err != nil {
		return fmt.Errorf("failed to increment login attempts: %w", err)
	}

	// 设置过期时间（如果是第一次失败）
	if attempts == 1 {
		if err := s.cache.Expire(ctx, key, s.lockoutDuration); err != nil {
			return fmt.Errorf("failed to set expiration: %w", err)
		}
	}

	// 如果达到最大尝试次数，锁定账户
	if attempts >= int64(s.maxAttempts) {
		lockKey := s.getLockKey(email)
		if err := s.cache.Set(ctx, lockKey, "locked", s.lockoutDuration); err != nil {
			return fmt.Errorf("failed to lock account: %w", err)
		}
	}

	return nil
}

// IsAccountLocked 检查账户是否被锁定
func (s *AccountLockoutService) IsAccountLocked(ctx context.Context, email string) (bool, time.Duration, error) {
	if !s.enabled {
		return false, 0, nil
	}

	lockKey := s.getLockKey(email)
	
	var locked string
	err := s.cache.Get(ctx, lockKey, &locked)
	if err != nil {
		// 键不存在，账户未锁定
		return false, 0, nil
	}

	// 获取剩余锁定时间
	ttl, err := s.cache.TTL(ctx, lockKey)
	if err != nil {
		return true, 0, err
	}

	return true, ttl, nil
}

// ResetFailedAttempts 重置登录失败次数（登录成功后调用）
func (s *AccountLockoutService) ResetFailedAttempts(ctx context.Context, email string) error {
	if !s.enabled {
		return nil
	}

	attemptsKey := s.getAttemptsKey(email)
	lockKey := s.getLockKey(email)

	// 删除失败次数记录
	if err := s.cache.Delete(ctx, attemptsKey); err != nil {
		return fmt.Errorf("failed to delete attempts: %w", err)
	}

	// 删除锁定记录
	if err := s.cache.Delete(ctx, lockKey); err != nil {
		return fmt.Errorf("failed to delete lock: %w", err)
	}

	return nil
}

// GetRemainingAttempts 获取剩余尝试次数
func (s *AccountLockoutService) GetRemainingAttempts(ctx context.Context, email string) (int, error) {
	if !s.enabled {
		return s.maxAttempts, nil
	}

	key := s.getAttemptsKey(email)
	
	var attempts int64
	err := s.cache.Get(ctx, key, &attempts)
	if err != nil {
		// 键不存在，返回最大尝试次数
		return s.maxAttempts, nil
	}

	remaining := s.maxAttempts - int(attempts)
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// getAttemptsKey 获取失败尝试次数的缓存键
func (s *AccountLockoutService) getAttemptsKey(email string) string {
	return fmt.Sprintf("login:attempts:%s", email)
}

// getLockKey 获取账户锁定的缓存键
func (s *AccountLockoutService) getLockKey(email string) string {
	return fmt.Sprintf("login:locked:%s", email)
}
