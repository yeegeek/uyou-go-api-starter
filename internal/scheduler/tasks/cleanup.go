// Package tasks 提供具体的定时任务实现
package tasks

import (
	"context"
	"log/slog"
)

// CleanupTask 清理任务：定期清理过期数据
//
// 示例：每小时执行一次，清理过期的刷新令牌、验证码等
type CleanupTask struct {
	logger *slog.Logger
}

// NewCleanupTask 创建清理任务
func NewCleanupTask(logger *slog.Logger) *CleanupTask {
	return &CleanupTask{
		logger: logger,
	}
}

// Name 返回任务名称
func (t *CleanupTask) Name() string {
	return "cleanup_expired_data"
}

// Run 执行清理任务
func (t *CleanupTask) Run(ctx context.Context) error {
	t.logger.Info("开始清理过期数据")

	// TODO: 实现具体的清理逻辑
	// 1. 清理过期的刷新令牌
	// 2. 清理过期的验证码
	// 3. 清理过期的会话数据
	// 4. 清理临时文件

	t.logger.Info("过期数据清理完成")
	return nil
}
