#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RUN_DIR="$ROOT_DIR/.run"
LOG_DIR="$RUN_DIR/logs"

AUTH_PID="$RUN_DIR/auth-api.pid"
SYNCNOTE_API_PID="$RUN_DIR/syncnote-api.pid"
SYNCNOTE_RPC_PID="$RUN_DIR/syncnote-rpc.pid"

AUTH_LOG="$LOG_DIR/auth-api.log"
SYNCNOTE_API_LOG="$LOG_DIR/syncnote-api.log"
SYNCNOTE_RPC_LOG="$LOG_DIR/syncnote-rpc.log"

DOCKER_COMPOSE_FILE="$ROOT_DIR/docker-compose.yml"
FORCE_KILL_CONFLICTS=0

mkdir -p "$LOG_DIR"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "[ERR] Missing command: $1"
    exit 1
  fi
}

compose() {
  if docker compose version >/dev/null 2>&1; then
    docker compose -f "$DOCKER_COMPOSE_FILE" "$@"
  elif command -v docker-compose >/dev/null 2>&1; then
    docker-compose -f "$DOCKER_COMPOSE_FILE" "$@"
  else
    echo "[ERR] docker compose is not available"
    exit 1
  fi
}

is_pid_alive() {
  local pid="$1"
  [[ -n "$pid" ]] && kill -0 "$pid" >/dev/null 2>&1
}

list_listen_pids_by_port() {
  local port="$1"
  lsof -t -iTCP:"$port" -sTCP:LISTEN 2>/dev/null || true
}

ensure_port_available() {
  local service_name="$1"
  local port="$2"

  local conflict_pids
  conflict_pids="$(list_listen_pids_by_port "$port")"
  if [[ -z "$conflict_pids" ]]; then
    return 0
  fi

  if [[ "$FORCE_KILL_CONFLICTS" == "1" ]]; then
    echo "[WARN] $service_name port $port is occupied, force cleaning stale listeners"
    while IFS= read -r pid; do
      [[ -z "$pid" ]] && continue
      if is_pid_alive "$pid"; then
        kill "$pid" >/dev/null 2>&1 || true
      fi
    done <<<"$conflict_pids"

    sleep 1
    local remains
    remains="$(list_listen_pids_by_port "$port")"
    if [[ -n "$remains" ]]; then
      while IFS= read -r pid; do
        [[ -z "$pid" ]] && continue
        if is_pid_alive "$pid"; then
          kill -9 "$pid" >/dev/null 2>&1 || true
        fi
      done <<<"$remains"
    fi

    sleep 0.3
    if [[ -n "$(list_listen_pids_by_port "$port")" ]]; then
      echo "[ERR] unable to free port $port for $service_name"
      return 1
    fi
    return 0
  fi

  echo "[ERR] $service_name cannot start: port $port is already in use"
  echo "[HINT] run: scripts/dev_services.sh up --force-kill-conflicts"
  echo "[HINT] or inspect manually: lsof -iTCP:$port -sTCP:LISTEN"
  return 1
}

read_pid() {
  local pid_file="$1"
  if [[ -f "$pid_file" ]]; then
    tr -d '[:space:]' <"$pid_file"
  fi
}

start_go_service() {
  local name="$1"
  local workdir="$2"
  local config_rel="$3"
  local pid_file="$4"
  local log_file="$5"

  local existing_pid
  existing_pid="$(read_pid "$pid_file")"
  if is_pid_alive "$existing_pid"; then
    echo "[OK] $name already running (pid=$existing_pid)"
    return 0
  fi

  rm -f "$pid_file"
  echo "[INFO] starting $name ..."
  (
    cd "$workdir"
    nohup go run . -f "$config_rel" >"$log_file" 2>&1 &
    echo "$!" >"$pid_file"
  )

  sleep 1
  local new_pid
  new_pid="$(read_pid "$pid_file")"
  if is_pid_alive "$new_pid"; then
    echo "[OK] $name started (pid=$new_pid)"
  else
    echo "[ERR] $name failed to start. Check log: $log_file"
    exit 1
  fi
}

stop_go_service() {
  local name="$1"
  local pid_file="$2"

  local pid
  pid="$(read_pid "$pid_file")"
  if ! is_pid_alive "$pid"; then
    rm -f "$pid_file"
    echo "[OK] $name is not running"
    return 0
  fi

  echo "[INFO] stopping $name (pid=$pid) ..."
  kill "$pid" >/dev/null 2>&1 || true

  for _ in {1..20}; do
    if ! is_pid_alive "$pid"; then
      break
    fi
    sleep 0.2
  done

  if is_pid_alive "$pid"; then
    echo "[WARN] force killing $name (pid=$pid)"
    kill -9 "$pid" >/dev/null 2>&1 || true
  fi

  rm -f "$pid_file"
  echo "[OK] $name stopped"
}

show_go_service() {
  local name="$1"
  local pid_file="$2"

  local pid
  pid="$(read_pid "$pid_file")"
  if is_pid_alive "$pid"; then
    echo "[RUNNING] $name (pid=$pid)"
  else
    echo "[STOPPED] $name"
  fi
}

wait_for_container_health() {
  local container="$1"
  local retries="${2:-30}"

  for ((i = 1; i <= retries; i++)); do
    local status
    status="$(docker inspect --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}' "$container" 2>/dev/null || true)"

    if [[ "$status" == "healthy" || "$status" == "running" ]]; then
      echo "[OK] $container is $status"
      return 0
    fi

    sleep 1
  done

  echo "[ERR] $container did not become healthy/running in time"
  return 1
}

up() {
  need_cmd docker
  need_cmd go
  need_cmd lsof

  echo "[STEP] starting infrastructure containers"
  compose up -d etcd redis mysql

  wait_for_container_health etcd 40
  wait_for_container_health redis 40
  wait_for_container_health mysql 60

  ensure_port_available "syncnote-rpc" 8080
  ensure_port_available "syncnote-api" 8888
  ensure_port_available "auth-api" 8889

  echo "[STEP] starting Go services"
  start_go_service "syncnote-rpc" "$ROOT_DIR/syncnote/rpc" "etc/syncnoterpc.yaml" "$SYNCNOTE_RPC_PID" "$SYNCNOTE_RPC_LOG"
  start_go_service "auth-api" "$ROOT_DIR/auth/api" "etc/auth-api.yaml" "$AUTH_PID" "$AUTH_LOG"
  start_go_service "syncnote-api" "$ROOT_DIR/syncnote/api" "etc/syncnote-api.yaml" "$SYNCNOTE_API_PID" "$SYNCNOTE_API_LOG"

  echo "[DONE] all services are up"
  status
}

down() {
  stop_go_service "syncnote-api" "$SYNCNOTE_API_PID"
  stop_go_service "auth-api" "$AUTH_PID"
  stop_go_service "syncnote-rpc" "$SYNCNOTE_RPC_PID"

  need_cmd docker
  echo "[STEP] stopping infrastructure containers"
  compose stop etcd redis mysql >/dev/null
  echo "[DONE] all services are down"
}

status() {
  echo "--- Go services ---"
  show_go_service "auth-api" "$AUTH_PID"
  show_go_service "syncnote-api" "$SYNCNOTE_API_PID"
  show_go_service "syncnote-rpc" "$SYNCNOTE_RPC_PID"

  echo
  echo "--- Containers ---"
  compose ps --status running || true
}

logs() {
  local target="${1:-all}"
  case "$target" in
    auth)
      tail -f "$AUTH_LOG"
      ;;
    api)
      tail -f "$SYNCNOTE_API_LOG"
      ;;
    rpc)
      tail -f "$SYNCNOTE_RPC_LOG"
      ;;
    infra)
      compose logs -f etcd redis mysql
      ;;
    all)
      echo "[INFO] logs directory: $LOG_DIR"
      echo "[INFO] use one of: auth | api | rpc | infra"
      ;;
    *)
      echo "Usage: $0 logs [auth|api|rpc|infra|all]"
      exit 1
      ;;
  esac
}

usage() {
  cat <<'EOF'
Usage:
  scripts/dev_services.sh up
  scripts/dev_services.sh up --force-kill-conflicts
  scripts/dev_services.sh down
  scripts/dev_services.sh restart
  scripts/dev_services.sh status
  scripts/dev_services.sh logs [auth|api|rpc|infra|all]

Notes:
  - Go service logs are stored in .run/logs/
  - Infra dependencies are managed by docker compose
EOF
}

cmd="${1:-status}"
arg2="${2:-}"

if [[ "$arg2" == "--force-kill-conflicts" ]]; then
  FORCE_KILL_CONFLICTS=1
fi

case "$cmd" in
  up)
    up
    ;;
  down)
    down
    ;;
  restart)
    down
    up
    ;;
  status)
    status
    ;;
  logs)
    logs "${2:-all}"
    ;;
  *)
    usage
    exit 1
    ;;
esac
