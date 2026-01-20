// Package messaging 提供事件发布订阅功能
package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

// Event 事件接口
type Event interface {
	EventType() string
	EventData() interface{}
}

// BaseEvent 基础事件
type BaseEvent struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	TraceID   string      `json:"trace_id,omitempty"`
}

// EventType 返回事件类型
func (e *BaseEvent) EventType() string {
	return e.Type
}

// EventData 返回事件数据
func (e *BaseEvent) EventData() interface{} {
	return e.Data
}

// EventBus 事件总线
type EventBus struct {
	mq *RabbitMQ
}

// NewEventBus 创建事件总线
func NewEventBus(mq *RabbitMQ) *EventBus {
	return &EventBus{
		mq: mq,
	}
}

// Publish 发布事件
func (eb *EventBus) Publish(ctx context.Context, event Event) error {
	// 构建事件消息
	eventMsg := &BaseEvent{
		Type:      event.EventType(),
		Data:      event.EventData(),
		Timestamp: time.Now(),
	}

	// 从上下文中获取 TraceID（如果存在）
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		eventMsg.TraceID = traceID
	}

	// 序列化事件
	body, err := json.Marshal(eventMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// 发布到消息队列
	routingKey := fmt.Sprintf("event.%s", event.EventType())
	if err := eb.mq.Publish(ctx, routingKey, body); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	slog.Info("Event published", "type", event.EventType(), "routing_key", routingKey)
	return nil
}

// Subscribe 订阅事件
func (eb *EventBus) Subscribe(ctx context.Context, handler EventHandler) error {
	return eb.mq.Consume(ctx, func(body []byte) error {
		// 反序列化事件
		var event BaseEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return fmt.Errorf("failed to unmarshal event: %w", err)
		}

		// 处理事件
		slog.Info("Event received", "type", event.Type, "timestamp", event.Timestamp)
		return handler.Handle(ctx, &event)
	})
}

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(ctx context.Context, event *BaseEvent) error
}

// EventHandlerFunc 事件处理器函数类型
type EventHandlerFunc func(ctx context.Context, event *BaseEvent) error

// Handle 实现 EventHandler 接口
func (f EventHandlerFunc) Handle(ctx context.Context, event *BaseEvent) error {
	return f(ctx, event)
}

// 预定义的事件类型常量
const (
	EventTypeUserCreated = "user.created"
	EventTypeUserUpdated = "user.updated"
	EventTypeUserDeleted = "user.deleted"
)

// UserCreatedEvent 用户创建事件
type UserCreatedEvent struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

// EventType 返回事件类型
func (e *UserCreatedEvent) EventType() string {
	return EventTypeUserCreated
}

// EventData 返回事件数据
func (e *UserCreatedEvent) EventData() interface{} {
	return e
}

// UserUpdatedEvent 用户更新事件
type UserUpdatedEvent struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

// EventType 返回事件类型
func (e *UserUpdatedEvent) EventType() string {
	return EventTypeUserUpdated
}

// EventData 返回事件数据
func (e *UserUpdatedEvent) EventData() interface{} {
	return e
}

// UserDeletedEvent 用户删除事件
type UserDeletedEvent struct {
	UserID uint `json:"user_id"`
}

// EventType 返回事件类型
func (e *UserDeletedEvent) EventType() string {
	return EventTypeUserDeleted
}

// EventData 返回事件数据
func (e *UserDeletedEvent) EventData() interface{} {
	return e
}
