// Package docs 由 swag init 生成的占位 stub。
//
// 真实部署步骤：
//   go install github.com/swaggo/swag/cmd/swag@latest
//   swag init -g main.go -o docs/
//
// 上面的命令会覆盖本文件，生成真正的 OpenAPI spec。
package docs

import "github.com/swaggo/swag"

// SwaggerInfo 由 swag init 注入
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:9765",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "fing API",
	Description:      "fing 后端脚手架 API 文档",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

const docTemplate = `{
  "swagger": "2.0",
  "info": {
    "title": "fing API",
    "description": "fing 后端脚手架 API 文档",
    "version": "1.0"
  },
  "host": "localhost:9765",
  "basePath": "/",
  "paths": {}
}`