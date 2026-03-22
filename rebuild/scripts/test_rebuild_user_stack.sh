#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
AUTH_SCRIPT="$ROOT_DIR/rebuild/scripts/test_authapi.sh"
USER_API_SCRIPT="$ROOT_DIR/rebuild/scripts/test_userapi.sh"
USER_RPC_SCRIPT="$ROOT_DIR/rebuild/scripts/test_userrpc.sh"
RUN_DIR="$ROOT_DIR/rebuild/.run"
LOG_DIR="$RUN_DIR/logs"

AUTH_PID_FILE="$RUN_DIR/authapi.pid"
USER_API_PID_FILE="$RUN_DIR/userapi.pid"
USER_RPC_PID_FILE="$RUN_DIR/userrpc.pid"

AUTH_LOG_FILE="$LOG_DIR/authapi.log"
USER_API_LOG_FILE="$LOG_DIR/userapi.log"
USER_RPC_LOG_FILE="$LOG_DIR/userrpc.log"
USER_RPC_CONFIG_FILE="${USER_RPC_CONFIG_FILE:-$ROOT_DIR/rebuild/user/rpc/etc/user.yaml}"
USER_RPC_DATASOURCE="${USER_RPC_DATASOURCE:-root:devpass123@tcp(127.0.0.1:3306)/syncnote?charset=utf8mb4&parseTime=true&loc=Local}"
USER_RPC_RUNTIME_CONFIG="$RUN_DIR/user-rpc.runtime.yaml"

MODE="${1:-all}"
AUTO_START_MISSING_SERVICES="${AUTO_START_MISSING_SERVICES:-1}"

mkdir -p "$LOG_DIR"

usage() {
  cat <<'EOF'
Usage:
  ./rebuild/scripts/test_rebuild_user_stack.sh [all|authapi|userapi|userrpc]

Env:
  AUTH_HOST      authapi address (default: http://127.0.0.1:8000)
  USER_API_HOST  userapi address (default: http://127.0.0.1:8888)
  RPC_ADDR       userrpc address (default: 127.0.0.1:8080)
  AUTO_START_MISSING_SERVICES auto start rebuild services if unavailable (default: 1)

Examples:
  ./rebuild/scripts/test_rebuild_user_stack.sh all
  AUTH_HOST=http://127.0.0.1:8000 ./rebuild/scripts/test_rebuild_user_stack.sh authapi
EOF
}

ensure_executable() {
  local script="$1"
  if [[ ! -f "$script" ]]; then
    echo "Error: missing script $script"
    exit 1
  fi
  if [[ ! -x "$script" ]]; then
    echo "Error: non-executable script $script"
    echo "Run: chmod +x rebuild/scripts/*.sh"
    exit 1
  fi
}

ensure_executable "$AUTH_SCRIPT"
ensure_executable "$USER_API_SCRIPT"
ensure_executable "$USER_RPC_SCRIPT"

is_pid_alive() {
  local pid="$1"
  [[ -n "$pid" ]] && kill -0 "$pid" >/dev/null 2>&1
}

read_pid() {
  local pid_file="$1"
  if [[ -f "$pid_file" ]]; then
    tr -d '[:space:]' <"$pid_file"
  fi
}

wait_http_ready() {
  local url="$1"
  local retries="${2:-40}"
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
  local retries="${3:-40}"
  for _ in $(seq 1 "$retries"); do
    if timeout 1 bash -c "</dev/tcp/$host/$port" >/dev/null 2>&1; then
      return 0
    fi
    sleep 0.3
  done
  return 1
}

prepare_userrpc_config() {
  if [[ ! -f "$USER_RPC_CONFIG_FILE" ]]; then
    echo "[ERR] userrpc config not found: $USER_RPC_CONFIG_FILE"
    return 1
  fi

  if grep -q '^DataSource:' "$USER_RPC_CONFIG_FILE"; then
    printf "%s" "$USER_RPC_CONFIG_FILE"
    return 0
  fi

  cp "$USER_RPC_CONFIG_FILE" "$USER_RPC_RUNTIME_CONFIG"
  {
    echo ""
    echo "DataSource: $USER_RPC_DATASOURCE"
  } >>"$USER_RPC_RUNTIME_CONFIG"

  printf "%s" "$USER_RPC_RUNTIME_CONFIG"
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

  if [[ "$AUTO_START_MISSING_SERVICES" != "1" ]]; then
    echo "[WARN] $name is not reachable and auto-start is disabled"
    return 1
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
    if ! wait_http_ready "$endpoint1" 40; then
      echo "[ERR] $name failed to become ready"
      echo "[HINT] see log: $log_file"
      return 1
    fi
  else
    if ! wait_tcp_ready "$endpoint1" "$endpoint2" 40; then
      echo "[ERR] $name failed to become ready"
      echo "[HINT] see log: $log_file"
      return 1
    fi
  fi

  echo "[OK] $name is ready"
}

ensure_services_ready() {
  local auth_host="${AUTH_HOST:-http://127.0.0.1:8000}"
  local user_api_host="${USER_API_HOST:-http://127.0.0.1:8888}"
  local rpc_addr="${RPC_ADDR:-127.0.0.1:8080}"
  local rpc_host="${rpc_addr%%:*}"
  local rpc_port="${rpc_addr##*:}"
  local user_rpc_start_cfg

  user_rpc_start_cfg="$(prepare_userrpc_config)" || return 1

  start_go_service_if_missing \
    "authapi" "http" "$auth_host/auth/login" "" \
    "$AUTH_PID_FILE" "$AUTH_LOG_FILE" \
    "go run ./rebuild/authapi/auth.go -f ./rebuild/authapi/etc/auth-api.yaml"

  start_go_service_if_missing \
    "userrpc" "tcp" "$rpc_host" "$rpc_port" \
    "$USER_RPC_PID_FILE" "$USER_RPC_LOG_FILE" \
    "go run ./rebuild/user/rpc/user.go -f $user_rpc_start_cfg"

  start_go_service_if_missing \
    "userapi" "http" "$user_api_host/api/user/me" "" \
    "$USER_API_PID_FILE" "$USER_API_LOG_FILE" \
    "go run ./rebuild/user/api/userapi.go -f ./rebuild/user/api/etc/userapi.yaml"
}

ensure_services_ready

auth_rc=0
userapi_rc=0
userrpc_rc=0

run_auth() {
  echo "========================================"
  echo "Running authapi smoke test"
  echo "AUTH_HOST=${AUTH_HOST:-http://127.0.0.1:8000}"
  echo "========================================"
  "$AUTH_SCRIPT"
}

run_userapi() {
  echo "========================================"
  echo "Running userapi smoke test"
  echo "AUTH_HOST=${AUTH_HOST:-http://127.0.0.1:8000}"
  echo "USER_API_HOST=${USER_API_HOST:-http://127.0.0.1:8888}"
  echo "========================================"
  "$USER_API_SCRIPT"
}

run_userrpc() {
  echo "========================================"
  echo "Running userrpc smoke test"
  echo "AUTH_HOST=${AUTH_HOST:-http://127.0.0.1:8000}"
  echo "RPC_ADDR=${RPC_ADDR:-127.0.0.1:8080}"
  echo "========================================"
  "$USER_RPC_SCRIPT"
}

case "$MODE" in
  all)
    run_auth || auth_rc=$?
    run_userapi || userapi_rc=$?
    run_userrpc || userrpc_rc=$?
    ;;
  authapi)
    run_auth || auth_rc=$?
    ;;
  userapi)
    run_userapi || userapi_rc=$?
    ;;
  userrpc)
    run_userrpc || userrpc_rc=$?
    ;;
  -h|--help|help)
    usage
    exit 0
    ;;
  *)
    echo "Error: invalid mode '$MODE'"
    usage
    exit 1
    ;;
esac

echo "========================================"
echo "Summary"
printf "authapi: %s\n" "$( [[ $auth_rc -eq 0 ]] && echo PASS || echo FAIL )"
printf "userapi: %s\n" "$( [[ $userapi_rc -eq 0 ]] && echo PASS || echo FAIL )"
printf "userrpc: %s\n" "$( [[ $userrpc_rc -eq 0 ]] && echo PASS || echo FAIL )"
echo "========================================"

if [[ $auth_rc -ne 0 || $userapi_rc -ne 0 || $userrpc_rc -ne 0 ]]; then
  exit 1
fi

echo "All requested rebuild tests passed."
