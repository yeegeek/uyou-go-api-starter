package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yeegeek/uyou-go-api-starter/internal/config"
	"log/slog"
)

// MockTask 用于测试的 Mock 任务
type MockTask struct {
	name        string
	runFunc     func(ctx context.Context) error
	runCount    int
	lastRunTime time.Time
}

func (m *MockTask) Name() string {
	return m.name
}

func (m *MockTask) Run(ctx context.Context) error {
	m.runCount++
	m.lastRunTime = time.Now()
	if m.runFunc != nil {
		return m.runFunc(ctx)
	}
	return nil
}

func TestNewScheduler(t *testing.T) {
	cfg := &config.Config{}
	logger := slog.Default()

	scheduler := NewScheduler(cfg, logger)
	assert.NotNil(t, scheduler)
	assert.NotNil(t, scheduler.cron)
	assert.Equal(t, cfg, scheduler.config)
	assert.Equal(t, logger, scheduler.logger)
	assert.NotNil(t, scheduler.tasks)
}

func TestScheduler_AddTask(t *testing.T) {
	cfg := &config.Config{}
	logger := slog.Default()
	scheduler := NewScheduler(cfg, logger)

	task := &MockTask{name: "test-task"}

	err := scheduler.AddTask("*/1 * * * * *", task) // 每秒执行
	assert.NoError(t, err)
	assert.Equal(t, task, scheduler.tasks["test-task"])
}

func TestScheduler_AddTask_InvalidSpec(t *testing.T) {
	cfg := &config.Config{}
	logger := slog.Default()
	scheduler := NewScheduler(cfg, logger)

	task := &MockTask{name: "test-task"}

	err := scheduler.AddTask("invalid-cron-spec", task)
	assert.Error(t, err)
}

func TestScheduler_GetTasks(t *testing.T) {
	cfg := &config.Config{}
	logger := slog.Default()
	scheduler := NewScheduler(cfg, logger)

	task1 := &MockTask{name: "task1"}
	task2 := &MockTask{name: "task2"}

	_ = scheduler.AddTask("*/1 * * * * *", task1)
	_ = scheduler.AddTask("*/2 * * * * *", task2)

	tasks := scheduler.GetTasks()
	assert.Equal(t, 2, len(tasks))
	assert.Equal(t, task1, tasks["task1"])
	assert.Equal(t, task2, tasks["task2"])
}

func TestScheduler_GetEntries(t *testing.T) {
	cfg := &config.Config{}
	logger := slog.Default()
	scheduler := NewScheduler(cfg, logger)

	task := &MockTask{name: "test-task"}
	_ = scheduler.AddTask("*/1 * * * * *", task)

	entries := scheduler.GetEntries()
	assert.Equal(t, 1, len(entries))
}

func TestScheduler_StartStop(t *testing.T) {
	cfg := &config.Config{}
	logger := slog.Default()
	scheduler := NewScheduler(cfg, logger)

	task := &MockTask{name: "test-task"}
	_ = scheduler.AddTask("*/1 * * * * *", task)

	// 启动调度器
	scheduler.Start()

	// 等待一段时间让任务执行
	time.Sleep(2 * time.Second)

	// 停止调度器
	scheduler.Stop()

	// 验证任务至少执行了一次
	assert.Greater(t, task.runCount, 0)
}

func TestScheduler_TaskExecution_Success(t *testing.T) {
	cfg := &config.Config{}
	logger := slog.Default()
	scheduler := NewScheduler(cfg, logger)

	task := &MockTask{
		name: "success-task",
		runFunc: func(ctx context.Context) error {
			return nil
		},
	}

	_ = scheduler.AddTask("*/1 * * * * *", task)
	scheduler.Start()

	time.Sleep(2 * time.Second)
	scheduler.Stop()

	assert.Greater(t, task.runCount, 0)
}

func TestScheduler_TaskExecution_Error(t *testing.T) {
	cfg := &config.Config{}
	logger := slog.Default()
	scheduler := NewScheduler(cfg, logger)

	expectedError := errors.New("task error")
	task := &MockTask{
		name: "error-task",
		runFunc: func(ctx context.Context) error {
			return expectedError
		},
	}

	_ = scheduler.AddTask("*/1 * * * * *", task)
	scheduler.Start()

	time.Sleep(2 * time.Second)
	scheduler.Stop()

	// 即使任务返回错误，也应该执行
	assert.Greater(t, task.runCount, 0)
}

func TestScheduler_TaskExecution_ContextCancellation(t *testing.T) {
	cfg := &config.Config{}
	logger := slog.Default()
	scheduler := NewScheduler(cfg, logger)

	task := &MockTask{
		name: "context-task",
		runFunc: func(ctx context.Context) error {
			// 检查上下文是否被取消
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return nil
			}
		},
	}

	_ = scheduler.AddTask("*/1 * * * * *", task)
	scheduler.Start()

	time.Sleep(1 * time.Second)
	scheduler.Stop()

	assert.Greater(t, task.runCount, 0)
}

func TestScheduler_MultipleTasks(t *testing.T) {
	cfg := &config.Config{}
	logger := slog.Default()
	scheduler := NewScheduler(cfg, logger)

	task1 := &MockTask{name: "task1"}
	task2 := &MockTask{name: "task2"}
	task3 := &MockTask{name: "task3"}

	_ = scheduler.AddTask("*/1 * * * * *", task1)
	_ = scheduler.AddTask("*/2 * * * * *", task2)
	_ = scheduler.AddTask("*/3 * * * * *", task3)

	scheduler.Start()
	time.Sleep(3 * time.Second)
	scheduler.Stop()

	// task1 应该执行最多（每秒）
	assert.Greater(t, task1.runCount, task2.runCount)
	assert.Greater(t, task2.runCount, task3.runCount)
}

func TestScheduler_Stop_GracefulShutdown(t *testing.T) {
	cfg := &config.Config{}
	logger := slog.Default()
	scheduler := NewScheduler(cfg, logger)

	task := &MockTask{
		name: "long-task",
		runFunc: func(ctx context.Context) error {
			time.Sleep(100 * time.Millisecond)
			return nil
		},
	}

	_ = scheduler.AddTask("*/1 * * * * *", task)
	scheduler.Start()

	time.Sleep(500 * time.Millisecond)

	// 停止应该等待当前任务完成
	startStop := time.Now()
	scheduler.Stop()
	stopDuration := time.Since(startStop)

	// 停止应该很快完成（优雅关闭）
	assert.Less(t, stopDuration, 2*time.Second)
}
