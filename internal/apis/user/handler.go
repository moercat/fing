package user

import (
	"fing/internal/model"
	"fing/internal/service/user"
	"fing/pkg/resp"
	"github.com/gin-gonic/gin"
)

// getProfile 获取当前用户资料
func (r *RouterUser) getProfile(c *gin.Context) {
	uid, _ := c.Get("user_id")
	view, err := user.GetProfile(uid.(uint))
	if err != nil {
		resp.Fail(c, err, "获取资料失败")
		return
	}
	resp.OK(c, view, "")
}

// updateProfile 修改资料
func (r *RouterUser) updateProfile(c *gin.Context) {
	var sv model.UpdateProfile
	if err := c.ShouldBind(&sv); err != nil {
		resp.Fail(c, err, "参数错误")
		return
	}
	uid, _ := c.Get("user_id")
	if err := user.UpdateProfile(uid.(uint), &sv); err != nil {
		resp.Fail(c, err, "更新失败")
		return
	}
	resp.OK(c, nil, "更新成功")
}

// changePassword 修改密码
func (r *RouterUser) changePassword(c *gin.Context) {
	var sv model.ChangePassword
	if err := c.ShouldBind(&sv); err != nil {
		resp.Fail(c, err, "参数错误")
		return
	}
	uid, _ := c.Get("user_id")
	if err := user.ChangePassword(uid.(uint), sv.OldPassword, sv.NewPassword); err != nil {
		resp.Fail(c, err, "修改失败")
		return
	}
	resp.OK(c, nil, "密码已更新")
}

// updateAvatar 修改头像
func (r *RouterUser) updateAvatar(c *gin.Context) {
	var sv model.UpdateAvatar
	if err := c.ShouldBind(&sv); err != nil {
		resp.Fail(c, err, "参数错误")
		return
	}
	uid, _ := c.Get("user_id")
	if err := user.UpdateAvatar(uid.(uint), sv.Avatar); err != nil {
		resp.Fail(c, err, "更新失败")
		return
	}
	resp.OK(c, nil, "头像已更新")
}

// listUsers 管理员列出用户
func (r *RouterUser) listUsers(c *gin.Context) {
	users, err := user.ListUsers()
	if err != nil {
		resp.Fail(c, err, "查询失败")
		return
	}
	resp.OK(c, users, "")
}

// updateUserRole 管理员修改用户角色
func (r *RouterUser) updateUserRole(c *gin.Context) {
	var sv model.UpdateRole
	if err := c.ShouldBind(&sv); err != nil {
		resp.Fail(c, err, "参数错误")
		return
	}
	id := c.Param("id")
	if err := user.UpdateRole(id, sv.Role); err != nil {
		resp.Fail(c, err, "更新失败")
		return
	}
	resp.OK(c, nil, "角色已更新")
}
