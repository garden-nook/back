UPDATE grid_cells SET shade_level = 1 WHERE shade_level IS NULL;
UPDATE grid_cells SET is_deleted = FALSE WHERE is_deleted IS NULL;

ALTER TABLE grid_cells
    ALTER COLUMN shade_level SET NOT NULL,
    ALTER COLUMN shade_level SET DEFAULT 1,         -- 1 = ясно / солнце
    ALTER COLUMN is_deleted SET NOT NULL,
    ALTER COLUMN is_deleted SET DEFAULT FALSE;