#!/usr/bin/env bash
set -euo pipefail

# One-time environment bootstrap for local development.
# Default Go module proxy is goproxy.cn for faster dependency download in CN networks.
GO_PROXY_VALUE="${GO_PROXY_VALUE:-https://goproxy.cn,direct}"
SKIP_GO_PROXY_SETUP="${SKIP_GO_PROXY_SETUP:-0}"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Error: missing command '$1'"
    exit 1
  fi
}

need_cmd go

if [[ "$SKIP_GO_PROXY_SETUP" == "1" ]]; then
  echo "[INFO] SKIP_GO_PROXY_SETUP=1, skip Go proxy setup"
  exit 0
fi

current_proxy="$(go env GOPROXY)"
if [[ "$current_proxy" == "$GO_PROXY_VALUE" ]]; then
  echo "[OK] GOPROXY already set: $current_proxy"
  exit 0
fi

echo "[INFO] setting GOPROXY => $GO_PROXY_VALUE"
go env -w "GOPROXY=$GO_PROXY_VALUE"
echo "[OK] GOPROXY now: $(go env GOPROXY)"
