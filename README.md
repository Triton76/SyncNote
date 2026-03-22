# SyncNote

SyncNote 是一个基于 Go + Vue 的协作笔记演示项目，当前提供认证、笔记、团队加入等基础能力，并附带前端调试控制台。

注意：当前版本暂不包含团队授权相关内容。

## 目录概览

- [rebuild](rebuild): 后端服务与脚本
- [frontend/vue-app](frontend/vue-app): Vue 前端控制台
- [docker-compose.yml](docker-compose.yml): 本地依赖服务（MySQL、Redis、Etcd）

## 当前功能

- 登录与注册
- 笔记创建、按 ID 查询、编辑、删除
- 团队创建、加入
- 笔记用户授权（给用户授予 read/write/admin）
- 前端 Dashboard 调试信息面板

## 快速启动

前置依赖：

- Go
- Node.js 与 npm
- Docker（含 docker compose）

一键启动（包含基础设施与后端，默认也会启动前端）：

cd /workspaces/SyncNote/rebuild/scripts
./start_all.sh

启动后默认地址：

- Frontend: http://127.0.0.1:5173
- Auth API: http://127.0.0.1:8000
- Syncnote API: http://127.0.0.1:8001
- User API: http://127.0.0.1:8003

## 停止服务

cd /workspaces/SyncNote/rebuild/scripts
./stop_all.sh

如需保留 MySQL/Redis/Etcd 不关闭：

STOP_INFRA=0 ./stop_all.sh

## 重置数据库（测试推荐）

已内置 SQL 重置脚本：

cd /workspaces/SyncNote/rebuild/scripts
./reset_sql.sh

可选环境变量：

- MYSQL_HOST
- MYSQL_PORT
- MYSQL_USER
- MYSQL_PASSWORD
- DB_NAME
- SCHEMA_SQL

示例：

MYSQL_PASSWORD=devpass123 DB_NAME=syncnote ./reset_sql.sh

## 前端说明

前端位于 [frontend/vue-app](frontend/vue-app)，采用组件化 Dashboard：

- 登录/注册页
- API 配置卡片
- 团队操作卡片
- 笔记侧栏与详情编辑
- 调试信息面板

开发模式运行：

cd /workspaces/SyncNote/frontend/vue-app
npm install
npm run dev

构建校验：

cd /workspaces/SyncNote/frontend/vue-app
npm run build

## 常用测试脚本

- [rebuild/scripts/test_rebuild_all.sh](rebuild/scripts/test_rebuild_all.sh): 全量构建/冒烟
- [rebuild/scripts/test_syncnote_optimistic_lock.sh](rebuild/scripts/test_syncnote_optimistic_lock.sh): 乐观锁场景
- [rebuild/scripts/test_userapi.sh](rebuild/scripts/test_userapi.sh): userapi 测试
- [rebuild/scripts/test_authapi.sh](rebuild/scripts/test_authapi.sh): authapi 测试

## 常见问题

1. docker compose 启动报容器名冲突

- 已在启动脚本内做冲突处理。
- 若仍冲突，可先手动查看并清理：

  docker ps -a

2. 前端跨域问题

- 项目默认通过 Vite 代理走同源。
- API Base 建议留空，使用代理模式。

3. 数据异常或脏数据

- 先执行 SQL 重置脚本，再重启服务。

## 后续建议

- 补充团队权限能力（按团队授权/撤销）
- 增加后端分页列表接口（笔记、团队）
- 增加端到端集成测试与 API 文档
