// Package middleware 提供增强的结构化日志中间件
package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// EnhancedLoggerMiddleware 增强的结构化日志中间件
func EnhancedLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求 ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// 生成追踪 ID
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Set("trace_id", traceID)
		c.Header("X-Trace-ID", traceID)

		// 记录开始时间
		start := time.Now()

		// 记录请求开始
		slog.Info("Request started",
			"request_id", requestID,
			"trace_id", traceID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)

		// 处理请求
		c.Next()

		// 计算请求延迟
		duration := time.Since(start)

		// 获取用户信息（如果已认证）
		userID, _ := c.Get("user_id")
		userEmail, _ := c.Get("user_email")

		// 记录请求完成
		logLevel := slog.LevelInfo
		if c.Writer.Status() >= 500 {
			logLevel = slog.LevelError
		} else if c.Writer.Status() >= 400 {
			logLevel = slog.LevelWarn
		}

		slog.Log(c.Request.Context(), logLevel, "Request completed",
			"request_id", requestID,
			"trace_id", traceID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
			"response_size", c.Writer.Size(),
			"user_id", userID,
			"user_email", userEmail,
			"errors", c.Errors.String(),
		)
	}
}

// GetRequestID 从上下文获取请求 ID
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return ""
}

// GetTraceID 从上下文获取追踪 ID
func GetTraceID(c *gin.Context) string {
	if traceID, exists := c.Get("trace_id"); exists {
		return traceID.(string)
	}
	return ""
}
