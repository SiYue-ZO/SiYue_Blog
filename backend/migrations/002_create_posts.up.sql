-- posts table
CREATE TABLE posts (
    id             BIGSERIAL PRIMARY KEY,
    title          VARCHAR(255)  NOT NULL,
    slug           VARCHAR(255)  NOT NULL,
    excerpt        VARCHAR(500),
    content        TEXT,
    cover_image    VARCHAR(512),
    status         post_status   NOT NULL DEFAULT 'draft',
    view_count     INTEGER       NOT NULL DEFAULT 0,
    is_pinned      BOOLEAN       NOT NULL DEFAULT FALSE,
    author_id      BIGINT        NOT NULL REFERENCES users(id),
    published_at   TIMESTAMPTZ,
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_posts_slug ON posts (slug);
CREATE INDEX idx_posts_status_published_at ON posts (status, published_at DESC)
    WHERE status = 'published';
CREATE INDEX idx_posts_author ON posts (author_id);
CREATE INDEX idx_posts_search ON posts USING GIN (
    (title || ' ' || COALESCE(excerpt, '')) gin_trgm_ops
);

CREATE TRIGGER trg_posts_updated_at
    BEFORE UPDATE ON posts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
