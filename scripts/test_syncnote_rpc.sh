#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if ! command -v go >/dev/null 2>&1; then
  echo "Error: go is not installed or not in PATH"
  exit 1
fi

RPC_ADDR="${RPC_ADDR:-127.0.0.1:8080}"

echo "Running RPC smoke test against $RPC_ADDR ..."
(
  cd "$ROOT_DIR"
  RPC_ADDR="$RPC_ADDR" go run ./scripts/rpc_smoke.go
)
