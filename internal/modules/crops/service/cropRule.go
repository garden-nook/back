package service

import (
	"context"
	"errors"
	"garden-nook/internal/modules/crops/dto"
	"garden-nook/internal/modules/crops/model"
	"garden-nook/internal/modules/crops/repository"
	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/helpers"
	"garden-nook/internal/pkg/response"
)

type CropRuleService struct {
	repo      *repository.CropRuleRepo
	ruleCache *RuleCache
	seh       *helpers.ServiceErrorHandler
}

func NewCropRuleService(
	repo *repository.CropRuleRepo,
	ruleCache *RuleCache,
	seh *helpers.ServiceErrorHandler,
) *CropRuleService {
	return &CropRuleService{
		repo:      repo,
		ruleCache: ruleCache,
		seh:       seh,
	}
}

func (s *CropRuleService) ListRules(ctx context.Context) ([]model.CropRule, int, error) {
	return s.repo.ListRules(ctx, nil)
}

func (s *CropRuleService) CreateRule(ctx context.Context, req dto.CreateRuleRequest) (*response.CreateUpdateIntId, error) {
	// Бизнес-правило: должен быть указан либо crop, либо family (но не оба одновременно на одной стороне)
	if req.SubjectCropID != nil && req.SubjectFamilyID != nil {
		return nil, apperrors.ErrBadRequest
	}
	if req.ContextCropID != nil && req.ContextFamilyID != nil {
		return nil, apperrors.ErrBadRequest
	}
	if req.SubjectCropID == nil && req.SubjectFamilyID == nil {
		return nil, apperrors.ErrBadRequest
	}
	if req.ContextCropID == nil && req.ContextFamilyID == nil {
		return nil, apperrors.ErrBadRequest
	}

	r, err := s.repo.CreateRule(ctx, req)
	if err != nil {
		return nil, s.seh.HandleError(err, "create rule")
	}

	err = s.ruleCache.Refresh(ctx)
	if err != nil {
		return nil, s.seh.HandleError(err, "refresh rule")
	}
	return &response.CreateUpdateIntId{Id: r}, nil
}

func (s *CropRuleService) DeleteRule(ctx context.Context, id int32) error {
	err := s.repo.DeleteRule(ctx, id)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return s.seh.HandleError(err, "delete rule")
	}

	err = s.ruleCache.Refresh(ctx)
	if err != nil {
		return s.seh.HandleError(err, "refresh rule")
	}
	return err
}
