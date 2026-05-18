-- tags table
CREATE TABLE tags (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(50)  NOT NULL,
    slug       VARCHAR(50)  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_tags_slug ON tags (slug);
CREATE UNIQUE INDEX idx_tags_name ON tags (name);

-- post_tags junction table
CREATE TABLE post_tags (
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    tag_id  BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);

CREATE INDEX idx_post_tags_tag ON post_tags (tag_id);
