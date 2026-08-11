UPDATE projects SET tech_stack = '{}' WHERE tech_stack IS NULL;

ALTER TABLE projects
    ALTER COLUMN tech_stack SET DEFAULT '{}',
    ALTER COLUMN tech_stack SET NOT NULL;
