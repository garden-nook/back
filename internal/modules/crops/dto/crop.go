package dto

import (
	"garden-nook/internal/modules/crops/model"
	"garden-nook/internal/modules/plot/enum"
)

type CreateCropRequest struct {
	Name              string          `json:"name" validate:"required,min=2,max=100"`
	FamilyID          int32           `json:"family_id" validate:"required,gte=1"`
	SoilTypeID        int32           `json:"soil_type_id" validate:"required,gte=1"`
	VegetationDaysAvg int32           `json:"vegetation_days_avg" validate:"required,gte=1,lte=1000"`
	SunNeeds          enum.ShadeLevel `json:"sun_needs" validate:"required,oneof=1 2 3"`
}

type UpdateCropRequest struct {
	Name              *string          `json:"name" validate:"omitempty,min=2,max=100"`
	FamilyID          *int32           `json:"family_id" validate:"omitempty,gte=1"`
	SoilTypeID        *int32           `json:"soil_type_id" validate:"required,gte=1"`
	VegetationDaysAvg *int32           `json:"vegetation_days_avg" validate:"omitempty,gte=1,lte=1000"`
	SunNeeds          *enum.ShadeLevel `json:"sun_needs" validate:"omitempty,oneof=1 2 3"`
}

type ListCropsFilter struct {
	FamilyID   *int32 `json:"family_id"`
	SoilTypeID *int32 `json:"soil_type_id"`
	Search     string `json:"search"`
}

type CropExtended struct {
	*model.Crop          `json:"crop"`
	*model.CropRelations `json:"crop_relations"`
}
