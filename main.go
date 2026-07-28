package main

import (
	"fing/internal"
	"fing/log"
	"fing/pkg/cobra"
	"fing/pkg/config"
	"fing/pkg/db"
	"fing/pkg/entity/usr"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化路由
	r := gin.Default()
	r = internal.InitRouter(r)

	// 2. 启动时自动建表（生产建议改用 migration 工具）
	autoMigrate()

	// 3. 注册定时任务
	registerTasks()

	// 4. 启动定时任务调度器
	cobra.Cobra()

	// 5. 启动 HTTP 服务
	addr := ":" + config.Config.Port
	logger.Infof("fing listening on %s (mode=%s)", addr, config.Config.Mode)
	if err := r.Run(addr); err != nil {
		panic(err)
	}
}

// autoMigrate 启动时自动迁移数据库表结构。
// 生产环境建议改用 golang-migrate 或 goose 等迁移工具。
func autoMigrate() {
	if err := db.Gain.AutoMigrate(&usr.UserInfo{}); err != nil {
		panic(fmt.Sprintf("auto migrate failed: %v", err))
	}
	logger.Info("auto migrate: user_info table ready")
}

// registerTasks 注册业务定时任务。
// 在这里加你自己的定时任务，周期和函数都自己定。
func registerTasks() {
	// 示例：每 30 分钟清理一次过期的密码重置 token
	// 你可以删掉这段，自己加任务
	cobra.Register("cleanup-expired-tokens", 30*time.Minute, func() {
		// 实现逻辑见 TUTORIAL.md 第 5 步
		logger.Info("[task] cleanup expired tokens (skeleton, implement me)")
	})
}