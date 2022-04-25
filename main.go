package main

import (
	"fing/internal"
	"fing/pkg/cobra"
	"fing/pkg/config"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 装载路由
	r = internal.InitRouter(r)

	// 定时任务
	cobra.Cobra()

	if err := r.Run(":" + config.Config.Port); err != nil {
		panic(err)
	}
}
