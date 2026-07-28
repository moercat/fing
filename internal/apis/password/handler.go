// Package password 提供密码重置与邮箱验证流程。
package password

import (
	"crypto/rand"
	"encoding/hex"
	"fing/internal/model"
	"fing/internal/service/password"
	"fing/pkg/email"
	"fing/pkg/resp"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type RouterPassword struct{}

// Router 注册公开的密码/邮箱相关路由
func (r *RouterPassword) Router(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("password/forgot", r.forgot)
		v1.POST("password/reset", r.reset)
		v1.POST("email/verify", r.verifyEmail)
	}
}

// forgot 申请重置密码
func (r *RouterPassword) forgot(c *gin.Context) {
	var sv model.ForgotPassword
	if err := c.ShouldBind(&sv); err != nil {
		resp.Fail(c, err, "参数错误")
		return
	}

	token, err := password.IssueResetToken(sv.Email)
	if err != nil {
		// 出于安全考虑，不暴露邮箱是否存在
		resp.OK(c, nil, "如果该邮箱已注册，重置链接已发送")
		return
	}

	// 发送邮件
	resetLink := fmt.Sprintf("https://your-domain.com/reset?token=%s", token)
	body := fmt.Sprintf("请在 30 分钟内点击以下链接重置密码：<br><a href=\"%s\">%s</a>", resetLink, resetLink)
	if err := email.SendMail(sv.Email, "User", "密码重置", body); err != nil {
		resp.Fail(c, err, "邮件发送失败")
		return
	}
	resp.OK(c, nil, "重置链接已发送到您的邮箱")
}

// reset 用 token 重置密码
func (r *RouterPassword) reset(c *gin.Context) {
	var sv model.ResetPassword
	if err := c.ShouldBind(&sv); err != nil {
		resp.Fail(c, err, "参数错误")
		return
	}
	if err := password.ResetWithToken(sv.Token, sv.NewPassword); err != nil {
		resp.Fail(c, err, "重置失败")
		return
	}
	resp.OK(c, nil, "密码已重置")
}

// verifyEmail 邮箱验证
func (r *RouterPassword) verifyEmail(c *gin.Context) {
	var sv model.VerifyEmail
	if err := c.ShouldBind(&sv); err != nil {
		resp.Fail(c, err, "参数错误")
		return
	}
	if err := password.VerifyEmail(sv.Token); err != nil {
		resp.Fail(c, err, "验证失败")
		return
	}
	resp.OK(c, nil, "邮箱已验证")
}

// generateToken 生成随机 token
func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// now is used in templates
var now = time.Now
