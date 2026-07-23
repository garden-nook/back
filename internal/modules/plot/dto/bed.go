package dto

type BedCreatedRequest struct {
	Name   string `json:"name"`
	XStart int    `json:"x_start"`
	YStart int    `json:"y_start"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type BedUpdatedRequest struct {
	BedID  string  `json:"bed_id"`
	Name   *string `json:"name,omitempty"`
	XStart *int    `json:"x_start,omitempty"`
	YStart *int    `json:"y_start,omitempty"`
	Width  *int    `json:"width,omitempty"`
	Height *int    `json:"height,omitempty"`
}

type BedDeletedRequest struct {
	BedID string `json:"bed_id"`
}

type CropPlantedRequest struct {
	BedID     string  `json:"bed_id"`
	CropID    int32   `json:"crop_id"`
	PlantDate *string `json:"plant_date,omitempty"`
}

type CropRemovedRequest struct {
	BedID     string  `json:"bed_id"`
	Harvested bool    `json:"harvested"`
	Date      *string `json:"date,omitempty"`
}
