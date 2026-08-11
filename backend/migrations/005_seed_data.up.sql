-- Safe starter taxonomy. Credentials and personal content are intentionally
-- not stored in source-controlled migrations.
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
('Git', 'git')
ON CONFLICT DO NOTHING;
