#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-devpass123}"
DB_NAME="${DB_NAME:-syncnote}"
SCHEMA_SQL="${SCHEMA_SQL:-$ROOT_DIR/rebuild/deploy/sql/newmodels.sql}"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Error: missing command '$1'"
    exit 1
  fi
}

need_cmd mysql

usage() {
  cat <<'EOF'
Usage:
  ./rebuild/scripts/reset_sql.sh

Behavior:
  1) Drop and recreate target database
  2) Import schema SQL file
  3) Print imported table count

Environment variables:
  MYSQL_HOST       MySQL host (default: 127.0.0.1)
  MYSQL_PORT       MySQL port (default: 3306)
  MYSQL_USER       MySQL user (default: root)
  MYSQL_PASSWORD   MySQL password (default: devpass123)
  DB_NAME          Database name (default: syncnote)
  SCHEMA_SQL       SQL file path (default: rebuild/deploy/sql/newmodels.sql)

Examples:
  ./rebuild/scripts/reset_sql.sh
  MYSQL_PASSWORD=xxx ./rebuild/scripts/reset_sql.sh
  DB_NAME=syncnote_test SCHEMA_SQL=./rebuild/deploy/sql/newmodels.sql ./rebuild/scripts/reset_sql.sh
EOF
}

case "${1:-}" in
  -h|--help|help)
    usage
    exit 0
    ;;
esac

if [[ ! -f "$SCHEMA_SQL" ]]; then
  echo "Error: schema SQL not found: $SCHEMA_SQL"
  exit 1
fi

mysql_exec() {
  MYSQL_PWD="$MYSQL_PASSWORD" mysql \
    -h "$MYSQL_HOST" \
    -P "$MYSQL_PORT" \
    -u "$MYSQL_USER" \
    "$@"
}

echo "========================================"
echo "SQL reset"
echo "========================================"
echo "- host: $MYSQL_HOST:$MYSQL_PORT"
echo "- user: $MYSQL_USER"
echo "- db:   $DB_NAME"
echo "- sql:  $SCHEMA_SQL"

echo "[INFO] checking MySQL connectivity ..."
mysql_exec -e "SELECT 1;" >/dev/null
echo "[OK] MySQL is reachable"

echo "[INFO] resetting database '$DB_NAME' ..."
mysql_exec -e "DROP DATABASE IF EXISTS \`$DB_NAME\`; CREATE DATABASE \`$DB_NAME\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
echo "[OK] database recreated"

echo "[INFO] importing schema SQL ..."
mysql_exec "$DB_NAME" <"$SCHEMA_SQL"
echo "[OK] schema import done"

table_count="$(mysql_exec -N -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='$DB_NAME';")"
echo "[OK] table count in '$DB_NAME': $table_count"

echo "Done."
