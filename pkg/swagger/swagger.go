// Package swagger 提供 Swagger UI 集成。
//
// 用法：
//   1. 在 API handler 上写 swag 注释
//   2. swag init 生成 docs/
//   3. main.go 加 swagger.Register(r)
//   4. 访问 http://localhost:9765/swagger/index.html
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

	// 触发 swag init 生成的 docs 包（如果没有会编译报错，先放占位）
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