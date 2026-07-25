package provider

import (
	"context"
	cropModel "garden-nook/internal/modules/crops/model"
	"garden-nook/internal/modules/plot/model"
)

type CropProvider interface {
	ListCropsFiltered(ctx context.Context, filter model.CropFilter) ([]model.CropInfo, error)
}

type RuleProvider interface {
	GetRulesBySubjectCropID(ctx context.Context, cropID int32) ([]cropModel.RuleInfo, error)
	GetRulesBySubjectFamilyID(ctx context.Context, familyID int32) ([]cropModel.RuleInfo, error)
}
