// Package user 提供用户资料修改、密码修改等 API。
package user

import (
	"fing/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type RouterUser struct{}

// Router 注册用户相关路由
func (r *RouterUser) Router(router *gin.Engine) {
	v2 := router.Group("/api/v2")
	v2.Use(middleware.AuthRequired())
	{
		// 资料
		v2.GET("profile", r.getProfile)
		v2.PUT("profile", r.updateProfile)
		v2.PUT("profile/password", r.changePassword)
		v2.PUT("profile/avatar", r.updateAvatar)

		// 角色管理（仅 admin）
		admin := v2.Group("")
		admin.Use(middleware.AdminRequired())
		{
			admin.GET("users", r.listUsers)
			admin.PUT("users/:id/role", r.updateUserRole)
		}
	}
}
