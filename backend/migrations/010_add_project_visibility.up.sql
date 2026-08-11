ALTER TABLE projects
    ADD COLUMN is_visible BOOLEAN NOT NULL DEFAULT TRUE;

CREATE INDEX idx_projects_visible_order
    ON projects (is_visible, sort_order, created_at DESC);
