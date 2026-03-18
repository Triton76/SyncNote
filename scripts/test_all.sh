#!/usr/bin/env bash
set -u

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
API_SCRIPT="$ROOT_DIR/scripts/test_syncnote_api.sh"
RPC_SCRIPT="$ROOT_DIR/scripts/test_syncnote_rpc.sh"

MODE="${1:-all}"

usage() {
  cat <<'EOF'
Usage:
  ./scripts/test_all.sh [all|api|rpc]

Env:
  API_HOST   API base URL for API smoke test (default: http://127.0.0.1:8888)
  USER_ID    Optional user id for API smoke test
  RPC_ADDR   RPC target for RPC smoke test (default: 127.0.0.1:8080)

Examples:
  ./scripts/test_all.sh
  ./scripts/test_all.sh api
  API_HOST=http://127.0.0.1:8888 RPC_ADDR=127.0.0.1:8080 ./scripts/test_all.sh all
EOF
}

if [[ ! -x "$API_SCRIPT" ]]; then
  echo "Error: missing or non-executable $API_SCRIPT"
  echo "Try: chmod +x scripts/test_syncnote_api.sh"
  exit 1
fi

if [[ ! -x "$RPC_SCRIPT" ]]; then
  echo "Error: missing or non-executable $RPC_SCRIPT"
  echo "Try: chmod +x scripts/test_syncnote_rpc.sh"
  exit 1
fi

run_api() {
  echo "========================================"
  echo "Running API smoke test"
  echo "API_HOST=${API_HOST:-http://127.0.0.1:8888}"
  echo "========================================"
  "$API_SCRIPT"
}

run_rpc() {
  echo "========================================"
  echo "Running RPC smoke test"
  echo "RPC_ADDR=${RPC_ADDR:-127.0.0.1:8080}"
  echo "========================================"
  "$RPC_SCRIPT"
}

api_rc=0
rpc_rc=0

case "$MODE" in
  all)
    run_api || api_rc=$?
    run_rpc || rpc_rc=$?
    ;;
  api)
    run_api || api_rc=$?
    ;;
  rpc)
    run_rpc || rpc_rc=$?
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
printf "API result: %s\n" "$( [[ $api_rc -eq 0 ]] && echo PASS || echo FAIL )"
printf "RPC result: %s\n" "$( [[ $rpc_rc -eq 0 ]] && echo PASS || echo FAIL )"
echo "========================================"

if [[ $api_rc -ne 0 || $rpc_rc -ne 0 ]]; then
  exit 1
fi

echo "All requested tests passed."
