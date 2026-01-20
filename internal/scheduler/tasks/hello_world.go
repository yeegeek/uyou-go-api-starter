// Package tasks 提供具体的定时任务实现
package tasks

import (
	"context"
	"log/slog"
)

// HelloWorldTask 示例任务：每分钟输出 "Hello World"
type HelloWorldTask struct {
	logger *slog.Logger
}

// NewHelloWorldTask 创建 Hello World 任务
func NewHelloWorldTask(logger *slog.Logger) *HelloWorldTask {
	return &HelloWorldTask{
		logger: logger,
	}
}

// Name 返回任务名称
func (t *HelloWorldTask) Name() string {
	return "hello_world"
}

// Run 执行任务
func (t *HelloWorldTask) Run(ctx context.Context) error {
	t.logger.Info("Hello World")
	return nil
}
