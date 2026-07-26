package service

import (
	"context"
	"garden-nook/internal/modules/plot/dto"
	"garden-nook/internal/modules/plot/repository"
	"garden-nook/internal/pkg/helpers"
)

type HistoryService struct {
	bedRepo  *repository.BedRepo
	plotRepo *repository.PlotRepo
	histRepo *repository.HistoryRepo
	seh      *helpers.ServiceErrorHandler
}

func NewHistoryService(
	bedRepo *repository.BedRepo,
	plotRepo *repository.PlotRepo,
	histRepo *repository.HistoryRepo,
	seh *helpers.ServiceErrorHandler,
) *HistoryService {
	return &HistoryService{
		bedRepo:  bedRepo,
		plotRepo: plotRepo,
		histRepo: histRepo,
		seh:      seh,
	}
}

func (s *HistoryService) GetBedCropHistory(ctx context.Context, bedID, ownerID string) ([]dto.BedCropHistoryEntry, error) {
	bed, err := s.bedRepo.GetBedByID(ctx, bedID)
	if err != nil {
		return nil, s.seh.HandleError(err, "get bed")
	}

	if _, err = s.plotRepo.GetPlotByOwnerAndID(ctx, bed.PlotID, ownerID); err != nil {
		return nil, s.seh.HandleError(err, "check plot ownership")
	}

	entries, err := s.histRepo.GetCropHistoryForBed(ctx, bed.PlotID, bed.XStart, bed.YStart, bed.Width, bed.Height)
	if err != nil {
		return nil, s.seh.HandleError(err, "get history for bed")
	}
	return entries, nil
}
