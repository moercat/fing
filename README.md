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

## 目的

本项目采用了一系列Golang中比较流行的组件，可以以本项目为基础快速搭建Restful Web API

## 特色
看到 gourouting 大佬的 singo，便想建立一个自己的 web 项目，基于大佬的 singo 项目进行重构

[Singo](https://github.com/gourouting/singo): 使用Singo开发Web服务: 用最简单的架构，实现够用的框架，服务海量用户

整合了开发API最基本的组件：

1. [Gin](https://github.com/gin-gonic/gin): 轻量级Web框架，Gin 是一个用 Go (Golang) 编写的 HTTP Web 框架。它具有类似 Martini 的 API，性能要好得多——速度提高了 40 倍.
2. [GORM](https://gorm.io/index.html): ORM工具。本项目需要配合Mysql使用
3. [XORM](https://xorm.io/index.html): ORM工具。本项目需要配合Mysql使用
4. [Gin-Session](https://github.com/gin-contrib/sessions): Gin框架提供的Session操作工具
5. [Go-Redis](https://github.com/go-redis/redis): Golang Redis客户端
6. [Configor](https://github.com/jinzhu/configor): 支持 YAML、JSON、TOML、Shell 环境的 Golang 配置工具（支持 Go 1.10+）
7. [Gin-Cors](https://github.com/gin-contrib/cors): Gin框架提供的跨域中间件
8. [Elastic](https://github.com/olivere/elastic): Elastic 是Go编程语言的Elasticsearch客户端 。
9. [Gomail](https://github.com/go-gomail/gomail): 在 Go 中发送电子邮件的最佳方式。

本项目已经预先实现了一些常用的代码方便参考和复用:

1. 用户模型 ```pkg/entity/usr/user.go```
2. 实现了```/api/v1/register```用户注册接口
3. 实现了```/api/v1/login```用户登录接口

以下模块需要需要登录后获取session:

4. 实现了```/api/v2/user_info```用户资料接口(需要登录后获取session)
5. 实现了```/api/v2/logout```用户登出接口(需要登录后获取session)

本项目已经预先创建了一系列文件夹划分出下列模块:

1. ```internal``` 文件夹就是每一个项目的业务层所在区域，所有业务都实现在该文件夹内

   1.```apis``` 文件夹就是接口文件位置，相当于其他框架中的controller，所有接口接参都实现在该文件夹内，每个模块单独一个文件夹，避免不同模块耦合度过高

   2.```model``` 文件夹就是模型文件位置，所有接参模型与返回模型都实现在该文件夹内，实现与业务代码的解耦

   3.```service``` 文件夹就是业务文件位置，所有业务逻辑都实现在该文件夹内，每个模块单独一个文件夹，避免不同模块耦合度过高

   4.```tools``` 文件夹就是工具文件位置，当前项目业务内使用的工具文件

2. ```log``` 文件夹负责日志记录，与之后可扩充的业务监控与其他流程追踪业务

3. ```pkg``` 文件夹实现所有项目的功能、中间件、定时任务、数据库链接等等

   1.```cobra``` 文件夹就是定时任务文件位置，通过 time.NewTicker 实现定时任务功能

   2.```config``` 文件夹就是配置文件位置，通过 configor 来获取 yaml 配置文件中的文件

   3.```db``` 文件夹就是数据库文件位置，gorm、xorm、redis、es 等数据库链接初始化

   4.```elastic``` 文件夹就是 es 操作文件位置，简单实现es 中用户的增删改查功能

   5.```email``` 文件夹就是 gomail 的实现文件位置，通过 gomail 实现邮件的发送

   6.```entity``` 文件夹就是数据库实体文件位置，通过每一个数据库一个文件夹来清晰化项目结构

   7.```middleware``` 文件夹就是中间件文件位置，cors 跨域、session、recover 错误捕获、auth 登陆验证

   8.```resp``` 文件夹就是自定义返回结构文件位置，实现通用成功返回、失败返回、分页返回、错误校验

## Configor

项目在启动的时候依赖于 ```configor```，根据环境变量加载配置文件 使用CONFIGOR_ENV 设置环境变量， 如果 CONFIGOR_ENV 没有设置，框架将会使用 development 作为默认环境变量，也就是 读取 ```config.development.yaml```

如果我们要区分开发环境与生产环境，通过下方命令实现读取 **config.production.yaml**
`
CONFIGOR_ENV=production go run main.go
`
yaml 文件配置

```shell
mode: dev  // 当前所在环境
secret: asdasggas // session 加密密钥
level: 4 // 日志等级
port: 9765  // 端口

dataSource:
  main: you_name:you_password@tcp(you_ip:you_port)/db_name?charset=utf8mb4 // 数据库连接

redis:
  addr: 127.0.0.1:6379  // redis 连接

es:  // es 连接所需配置
  esUrl: you_es 
  esUsername: ""
  esPassword: ""

email: // 发送邮件 email 所需配置
  host: smtp.exmail.qq.com
  name: my_name
  email: my_email
  password: ""
```

## Go Mod

本项目使用[Go Mod](https://github.com/golang/go/wiki/Modules)管理依赖。

```shell
go mod tidy // 自动管理
go run main.go 
```

## 运行

```shell
go run main.go
```

## 部署

***linux***

```shell
GOOS=linux GOARCH=amd64 go build -o admin
```

项目运行后启动在9765端口（可以修改，参考gin文档)
