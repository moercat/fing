# API 文档

## 基本信息

- 基础URL: `http://localhost:9765`
- 内容类型: `application/json`

## 公共接口

### 健康检查
- **GET** `/health`
- 检查服务的运行状态
- 响应示例：
```json
{
  "code": 200,
  "data": {
    "status": "healthy",
    "checks": {
      "database": {
        "status": "up",
        "duration": 2
      },
      "redis": {
        "status": "up",
        "duration": 1
      }
    },
    "timestamp": 1678886400
  },
  "msg": "服务健康"
}
```

## 用户接口 (V1)

### 检查服务状态
- **GET** `/api/v1/ping`
- 检查 API 服务是否正常运行
- 响应示例：
```json
{
  "code": 200,
  "data": "",
  "msg": "pong"
}
```

### 用户注册
- **POST** `/api/v1/register`
- 注册新用户
- 请求体：
```json
{
  "nickname": "用户名",
  "user_name": "唯一用户名",
  "password": "密码",
  "re_password": "确认密码"
}
```
- 响应示例：
```json
{
  "code": 200,
  "data": null,
  "msg": ""
}
```

### 用户登录
- **POST** `/api/v1/login`
- 用户登录
- 请求体：
```json
{
  "user_name": "用户名",
  "password": "密码"
}
```
- 响应示例：
```json
{
  "code": 200,
  "data": {
    "id": 1,
    "user_name": "用户名",
    "nickname": "昵称",
    "status": 0,
    "avatar": "",
    "created_at": 1678886400
  },
  "msg": ""
}
```

## 用户接口 (V2) - 需要认证

### 获取用户信息
- **GET** `/api/v2/user_info`
- 需要有效的会话
- 响应示例：
```json
{
  "code": 200,
  "data": {
    "id": 1,
    "user_name": "用户名",
    "nickname": "昵称",
    "status": 0,
    "avatar": "",
    "created_at": 1678886400
  },
  "msg": ""
}
```

### 用户登出
- **DELETE** `/api/v2/logout`
- 清除用户会话
- 响应示例：
```json
{
  "code": 200,
  "data": null,
  "msg": "登出成功"
}
```

## 错误码

- `200`: 成功
- `400`: 请求参数错误
- `401`: 未授权访问
- `403`: 禁止访问
- `404`: 资源不存在
- `422`: 参数校验失败
- `500`: 服务器内部错误