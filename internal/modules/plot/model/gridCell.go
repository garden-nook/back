package model

import "garden-nook/internal/modules/plot/enum"

type GridCell struct {
	PlotID     string          `db:"plot_id" json:"-"`
	XIndex     int             `db:"x_index" json:"x_index"`
	YIndex     int             `db:"y_index" json:"y_index"`
	ShadeLevel enum.ShadeLevel `db:"shade_level" json:"shade_level"`
}

type ShadeGroup struct {
	ShadeLevel enum.ShadeLevel `json:"shade_level"`
	Cells      []CellCoord     `json:"cells"`
}

type CellCoord struct {
	X int `json:"x"`
	Y int `json:"y"`
}
