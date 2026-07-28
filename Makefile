# fing Makefile — 一键 build / test / run / deploy / release
#
# 常用命令：
#   make help       显示所有命令
#   make build      本地编译
#   make test       跑测试
#   make run        本地启动（dev 配置）
#   make docker     构建 docker 镜像
#   make up         docker-compose 一键拉起 MySQL+Redis+fing
#   make down       停止
#   make release    打 tag 触发自动发布
#   make snapshot   本地试发布（不发 GitHub）

GO          ?= go
BINARY      ?= fing
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT      ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE        ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS     ?= -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

# 是否安装 goreleaser（可选）
GORELEASER  ?= $(shell command -v goreleaser 2>/dev/null)

.PHONY: help
help: ## 显示帮助
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## 本地编译二进制
	$(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY) .

.PHONY: run
run: ## 本地启动（dev 环境）
	CONFIGOR_ENV=development $(GO) run -trimpath -ldflags "$(LDFLAGS)" .

.PHONY: test
test: ## 跑测试
	$(GO) test -race -coverprofile=coverage.out ./...

.PHONY: vet
vet: ## 静态检查
	$(GO) vet ./...

.PHONY: check
check: vet test ## vet + test

.PHONY: tidy
tidy: ## 整理依赖
	$(GO) mod tidy

.PHONY: clean
clean: ## 清理构建产物
	rm -f $(BINARY) coverage.out
	$(GO) clean -cache

# === Docker ===

.PHONY: docker
docker: ## 构建 docker 镜像
	docker build -t fing/fing:$(VERSION) -t fing/fing:latest .

.PHONY: docker-buildx
docker-buildx: ## 多架构构建（linux/amd64 + linux/arm64）
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t fing/fing:$(VERSION) \
		-t fing/fing:latest \
		--push .

.PHONY: up
up: ## docker-compose 一键启动（MySQL+Redis+fing）
	docker compose up -d

.PHONY: down
down: ## 停止
	docker compose down

.PHONY: logs
logs: ## 查看 fing 日志
	docker compose logs -f fing

.PHONY: restart
restart: ## 重启 fing 服务
	docker compose restart fing

# === 发布 ===

.PHONY: release
release: ## 打 tag 并推送（触发 GitHub Actions 自动发布）
	@if [ -z "$(TAG)" ]; then \
		echo "用法: make release TAG=v1.0.0"; \
		exit 1; \
	fi
	git tag -a $(TAG) -m "Release $(TAG)"
	git push origin $(TAG)

.PHONY: snapshot
snapshot: ## goreleaser 本地试发布（不推到 GitHub）
	@if [ -z "$(GORELEASER)" ]; then \
		echo "需要先安装 goreleaser：go install github.com/goreleaser/goreleaser@latest"; \
		exit 1; \
	fi
	goreleaser release --snapshot --clean