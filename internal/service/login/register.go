package login

import (
	"errors"
	"fing/internal/model"
	"fing/pkg/db"
	"fing/pkg/elastic"
	"fing/pkg/entity/usr"
	"golang.org/x/crypto/bcrypt"
)

const (
	PassWordCost = 12 // 加密等级
)

// Register 用户注册
func Register(ser *model.Register) error {

	// 表单验证
	if ser.RePassword != ser.Password {
		return errors.New("两次输入的密码不相同")
	}

	exist, _ := db.Main.Where("nickname = ?", ser.Nickname).Exist(&usr.UserInfo{})
	if exist {
		return errors.New("昵称被占用")
	}

	exist, _ = db.Main.Where("user_name = ?", ser.UserName).Exist(&usr.UserInfo{})
	if exist {
		return errors.New("用户名已经注册")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(ser.Password), PassWordCost)
	if err != nil {
		return err
	}

	// 创建用户
	if err := db.Gain.Create(&usr.UserInfo{
		Nickname: ser.Nickname,
		UserName: ser.UserName,
		Password: string(bytes),
		Status:   0,
	}).Error; err != nil {
		return errors.New("注册失败")
	}

	var u usr.UserInfo
	db.Gain.Model(usr.UserInfo{
		Nickname: ser.Nickname,
		UserName: ser.UserName,
		Password: string(bytes),
		Status:   0,
	}).First(&u)

	elastic.Index(elastic.UserInfoIndex).AddUser(&u)

	return nil
}
