// Package scheduler 提供定时任务管理功能
package scheduler

import (
	"log/slog"

	"github.com/uyou/uyou-go-api-starter/internal/config"
)

// Manager 任务管理器
type Manager struct {
	scheduler *Scheduler
	config    *config.Config
	logger    *slog.Logger
}

// NewManager 创建任务管理器
func NewManager(cfg *config.Config, logger *slog.Logger) *Manager {
	return &Manager{
		scheduler: NewScheduler(cfg, logger),
		config:    cfg,
		logger:    logger,
	}
}

// RegisterTasks 注册所有定时任务
//
// 这里集中管理所有定时任务的注册
func (m *Manager) RegisterTasks(tasks []TaskConfig) error {
	for _, taskConfig := range tasks {
		if err := m.scheduler.AddTask(taskConfig.Spec, taskConfig.Task); err != nil {
			m.logger.Error("注册定时任务失败",
				"task", taskConfig.Task.Name(),
				"spec", taskConfig.Spec,
				"error", err,
			)
			return err
		}
	}
	return nil
}

// Start 启动任务管理器
func (m *Manager) Start() {
	m.scheduler.Start()
}

// Stop 停止任务管理器
func (m *Manager) Stop() {
	m.scheduler.Stop()
}

// GetScheduler 获取调度器实例
func (m *Manager) GetScheduler() *Scheduler {
	return m.scheduler
}

// TaskConfig 任务配置
type TaskConfig struct {
	Spec string // cron 表达式
	Task Task   // 任务实例
}
