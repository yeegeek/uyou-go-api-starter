# RabbitMQ 微服务架构模式指南

## 概述

在微服务架构中，RabbitMQ 可以用于服务间的异步通信和事件驱动。根据业务需求，可以选择不同的使用模式。

## 重要说明：RabbitMQ 连接方式

**所有微服务都连接到同一个 RabbitMQ 实例！**

RabbitMQ 是一个独立的中间件服务，所有微服务（中心服务和工作服务）都连接到同一个 RabbitMQ 实例。**只需要配置正确的 URL，发布和订阅行为由代码逻辑控制。**

```
┌──────────┐      ┌──────────┐      ┌──────────┐
│ 中心服务 │      │ 工作服务A│      │ 工作服务B│
│          │      │          │      │          │
└────┬─────┘      └────┬─────┘      └────┬─────┘
     │                │                │
     └────────┬───────┴────────────────┘
              │
         ┌────▼────┐
         │ RabbitMQ│  ← 所有服务连接到这里（共享）
         │ (独立服务)│
         └─────────┘
```

**连接配置（所有服务相同）：**
```yaml
rabbitmq:
  enabled: true
  url: "amqp://guest:guest@rabbitmq:5672/"  # 同一个 RabbitMQ
  exchange: "uyou_events"                    # 同一个 Exchange
  exchange_type: "topic"
  queue: "uyou_queue"
  routing_key: "uyou.#"
```

**发布/订阅控制：**
- 通过代码逻辑控制是否调用 `Publish()` 或 `Subscribe()`
- 不需要配置限制，灵活使用

---

## 三种主要模式

### 模式 1：混合模式（推荐用于大多数场景）⭐

**适用场景：**
- 中心服务需要发布事件给其他服务
- 中心服务也需要订阅其他服务或自己的事件（如任务完成通知、状态更新等）
- 其他微服务只需要响应事件，不需要发布事件
- **这是最灵活的架构模式**

**架构图：**
```
┌─────────────────┐
│   中心服务      │ ──发布事件──> RabbitMQ
│  (API Gateway)  │ <──订阅事件──┘
│ 发布+订阅       │                    │
└─────────────────┘                    │
                                       ├──> ┌──────────┐
                                       │    │ 服务 A   │
                                       │    │consume_only│
                                       ├──> ┌──────────┘
                                       │    ┌──────────┐
                                       │    │ 服务 B   │
                                       └──> │consume_only│
                                            └──────────┘
```

**配置示例：**

**所有服务配置相同** (`config.yaml`):
```yaml
rabbitmq:
  enabled: true
  url: "amqp://guest:guest@rabbitmq:5672/"  # 同一个 RabbitMQ
  exchange: "uyou_events"                    # 同一个 Exchange
  exchange_type: "topic"
  queue: "uyou_queue"
  routing_key: "uyou.#"
  prefetch_count: 10
```

**区别在于代码使用：**
- 中心服务：代码中调用 `Publish()` 和 `Subscribe()`
- 工作服务：代码中根据需要调用 `Publish()` 和/或 `Subscribe()`

**使用场景示例：**
- 中心服务发布：`user.created` → 工作服务处理
- 工作服务完成：`task.completed` → 中心服务订阅并更新状态
- 中心服务发布：`notification.send` → 通知服务处理
- 通知服务完成：`notification.sent` → 中心服务订阅并记录日志

**优点：**
- ✅ 中心服务可以接收反馈和状态更新
- ✅ 支持双向通信
- ✅ 工作服务保持简单（只处理任务）
- ✅ 适合大多数业务场景

**缺点：**
- ⚠️ 中心服务需要处理更多逻辑
- ⚠️ 需要合理设计事件类型避免循环

---

### 模式 2：双向通信模式（工作服务也需要发布）⭐

**适用场景：**
- 中心服务发布任务给工作服务
- 工作服务处理完成后，也需要发布事件（如任务完成、状态更新）
- 中心服务订阅工作服务的事件
- **这是最灵活的架构模式**

**架构图：**
```
┌─────────────────┐      ┌──────────┐      ┌──────────┐
│   中心服务      │      │ 工作服务A│      │ 工作服务B│
│ 发布+订阅       │      │ 发布+订阅│      │ 仅订阅   │
└────┬────────────┘      └────┬─────┘      └────┬─────┘
     │                        │                │
     └──────────┬──────────────┴────────────────┘
                │
         ┌──────▼──────┐
         │  RabbitMQ   │
         │  (共享实例) │
         └─────────────┘
```

**配置示例：**

**所有服务配置相同** (`config.yaml`):
```yaml
rabbitmq:
  enabled: true
  url: "amqp://guest:guest@rabbitmq:5672/"  # 同一个 RabbitMQ
  exchange: "uyou_events"                    # 同一个 Exchange
  exchange_type: "topic"
  queue: "uyou_queue"
  routing_key: "uyou.#"
  prefetch_count: 10
```

**区别在于代码使用：**
- 中心服务：代码中调用 `Publish()` 和 `Subscribe()`
- 工作服务A：代码中调用 `Publish()` 和 `Subscribe()`（双向）
- 工作服务B：代码中只调用 `Subscribe()`（只订阅）

**事件流示例：**
1. 中心服务发布：`task.assign` → 工作服务A订阅并处理
2. 工作服务A完成：`task.completed` → 中心服务订阅并更新状态
3. 中心服务发布：`notification.send` → 工作服务B订阅并发送
4. 工作服务A发布：`log.created` → 日志服务订阅并记录

**优点：**
- ✅ 完全双向通信
- ✅ 工作服务可以反馈结果
- ✅ 支持复杂的事件流
- ✅ 服务间解耦

**缺点：**
- ⚠️ 需要合理设计事件类型
- ⚠️ 需要防止事件循环
- ⚠️ 监控和调试更复杂

---

### 模式 3：纯中心化事件发布模式（简单场景）

**适用场景：**
- 有一个中心服务（如 API Gateway、主服务）负责所有业务操作
- 其他微服务只需要响应事件，不需要发布事件
- 希望集中管理事件流

**架构图：**
```
┌─────────────────┐
│   中心服务      │
│  (API Gateway)  │ ──发布事件──> RabbitMQ
│  (只发布)       │
└─────────────────┘                    │
                                       ├──> ┌──────────┐
                                       │    │ 服务 A   │
                                       │    │(只订阅)  │
                                       ├──> ┌──────────┘
                                       │    ┌──────────┐
                                       │    │ 服务 B   │
                                       └──> │(只订阅)  │
                                            └──────────┘
```

**配置示例：**

**所有服务配置相同** (`config.yaml`):
```yaml
rabbitmq:
  enabled: true
  url: "amqp://guest:guest@rabbitmq:5672/"
  exchange: "uyou_events"
  exchange_type: "topic"
  queue: "uyou_queue"
  routing_key: "uyou.#"
  prefetch_count: 10
```

**区别在于代码使用：**
- 中心服务：代码中只调用 `Publish()`（只发布）
- 工作服务：代码中只调用 `Subscribe()`（只订阅）

**优点：**
- ✅ 简单清晰，易于管理
- ✅ 事件流集中，便于监控和调试
- ✅ 适合小型到中型微服务架构

**缺点：**
- ❌ 中心服务可能成为瓶颈
- ❌ 服务间耦合度较高
- ❌ 扩展性有限

---

### 模式 4：分布式事件驱动模式（复杂场景）

**适用场景：**
- 多个微服务都需要发布和订阅事件
- 服务间需要双向通信
- 希望服务高度自治

**架构图：**
```
┌──────────┐      ┌──────────┐      ┌──────────┐
│ 服务 A   │      │ 服务 B   │      │ 服务 C   │
│发布+订阅 │      │发布+订阅 │      │发布+订阅 │
└────┬─────┘      └────┬─────┘      └────┬─────┘
     │                 │                 │
     └─────────┬───────┴─────────┬───────┘
               │                 │
          ┌────▼─────────────────▼────┐
          │      RabbitMQ              │
          │   (Event Bus)              │
          └────────────────────────────┘
```

**配置示例：**

所有服务使用相同配置：
```yaml
rabbitmq:
  enabled: true
  url: "amqp://guest:guest@rabbitmq:5672/"
  exchange: "uyou_events"
  exchange_type: "topic"
  queue: "uyou_queue"
  routing_key: "uyou.#"
  
  
```

**优点：**
- ✅ 服务高度自治
- ✅ 更好的扩展性
- ✅ 服务间解耦
- ✅ 适合大型微服务架构

**缺点：**
- ❌ 事件流分散，监控复杂
- ❌ 需要更好的治理机制
- ❌ 可能出现事件循环

---

## 使用建议

### 1. 根据服务类型选择模式

| 服务类型 | 推荐模式 | 配置 | 说明 |
|---------|---------|------|------|
| **服务类型** | **配置** | **代码使用** | **说明** |
|------------|---------|------------|---------|
| **中心服务** | 相同配置 | `Publish()` + `Subscribe()` | 发布任务，接收反馈 |
| **工作服务（需要反馈）** | 相同配置 | `Publish()` + `Subscribe()` | 订阅任务，发布完成事件 |
| **工作服务（只处理）** | 相同配置 | `Subscribe()` | 只订阅事件，处理任务 |
| **纯发布服务** | 相同配置 | `Publish()` | 只发布事件 |

### 2. 代码使用示例

#### 中心服务（发布+订阅）⭐ 推荐模式

```go
// 初始化 EventBus
if cfg.RabbitMQ.Enabled {
    mq, err := messaging.NewRabbitMQ(&cfg.RabbitMQ)
    if err != nil {
        log.Fatal("Failed to connect to RabbitMQ:", err)
    }
    eventBus := messaging.NewEventBus(mq)
    
    // 1. 发布事件（给工作服务）
    event := &messaging.UserCreatedEvent{
        UserID: user.ID,
        Email:  user.Email,
        Name:   user.Name,
    }
    eventBus.Publish(ctx, event)
    
    // 2. 订阅事件（接收工作服务的反馈）
    ctx := context.Background()
    go func() {
        eventBus.Subscribe(ctx, messaging.EventHandlerFunc(func(ctx context.Context, event *messaging.BaseEvent) error {
            switch event.Type {
            case "task.completed":
                // 处理任务完成事件
                return handleTaskCompleted(ctx, event)
            case "notification.sent":
                // 处理通知发送完成事件
                return handleNotificationSent(ctx, event)
            }
            return nil
        }))
    }()
}
```

#### 工作服务（订阅+发布）⭐ 双向通信

```go
// 初始化 EventBus
if cfg.RabbitMQ.Enabled {
    mq, err := messaging.NewRabbitMQ(&cfg.RabbitMQ)
    if err != nil {
        log.Fatal("Failed to connect to RabbitMQ:", err)
    }
    eventBus := messaging.NewEventBus(mq)
    
    // 1. 订阅任务事件
    ctx := context.Background()
    go func() {
        eventBus.Subscribe(ctx, messaging.EventHandlerFunc(func(ctx context.Context, event *messaging.BaseEvent) error {
            switch event.Type {
            case messaging.EventTypeUserCreated:
                // 处理用户创建任务
                result := processUserCreated(ctx, event)
                
                // 2. 发布任务完成事件（反馈给中心服务）
                completedEvent := &TaskCompletedEvent{
                    TaskID: result.TaskID,
                    Status: "completed",
                }
                eventBus.Publish(ctx, completedEvent)
                return nil
            }
            return nil
        }))
    }()
}
```

#### 工作服务（仅订阅）

```go
// 初始化 EventBus
if cfg.RabbitMQ.Enabled {
    mq, err := messaging.NewRabbitMQ(&cfg.RabbitMQ)
    if err != nil {
        log.Fatal("Failed to connect to RabbitMQ:", err)
    }
    eventBus := messaging.NewEventBus(mq)
    
    // 只订阅事件，不发布
    ctx := context.Background()
    go func() {
        eventBus.Subscribe(ctx, messaging.EventHandlerFunc(func(ctx context.Context, event *messaging.BaseEvent) error {
            switch event.Type {
            case messaging.EventTypeUserCreated:
                // 处理用户创建事件
                return handleUserCreated(ctx, event)
            }
            return nil
        }))
    }()
}
```

### 3. 环境变量配置

**所有服务配置相同** (`.env`):
```bash
RABBITMQ_ENABLED=true
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
RABBITMQ_EXCHANGE=uyou_events
RABBITMQ_EXCHANGE_TYPE=topic
RABBITMQ_QUEUE=uyou_queue
RABBITMQ_ROUTING_KEY=uyou.#
RABBITMQ_PREFETCH_COUNT=10
```

**区别在于代码使用：**
- 中心服务：代码中调用 `Publish()` 和 `Subscribe()`
- 工作服务：根据需求在代码中调用 `Publish()` 和/或 `Subscribe()`

### 4. 最佳实践

1. **事件命名规范**
   - 使用点分隔的命名：`service.action`（如 `user.created`）
   - 保持事件类型的一致性

2. **路由键设计**
   - 使用通配符：`uyou.#` 匹配所有 `uyou.` 开头的事件
   - 使用 `service.action` 格式便于路由

3. **错误处理**
   - 实现重试机制
   - 使用死信队列处理失败消息
   - 记录所有事件处理日志

4. **监控和追踪**
   - 在事件中添加 TraceID
   - 监控消息队列长度
   - 记录事件发布和消费的指标

5. **性能优化**
   - 设置合适的 `prefetch_count`
   - 使用批量处理
   - 避免阻塞事件处理

---

## 总结

### 推荐配置

**所有服务使用相同的配置：**
```yaml
rabbitmq:
  enabled: true
  url: "amqp://guest:guest@rabbitmq:5672/"
  exchange: "uyou_events"
  exchange_type: "topic"
  queue: "uyou_queue"
  routing_key: "uyou.#"
  prefetch_count: 10
```

**区别在于代码使用：**

#### 场景 1：工作服务需要反馈（推荐）⭐

**中心服务代码：**
- 调用 `Publish()` 发布任务事件
- 调用 `Subscribe()` 接收完成反馈

**工作服务代码：**
- 调用 `Subscribe()` 订阅任务
- 调用 `Publish()` 发布完成事件

#### 场景 2：工作服务只处理任务

**中心服务代码：**
- 调用 `Publish()` 发布任务事件
- 调用 `Subscribe()` 接收其他事件

**工作服务代码：**
- 调用 `Subscribe()` 订阅任务
- 不调用 `Publish()`（通过 HTTP/gRPC 反馈结果）

### 代码使用矩阵

| 服务类型 | 调用 Publish() | 调用 Subscribe() | 说明 |
|---------|---------------|-----------------|------|
| **中心服务** | ✅ | ✅ | 发布任务，接收反馈 |
| **工作服务（需要反馈）** | ✅ | ✅ | 订阅任务，发布完成事件 |
| **工作服务（只处理）** | ❌ | ✅ | 只订阅任务，通过 HTTP/gRPC 反馈 |
| **纯发布服务** | ✅ | ❌ | 只发布事件 |

### 最佳实践

1. **中心服务**：代码中调用 `Publish()` 和 `Subscribe()`
   - 发布任务事件
   - 订阅任务完成、状态更新等事件
   - 实现完整的双向通信

2. **工作服务（需要反馈）**：代码中调用 `Publish()` 和 `Subscribe()`
   - 订阅任务事件
   - 处理完成后发布完成事件
   - 通过 RabbitMQ 反馈结果（推荐）

3. **工作服务（只处理）**：代码中只调用 `Subscribe()`
   - 只关注处理任务
   - 保持服务简单
   - 通过 HTTP/gRPC 反馈结果（如果需要）

3. **事件命名规范**
   - 中心服务发布：`task.*`, `user.*`, `order.*`
   - 工作服务反馈：`task.completed`, `task.failed`（通过 HTTP 回调或 gRPC）

**核心原则：所有服务连接到同一个 RabbitMQ，通过代码逻辑控制发布/订阅行为，配置保持简单。**
