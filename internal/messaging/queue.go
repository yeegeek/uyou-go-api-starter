// Package messaging 提供消息队列抽象接口
package messaging

import (
	"context"
)

// MessageHandler 消息处理器函数类型
type MessageHandler func(ctx context.Context, message []byte) error

// MessageQueue 消息队列接口，支持多种实现（RabbitMQ、AWS SNS+SQS 等）
type MessageQueue interface {
	// Publish 发布消息到指定主题/路由键
	// topic: 主题或路由键
	// message: 消息内容（通常是 JSON 序列化的字节数组）
	Publish(ctx context.Context, topic string, message []byte) error

	// Subscribe 订阅消息
	// topic: 主题或路由键（可选，某些实现可能忽略）
	// handler: 消息处理函数
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error

	// Close 关闭消息队列连接
	Close() error

	// HealthCheck 检查消息队列连接健康状态
	HealthCheck() error
}
