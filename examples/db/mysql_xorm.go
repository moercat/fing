//go:build ignore
// +build ignore

// Package main 演示 fing 的 MySQL（XORM）使用方式。
//
// 注意：fing 同时引入了 GORM 和 XORM 两个 ORM。
//   - 业务代码建议用 GORM（更现代，文档多）
//   - XORM 用于已有遗留代码 / 特殊场景
package main

import (
	"fmt"

	"fing/pkg/db"
	"fing/pkg/entity/usr"
)

// ============================================================
// 基础 CRUD
// ============================================================

func exampleXormCreate() {
	user := new(usr.UserInfo)
	user.UserName = "alice"
	user.Nickname = "Alice"
	user.Email = "alice@example.com"

	affected, err := db.Main.Insert(user)
	if err != nil {
		fmt.Println("insert error:", err)
		return
	}
	fmt.Println("inserted:", affected)
}

func exampleXormRead() {
	user := new(usr.UserInfo)
	// 按主键查
	has, err := db.Main.ID(1).Get(user)
	if err != nil {
		fmt.Println("query error:", err)
		return
	}
	if !has {
		fmt.Println("not found")
		return
	}
	fmt.Printf("user: %+v\n", user)
}

func exampleXormUpdate() {
	// 按条件更新
	user := new(usr.UserInfo)
	user.Nickname = "Alice-New"
	affected, err := db.Main.ID(1).Cols("nickname").Update(user)
	if err != nil {
		fmt.Println("update error:", err)
		return
	}
	fmt.Println("updated rows:", affected)
}

func exampleXormDelete() {
	// 软删（如果有 DeletedAt 字段）
	affected, err := db.Main.ID(1).Delete(&usr.UserInfo{})
	if err != nil {
		fmt.Println("delete error:", err)
		return
	}
	fmt.Println("deleted:", affected)
}

// ============================================================
// 查询
// ============================================================

func exampleXormQuery() {
	var users []usr.UserInfo

	// Where 条件
	err := db.Main.Where("user_name = ?", "alice").Find(&users)
	if err != nil {
		fmt.Println("find error:", err)
		return
	}
	fmt.Printf("found %d users\n", len(users))

	// 多条件
	err = db.Main.Where("status = ? AND role = ?", 0, "user").Find(&users)
	if err != nil {
		fmt.Println("find error:", err)
		return
	}

	// 排序 + 限制
	err = db.Main.OrderBy("created_at DESC").Limit(10, 0).Find(&users)
	if err != nil {
		fmt.Println("find error:", err)
		return
	}

	// 统计
	count, err := db.Main.Where("role = ?", "user").Count(&usr.UserInfo{})
	if err != nil {
		fmt.Println("count error:", err)
		return
	}
	fmt.Println("user count:", count)
}

// ============================================================
// 何时用 XORM
// ============================================================
//
// 推荐：业务代码统一用 GORM（fing 内置示例都用 GORM）
// XORM 适合：
//   - 你接手的老项目原本就用 XORM
//   - 需要 XORM 特有功能（如 SQL builder 链式）
//   - 性能关键场景（XORM 略快但不明显）
//
// 混用注意：不要在同一事务里同时用 GORM 和 XORM 操作同一行
