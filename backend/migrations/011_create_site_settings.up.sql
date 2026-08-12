CREATE TABLE site_settings (
    id                     SMALLINT    PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    project_page_visible   BOOLEAN     NOT NULL DEFAULT TRUE,
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO site_settings (id, project_page_visible)
VALUES (1, TRUE)
ON CONFLICT (id) DO NOTHING;

CREATE TRIGGER trg_site_settings_updated_at
    BEFORE UPDATE ON site_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
