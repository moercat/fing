// Package middleware 操作日志（审计日志）中间件。
//
// 自动记录每个请求：
//   - 方法 / 路径 / IP / 用户 ID
//   - 状态码 / 耗时
//   - TraceID
//
// 写入 pkg/audit 包（默认内存环形缓冲，可扩展到文件 / MySQL / ES）。
package middleware

import (
	"bytes"
	"io"
	"time"

	"fing/pkg/audit"
	"github.com/gin-gonic/gin"
)

// Audit 记录每个请求的审计日志
func Audit() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 读 body（用于审计 POST 内容，可选）
		var bodyBytes []byte
		if c.Request.Body != nil && c.Request.Method != "GET" {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		uid, _ := c.Get("user_id")
		uidUint, _ := uid.(uint)

		audit.Log(audit.Entry{
			Time:     start,
			Method:   c.Request.Method,
			Path:     c.Request.URL.Path,
			IP:       c.ClientIP(),
			UserID:   uidUint,
			Status:   c.Writer.Status(),
			Latency:  time.Since(start),
			TraceID:  c.GetString("trace_id"),
			UA:       c.Request.UserAgent(),
			BodySize: len(bodyBytes),
		})
	}
}