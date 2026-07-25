package model

import "garden-nook/internal/modules/crops/enum"

type CropRule struct {
	ID              int32                `json:"id" db:"rule_id"`
	SubjectCropID   *int32               `json:"subject_crop_id,omitempty" db:"subject_crop_id"`
	SubjectFamilyID *int32               `json:"subject_family_id,omitempty" db:"subject_family_id"`
	ContextType     enum.RuleContextType `json:"context_type" db:"context_type"`
	ContextCropID   *int32               `json:"context_crop_id,omitempty" db:"context_crop_id"`
	ContextFamilyID *int32               `json:"context_family_id,omitempty" db:"context_family_id"`
	ReturnAfterDays int32                `json:"return_after_days" db:"return_after_days"`
	ScoreModifier   int32                `json:"score_modifier" db:"score_modifier"`
	Explanation     string               `json:"explanation" db:"explanation"`
	Priority        int32                `json:"priority" db:"priority"`
}

type CropRelation struct {
	CropID   int32  `json:"crop_id"`
	CropName string `json:"crop_name"`
	Score    int32  `json:"-"`
}

type FamilyRelation struct {
	FamilyID   int32  `json:"family_id"`
	FamilyName string `json:"family_name"`
	Score      int32  `json:"-"`
}

type CropRelations struct {
	GoodPredecessors []CropRelation `json:"good_predecessors"`
	BadPredecessors  []CropRelation `json:"bad_predecessors"`
	GoodSuccessors   []CropRelation `json:"good_successors"`
	BadSuccessors    []CropRelation `json:"bad_successors"`
	GoodCompanions   []CropRelation `json:"good_companions"`
	BadCompanions    []CropRelation `json:"bad_companions"`

	GoodPredecessorFamilies []FamilyRelation `json:"good_predecessor_families"`
	BadPredecessorFamilies  []FamilyRelation `json:"bad_predecessor_families"`
	GoodSuccessorFamilies   []FamilyRelation `json:"good_successor_families"`
	BadSuccessorFamilies    []FamilyRelation `json:"bad_successor_families"`
	GoodCompanionFamilies   []FamilyRelation `json:"good_companion_families"`
	BadCompanionFamilies    []FamilyRelation `json:"bad_companion_families"`
}

type RuleInfo struct {
	SubjectCropID   *int32
	SubjectFamilyID *int32
	ContextCropID   *int32
	ContextFamilyID *int32
	ReturnAfterDays int32
	ScoreModifier   int32
	Explanation     string
}
