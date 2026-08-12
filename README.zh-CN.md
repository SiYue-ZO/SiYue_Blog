# SiYue Blog

一个可直接自托管的现代个人博客系统。公开站点、内容管理后台、Go API、PostgreSQL、自动 HTTPS、数据库迁移和备份均已打通。

[English](README.md) · [生产部署](docs/DEPLOYMENT.zh-CN.md) · [API 概览](docs/API.zh-CN.md)

## 已实现功能

- 首页、文章列表/搜索/标签筛选、文章详情、项目、标签、关于、联系页面
- PostgreSQL 作为生产内容唯一数据源；本地 Markdown 只作为开发回退
- Markdown 安全渲染、文章目录、相关文章、阅读时间与响应式深浅色界面
- 管理后台：登录、统计、首页草稿与发布、文章/标签/项目 CRUD、项目公开展示开关、联系页社交账号、联系消息处理
- Sitemap、RSS、robots.txt、Open Graph、结构化数据、站点清单和错误页
- 短效访问令牌、可轮换且服务端可撤销的刷新会话、HttpOnly Cookie、角色校验
- 输入校验、速率限制、请求体限制、CORS/CSP/安全响应头、数据库内容过滤
- 启动时自动执行幂等 SQL 迁移并从环境变量引导管理员账号
- Docker Compose、Caddy 自动 HTTPS、PostgreSQL 定时备份和 GitHub Actions CI/CD

## 技术栈

- 前端：Next.js 16、React 19、TypeScript 6、Tailwind CSS 4、Framer Motion
- 后端：Go 1.26、Gin、GORM、PostgreSQL 18、JWT
- 运维：Docker Compose、Caddy、GitHub Actions、GHCR

## 本地启动

需要 Docker Engine 及 Compose 插件。

```bash
docker compose up --build
```

启动完成后：

- 站点：http://localhost:3000
- 后台：http://localhost:3000/admin/login
- API 健康检查：http://localhost:8080/api/v1/health
- 本地管理员：`admin@example.local` / `local-admin-password-change-me`

这些账号和密钥只用于回环地址上的本地开发，不能复制到生产环境。后端首次启动会自动执行 `backend/migrations` 中尚未执行的迁移。

停止服务：

```bash
docker compose down
```

如需连同本地数据库一起重建，可在明确不需要数据时执行 `docker compose down -v`。

## 不使用 Docker 开发

要求 Node.js 24、Go 1.26 和 PostgreSQL 18。

```bash
# 后端：先根据 backend/config.example.yaml 创建 backend/config.yaml
cd backend
go run ./cmd/server

# 前端：另开终端
cd frontend
corepack enable
pnpm install --frozen-lockfile
pnpm dev
```

浏览器请求使用同源 `/api`，Next.js 服务端通过 `API_INTERNAL_URL` 连接 Go API。本地直接运行前端时，未设置该变量会回退到 `http://localhost:8080/api/v1`。

## 质量检查

```bash
cd backend
gofmt -w .
go vet ./...
go test -race ./...

cd ../frontend
corepack enable
pnpm install --frozen-lockfile
pnpm check
```

## 生产部署

1. 准备一台安装 Docker 的 Linux 服务器，并将域名的 A/AAAA 记录指向它。
2. 放行 TCP 80、TCP 443 和 UDP 443。
3. 复制 `.env.example` 为 `.env`，填入真实域名、个人资料及全新的长随机密钥。
4. 运行：

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

Caddy 会自动申请和续期 TLS 证书；后端会自动迁移数据库，并仅在账号不存在时创建环境变量指定的管理员。后续修改环境变量不会覆盖现有管理员密码、角色或状态。生产环境缺少必要密钥、Cookie 未启用 Secure 或仍使用不安全配置时，后端会拒绝启动。

完整的上线、备份恢复、升级、回滚与 GitHub Actions 配置见 [生产部署指南](docs/DEPLOYMENT.zh-CN.md)。

## 目录

```text
frontend/                 Next.js 公开站点与管理后台
backend/                  Go API、领域逻辑与嵌入式迁移
backend/migrations/       按版本执行的 PostgreSQL SQL
caddy/                    生产反向代理与 HTTPS 配置
scripts/                  数据库备份脚本
.github/workflows/        CI、镜像发布及可选 SSH 部署
docker-compose.yml        本地开发环境
docker-compose.prod.yml   生产环境
```

## 发布内容

登录 `/admin/login` 后可直接创建文章、标签和项目。文章正文使用 Markdown，封面图填写公网 HTTPS URL；文章状态设为 `published` 后才会出现在公开站点。联系表单消息在后台“留言”页面处理。

## 许可证

本项目当前未附带开源许可证；在添加许可证前保留全部权利。
