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

type CropFamilyService struct {
	repo *repository.CropFamilyRepo
	seh  *helpers.ServiceErrorHandler
}

func NewCropFamilyService(
	repo *repository.CropFamilyRepo,
	seh *helpers.ServiceErrorHandler,
) *CropFamilyService {
	return &CropFamilyService{repo: repo, seh: seh}
}

func (s *CropFamilyService) ListFamilies(ctx context.Context) ([]model.CropFamily, int, error) {
	return s.repo.ListFamilies(ctx, nil)
}

func (s *CropFamilyService) GetFamily(ctx context.Context, id int32) (*model.CropFamily, error) {
	return s.repo.GetFamilyByID(ctx, id)
}

func (s *CropFamilyService) CreateFamily(ctx context.Context, req dto.CreateFamilyRequest) (*response.CreateUpdateIntId, error) {
	f, err := s.repo.CreateFamily(ctx, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrConflict) {
			return nil, err
		}
		return nil, s.seh.HandleError(err, "update family")
	}
	return &response.CreateUpdateIntId{Id: f}, nil
}

func (s *CropFamilyService) UpdateFamily(ctx context.Context, id int32, req dto.UpdateFamilyRequest) (*response.CreateUpdateIntId, error) {
	f, err := s.repo.UpdateFamily(ctx, id, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) || errors.Is(err, apperrors.ErrConflict) {
			return nil, err
		}
		return nil, s.seh.HandleError(err, "update family")
	}
	return &response.CreateUpdateIntId{Id: f}, nil
}

func (s *CropFamilyService) DeleteFamily(ctx context.Context, id int32) error {
	err := s.repo.DeleteFamily(ctx, id)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) && !errors.Is(err, apperrors.ErrConflict) {
		return s.seh.HandleError(err, "delete family")
	}
	return err
}
