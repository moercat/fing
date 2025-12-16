# fing

Simple architecture, clear directory structure, giving users an unusual feeling

简单的架构，清晰的目录结构，带给使用者非一般的感觉

![example](https://img.shields.io/badge/Go-1.16-blue)
![example](https://img.shields.io/badge/Gin%20-1.7.1-lightgrey)
![example](https://img.shields.io/badge/Redis-go--redis-red)
![example](https://img.shields.io/badge/Gorm%20-gorm.io-red)
![example](https://img.shields.io/badge/Xorm%20-xorm.io-red)
![example](https://img.shields.io/badge/Elastic-olivere-blue)
![example](https://img.shields.io/badge/License-MIT-green)

## 项目介绍

本项目是一个基于 Go 语言开发的 Web 应用程序，采用了一系列 Golang 中流行的组件，可以以本项目为基础快速搭建 Restful Web API。项目具备完整的用户认证系统，包括注册、登录、会话管理等功能。

## 功能特性

1. **用户系统**：
   - 用户注册（带密码确认验证，使用 bcrypt 加密）
   - 用户登录认证（基于 Session）
   - 用户信息查询
   - 用户登出功能

2. **技术集成**：
   - [Gin](https://github.com/gin-gonic/gin): 轻量级高性能 Web 框架
   - [GORM](https://gorm.io/index.html) 和 [XORM](https://xorm.io/index.html): ORM 工具，支持 MySQL
   - [Go-Redis](https://github.com/go-redis/redis): Redis 客户端
   - [Elasticsearch](https://github.com/olivere/elastic): 搜索引擎客户端
   - [Gin-Sessions](https://github.com/gin-contrib/sessions): 会话管理
   - [Configor](https://github.com/jinzhu/configor): 支持多种格式的配置管理
   - [Cors](https://github.com/gin-contrib/cors): 跨域中间件
   - [Gomail](https://github.com/go-gomail/gomail): 邮件发送功能

3. **架构特色**：
   - 清晰的分层架构（API 层、Service 层、Model 层、Entity 层）
   - 完善的中间件体系（CORS、Session、Auth、Recover、Logger）
   - 统一的响应格式封装
   - 模块化的代码组织结构
   - 标准化的错误处理机制
   - 结构化日志记录

## 目录结构

```
fing/
├── config*.yaml          # 配置文件
├── go.mod/go.sum        # Go 依赖管理
├── main.go              # 主程序入口
├── Dockerfile           # Docker 配置文件
├── docker-compose.yml   # Docker Compose 配置文件
├── Makefile             # 构建和部署脚本
├── APIDOC.md            # API 文档
├── .env.example         # 环境变量示例
├── internal/            # 内部业务代码
│   ├── router.go        # 路由初始化
│   ├── apis/            # API 接口实现
│   │   └── login/       # 登录相关接口
│   ├── model/           # 数据模型（序列化器）
│   ├── service/         # 业务逻辑层
│   └── tools/           # 工具函数
├── pkg/                 # 公共包
│   ├── cobra/           # 命令行工具（定时任务）
│   ├── config/          # 配置管理
│   ├── db/              # 数据库连接管理（MySQL、Redis、ES）
│   ├── elastic/         # Elasticsearch 操作
│   ├── email/           # 邮件功能
│   ├── entity/          # 数据库实体模型
│   ├── errors/          # 标准错误定义
│   ├── health/          # 健康检查
│   ├── middleware/      # 中间件（认证、CORS、Session等）
│   ├── resp/            # 统一响应格式封装
│   └── tools/           # 工具函数
└── log/                 # 日志相关
```

## API 接口

### 公共接口

- `GET /health` - 服务健康检查
- `GET /api/v1/ping` - 服务状态检查

### 用户接口 (V1)

- `POST /api/v1/register` - 用户注册
- `POST /api/v1/login` - 用户登录

### 用户接口 (V2) - 需要认证

- `GET /api/v2/user_info` - 用户信息
- `DELETE /api/v2/logout` - 用户登出

## 环境配置

项目依赖于 `configor`，根据环境变量加载配置文件。使用 `CONFIGOR_ENV` 设置环境变量，如果没有设置 `CONFIGOR_ENV`，框架将使用 `development` 作为默认环境，即读取 `config.development.yaml`。

区分开发环境与生产环境的命令：
```
CONFIGOR_ENV=production go run main.go
```

### 配置文件说明

配置文件示例 (`config.yaml`)：

```yaml
mode: dev          # 当前所在环境
secret: your_secret_key  # session 加密密钥，不要使用默认值!
level: 4          # 日志等级
port: 9765        # 服务端口

dataSource:
  main: you_name:you_password@tcp(you_ip:you_port)/db_name?charset=utf8mb4  # 数据库连接

redis:
  addr: 127.0.0.1:6379  # Redis 连接地址

es:  # Elasticsearch 连接配置
  esUrl: you_es
  esUsername: ""
  esPassword: ""

email:  # 邮件发送配置
  host: smtp.exmail.qq.com
  name: my_name
  email: my_email
  password: ""
```

## 快速开始

### 开发环境准备

```bash
# 1. 下载并安装 Go 1.16+
# 2. 克隆项目
git clone https://github.com/moercat/fing.git
cd fing

# 3. 安装依赖
go mod tidy

# 4. 复制环境变量文件
cp .env.example .env

# 5. 修改配置文件（config.yaml）连接到你的数据库等服务
# 6. 启动服务
go run main.go
```

### 使用 Docker 运行

```bash
# 构建并启动服务
docker-compose up -d
```

### 使用 Makefile

```bash
# 构建项目
make build

# 运行项目
make run

# 运行测试
make test

# 构建生产环境二进制文件
make build-prod

# Docker 构建
make docker-build

# 使用 Docker Compose 启动
make docker-up
```

## 构建部署

### 本地构建

Linux 环境下构建：
```bash
GOOS=linux GOARCH=amd64 go build -o fing
```

### Docker 部署

```bash
# 构建 Docker 镜像
docker build -t fing-app .

# 运行容器
docker run -d -p 9765:9765 fing-app
```

## 项目特性

1. **安全性**：
   - 密码使用 bcrypt 加密存储
   - 配置验证，防止使用默认安全密钥
   - 完善的输入验证
   - 会话安全设置

2. **可维护性**：
   - 清晰的分层架构
   - 标准化的错误处理
   - 统一的响应格式
   - 详细的 API 文档

3. **生产就绪**：
   - 结构化日志记录
   - 健康检查接口
   - 容器化部署支持
   - 性能监控

4. **开发体验**：
   - 基于 Makefile 的一键操作
   - Docker 和 Docker Compose 支持
   - 详细的文档和示例

## 许可证

MIT License
