# fing Dockerfile — 多阶段构建，最终镜像只有二进制 + 配置
#
# 用法：
#   docker build -t fing-app .
#   docker compose up -d

# ---- 构建阶段 ----
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 单独缓存依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 静态链接，最终镜像小
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o fing .

# ---- 运行阶段 ----
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata && \
    adduser -D -s /bin/sh finguser

WORKDIR /app

COPY --from=builder /app/fing .
COPY --from=builder /app/config.*.yaml ./
COPY --from=builder /app/.env.example ./.env.example

RUN chown -R finguser:finguser /app
USER finguser

EXPOSE 9765

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:9765/health || exit 1

CMD ["./fing"]