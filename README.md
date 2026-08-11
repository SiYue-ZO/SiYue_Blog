# SiYue Blog

A production-oriented, self-hosted personal blog with a public website, database-backed content management, a Go API, automatic HTTPS, migrations, and backups.

[简体中文](README.zh-CN.md) · [Deployment guide (Chinese)](docs/DEPLOYMENT.zh-CN.md) · [API overview (Chinese)](docs/API.zh-CN.md)

## What is included

- Responsive public pages for posts, tags, projects, profile, and contact
- Database-backed admin UI for posts, tags, projects (including visibility control), messages, and metrics
- Safe Markdown rendering, search, table of contents, RSS, sitemap, and structured SEO data
- Rotating and revocable refresh sessions stored in HttpOnly cookies
- Validation, rate limiting, CSP/CORS/security headers, and role-based authorization
- Embedded SQL migrations and environment-driven admin bootstrap
- Docker Compose, Caddy automatic TLS, PostgreSQL backups, and GitHub Actions CI/CD

## Stack

- Next.js 16, React 19, TypeScript 6, Tailwind CSS 4
- pnpm 11 with a frozen lockfile
- Go 1.26, Gin, GORM, PostgreSQL 18
- Docker Compose, Caddy, GitHub Actions, GHCR

## Run locally

```bash
docker compose up --build
```

Open <http://localhost:3000>. The development admin is `admin@example.local` / `local-admin-password-change-me`; never reuse these credentials in production.

For native development and verification commands, see the [Chinese README](README.zh-CN.md).

## Deploy

Point a domain at a Linux server with Docker, copy `.env.example` to `.env`, replace every secret and profile value, then run:

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

See the [deployment guide](docs/DEPLOYMENT.zh-CN.md) before exposing the service publicly.

## License

No open-source license is currently granted. All rights are reserved until a license is added.
