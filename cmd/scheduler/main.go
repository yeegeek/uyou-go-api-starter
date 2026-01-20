// Package main 定时任务调度器入口
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/uyou/uyou-go-api-starter/internal/config"
	"github.com/uyou/uyou-go-api-starter/internal/scheduler"
	"github.com/uyou/uyou-go-api-starter/internal/scheduler/tasks"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		slog.Error("加载配置失败", "error", err)
		os.Exit(1)
	}

	// 初始化日志
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	logger.Info("定时任务调度器启动中...",
		"app_name", cfg.App.Name,
		"environment", cfg.App.Environment,
	)

	// 创建任务管理器
	manager := scheduler.NewManager(cfg, logger)

	// 注册定时任务
	taskConfigs := []scheduler.TaskConfig{
		{
			// 每分钟执行一次 Hello World
			Spec: "0 */1 * * * *",
			Task: tasks.NewHelloWorldTask(logger),
		},
		{
			// 每小时执行一次清理任务
			Spec: "0 0 */1 * * *",
			Task: tasks.NewCleanupTask(logger),
		},
		{
			// 每天凌晨 2 点执行统计任务
			Spec: "0 0 2 * * *",
			Task: tasks.NewStatisticsTask(logger),
		},
	}

	if err := manager.RegisterTasks(taskConfigs); err != nil {
		logger.Error("注册定时任务失败", "error", err)
		os.Exit(1)
	}

	// 启动调度器
	manager.Start()

	logger.Info("定时任务调度器已启动",
		"tasks_count", len(taskConfigs),
	)

	// 打印已注册的任务信息
	for _, taskConfig := range taskConfigs {
		logger.Info("已注册任务",
			"task", taskConfig.Task.Name(),
			"schedule", taskConfig.Spec,
		)
	}

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("收到停止信号，开始优雅关闭...")

	// 停止调度器
	manager.Stop()

	logger.Info("定时任务调度器已停止")
}

// 示例：如何在 API 服务器中集成定时任务
//
// 在 cmd/server/main.go 中添加以下代码：
//
// ```go
// import (
//     "github.com/uyou/uyou-go-api-starter/internal/scheduler"
//     "github.com/uyou/uyou-go-api-starter/internal/scheduler/tasks"
// )
//
// func main() {
//     // ... 现有的初始化代码 ...
//
//     // 创建并启动定时任务管理器
//     taskManager := scheduler.NewManager(cfg, logger)
//     taskConfigs := []scheduler.TaskConfig{
//         {
//             Spec: "0 */1 * * * *",
//             Task: tasks.NewHelloWorldTask(logger),
//         },
//     }
//     taskManager.RegisterTasks(taskConfigs)
//     taskManager.Start()
//
//     // 在优雅关闭时停止任务管理器
//     defer taskManager.Stop()
//
//     // ... 现有的服务器启动代码 ...
// }
// ```
