// Package middleware 限流中间件。
//
// 用法：
//
//	r.Use(middleware.RateLimit(100, time.Minute))  // 每 IP 每分钟 100 次
//
// 实现：基于内存的滑动窗口（sync.Map + slice 时间戳）。
// 适用单机；分布式请用 Redis 令牌桶（见 pkg/middleware/ratelimit_redis.go）。
package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type window struct {
	mu     sync.Mutex
	stamps []time.Time
}

// RateLimit 基于 IP 的滑动窗口限流。
//   - limit：窗口内最大请求数
//   - win：窗口时长
//
// 超出限制返回 429。
func RateLimit(limit int, win time.Duration) gin.HandlerFunc {
	var (
		mu sync.RWMutex
		w  = make(map[string]*window)
	)

	// 定期清理过期条目（每 5 分钟）
	go func() {
		t := time.NewTicker(5 * time.Minute)
		defer t.Stop()
		for range t.C {
			cutoff := time.Now().Add(-win)
			mu.Lock()
			for k, v := range w {
				v.mu.Lock()
				if len(v.stamps) == 0 || v.stamps[len(v.stamps)-1].Before(cutoff) {
					delete(w, k)
				}
				v.mu.Unlock()
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		cutoff := now.Add(-win)

		mu.RLock()
		win, ok := w[ip]
		mu.RUnlock()
		if !ok {
			win = &window{}
			mu.Lock()
			win, ok = w[ip]
			if !ok {
				w[ip] = win
			}
			mu.Unlock()
		}

		win.mu.Lock()
		// 滑出窗口外的时间戳
		i := 0
		for ; i < len(win.stamps); i++ {
			if win.stamps[i].After(cutoff) {
				break
			}
		}
		win.stamps = win.stamps[i:]
		if len(win.stamps) >= limit {
			win.mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求过于频繁",
			})
			return
		}
		win.stamps = append(win.stamps, now)
		win.mu.Unlock()

		c.Next()
	}
}