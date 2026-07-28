# fing

> Gin + GORM + Redis 的 Go 后端脚手架。**Fork 之后改几行就能跑。**

![Go Version](https://img.shields.io/badge/Go-1.21-blue)
![Gin](https://img.shields.io/badge/Gin-1.10-lightgrey)
![Gorm](https://img.shields.io/badge/Gorm-gorm.io-red)
![License](https://img.shields.io/badge/License-MIT-green)
[![CI](https://github.com/moercat/fing/actions/workflows/go.yml/badge.svg)](https://github.com/moercat/fing/actions/workflows/go.yml)
![Stars](https://img.shields.io/badge/Stars-46-brightgreen)

## 🚀 Fork 之后 30 秒启动

### 方式一：Docker（推荐，零配置）

```bash
git clone https://github.com/<you>/fing.git myapp
cd myapp
cp .env.example .env
docker compose up -d
curl http://localhost:9765/health
# → {"code":0,"msg":"ok","data":""}
```

自动启动：
- MySQL 8.0（端口 3306，自动建库 fing）
- Redis 7（端口 6379）
- fing 服务（端口 9765）
- 自动 `AutoMigrate` 建表

### 方式二：本地开发

```bash
# 1. 启动 MySQL + Redis（你已有，或用 docker compose up -d mysql redis）
# 2. 准备数据库
mysql -u root -p -e "CREATE DATABASE fing DEFAULT CHARSET utf8mb4"

# 3. 启动 fing
go mod download
go run main.go
```

## 📋 fork 之后需要改哪些地方

打开 IDE 全局搜索替换：

| 改什么 | 在哪里 | 改成什么 |
| --- | --- | --- |
| 模块名 | `go.mod` 第 1 行 `module fing` | `module myapp` |
| 数据库名 | `.env.example` 和 `config.*.yaml` | 你的库名 |
| 业务端口 | `.env.example` `APP_PORT` | 你的端口 |
| JWT 密钥 | `.env.example` `APP_SECRET` | 32+ 位随机字符串 |
| 业务模块 | `internal/apis/login`、`internal/service/login` | 你的业务名 |
| Dockerfile 镜像名 | 最后一行 `docker build -t fing-app` | `docker build -t myapp` |
| docker-compose 服务名 | `services.fing` | `services.myapp` |

**仅此而已**，剩下的脚手架代码都不用动。

## ✅ 内置能力（开箱即用）

| 模块 | 接口 | 说明 |
| --- | --- | --- |
| 健康检查 | `GET /health` | 验证服务存活 |
| 注册 | `POST /api/v1/register` | 用户名+邮箱+密码 |
| 登录 | `POST /api/v1/login` | Session 模式（Web） |
| 登录 JWT | `POST /api/v1/login/jwt` | 移动端 / 前后端分离 |
| 密码重置申请 | `POST /api/v1/password/forgot` | 发邮件 |
| 密码重置 | `POST /api/v1/password/reset` | 用 token 重置 |
| 邮箱验证 | `POST /api/v1/email/verify` | 验证邮箱 |
| 当前用户 | `GET /api/v2/user_info` | Session |
| 资料 | `GET / PUT /api/v2/profile` | 昵称/邮箱 |
| 修改密码 | `PUT /api/v2/profile/password` | 登录态改密码 |
| 修改头像 | `PUT /api/v2/profile/avatar` | URL 形式 |
| 用户列表 | `GET /api/v2/users` | admin |
| 修改角色 | `PUT /api/v2/users/:id/role` | admin |
| 登出 | `DELETE /api/v2/logout` | 清 Session |

中间件齐全：`CORS` / `Session` / `AuthRequired` / `AdminRequired` / `LoggerToFile` / `Recover` / `Cover`

## 🧩 内置能力（脚手架自带）

| 能力 | 在哪里 | 怎么用 |
| --- | --- | --- |
| MySQL（GORM） | `pkg/db.Gain` | `db.Gain.Where(...).Find(&users)` |
| MySQL（XORM） | `pkg/db.Main` | 兼容老代码用，新代码用 GORM |
| Redis | `pkg/db.RedisClient` | `db.RedisClient.Set(key, val, ttl).Err()` |
| 邮件发送 | `pkg/email.SendMail(to, name, subj, body)` | 一行调用 |
| JWT | `pkg/jwt.Sign(secret, uid, name, role, ttl)` | 签发 + 解析 |
| 定时任务 | `cobra.Register(name, interval, fn)` | 注册式 |
| 配置加载 | `pkg/config.Config` | YAML + 环境变量 |
| 统一响应 | `resp.OK(c, data, msg)` / `resp.Fail(c, err, msg)` | 全局 |
| 错误包装 | `errors.New(422, "msg")` / `errors.Wrap(err, 500, "msg")` | 带错误码 |

## 📁 项目结构（fork 后照着加业务就行）

```
fing/
├── main.go                  ← 入口（AutoMigrate + 注册任务 + 启动服务）
├── config.{dev,docker,prod}.yaml  ← 3 套环境配置
├── .env.example             ← 环境变量模板
├── Dockerfile               ← 多阶段构建
├── docker-compose.yml       ← 一键拉起 MySQL + Redis + fing
├── Makefile                 ← make build / run / test / docker-up
├── APIDOC.md                ← 接口详情
├── examples/client.go       ← Go 客户端调用示例（//go:build ignore）
├── internal/
│   ├── router.go            ← 路由注册（公共 / 业务 / admin 分组）
│   ├── apis/                ← HTTP handler 层（薄）
│   │   ├── login/           ← 注册 / 登录 / JWT
│   │   ├── user/            ← 资料 / 密码 / 角色
│   │   └── password/        ← 密码重置 / 邮箱验证
│   ├── service/             ← 业务逻辑层（厚）
│   │   ├── login/
│   │   ├── user/
│   │   └── password/
│   ├── model/               ← HTTP 请求/响应结构体
│   └── tools/               ← Session 等工具
└── pkg/                     ← 可复用工具包（业务也可 import）
    ├── jwt/                 ← JWT 签发/解析
    ├── middleware/          ← Auth/CORS/Logger/Recover
    ├── cobra/               ← 定时任务（注册式）
    ├── db/                  ← GORM + XORM + Redis 客户端
    ├── email/               ← 邮件
    ├── config/              ← 配置加载
    ├── entity/usr/          ← 数据模型
    ├── errors/              ← 错误包装
    ├── resp/                ← 统一响应
    └── health/              ← 健康检查
```

## 🛠️ 加一个业务模块（比如文章）

照着 `internal/apis/login` 复制一个：

```bash
# 1. 加数据模型
cp pkg/entity/usr/user.go pkg/entity/article/article.go
# 改字段

# 2. 加业务逻辑
mkdir -p internal/service/article
cp -r internal/service/login/. internal/service/article/
# 改函数名

# 3. 加 HTTP handler
mkdir -p internal/apis/article
cp -r internal/apis/login/. internal/apis/article/
# 改路由

# 4. 注册到 router
# internal/router.go 加：
#   new(article.RouterArticle).Router(r)

# 5. 在 main.go 加 AutoMigrate
#   db.Gain.AutoMigrate(&article.Article{})
```

## 🧪 怎么验证 fork 后的项目

```bash
# 1. 启动
docker compose up -d

# 2. 健康检查
curl http://localhost:9765/health

# 3. 注册测试用户
curl -X POST http://localhost:9765/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"user_name":"alice","password":"alice12345","re_password":"alice12345","nickname":"Alice","email":"alice@example.com"}'

# 4. 登录拿 JWT
curl -X POST http://localhost:9765/api/v1/login/jwt \
  -H "Content-Type: application/json" \
  -d '{"user_name":"alice","password":"alice12345"}'

# 5. 调受保护接口
TOKEN="..."  # 上一步返回的 token
curl -H "Authorization: Bearer $TOKEN" http://localhost:9765/api/v2/profile
```

## ⚙️ 配置项

通过 `configor` 自动加载：
- 有 `CONFIGOR_ENV=docker` → `config.docker.yaml`
- 没设或 `dev` → `config.development.yaml`
- 有 `CONFIGOR_ENV=production` → `config.production.yaml`

YAML 里的字段会被环境变量覆盖（同名大写）：
```yaml
port: 9765      # APP_PORT=8888 会覆盖
```

## 📜 许可证

MIT