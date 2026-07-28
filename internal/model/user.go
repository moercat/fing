package model

// UserView 用户序列化器
type UserView struct {
	ID            uint   `json:"id"`
	UserName      string `json:"user_name"`
	Nickname      string `json:"nickname"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Role          string `json:"role"`
	Status        int    `json:"status"`
	Avatar        string `json:"avatar"`
	CreatedAt     int64  `json:"created_at"`
}

// Login 管理用户登录的服务
type Login struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=5,max=30"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=40"`
}

// Register 管理用户注册服务
type Register struct {
	Nickname   string `form:"nickname" json:"nickname" binding:"required,min=2,max=30"`
	UserName   string `form:"user_name" json:"user_name" binding:"required,min=5,max=30"`
	Password   string `form:"password" json:"password" binding:"required,min=8,max=40"`
	RePassword string `form:"re_password" json:"re_password" binding:"required,min=8,max=40"`
	Email      string `form:"email" json:"email" binding:"required,email"`
}

// UpdateProfile 修改资料请求
type UpdateProfile struct {
	Nickname string `form:"nickname" json:"nickname" binding:"omitempty,min=2,max=30"`
	Email    string `form:"email" json:"email" binding:"omitempty,email"`
}

// ChangePassword 修改密码请求
type ChangePassword struct {
	OldPassword string `form:"old_password" json:"old_password" binding:"required,min=8,max=40"`
	NewPassword string `form:"new_password" json:"new_password" binding:"required,min=8,max=40"`
}

// UpdateAvatar 修改头像请求
type UpdateAvatar struct {
	Avatar string `form:"avatar" json:"avatar" binding:"required,url"`
}

// UpdateRole 修改用户角色请求
type UpdateRole struct {
	Role string `form:"role" json:"role" binding:"required,oneof=admin user"`
}

// ForgotPassword 申请重置密码请求
type ForgotPassword struct {
	Email string `form:"email" json:"email" binding:"required,email"`
}

// ResetPassword 用 token 重置密码请求
type ResetPassword struct {
	Token       string `form:"token" json:"token" binding:"required"`
	NewPassword string `form:"new_password" json:"new_password" binding:"required,min=8,max=40"`
}

// VerifyEmail 邮箱验证请求
type VerifyEmail struct {
	Token string `form:"token" json:"token" binding:"required"`
}
