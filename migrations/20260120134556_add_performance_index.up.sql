-- Migration: add_performance_index
-- Created: 2026-01-20T13:45:56Z
-- Description: Add description here

BEGIN;

-- Add your migration SQL here
-- 为常用查询添加索引以优化性能

-- 用户表索引
-- -- email 索引已存在（uniqueIndex）
-- CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC);
-- CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- -- 刷新令牌表索引
-- CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
-- CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
-- CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_family ON refresh_tokens(token_family);
-- CREATE INDEX IF NOT EXISTS idx_refresh_tokens_revoked ON refresh_tokens(revoked);

-- -- 用户角色关联表索引
-- CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
-- CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);

-- -- 复合索引：用于查询未过期且未撤销的令牌
-- CREATE INDEX IF NOT EXISTS idx_refresh_tokens_active ON refresh_tokens(user_id, expires_at, revoked) 
-- WHERE revoked = false AND expires_at > NOW();

-- COMMENT ON INDEX idx_users_created_at IS '用户创建时间索引，用于按时间排序查询';
-- COMMENT ON INDEX idx_refresh_tokens_user_id IS '刷新令牌用户ID索引，用于查询用户的所有令牌';
-- COMMENT ON INDEX idx_refresh_tokens_expires_at IS '刷新令牌过期时间索引，用于清理过期令牌';
-- COMMENT ON INDEX idx_refresh_tokens_active IS '活跃刷新令牌复合索引，优化令牌验证查询';


COMMIT;
