# 本地开发指南

## 1. 前提条件

- Docker Engine + Compose 插件
- 或者本机安装
  - Node.js 24
  - Go 1.26
  - PostgreSQL 18
  - pnpm 11（通过 corepack 使用）

## 2. 推荐：Docker Compose 开发

```bash
cd /home/sy/Desktop/blog/SiYue_Blog
docker compose up --build
```

- 访问网站：`http://localhost:3000`
- 管理后台：`http://localhost:3000/admin/login`
- 后端健康检查：`http://localhost:8080/api/v1/health`

默认本地管理员：
- 用户名：`admin@example.local`
- 密码：`local-admin-password-change-me`

停止服务：

```bash
docker compose down
```

如果需要删除本地数据库数据：

```bash
docker compose down -v
```

## 3. 原生开发（不使用 Docker）

### 3.1 后端

```bash
cd /home/sy/Desktop/blog/SiYue_Blog/backend
cp config.example.yaml config.yaml
go run ./cmd/server
```

后端默认监听 `http://localhost:8080`。

### 3.2 前端

```bash
cd /home/sy/Desktop/blog/SiYue_Blog/frontend
corepack enable
pnpm install --frozen-lockfile
pnpm dev
```

前端默认监听 `http://localhost:3000`。

> 本地直接运行前端时，如果未设置 `API_INTERNAL_URL`，会回退到 `http://localhost:8080/api/v1`。

## 4. 调试与检查

- 查看 Docker 日志：
  - `docker compose logs -f backend`
  - `docker compose logs -f frontend`
- 后端代码检查：
  - `gofmt -w .`
  - `go vet ./...`
  - `go test -race ./...`
- 前端代码检查：
  - `pnpm lint`
  - `pnpm typecheck`
  - `pnpm test`

## 5. 快速排错

- 若 Docker 启动失败，先确认 `docker compose logs` 中的错误信息。
- 如果后端无法连接数据库，检查 PostgreSQL 服务是否已健康启动。
- 如果前端请求 API 失败，确认前端是否已正确访问 `http://localhost:8080`。
