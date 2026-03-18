#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SYNCNOTE_DIR="$ROOT_DIR/syncnote"
API_DIR="$SYNCNOTE_DIR/api"
RPC_DIR="$SYNCNOTE_DIR/rpc"

MODE="${1:-all}"
RUN_TESTS="${2:-test}"

usage() {
  cat <<'EOF'
Usage:
  ./scripts/gen.sh [all|api|rpc] [test|no-test]

Examples:
  ./scripts/gen.sh
  ./scripts/gen.sh api
  ./scripts/gen.sh rpc no-test
EOF
}

if ! command -v goctl >/dev/null 2>&1; then
  echo "Error: goctl is not installed or not in PATH"
  exit 1
fi

if ! command -v go >/dev/null 2>&1; then
  echo "Error: go is not installed or not in PATH"
  exit 1
fi

case "$MODE" in
  all)
    echo "[1/3] Generating API from SyncNote.api..."
    (cd "$API_DIR" && goctl api go -api SyncNote.api -dir .)

    echo "[2/3] Generating RPC from syncNoterpc.proto..."
    (cd "$RPC_DIR" && goctl rpc protoc syncNoterpc.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.)
    ;;
  api)
    echo "[1/2] Generating API from SyncNote.api..."
    (cd "$API_DIR" && goctl api go -api SyncNote.api -dir .)
    ;;
  rpc)
    echo "[1/2] Generating RPC from syncNoterpc.proto..."
    (cd "$RPC_DIR" && goctl rpc protoc syncNoterpc.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.)
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

if [[ "$RUN_TESTS" == "test" ]]; then
  echo "[final] Running go test ./... in syncnote..."
  (cd "$SYNCNOTE_DIR" && go test ./...)
elif [[ "$RUN_TESTS" == "no-test" ]]; then
  echo "[final] Skipping tests"
else
  echo "Error: second arg must be 'test' or 'no-test'"
  usage
  exit 1
fi

echo "Done."
