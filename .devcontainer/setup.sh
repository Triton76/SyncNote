#!/bin/bash
set -e

echo "🚀 开始配置 SyncNote 开发环境..."

# 1. 下载 Go 依赖
echo "📦 下载 Go 模块依赖..."
cd /workspaces/SyncNote
go mod download

# 1.1 启动本地依赖服务（优先复用仓库内 docker-compose）
if command -v docker >/dev/null 2>&1 && [ -f /workspaces/SyncNote/docker-compose.yml ]; then
    echo "🐳 启动 MySQL 依赖容器..."
    docker compose up -d mysql >/dev/null 2>&1 || true
fi

# 2. 验证工具安装
echo "🔧 验证工具安装..."
go version
goctl --help >/dev/null
protolint version
golangci-lint version
mysql --version

# 3. 初始化数据库（等待 MySQL 就绪）
echo "🗄️ 等待 MySQL 就绪..."
MYSQL_READY=0
for i in {1..30}; do
    if mysql -h 127.0.0.1 -u root -pdevpass123 -e "SELECT 1" &>/dev/null; then
        MYSQL_READY=1
        echo "✅ MySQL 已就绪！"
        break
    fi
    if mysql -h 127.0.0.1 -u root -e "SELECT 1" &>/dev/null; then
        MYSQL_READY=1
        echo "✅ MySQL 已就绪！"
        break
    fi
    echo "⏳ 等待 MySQL 启动... ($i/30)"
    sleep 2
done

if [ "$MYSQL_READY" -ne 1 ]; then
    echo "❌ MySQL 在预期时间内未就绪，请检查 docker compose logs mysql"
    exit 1
fi

# 4. 创建数据库（如果不存在）
echo "📝 创建数据库..."
mysql -h 127.0.0.1 -u root -pdevpass123 -e "CREATE DATABASE IF NOT EXISTS syncnote DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 5. 执行数据库迁移（如果有迁移脚本）
# echo "🔄 执行数据库迁移..."
# mysql -h 127.0.0.1 -u root -pdevpass123 syncnote < /workspaces/SyncNote/auth/auth.sql

# 6. 生成 .env 文件（可选）
echo "📄 生成环境变量文件..."
cat > /workspaces/SyncNote/.env <<EOF
GO_ENV=development
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=devpass123
DB_NAME=syncnote
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