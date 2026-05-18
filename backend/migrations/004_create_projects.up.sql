-- projects table
CREATE TABLE projects (
    id           BIGSERIAL PRIMARY KEY,
    title        VARCHAR(255)  NOT NULL,
    description  TEXT,
    repo_url     VARCHAR(512),
    demo_url     VARCHAR(512),
    image_url    VARCHAR(512),
    tech_stack   TEXT[],
    is_featured  BOOLEAN       NOT NULL DEFAULT FALSE,
    sort_order   INTEGER       NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_projects_featured ON projects (is_featured, sort_order)
    WHERE is_featured = TRUE;

CREATE TRIGGER trg_projects_updated_at
    BEFORE UPDATE ON projects
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
