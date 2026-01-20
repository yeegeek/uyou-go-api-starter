# 定时任务使用指南

本文档介绍如何在 UYou Go API Starter 中使用定时任务功能。

## 概述

项目集成了 **robfig/cron/v3** 定时任务调度库，提供了强大的 cron 表达式支持和任务管理功能。

### 核心特性

- ✅ 支持秒级精度的 cron 表达式（6 字段格式）
- ✅ 自动 panic 恢复机制
- ✅ 结构化日志记录
- ✅ 任务执行时间统计
- ✅ 优雅关闭支持

---

## 快速开始

### 1. 运行调度器

```bash
# 独立运行定时任务调度器
make scheduler
```

### 2. 查看日志

调度器启动后，会输出以下信息：

```json
{
  "time": "2026-01-20T10:00:00Z",
  "level": "INFO",
  "msg": "定时任务调度器启动中...",
  "app_name": "UYou API Starter API",
  "environment": "development"
}

{
  "time": "2026-01-20T10:00:00Z",
  "level": "INFO",
  "msg": "已注册任务",
  "task": "hello_world",
  "schedule": "0 */1 * * * *"
}

{
  "time": "2026-01-20T10:01:00Z",
  "level": "INFO",
  "msg": "定时任务开始执行",
  "task": "hello_world",
  "time": "2026-01-20T10:01:00Z"
}

{
  "time": "2026-01-20T10:01:00Z",
  "level": "INFO",
  "msg": "Hello World"
}

{
  "time": "2026-01-20T10:01:00Z",
  "level": "INFO",
  "msg": "定时任务执行成功",
  "task": "hello_world",
  "duration": "1.234ms"
}
```

---

## Cron 表达式格式

### 6 字段格式（支持秒级）

```
秒 分 时 日 月 周

字段说明：
- 秒：0-59
- 分：0-59
- 时：0-23
- 日：1-31
- 月：1-12
- 周：0-6（0 = 周日）
```

### 常用表达式示例

| 表达式 | 说明 |
|:-------|:-----|
| `*/1 * * * * *` | 每秒执行 |
| `0 */1 * * * *` | 每分钟执行 |
| `0 0 */1 * * *` | 每小时执行 |
| `0 0 8 * * *` | 每天 8:00 执行 |
| `0 0 2 * * *` | 每天凌晨 2:00 执行 |
| `0 0 8 * * 1` | 每周一 8:00 执行 |
| `0 0 8 1 * *` | 每月 1 号 8:00 执行 |
| `0 0 8-18 * * *` | 每天 8:00-18:00 每小时执行 |
| `0 */30 * * * *` | 每 30 分钟执行 |
| `0 0 8,12,18 * * *` | 每天 8:00、12:00、18:00 执行 |

### 特殊字符

- `*`：匹配任意值
- `/`：步长（如 `*/5` 表示每 5 个单位）
- `,`：列举多个值（如 `1,3,5`）
- `-`：范围（如 `1-5`）

---

## 创建自定义任务

### 步骤 1：实现 Task 接口

在 `internal/scheduler/tasks/` 目录下创建新文件：

```go
// internal/scheduler/tasks/my_task.go
package tasks

import (
    "context"
    "log/slog"
)

// MyTask 自定义任务
type MyTask struct {
    logger *slog.Logger
}

// NewMyTask 创建任务实例
func NewMyTask(logger *slog.Logger) *MyTask {
    return &MyTask{
        logger: logger,
    }
}

// Name 返回任务名称（必须实现）
func (t *MyTask) Name() string {
    return "my_custom_task"
}

// Run 执行任务逻辑（必须实现）
func (t *MyTask) Run(ctx context.Context) error {
    t.logger.Info("开始执行自定义任务")
    
    // TODO: 实现你的业务逻辑
    // 例如：
    // - 清理过期数据
    // - 发送邮件通知
    // - 生成统计报表
    // - 同步第三方数据
    
    t.logger.Info("自定义任务执行完成")
    return nil
}
```

### 步骤 2：注册任务

在 `cmd/scheduler/main.go` 中注册任务：

```go
taskConfigs := []scheduler.TaskConfig{
    // 现有任务...
    
    // 添加你的任务
    {
        Spec: "0 0 3 * * *",  // 每天凌晨 3 点执行
        Task: tasks.NewMyTask(logger),
    },
}
```

### 步骤 3：重启调度器

```bash
make scheduler
```

---

## 内置示例任务

### 1. Hello World 任务

**文件**：`internal/scheduler/tasks/hello_world.go`

**功能**：每分钟输出 "Hello World" 日志

**调度**：`0 */1 * * * *`（每分钟执行）

**用途**：测试调度器是否正常工作

---

### 2. 清理任务

**文件**：`internal/scheduler/tasks/cleanup.go`

**功能**：定期清理过期数据

**调度**：`0 0 */1 * * *`（每小时执行）

**用途**：
- 清理过期的刷新令牌
- 清理过期的验证码
- 清理过期的会话数据
- 清理临时文件

**扩展示例**：

```go
func (t *CleanupTask) Run(ctx context.Context) error {
    t.logger.Info("开始清理过期数据")
    
    // 清理过期的刷新令牌
    if err := t.cleanExpiredTokens(ctx); err != nil {
        t.logger.Error("清理令牌失败", "error", err)
    }
    
    // 清理过期的验证码
    if err := t.cleanExpiredCodes(ctx); err != nil {
        t.logger.Error("清理验证码失败", "error", err)
    }
    
    t.logger.Info("过期数据清理完成")
    return nil
}

func (t *CleanupTask) cleanExpiredTokens(ctx context.Context) error {
    // 实现清理逻辑
    // 例如：DELETE FROM refresh_tokens WHERE expires_at < NOW()
    return nil
}
```

---

### 3. 统计任务

**文件**：`internal/scheduler/tasks/statistics.go`

**功能**：生成每日统计数据

**调度**：`0 0 2 * * *`（每天凌晨 2 点执行）

**用途**：
- 统计新增用户数
- 统计活跃用户数
- 统计消息发送量
- 统计交易金额
- 生成报表并发送给管理员

---

## 在 API 服务器中集成定时任务

如果你想在 API 服务器中同时运行定时任务，可以这样做：

### 修改 `cmd/server/main.go`

```go
package main

import (
    // ... 现有的 import
    "github.com/uyou/uyou-go-api-starter/internal/scheduler"
    "github.com/uyou/uyou-go-api-starter/internal/scheduler/tasks"
)

func main() {
    // ... 现有的初始化代码 ...
    
    // 创建并启动定时任务管理器
    taskManager := scheduler.NewManager(cfg, logger)
    taskConfigs := []scheduler.TaskConfig{
        {
            Spec: "0 */1 * * * *",
            Task: tasks.NewHelloWorldTask(logger),
        },
        {
            Spec: "0 0 */1 * * *",
            Task: tasks.NewCleanupTask(logger),
        },
    }
    
    if err := taskManager.RegisterTasks(taskConfigs); err != nil {
        logger.Error("注册定时任务失败", "error", err)
        os.Exit(1)
    }
    
    taskManager.Start()
    logger.Info("定时任务已启动")
    
    // ... 现有的服务器启动代码 ...
    
    // 在优雅关闭时停止任务管理器
    go func() {
        <-quit
        logger.Info("正在关闭定时任务...")
        taskManager.Stop()
    }()
}
```

---

## 高级用法

### 1. 任务依赖注入

如果任务需要访问数据库或其他服务：

```go
type CleanupTask struct {
    logger *slog.Logger
    db     *gorm.DB
    redis  *redis.Client
}

func NewCleanupTask(logger *slog.Logger, db *gorm.DB, redis *redis.Client) *CleanupTask {
    return &CleanupTask{
        logger: logger,
        db:     db,
        redis:  redis,
    }
}

func (t *CleanupTask) Run(ctx context.Context) error {
    // 使用数据库
    result := t.db.WithContext(ctx).
        Where("expires_at < ?", time.Now()).
        Delete(&RefreshToken{})
    
    t.logger.Info("清理完成", "deleted", result.RowsAffected)
    return nil
}
```

### 2. 错误处理和重试

```go
func (t *MyTask) Run(ctx context.Context) error {
    maxRetries := 3
    
    for i := 0; i < maxRetries; i++ {
        if err := t.doWork(ctx); err != nil {
            t.logger.Error("任务执行失败", 
                "attempt", i+1, 
                "error", err,
            )
            
            if i < maxRetries-1 {
                time.Sleep(time.Second * time.Duration(i+1))
                continue
            }
            return err
        }
        return nil
    }
    
    return fmt.Errorf("任务失败，已重试 %d 次", maxRetries)
}
```

### 3. 任务超时控制

```go
func (t *MyTask) Run(ctx context.Context) error {
    // 设置 5 分钟超时
    ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
    defer cancel()
    
    done := make(chan error, 1)
    
    go func() {
        done <- t.doWork(ctx)
    }()
    
    select {
    case err := <-done:
        return err
    case <-ctx.Done():
        return fmt.Errorf("任务超时: %w", ctx.Err())
    }
}
```

### 4. 分布式锁（防止重复执行）

```go
func (t *MyTask) Run(ctx context.Context) error {
    lockKey := fmt.Sprintf("task:lock:%s", t.Name())
    
    // 尝试获取锁（30 秒过期）
    ok, err := t.redis.SetNX(ctx, lockKey, "1", 30*time.Second).Result()
    if err != nil {
        return err
    }
    
    if !ok {
        t.logger.Info("任务已在其他节点执行，跳过")
        return nil
    }
    
    defer t.redis.Del(ctx, lockKey)
    
    // 执行任务
    return t.doWork(ctx)
}
```

---

## 配置

### 配置文件

在 `configs/config.yaml` 中配置定时任务：

```yaml
scheduler:
  enabled: true
  timezone: "Asia/Shanghai"
```

### 环境变量

```bash
SCHEDULER_ENABLED=true
SCHEDULER_TIMEZONE=Asia/Shanghai
```

---

## 最佳实践

### 1. 任务命名规范

- 使用小写下划线命名：`cleanup_expired_tokens`
- 名称应清晰描述任务功能
- 避免使用缩写

### 2. 日志记录

- 任务开始和结束都应记录日志
- 记录关键操作和结果
- 错误日志应包含足够的上下文信息

### 3. 性能考虑

- 避免在任务中执行长时间阻塞操作
- 大批量数据处理应分批进行
- 使用 context 控制超时

### 4. 错误处理

- 任务失败应返回错误，而不是 panic
- 记录详细的错误信息
- 考虑实现重试机制

### 5. 测试

为每个任务编写单元测试：

```go
func TestCleanupTask_Run(t *testing.T) {
    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
    task := NewCleanupTask(logger)
    
    ctx := context.Background()
    err := task.Run(ctx)
    
    assert.NoError(t, err)
}
```

---

## 常见问题

### Q1: 如何查看任务的下次执行时间？

```go
scheduler := manager.GetScheduler()
entries := scheduler.GetEntries()

for _, entry := range entries {
    fmt.Printf("下次执行时间: %s\n", entry.Next)
}
```

### Q2: 如何动态添加任务？

```go
scheduler := manager.GetScheduler()
newTask := tasks.NewMyTask(logger)
scheduler.AddTask("0 0 4 * * *", newTask)
```

### Q3: 如何临时禁用某个任务？

在 `cmd/scheduler/main.go` 中注释掉对应的任务配置：

```go
taskConfigs := []scheduler.TaskConfig{
    // {
    //     Spec: "0 */1 * * * *",
    //     Task: tasks.NewHelloWorldTask(logger),
    // },
}
```

### Q4: 任务执行失败会影响其他任务吗？

不会。每个任务独立执行，一个任务失败不会影响其他任务。调度器会自动恢复 panic。

---

## 总结

UYou Go API Starter 的定时任务系统提供了：

- ✅ 简单易用的任务接口
- ✅ 强大的 cron 表达式支持
- ✅ 完善的日志和错误处理
- ✅ 灵活的任务管理

你可以轻松地添加自定义任务，满足各种定时执行需求。

---

## 参考资源

- [robfig/cron 官方文档](https://pkg.go.dev/github.com/robfig/cron/v3)
- [Cron 表达式在线生成器](https://crontab.guru/)
