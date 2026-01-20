// Package metrics 提供 Prometheus 监控指标
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTPRequestsTotal HTTP 请求总数
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "HTTP 请求总数",
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestDuration HTTP 请求延迟（秒）
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP 请求延迟（秒）",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// HTTPRequestSize HTTP 请求大小（字节）
	HTTPRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP 请求大小（字节）",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7),
		},
		[]string{"method", "path"},
	)

	// HTTPResponseSize HTTP 响应大小（字节）
	HTTPResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP 响应大小（字节）",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7),
		},
		[]string{"method", "path"},
	)

	// DatabaseQueriesTotal 数据库查询总数
	DatabaseQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_queries_total",
			Help: "数据库查询总数",
		},
		[]string{"operation", "table"},
	)

	// DatabaseQueryDuration 数据库查询延迟（秒）
	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "数据库查询延迟（秒）",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "table"},
	)

	// CacheHitsTotal 缓存命中总数
	CacheHitsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "缓存命中总数",
		},
		[]string{"cache_name"},
	)

	// CacheMissesTotal 缓存未命中总数
	CacheMissesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "缓存未命中总数",
		},
		[]string{"cache_name"},
	)

	// GRPCRequestsTotal gRPC 请求总数
	GRPCRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "gRPC 请求总数",
		},
		[]string{"method", "status"},
	)

	// GRPCRequestDuration gRPC 请求延迟（秒）
	GRPCRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "gRPC 请求延迟（秒）",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	// MessageQueuePublishedTotal 消息队列发布总数
	MessageQueuePublishedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "message_queue_published_total",
			Help: "消息队列发布总数",
		},
		[]string{"event_type"},
	)

	// MessageQueueConsumedTotal 消息队列消费总数
	MessageQueueConsumedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "message_queue_consumed_total",
			Help: "消息队列消费总数",
		},
		[]string{"event_type", "status"},
	)

	// ActiveConnections 活跃连接数
	ActiveConnections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "活跃连接数",
		},
		[]string{"type"},
	)

	// ErrorsTotal 错误总数
	ErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "错误总数",
		},
		[]string{"type", "code"},
	)
)

// RecordHTTPRequest 记录 HTTP 请求指标
func RecordHTTPRequest(method, path, status string, duration, requestSize, responseSize float64) {
	HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
	HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
	HTTPRequestSize.WithLabelValues(method, path).Observe(requestSize)
	HTTPResponseSize.WithLabelValues(method, path).Observe(responseSize)
}

// RecordDatabaseQuery 记录数据库查询指标
func RecordDatabaseQuery(operation, table string, duration float64) {
	DatabaseQueriesTotal.WithLabelValues(operation, table).Inc()
	DatabaseQueryDuration.WithLabelValues(operation, table).Observe(duration)
}

// RecordCacheHit 记录缓存命中
func RecordCacheHit(cacheName string) {
	CacheHitsTotal.WithLabelValues(cacheName).Inc()
}

// RecordCacheMiss 记录缓存未命中
func RecordCacheMiss(cacheName string) {
	CacheMissesTotal.WithLabelValues(cacheName).Inc()
}

// RecordGRPCRequest 记录 gRPC 请求指标
func RecordGRPCRequest(method, status string, duration float64) {
	GRPCRequestsTotal.WithLabelValues(method, status).Inc()
	GRPCRequestDuration.WithLabelValues(method).Observe(duration)
}

// RecordMessagePublished 记录消息发布
func RecordMessagePublished(eventType string) {
	MessageQueuePublishedTotal.WithLabelValues(eventType).Inc()
}

// RecordMessageConsumed 记录消息消费
func RecordMessageConsumed(eventType, status string) {
	MessageQueueConsumedTotal.WithLabelValues(eventType, status).Inc()
}

// SetActiveConnections 设置活跃连接数
func SetActiveConnections(connType string, count float64) {
	ActiveConnections.WithLabelValues(connType).Set(count)
}

// RecordError 记录错误
func RecordError(errorType, code string) {
	ErrorsTotal.WithLabelValues(errorType, code).Inc()
}
