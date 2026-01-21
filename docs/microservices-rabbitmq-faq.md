# RabbitMQ 微服务架构 FAQ

## Q1: 多个微服务可以连接到同一个 RabbitMQ 吗？

**A: 可以！** RabbitMQ 是一个独立的中间件服务，所有微服务都连接到同一个 RabbitMQ 实例。

```
所有服务 → 同一个 RabbitMQ 实例
```

**配置示例（所有服务相同）：**
```yaml
rabbitmq:
  enabled: true
  url: "amqp://guest:guest@rabbitmq:5672/"  # 同一个 RabbitMQ
  exchange: "uyou_events"                    # 同一个 Exchange
  exchange_type: "topic"
  queue: "uyou_queue"
  routing_key: "uyou.#"
```

**区别在于代码使用：**
- 中心服务：代码中调用 `Publish()` 和 `Subscribe()`
- 工作服务A：代码中调用 `Publish()` 和 `Subscribe()`（如果需要发布）
- 工作服务B：代码中只调用 `Subscribe()`（如果只订阅）

---

## Q2: 工作服务可以发布事件到中心服务的 RabbitMQ 吗？

**A: 可以！** 工作服务连接到同一个 RabbitMQ，可以发布事件。

**架构：**
```
工作服务 → 发布事件 → RabbitMQ → 中心服务订阅
```

**配置（所有服务相同）：**
```yaml
rabbitmq:
  enabled: true
  url: "amqp://guest:guest@rabbitmq:5672/"  # 同一个 RabbitMQ
  exchange: "uyou_events"
  exchange_type: "topic"
  queue: "uyou_queue"
  routing_key: "uyou.#"
```

**代码示例：**
```go
// 工作服务处理完任务后发布事件
eventBus.Publish(ctx, &TaskCompletedEvent{
    TaskID: taskID,
    Status: "completed",
})
```

---

## Q3: 如何防止事件循环？

**A: 使用事件命名规范和路由键设计**

**建议：**
1. 使用服务前缀：`service.action`
   - 中心服务：`center.*`
   - 工作服务A：`worker-a.*`
   - 工作服务B：`worker-b.*`

2. 使用事件类型区分：
   - 任务事件：`task.*`
   - 反馈事件：`feedback.*`
   - 通知事件：`notification.*`

3. 路由键设计：
   ```yaml
   # 中心服务订阅工作服务的反馈
   routing_key: "feedback.#"
   
   # 工作服务订阅中心服务的任务
   routing_key: "task.#"
   ```

---

## Q4: 如何确保消息不丢失？

**A: 使用持久化和确认机制**

**配置：**
```yaml
rabbitmq:
  exchange: "uyou_events"
  exchange_type: "topic"
  queue: "uyou_queue"
  # Exchange 和 Queue 都设置为持久化（代码中已实现）
```

**代码中已实现：**
- Exchange 持久化：`durable: true`
- Queue 持久化：`durable: true`
- 消息持久化：`DeliveryMode: amqp.Persistent`
- 手动确认：`autoAck: false`

---

## Q5: 如何处理消息处理失败？

**A: 使用 Nack 和重试机制**

**代码中已实现：**
```go
if err := handler(msg.Body); err != nil {
    // 拒绝消息并重新入队
    msg.Nack(false, true)
} else {
    // 确认消息
    msg.Ack(false)
}
```

**建议：**
- 实现重试次数限制
- 使用死信队列（DLQ）处理失败消息
- 记录失败日志

---

## Q6: 如何监控消息队列？

**A: 使用 RabbitMQ 管理界面和 Prometheus**

**访问管理界面：**
- URL: `http://localhost:15672`
- 用户名：`guest`
- 密码：`guest`

**监控指标：**
- 队列长度
- 消息发布/消费速率
- 连接数
- 消费者数量

---

## Q7: 多个工作服务如何避免重复处理？

**A: 使用工作队列模式（Work Queue）**

**配置：**
```yaml
rabbitmq:
  queue: "task_queue"  # 同一个队列
  prefetch_count: 1    # 每次只处理一个消息
```

**工作原理：**
- 多个工作服务订阅同一个队列
- RabbitMQ 自动分发消息（轮询）
- 每个消息只被一个服务处理

---

## Q8: 如何实现事件广播？

**A: 使用 Fanout Exchange**

**配置：**
```yaml
rabbitmq:
  exchange_type: "fanout"  # 广播模式
  routing_key: ""          # Fanout 不需要路由键
```

**工作原理：**
- 发布到 Fanout Exchange 的消息会广播到所有绑定的队列
- 每个服务有自己的队列
- 所有服务都会收到消息

---

## Q9: 如何实现优先级队列？

**A: 使用队列参数和消息优先级**

**配置：**
```yaml
rabbitmq:
  queue: "priority_queue"
  # 在代码中设置队列参数
```

**代码示例：**
```go
args := amqp.Table{
    "x-max-priority": 10,  // 最大优先级
}
channel.QueueDeclare("priority_queue", true, false, false, false, args)

// 发布消息时设置优先级
amqp.Publishing{
    Priority: 5,  // 优先级 0-10
    Body: body,
}
```

---

## Q10: 如何实现延迟消息？

**A: 使用 RabbitMQ Delayed Message Plugin**

**安装插件后配置：**
```yaml
rabbitmq:
  exchange_type: "x-delayed-message"
  # 在消息中设置延迟时间
```

**代码示例：**
```go
headers := amqp.Table{
    "x-delay": 5000,  // 延迟 5 秒
}
amqp.Publishing{
    Headers: headers,
    Body: body,
}
```

---

## 总结

- ✅ 所有服务连接到同一个 RabbitMQ
- ✅ 工作服务可以发布事件
- ✅ 通过配置控制服务角色
- ✅ 使用命名规范防止事件循环
- ✅ 使用持久化确保消息不丢失
