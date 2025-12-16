# 使用官方 Golang 镜像作为基础镜像
FROM golang:1.19-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 使用官方 Alpine 镜像作为最终镜像
FROM alpine:latest

# 安装 ca-certificates，以便应用可以连接到 HTTPS 端点
RUN apk --no-cache add ca-certificates

# 创建非 root 用户
RUN adduser -D -s /bin/sh appuser

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件
COPY --from=builder /app/config.yaml .

# 更改文件所有者
RUN chown -R appuser:appuser /root/
USER appuser

# 暴露端口
EXPOSE 9765

# 启动应用
CMD ["./main"]