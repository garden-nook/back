package dto

import "time"

type AddHistoryRequest struct {
	PlotID    string
	XIndex    int
	YIndex    int
	CropID    int32
	FamilyID  int32
	BedID     *string
	PlantDate time.Time
}

type HarvestRequest struct {
	PlotID      string
	XIndex      int
	YIndex      int
	HarvestDate time.Time
	Metadata    *string
}

type HistoryFilter struct {
	BedID    *string
	CropID   *int32
	FamilyID *int32
	FromDate *time.Time
	ToDate   *time.Time
}
