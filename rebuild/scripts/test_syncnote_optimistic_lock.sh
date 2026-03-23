#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
AUTH_HOST="${AUTH_HOST:-http://127.0.0.1:8000}"
RPC_ADDR="${RPC_ADDR:-127.0.0.1:8080}"
AUTH_EMAIL="${AUTH_EMAIL:-rebuild_syncnote_lock_$(date +%s)@example.com}"
AUTH_PASSWORD="${AUTH_PASSWORD:-123456}"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Error: missing command '$1'"
    exit 1
  fi
}

need_cmd curl
need_cmd jq
need_cmd go

wait_http_ready() {
  local url="$1"
  local retries="${2:-20}"
  for _ in $(seq 1 "$retries"); do
    if curl -sS -m 1 "$url" >/dev/null 2>&1; then
      return 0
    fi
    sleep 0.3
  done
  return 1
}

wait_tcp_ready() {
  local host="$1"
  local port="$2"
  local retries="${3:-20}"
  for _ in $(seq 1 "$retries"); do
    if timeout 1 bash -c "</dev/tcp/$host/$port" >/dev/null 2>&1; then
      return 0
    fi
    sleep 0.3
  done
  return 1
}

rpc_host="${RPC_ADDR%%:*}"
rpc_port="${RPC_ADDR##*:}"

if ! wait_http_ready "$AUTH_HOST/auth/login" 10; then
  echo "authapi is not reachable at $AUTH_HOST"
  echo "Please start rebuild/authapi first."
  exit 1
fi

if ! wait_tcp_ready "$rpc_host" "$rpc_port" 10; then
  echo "syncnote rpc is not listening at $RPC_ADDR"
  echo "Please start rebuild/syncnote/rpc first."
  exit 1
fi

echo "[1/3] Register user via authapi"
register_resp="$({
  curl -sS -X POST "$AUTH_HOST/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$AUTH_EMAIL\",\"password\":\"$AUTH_PASSWORD\",\"captcha\":\"\"}"
})"

echo "$register_resp" | jq . >/dev/null
user_id="$(echo "$register_resp" | jq -r '.userId // .UserId // empty')"
if [[ -z "$user_id" ]]; then
  echo "Register failed: $register_resp"
  exit 1
fi

echo "[2/3] Login via authapi (sanity check)"
login_resp="$({
  curl -sS -X POST "$AUTH_HOST/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$AUTH_EMAIL\",\"password\":\"$AUTH_PASSWORD\",\"captcha\":\"\"}"
})"

echo "$login_resp" | jq . >/dev/null
token="$(echo "$login_resp" | jq -r '.token // .Token // empty')"
if [[ -z "$token" ]]; then
  echo "Login failed: $login_resp"
  exit 1
fi

echo "[3/3] Run syncnote optimistic lock smoke test"
(
  cd "$ROOT_DIR"
  USER_ID="$user_id" RPC_ADDR="$RPC_ADDR" go run ./rebuild/scripts/syncnote_optimistic_lock_smoke.go
)

echo "syncnote optimistic lock script passed"
echo "userId=$user_id"
echo "email=$AUTH_EMAIL"