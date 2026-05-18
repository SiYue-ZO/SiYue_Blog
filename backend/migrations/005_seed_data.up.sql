-- Seed admin user (password: admin123)
INSERT INTO users (username, email, password_hash, role, status) VALUES
('admin', 'admin@siyue.dev', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin', 'active');

-- Seed tags
INSERT INTO tags (name, slug) VALUES
('Next.js', 'nextjs'),
('React', 'react'),
('Go', 'golang'),
('TypeScript', 'typescript'),
('Docker', 'docker'),
('PostgreSQL', 'postgresql'),
('Tailwind CSS', 'tailwindcss'),
('API', 'api'),
('DevOps', 'devops'),
('Git', 'git');

-- Seed posts
INSERT INTO posts (title, slug, excerpt, status, author_id, published_at) VALUES
('使用 Next.js 15 构建现代博客', 'building-blog-with-nextjs-15', '探索 Next.js 15 的新特性，以及如何利用 App Router 构建一个高性能的博客系统。', 'published', 1, NOW()),
('Go 语言实战：构建 RESTful API', 'go-restful-api', '学习如何使用 Go 和 Gin 框架构建一个生产级别的 RESTful API。', 'published', 1, NOW() - INTERVAL '5 days'),
('Docker 容器化最佳实践', 'docker-best-practices', '深入了解 Docker 容器化的最佳实践，包括多阶段构建、安全性和性能优化。', 'published', 1, NOW() - INTERVAL '10 days');

-- Seed post_tags
INSERT INTO post_tags (post_id, tag_id) VALUES
(1, 1), (1, 2),  -- Next.js, React
(2, 3), (2, 8),  -- Go, API
(3, 5), (3, 9);  -- Docker, DevOps

-- Seed projects
INSERT INTO projects (title, description, repo_url, demo_url, tech_stack, is_featured, sort_order) VALUES
('SiYue Blog', '一个现代化的个人博客系统，使用 Next.js 15 和 Go 构建', 'https://github.com/siyue/blog', 'https://siyue.dev', ARRAY['Next.js', 'Go', 'PostgreSQL', 'Docker'], true, 1),
('Task Manager API', '一个高性能的任务管理 RESTful API', 'https://github.com/siyue/task-manager', '', ARRAY['Go', 'Gin', 'PostgreSQL', 'Redis'], true, 2),
('React Component Library', '一套可复用的 React 组件库', 'https://github.com/siyue/ui-components', 'https://ui.siyue.dev', ARRAY['React', 'TypeScript', 'Tailwind CSS'], true, 3);
