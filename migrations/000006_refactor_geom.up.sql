DROP INDEX IF EXISTS idx_beds_ui_geom_gist;

ALTER TABLE beds_ui
    DROP COLUMN IF EXISTS geom,
    ADD COLUMN x_start INT NOT NULL DEFAULT 0,
    ADD COLUMN y_start INT NOT NULL DEFAULT 0,
    ADD COLUMN width INT NOT NULL DEFAULT 1,
    ADD COLUMN height INT NOT NULL DEFAULT 1,
    ADD COLUMN current_crop_id INT REFERENCES crops(id),
    ADD COLUMN plant_date DATE;

CREATE INDEX IF NOT EXISTS idx_beds_ui_current_crop_id ON beds_ui (current_crop_id);

DROP INDEX IF EXISTS idx_objects_ui_geom_gist;

ALTER TABLE objects_ui
    DROP COLUMN IF EXISTS geom,
    ADD COLUMN x_start INT NOT NULL DEFAULT 0,
    ADD COLUMN y_start INT NOT NULL DEFAULT 0,
    ADD COLUMN width INT NOT NULL DEFAULT 1,
    ADD COLUMN height INT NOT NULL DEFAULT 1;