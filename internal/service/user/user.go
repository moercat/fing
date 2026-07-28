// Package user 实现用户资料的增删改查逻辑。
package user

import (
	"fing/internal/model"
	"fing/pkg/db"
	"fing/pkg/entity/usr"
	"fing/pkg/errors"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// GetProfile 获取用户资料
func GetProfile(uid uint) (model.UserView, error) {
	var u usr.UserInfo
	if err := db.Gain.First(&u, uid).Error; err != nil {
		return model.UserView{}, errors.Wrap(err, 500, "查询失败")
	}
	return toView(u), nil
}

// UpdateProfile 更新昵称/邮箱
func UpdateProfile(uid uint, p *model.UpdateProfile) error {
	updates := map[string]interface{}{}
	if p.Nickname != "" {
		updates["nickname"] = p.Nickname
	}
	if p.Email != "" {
		updates["email"] = p.Email
		updates["email_verified"] = false // 邮箱变更需重新验证
	}
	if len(updates) == 0 {
		return errors.New(422, "无更新内容")
	}
	if err := db.Gain.Model(&usr.UserInfo{}).Where("id = ?", uid).Updates(updates).Error; err != nil {
		return errors.Wrap(err, 500, "更新失败")
	}
	return nil
}

// ChangePassword 修改密码
func ChangePassword(uid uint, oldPwd, newPwd string) error {
	var u usr.UserInfo
	if err := db.Gain.First(&u, uid).Error; err != nil {
		return errors.Wrap(err, 500, "用户不存在")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldPwd)); err != nil {
		return errors.New(401, "原密码错误")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPwd), 12)
	if err != nil {
		return errors.Wrap(err, 500, "密码加密失败")
	}
	if err := db.Gain.Model(&u).Update("password", string(hash)).Error; err != nil {
		return errors.Wrap(err, 500, "更新失败")
	}
	return nil
}

// UpdateAvatar 更新头像
func UpdateAvatar(uid uint, avatar string) error {
	if err := db.Gain.Model(&usr.UserInfo{}).Where("id = ?", uid).Update("avatar", avatar).Error; err != nil {
		return errors.Wrap(err, 500, "更新失败")
	}
	return nil
}

// ListUsers 列出所有用户
func ListUsers() ([]model.UserView, error) {
	var users []usr.UserInfo
	if err := db.Gain.Find(&users).Error; err != nil {
		return nil, errors.Wrap(err, 500, "查询失败")
	}
	views := make([]model.UserView, 0, len(users))
	for _, u := range users {
		views = append(views, toView(u))
	}
	return views, nil
}

// UpdateRole 管理员更新用户角色
func UpdateRole(id, role string) error {
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errors.New(422, "无效的用户ID")
	}
	if err := db.Gain.Model(&usr.UserInfo{}).Where("id = ?", uid).Update("role", role).Error; err != nil {
		return errors.Wrap(err, 500, "更新失败")
	}
	return nil
}

func toView(u usr.UserInfo) model.UserView {
	return model.UserView{
		ID:            u.ID,
		UserName:      u.UserName,
		Nickname:      u.Nickname,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		Role:          u.Role,
		Status:        u.Status,
		Avatar:        u.Avatar,
		CreatedAt:     u.CreatedAt.Unix(),
	}
}
