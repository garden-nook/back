package service

import (
	"context"
	"errors"
	"garden-nook/internal/modules/crops/dto"
	"garden-nook/internal/modules/crops/model"
	"garden-nook/internal/modules/crops/repository"
	plotModel "garden-nook/internal/modules/plot/model"
	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/helpers"
	"garden-nook/internal/pkg/response"

	"golang.org/x/sync/errgroup"
)

type CropService struct {
	cropRepo       *repository.CropRepo
	cropFamilyRepo *repository.CropFamilyRepo
	cropRuleRepo   *repository.CropRuleRepo
	seh            *helpers.ServiceErrorHandler
}

func NewCropService(
	cropRepo *repository.CropRepo,
	cropFamilyRepo *repository.CropFamilyRepo,
	cropRuleRepo *repository.CropRuleRepo,
	seh *helpers.ServiceErrorHandler,
) *CropService {
	return &CropService{
		cropRepo:       cropRepo,
		cropFamilyRepo: cropFamilyRepo,
		cropRuleRepo:   cropRuleRepo,
		seh:            seh,
	}
}

func (s *CropService) ListCrops(ctx context.Context, f dto.ListCropsFilter) ([]model.Crop, int, error) {
	return s.cropRepo.ListCrops(ctx, f, nil)
}

func (s *CropService) ListCropsFiltered(ctx context.Context, filter plotModel.CropFilter) ([]plotModel.CropInfo, error) {
	crops, err := s.cropRepo.ListCropsFiltered(ctx, filter.SoilTypeID, filter.SunNeeds, filter.Search, filter.Limit)
	if err != nil {
		return nil, err
	}
	result := make([]plotModel.CropInfo, len(crops))
	for i, c := range crops {
		result[i] = plotModel.CropInfo{
			ID:         c.ID,
			Name:       c.Name,
			FamilyID:   c.FamilyID,
			FamilyName: *c.FamilyName,
			SunNeeds:   c.SunNeeds,
			SoilTypeID: c.SoilTypeID,
		}
	}
	return result, nil
}

func (s *CropService) GetCrop(ctx context.Context, id int32) (*dto.CropExtended, error) {
	var (
		c  *model.Crop
		cr *model.CropRelations
	)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		c, err = s.cropRepo.GetCropByID(ctx, id)
		return err
	})
	g.Go(func() error {
		var err error
		cr, err = s.cropRuleRepo.GetCropRelations(ctx, id)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return &dto.CropExtended{Crop: c, CropRelations: cr}, nil
}

func (s *CropService) CreateCrop(ctx context.Context, req dto.CreateCropRequest) (*response.CreateUpdateIntId, error) {
	if _, err := s.cropFamilyRepo.GetFamilyByID(ctx, req.FamilyID); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrBadRequest
		}
		return nil, err
	}
	c, err := s.cropRepo.CreateCrop(ctx, req)
	if err != nil {
		return nil, s.seh.HandleError(err, "create crop")
	}
	return &response.CreateUpdateIntId{Id: c}, nil
}

func (s *CropService) UpdateCrop(ctx context.Context, id int32, req dto.UpdateCropRequest) (*response.CreateUpdateIntId, error) {
	if req.FamilyID != nil {
		if _, err := s.cropFamilyRepo.GetFamilyByID(ctx, *req.FamilyID); err != nil {
			if errors.Is(err, apperrors.ErrNotFound) {
				return nil, apperrors.ErrBadRequest
			}
			return nil, err
		}
	}
	c, err := s.cropRepo.UpdateCrop(ctx, id, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, err
		}
		return nil, s.seh.HandleError(err, "update crop")
	}
	return &response.CreateUpdateIntId{Id: c}, nil
}

func (s *CropService) DeleteCrop(ctx context.Context, id int32) error {
	err := s.cropRepo.SoftDeleteCrop(ctx, id)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return s.seh.HandleError(err, "delete crop")
	}
	return err
}
