DROP INDEX IF EXISTS idx_projects_visible_order;
ALTER TABLE projects DROP COLUMN IF EXISTS is_visible;
