package payload

type PlotCreated struct {
	Name         string  `json:"name"`
	SoilTypeID   int32   `json:"soil_type_id"`
	WidthMeters  float64 `json:"width_meters"`
	HeightMeters float64 `json:"height_meters"`
	GridCellSize float64 `json:"grid_cell_size"`
	GridCols     int     `json:"grid_cols"`
	GridRows     int     `json:"grid_rows"`
}

type PlotUpdated struct {
	Name       *string `json:"name,omitempty"`
	SoilTypeID *int32  `json:"soil_type_id,omitempty"`
}
