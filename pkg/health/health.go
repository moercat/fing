package health

import (
	"fing/pkg/db"
	"fing/pkg/resp"
	"github.com/gin-gonic/gin"
	"time"
)

// HealthCheck 健康检查接口
func HealthCheck(c *gin.Context) {
	// 检查数据库连接
	dbStartTime := time.Now()
	var count int64
	result := db.Gain.Model(&count).Limit(1).Count(&count)
	dbDuration := time.Since(dbStartTime)

	if result.Error != nil {
		resp.Fail(c, result.Error, "数据库连接失败", 503)
		return
	}

	// 检查Redis连接
	redisStartTime := time.Now()
	_, err := db.RedisClient.Ping().Result()
	redisDuration := time.Since(redisStartTime)

	if err != nil {
		resp.Fail(c, err, "Redis连接失败", 503)
		return
	}

	// 返回健康状态
	resp.OK(c, gin.H{
		"status": "healthy",
		"checks": gin.H{
			"database": gin.H{
				"status":   "up",
				"duration": dbDuration.Milliseconds(),
			},
			"redis": gin.H{
				"status":   "up",
				"duration": redisDuration.Milliseconds(),
			},
		},
		"timestamp": time.Now().Unix(),
	}, "服务健康")
}
