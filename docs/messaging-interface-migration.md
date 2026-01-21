# 消息队列接口重构说明

**日期**: 2026-01-20  
**版本**: v2.0.0

---

## 重构概述

已将 RabbitMQ 的具体实现抽象为 `MessageQueue` 接口，使 `EventBus` 可以支持多种消息队列实现（RabbitMQ、AWS SNS+SQS 等）。

## 变更内容

### 1. 新增接口

**文件**: `internal/messaging/queue.go`

定义了 `MessageQueue` 接口：

```go
type MessageQueue interface {
    Publish(ctx context.Context, topic string, message []byte) error
    Subscribe(ctx context.Context, topic string, handler MessageHandler) error
    Close() error
    HealthCheck() error
}
```

### 2. RabbitMQ 实现接口

**文件**: `internal/messaging/rabbitmq.go`

- ✅ `RabbitMQ` 结构体现在实现 `MessageQueue` 接口
- ✅ 新增 `Subscribe` 方法（实现接口）
- ✅ 保留 `Consume` 方法（向后兼容，已标记为 Deprecated）

### 3. EventBus 使用接口

**文件**: `internal/messaging/events.go`

- ✅ `EventBus.mq` 从 `*RabbitMQ` 改为 `MessageQueue` 接口
- ✅ `NewEventBus` 现在接受 `MessageQueue` 接口
- ✅ 新增 `NewEventBusWithRabbitMQ` 函数（向后兼容）

### 4. 工厂函数

**文件**: `internal/messaging/factory.go`（新建）

提供工厂函数，便于根据配置选择实现：

```go
mq, err := messaging.NewMessageQueue(&cfg.RabbitMQ, messaging.ProviderRabbitMQ)
```

## 使用方式

### 方式 1: 直接使用接口（推荐）

```go
import "github.com/yeegeek/uyou-go-api-starter/internal/messaging"

// 创建 RabbitMQ 实例（实现 MessageQueue 接口）
mq, err := messaging.NewRabbitMQ(&cfg.RabbitMQ)
if err != nil {
    log.Fatal(err)
}
defer mq.Close()

// 使用接口创建 EventBus
eventBus := messaging.NewEventBus(mq)

// 发布事件
event := &messaging.UserCreatedEvent{
    UserID: 1,
    Email:  "user@example.com",
    Name:   "User Name",
}
err = eventBus.Publish(ctx, event)
```

### 方式 2: 使用工厂函数（推荐用于多提供商）

```go
import "github.com/yeegeek/uyou-go-api-starter/internal/messaging"

// 使用工厂函数创建（未来可以切换提供商）
mq, err := messaging.NewMessageQueue(&cfg.RabbitMQ, messaging.ProviderRabbitMQ)
if err != nil {
    log.Fatal(err)
}
defer mq.Close()

eventBus := messaging.NewEventBus(mq)
```

### 方式 3: 向后兼容（仍可用）

```go
// 旧代码仍然可以工作
mq, err := messaging.NewRabbitMQ(&cfg.RabbitMQ)
eventBus := messaging.NewEventBusWithRabbitMQ(mq)  // 或直接使用 NewEventBus(mq)
```

## 未来扩展

### 添加 AWS SNS+SQS 实现

```go
// internal/messaging/aws/sns.go
type AWSSNSQueue struct {
    snsClient *sns.Client
    sqsClient *sqs.Client
    // ...
}

func (a *AWSSNSQueue) Publish(ctx context.Context, topic string, message []byte) error {
    // 实现 SNS 发布
}

func (a *AWSSNSQueue) Subscribe(ctx context.Context, topic string, handler MessageHandler) error {
    // 实现 SQS 订阅
}

// 在 factory.go 中添加
case ProviderAWSSNS:
    return NewAWSSNSQueue(cfg), nil
```

### 使用示例

```go
// 切换到 AWS SNS+SQS
mq, err := messaging.NewMessageQueue(&cfg.RabbitMQ, messaging.ProviderAWSSNS)
eventBus := messaging.NewEventBus(mq)  // 无需修改 EventBus 代码
```

## 向后兼容性

✅ **完全向后兼容**

- `NewRabbitMQ` 函数保持不变
- `RabbitMQ` 结构体的所有方法保持不变
- `Consume` 方法仍然可用（已标记为 Deprecated）
- `NewEventBusWithRabbitMQ` 提供向后兼容

## 测试

所有测试已更新并通过：

```bash
go test ./internal/messaging/... -v
```

## 迁移指南

### 现有代码无需修改

如果您的代码使用以下方式，**无需修改**：

```go
mq, err := messaging.NewRabbitMQ(&cfg.RabbitMQ)
eventBus := messaging.NewEventBus(mq)  // ✅ 仍然工作
```

### 推荐迁移（可选）

为了更好的可扩展性，建议使用工厂函数：

```go
// 旧代码
mq, err := messaging.NewRabbitMQ(&cfg.RabbitMQ)
eventBus := messaging.NewEventBus(mq)

// 新代码（推荐）
mq, err := messaging.NewMessageQueue(&cfg.RabbitMQ, messaging.ProviderRabbitMQ)
eventBus := messaging.NewEventBus(mq)
```

## 接口设计说明

### MessageQueue 接口方法

1. **Publish**: 发布消息到主题/路由键
   - `topic`: 主题或路由键（如 "event.user.created"）
   - `message`: 消息内容（JSON 字节数组）

2. **Subscribe**: 订阅消息
   - `topic`: 主题或路由键（RabbitMQ 中可能被忽略，因为队列已绑定）
   - `handler`: 消息处理函数

3. **Close**: 关闭连接
   - 清理资源

4. **HealthCheck**: 健康检查
   - 验证连接状态

### 设计考虑

- **topic 参数**: 在 RabbitMQ 中，队列已经绑定到交换机，所以 topic 参数会被忽略。但在其他实现（如 SNS+SQS）中，topic 用于路由。
- **MessageHandler**: 使用函数类型，便于实现。
- **向后兼容**: 保留所有旧方法，确保现有代码无需修改。

---

**文档版本**: v1.0  
**最后更新**: 2026-01-20
