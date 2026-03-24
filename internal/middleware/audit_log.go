package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// responseBodyWriter is a custom response writer to capture the response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// AuditLog 零信任访问审计日志中间件，用于记录敏感数据访问（特别是客户档案和认证信息）
func AuditLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Read request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody)) // Restore the body
		}

		// Use custom response writer to capture response
		w := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = w

		// Process request
		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)

		userID := c.GetUint64("user_id")

		// 我们主要审计涉及敏感数据的接口，如客户端详情、认证信息
		path := c.Request.URL.Path
		isSensitive := false
		if path == "/api/clients/detail" || path == "/api/clients/auth" || path == "/api/clients/list" {
			isSensitive = true
		}

		if isSensitive {
			// 在生产环境中，可以将日志发送到专门的审计存储如 SLS、Elasticsearch，或单独的审计表
			logrus.WithFields(logrus.Fields{
				"event":        "audit_log",
				"user_id":      userID,
				"client_ip":    c.ClientIP(),
				"method":       c.Request.Method,
				"path":         path,
				"status_code":  c.Writer.Status(),
				"latency_ms":   latency.Milliseconds(),
				"request_body": string(requestBody),
				// "response_body": w.body.String(), // 如果响应体过大，可选择截断或不记录
			}).Info("Sensitive Data Access")
		}
	}
}
