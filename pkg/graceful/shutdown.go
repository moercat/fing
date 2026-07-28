// Package graceful 提供 HTTP 服务优雅关闭。
//
// 监听 SIGINT / SIGTERM 信号，到期后：
//   1. 停止接收新请求（http.Server.Shutdown）
//   2. 等待正在处理的请求完成（最多 timeout）
//   3. 执行注册的 cleanup 钩子（如关闭 DB / Redis 连接）
//
// 用法：
//
//	graceful.Run(graceful.Config{
//	    Addr:    ":9765",
//	    Handler: r,
//	    Timeout: 30 * time.Second,
//	    Cleanups: []func(){
//	        closeDB,
//	        closeRedis,
//	    },
//	})
package graceful

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Config 优雅关闭配置
type Config struct {
	Addr     string          // 监听地址
	Handler  http.Handler    // HTTP 处理器
	Timeout  time.Duration   // 关闭超时（默认 30s）
	Cleanups []func()        // 关闭前执行的清理钩子
}

// Run 启动服务，阻塞直到收到关闭信号
func Run(cfg Config) {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: cfg.Handler,
	}

	// 监听关闭信号
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// 启动服务
	go func() {
		log.Printf("[graceful] listening on %s", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[graceful] server error: %v", err)
		}
	}()

	// 等待信号
	<-stop
	log.Printf("[graceful] shutdown signal received, waiting %s ...", cfg.Timeout)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	// 关闭 HTTP 服务（等待正在处理的请求）
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[graceful] http shutdown error: %v", err)
	}

	// 执行清理钩子
	for i, fn := range cfg.Cleanups {
		log.Printf("[graceful] running cleanup #%d", i+1)
		fn()
	}

	log.Println("[graceful] bye")
}