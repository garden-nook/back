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

type SoilTypeService struct {
	repo *repository.SoilTypeRepo
	seh  *helpers.ServiceErrorHandler
}

func NewSoilTypeService(
	repo *repository.SoilTypeRepo,
	seh *helpers.ServiceErrorHandler,
) *SoilTypeService {
	return &SoilTypeService{repo: repo, seh: seh}
}

func (s *SoilTypeService) ListSoilTypes(ctx context.Context) ([]model.SoilType, int, error) {
	return s.repo.ListSoilTypes(ctx, nil)
}

func (s *SoilTypeService) GetSoilType(ctx context.Context, id int32) (*model.SoilType, error) {
	return s.repo.GetSoilTypeByID(ctx, id)
}

func (s *SoilTypeService) CreateSoilType(ctx context.Context, req dto.CreateSoilTypeRequest) (*response.CreateUpdateIntId, error) {
	f, err := s.repo.CreateSoilType(ctx, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrConflict) {
			return nil, err
		}
		return nil, s.seh.HandleError(err, "create soil type")
	}
	return &response.CreateUpdateIntId{Id: f}, nil
}

func (s *SoilTypeService) UpdateSoilType(ctx context.Context, id int32, req dto.UpdateSoilTypeRequest) (*response.CreateUpdateIntId, error) {
	f, err := s.repo.UpdateSoilType(ctx, id, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) || errors.Is(err, apperrors.ErrConflict) {
			return nil, err
		}
		return nil, s.seh.HandleError(err, "update soil type")
	}
	return &response.CreateUpdateIntId{Id: f}, nil
}

func (s *SoilTypeService) DeleteSoilType(ctx context.Context, id int32) error {
	err := s.repo.DeleteSoilType(ctx, id)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) && !errors.Is(err, apperrors.ErrConflict) {
		return s.seh.HandleError(err, "delete soil type")
	}
	return err
}
