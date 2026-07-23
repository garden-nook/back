package plots

// Типы событий (enum в БД)
const (
	EventPlotCreated   = 1
	EventPlotUpdated   = 2
	EventPlotDeleted   = 3
	EventBedCreated    = 10
	EventBedUpdated    = 11
	EventBedDeleted    = 12
	EventObjectCreated = 20
	EventObjectUpdated = 21
	EventObjectDeleted = 22
	EventCropPlanted   = 30
	EventCropHarvested = 31
)

// Payload'ы событий

type PlotCreatedPayload struct {
	PlotID       string  `json:"plot_id"`
	OwnerID      string  `json:"owner_id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Boundary     string  `json:"boundary"`
	GridCellSize float64 `json:"grid_cell_size"`
	GridCols     int     `json:"grid_cols"`
	GridRows     int     `json:"grid_rows"`
}

type BedCreatedPayload struct {
	BedID  string `json:"bed_id"`
	PlotID string `json:"plot_id"`
	Name   string `json:"name"`
	Geom   string `json:"geom"`
}

type BedUpdatedPayload struct {
	BedID   string  `json:"bed_id"`
	PlotID  string  `json:"plot_id"`
	Name    *string `json:"name,omitempty"`
	NewGeom *string `json:"new_geom,omitempty"`
}

type BedDeletedPayload struct {
	BedID  string `json:"bed_id"`
	PlotID string `json:"plot_id"`
}

type ObjectCreatedPayload struct {
	ObjectID   string `json:"object_id"`
	PlotID     string `json:"plot_id"`
	Name       string `json:"name"`
	ObjectType int    `json:"object_type"`
	Geom       string `json:"geom"`
}

type CropPlantedPayload struct {
	PlotID    string   `json:"plot_id"`
	BedID     string   `json:"bed_id"`
	CropID    int      `json:"crop_id"`
	PlantDate string   `json:"plant_date"`
	CellIDs   []string `json:"cell_ids"` // Список затронутых ячеек
}

type CropHarvestedPayload struct {
	PlotID      string   `json:"plot_id"`
	BedID       string   `json:"bed_id"`
	HarvestDate string   `json:"harvest_date"`
	YieldKg     *float64 `json:"yield_kg,omitempty"`
	CellIDs     []string `json:"cell_ids"`
}
