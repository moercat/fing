# fing

> 🎯 **后端新手的第一个 Go 项目** — 模块齐全 + 易懂 + Fork 即可跑

> 一个 Gin + GORM + Redis 后端脚手架。每个模块（登录 / 资料 / 密码重置 / 邮箱验证 / RBAC / 定时任务）都是教科书式的"小型完整示例"。新人 fork 之后改 6 处就能变成自己的项目。

![Go Version](https://img.shields.io/badge/Go-1.21-blue)
![Gin](https://img.shields.io/badge/Gin-1.10-lightgrey)
![Gorm](https://img.shields.io/badge/Gorm-gorm.io-red)
![License](https://img.shields.io/badge/License-MIT-green)
[![CI](https://github.com/moercat/fing/actions/workflows/go.yml/badge.svg)](https://github.com/moercat/fing/actions/workflows/go.yml)
[![Release](https://github.com/moercat/fing/actions/workflows/release.yml/badge.svg)](https://github.com/moercat/fing/releases)
![Stars](https://img.shields.io/badge/Stars-46-brightgreen)
![Docker](https://img.shields.io/badge/Docker-ghcr.io-blue)

## ✨ 为什么是 fing

- 🎓 **新手友好** — 每个模块 ≤ 100 行，注释清楚每步在干什么
- 📦 **模块齐全** — 12 个用户中心必备 API，开箱即用
- 🚀 **Fork 即可跑** — `docker compose up -d` 30 秒启动
- 🏗️ **架构清晰** — API / Service / Entity 三层分离，照着加业务就行
- 🔐 **鉴权完备** — Session + JWT 双模式 + RBAC 角色权限

## 🚀 30 秒启动

```bash
git clone https://github.com/<you>/fing.git myapp
cd myapp
bash scripts/quickstart.sh   # 一键拉起 + 验证
```

或手动：

```bash
cp .env.example .env
docker compose up -d
curl http://localhost:9765/health
```

## 📚 内置模块（12 个 API）

| 模块 | 接口 |
| --- | --- |
| 注册 | `POST /api/v1/register` |
| 登录（Session） | `POST /api/v1/login` |
| 登录（JWT） | `POST /api/v1/login/jwt` |
| 密码重置申请 | `POST /api/v1/password/forgot` |
| 密码重置 | `POST /api/v1/password/reset` |
| 邮箱验证 | `POST /api/v1/email/verify` |
| 当前用户 | `GET /api/v2/user_info` |
| 资料 | `GET PUT /api/v2/profile` |
| 修改密码 | `PUT /api/v2/profile/password` |
| 修改头像 | `PUT /api/v2/profile/avatar` |
| 用户列表（admin） | `GET /api/v2/users` |
| 修改角色（admin） | `PUT /api/v2/users/:id/role` |

## 🧩 内置能力

| 能力 | 用途 |
| --- | --- |
| GORM（MySQL） | 主存储 |
| Redis | 缓存 / Session / 限流 |
| JWT | 移动端鉴权 |
| 邮件 | 密码重置 / 通知 |
| 定时任务（注册式） | 周期任务 |
| 中间件 | CORS / Logger / Recover / Auth / Admin |
| 配置 | YAML + 环境变量 |
| 统一响应 / 错误码 | 全局封装 |

## 📁 项目结构（核心 3 层）

```
fing/
├── main.go                  ← 5 步启动入口
├── internal/
│   ├── apis/<模块>/         ← HTTP handler（薄）
│   ├── service/<模块>/      ← 业务逻辑（厚）
│   └── model/               ← 请求/响应结构体
└── pkg/                     ← 可复用工具
    ├── db/                  ← GORM + Redis
    ├── jwt, middleware, email, cobra, config, ...
    └── entity/usr/          ← 数据模型
```

**加业务模块**（比如文章）：
```bash
# 1. 复制 login 模块，改名字
cp -r internal/apis/login internal/apis/article
cp -r internal/service/login internal/service/article
cp pkg/entity/usr/user.go pkg/entity/article/article.go
# 2. 在 router.go 注册，在 main.go AutoMigrate
# 完成
```

详细 fork 指引见 [README §fork 改什么](#📋-fork-之后要改哪些地方)。