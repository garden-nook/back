ALTER TABLE event_store       DROP CONSTRAINT IF EXISTS event_store_plot_id_fkey;
ALTER TABLE grid_cells        DROP CONSTRAINT IF EXISTS grid_cells_plot_id_fkey;
ALTER TABLE cell_crop_history DROP CONSTRAINT IF EXISTS cell_crop_history_plot_id_fkey;
ALTER TABLE beds_ui           DROP CONSTRAINT IF EXISTS beds_ui_plot_id_fkey;
ALTER TABLE objects_ui        DROP CONSTRAINT IF EXISTS objects_ui_plot_id_fkey;
ALTER TABLE timeline_events   DROP CONSTRAINT IF EXISTS timeline_events_plot_id_fkey;
ALTER TABLE snapshots         DROP CONSTRAINT IF EXISTS snapshots_plot_id_fkey;

ALTER TABLE event_store       ADD CONSTRAINT event_store_plot_id_fkey       FOREIGN KEY (plot_id) REFERENCES plots(plot_id) ON DELETE CASCADE;
ALTER TABLE grid_cells        ADD CONSTRAINT grid_cells_plot_id_fkey        FOREIGN KEY (plot_id) REFERENCES plots(plot_id) ON DELETE CASCADE;
ALTER TABLE cell_crop_history ADD CONSTRAINT cell_crop_history_plot_id_fkey FOREIGN KEY (plot_id) REFERENCES plots(plot_id) ON DELETE CASCADE;
ALTER TABLE beds_ui           ADD CONSTRAINT beds_ui_plot_id_fkey           FOREIGN KEY (plot_id) REFERENCES plots(plot_id) ON DELETE CASCADE;
ALTER TABLE objects_ui        ADD CONSTRAINT objects_ui_plot_id_fkey        FOREIGN KEY (plot_id) REFERENCES plots(plot_id) ON DELETE CASCADE;
ALTER TABLE timeline_events   ADD CONSTRAINT timeline_events_plot_id_fkey   FOREIGN KEY (plot_id) REFERENCES plots(plot_id) ON DELETE CASCADE;
ALTER TABLE snapshots         ADD CONSTRAINT snapshots_plot_id_fkey         FOREIGN KEY (plot_id) REFERENCES plots(plot_id) ON DELETE CASCADE;

DELETE FROM plots WHERE is_deleted = TRUE;

ALTER TABLE plots DROP COLUMN IF EXISTS is_deleted;

ALTER TABLE grid_cells DROP COLUMN IF EXISTS is_deleted;