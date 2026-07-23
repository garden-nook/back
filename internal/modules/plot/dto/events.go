package dto

import (
	"encoding/json"
	"garden-nook/internal/modules/plot/enum"
	"garden-nook/internal/modules/plot/model"
)

type PlotEvents struct {
	Events []Event `json:"events"`
}

type Event struct {
	Type    enum.EventType  `json:"type"`
	Payload json.RawMessage `json:"payload" swaggertype:"object"`
}

type BedCreatedRequest struct {
	Name   string `json:"name"`
	XStart int    `json:"x_start"`
	YStart int    `json:"y_start"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type BedUpdatedRequest struct {
	Name   *string `json:"name,omitempty"`
	XStart *int    `json:"x_start,omitempty"`
	YStart *int    `json:"y_start,omitempty"`
	Width  *int    `json:"width,omitempty"`
	Height *int    `json:"height,omitempty"`
}

type BedDeletedRequest struct {
	BedID string `json:"bed_id"`
}

type ObjectCreatedRequest struct {
	Name       string `json:"name"`
	ObjectType int32  `json:"object_type"`
	XStart     int    `json:"x_start"`
	YStart     int    `json:"y_start"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
}

type ObjectUpdatedRequest struct {
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

type CropPlantedRequest struct {
	BedID     string  `json:"bed_id"`
	CropID    int32   `json:"crop_id"`
	PlantDate *string `json:"plant_date,omitempty"`
}

type CropRemovedRequest struct {
	BedID     string  `json:"bed_id"`
	Harvested bool    `json:"harvested"`
	Date      *string `json:"date,omitempty"`
	Metadata  *string `json:"metadata,omitempty"`
}

type CellShadeUpdatedRequest struct {
	ShadeGroups []model.ShadeGroup `json:"shade_groups"`
}

type PlotResizedRequest struct {
	NewWidthMeters  float64 `json:"new_width_meters"`
	NewHeightMeters float64 `json:"new_height_meters"`
}
