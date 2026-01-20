// Package messaging 提供消息队列功能
package messaging

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/uyou/uyou-go-api-starter/internal/config"
)

// RabbitMQ RabbitMQ 客户端
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  *config.RabbitMQConfig
}

// NewRabbitMQ 创建 RabbitMQ 客户端
func NewRabbitMQ(cfg *config.RabbitMQConfig) (*RabbitMQ, error) {
	// 连接到 RabbitMQ
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// 创建通道
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// 声明交换机
	err = channel.ExchangeDeclare(
		cfg.Exchange,     // 交换机名称
		cfg.ExchangeType, // 交换机类型
		true,             // 持久化
		false,            // 自动删除
		false,            // 内部交换机
		false,            // 不等待
		nil,              // 参数
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// 声明队列
	_, err = channel.QueueDeclare(
		cfg.Queue, // 队列名称
		true,      // 持久化
		false,     // 自动删除
		false,     // 独占
		false,     // 不等待
		nil,       // 参数
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// 绑定队列到交换机
	err = channel.QueueBind(
		cfg.Queue,      // 队列名称
		cfg.RoutingKey, // 路由键
		cfg.Exchange,   // 交换机名称
		false,          // 不等待
		nil,            // 参数
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	// 设置 QoS（服务质量）
	err = channel.Qos(
		cfg.PrefetchCount, // 预取数量
		0,                 // 预取大小
		false,             // 全局
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	slog.Info("Connected to RabbitMQ", "exchange", cfg.Exchange, "queue", cfg.Queue)

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
		config:  cfg,
	}, nil
}

// Close 关闭连接
func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			return err
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Publish 发布消息
func (r *RabbitMQ) Publish(ctx context.Context, routingKey string, body []byte) error {
	return r.channel.PublishWithContext(
		ctx,
		r.config.Exchange, // 交换机
		routingKey,        // 路由键
		false,             // 强制
		false,             // 立即
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // 持久化消息
			Timestamp:    time.Now(),
		},
	)
}

// Consume 消费消息
func (r *RabbitMQ) Consume(ctx context.Context, handler func([]byte) error) error {
	// 开始消费
	msgs, err := r.channel.Consume(
		r.config.Queue, // 队列名称
		"",             // 消费者标签
		false,          // 自动确认
		false,          // 独占
		false,          // 不等待
		false,          // 无本地
		nil,            // 参数
	)
	if err != nil {
		return fmt.Errorf("failed to consume messages: %w", err)
	}

	slog.Info("Started consuming messages", "queue", r.config.Queue)

	// 处理消息
	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping message consumer")
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("message channel closed")
			}

			// 处理消息
			if err := handler(msg.Body); err != nil {
				slog.Error("Failed to handle message", "error", err)
				// 拒绝消息并重新入队
				msg.Nack(false, true)
			} else {
				// 确认消息
				msg.Ack(false)
			}
		}
	}
}

// HealthCheck 检查连接健康状态
func (r *RabbitMQ) HealthCheck() error {
	if r.conn == nil || r.conn.IsClosed() {
		return fmt.Errorf("connection is closed")
	}
	if r.channel == nil {
		return fmt.Errorf("channel is nil")
	}
	return nil
}
