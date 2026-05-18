# SiYue Blog

A modern personal blog and portfolio built with Next.js 15 and Go.

[中文文档](README.zh-CN.md)

## Tech Stack

### Frontend
- **Next.js 15** - React framework with App Router
- **TypeScript** - Type safety
- **Tailwind CSS** - Utility-first CSS
- **shadcn/ui** - UI components
- **Framer Motion** - Animations
- **MDX** - Blog content with code highlighting

### Backend
- **Go** - Programming language
- **Gin** - HTTP framework
- **PostgreSQL** - Database
- **GORM** - ORM
- **JWT** - Authentication

### DevOps
- **Docker** - Containerization
- **Nginx** - Reverse proxy
- **GitHub Actions** - CI/CD

## Getting Started

### Prerequisites
- Node.js 20+
- Go 1.23+
- PostgreSQL 16+
- Docker (optional)

### Quick Start with Docker

```bash
# Clone the repository
git clone https://github.com/siyue/blog.git
cd blog

# Start all services
docker-compose up -d

# Access the application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080/api/v1
# Admin Panel: http://localhost:3000/admin
```

### Manual Development Setup

#### Backend

```bash
cd backend

# Copy config template
cp config.example.yaml config.yaml

# Edit config.yaml with your database credentials

# Run database migrations
psql -U siyue -d siyue_blog -f migrations/001_create_users.up.sql
psql -U siyue -d siyue_blog -f migrations/002_create_posts.up.sql
psql -U siyue -d siyue_blog -f migrations/003_create_tags.up.sql
psql -U siyue -d siyue_blog -f migrations/004_create_projects.up.sql
psql -U siyue -d siyue_blog -f migrations/005_seed_data.up.sql

# Start the server
go run ./cmd/server
```

#### Frontend

```bash
cd frontend

# Install dependencies
npm install

# Copy environment template
cp .env.local.example .env.local

# Start development server
npm run dev
```

## Project Structure

```
SiYue_Blog/
├── frontend/              # Next.js 15 application
│   ├── src/
│   │   ├── app/           # App Router pages
│   │   ├── components/    # React components
│   │   ├── lib/           # Utility functions
│   │   └── content/       # MDX blog posts
│   └── public/            # Static assets
├── backend/               # Go + Gin API
│   ├── cmd/               # Entry points
│   ├── internal/          # Business logic
│   │   ├── entity/        # Domain entities
│   │   ├── usecase/       # Use cases
│   │   ├── repo/          # Repositories
│   │   └── controller/    # HTTP handlers
│   ├── pkg/               # Shared packages
│   └── migrations/        # SQL migrations
├── nginx/                 # Nginx configuration
├── docker-compose.yml     # Development environment
└── docker-compose.prod.yml # Production environment
```

## API Documentation

### Public API (`/api/v1/`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/auth/login` | User login |
| POST | `/auth/register` | User registration |
| POST | `/auth/refresh` | Refresh access token |
| GET | `/posts` | List posts (with pagination) |
| GET | `/posts/featured` | List featured posts |
| GET | `/posts/:slug` | Get post by slug |
| GET | `/tags` | List all tags |
| GET | `/projects` | List all projects |
| GET | `/projects/featured` | List featured projects |

### Admin API (`/api/admin/`)

Requires JWT authentication with admin role.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/posts` | List all posts (including drafts) |
| POST | `/posts` | Create new post |
| PUT | `/posts/:id` | Update post |
| DELETE | `/posts/:id` | Delete post |

## Features

### Blog System
- MDX support with custom components
- Code syntax highlighting
- Table of contents generation
- Reading progress indicator
- Tag-based filtering
- Full-text search

### SEO
- Automatic sitemap generation
- robots.txt
- Open Graph metadata
- Dynamic meta tags per page

### Performance
- Static site generation (SSG)
- Image optimization
- Gzip compression
- Static asset caching

### Security
- JWT authentication
- CORS protection
- Rate limiting
- Security headers

## Environment Variables

### Backend
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

### Frontend
```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_SITE_URL=http://localhost:3000
```

## Deployment

### Production with Docker

```bash
# Set environment variables
export DB_PASSWORD=your_secure_password
export JWT_ACCESS_SECRET=your_access_secret
export JWT_REFRESH_SECRET=your_refresh_secret
export NEXT_PUBLIC_API_URL=https://your-domain.com/api/v1
export NEXT_PUBLIC_SITE_URL=https://your-domain.com

# Start production environment
docker-compose -f docker-compose.prod.yml up -d
```

### CI/CD

The project includes GitHub Actions workflows for:
- Linting (golangci-lint, ESLint)
- Testing (Go tests, frontend build)
- Building and pushing Docker images

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Next.js](https://nextjs.org/)
- [Gin](https://gin-gonic.com/)
- [Tailwind CSS](https://tailwindcss.com/)
- [shadcn/ui](https://ui.shadcn.com/)
- [Framer Motion](https://www.framer.com/motion/)
