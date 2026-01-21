-- Migration: add_performance_index (rollback)
-- Created: 2026-01-20T13:45:56Z

BEGIN;

-- Add your rollback SQL here
-- 删除性能优化索引

-- DROP INDEX IF EXISTS idx_users_created_at;
-- DROP INDEX IF EXISTS idx_users_deleted_at;
-- DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
-- DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
-- DROP INDEX IF EXISTS idx_refresh_tokens_token_family;
-- DROP INDEX IF EXISTS idx_refresh_tokens_revoked;
-- DROP INDEX IF EXISTS idx_user_roles_user_id;
-- DROP INDEX IF EXISTS idx_user_roles_role_id;
-- DROP INDEX IF EXISTS idx_refresh_tokens_active;

COMMIT;
