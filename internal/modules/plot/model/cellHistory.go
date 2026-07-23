package model

import "time"

type CellHistory struct {
	HistoryID   string     `json:"history_id" db:"history_id"`
	PlotID      string     `json:"plot_id" db:"plot_id"`
	XIndex      int        `json:"x_index" db:"x_index"`
	YIndex      int        `json:"y_index" db:"y_index"`
	CropID      int        `json:"crop_id" db:"crop_id"`
	FamilyID    int        `json:"family_id" db:"family_id"`
	BedID       *string    `json:"bed_id,omitempty" db:"bed_id"`
	PlantDate   time.Time  `json:"plant_date" db:"plant_date"`
	HarvestDate *time.Time `json:"harvest_date,omitempty" db:"harvest_date"`
	Metadata    *string    `json:"metadata" db:"metadata"`
}
