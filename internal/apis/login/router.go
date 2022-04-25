package login

import (
	"fing/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type RouterLogin struct{}

func (r *RouterLogin) Router(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("ping", r.ping)
		v1.POST("register", r.register)
		v1.POST("login", r.login)
	}

	v2 := router.Group("/api/v2")
	//需要登录认证
	v2.Use(middleware.AuthRequired())
	{
		v2.GET("user_info", r.userInfo)
		v2.DELETE("logout", r.logout)
	}

}
