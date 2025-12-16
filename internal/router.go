package internal

import (
	"fing/internal/apis/login"
	"fing/pkg/health"
	"fing/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// InitRouter 路由示例
func InitRouter(r *gin.Engine) *gin.Engine {
	// 公共路由，不需要认证
	publicRouter(r)

	// 中间件, 顺序不能改
	r.Use(middleware.LoggerToFile(), middleware.Session(), middleware.Cover, middleware.Cors(), middleware.CurrentUser())

	// 日常业务路由
	normalRouter(r)

	// 可以按照当前业务路由自定义自己的分组
	// authRouter(r)

	return r
}

func publicRouter(r *gin.Engine) {
	// 健康检查等公共接口
	r.GET("/health", health.HealthCheck)
}

func normalRouter(r *gin.Engine) {

	new(login.RouterLogin).Router(r)

}
