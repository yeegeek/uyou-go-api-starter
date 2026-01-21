// Package messaging 提供消息队列工厂函数
package messaging

import (
	"fmt"

	"github.com/yeegeek/uyou-go-api-starter/internal/config"
)

// MessageQueueProvider 消息队列提供商类型
type MessageQueueProvider string

const (
	// ProviderRabbitMQ RabbitMQ 提供商
	ProviderRabbitMQ MessageQueueProvider = "rabbitmq"
	// ProviderAWSSNS AWS SNS+SQS 提供商（未来实现）
	ProviderAWSSNS MessageQueueProvider = "aws-sns"
	// ProviderGCPPubSub Google Cloud Pub/Sub 提供商（未来实现）
	ProviderGCPPubSub MessageQueueProvider = "gcp-pubsub"
)

// NewMessageQueue 根据配置创建消息队列实例
// 根据配置中的 provider 字段选择实现
func NewMessageQueue(cfg *config.RabbitMQConfig, provider MessageQueueProvider) (MessageQueue, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("message queue is not enabled")
	}

	switch provider {
	case ProviderRabbitMQ:
		return NewRabbitMQ(cfg)
	case ProviderAWSSNS:
		// TODO: 实现 AWS SNS+SQS
		return nil, fmt.Errorf("AWS SNS provider not yet implemented")
	case ProviderGCPPubSub:
		// TODO: 实现 GCP Pub/Sub
		return nil, fmt.Errorf("GCP Pub/Sub provider not yet implemented")
	default:
		// 默认使用 RabbitMQ
		return NewRabbitMQ(cfg)
	}
}

// NewMessageQueueFromConfig 从完整配置创建消息队列实例
// 根据配置中的 provider 字段选择实现（如果配置中有的话）
func NewMessageQueueFromConfig(cfg *config.Config) (MessageQueue, error) {
	if !cfg.RabbitMQ.Enabled {
		return nil, fmt.Errorf("message queue is not enabled")
	}

	// 默认使用 RabbitMQ，未来可以从配置中读取 provider
	// provider := MessageQueueProvider(cfg.MessageQueue.Provider) // 如果配置中有
	provider := ProviderRabbitMQ

	return NewMessageQueue(&cfg.RabbitMQ, provider)
}
