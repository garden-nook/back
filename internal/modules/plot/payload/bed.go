package payload

type BedCreated struct {
	BedID  string `json:"bed_id"`
	Name   string `json:"name"`
	XStart int    `json:"x_start"`
	YStart int    `json:"y_start"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type BedUpdated struct {
	BedID  string `json:"bed_id"`
	Name   string `json:"name"`
	XStart int    `json:"x_start"`
	YStart int    `json:"y_start"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type BedDeleted struct {
	BedID string `json:"bed_id"`
}

type CropPlanted struct {
	BedID     string `json:"bed_id"`
	CropID    int32  `json:"crop_id"`
	PlantDate string `json:"plant_date"`
}

type CropRemoved struct {
	BedID     string `json:"bed_id"`
	CropID    int32  `json:"crop_id"`
	Harvested bool   `json:"harvested"`
	Date      string `json:"date"`
}
