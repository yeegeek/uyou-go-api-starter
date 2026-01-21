package scheduler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yeegeek/uyou-go-api-starter/internal/config"
	"log/slog"
)

func TestNewManager(t *testing.T) {
	cfg := &config.Config{}
	logger := slog.Default()

	manager := NewManager(cfg, logger)
	assert.NotNil(t, manager)
	assert.NotNil(t, manager.scheduler)
	assert.Equal(t, cfg, manager.config)
	assert.Equal(t, logger, manager.logger)
}

func TestManager_RegisterTasks(t *testing.T) {
	cfg := &config.Config{}
	logger := slog.Default()
	manager := NewManager(cfg, logger)

	task1 := &MockTask{name: "task1"}
	task2 := &MockTask{name: "task2"}

	tasks := []TaskConfig{
		{Spec: "*/1 * * * * *", Task: task1},
		{Spec: "*/2 * * * * *", Task: task2},
	}

	err := manager.RegisterTasks(tasks)
	assert.NoError(t, err)

	registeredTasks := manager.GetScheduler().GetTasks()
	assert.Equal(t, 2, len(registeredTasks))
}

func TestManager_RegisterTasks_InvalidSpec(t *testing.T) {
	cfg := &config.Config{}
	logger := slog.Default()
	manager := NewManager(cfg, logger)

	task := &MockTask{name: "task1"}

	tasks := []TaskConfig{
		{Spec: "invalid-spec", Task: task},
	}

	err := manager.RegisterTasks(tasks)
	assert.Error(t, err)
}

func TestManager_StartStop(t *testing.T) {
	cfg := &config.Config{}
	logger := slog.Default()
	manager := NewManager(cfg, logger)

	task := &MockTask{name: "test-task"}
	tasks := []TaskConfig{
		{Spec: "*/1 * * * * *", Task: task},
	}

	_ = manager.RegisterTasks(tasks)
	manager.Start()

	// 等待任务执行
	time.Sleep(2 * time.Second)

	manager.Stop()

	assert.Greater(t, task.runCount, 0)
}

func TestManager_GetScheduler(t *testing.T) {
	cfg := &config.Config{}
	logger := slog.Default()
	manager := NewManager(cfg, logger)

	scheduler := manager.GetScheduler()
	assert.NotNil(t, scheduler)
	assert.Equal(t, manager.scheduler, scheduler)
}
