package usr

import "gorm.io/gorm"

// UserInfo 用户模型
type UserInfo struct {
	gorm.Model
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Status   int    `json:"status"`
	Avatar   string `json:"avatar" gorm:"size:1000"`
}

func (i UserInfo) TableName() string {
	return "user_info"
}
