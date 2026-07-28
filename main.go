package main

import (
	"fing/internal"
	"fing/log"
	"fing/pkg/cobra"
	"fing/pkg/config"
	"fing/pkg/db"
	"fing/pkg/entity/usr"
	"fing/pkg/graceful"
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

	// 5. 启动 HTTP 服务（优雅关闭）
	addr := ":" + config.Config.Port
	logger.Infof("fing listening on %s (mode=%s)", addr, config.Config.Mode)
	graceful.Run(graceful.Config{
		Addr:    addr,
		Handler: r,
		Timeout: 30 * time.Second,
		Cleanups: []func(){
			closeDB,
			closeRedis,
		},
	})
}

// closeDB 关闭 DB 连接
func closeDB() {
	sqlDB, err := db.Gain.DB()
	if err != nil {
		logger.Errorf("get sql.DB failed: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		logger.Errorf("close DB failed: %v", err)
	}
}

// closeRedis 关闭 Redis 连接
func closeRedis() {
	if err := db.RedisClient.Close(); err != nil {
		logger.Errorf("close Redis failed: %v", err)
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
func registerTasks() {
	// 每 30 分钟清理过期的密码重置 token
	cobra.Register("cleanup-expired-tokens", 30*time.Minute, func() {
		now := time.Now().Unix()
		// 只清空 reset_token 和 reset_expires，保留用户记录
		res := db.Gain.Model(&usr.UserInfo{}).
			Where("reset_token <> ? AND reset_expires > 0 AND reset_expires < ?", "", now).
			Updates(map[string]interface{}{
				"reset_token":   "",
				"reset_expires": 0,
			})
		if res.Error != nil {
			logger.Errorf("[task] cleanup expired tokens failed: %v", res.Error)
			return
		}
		logger.Infof("[task] cleaned %d expired reset tokens", res.RowsAffected)
	})
}