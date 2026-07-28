#!/usr/bin/env bash
# fing API 调用示例（cURL）
# 假设服务监听 http://localhost:9765

set -euo pipefail

BASE="http://localhost:9765"

echo "=== 1. 健康检查 ==="
curl -s "$BASE/health" | jq

echo ""
echo "=== 2. 注册 ==="
curl -s -X POST "$BASE/api/v1/register" \
  -H "Content-Type: application/json" \
  -d '{
    "user_name": "alice",
    "password": "alice12345",
    "re_password": "alice12345",
    "nickname": "Alice",
    "email": "alice@example.com"
  }' | jq

echo ""
echo "=== 3. 登录（Session） ==="
curl -s -c cookies.txt -X POST "$BASE/api/v1/login" \
  -H "Content-Type: application/json" \
  -d '{"user_name":"alice","password":"alice12345"}' | jq

echo ""
echo "=== 4. 登录（JWT） ==="
LOGIN_RESP=$(curl -s -X POST "$BASE/api/v1/login/jwt" \
  -H "Content-Type: application/json" \
  -d '{"user_name":"alice","password":"alice12345"}')
echo "$LOGIN_RESP" | jq

TOKEN=$(echo "$LOGIN_RESP" | jq -r '.data.token')
echo "JWT: $TOKEN"

echo ""
echo "=== 5. 用 JWT 访问受保护接口 ==="
curl -s -X GET "$BASE/api/v2/profile" \
  -H "Authorization: Bearer $TOKEN" | jq

echo ""
echo "=== 6. 修改资料 ==="
curl -s -X PUT "$BASE/api/v2/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"nickname":"Alice2"}' | jq

echo ""
echo "=== 7. 修改密码 ==="
curl -s -X PUT "$BASE/api/v2/profile/password" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"old_password":"alice12345","new_password":"newpass12345"}' | jq

echo ""
echo "=== 8. 申请密码重置 ==="
curl -s -X POST "$BASE/api/v1/password/forgot" \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com"}' | jq

# 上面这步会发邮件，邮件里有 token，把 token 复制到这里：
# TOKEN=$(...)
# curl -X POST "$BASE/api/v1/password/reset" \
#   -H "Content-Type: application/json" \
#   -d "{\"token\":\"$TOKEN\",\"new_password\":\"forgot12345\"}"
