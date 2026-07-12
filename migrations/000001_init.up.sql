CREATE EXTENSION IF NOT EXISTS "postgis";

CREATE TABLE admins (
    admin_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE soil_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE plots (
    plot_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    description TEXT,
    owner_id UUID NOT NULL REFERENCES users(user_id),
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    soil_type INT NOT NULL REFERENCES soil_types(id),
    boundary GEOMETRY(Polygon, 3857) NOT NULL,
    area_sq_m DECIMAL(12,2),
    grid_cell_size DOUBLE PRECISION NOT NULL DEFAULT 0.5,
    grid_cols INT NOT NULL,
    grid_rows INT NOT NULL
);
CREATE INDEX idx_plots_owner_id_is_deleted ON plots (owner_id, is_deleted);
CREATE INDEX idx_plots_soil_type ON plots (soil_type);
CREATE INDEX idx_plots_boundary_gist ON plots USING GIST(boundary);

CREATE TABLE crop_families (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE crops (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    family_id INT NOT NULL REFERENCES crop_families(id),
    vegetation_days_avg INT NOT NULL,
    sun_needs INT, -- shade, partial, full
    is_deleted BOOLEAN DEFAULT FALSE
);
CREATE INDEX idx_crops_family_id ON crops (family_id);

CREATE TABLE crop_rules (
    rule_id SERIAL PRIMARY KEY,
    subject_crop_id INT REFERENCES crops(id),
    subject_family_id INT REFERENCES crop_families(id),
    context_type INT NOT NULL,
    context_crop_id INT REFERENCES crops(id),
    context_family_id INT REFERENCES crop_families(id),
    return_after_days INT NOT NULL,
    score_modifier INT NOT NULL,
    explanation TEXT NOT NULL,
    priority INT DEFAULT 1,
    CHECK (
        (subject_crop_id IS NOT NULL OR subject_family_id IS NOT NULL) AND
        (context_crop_id IS NOT NULL OR context_family_id IS NOT NULL)
    )
);
CREATE INDEX idx_crop_rules_subject_crop_id ON crop_rules (subject_crop_id);
CREATE INDEX idx_crop_rules_subject_family_id ON crop_rules (subject_family_id);
CREATE INDEX idx_crop_rules_context_crop_id ON crop_rules (context_crop_id);
CREATE INDEX idx_crop_rules_context_family_id ON crop_rules (context_family_id);

CREATE TABLE event_store (
    event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plot_id UUID NOT NULL REFERENCES plots(plot_id),
    event_type INT NOT NULL,
    payload JSONB NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    sequence_number BIGINT NOT NULL,
    UNIQUE (plot_id, sequence_number)
);
CREATE INDEX idx_event_store_plot_id_sequence_number ON event_store (plot_id, sequence_number);
CREATE INDEX idx_event_store_event_type_occurred_at ON event_store (event_type, occurred_at);

CREATE TABLE grid_cells (
    plot_id UUID REFERENCES plots(plot_id),
    x_index INT,
    y_index INT,
    geom GEOMETRY(Polygon, 3857) NOT NULL,
    shade_level INT, -- shade, partial, full
    current_crop_id INT REFERENCES crops(id),
    is_occupied BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (plot_id, x_index, y_index)
);
CREATE INDEX idx_grid_cells_plot_id_x_index_y_index ON grid_cells (plot_id, x_index, y_index);
CREATE INDEX idx_grid_cells_current_crop_id ON grid_cells (current_crop_id);
CREATE INDEX idx_grid_cells_geom_gist ON grid_cells USING GIST(geom);

CREATE TABLE cell_crop_history (
    history_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plot_id UUID NOT NULL REFERENCES plots(plot_id),
    x_index INT NOT NULL,
    y_index INT NOT NULL,
    crop_id INT NOT NULL REFERENCES crops(id),
    family_id INT NOT NULL REFERENCES crop_families(id),
    bed_id UUID,
    plant_date DATE NOT NULL,
    harvest_date DATE,
    metadata JSONB,
    FOREIGN KEY (plot_id, x_index, y_index)
       REFERENCES grid_cells(plot_id, x_index, y_index)
       ON DELETE RESTRICT
);
CREATE INDEX idx_cell_crop_history_plot_id_x_index_y_index_plant_date_desc ON cell_crop_history (plot_id, x_index, y_index, plant_date DESC);
CREATE INDEX idx_cell_crop_history_plot_id_family_id_plant_date_desc ON cell_crop_history (plot_id, family_id, plant_date DESC);
CREATE INDEX idx_cell_crop_history_crop_id ON cell_crop_history (crop_id);

CREATE TABLE beds_ui (
    bed_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plot_id UUID NOT NULL REFERENCES plots(plot_id),
    name VARCHAR(100) NOT NULL,
    geom GEOMETRY(Polygon, 3857) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE
);
CREATE INDEX idx_beds_ui_plot_id_is_deleted ON beds_ui (plot_id, is_deleted);
CREATE INDEX idx_beds_ui_geom_gist ON beds_ui USING GIST(geom);

CREATE TABLE objects_ui (
    object_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plot_id UUID NOT NULL REFERENCES plots(plot_id),
    name VARCHAR(100) NOT NULL,
    object_type INT,
    geom GEOMETRY(Polygon, 3857) NOT NULL
);
CREATE INDEX idx_objects_ui_plot_id ON objects_ui (plot_id);
CREATE INDEX idx_objects_ui_geom_gist ON objects_ui USING GIST(geom);

CREATE TABLE timeline_events (
    event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plot_id UUID NOT NULL REFERENCES plots(plot_id),
    event_type INT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    display_title VARCHAR(255) NOT NULL,
    display_description TEXT,
    affected_cells TEXT[],
    metadata JSONB,
    source_event_ids UUID[]
);
CREATE INDEX idx_timeline_events_plot_id_occurred_at_desc ON timeline_events (plot_id, occurred_at DESC);
CREATE INDEX idx_timeline_gin_cells ON timeline_events USING GIN(affected_cells);

CREATE TABLE snapshots (
    snapshot_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plot_id UUID NOT NULL REFERENCES plots(plot_id),
    snapshot_date DATE NOT NULL,
    snapshot_type INT NOT NULL,
    grid_state JSONB NOT NULL,
    beds_state JSONB NOT NULL,
    objects_state JSONB NOT NULL,
    last_event_sequence BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (plot_id, snapshot_date, snapshot_type)
);
CREATE INDEX idx_snapshots_plot_id_snapshot_date_desc ON snapshots (plot_id, snapshot_date DESC);