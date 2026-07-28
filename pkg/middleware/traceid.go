// Package middleware TraceID 中间件。
//
// 每个请求生成一个唯一 TraceID：
//   - 优先使用上游传入的 X-Trace-Id 头
//   - 否则生成一个 UUID-like ID
//
// 注入到：
//   - gin context（c.Set("trace_id", ...)）
//   - 响应头（X-Trace-Id）
//   - 日志（自动带上）
package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

// TraceID 生成或透传 TraceID，写入上下文和响应头
func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Trace-Id")
		if id == "" {
			id = newTraceID()
		}
		c.Set("trace_id", id)
		c.Writer.Header().Set("X-Trace-Id", id)
		c.Next()
	}
}

func newTraceID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}