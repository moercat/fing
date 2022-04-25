package tools

import (
	"fing/pkg/entity/usr"
	"github.com/gin-gonic/gin"
)

// CurrentUser 获取当前用户
func CurrentUser(c *gin.Context) *usr.UserInfo {
	user, exist := c.Get("user")
	if !exist {
		return nil
	}

	if u, ok := user.(*usr.UserInfo); ok {
		return u
	}

	return nil
}
