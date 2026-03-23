# SyncNote Frontend

## Vue 前端（推荐）

新前端位于 `frontend/vue-app`，由 Vite + Vue 3 脚手架生成。

### 已实现功能

- 用户注册、登录、登出
- 创建笔记、按 `note_id` 查询笔记
- 笔记前端分页（可配置每页大小）
- 创建团队、加入团队
- 授予笔记权限（read/write/admin）
- API Host 配置保存（localStorage）

### 本地运行

```bash
cd frontend/vue-app
npm install
npm run dev -- --host 0.0.0.0 --port 5173
```

浏览器访问：

```text
http://127.0.0.1:5173
```

### 默认后端地址

- Auth Base: `http://127.0.0.1:8000`
- SyncNote Base: `http://127.0.0.1:8001`

也可通过环境变量覆盖：

```bash
VITE_AUTH_BASE=http://127.0.0.1:8000
VITE_SYNC_BASE=http://127.0.0.1:8001
```

## 旧版静态 Demo

旧静态页面文件仍在 `frontend/index.html`、`frontend/app.js`、`frontend/styles.css`。
