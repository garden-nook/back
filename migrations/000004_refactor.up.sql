ALTER TABLE plots DROP COLUMN IF EXISTS area_sq_m;

DROP INDEX IF EXISTS idx_grid_cells_geom_gist;
DROP INDEX IF EXISTS idx_grid_cells_current_crop_id;
ALTER TABLE grid_cells
    DROP COLUMN IF EXISTS geom,
    DROP COLUMN IF EXISTS current_crop_id,
    DROP COLUMN IF EXISTS is_occupied;

DELETE FROM cell_crop_history WHERE harvest_date IS NULL;
DROP INDEX IF EXISTS idx_cell_crop_history_plot_id_family_id_plant_date_desc;
ALTER TABLE cell_crop_history
    DROP COLUMN IF EXISTS family_id,
    DROP COLUMN IF EXISTS bed_id,
    DROP COLUMN IF EXISTS metadata,
    ALTER COLUMN harvest_date SET NOT NULL;

DROP INDEX IF EXISTS idx_beds_ui_plot_id_is_deleted;
CREATE INDEX IF NOT EXISTS idx_beds_ui_plot_id ON beds_ui (plot_id);
ALTER TABLE beds_ui DROP COLUMN IF EXISTS is_deleted;

DROP INDEX IF EXISTS idx_timeline_gin_cells;
ALTER TABLE timeline_events
    DROP COLUMN IF EXISTS affected_cells,
    DROP COLUMN IF EXISTS metadata,
    DROP COLUMN IF EXISTS source_event_ids,
    ADD COLUMN source_event_id UUID REFERENCES event_store(event_id);
CREATE INDEX IF NOT EXISTS idx_timeline_events_source_event_id ON timeline_events (source_event_id);