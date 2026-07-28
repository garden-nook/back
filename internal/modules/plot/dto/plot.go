package dto

type CreatePlotRequest struct {
	Name         string
	SoilTypeID   int32
	WidthMeters  float64 // Ширина в метрах
	HeightMeters float64 // Высота в метрах
}

type UpdatePlotRequest struct {
	Name       *string
	SoilTypeID *int32
}
