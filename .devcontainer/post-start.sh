#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="/workspaces/SyncNote"
RUN_DIR="$ROOT_DIR/rebuild/.run"
LOG_DIR="$RUN_DIR/logs"
mkdir -p "$LOG_DIR"

wait_port() {
  local name="$1"
  local host="$2"
  local port="$3"
  for _ in $(seq 1 60); do
    if timeout 1 bash -lc "</dev/tcp/$host/$port" >/dev/null 2>&1; then
      echo "[OK] $name ready at $host:$port"
      return 0
    fi
    sleep 0.5
  done
  echo "[ERR] $name not ready at $host:$port"
  return 1
}

ensure_container_running() {
  local name="$1"
  local compose_service="$2"

  if docker ps --format '{{.Names}}' | grep -qx "$name"; then
    echo "[OK] container $name already running"
    return 0
  fi

  if docker ps -a --format '{{.Names}}' | grep -qx "$name"; then
    echo "[INFO] starting existing container $name"
    docker start "$name" >/dev/null || true
    return 0
  fi

  echo "[INFO] creating container via compose service $compose_service"
  (cd "$ROOT_DIR" && docker compose up -d "$compose_service" >/dev/null) || true
}

start_if_missing() {
  local name="$1"
  local port="$2"
  local cmd="$3"
  local pid_file="$RUN_DIR/${name}.pid"
  local log_file="$LOG_DIR/${name}.log"

  if ss -ltn "sport = :$port" | grep -q ":$port"; then
    echo "[OK] $name already listening on :$port"
    return 0
  fi

  if [[ -f "$pid_file" ]]; then
    local old_pid
    old_pid="$(tr -d '[:space:]' <"$pid_file" || true)"
    if [[ -n "$old_pid" ]] && kill -0 "$old_pid" >/dev/null 2>&1; then
      echo "[INFO] stopping stale $name process pid=$old_pid"
      kill "$old_pid" >/dev/null 2>&1 || true
      sleep 1
    fi
  fi

  echo "[INFO] starting $name ..."
  (cd "$ROOT_DIR" && nohup bash -lc "$cmd" >"$log_file" 2>&1 & echo $! >"$pid_file")

  for _ in $(seq 1 40); do
    if ss -ltn "sport = :$port" | grep -q ":$port"; then
      echo "[OK] $name is ready on :$port"
      return 0
    fi
    sleep 0.3
  done

  echo "[ERR] $name failed to start, see $log_file"
  tail -n 30 "$log_file" 2>/dev/null || true
  return 1
}

echo "[devcontainer] Ensuring docker dependencies are up..."
if command -v docker >/dev/null 2>&1 && [[ -f "$ROOT_DIR/docker-compose.yml" ]]; then
  ensure_container_running "etcd" "etcd"
  ensure_container_running "redis" "redis"
  ensure_container_running "mysql" "mysql"
else
  echo "[WARN] docker or docker-compose.yml not available, skip dependencies"
fi

wait_port "etcd" "127.0.0.1" "2379"
wait_port "redis" "127.0.0.1" "6379"
wait_port "mysql" "127.0.0.1" "3306"

echo "[devcontainer] Ensuring application services are up..."
start_if_missing "authapi" 8000 "go run ./rebuild/authapi/auth.go -f ./rebuild/authapi/etc/auth-api.yaml"
start_if_missing "syncnote-rpc" 8002 "go run ./rebuild/syncnote/rpc/syncnote.go -f ./rebuild/syncnote/rpc/etc/syncnote.yaml"
start_if_missing "syncnote-api" 8001 "go run ./rebuild/syncnote/api/syncnote.go -f ./rebuild/syncnote/api/etc/syncnote.yaml"
start_if_missing "user-rpc" 8004 "go run ./rebuild/user/rpc/user.go -f ./rebuild/user/rpc/etc/user.yaml"
start_if_missing "user-api" 8003 "go run ./rebuild/user/api/userapi.go -f ./rebuild/user/api/etc/userapi.yaml"

if [[ -f "$ROOT_DIR/frontend/vue-app/package.json" ]] && command -v npm >/dev/null 2>&1; then
  start_if_missing "vue-web" 5173 "npm --prefix ./frontend/vue-app run dev -- --host 0.0.0.0 --port 5173"
fi

echo "[devcontainer] Services are ready."
