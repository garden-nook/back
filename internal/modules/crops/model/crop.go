package model

import (
	"garden-nook/internal/modules/plot/enum"
)

type Crop struct {
	ID                int32           `json:"id" db:"id"`
	Name              string          `json:"name" db:"name"`
	Description       *string         `json:"description,omitempty" db:"description"`
	FamilyID          int32           `json:"family_id" db:"family_id"`
	FamilyName        *string         `json:"family_name,omitempty" db:"family_name"`
	VegetationDaysAvg int32           `json:"vegetation_days_avg" db:"vegetation_days_avg"`
	SunNeeds          enum.ShadeLevel `json:"sun_needs" db:"sun_needs"`
	SoilTypeID        int32           `json:"soil_type_id" db:"soil_type_id"`
	SoilName          *string         `json:"soil_name,omitempty" db:"soil_name"`
}
