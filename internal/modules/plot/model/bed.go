package model

import "time"

type Bed struct {
	BedID         string     `db:"bed_id" json:"bed_id"`
	PlotID        string     `db:"plot_id" json:"-"`
	Name          string     `db:"name" json:"name"`
	XStart        int        `db:"x_start" json:"x_start"`
	YStart        int        `db:"y_start" json:"y_start"`
	Width         int        `db:"width" json:"width"`
	Height        int        `db:"height" json:"height"`
	CurrentCropID *int32     `db:"current_crop_id" json:"current_crop_id"`
	PlantDate     *time.Time `db:"plant_date" json:"plant_date"`
}
