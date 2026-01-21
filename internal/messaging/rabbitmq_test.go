package messaging

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yeegeek/uyou-go-api-starter/internal/config"
)

func TestNewRabbitMQ(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.RabbitMQConfig
		wantErr bool
		errMsg  string
	}{
		{
			// 注意：本测试在“RabbitMQ 未运行”的机器上用于验证失败分支；
			// 如果本机恰好有 RabbitMQ 运行，则会返回成功，这种情况下我们跳过该用例。
			name: "valid config but no rabbitmq server (expected failure if rabbitmq is down)",
			cfg: &config.RabbitMQConfig{
				URL:          "amqp://guest:guest@localhost:5672/",
				Exchange:     "test_exchange",
				ExchangeType: "topic",
				Queue:        "test_queue",
				RoutingKey:   "test.routing.key",
				PrefetchCount: 10,
			},
			wantErr: true,
			errMsg:  "failed to connect to RabbitMQ",
		},
		{
			name: "invalid URL",
			cfg: &config.RabbitMQConfig{
				URL:          "invalid-url",
				Exchange:     "test_exchange",
				ExchangeType: "topic",
				Queue:        "test_queue",
				RoutingKey:   "test.routing.key",
				PrefetchCount: 10,
			},
			wantErr: true,
			errMsg:  "failed to connect to RabbitMQ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mq, err := NewRabbitMQ(tt.cfg)
			if tt.wantErr {
				if err == nil && mq != nil {
					// 本机 RabbitMQ 可用时，失败分支不可复现，跳过即可
					_ = mq.Close()
					t.Skip("RabbitMQ is available on localhost; skipping failure-path test case")
				}
				assert.Error(t, err)
				assert.Nil(t, mq)
				if err != nil {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, mq)
				if mq != nil {
					defer mq.Close()
				}
			}
		})
	}
}

func TestRabbitMQ_HealthCheck(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		Exchange:     "test_exchange",
		ExchangeType: "topic",
		Queue:        "test_queue",
		RoutingKey:   "test.routing.key",
		PrefetchCount: 10,
	}

	mq, err := NewRabbitMQ(cfg)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer mq.Close()

	err = mq.HealthCheck()
	assert.NoError(t, err)
}

func TestRabbitMQ_Publish(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		Exchange:     "test_exchange",
		ExchangeType: "topic",
		Queue:        "test_queue",
		RoutingKey:   "test.routing.key",
		PrefetchCount: 10,
	}

	mq, err := NewRabbitMQ(cfg)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer mq.Close()

	ctx := context.Background()
	body := []byte("test message")

	err = mq.Publish(ctx, "test.routing.key", body)
	assert.NoError(t, err)
}

func TestRabbitMQ_ImplementsMessageQueue(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		Exchange:     "test_exchange",
		ExchangeType: "topic",
		Queue:        "test_queue",
		RoutingKey:   "test.routing.key",
		PrefetchCount: 10,
	}

	mq, err := NewRabbitMQ(cfg)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer mq.Close()

	// 验证 RabbitMQ 实现了 MessageQueue 接口
	var _ MessageQueue = mq
	
	// 测试接口方法
	ctx := context.Background()
	
	// 测试 Publish（接口方法）
	err = mq.Publish(ctx, "test.topic", []byte("test"))
	assert.NoError(t, err)
	
	// 测试 HealthCheck（接口方法）
	err = mq.HealthCheck()
	assert.NoError(t, err)
	
	// 测试 Close（接口方法）
	err = mq.Close()
	assert.NoError(t, err)
}

func TestRabbitMQ_Close(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		Exchange:     "test_exchange",
		ExchangeType: "topic",
		Queue:        "test_queue",
		RoutingKey:   "test.routing.key",
		PrefetchCount: 10,
	}

	mq, err := NewRabbitMQ(cfg)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}

	err = mq.Close()
	assert.NoError(t, err)
}

func TestNewEventBus(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		Exchange:     "test_exchange",
		ExchangeType: "topic",
		Queue:        "test_queue",
		RoutingKey:   "test.routing.key",
		PrefetchCount: 10,
	}

	mq, err := NewRabbitMQ(cfg)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer mq.Close()

	eventBus := NewEventBus(mq)
	assert.NotNil(t, eventBus)
	assert.NotNil(t, eventBus.mq)
	
	// 验证 mq 实现了 MessageQueue 接口
	var _ MessageQueue = mq
}

func TestEventBus_Publish(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		Exchange:     "test_exchange",
		ExchangeType: "topic",
		Queue:        "test_queue",
		RoutingKey:   "test.routing.key",
		PrefetchCount: 10,
	}

	mq, err := NewRabbitMQ(cfg)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer mq.Close()

	eventBus := NewEventBus(mq)
	ctx := context.Background()

	event := &UserCreatedEvent{
		UserID: 1,
		Email:   "test@example.com",
		Name:    "Test User",
	}

	err = eventBus.Publish(ctx, event)
	assert.NoError(t, err)
}

func TestEventBus_Publish_WithTraceID(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		Exchange:     "test_exchange",
		ExchangeType: "topic",
		Queue:        "test_queue",
		RoutingKey:   "test.routing.key",
		PrefetchCount: 10,
	}

	mq, err := NewRabbitMQ(cfg)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer mq.Close()

	eventBus := NewEventBus(mq)
	ctx := context.WithValue(context.Background(), "trace_id", "test-trace-id")

	event := &UserCreatedEvent{
		UserID: 1,
		Email:   "test@example.com",
		Name:    "Test User",
	}

	err = eventBus.Publish(ctx, event)
	assert.NoError(t, err)
}

func TestUserCreatedEvent(t *testing.T) {
	event := &UserCreatedEvent{
		UserID: 1,
		Email:  "test@example.com",
		Name:   "Test User",
	}

	assert.Equal(t, EventTypeUserCreated, event.EventType())
	assert.Equal(t, event, event.EventData())
}

func TestUserUpdatedEvent(t *testing.T) {
	event := &UserUpdatedEvent{
		UserID: 1,
		Email:  "test@example.com",
		Name:   "Test User",
	}

	assert.Equal(t, EventTypeUserUpdated, event.EventType())
	assert.Equal(t, event, event.EventData())
}

func TestUserDeletedEvent(t *testing.T) {
	event := &UserDeletedEvent{
		UserID: 1,
	}

	assert.Equal(t, EventTypeUserDeleted, event.EventType())
	assert.Equal(t, event, event.EventData())
}

func TestEventHandlerFunc(t *testing.T) {
	var called bool
	handler := EventHandlerFunc(func(ctx context.Context, event *BaseEvent) error {
		called = true
		return nil
	})

	ctx := context.Background()
	event := &BaseEvent{
		Type:      "test.event",
		Data:      map[string]interface{}{"key": "value"},
		Timestamp: time.Now(),
	}

	err := handler.Handle(ctx, event)
	assert.NoError(t, err)
	assert.True(t, called)
}
