package crops

import "time"

// CropFamily — доменная модель семейства культур.
type CropFamily struct {
	ID          int32  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description,omitempty" db:"description"`
}

// Crop — доменная модель культуры.
type Crop struct {
	ID                int32     `json:"id" db:"id"`
	Name              string    `json:"name" db:"name"`
	FamilyID          int32     `json:"family_id" db:"family_id"`
	FamilyName        string    `json:"family_name,omitempty" db:"family_name"`
	VegetationDaysAvg int32     `json:"vegetation_days_avg" db:"vegetation_days_avg"`
	SunNeeds          SunNeeds  `json:"sun_needs" db:"sun_needs"`
	IsDeleted         bool      `json:"-" db:"is_deleted"`
	CreatedAt         time.Time `json:"created_at,omitempty" db:"created_at"`
}

// CropRule — правило совместимости (севооборот/аллелопатия).
type CropRule struct {
	ID              int32       `json:"id" db:"rule_id"`
	SubjectCropID   *int32      `json:"subject_crop_id,omitempty" db:"subject_crop_id"`
	SubjectFamilyID *int32      `json:"subject_family_id,omitempty" db:"subject_family_id"`
	ContextType     ContextType `json:"context_type" db:"context_type"`
	ContextCropID   *int32      `json:"context_crop_id,omitempty" db:"context_crop_id"`
	ContextFamilyID *int32      `json:"context_family_id,omitempty" db:"context_family_id"`
	ReturnAfterDays int32       `json:"return_after_days" db:"return_after_days"`
	ScoreModifier   int32       `json:"score_modifier" db:"score_modifier"`
	Explanation     string      `json:"explanation" db:"explanation"`
	Priority        int32       `json:"priority" db:"priority"`
}

// SunNeeds Enum: 1=shade, 2=partial, 3=full
type SunNeeds int32

const (
	SunNeedsShade   SunNeeds = 1
	SunNeedsPartial SunNeeds = 2
	SunNeedsFull    SunNeeds = 3
)

// ContextType Enum: 1=PREDECESSOR, 2=NEIGHBOR, 3=SELF
type ContextType int32

const (
	ContextPredecessor ContextType = 1
	ContextNeighbor    ContextType = 2
	ContextSelf        ContextType = 3
)
