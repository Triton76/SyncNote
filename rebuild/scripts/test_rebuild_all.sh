#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
MODE="${1:-all}"
AUTO_RUN_SMOKE="${AUTO_RUN_SMOKE:-1}"

USER_STACK_SCRIPT="$ROOT_DIR/rebuild/scripts/test_rebuild_user_stack.sh"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Error: missing command '$1'"
    exit 1
  fi
}

need_cmd go
need_cmd bash

usage() {
  cat <<'EOF'
Usage:
  ./rebuild/scripts/test_rebuild_all.sh [all|build|smoke]

Modes:
  all   : run build checks and smoke checks
  build : run go test build checks only
  smoke : run smoke checks only

Env:
  AUTO_RUN_SMOKE  run smoke suite in all mode (1/0, default 1)
  AUTH_HOST       forwarded to smoke scripts
  USER_API_HOST   forwarded to smoke scripts
  RPC_ADDR        forwarded to smoke scripts

Examples:
  ./rebuild/scripts/test_rebuild_all.sh
  ./rebuild/scripts/test_rebuild_all.sh build
  AUTO_RUN_SMOKE=0 ./rebuild/scripts/test_rebuild_all.sh all
EOF
}

run_build_suite() {
  echo "========================================"
  echo "Build checks"
  echo "========================================"

  (
    cd "$ROOT_DIR"
    go test ./rebuild/authapi/...
    go test ./rebuild/user/...
    go test ./rebuild/syncnote/api/...
    go test ./rebuild/syncnote/rpc/...
    go test ./rebuild/common/...
    go test ./rebuild/pkg/...
  )

  echo "Build checks passed."
}

run_smoke_suite() {
  if [[ ! -f "$USER_STACK_SCRIPT" ]]; then
    echo "Error: missing smoke script: $USER_STACK_SCRIPT"
    exit 1
  fi
  if [[ ! -x "$USER_STACK_SCRIPT" ]]; then
    echo "Error: non-executable smoke script: $USER_STACK_SCRIPT"
    echo "Run: chmod +x rebuild/scripts/*.sh"
    exit 1
  fi

  echo "========================================"
  echo "Smoke checks"
  echo "========================================"

  "$USER_STACK_SCRIPT" all

  echo "Smoke checks passed."
}

case "$MODE" in
  all)
    run_build_suite
    if [[ "$AUTO_RUN_SMOKE" == "1" ]]; then
      run_smoke_suite
    else
      echo "AUTO_RUN_SMOKE=0, skip smoke checks."
    fi
    ;;
  build)
    run_build_suite
    ;;
  smoke)
    run_smoke_suite
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

echo "All requested rebuild project tests passed."
