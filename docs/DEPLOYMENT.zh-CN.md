# 生产部署与运维

## 1. 服务器准备

推荐使用受支持的 Linux 发行版，安装 Docker Engine 与 Compose 插件。域名 A/AAAA 记录需指向服务器公网地址，防火墙和云安全组放行 TCP 80、TCP 443、UDP 443。数据库与应用端口不应暴露到公网。

在服务器上克隆仓库并创建配置：

```bash
cp .env.example .env
chmod 600 .env
```

至少替换以下值：`DOMAIN`、`DB_PASSWORD`、两个不同的 `JWT_*_SECRET`、`ADMIN_EMAIL`、`ADMIN_PASSWORD`、作者资料。JWT 密钥建议分别使用 64 个以上随机字符，管理员密码使用密码管理器生成。

### 从 PostgreSQL 16 升级

生产编排使用 PostgreSQL 18 和全新的 `postgres_data_v18` 卷。旧版 Compose 的 PG16 卷不会被自动复用，这是为了避免 PostgreSQL 跨大版本直接读取数据目录造成损坏。已有数据时必须先在 PG16 容器中执行完整 `pg_dump`，再导入全新的 PG18 临时库；核对用户、文章、标签、项目和留言数量并验证管理员登录后才切换。也可以按 PostgreSQL 官方流程使用 `pg_upgrade`。切勿把旧 PG16 named volume 改名后直接挂到 PG18；确认恢复和异地备份前不要删除旧卷。

## 2. 首次上线

```bash
docker compose -f docker-compose.prod.yml config
docker compose -f docker-compose.prod.yml up -d --build
docker compose -f docker-compose.prod.yml ps
docker compose -f docker-compose.prod.yml logs --tail=100 backend caddy
curl --fail https://你的域名/api/v1/health
```

后端在监听 HTTP 前自动加锁并执行未应用的 SQL 迁移，然后仅在邮箱不存在时创建 `.env` 指定的管理员。后续修改 `ADMIN_PASSWORD` 不会重置已有账号；密码轮换应通过显式管理流程完成。Caddy 根据 `DOMAIN` 自动申请证书。请勿在另一套环境中复用数据库密码、JWT 密钥或管理员密码。

## 3. 发布与升级

手动升级：

```bash
git pull --ff-only
docker compose -f docker-compose.prod.yml build
docker compose -f docker-compose.prod.yml up -d --remove-orphans
curl --fail https://你的域名/api/v1/health
```

生产 Compose 支持不可变镜像：设置 `BACKEND_IMAGE`、`FRONTEND_IMAGE`、`RELEASE_TAG` 后执行 `pull` 与 `up -d`。回滚时将 `RELEASE_TAG` 改为上一个 Git commit SHA 并重新启动。数据库迁移应保持向后兼容；回滚镜像前先确认旧版本能够读取新 schema。

## 4. GitHub Actions 自动发布

`CI` 在前后端检查通过后触发 `Publish and Deploy`，将带 commit SHA 和 `latest` 标签的镜像推送到 GHCR。

仓库 Variables：

- 必需：`DOMAIN`
- 站点资料：`SITE_NAME`、`SITE_DESCRIPTION`、`AUTHOR_NAME`、`AUTHOR_BIO`、`AUTHOR_EMAIL`、`AUTHOR_LOCATION`、`AUTHOR_GITHUB`、`AUTHOR_TWITTER`
- 自动 SSH 部署：`ENABLE_SSH_DEPLOY=true`、`DEPLOY_PATH`
- 私有 GHCR 包：可选 `GHCR_USERNAME`，默认使用仓库 owner

自动部署 Secrets：

- `DEPLOY_HOST`、`DEPLOY_USER`、`DEPLOY_SSH_KEY`
- `DEPLOY_KNOWN_HOSTS`：提前用可信渠道取得服务器 host key，不能在工作流中关闭校验
- `DEPLOY_ENV`：完整的生产 `.env` 内容；工作流以 `600` 权限写入服务器
- 私有 GHCR 包还需 `GHCR_READ_TOKEN`（至少具有 `read:packages` 权限）；公开包可不配置

服务器只需预先安装 Docker 并允许部署用户运行 Docker。工作流会创建 `DEPLOY_PATH`，同步生产 Compose、Caddy 与备份脚本，写入 `.env`，登录私有 GHCR（如配置），然后拉取不可变 commit SHA 镜像并检查公网健康状态。

## 5. 备份与恢复

`backup` 服务每天生成一个 `/backups/siyue_blog_UTC时间.sql.gz`，默认保留 14 天。备份保存在 Docker `backups` volume 中；这不是异地备份。应定期将其同步到加密对象存储，并实际演练恢复。

列出备份：

```bash
docker compose -f docker-compose.prod.yml exec backup find /backups -maxdepth 1 -type f -name 'siyue_blog_*.sql.gz' -print
```

恢复会覆盖逻辑数据，应先停止写入、额外备份并确认目标数据库。一个安全流程是把选中的备份复制到主机，在临时数据库中解压并用 `psql` 导入，验证文章、管理员登录和消息数量后再切换；不要未经演练直接导入生产库。

## 6. 日常检查

```bash
docker compose -f docker-compose.prod.yml ps
docker compose -f docker-compose.prod.yml logs --since=1h backend frontend caddy
docker compose -f docker-compose.prod.yml exec postgres pg_isready
curl --fail https://你的域名/api/v1/health
```

建议同时配置外部 HTTPS 可用性监控、磁盘/内存告警、证书到期告警和备份成功告警。Compose 已限制日志文件大小，但仍需监控 Docker volume 与系统磁盘。

## 7. 安全检查清单

- `.env` 未提交到 Git，权限为 `600`
- 80/443 之外的应用端口未公开，SSH 禁止密码登录
- `ADMIN_PASSWORD`、数据库密码、两个 JWT 密钥互不相同
- 管理员邮箱和作者联系方式已替换，公开页面没有占位信息
- 定期更新基础镜像与依赖，并让 CI 的测试和审计保持通过
- 备份已复制到异地，且最近一次恢复演练成功

仓库中的 `nginx/` 是旧配置，仅为历史兼容保留；生产编排实际使用 `caddy/Caddyfile`。
