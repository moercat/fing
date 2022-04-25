package login

import (
	"fing/internal/model"
	"fing/internal/service/login"
	"fing/internal/tools"
	"fing/pkg/resp"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// register 用户注册接口
func (r *RouterLogin) register(c *gin.Context) {
	var sv model.Register
	if err := c.ShouldBind(&sv); err != nil {
		resp.Fail(c, err, "")
		return
	}

	err := login.Register(&sv)
	resp.CkFail(500, err)

	resp.OK(c, nil, "")
}

// login 用户登录接口
func (r *RouterLogin) login(c *gin.Context) {
	var sv model.Login
	if err := c.ShouldBind(&sv); err == nil {
		resp.Fail(c, err, "")
		return
	}

	usr, err := login.Login(&sv)
	resp.CkFail(500, err)

	s := sessions.Default(c)
	s.Clear()
	s.Set("user_id", usr.ID)
	err = s.Save()
	if err != nil {
		resp.Fail(c, err, "")
		return
	}

	resp.OK(c, usr, "")
}

// userInfo 用户详情
func (r *RouterLogin) userInfo(c *gin.Context) {

	res := tools.CurrentUser(c)

	resp.OK(c, res, "")
}

// logout 用户登出
func (r *RouterLogin) logout(c *gin.Context) {

	s := sessions.Default(c)
	s.Clear()
	err := s.Save()
	if err != nil {
		resp.Fail(c, err, "")
		return
	}

	resp.OK(c, nil, "登出成功")
}

// ping 状态检查页面
func (r *RouterLogin) ping(c *gin.Context) {

	resp.OK(c, "", "pong")
}
