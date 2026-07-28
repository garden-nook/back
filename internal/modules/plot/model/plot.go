package model

const DefaultGridCellSize = 0.5

type Plot struct {
	PlotID       string  `json:"plot_id" db:"plot_id"`
	Name         string  `json:"name" db:"name"`
	SoilTypeID   int32   `json:"soil_type" db:"soil_type"`
	SoilTypeName string  `json:"soil_name" db:"soil_name"`
	GridCellSize float64 `json:"grid_cell_size" db:"grid_cell_size"`
	GridCols     int     `json:"grid_cols" db:"grid_cols"`
	GridRows     int     `json:"grid_rows" db:"grid_rows"`
	BedCount     int     `json:"bed_count" db:"bed_count"`
	CropCount    int     `json:"crop_count" db:"crop_count"`
}

type PlotStructure struct {
	Plot        Plot         `json:"plot"`
	Beds        []Bed        `json:"beds"`
	Objects     []Object     `json:"objects"`
	ShadeGroups []ShadeGroup `json:"shade_groups"`
}

type CreatePlotModel struct {
	Name           string
	SoilTypeID     int32
	BoundaryWidth  float64
	BoundaryHeight float64
	GridCellSize   float64
	GridCols       int
	GridRows       int
}
