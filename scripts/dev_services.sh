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

first_listen_pid_by_port() {
  local port="$1"
  list_listen_pids_by_port "$port" | head -n 1
}

pid_cwd() {
  local pid="$1"
  readlink "/proc/$pid/cwd" 2>/dev/null || true
}

resolve_port_conflict() {
  local service_name="$1"
  local workdir="$2"
  local port="$3"
  local pid_file="$4"

  local conflict_pids
  conflict_pids="$(list_listen_pids_by_port "$port")"
  if [[ -z "$conflict_pids" ]]; then
    return 0
  fi

  local first_pid
  first_pid="$(first_listen_pid_by_port "$port")"
  if [[ -n "$first_pid" ]]; then
    local first_cwd
    first_cwd="$(pid_cwd "$first_pid")"
    if [[ "$first_cwd" == "$workdir" ]]; then
      echo "[OK] $service_name already running (pid=$first_pid, adopted by port=$port)"
      echo "$first_pid" >"$pid_file"
      return 2
    fi
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
  local port="$4"
  local pid_file="$5"
  local log_file="$6"

  local existing_pid
  existing_pid="$(read_pid "$pid_file")"
  if is_pid_alive "$existing_pid"; then
    echo "[OK] $name already running (pid=$existing_pid)"
    return 0
  fi

  local conflict_state=0
  resolve_port_conflict "$name" "$workdir" "$port" "$pid_file" || conflict_state=$?
  if [[ "$conflict_state" == "2" ]]; then
    return 0
  elif [[ "$conflict_state" != "0" ]]; then
    return 1
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
  local workdir="$3"
  local port="$4"

  local pid
  pid="$(read_pid "$pid_file")"
  if ! is_pid_alive "$pid"; then
    local adopted_pid
    adopted_pid="$(first_listen_pid_by_port "$port")"
    if [[ -n "$adopted_pid" && "$(pid_cwd "$adopted_pid")" == "$workdir" ]]; then
      pid="$adopted_pid"
      echo "$pid" >"$pid_file"
    else
      rm -f "$pid_file"
      echo "[OK] $name is not running"
      return 0
    fi
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
  local workdir="$3"
  local port="$4"

  local pid
  pid="$(read_pid "$pid_file")"
  if is_pid_alive "$pid"; then
    echo "[RUNNING] $name (pid=$pid)"
  else
    local adopted_pid
    adopted_pid="$(first_listen_pid_by_port "$port")"
    if [[ -n "$adopted_pid" && "$(pid_cwd "$adopted_pid")" == "$workdir" ]]; then
      echo "$adopted_pid" >"$pid_file"
      echo "[RUNNING] $name (pid=$adopted_pid, adopted by port=$port)"
    else
      echo "[STOPPED] $name"
    fi
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

apply_sql_file() {
  local sql_file="$1"
  if [[ ! -f "$sql_file" ]]; then
    echo "[WARN] schema file not found: $sql_file"
    return 0
  fi

  echo "[INFO] applying schema: ${sql_file#$ROOT_DIR/}"
  mysql -h127.0.0.1 -P3306 -uroot -pdevpass123 syncnote <"$sql_file"
}

bootstrap_mysql_schema() {
  echo "[STEP] bootstrapping MySQL schema"
  apply_sql_file "$ROOT_DIR/auth/auth.sql"
  apply_sql_file "$ROOT_DIR/syncnote/rpc/internal/model/notesmodel.sql"
  apply_sql_file "$ROOT_DIR/syncnote/rpc/internal/model/collaboration.sql"
  echo "[OK] MySQL schema is ready"
}

up() {
  need_cmd docker
  need_cmd go
  need_cmd lsof
  need_cmd mysql

  echo "[STEP] starting infrastructure containers"
  compose up -d etcd redis mysql

  wait_for_container_health etcd 40
  wait_for_container_health redis 40
  wait_for_container_health mysql 60

  bootstrap_mysql_schema

  echo "[STEP] starting Go services"
  start_go_service "syncnote-rpc" "$ROOT_DIR/syncnote/rpc" "etc/syncnoterpc.yaml" 8080 "$SYNCNOTE_RPC_PID" "$SYNCNOTE_RPC_LOG"
  start_go_service "auth-api" "$ROOT_DIR/auth/api" "etc/auth-api.yaml" 8889 "$AUTH_PID" "$AUTH_LOG"
  start_go_service "syncnote-api" "$ROOT_DIR/syncnote/api" "etc/syncnote-api.yaml" 8888 "$SYNCNOTE_API_PID" "$SYNCNOTE_API_LOG"

  echo "[DONE] all services are up"
  status
}

down() {
  stop_go_service "syncnote-api" "$SYNCNOTE_API_PID" "$ROOT_DIR/syncnote/api" 8888
  stop_go_service "auth-api" "$AUTH_PID" "$ROOT_DIR/auth/api" 8889
  stop_go_service "syncnote-rpc" "$SYNCNOTE_RPC_PID" "$ROOT_DIR/syncnote/rpc" 8080

  need_cmd docker
  echo "[STEP] stopping infrastructure containers"
  compose stop etcd redis mysql >/dev/null
  echo "[DONE] all services are down"
}

status() {
  echo "--- Go services ---"
  show_go_service "auth-api" "$AUTH_PID" "$ROOT_DIR/auth/api" 8889
  show_go_service "syncnote-api" "$SYNCNOTE_API_PID" "$ROOT_DIR/syncnote/api" 8888
  show_go_service "syncnote-rpc" "$SYNCNOTE_RPC_PID" "$ROOT_DIR/syncnote/rpc" 8080

  echo
  echo "--- Containers ---"
  compose ps --status running || true
}

show_go_service_ps() {
  local name="$1"
  local pid_file="$2"
  local workdir="$3"
  local port="$4"

  local pid
  pid="$(read_pid "$pid_file")"
  if ! is_pid_alive "$pid"; then
    local adopted_pid
    adopted_pid="$(first_listen_pid_by_port "$port")"
    if [[ -n "$adopted_pid" && "$(pid_cwd "$adopted_pid")" == "$workdir" ]]; then
      pid="$adopted_pid"
      echo "$pid" >"$pid_file"
    else
      printf '%-14s %-8s %-6s %-35s %s\n' "$name" "-" "$port" "$workdir" "stopped"
      return 0
    fi
  fi

  local cwd
  cwd="$(pid_cwd "$pid")"
  local cmd
  cmd="$(ps -p "$pid" -o args= 2>/dev/null || true)"
  if [[ -z "$cmd" ]]; then
    cmd="unknown"
  fi
  printf '%-14s %-8s %-6s %-35s %s\n' "$name" "$pid" "$port" "$cwd" "$cmd"
}

ps_services() {
  echo "SERVICE        PID      PORT   CWD                                 CMD"
  show_go_service_ps "auth-api" "$AUTH_PID" "$ROOT_DIR/auth/api" 8889
  show_go_service_ps "syncnote-api" "$SYNCNOTE_API_PID" "$ROOT_DIR/syncnote/api" 8888
  show_go_service_ps "syncnote-rpc" "$SYNCNOTE_RPC_PID" "$ROOT_DIR/syncnote/rpc" 8080
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
  scripts/dev_services.sh ps
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
  ps)
    ps_services
    ;;
  logs)
    logs "${2:-all}"
    ;;
  *)
    usage
    exit 1
    ;;
esac
