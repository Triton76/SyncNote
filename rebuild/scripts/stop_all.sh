#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
RUN_DIR="$ROOT_DIR/rebuild/.run"

AUTH_PID_FILE="$RUN_DIR/authapi.pid"
USER_RPC_PID_FILE="$RUN_DIR/userrpc.pid"
USER_API_PID_FILE="$RUN_DIR/userapi.pid"
SYNC_RPC_PID_FILE="$RUN_DIR/syncnote-rpc.pid"
SYNC_API_PID_FILE="$RUN_DIR/syncnote-api.pid"
FRONTEND_PID_FILE="$RUN_DIR/frontend.pid"

# By default stop docker-compose infra too; set STOP_INFRA=0 to keep mysql/redis/etcd.
STOP_INFRA="${STOP_INFRA:-1}"
GRACE_SECONDS="${GRACE_SECONDS:-8}"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Error: missing command '$1'"
    exit 1
  fi
}

need_cmd bash
need_cmd docker

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

stop_pid_file_service() {
  local name="$1"
  local pid_file="$2"

  local pid
  pid="$(read_pid "$pid_file")"

  if [[ -z "$pid" ]]; then
    echo "[INFO] $name: no pid file"
    rm -f "$pid_file"
    return 0
  fi

  if ! is_pid_alive "$pid"; then
    echo "[INFO] $name: stale pid ($pid), cleaning"
    rm -f "$pid_file"
    return 0
  fi

  echo "[INFO] stopping $name (pid=$pid)"
  kill "$pid" >/dev/null 2>&1 || true

  for _ in $(seq 1 "$GRACE_SECONDS"); do
    if ! is_pid_alive "$pid"; then
      rm -f "$pid_file"
      echo "[OK] $name stopped"
      return 0
    fi
    sleep 1
  done

  echo "[WARN] $name did not exit in time, force killing"
  kill -9 "$pid" >/dev/null 2>&1 || true
  rm -f "$pid_file"
  echo "[OK] $name stopped (SIGKILL)"
}

# Fallback: if pid file is missing but process still exists, terminate by signature.
kill_by_pattern() {
  local name="$1"
  local pattern="$2"
  local pids

  pids="$(pgrep -f "$pattern" || true)"
  if [[ -z "$pids" ]]; then
    return 0
  fi

  echo "[INFO] stopping $name by pattern"
  # shellcheck disable=SC2086
  kill $pids >/dev/null 2>&1 || true
  sleep 1
  pids="$(pgrep -f "$pattern" || true)"
  if [[ -n "$pids" ]]; then
    # shellcheck disable=SC2086
    kill -9 $pids >/dev/null 2>&1 || true
  fi
}

echo "========================================"
echo "Step 1/2: Stop app services"
echo "========================================"

stop_pid_file_service "frontend" "$FRONTEND_PID_FILE"
stop_pid_file_service "syncnote-api" "$SYNC_API_PID_FILE"
stop_pid_file_service "syncnote-rpc" "$SYNC_RPC_PID_FILE"
stop_pid_file_service "userapi" "$USER_API_PID_FILE"
stop_pid_file_service "userrpc" "$USER_RPC_PID_FILE"
stop_pid_file_service "authapi" "$AUTH_PID_FILE"

kill_by_pattern "frontend" "vite.*--port"
kill_by_pattern "syncnote-api" "go run ./rebuild/syncnote/api/syncnote.go|/exe/syncnote -f ./rebuild/syncnote/api/etc/syncnote.yaml"
kill_by_pattern "syncnote-rpc" "go run ./rebuild/syncnote/rpc/syncnote.go|/exe/syncnote -f ./rebuild/syncnote/rpc/etc/syncnote.yaml"
kill_by_pattern "userapi" "go run ./rebuild/user/api/userapi.go|/exe/userapi -f ./rebuild/user/api/etc/userapi.yaml"
kill_by_pattern "userrpc" "go run ./rebuild/user/rpc/user.go|/exe/user -f ./rebuild/user/rpc/etc/user.yaml"
kill_by_pattern "authapi" "go run ./rebuild/authapi/auth.go|/exe/auth -f ./rebuild/authapi/etc/auth-api.yaml"

if [[ "$STOP_INFRA" == "1" ]]; then
  echo "========================================"
  echo "Step 2/2: Stop infra (docker compose)"
  echo "========================================"
  (
    cd "$ROOT_DIR"
    docker compose down
  )
else
  echo "[INFO] STOP_INFRA=0, keep docker infra running"
fi

echo "All requested services are stopped."
