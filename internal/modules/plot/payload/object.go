package payload

type ObjectCreated struct {
	ObjectID   string `json:"object_id"`
	Name       string `json:"name"`
	ObjectType int32  `json:"object_type"`
	XStart     int    `json:"x_start"`
	YStart     int    `json:"y_start"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
}

type ObjectUpdated struct {
	ObjectID   string `json:"object_id"`
	Name       string `json:"name"`
	ObjectType int32  `json:"object_type"`
	XStart     int    `json:"x_start"`
	YStart     int    `json:"y_start"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
}

type ObjectDeleted struct {
	ObjectID string `json:"object_id"`
}
