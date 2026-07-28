// Package password 实现密码重置与邮箱验证的业务逻辑。
package password

import (
	"crypto/rand"
	"encoding/hex"
	"fing/pkg/db"
	"fing/pkg/entity/usr"
	"fing/pkg/errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	tokenTTL = 30 * time.Minute
)

// IssueResetToken 签发密码重置 token，写入 user 表。
// 返回 token 与可能错误。如果邮箱不存在也返回 nil error，避免邮箱枚举攻击。
func IssueResetToken(email string) (string, error) {
	var u usr.UserInfo
	if err := db.Gain.Where("email = ?", email).First(&u).Error; err != nil {
		return "", nil
	}

	token := newToken()
	expires := time.Now().Add(tokenTTL).Unix()
	if err := db.Gain.Model(&u).Updates(map[string]interface{}{
		"reset_token":   token,
		"reset_expires": expires,
	}).Error; err != nil {
		return "", errors.Wrap(err, 500, "生成令牌失败")
	}
	return token, nil
}

// ResetWithToken 用 token 重置密码
func ResetWithToken(token, newPassword string) error {
	var u usr.UserInfo
	if err := db.Gain.Where("reset_token = ?", token).First(&u).Error; err != nil {
		return errors.New(404, "无效的重置令牌")
	}
	if time.Now().Unix() > u.ResetExpires {
		return errors.New(410, "重置令牌已过期")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return errors.Wrap(err, 500, "密码加密失败")
	}
	if err := db.Gain.Model(&u).Updates(map[string]interface{}{
		"password":      string(hash),
		"reset_token":   "",
		"reset_expires": 0,
	}).Error; err != nil {
		return errors.Wrap(err, 500, "更新失败")
	}
	return nil
}

// VerifyEmail 用 token 验证邮箱
func VerifyEmail(token string) error {
	// 简化实现：用 reset_token 字段同时复用做验证 token
	var u usr.UserInfo
	if err := db.Gain.Where("reset_token = ?", token).First(&u).Error; err != nil {
		return errors.New(404, "无效的验证令牌")
	}
	if err := db.Gain.Model(&u).Updates(map[string]interface{}{
		"email_verified": true,
		"reset_token":    "",
		"reset_expires":  0,
	}).Error; err != nil {
		return errors.Wrap(err, 500, "验证失败")
	}
	return nil
}

func newToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
