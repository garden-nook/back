package model

type SoilType struct {
	ID          int32   `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description,omitempty" db:"description"`
}
