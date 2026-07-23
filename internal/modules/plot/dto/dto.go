package dto

import "time"

// ---------- BEDS ----------

type CreateBedRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	Geom string `json:"geom" validate:"required"` // GeoJSON Polygon
}

type UpdateBedRequest struct {
	Name *string `json:"name" validate:"omitempty,min=1,max=100"`
	Geom *string `json:"geom" validate:"omitempty"`
}

type MoveBedRequest struct {
	NewGeom string `json:"new_geom" validate:"required"` // GeoJSON Polygon
}

// ---------- OBJECTS ----------

type CreateObjectRequest struct {
	Name       string `json:"name" validate:"required,min=1,max=100"`
	ObjectType int    `json:"object_type" validate:"required,oneof=1 2 3 4"` // 1=tree, 2=building, 3=path, 4=shade
	Geom       string `json:"geom" validate:"required"`
}

type UpdateObjectRequest struct {
	Name *string `json:"name" validate:"omitempty,min=1,max=100"`
	Geom *string `json:"geom" validate:"omitempty"`
}

// ---------- PLANTINGS ----------

type PlantCropRequest struct {
	BedID     string `json:"bed_id" validate:"required,uuid"`
	CropID    int    `json:"crop_id" validate:"required,gte=1"`
	PlantDate string `json:"plant_date" validate:"required,datetime=2006-01-02"`
}

type HarvestCropRequest struct {
	BedID       string   `json:"bed_id" validate:"required,uuid"`
	HarvestDate string   `json:"harvest_date" validate:"required,datetime=2006-01-02"`
	YieldKg     *float64 `json:"yield_kg" validate:"omitempty,gte=0"`
}

// ---------- HISTORY & TIMELINE ----------

type TimelineFilter struct {
	From  *time.Time `json:"from"`
	To    *time.Time `json:"to"`
	Limit int        `json:"limit"`
}

type PlotStateAtDateRequest struct {
	Date string `json:"date" validate:"required,datetime=2006-01-02"` // YYYY-MM-DD
}
