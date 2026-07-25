package dto

import (
	"garden-nook/internal/modules/crops/enum"
)

type CreateRuleRequest struct {
	SubjectCropID   *int32               `json:"subject_crop_id"`
	SubjectFamilyID *int32               `json:"subject_family_id"`
	ContextType     enum.RuleContextType `json:"context_type" validate:"required,oneof=1 2"`
	ContextCropID   *int32               `json:"context_crop_id"`
	ContextFamilyID *int32               `json:"context_family_id"`
	ReturnAfterDays int32                `json:"return_after_days" validate:"gte=0,lte=3650"`
	ScoreModifier   int32                `json:"score_modifier" validate:"gte=-100,lte=100"`
	Explanation     string               `json:"explanation" validate:"required,min=5,max=1000"`
	Priority        int32                `json:"priority" validate:"gte=1,lte=100"`
}
