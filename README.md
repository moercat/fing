# fing

> Gin + GORM + Redis + Sessions/JWT 的用户中心后端模板

![Go Version](https://img.shields.io/badge/Go-1.21-blue)
![Gin](https://img.shields.io/badge/Gin-1.10-lightgrey)
![Gorm](https://img.shields.io/badge/Gorm-gorm.io-red)
![License](https://img.shields.io/badge/License-MIT-green)
[![CI](https://github.com/moercat/fing/actions/workflows/go.yml/badge.svg)](https://github.com/moercat/fing/actions/workflows/go.yml)
![Stars](https://img.shields.io/badge/Stars-46-brightgreen)

## 这是什么

一个开箱即用的**用户中心**后端模板。覆盖注册、登录（Session + JWT 双模式）、密码重置、邮箱验证、RBAC 角色权限、资料修改、密码修改、头像修改。**复制即用**做你自己的业务后端。

适合做：
- 个人/中小团队的 Web 后端基座
- 学习 Gin + GORM 的工程模板
- 面试展示项目（覆盖常见后端能力）

## 功能特性

| 模块 | 能力 | 状态 |
| --- | --- | --- |
| 注册 | 邮箱+密码+昵称，bcrypt 加密 | ✅ |
| 登录 | Session / JWT 双模式 | ✅ |
| 登出 | 清空 Session | ✅ |
| 资料 | 查看/修改昵称、邮箱、头像 | ✅ |
| 密码 | 登录态修改 + 邮箱重置（30min token） | ✅ |
| 邮箱验证 | 注册后验证邮箱 | ✅ |
| 角色权限 | `user` / `admin` RBAC 中间件 | ✅ |
| 健康检查 | `/health` | ✅ |
| Redis 缓存 | session 存储 | ✅ |
| MySQL 持久化 | GORM 1.25 | ✅ |
| 邮件发送 | SMTP 客户端 | ✅ |
| 结构化日志 | 自定义 logger + 文件输出 | ✅ |
| 错误处理 | 统一错误码 + 响应包装 | ✅ |
| CORS | 中间件 | ✅ |
| API 文档 | 内联代码注释 | ✅ |
| Swagger 自动生成 | — | ❌ 计划中 |
| Elasticsearch 集成 | 高级搜索 | ❌ 未启用 |
| XORM 接入 | — | ❌ 未启用（与 GORM 重复） |
| 限流 | — | ❌ 计划中 |

## 快速开始

### Docker 一键运行

```bash
# 1. 启动 MySQL + Redis
docker compose up -d

# 2. 启动服务
go run main.go

# 3. 注册用户
curl -X POST http://localhost:9765/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "user_name": "testuser",
    "password": "testpass123",
    "re_password": "testpass123",
    "nickname": "测试用户",
    "email": "test@example.com"
  }'

# 4. 登录获取 JWT
curl -X POST http://localhost:9765/api/v1/login/jwt \
  -H "Content-Type: application/json" \
  -d '{"user_name":"testuser","password":"testpass123"}'
```

### 本地开发

```bash
# 1. 准备 MySQL + Redis
mysql -u root -p -e "CREATE DATABASE fing DEFAULT CHARSET utf8mb4"

# 2. 配置
cp .env.example .env
vim .env   # 填入你的 DB/Redis/SMTP 配置

# 3. 启动
go mod download
go run main.go
```

服务默认监听 `:9765`，访问 `http://localhost:9765/health` 验证。

## API 接口

### 公开接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/health` | 健康检查 |
| GET | `/api/v1/ping` | 服务存活 |
| POST | `/api/v1/register` | 用户注册 |
| POST | `/api/v1/login` | 登录（Session） |
| POST | `/api/v1/login/jwt` | 登录（JWT） |
| POST | `/api/v1/password/forgot` | 申请密码重置邮件 |
| POST | `/api/v1/password/reset` | 用 token 重置密码 |
| POST | `/api/v1/email/verify` | 验证邮箱 |

### 需要登录（`/api/v2/*`）

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/v2/user_info` | 当前用户信息（Session） |
| GET | `/api/v2/profile` | 获取资料 |
| PUT | `/api/v2/profile` | 修改昵称/邮箱 |
| PUT | `/api/v2/profile/password` | 修改密码 |
| PUT | `/api/v2/profile/avatar` | 修改头像 |
| DELETE | `/api/v2/logout` | 登出（Session） |

### 需要 admin 角色

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/v2/users` | 用户列表 |
| PUT | `/api/v2/users/:id/role` | 修改用户角色 |

## 技术栈

| 组件 | 选型 | 用途 |
| --- | --- | --- |
| Web 框架 | [Gin](https://github.com/gin-gonic/gin) 1.10 | 路由 / 中间件 |
| ORM | [GORM](https://gorm.io) 1.25 | MySQL 持久化 |
| 数据库 | MySQL 5.7+ | 主存储 |
| 缓存 | Redis 6+ | Session 存储 |
| 会话 | gin-contrib/sessions | Web 端 |
| 鉴权 | golang-jwt/jwt v5 | API 端 |
| 邮件 | go-gomail/gomail | 通知 |
| 配置 | jinzhu/configor | YAML + 环境变量 |
| 日志 | 自定义 logger | 结构化日志 |
| 密码 | golang.org/x/crypto/bcrypt | 加密 |

## 目录结构

```
fing/
├── main.go                 ← 入口
├── config.{dev,prod}.yaml  ← 配置文件
├── .env.example            ← 环境变量模板
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── APIDOC.md               ← 接口详情
├── examples/               ← 客户端示例
│   ├── curl.sh             ← cURL 调用示例
│   └── client.go           ← Go client 示例
├── log/                    ← 自定义 logger
└── internal/
    ├── router.go           ← 路由注册
    ├── apis/               ← HTTP handlers
    │   ├── login/          ← 注册/登录/JWT
    │   ├── user/           ← 资料/密码/角色
    │   └── password/       ← 密码重置/邮箱验证
    ├── service/            ← 业务逻辑
    │   ├── login/
    │   ├── user/
    │   └── password/
    ├── model/              ← 请求/响应序列化
    └── tools/              ← 工具函数
└── pkg/                    ← 可复用包
    ├── jwt/                ← JWT 签发/解析
    ├── middleware/         ← CORS/Auth/Logger/Recover
    ├── resp/               ← 统一响应
    ├── errors/             ← 统一错误
    ├── db/                 ← MySQL + Redis 客户端
    ├── email/              ← 邮件发送
    ├── config/             ← 配置加载
    ├── entity/usr/         ← 数据模型
    ├── health/             ← 健康检查
    └── cobra/              ← 定时任务
```

## 配置

通过 `configor` 加载：

```yaml
mode: dev                 # dev | prod
secret: your-secret-key   # JWT/session 加密密钥（生产必须改）
level: 4                  # 日志等级
port: 9765

dataSource:
  main: user:pass@tcp(localhost:3306)/fing?charset=utf8mb4&parseTime=True&loc=Local

redis:
  addr: localhost:6379

email:
  host: smtp.gmail.com
  name: "Your Name"
  email: your@email.com
  password: your-app-password
```

生产环境务必：
- 改 `secret` 为强随机字符串
- 用 `CONFIGOR_ENV=production` 加载 `config.production.yaml`
- 邮件密码用应用专用密码（Gmail/QQ 都需）

## 路线图

- [ ] Swagger 文档自动生成
- [ ] 请求限流（按 IP / 按用户）
- [ ] 单元测试覆盖（service 层）
- [ ] Docker 镜像推送到 GHCR
- [ ] 接入 ES 全文搜索
- [ ] 接入 OAuth2（GitHub/Google 登录）

## 许可证

MIT
