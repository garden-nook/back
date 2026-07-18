package crops

type CreateSoilTypeRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=2000"`
}

type UpdateSoilTypeRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=2,max=100"`
	Description *string `json:"description" validate:"omitempty,max=2000"`
}

type CreateFamilyRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=2000"`
}

type UpdateFamilyRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=2,max=100"`
	Description *string `json:"description" validate:"omitempty,max=2000"`
}

type CreateCropRequest struct {
	Name              string   `json:"name" validate:"required,min=2,max=100"`
	FamilyID          int32    `json:"family_id" validate:"required,gte=1"`
	SoilTypeID        int32    `json:"soil_type_id" validate:"required,gte=1"`
	VegetationDaysAvg int32    `json:"vegetation_days_avg" validate:"required,gte=1,lte=1000"`
	SunNeeds          SunNeeds `json:"sun_needs" validate:"required,oneof=1 2 3"`
}

type UpdateCropRequest struct {
	Name              *string   `json:"name" validate:"omitempty,min=2,max=100"`
	FamilyID          *int32    `json:"family_id" validate:"omitempty,gte=1"`
	SoilTypeID        *int32    `json:"soil_type_id" validate:"required,gte=1"`
	VegetationDaysAvg *int32    `json:"vegetation_days_avg" validate:"omitempty,gte=1,lte=1000"`
	SunNeeds          *SunNeeds `json:"sun_needs" validate:"omitempty,oneof=1 2 3"`
}

type ListCropsFilter struct {
	FamilyID   *int32 `json:"family_id"`
	SoilTypeID *int32 `json:"soil_type_id"`
	Search     string `json:"search"`
}

type CreateRuleRequest struct {
	SubjectCropID   *int32      `json:"subject_crop_id"`
	SubjectFamilyID *int32      `json:"subject_family_id"`
	ContextType     ContextType `json:"context_type" validate:"required,oneof=1 2 3"`
	ContextCropID   *int32      `json:"context_crop_id"`
	ContextFamilyID *int32      `json:"context_family_id"`
	ReturnAfterDays int32       `json:"return_after_days" validate:"gte=0,lte=3650"`
	ScoreModifier   int32       `json:"score_modifier" validate:"gte=-100,lte=100"`
	Explanation     string      `json:"explanation" validate:"required,min=5,max=1000"`
	Priority        int32       `json:"priority" validate:"gte=1,lte=100"`
}
