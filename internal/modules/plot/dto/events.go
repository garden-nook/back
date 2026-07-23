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

type CellShadeUpdatedRequest struct {
	ShadeGroups []model.ShadeGroup `json:"shade_groups"`
}

type PlotResizedRequest struct {
	NewWidthMeters  float64 `json:"new_width_meters"`
	NewHeightMeters float64 `json:"new_height_meters"`
}
