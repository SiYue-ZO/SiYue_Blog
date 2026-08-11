UPDATE users
SET password_hash = '!disabled-legacy-development-credential!'
WHERE username = 'admin'
  AND email = 'admin@siyue.dev'
  AND password_hash = '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy';

-- Remove only the untouched placeholder content shipped by the old seed.
DELETE FROM posts
WHERE slug IN ('building-blog-with-nextjs-15', 'go-restful-api', 'docker-best-practices')
  AND COALESCE(BTRIM(content), '') = '';

DELETE FROM projects
WHERE repo_url IN (
    'https://github.com/siyue/blog',
    'https://github.com/siyue/task-manager',
    'https://github.com/siyue/ui-components'
);

CREATE TYPE contact_status AS ENUM ('unread', 'read', 'archived');

CREATE TABLE contact_messages (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(100)   NOT NULL,
    email      VARCHAR(255)   NOT NULL,
    message    TEXT           NOT NULL,
    status     contact_status NOT NULL DEFAULT 'unread',
    ip_hash    VARCHAR(64)     NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_contact_messages_status_created
    ON contact_messages (status, created_at DESC);

CREATE TRIGGER trg_contact_messages_updated_at
    BEFORE UPDATE ON contact_messages
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
