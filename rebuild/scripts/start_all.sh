#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
RUN_DIR="$ROOT_DIR/rebuild/.run"
LOG_DIR="$RUN_DIR/logs"
INIT_ENV_SCRIPT="$ROOT_DIR/rebuild/scripts/init_env.sh"

AUTH_PID_FILE="$RUN_DIR/authapi.pid"
USER_RPC_PID_FILE="$RUN_DIR/userrpc.pid"
USER_API_PID_FILE="$RUN_DIR/userapi.pid"
SYNC_RPC_PID_FILE="$RUN_DIR/syncnote-rpc.pid"
SYNC_API_PID_FILE="$RUN_DIR/syncnote-api.pid"
FRONTEND_PID_FILE="$RUN_DIR/frontend.pid"

AUTH_LOG_FILE="$LOG_DIR/authapi.log"
USER_RPC_LOG_FILE="$LOG_DIR/userrpc.log"
USER_API_LOG_FILE="$LOG_DIR/userapi.log"
SYNC_RPC_LOG_FILE="$LOG_DIR/syncnote-rpc.log"
SYNC_API_LOG_FILE="$LOG_DIR/syncnote-api.log"
FRONTEND_LOG_FILE="$LOG_DIR/frontend.log"

FRONTEND_PORT="${FRONTEND_PORT:-5173}"
START_FRONTEND="${START_FRONTEND:-1}"
AUTO_NPM_INSTALL="${AUTO_NPM_INSTALL:-1}"
STARTUP_RETRIES="${STARTUP_RETRIES:-300}"
FRONTEND_RETRIES="${FRONTEND_RETRIES:-300}"

mkdir -p "$LOG_DIR"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Error: missing command '$1'"
    exit 1
  fi
}

need_cmd bash
need_cmd go
need_cmd curl
need_cmd docker

if [[ "$START_FRONTEND" == "1" ]]; then
  need_cmd npm
fi

if [[ -f "$INIT_ENV_SCRIPT" ]]; then
  if [[ ! -x "$INIT_ENV_SCRIPT" ]]; then
    chmod +x "$INIT_ENV_SCRIPT"
  fi
  "$INIT_ENV_SCRIPT"
else
  echo "[WARN] init script not found: $INIT_ENV_SCRIPT"
fi

read_pid() {
  local pid_file="$1"
  if [[ -f "$pid_file" ]]; then
    tr -d '[:space:]' <"$pid_file"
  fi
}

is_pid_alive() {
  local pid="$1"
  [[ -n "$pid" ]] && kill -0 "$pid" >/dev/null 2>&1
}

wait_http_ready() {
  local url="$1"
  local retries="${2:-50}"
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
  local retries="${3:-50}"
  for _ in $(seq 1 "$retries"); do
    if timeout 1 bash -c "</dev/tcp/$host/$port" >/dev/null 2>&1; then
      return 0
    fi
    sleep 0.3
  done
  return 1
}

start_go_service_if_missing() {
  local name="$1"
  local check_type="$2"
  local endpoint1="$3"
  local endpoint2="$4"
  local pid_file="$5"
  local log_file="$6"
  local go_cmd="$7"

  local ready=1
  if [[ "$check_type" == "http" ]]; then
    if wait_http_ready "$endpoint1" 2; then
      ready=0
    fi
  else
    if wait_tcp_ready "$endpoint1" "$endpoint2" 2; then
      ready=0
    fi
  fi

  if [[ "$ready" -eq 0 ]]; then
    echo "[OK] $name already reachable"
    return 0
  fi

  local existing_pid
  existing_pid="$(read_pid "$pid_file")"

  if is_pid_alive "$existing_pid"; then
    echo "[INFO] $name process exists (pid=$existing_pid), waiting ready"
  else
    echo "[INFO] starting $name ..."
    (
      cd "$ROOT_DIR"
      nohup bash -lc "$go_cmd" >"$log_file" 2>&1 &
      echo "$!" >"$pid_file"
    )
  fi

  if [[ "$check_type" == "http" ]]; then
    if ! wait_http_ready "$endpoint1" "$STARTUP_RETRIES"; then
      echo "[ERR] $name failed to become ready"
      echo "[HINT] log: $log_file"
      return 1
    fi
  else
    if ! wait_tcp_ready "$endpoint1" "$endpoint2" "$STARTUP_RETRIES"; then
      echo "[ERR] $name failed to become ready"
      echo "[HINT] log: $log_file"
      return 1
    fi
  fi

  echo "[OK] $name is ready"
}

start_frontend_if_missing() {
  local host="127.0.0.1"
  local port="$FRONTEND_PORT"
  local endpoint="http://$host:$port"

  if wait_http_ready "$endpoint" 2; then
    echo "[OK] frontend already reachable at $endpoint"
    return 0
  fi

  local existing_pid
  existing_pid="$(read_pid "$FRONTEND_PID_FILE")"

  if is_pid_alive "$existing_pid"; then
    echo "[INFO] frontend process exists (pid=$existing_pid), waiting ready"
  else
    if [[ "$AUTO_NPM_INSTALL" == "1" && ! -d "$ROOT_DIR/frontend/vue-app/node_modules" ]]; then
      echo "[INFO] frontend deps not found, running npm install ..."
      (
        cd "$ROOT_DIR/frontend/vue-app"
        npm install --no-fund --no-audit
      )
    fi

    echo "[INFO] starting frontend ..."
    (
      cd "$ROOT_DIR/frontend/vue-app"
      nohup bash -lc "npm run dev -- --host 0.0.0.0 --port $port" >"$FRONTEND_LOG_FILE" 2>&1 &
      echo "$!" >"$FRONTEND_PID_FILE"
    )
  fi

  if ! wait_http_ready "$endpoint" "$FRONTEND_RETRIES"; then
    echo "[ERR] frontend failed to become ready"
    echo "[HINT] log: $FRONTEND_LOG_FILE"
    return 1
  fi

  echo "[OK] frontend is ready"
}

echo "========================================"
echo "Step 1/3: Start infra services via docker compose"
echo "========================================"
(
  cd "$ROOT_DIR"
  docker compose up -d
)

echo "========================================"
echo "Step 2/3: Start backend services"
echo "========================================"
start_go_service_if_missing \
  "authapi" "http" "http://127.0.0.1:8000/auth/login" "" \
  "$AUTH_PID_FILE" "$AUTH_LOG_FILE" \
  "go run ./rebuild/authapi/auth.go -f ./rebuild/authapi/etc/auth-api.yaml"

start_go_service_if_missing \
  "userrpc" "tcp" "127.0.0.1" "8004" \
  "$USER_RPC_PID_FILE" "$USER_RPC_LOG_FILE" \
  "go run ./rebuild/user/rpc/user.go -f ./rebuild/user/rpc/etc/user.yaml"

start_go_service_if_missing \
  "userapi" "http" "http://127.0.0.1:8003/api/user/me" "" \
  "$USER_API_PID_FILE" "$USER_API_LOG_FILE" \
  "go run ./rebuild/user/api/userapi.go -f ./rebuild/user/api/etc/userapi.yaml"

start_go_service_if_missing \
  "syncnote-rpc" "tcp" "127.0.0.1" "8002" \
  "$SYNC_RPC_PID_FILE" "$SYNC_RPC_LOG_FILE" \
  "go run ./rebuild/syncnote/rpc/syncnote.go -f ./rebuild/syncnote/rpc/etc/syncnote.yaml"

start_go_service_if_missing \
  "syncnote-api" "http" "http://127.0.0.1:8001/note/list" "" \
  "$SYNC_API_PID_FILE" "$SYNC_API_LOG_FILE" \
  "go run ./rebuild/syncnote/api/syncnote.go -f ./rebuild/syncnote/api/etc/syncnote.yaml"

if [[ "$START_FRONTEND" == "1" ]]; then
  echo "========================================"
  echo "Step 3/3: Start frontend"
  echo "========================================"
  start_frontend_if_missing
else
  echo "[INFO] START_FRONTEND=0, skip frontend"
fi

echo ""
echo "All services are up."
echo "- Frontend:     http://127.0.0.1:${FRONTEND_PORT}"
echo "- Auth API:     http://127.0.0.1:8000"
echo "- Syncnote API: http://127.0.0.1:8001"
echo "- User API:     http://127.0.0.1:8003"
echo ""
echo "Logs directory: $LOG_DIR"
