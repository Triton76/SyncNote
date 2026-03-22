#!/usr/bin/env bash
set -euo pipefail

# 服务地址
AUTH_HOST="${AUTH_HOST:-http://127.0.0.1:8000}"
USER_HOST="${USER_HOST:-http://127.0.0.1:8888}"

# 测试用的随机邮箱和密码 (避免重复注册报错)
TEST_EMAIL="${TEST_EMAIL:-test_$(date +%s)@syncnote.dev}"
TEST_PASS="${TEST_PASS:-SecurePass123!}"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "❌ 缺少命令: $1"
    exit 1
  fi
}

need_cmd curl
need_cmd jq

http_json_request() {
  local method="$1"
  local url="$2"
  local token="${3:-}"
  local body="${4:-}"

  local tmp_body
  tmp_body="$(mktemp)"

  local code
  if [ -n "$token" ]; then
    if [ -n "$body" ]; then
      code="$(curl -sS -o "$tmp_body" -w "%{http_code}" -X "$method" "$url" -H "Authorization: Bearer $token" -H "Content-Type: application/json" -d "$body")"
    else
      code="$(curl -sS -o "$tmp_body" -w "%{http_code}" -X "$method" "$url" -H "Authorization: Bearer $token")"
    fi
  else
    if [ -n "$body" ]; then
      code="$(curl -sS -o "$tmp_body" -w "%{http_code}" -X "$method" "$url" -H "Content-Type: application/json" -d "$body")"
    else
      code="$(curl -sS -o "$tmp_body" -w "%{http_code}" -X "$method" "$url")"
    fi
  fi

  local resp
  resp="$(cat "$tmp_body")"
  rm -f "$tmp_body"

  if [ "$code" -lt 200 ] || [ "$code" -ge 300 ]; then
    echo "❌ 请求失败: $method $url" >&2
    echo "HTTP状态码: $code" >&2
    if [ -n "$resp" ]; then
      echo "响应内容: $resp" >&2
    else
      echo "响应内容: <empty>" >&2
    fi
    return 1
  fi

  if ! echo "$resp" | jq . >/dev/null 2>&1; then
    echo "❌ 返回非 JSON: $method $url" >&2
    echo "HTTP状态码: $code" >&2
    if [ -n "$resp" ]; then
      echo "响应内容: $resp" >&2
    else
      echo "响应内容: <empty>" >&2
    fi
    return 1
  fi

  printf "%s" "$resp"
}

echo "🚀 开始测试 SyncNote 重建版逻辑..."
echo "测试邮箱: $TEST_EMAIL"

echo -e "\n--- [1/6] 测试注册 (RegisterLogic) ---"
REGISTER_RESP="$(http_json_request "POST" "$AUTH_HOST/auth/register" "" "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASS\",\"captcha\":\"\"}")"
echo "响应内容: $REGISTER_RESP"

USER_ID="$(echo "$REGISTER_RESP" | jq -r '.userId // empty')"
if [ -z "$USER_ID" ]; then
  echo "❌ 注册失败: 未返回 userId"
  echo "详细错误: $REGISTER_RESP"
  exit 1
fi
echo "✅ 注册成功: UserId = $USER_ID"

echo -e "\n--- [2/6] 测试登录 (LoginLogic) ---"
LOGIN_RESP="$(http_json_request "POST" "$AUTH_HOST/auth/login" "" "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASS\",\"captcha\":\"\"}")"
echo "响应内容: $LOGIN_RESP"

TOKEN="$(echo "$LOGIN_RESP" | jq -r '.token // empty')"
if [ -z "$TOKEN" ]; then
  echo "❌ 登录失败: 未返回 token"
  echo "详细错误: $LOGIN_RESP"
  exit 1
fi
echo "✅ 登录成功: Token 已获取 (长度: ${#TOKEN})"

echo -e "\n--- [3/6] 测试获取自身信息 (GetSelfInfoLogic) ---"
SELF_INFO_RESP="$(http_json_request "GET" "$USER_HOST/api/user/me" "$TOKEN")"
echo "响应内容: $SELF_INFO_RESP"

RET_USER_ID="$(echo "$SELF_INFO_RESP" | jq -r '.userId // empty')"
RET_USERNAME="$(echo "$SELF_INFO_RESP" | jq -r '.username // empty')"
if [ "$RET_USER_ID" = "$USER_ID" ] && [ -n "$RET_USERNAME" ]; then
  echo "✅ 获取自身信息成功: Username = $RET_USERNAME"
else
  echo "❌ 获取自身信息失败: ID不匹配或字段缺失"
  echo "期望 ID: $USER_ID, 返回 ID: $RET_USER_ID"
  exit 1
fi

echo -e "\n--- [4/6] 测试修改自身信息 (EditSelfInfoLogic) ---"
NEW_BIO="这是通过 API 测试更新的简介"
NEW_AVATAR="http://example.com/new_avatar.png"
EDIT_RESP="$(http_json_request "POST" "$USER_HOST/api/user/edit/me" "$TOKEN" "{\"username\":\"UpdatedUser_$USER_ID\",\"synopsis\":\"$NEW_BIO\",\"avatarUrl\":\"$NEW_AVATAR\"}")"
echo "响应内容: $EDIT_RESP"

if [ "$EDIT_RESP" != "{}" ]; then
  echo "⚠️ Edit 接口返回非空对象，继续校验读取结果"
fi

VERIFY_RESP="$(http_json_request "GET" "$USER_HOST/api/user/me" "$TOKEN")"
VERIFY_BIO="$(echo "$VERIFY_RESP" | jq -r '.synopsis // empty')"
if [ "$VERIFY_BIO" = "$NEW_BIO" ]; then
  echo "✅ 数据验证成功: 简介已更新为 '$NEW_BIO'"
else
  echo "⚠️ 警告: 修改请求成功但数据未同步 (当前简介: $VERIFY_BIO)"
fi

echo -e "\n--- [5/6] 测试搜索用户 (SearchUserLogic) ---"
SEARCH_RESP="$(http_json_request "GET" "$USER_HOST/api/user/search?email=$TEST_EMAIL" "$TOKEN")"
echo "响应内容: $SEARCH_RESP"

FOUND_EMAIL="$(echo "$SEARCH_RESP" | jq -r '.infoList[0].email // empty')"
TOTAL="$(echo "$SEARCH_RESP" | jq -r '.total // 0')"
if [ "$FOUND_EMAIL" = "$TEST_EMAIL" ] && [ "$TOTAL" = "1" ]; then
  echo "✅ 搜索用户成功: 找到目标邮箱"
else
  echo "❌ 搜索用户失败"
  exit 1
fi

echo -e "\n--- [6/6] 测试获取指定用户信息 (GetUserInfoLogic) ---"
GET_USER_RESP="$(http_json_request "GET" "$USER_HOST/api/user/$USER_ID/info" "$TOKEN")"
echo "响应内容: $GET_USER_RESP"

RESP_ID="$(echo "$GET_USER_RESP" | jq -r '.userId // empty')"
if [ "$RESP_ID" = "$USER_ID" ]; then
  echo "✅ 获取指定用户信息成功"
else
  echo "❌ 获取指定用户信息失败"
  exit 1
fi

echo -e "\n🎉 脚本执行完成"
