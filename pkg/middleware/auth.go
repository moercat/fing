package middleware

import (
	"fing/internal/service/login"
	"fing/pkg/entity/usr"
	"fing/pkg/resp"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// CurrentUser 获取登录用户
func CurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		uid := session.Get("user_id")
		if uid != nil {
			user, err := login.GetUser(uid)
			if err == nil {
				c.Set("user", &user)
			}
		}
		c.Next()
	}
}

// AuthRequired 需要登录
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get("user"); user != nil {
			if _, ok := user.(*usr.UserInfo); ok {
				c.Next()
				return
			}
		}

		resp.Fail(c, nil, "未登录", 401)
		c.Abort()
	}
}
