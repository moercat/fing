#!/usr/bin/env bash
# fing 一键启动脚本
# 用法：bash scripts/quickstart.sh
#
# 作用：
#   1. 复制 .env.example -> .env（如不存在）
#   2. docker compose up -d 拉起 MySQL + Redis + fing
#   3. 等服务就绪
#   4. 验证 /health 接口

set -euo pipefail

cd "$(dirname "$0")/.."
ROOT=$(pwd)

echo "==> fing quickstart"
echo ""

# 1. .env
if [ ! -f .env ]; then
  echo "[1/4] cp .env.example .env"
  cp .env.example .env
else
  echo "[1/4] .env exists, skip"
fi

# 2. docker compose
echo "[2/4] docker compose up -d"
docker compose up -d

# 3. wait for healthy
echo "[3/4] waiting for services..."
for i in $(seq 1 60); do
  if curl -sf http://localhost:9765/health >/dev/null 2>&1; then
    echo "    ready after ${i}s"
    break
  fi
  sleep 1
  if [ "$i" = "60" ]; then
    echo "    TIMEOUT: service not ready in 60s"
    echo "    run 'docker compose logs fing' to debug"
    exit 1
  fi
done

# 4. verify
echo "[4/4] verify /health"
RESP=$(curl -s http://localhost:9765/health)
echo "    response: $RESP"

echo ""
echo "✅ fing is running at http://localhost:9765"
echo ""
echo "Next steps:"
echo "  • API docs (if generated): http://localhost:9765/swagger/index.html"
echo "  • Logs:                     docker compose logs -f fing"
echo "  • Stop:                     docker compose down"
echo "  • Re-run quickstart:        bash scripts/quickstart.sh"