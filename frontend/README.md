# SyncNote Demo Frontend

## 功能

- 用户注册 / 登录并获取 token
- 创建笔记
- 读取笔记（包含 403 禁止访问验证）
- 保存笔记（含 expectedVersion）
- 获取当前用户笔记列表

## 启动方式

在仓库根目录执行：

```bash
python3 -m http.server 5500
```

然后在浏览器打开：

```text
http://127.0.0.1:5500/frontend/
```

## 默认后端地址

- Auth Host: `http://127.0.0.1:8889`
- API Host: `http://127.0.0.1:8888`

你可以在页面上修改并保存地址。
