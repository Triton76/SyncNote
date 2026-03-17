#!/bin/bash
echo "🧹 正在清理测试数据..."

mysql -h 127.0.0.1 -P 3306 -u root -pdevpass123 syncnote < ./test.delete.sql

if [ $? -eq 0 ]; then
  echo "✅ 清理完成！"
else
  echo "❌ 清理失败，请检查数据库连接或 SQL 语法。"
fi