// Package health 提供增强的健康检查功能
package health

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"

	"github.com/yeegeek/uyou-go-api-starter/internal/config"
)

// ComponentCheck 组件检查结果
type ComponentCheck struct {
	Status  HealthStatus `json:"status"`
	Message string       `json:"message,omitempty"`
	Latency string       `json:"latency,omitempty"` // 响应延迟
}

// EnhancedHealthResponse 增强的健康检查响应
type EnhancedHealthResponse struct {
	Status     HealthStatus              `json:"status"`
	Timestamp  time.Time                 `json:"timestamp"`
	Version    string                    `json:"version"`
	Components map[string]ComponentCheck `json:"components"`
}

// EnhancedHealthHandler 增强的健康检查处理器
type EnhancedHealthHandler struct {
	db          *gorm.DB
	redisClient *redis.Client
	mongoClient *mongo.Client
	cfg         *config.Config
}

// NewEnhancedHealthHandler 创建增强的健康检查处理器
func NewEnhancedHealthHandler(
	db *gorm.DB,
	redisClient *redis.Client,
	mongoClient *mongo.Client,
	cfg *config.Config,
) *EnhancedHealthHandler {
	return &EnhancedHealthHandler{
		db:          db,
		redisClient: redisClient,
		mongoClient: mongoClient,
		cfg:         cfg,
	}
}

// HandleEnhancedHealth 处理增强的健康检查请求
func (h *EnhancedHealthHandler) HandleEnhancedHealth(c *gin.Context) {
	ctx := c.Request.Context()
	response := EnhancedHealthResponse{
		Timestamp:  time.Now(),
		Version:    h.cfg.App.Version,
		Components: make(map[string]ComponentCheck),
	}

	// 检查 PostgreSQL
	if h.db != nil {
		response.Components["database"] = h.checkPostgreSQL(ctx)
	}

	// 检查 Redis（如果启用）
	if h.cfg.Redis.Enabled && h.redisClient != nil {
		response.Components["redis"] = h.checkRedis(ctx)
	}

	// 检查 MongoDB（如果启用）
	if h.cfg.MongoDB.Enabled && h.mongoClient != nil {
		response.Components["mongodb"] = h.checkMongoDB(ctx)
	}

	// 确定整体健康状态
	response.Status = h.determineOverallStatus(response.Components)

	// 根据健康状态返回相应的 HTTP 状态码
	statusCode := http.StatusOK
	if response.Status == StatusDegraded {
		statusCode = http.StatusOK // 降级但仍可用
	} else if response.Status == StatusUnhealthy {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// checkPostgreSQL 检查 PostgreSQL 连接
func (h *EnhancedHealthHandler) checkPostgreSQL(ctx context.Context) ComponentCheck {
	start := time.Now()

	sqlDB, err := h.db.DB()
	if err != nil {
		return ComponentCheck{
			Status:  StatusUnhealthy,
			Message: "Failed to get database instance: " + err.Error(),
		}
	}

	// 设置超时
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		return ComponentCheck{
			Status:  StatusUnhealthy,
			Message: "Database ping failed: " + err.Error(),
		}
	}

	// 检查连接池状态
	stats := sqlDB.Stats()
	if stats.OpenConnections >= stats.MaxOpenConnections {
		return ComponentCheck{
			Status:  StatusDegraded,
			Message: "Connection pool exhausted",
			Latency: time.Since(start).String(),
		}
	}

	return ComponentCheck{
		Status:  StatusHealthy,
		Message: "Connected",
		Latency: time.Since(start).String(),
	}
}

// checkRedis 检查 Redis 连接
func (h *EnhancedHealthHandler) checkRedis(ctx context.Context) ComponentCheck {
	start := time.Now()

	// 设置超时
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := h.redisClient.Ping(pingCtx).Err(); err != nil {
		return ComponentCheck{
			Status:  StatusUnhealthy,
			Message: "Redis ping failed: " + err.Error(),
		}
	}

	// 检查连接池状态
	stats := h.redisClient.PoolStats()
	if stats.TotalConns >= uint32(h.cfg.Redis.PoolSize) {
		return ComponentCheck{
			Status:  StatusDegraded,
			Message: "Connection pool near capacity",
			Latency: time.Since(start).String(),
		}
	}

	return ComponentCheck{
		Status:  StatusHealthy,
		Message: "Connected",
		Latency: time.Since(start).String(),
	}
}

// checkMongoDB 检查 MongoDB 连接
func (h *EnhancedHealthHandler) checkMongoDB(ctx context.Context) ComponentCheck {
	start := time.Now()

	// 设置超时
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := h.mongoClient.Ping(pingCtx, nil); err != nil {
		return ComponentCheck{
			Status:  StatusUnhealthy,
			Message: "MongoDB ping failed: " + err.Error(),
		}
	}

	return ComponentCheck{
		Status:  StatusHealthy,
		Message: "Connected",
		Latency: time.Since(start).String(),
	}
}

// determineOverallStatus 确定整体健康状态
func (h *EnhancedHealthHandler) determineOverallStatus(components map[string]ComponentCheck) HealthStatus {
	hasUnhealthy := false
	hasDegraded := false

	for _, check := range components {
		if check.Status == StatusUnhealthy {
			hasUnhealthy = true
		} else if check.Status == StatusDegraded {
			hasDegraded = true
		}
	}

	if hasUnhealthy {
		return StatusUnhealthy
	}
	if hasDegraded {
		return StatusDegraded
	}
	return StatusHealthy
}

// HandleLiveness 处理存活探针（Kubernetes liveness probe）
// 仅检查应用是否运行，不检查依赖服务
func (h *EnhancedHealthHandler) HandleLiveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": time.Now(),
	})
}

// HandleReadiness 处理就绪探针（Kubernetes readiness probe）
// 检查应用是否准备好接收流量
func (h *EnhancedHealthHandler) HandleReadiness(c *gin.Context) {
	ctx := c.Request.Context()

	// 检查关键依赖（数据库）
	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil || sqlDB.Ping() != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":    "not ready",
				"timestamp": time.Now(),
				"reason":    "database not available",
			})
			return
		}
	}

	// 检查 Redis（如果启用且关键）
	if h.cfg.Redis.Enabled && h.redisClient != nil {
		if err := h.redisClient.Ping(ctx).Err(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":    "not ready",
				"timestamp": time.Now(),
				"reason":    "redis not available",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"timestamp": time.Now(),
	})
}

// GetDatabaseStats 获取数据库连接池统计信息
func (h *EnhancedHealthHandler) GetDatabaseStats() (*sql.DBStats, error) {
	if h.db == nil {
		return nil, nil
	}

	sqlDB, err := h.db.DB()
	if err != nil {
		return nil, err
	}

	stats := sqlDB.Stats()
	return &stats, nil
}
