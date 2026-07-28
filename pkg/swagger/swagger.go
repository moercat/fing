// Package swagger 提供 Swagger UI 集成。
//
// 用法：
//   1. 在 API handler 上写 swag 注释（参考 examples/）
//   2. 安装 swag CLI: go install github.com/swaggo/swag/cmd/swag@latest
//   3. 生成文档: swag init -g main.go -o docs/
//   4. 访问 http://localhost:9765/swagger/index.html
//
// 注意：docs/docs.go 是 swag 生成的占位 stub，路径全空。
// 真实部署前必须先跑 swag init，否则 API 文档里看不到接口。
//
// handler 注释示例：
//
//	// @Summary  用户注册
//	// @Tags     auth
//	// @Accept   json
//	// @Produce  json
//	// @Param    req  body  model.Register  true  "注册参数"
//	// @Success  200  {object}  resp.Response
//	// @Router   /api/v1/register [post]
package swagger

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// docs 包由 swag init 自动生成
	_ "fing/docs"
)

// Register 挂载 /swagger/index.html 路由
func Register(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// @title           fing API
// @version         1.0
// @description     fing 后端脚手架 API 文档
// @host            localhost:9765
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
type swaggerInfo struct{}