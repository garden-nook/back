package services

import (
	"context"
	"errors"
	"garden-nook/internal/modules/plot/models"
	"garden-nook/internal/modules/plot/repositories"
	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/response"
	"log/slog"
	"math"
)

type Service struct {
	repo *repositories.Repository
	log  *slog.Logger
}

func NewService(repo *repositories.Repository, log *slog.Logger) *Service {
	return &Service{repo: repo, log: log}
}

// =====================================================================
// PLOTS
// =====================================================================

func (s *Service) ListPlots(ctx context.Context, ownerID string) ([]models.Plot, int, error) {
	return s.repo.ListPlots(ctx, ownerID, nil)
}

func (s *Service) CreatePlot(ctx context.Context, ownerID string, req models.CreatePlotRequest) (*response.CreateUpdateUuidId, error) {
	cellSize := models.DefaultGridCellSize
	cols := int(math.Ceil(req.WidthMeters / cellSize))
	rows := int(math.Ceil(req.HeightMeters / cellSize))
	area := float64(cols*rows) * cellSize * cellSize

	plot := &models.CreatePlotModel{
		Name:           req.Name,
		SoilTypeID:     req.SoilTypeID,
		BoundaryWidth:  float64(cols) * cellSize,
		BoundaryHeight: float64(rows) * cellSize,
		AreaSqM:        area,
		GridCellSize:   models.DefaultGridCellSize,
		GridCols:       cols,
		GridRows:       rows,
	}

	id, err := s.repo.CreatePlot(ctx, plot, ownerID)
	if err != nil {
		if errors.Is(err, apperrors.ErrConflict) {
			return nil, err
		}
		s.log.Error("create plot failed", "err", err)
		return nil, apperrors.ErrInternal
	}

	return &response.CreateUpdateUuidId{Id: id}, nil
}

func (s *Service) UpdatePlot(ctx context.Context, id, ownerID string, req models.UpdatePlotRequest) (*response.CreateUpdateUuidId, error) {
	id, err := s.repo.UpdatePlot(ctx, id, ownerID, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) || errors.Is(err, apperrors.ErrConflict) {
			return nil, err
		}
		s.log.Error("update plot failed", "err", err, "id", id)
		return nil, apperrors.ErrInternal
	}
	return &response.CreateUpdateUuidId{Id: id}, nil
}

func (s *Service) DeletePlot(ctx context.Context, plotID, ownerID string) error {
	err := s.repo.SoftDeletePlot(ctx, plotID, ownerID)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) && !errors.Is(err, apperrors.ErrConflict) {
		s.log.Error("delete plot failed", "err", err, "id", plotID)
		return apperrors.ErrInternal
	}
	return nil
}
