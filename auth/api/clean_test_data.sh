#!/bin/bash
set -e

echo "🧹 正在清理测试数据..."

REDIS_HOST="${REDIS_HOST:-127.0.0.1}"
REDIS_PORT="${REDIS_PORT:-6379}"
MYSQL_HOST="127.0.0.1"
MYSQL_PORT="3306"
MYSQL_USER="root"
MYSQL_PASS="devpass123"
MYSQL_DB="syncnote"

# Step 1: Clear MySQL test data
echo "[1/2] Clearing MySQL test data..."
mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_PASS" "$MYSQL_DB" < ./test.delete.sql
if [ $? -eq 0 ]; then
  echo "✅ MySQL cleanup done."
else
  echo "❌ MySQL cleanup failed, check connection or SQL syntax."
  exit 1
fi

# Step 2: Clear Redis cache via Docker or local redis-cli
echo "[2/2] Clearing Redis cache..."

REDIS_CLEARED=0

# Try Docker first (preferred for containerized Redis)
if command -v docker >/dev/null 2>&1; then
  if docker exec redis redis-cli FLUSHDB 2>/dev/null; then
    echo "✅ Redis cache cleared via Docker."
    REDIS_CLEARED=1
  fi
fi

# Fallback to local redis-cli if Docker approach failed
if [ $REDIS_CLEARED -eq 0 ] && command -v redis-cli >/dev/null 2>&1; then
  if redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" FLUSHDB 2>/dev/null; then
    echo "✅ Redis cache cleared via redis-cli."
    REDIS_CLEARED=1
  fi
fi

# If both failed, print helpful message
if [ $REDIS_CLEARED -eq 0 ]; then
  echo "⚠️  Could not clear Redis cache automatically."
  echo "   You can manually clear it with one of these commands:"
  echo "   - Docker:    docker exec redis redis-cli FLUSHDB"
  echo "   - redis-cli: redis-cli -h $REDIS_HOST -p $REDIS_PORT FLUSHDB"
fi

echo "✅ All test data cleared successfully!"