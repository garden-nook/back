package dto

type ObjectCreatedRequest struct {
	Name       string `json:"name"`
	ObjectType int32  `json:"object_type"`
	XStart     int    `json:"x_start"`
	YStart     int    `json:"y_start"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
}

type ObjectUpdatedRequest struct {
	ObjectID   string  `json:"object_id"`
	Name       *string `json:"name,omitempty"`
	ObjectType *int32  `json:"object_type,omitempty"`
	XStart     *int    `json:"x_start,omitempty"`
	YStart     *int    `json:"y_start,omitempty"`
	Width      *int    `json:"width,omitempty"`
	Height     *int    `json:"height,omitempty"`
}

type ObjectDeletedRequest struct {
	ObjectID string `json:"object_id"`
}
