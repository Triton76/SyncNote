#!/bin/bash
set -e

MYSQL_HOST="127.0.0.1"
MYSQL_PORT="3306"
MYSQL_USER="root"
MYSQL_PASSWORD="devpass123"
MYSQL_DATABASE="syncnote"

echo "🚀 开始配置 SyncNote 开发环境..."

# 1. 下载 Go 依赖
echo "📦 下载 Go 模块依赖..."
cd /workspaces/SyncNote
go mod download

# 1.2 安装前端依赖（如果 Vue 项目已创建）
if [ -f /workspaces/SyncNote/frontend/vue-app/package.json ] && command -v npm >/dev/null 2>&1; then
    echo "🌐 安装 Vue 前端依赖..."
    cd /workspaces/SyncNote/frontend/vue-app
    npm install
    cd /workspaces/SyncNote
fi

# 1.1 启动本地依赖服务（优先复用仓库内 docker-compose）
if command -v docker >/dev/null 2>&1 && [ -f /workspaces/SyncNote/docker-compose.yml ]; then
    echo "🐳 启动依赖容器 (etcd/redis/mysql)..."
    docker compose up -d etcd redis mysql >/dev/null 2>&1 || true
fi

# 2. 验证工具安装
echo "🔧 验证工具安装..."
go version
goctl --help >/dev/null
protolint version
golangci-lint version
mysql --version

# 3. 初始化数据库（等待 etcd/Redis/MySQL 就绪）
echo "🗄️ 等待依赖服务就绪..."

wait_for_port() {
    local name="$1"
    local host="$2"
    local port="$3"

    for i in {1..30}; do
        if (echo > "/dev/tcp/${host}/${port}") >/dev/null 2>&1; then
            echo "✅ ${name} 已就绪！"
            return 0
        fi
        echo "⏳ 等待 ${name} 启动... (${i}/30)"
        sleep 2
    done

    echo "❌ ${name} 在预期时间内未就绪"
    return 1
}

wait_for_mysql() {
    for i in {1..30}; do
        if mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" -e "SELECT 1" >/dev/null 2>&1; then
            echo "✅ MySQL 已就绪！"
            return 0
        fi
        if mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -e "SELECT 1" >/dev/null 2>&1; then
            echo "✅ MySQL 已就绪！"
            return 0
        fi
        echo "⏳ 等待 MySQL 启动... (${i}/30)"
        sleep 2
    done

    echo "❌ MySQL 在预期时间内未就绪，请检查 docker compose logs mysql"
    return 1
}

wait_for_port "etcd" "127.0.0.1" "2379"
wait_for_port "Redis" "127.0.0.1" "6379"
wait_for_mysql

# 4. 创建数据库（如果不存在）
echo "📝 创建数据库..."
mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" -e "CREATE DATABASE IF NOT EXISTS ${MYSQL_DATABASE} DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 5. 执行数据库迁移（如果有迁移脚本）
# echo "🔄 执行数据库迁移..."
# mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" < /workspaces/SyncNote/auth/auth.sql

# 6. 生成 .env 文件（可选）
echo "📄 生成环境变量文件..."
cat > /workspaces/SyncNote/.env <<EOF
GO_ENV=development
DB_HOST=$MYSQL_HOST
DB_PORT=$MYSQL_PORT
DB_USER=$MYSQL_USER
DB_PASSWORD=$MYSQL_PASSWORD
DB_NAME=$MYSQL_DATABASE
EOF

echo "✅ 开发环境配置完成！"
echo "================================"
echo "📌 可用命令："
echo "  - go run auth.go           # 启动服务"
echo "  - goctl api go -api ...    # 生成 API 代码"
echo "  - goctl model mysql ...    # 生成 Model 代码"
echo "  - protolint lint .         # 检查 Proto 文件"
echo "  - golangci-lint run        # 检查 Go 代码"
echo "================================"