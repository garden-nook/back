package payload

type BedCreated struct {
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
