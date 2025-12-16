package login

import (
	"fing/internal/model"
	"fing/pkg/db"
	"fing/pkg/elastic"
	"fing/pkg/entity/usr"
	"fing/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	PassWordCost = 12 // 加密等级
)

// Register 用户注册
func Register(ser *model.Register) error {

	// 表单验证
	if ser.RePassword != ser.Password {
		return errors.New(422, "两次输入的密码不相同")
	}

	exist, _ := db.Main.Where("nickname = ?", ser.Nickname).Exist(&usr.UserInfo{})
	if exist {
		return errors.New(422, "昵称被占用")
	}

	exist, _ = db.Main.Where("user_name = ?", ser.UserName).Exist(&usr.UserInfo{})
	if exist {
		return errors.New(422, "用户名已经注册")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(ser.Password), PassWordCost)
	if err != nil {
		return errors.Wrap(err, 500, "密码加密失败")
	}

	// 创建用户
	user := usr.UserInfo{
		Nickname: ser.Nickname,
		UserName: ser.UserName,
		Password: string(bytes),
		Status:   0,
	}

	if err := db.Gain.Create(&user).Error; err != nil {
		return errors.Wrap(err, 500, "注册失败")
	}

	elastic.Index(elastic.UserInfoIndex).AddUser(&user)

	return nil
}
