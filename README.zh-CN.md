# SiYue Blog

一个使用 Next.js 15 和 Go 构建的现代化个人博客系统。

## 技术栈

### 前端
- **Next.js 15** - React 框架，使用 App Router
- **TypeScript** - 类型安全
- **Tailwind CSS** - 工具类优先的 CSS 框架
- **shadcn/ui** - UI 组件库
- **Framer Motion** - 动画库
- **MDX** - 博客内容，支持代码高亮

### 后端
- **Go** - 编程语言
- **Gin** - HTTP 框架
- **PostgreSQL** - 数据库
- **GORM** - ORM 框架
- **JWT** - 身份认证

### DevOps
- **Docker** - 容器化
- **Nginx** - 反向代理
- **GitHub Actions** - CI/CD

## 快速开始

### 环境要求
- Node.js 20+
- Go 1.23+
- PostgreSQL 16+
- Docker（可选）

### 使用 Docker 快速启动

```bash
# 克隆仓库
git clone https://github.com/siyue/blog.git
cd blog

# 启动所有服务
docker-compose up -d

# 访问应用
# 前端: http://localhost:3000
# 后端 API: http://localhost:8080/api/v1
# 管理后台: http://localhost:3000/admin
```

### 手动开发环境搭建

#### 后端

```bash
cd backend

# 复制配置模板
cp config.example.yaml config.yaml

# 编辑 config.yaml，配置数据库连接信息

# 执行数据库迁移
psql -U siyue -d siyue_blog -f migrations/001_create_users.up.sql
psql -U siyue -d siyue_blog -f migrations/002_create_posts.up.sql
psql -U siyue -d siyue_blog -f migrations/003_create_tags.up.sql
psql -U siyue -d siyue_blog -f migrations/004_create_projects.up.sql
psql -U siyue -d siyue_blog -f migrations/005_seed_data.up.sql

# 启动服务器
go run ./cmd/server
```

#### 前端

```bash
cd frontend

# 安装依赖
npm install

# 复制环境变量模板
cp .env.local.example .env.local

# 启动开发服务器
npm run dev
```

## 项目结构

```
SiYue_Blog/
├── frontend/              # Next.js 15 应用
│   ├── src/
│   │   ├── app/           # App Router 页面
│   │   ├── components/    # React 组件
│   │   ├── lib/           # 工具函数
│   │   └── content/       # MDX 博客文章
│   └── public/            # 静态资源
├── backend/               # Go + Gin API
│   ├── cmd/               # 入口文件
│   ├── internal/          # 业务逻辑
│   │   ├── entity/        # 领域实体
│   │   ├── usecase/       # 用例层
│   │   ├── repo/          # 仓储层
│   │   └── controller/    # HTTP 处理器
│   ├── pkg/               # 共享包
│   └── migrations/        # SQL 迁移文件
├── nginx/                 # Nginx 配置
├── docker-compose.yml     # 开发环境
└── docker-compose.prod.yml # 生产环境
```

## API 文档

### 公开 API (`/api/v1/`)

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/health` | 健康检查 |
| POST | `/auth/login` | 用户登录 |
| POST | `/auth/register` | 用户注册 |
| POST | `/auth/refresh` | 刷新访问令牌 |
| GET | `/posts` | 文章列表（分页） |
| GET | `/posts/featured` | 精选文章 |
| GET | `/posts/:slug` | 文章详情 |
| GET | `/tags` | 标签列表 |
| GET | `/projects` | 项目列表 |
| GET | `/projects/featured` | 精选项目 |

### 管理 API (`/api/admin/`)

需要 JWT 身份认证和管理员角色。

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/posts` | 获取所有文章（含草稿） |
| POST | `/posts` | 创建新文章 |
| PUT | `/posts/:id` | 更新文章 |
| DELETE | `/posts/:id` | 删除文章 |

## 功能特性

### 博客系统
- MDX 支持，可自定义组件
- 代码语法高亮
- 自动生成目录
- 阅读进度指示器
- 基于标签的筛选
- 全文搜索

### SEO 优化
- 自动生成 sitemap.xml
- robots.txt
- Open Graph 元数据
- 每个页面的动态 meta 标签

### 性能优化
- 静态站点生成（SSG）
- 图片优化
- Gzip 压缩
- 静态资源缓存

### 安全特性
- JWT 身份认证
- CORS 保护
- 速率限制
- 安全响应头

## 环境变量

### 后端
```yaml
# config.yaml
app:
  name: siyue-blog
  env: development

http:
  addr: ":8080"

postgres:
  host: localhost
  port: 5432
  user: siyue
  password: your_password
  dbname: siyue_blog

jwt:
  access_secret: your_access_secret
  refresh_secret: your_refresh_secret
```

### 前端
```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_SITE_URL=http://localhost:3000
```

## 部署

### 使用 Docker 进行生产部署

```bash
# 设置环境变量
export DB_PASSWORD=your_secure_password
export JWT_ACCESS_SECRET=your_access_secret
export JWT_REFRESH_SECRET=your_refresh_secret
export NEXT_PUBLIC_API_URL=https://your-domain.com/api/v1
export NEXT_PUBLIC_SITE_URL=https://your-domain.com

# 启动生产环境
docker-compose -f docker-compose.prod.yml up -d
```

### CI/CD

项目包含 GitHub Actions 工作流：
- 代码检查（golangci-lint, ESLint）
- 测试（Go 测试，前端构建）
- 构建和推送 Docker 镜像

## 贡献指南

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 开源协议

本项目基于 MIT 协议开源 - 详见 [LICENSE](LICENSE) 文件

## 致谢

- [Next.js](https://nextjs.org/)
- [Gin](https://gin-gonic.com/)
- [Tailwind CSS](https://tailwindcss.com/)
- [shadcn/ui](https://ui.shadcn.com/)
- [Framer Motion](https://www.framer.com/motion/)
