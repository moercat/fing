// Package notify 提供站内通知 API。
package notify

import (
	"fing/pkg/middleware"
	"fing/pkg/notify"
	"fing/pkg/resp"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RouterNotify struct{}

func (r *RouterNotify) Router(router *gin.Engine) {
	g := router.Group("/api/v2/notify")
	g.Use(middleware.AuthRequired())
	{
		g.GET("", r.list)
		g.GET("unread", r.unreadCount)
		g.POST("read", r.markAllRead)
	}
}

func (r *RouterNotify) list(c *gin.Context) {
	uid, _ := c.Get("user_id")
	n, _ := strconv.Atoi(c.Query("n"))
	if n == 0 {
		n = 20
	}
	resp.OK(c, notify.List(uid.(uint), n), "")
}

func (r *RouterNotify) unreadCount(c *gin.Context) {
	uid, _ := c.Get("user_id")
	resp.OK(c, gin.H{"count": notify.UnreadCount(uid.(uint))}, "")
}

func (r *RouterNotify) markAllRead(c *gin.Context) {
	uid, _ := c.Get("user_id")
	notify.MarkAllRead(uid.(uint))
	resp.OK(c, nil, "已全部已读")
}