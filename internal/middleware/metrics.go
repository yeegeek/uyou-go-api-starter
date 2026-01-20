// Package middleware 提供 Prometheus 监控中间件
package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uyou/uyou-go-api-starter/internal/metrics"
)

// MetricsMiddleware Prometheus 监控中间件
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		start := time.Now()

		// 获取请求大小
		requestSize := float64(c.Request.ContentLength)

		// 处理请求
		c.Next()

		// 计算请求延迟
		duration := time.Since(start).Seconds()

		// 获取响应大小
		responseSize := float64(c.Writer.Size())

		// 获取状态码
		status := strconv.Itoa(c.Writer.Status())

		// 记录指标
		metrics.RecordHTTPRequest(
			c.Request.Method,
			c.FullPath(),
			status,
			duration,
			requestSize,
			responseSize,
		)
	}
}
