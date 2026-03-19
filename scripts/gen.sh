#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SYNCNOTE_DIR="$ROOT_DIR/syncnote"
API_DIR="$SYNCNOTE_DIR/api"
RPC_DIR="$SYNCNOTE_DIR/rpc"
# 定义 SQL 目录
SQL_DIR="$RPC_DIR/internal/model/"

MODE="${1:-all}"
RUN_TESTS="${2:-test}"

usage() {
  cat <<'EOF'
Usage:
  ./scripts/gen.sh [all|api|rpc|model] [test|no-test]

Examples:
  ./scripts/gen.sh                  # 生成所有 (API + RPC + Model)
  ./scripts/gen.sh model            # 仅生成 Model (自动扫描 sql 目录下所有 .sql 文件)
EOF
}

if ! command -v goctl >/dev/null 2>&1; then
  echo "Error: goctl not found."
  exit 1
fi

# 生成 Model 的函数 (支持处理多个文件)
generate_models() {
  if [[ ! -d "$SQL_DIR" ]]; then
    echo "Warning: SQL directory not found: $SQL_DIR"
    return
  fi

  # 查找目录下所有的 .sql 文件
  sql_files=$(find "$SQL_DIR" -maxdepth 1 -name "*.sql" -type f | sort)
  
  if [[ -z "$sql_files" ]]; then
    echo "Warning: No .sql files found in $SQL_DIR"
    return
  fi

  echo "Found SQL files:"
  echo "$sql_files"
  echo "----------------"

  for sql_file in $sql_files; do
    filename=$(basename "$sql_file")
    echo "Generating Model from: $filename ..."
    
    # 关键点：goctl model mysql ddl 支持单个文件输入
    # 它会自动解析文件内的所有 CREATE TABLE 语句
    # --dir 指定输出到同一个目录，goctl 会自动处理包名和文件追加
    (cd "$RPC_DIR" && goctl model mysql ddl -src="$sql_file" -dir="./internal/model" --cache --style=goZero)
    
    echo "✅ Generated for $filename"
  done
}

case "$MODE" in
  all)
    echo "=== Starting Full Generation ==="
    
    echo "[1/4] Generating Models..."
    generate_models

    echo "[2/4] Generating API..."
    (cd "$API_DIR" && goctl api go -api SyncNote.api -dir .)

    echo "[3/4] Generating RPC..."
    (cd "$RPC_DIR" && goctl rpc protoc syncNoterpc.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.)

    echo "[4/4] Formatting code..."
    (cd "$ROOT_DIR" && go fmt ./...)
    ;;

  api)
    echo "Generating API..."
    (cd "$API_DIR" && goctl api go -api SyncNote.api -dir .)
    (cd "$ROOT_DIR" && go fmt ./...)
    ;;

  rpc)
    echo "Generating RPC..."
    (cd "$RPC_DIR" && goctl rpc protoc syncNoterpc.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.)
    (cd "$ROOT_DIR" && go fmt ./...)
    ;;

  model)
    echo "Generating Models..."
    generate_models
    (cd "$ROOT_DIR" && go fmt ./...)
    ;;

  *)
    echo "Error: invalid mode '$MODE'"
    usage
    exit 1
    ;;
esac

if [[ "$RUN_TESTS" == "test" ]]; then
  echo "Running tests..."
  (cd "$SYNCNOTE_DIR" && go test ./... || echo "Tests finished with errors (check DB connection)")
fi

echo "✅ Done."