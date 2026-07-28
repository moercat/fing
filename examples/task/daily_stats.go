//go:build ignore
// +build ignore

// Package main 演示 fing 的定时任务用法。
//
// 定时任务要做什么？这里列了 3 个常见场景：
//  1. 清理过期数据（每天凌晨）
//  2. 同步数据（每小时）
//  3. 发送日报（每天 9 点）
//
// 所有任务在 main.go 启动时注册，cobra.Cobra() 启动后会自动按周期执行。
//
// 把这个文件 import 到 main.go 即可生效：
//
//	import _ "fing/examples/task"
package main

import (
	"fmt"
	"time"

	"fing/pkg/cobra"
)

func init() {
	// 任务 1：每 5 分钟清理一次过期 token（演示用，间隔很短）
	cobra.Register("cleanup-expired-tokens", 5*time.Minute, cleanupExpiredTokens)

	// 任务 2：每小时同步一次数据
	cobra.Register("sync-data", time.Hour, syncData)

	// 任务 3：每天 9 点发日报
	cobra.Register("daily-report", 24*time.Hour, dailyReport)
}

// cleanupExpiredTokens 清理过期的密码重置 token
func cleanupExpiredTokens() {
	// 业务代码写这里
	// 例如：db.Gain.Where("reset_expires < ?", time.Now().Unix()).
	//     Delete(&usr.UserInfo{}, "reset_token != ''")
	fmt.Println("[task] cleanup expired tokens")
}

// syncData 同步外部数据
func syncData() {
	// 业务代码写这里
	fmt.Println("[task] sync external data")
}

// dailyReport 发送日报
func dailyReport() {
	hour := time.Now().Hour()
	if hour != 9 {
		// 不是 9 点就跳过（实现每天 9 点运行）
		return
	}
	// 业务代码写这里
	fmt.Println("[task] send daily report")
}
