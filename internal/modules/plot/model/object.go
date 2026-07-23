package model

type Object struct {
	ObjectID   string `db:"object_id" json:"object_id"`
	PlotID     string `db:"plot_id" json:"-"`
	Name       string `db:"name" json:"name"`
	ObjectType int32  `db:"object_type" json:"object_type"`
	XStart     int    `db:"x_start" json:"x_start"`
	YStart     int    `db:"y_start" json:"y_start"`
	Width      int    `db:"width" json:"width"`
	Height     int    `db:"height" json:"height"`
}
