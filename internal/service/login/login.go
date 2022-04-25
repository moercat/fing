package login

import (
	"errors"
	"fing/internal/model"
	"fing/pkg/db"
	"fing/pkg/entity/usr"
	"golang.org/x/crypto/bcrypt"
)

// GetUser 用ID获取用户
func GetUser(ID interface{}) (usr.UserInfo, error) {
	var userInfo usr.UserInfo

	result := db.Gain.First(&userInfo, ID)

	return userInfo, result.Error
}

// Login 用户登录函数
func Login(ser *model.Login) (model.UserView, error) {
	var (
		user  usr.UserInfo
		uView model.UserView
	)

	if err := db.Gain.Where("user_name = ?", ser.UserName).First(&user).Error; err != nil {
		return uView, errors.New("账号或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ser.Password)); err != nil {
		return uView, errors.New("账号或密码错误")
	}

	uView = model.UserView{
		ID:        user.ID,
		UserName:  user.UserName,
		Nickname:  user.Nickname,
		Status:    user.Status,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt.Unix(),
	}

	return uView, nil
}
