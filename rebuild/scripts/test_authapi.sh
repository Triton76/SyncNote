#!/usr/bin/env bash
set -euo pipefail

AUTH_HOST="${AUTH_HOST:-http://127.0.0.1:8000}"
AUTH_EMAIL="${AUTH_EMAIL:-rebuild_auth_$(date +%s)@example.com}"
AUTH_PASSWORD="${AUTH_PASSWORD:-123456}"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Error: missing command '$1'"
    exit 1
  fi
}

need_cmd curl
need_cmd jq

echo "[1/3] Register on authapi: $AUTH_HOST"
register_resp="$({
  curl -sS -X POST "$AUTH_HOST/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$AUTH_EMAIL\",\"password\":\"$AUTH_PASSWORD\",\"captcha\":\"\"}"
} )"

echo "$register_resp" | jq . >/dev/null
user_id="$(echo "$register_resp" | jq -r '.userId // .UserId // empty')"
if [[ -z "$user_id" ]]; then
  echo "Register failed: $register_resp"
  exit 1
fi

echo "[2/3] Login on authapi"
login_resp="$({
  curl -sS -X POST "$AUTH_HOST/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$AUTH_EMAIL\",\"password\":\"$AUTH_PASSWORD\",\"captcha\":\"\"}"
} )"

echo "$login_resp" | jq . >/dev/null
token="$(echo "$login_resp" | jq -r '.token // .Token // empty')"
expire_in="$(echo "$login_resp" | jq -r '.expirein // .expireIn // .ExpireIn // 0')"
if [[ -z "$token" || "$expire_in" == "0" ]]; then
  echo "Login failed: $login_resp"
  exit 1
fi

echo "[3/3] Validate token format"
if [[ "$token" != *.*.* ]]; then
  echo "Invalid JWT token format: $token"
  exit 1
fi

echo "authapi smoke test passed"
echo "userId=$user_id"
echo "email=$AUTH_EMAIL"
