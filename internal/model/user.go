package model

// UserView 用户序列化器
type UserView struct {
	ID        uint   `json:"id"`
	UserName  string `json:"user_name"`
	Nickname  string `json:"nickname"`
	Status    int    `json:"status"`
	Avatar    string `json:"avatar"`
	CreatedAt int64  `json:"created_at"`
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
}
