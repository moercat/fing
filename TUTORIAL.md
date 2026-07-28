# fing 后端开发教程

> 本教程帮助新人通过 fing 项目**学会 Go 后端开发的基本思路**。
> 跟着这个顺序读 + 写，2-3 天能搭一个完整后端。

## 🎯 学习目标

读完后你应该会：

- 知道后端项目的标准目录结构（API 层 / Service 层 / Model 层 / 工具层）
- 会写 RESTful API（Gin）
- 会操作 MySQL（GORM 为主，XORM 备用）
- 会用 Redis 做缓存 / 限流 / 分布式锁
- 会注册和写定时任务
- 会发邮件
- 会用 JWT / Session 做用户鉴权
- 会做 RBAC 权限控制

## 📚 推荐阅读顺序

### 第 1 步：理解项目结构（30 分钟）

看 [README.md](README.md) 的「目录结构」一节，跑一下项目：
```bash
docker compose up -d
go run main.go
curl http://localhost:9765/health
```

### 第 2 步：跟一个完整 API（1 小时）

跟一个最简单的 `register` API，从路由到数据库走一遍：

```
HTTP 请求
    ↓
internal/router.go          ← 路由注册：URL 映射到 handler
    ↓
internal/apis/login/login.go      ← handler：解析参数、调用 service、返回响应
    ↓
internal/service/login/register.go ← service：业务逻辑（密码加密、查重、写库）
    ↓
pkg/entity/usr/user.go       ← entity：数据模型（GORM tag 决定表结构）
    ↓
pkg/db/db.go                ← db：GORM 连接
    ↓
MySQL
```

**动手**：仿照 `register` 写一个 `logout-all-devices` API（注销所有设备）。

### 第 3 步：学数据库（GORM）— [examples/db/mysql_gorm.go](examples/db/mysql_gorm.go)

**必看**：
- `exampleCreate` / `exampleRead` / `exampleUpdate` / `exampleDelete` — 基础 CRUD
- `exampleWhere` — 各种查询
- `exampleTransaction` — 事务（最容易出错的地方）
- `examplePagination` — 分页（业务必备）

**练手**：给 user 表加一个 `last_login_at` 字段，记录用户最近登录时间。

### 第 4 步：学缓存（Redis）— [examples/db/redis_basic.go](examples/db/redis_basic.go)

**必看**：
- `exampleString` / `exampleHash` — 基础数据类型
- `exampleCounter` — INCR 用法（限流、点赞）
- `exampleLock` — 分布式锁
- `getUserWithCache` — Cache-Aside 缓存模式

**练手**：给 user_info 接口加 5 分钟 Redis 缓存。

### 第 5 步：学定时任务（cobra）— [examples/task/daily_stats.go](examples/task/daily_stats.go)

**核心改动**（已自动应用，见 [pkg/cobra/cobra.go](pkg/cobra/cobra.go)）：

老的写法（不推荐）：
```go
func Cobra() {
    go Ticker()  // 写死的 ticker
}
```

新的写法（推荐）：
```go
// 1. 在 examples/task/ 注册
cobra.Register("daily-stats", 24*time.Hour, func() {
    // 你的业务代码
})

// 2. main.go 自动启动
cobra.Cobra()
```

**练手**：写一个 `purge-old-sessions` 任务，每 10 分钟清理一次过期 session。

### 第 6 步：学邮件（email）— [examples/email/send.go](examples/email/send.go)

**核心**：业务代码调 `email.SendMail(to, name, subject, htmlBody)`，邮件内容用 `html/template` 渲染。

**练手**：注册成功时给用户发一封欢迎邮件。

### 第 7 步：学鉴权（JWT + Session）— [pkg/jwt/jwt.go](pkg/jwt/jwt.go)

**两条路线**：
- Web 端用 **Session**（`internal/router.go` 里 `middleware.Session()`）
- 移动端 / 前后端分离用 **JWT**（`/api/v1/login/jwt` 接口）

**练手**：给 `/api/v2/profile` 加 JWT 鉴权（参考 `/api/v1/login/jwt` 流程）。

### 第 8 步：学权限（RBAC）— [pkg/middleware/auth.go](pkg/middleware/auth.go)

**两种中间件**：
- `AuthRequired()` — 必须登录
- `AdminRequired()` — 必须 admin 角色

**用法**：
```go
v2 := router.Group("/api/v2")
v2.Use(middleware.AuthRequired()) {
    // 需要登录的接口
    admin := v2.Group("")
    admin.Use(middleware.AdminRequired()) {
        // 需要 admin 的接口
    }
}
```

**练手**：加一个"用户禁用"接口（admin 才能调用）。

## 🗺️ 在 fing 项目里"东西都放哪里"

| 你要做的事 | 文件位置 | 参考示例 |
| --- | --- | --- |
| 加一个 API | `internal/apis/<模块>/` + `internal/service/<模块>/` | `internal/apis/login/login.go` |
| 加一个数据表 | `pkg/entity/<模块>/` | `pkg/entity/usr/user.go` |
| 写 SQL（CRUD） | 在 `internal/service/` 里用 `db.Gain.*` | [examples/db/mysql_gorm.go](examples/db/mysql_gorm.go) |
| 用 Redis | 在 `internal/service/` 里用 `db.RedisClient.*` | [examples/db/redis_basic.go](examples/db/redis_basic.go) |
| 加定时任务 | `examples/task/` 或 `pkg/task/` | [examples/task/daily_stats.go](examples/task/daily_stats.go) |
| 发邮件 | `internal/service/` 里用 `email.SendMail` | [examples/email/send.go](examples/email/send.go) |
| 鉴权中间件 | `pkg/middleware/` | `pkg/middleware/auth.go` |
| 工具函数 | `pkg/<功能名>/` | `pkg/jwt/`、`pkg/errors/` |
| 启动初始化 | `main.go` | `main.go` |
| 配置项 | `pkg/config/config.go` + `config.{env}.yaml` | `config.development.yaml` |

## 🧪 怎么测试

```bash
# 单元测试
go test ./...

# 跑起来手动测
go run main.go

# 用 examples/curl.sh 一把梭
bash examples/curl.sh
```

## 🚀 完整练手项目（建议顺序）

按这个顺序做 5 个练习，就掌握了 80% 的后端开发套路：

1. **加一个"用户禁用"接口**（admin 权限）— 练路由 + 鉴权 + service + DB
2. **加 Redis 缓存 user_info** — 练 Redis + 缓存模式
3. **加定时清理过期 token** — 练定时任务
4. **加注册欢迎邮件** — 练邮件 + html/template
5. **加 JWT 登录** — 练 JWT

每个练习 1-2 小时，全部做完 ≈ 1 周成为后端熟手。

## 📖 进阶阅读

- [Go 官方 Effective Go](https://go.dev/doc/effective_go)
- [Gin 官方文档](https://gin-gonic.com/zh-cn/docs/)
- [GORM 官方文档](https://gorm.io/zh_CN/docs/)
- [Redis 命令参考](https://redis.io/commands/)
- [JWT 介绍](https://jwt.io/introduction)

## ❓ 常见问题

### Q: 业务代码应该放 apis/ 还是 service/？

**A**:
- `apis/` 只做：解析参数、调用 service、构造响应。**不超过 30 行**。
- `service/` 做：业务逻辑、数据库操作、调用其他 service。
- 业务逻辑复杂时，service 可以再拆 service。

### Q: entity 和 model 有什么区别？

**A**:
- `pkg/entity/` 是数据库表结构（GORM tag）
- `internal/model/` 是 HTTP 请求/响应的结构体（JSON tag + 校验 tag）
- 两者通过 service 层做转换

### Q: 什么时候用 Redis？

**A**:
- ✅ 缓存读多写少的数据
- ✅ 计数器（views / likes）
- ✅ 分布式锁
- ✅ Session 存储
- ❌ 不要把 Redis 当主存（重启会丢）

### Q: GORM 和 XORM 选哪个？

**A**: 默认用 GORM（社区活跃、文档全）。XORM 留着兼容老代码。
