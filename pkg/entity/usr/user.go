package usr

import "gorm.io/gorm"

// UserInfo 用户模型
type UserInfo struct {
	gorm.Model
	UserName      string `json:"user_name" gorm:"size:30;uniqueIndex"`
	Password      string `json:"-" gorm:"size:100"`
	Nickname      string `json:"nickname" gorm:"size:30"`
	Email         string `json:"email" gorm:"size:100;uniqueIndex"`
	EmailVerified bool   `json:"email_verified" gorm:"default:false"`
	Role          string `json:"role" gorm:"size:20;default:user"`
	Status        int    `json:"status" gorm:"default:0"`
	Avatar        string `json:"avatar" gorm:"size:1000"`
	ResetToken    string `json:"-" gorm:"size:100;index"`
	ResetExpires  int64  `json:"-"`
}

func (i UserInfo) TableName() string {
	return "user_info"
}

// IsAdmin 是否为管理员
func (i UserInfo) IsAdmin() bool {
	return i.Role == "admin"
}
