package internal

import (
	"fing/internal/apis/login"
	"fing/internal/apis/notify"
	"fing/internal/apis/password"
	"fing/internal/apis/upload"
	"fing/internal/apis/user"
	"fing/pkg/health"
	"fing/pkg/middleware"
	"fing/pkg/swagger"
	"github.com/gin-gonic/gin"
)

// InitRouter 路由示例
func InitRouter(r *gin.Engine) *gin.Engine {
	// 公共路由，不需要认证
	publicRouter(r)

	// 全局中间件（顺序重要）
	r.Use(
		middleware.TraceID(),                 // 1. 每个请求分配 TraceID
		middleware.RateLimit(100, 1<<63-1),  // 2. IP 限流（默认关闭，main.go 配置）
		middleware.LoggerToFile(),           // 3. 请求日志
		middleware.Session(),                // 4. Session
		middleware.Cover,                    // 5. Panic 恢复
		middleware.Cors(),                   // 6. CORS
		middleware.Audit(),                  // 7. 操作日志（审计）
		middleware.CurrentUser(),            // 8. 当前用户
	)

	// 日常业务路由
	normalRouter(r)

	return r
}

func publicRouter(r *gin.Engine) {
	// 健康检查等公共接口
	r.GET("/health", health.HealthCheck)

	// Swagger UI
	swagger.Register(r)
}

func normalRouter(r *gin.Engine) {
	new(login.RouterLogin).Router(r)
	new(user.RouterUser).Router(r)
	new(password.RouterPassword).Router(r)
	new(upload.RouterUpload).Router(r)
	new(notify.RouterNotify).Router(r)
}
