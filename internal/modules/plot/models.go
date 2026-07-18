package plots

import (
	"time"
)

const DefaultGridCellSize = 0.5

type Plot struct {
	PlotID       string  `json:"plot_id" db:"plot_id"`
	Name         string  `json:"name" db:"name"`
	SoilTypeID   int32   `json:"soil_type" db:"soil_type"`
	SoilTypeName string  `json:"soil_name" db:"soil_name"`
	AreaSqM      float64 `json:"area_sq_m" db:"area_sq_m"`
	GridCellSize float64 `json:"grid_cell_size" db:"grid_cell_size"`
	GridCols     int     `json:"grid_cols" db:"grid_cols"`
	GridRows     int     `json:"grid_rows" db:"grid_rows"`
}

type CreatePlotModel struct {
	Name           string
	SoilTypeID     int32
	BoundaryWidth  float64
	BoundaryHeight float64
	AreaSqM        float64
	GridCellSize   float64
	GridCols       int
	GridRows       int
}

type Bed struct {
	BedID     string    `json:"bed_id" db:"bed_id"`
	PlotID    string    `json:"plot_id" db:"plot_id"`
	Name      string    `json:"name" db:"name"`
	Geom      string    `json:"geom" db:"geom"` // GeoJSON
	IsDeleted bool      `json:"-" db:"is_deleted"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UIObject struct {
	ObjectID   string    `json:"object_id" db:"object_id"`
	PlotID     string    `json:"plot_id" db:"plot_id"`
	Name       string    `json:"name" db:"name"`
	ObjectType int       `json:"object_type" db:"object_type"` // 1=tree, 2=building, 3=path, 4=shade
	Geom       string    `json:"geom" db:"geom"`               // GeoJSON
	IsDeleted  bool      `json:"-" db:"is_deleted"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type GridCell struct {
	PlotID        string  `json:"plot_id" db:"plot_id"`
	XIndex        int     `json:"x_index" db:"x_index"`
	YIndex        int     `json:"y_index" db:"y_index"`
	CurrentBedID  *string `json:"current_bed_id,omitempty" db:"current_bed_id"`
	CurrentCropID *int    `json:"current_crop_id,omitempty" db:"current_crop_id"`
	IsOccupied    bool    `json:"is_occupied" db:"is_occupied"`
	IsShaded      bool    `json:"is_shaded" db:"is_shaded"` // Затенение от объектов
}

type CellHistory struct {
	HistoryID   string     `json:"history_id" db:"history_id"`
	PlotID      string     `json:"plot_id" db:"plot_id"`
	XIndex      int        `json:"x_index" db:"x_index"`
	YIndex      int        `json:"y_index" db:"y_index"`
	CropID      int        `json:"crop_id" db:"crop_id"`
	CropName    string     `json:"crop_name" db:"crop_name"`
	FamilyID    int        `json:"family_id" db:"family_id"`
	FamilyName  string     `json:"family_name" db:"family_name"`
	BedID       *string    `json:"bed_id,omitempty" db:"bed_id"`
	BedName     *string    `json:"bed_name,omitempty" db:"bed_name"`
	PlantDate   time.Time  `json:"plant_date" db:"plant_date"`
	HarvestDate *time.Time `json:"harvest_date,omitempty" db:"harvest_date"`
	YieldKg     *float64   `json:"yield_kg,omitempty" db:"yield_kg"`
}

type PlotState struct {
	Plot    *Plot      `json:"plot"`
	Beds    []Bed      `json:"beds"`
	Objects []UIObject `json:"objects"`
	Grid    []GridCell `json:"grid"`
}

type TimelineEvent struct {
	EventID       string    `json:"event_id" db:"event_id"`
	EventType     int       `json:"event_type" db:"event_type"`
	OccurredAt    time.Time `json:"occurred_at" db:"occurred_at"`
	DisplayTitle  string    `json:"display_title" db:"display_title"`
	DisplayDesc   string    `json:"display_description,omitempty" db:"display_description"`
	AffectedCells []string  `json:"affected_cells,omitempty" db:"affected_cells"`
}
