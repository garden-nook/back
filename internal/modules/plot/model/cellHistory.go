package model

import "time"

type CellHistory struct {
	PlotID      string    `db:"plot_id"`
	XIndex      int       `db:"x_index"`
	YIndex      int       `db:"y_index"`
	CropID      int32     `db:"crop_id"`
	PlantDate   time.Time `db:"plant_date"`
	HarvestDate time.Time `db:"harvest_date"`
}
