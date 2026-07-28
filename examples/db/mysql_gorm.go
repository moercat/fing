//go:build ignore
// +build ignore

// Package main 演示 fing 的 MySQL（GORM）使用方式。
//
// 业务代码统一从 fing/pkg/db 拿 *gorm.DB 实例（变量名 db.Gain），
// 不要自己在业务代码里 NewConnection。
package main

import (
	"errors"
	"fmt"
	"time"

	"fing/internal/model"
	"fing/pkg/db"
	"fing/pkg/entity/usr"
	"gorm.io/gorm"
)

// ============================================================
// 基础 CRUD
// ============================================================

// exampleCreate 插入一条用户
func exampleCreate() {
	user := usr.UserInfo{
		UserName: "alice",
		Nickname: "Alice",
		Email:    "alice@example.com",
		Password: "hashed-password", // 实际用 bcrypt.GenerateFromPassword
		Role:     "user",
	}
	// 自动建表需要 AutoMigrate 调过（main 启动时）
	if err := db.Gain.Create(&user).Error; err != nil {
		fmt.Println("create error:", err)
		return
	}
	fmt.Println("created user id:", user.ID)
}

// exampleRead 按 ID 查一条
func exampleRead() {
	var u usr.UserInfo
	if err := db.Gain.First(&u, 1).Error; err != nil {
		fmt.Println("not found:", err)
		return
	}
	fmt.Printf("user: %+v\n", u)
}

// exampleUpdate 更新字段
func exampleUpdate() {
	// 1. 更新单个字段
	db.Gain.Model(&usr.UserInfo{}).
		Where("id = ?", 1).
		Update("nickname", "Alice2")

	// 2. 更新多个字段（用 struct，零值不更新）
	db.Gain.Model(&usr.UserInfo{}).
		Where("id = ?", 1).
		Updates(usr.UserInfo{Nickname: "Alice3", Avatar: "https://..."})

	// 3. 更新多个字段（用 map，可以更新零值）
	db.Gain.Model(&usr.UserInfo{}).
		Where("id = ?", 1).
		Updates(map[string]interface{}{"status": 0})
}

// exampleDelete 软删除（GORM 默认带 DeletedAt 字段就软删）
func exampleDelete() {
	if err := db.Gain.Delete(&usr.UserInfo{}, 1).Error; err != nil {
		fmt.Println("delete error:", err)
	}
}

// ============================================================
// 查询模式
// ============================================================

// exampleWhere 各种查询方式
func exampleWhere() {
	var users []usr.UserInfo

	// 1. 条件查询
	db.Gain.Where("user_name = ?", "alice").Find(&users)

	// 2. IN 查询
	db.Gain.Where("id IN ?", []uint{1, 2, 3}).Find(&users)

	// 3. LIKE 查询
	db.Gain.Where("nickname LIKE ?", "%ali%").Find(&users)

	// 4. 组合条件
	db.Gain.Where("status = ? AND role = ?", 0, "user").Find(&users)

	// 5. 排序 + 限制
	db.Gain.Order("created_at DESC").Limit(10).Find(&users)
}

// exampleCount 统计
func exampleCount() {
	var count int64
	db.Gain.Model(&usr.UserInfo{}).Where("role = ?", "user").Count(&count)
	fmt.Println("user count:", count)
}

// examplePagination 分页
func examplePagination() {
	page := 1
	pageSize := 20
	var users []model.UserView

	db.Gain.Model(&usr.UserInfo{}).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users)
}

// ============================================================
// 事务（重要！）
// ============================================================

// exampleTransaction 事务示例：转账（A 给 B 转 100）
func exampleTransaction() error {
	// 方式 1：手动事务
	tx := db.Gain.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	// ... 操作 1
	if err := tx.Exec("UPDATE user_info SET balance = balance - 100 WHERE id = ?", 1).Error; err != nil {
		tx.Rollback()
		return err
	}
	// ... 操作 2
	if err := tx.Exec("UPDATE user_info SET balance = balance + 100 WHERE id = ?", 2).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// exampleTransactionAuto 自动事务：任何一个 step 失败就回滚
func exampleTransactionAuto() (err error) {
	return db.Gain.Transaction(func(tx *gorm.DB) error {
		// 在事务里做多个操作
		if err := tx.Create(&usr.UserInfo{UserName: "u1"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&usr.UserInfo{UserName: "u2"}).Error; err != nil {
			return err
		}
		return nil
	})
}

// exampleErrorHandling 常见错误处理
func exampleErrorHandling() {
	var u usr.UserInfo
	err := db.Gain.First(&u, 99999).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("user not found")
	}
}

// ============================================================
// 性能：批量插入 / 索引
// ============================================================

// exampleBatchInsert 批量插入 1000 条
func exampleBatchInsert() {
	users := make([]usr.UserInfo, 1000)
	for i := range users {
		users[i] = usr.UserInfo{
			UserName: fmt.Sprintf("user_%d", i),
			Email:    fmt.Sprintf("u%d@example.com", i),
			Password: "x",
		}
	}
	// 1000 条/批
	db.Gain.CreateInBatches(users, 1000)
}

// exampleIndex 索引（迁移时建）
func exampleIndex() {
	// 在 main 启动时调一次
	db.Gain.AutoMigrate(&usr.UserInfo{})

	// 已有的索引定义（看 entity/usr/user.go 的 struct tag）
	// UserName: gorm:"size:30;uniqueIndex"
	// Email:    gorm:"size:100;uniqueIndex"
}

// ============================================================
// 时间处理
// ============================================================

func exampleTimeQuery() {
	// 查最近 7 天注册的用户
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	var users []usr.UserInfo
	db.Gain.Where("created_at > ?", sevenDaysAgo).Find(&users)
}

// 引入 gorm 错误包
var _ = gorm.ErrRecordNotFound
