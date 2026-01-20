// Package scheduler 提供定时任务调度功能
//
// 基于 robfig/cron/v3 实现，支持标准 cron 表达式和任务管理
package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/uyou/uyou-go-api-starter/internal/config"
)

// Task 定义定时任务接口
type Task interface {
	// Name 返回任务名称
	Name() string
	// Run 执行任务逻辑
	Run(ctx context.Context) error
}

// Scheduler 定时任务调度器
type Scheduler struct {
	cron   *cron.Cron
	config *config.Config
	logger *slog.Logger
	tasks  map[string]Task
}

// NewScheduler 创建新的调度器实例
func NewScheduler(cfg *config.Config, logger *slog.Logger) *Scheduler {
	// 创建 cron 实例，使用秒级精度
	c := cron.New(
		cron.WithSeconds(),                                  // 支持秒级调度
		cron.WithChain(cron.Recover(cron.DefaultLogger)),   // 自动恢复 panic
		cron.WithLogger(cron.VerbosePrintfLogger(logger)),  // 使用自定义日志
	)

	return &Scheduler{
		cron:   c,
		config: cfg,
		logger: logger,
		tasks:  make(map[string]Task),
	}
}

// AddTask 添加定时任务
//
// spec: cron 表达式（支持秒级，格式：秒 分 时 日 月 周）
// task: 任务实例
//
// 示例：
//   - "*/1 * * * * *" - 每秒执行
//   - "0 */1 * * * *" - 每分钟执行
//   - "0 0 */1 * * *" - 每小时执行
//   - "0 0 8 * * *" - 每天 8:00 执行
func (s *Scheduler) AddTask(spec string, task Task) error {
	// 记录任务
	s.tasks[task.Name()] = task

	// 包装任务执行逻辑
	_, err := s.cron.AddFunc(spec, func() {
		ctx := context.Background()
		startTime := time.Now()

		s.logger.Info("定时任务开始执行",
			"task", task.Name(),
			"time", startTime.Format(time.RFC3339),
		)

		// 执行任务
		if err := task.Run(ctx); err != nil {
			s.logger.Error("定时任务执行失败",
				"task", task.Name(),
				"error", err,
				"duration", time.Since(startTime),
			)
		} else {
			s.logger.Info("定时任务执行成功",
				"task", task.Name(),
				"duration", time.Since(startTime),
			)
		}
	})

	if err != nil {
		s.logger.Error("添加定时任务失败",
			"task", task.Name(),
			"spec", spec,
			"error", err,
		)
		return err
	}

	s.logger.Info("定时任务已添加",
		"task", task.Name(),
		"spec", spec,
	)

	return nil
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.logger.Info("定时任务调度器启动",
		"tasks_count", len(s.tasks),
	)
	s.cron.Start()
}

// Stop 停止调度器（优雅关闭）
func (s *Scheduler) Stop() {
	s.logger.Info("定时任务调度器停止中...")
	ctx := s.cron.Stop()
	<-ctx.Done()
	s.logger.Info("定时任务调度器已停止")
}

// GetTasks 获取所有已注册的任务
func (s *Scheduler) GetTasks() map[string]Task {
	return s.tasks
}

// GetEntries 获取所有调度条目（用于查看下次执行时间）
func (s *Scheduler) GetEntries() []cron.Entry {
	return s.cron.Entries()
}
