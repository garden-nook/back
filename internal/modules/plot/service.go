package plots

import (
	"context"
	"errors"
	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/response"
	"log/slog"
	"math"
)

type Service struct {
	repo *Repository
	//projector *Projector
	log *slog.Logger
}

func NewService(repo *Repository /*projector *Projector,*/, log *slog.Logger) *Service {
	return &Service{repo: repo /*projector: projector,*/, log: log}
}

// =====================================================================
// PLOTS
// =====================================================================

func (s *Service) ListPlots(ctx context.Context, ownerID string) ([]Plot, int, error) {
	return s.repo.ListPlots(ctx, ownerID, nil)
}

func (s *Service) CreatePlot(ctx context.Context, ownerID string, req CreatePlotRequest) (*response.CreateUpdateUuidId, error) {
	cellSize := DefaultGridCellSize
	cols := int(math.Ceil(req.WidthMeters / cellSize))
	rows := int(math.Ceil(req.HeightMeters / cellSize))
	area := float64(cols*rows) * cellSize * cellSize

	plot := &CreatePlotModel{
		Name:           req.Name,
		SoilTypeID:     req.SoilTypeID,
		BoundaryWidth:  float64(cols) * cellSize,
		BoundaryHeight: float64(rows) * cellSize,
		AreaSqM:        area,
		GridCellSize:   DefaultGridCellSize,
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

	//// Генерируем событие
	//payload := PlotCreatedPayload{
	//	PlotID:       plotID,
	//	OwnerID:      ownerID,
	//	Name:         req.Name,
	//	Description:  req.Description,
	//	Boundary:     req.Boundary,
	//	GridCellSize: req.GridCellSize,
	//	GridCols:     cols,
	//	GridRows:     rows,
	//}
	//if err = s.appendAndProject(ctx, plotID, EventPlotCreated, payload); err != nil {
	//	s.log.Error("append event failed", "err", err)
	//	return nil, apperrors.ErrInternal
	//}

	return &response.CreateUpdateUuidId{Id: id}, nil
}

func (s *Service) UpdatePlot(ctx context.Context, id, ownerID string, req UpdatePlotRequest) (*response.CreateUpdateUuidId, error) {
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
	// Проверяем владельца
	err := s.repo.SoftDeletePlot(ctx, plotID, ownerID)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) && !errors.Is(err, apperrors.ErrConflict) {
		s.log.Error("delete plot failed", "err", err, "id", plotID)
		return apperrors.ErrInternal
	}
	return nil

	//payload := map[string]string{"plot_id": plotID}
	//return s.appendAndProject(ctx, plotID, EventPlotDeleted, payload)
}

//// =====================================================================
//// BEDS
//// =====================================================================
//
//func (s *Service) CreateBed(ctx context.Context, plotID, ownerID string, req CreateBedRequest) (*Bed, error) {
//	// Проверяем, что plot принадлежит user
//	plot, err := s.repo.GetPlotByID(ctx, plotID, ownerID)
//	if err != nil {
//		return nil, err
//	}
//	if plot.OwnerID != ownerID {
//		return nil, apperrors.ErrForbidden
//	}
//
//	bedID := uuid.New().String()
//	bed := &Bed{
//		BedID:  bedID,
//		PlotID: plotID,
//		Name:   req.Name,
//		Geom:   req.Geom,
//	}
//
//	payload := BedCreatedPayload{
//		BedID:  bedID,
//		PlotID: plotID,
//		Name:   req.Name,
//		Geom:   req.Geom,
//	}
//
//	if err := s.appendAndProject(ctx, plotID, EventBedCreated, payload); err != nil {
//		s.log.Error("create bed failed", "err", err)
//		return nil, apperrors.ErrInternal
//	}
//
//	// Привязываем ячейки к грядке
//	cellIDs, err := s.repo.GetCellIDsInPolygon(ctx, plotID, req.Geom)
//	if err != nil {
//		s.log.Error("get cells failed", "err", err)
//		return bed, nil // Bed создан, но ячейки не привязаны
//	}
//	if err := s.repo.UpdateCellsForBed(ctx, plotID, cellIDs, bedID); err != nil {
//		s.log.Error("update cells failed", "err", err)
//	}
//
//	return bed, nil
//}
//
//func (s *Service) UpdateBed(ctx context.Context, plotID, bedID, ownerID string, req UpdateBedRequest) (*Bed, error) {
//	plot, err := s.repo.GetPlotByID(ctx, plotID, ownerID)
//	if err != nil {
//		return nil, err
//	}
//	if plot.OwnerID != ownerID {
//		return nil, apperrors.ErrForbidden
//	}
//
//	payload := BedUpdatedPayload{
//		BedID:   bedID,
//		PlotID:  plotID,
//		Name:    req.Name,
//		NewGeom: req.Geom,
//	}
//
//	if err := s.appendAndProject(ctx, plotID, EventBedUpdated, payload); err != nil {
//		return nil, apperrors.ErrInternal
//	}
//
//	return s.repo.GetBedByID(ctx, bedID, plotID)
//}
//
//func (s *Service) DeleteBed(ctx context.Context, plotID, bedID, ownerID string) error {
//	plot, err := s.repo.GetPlotByID(ctx, plotID, ownerID)
//	if err != nil {
//		return err
//	}
//	if plot.OwnerID != ownerID {
//		return apperrors.ErrForbidden
//	}
//
//	// Очищаем ячейки
//	if err := s.repo.ClearCellsForBed(ctx, plotID, bedID); err != nil {
//		s.log.Error("clear cells failed", "err", err)
//	}
//
//	payload := BedDeletedPayload{BedID: bedID, PlotID: plotID}
//	return s.appendAndProject(ctx, plotID, EventBedDeleted, payload)
//}
//
//func (s *Service) ListBeds(ctx context.Context, plotID, ownerID string) ([]Bed, error) {
//	plot, err := s.repo.GetPlotByID(ctx, plotID, ownerID)
//	if err != nil {
//		return nil, err
//	}
//	if plot.OwnerID != ownerID {
//		return nil, apperrors.ErrForbidden
//	}
//	return s.repo.ListBedsByPlot(ctx, plotID)
//}
//
//// =====================================================================
//// OBJECTS
//// =====================================================================
//
//func (s *Service) CreateObject(ctx context.Context, plotID, ownerID string, req CreateObjectRequest) (*UIObject, error) {
//	plot, err := s.repo.GetPlotByID(ctx, plotID, ownerID)
//	if err != nil {
//		return nil, err
//	}
//	if plot.OwnerID != ownerID {
//		return nil, apperrors.ErrForbidden
//	}
//
//	objectID := uuid.New().String()
//	payload := ObjectCreatedPayload{
//		ObjectID:   objectID,
//		PlotID:     plotID,
//		Name:       req.Name,
//		ObjectType: req.ObjectType,
//		Geom:       req.Geom,
//	}
//
//	if err := s.appendAndProject(ctx, plotID, EventObjectCreated, payload); err != nil {
//		return nil, apperrors.ErrInternal
//	}
//
//	return &UIObject{
//		ObjectID:   objectID,
//		PlotID:     plotID,
//		Name:       req.Name,
//		ObjectType: req.ObjectType,
//		Geom:       req.Geom,
//	}, nil
//}
//
//func (s *Service) ListObjects(ctx context.Context, plotID, ownerID string) ([]UIObject, error) {
//	plot, err := s.repo.GetPlotByID(ctx, plotID, ownerID)
//	if err != nil {
//		return nil, err
//	}
//	if plot.OwnerID != ownerID {
//		return nil, apperrors.ErrForbidden
//	}
//	return s.repo.ListObjectsByPlot(ctx, plotID)
//}
//
//func (s *Service) DeleteObject(ctx context.Context, plotID, objectID, ownerID string) error {
//	plot, err := s.repo.GetPlotByID(ctx, plotID, ownerID)
//	if err != nil {
//		return err
//	}
//	if plot.OwnerID != ownerID {
//		return apperrors.ErrForbidden
//	}
//
//	payload := map[string]string{"object_id": objectID, "plot_id": plotID}
//	return s.appendAndProject(ctx, plotID, EventObjectDeleted, payload)
//}
//
//// =====================================================================
//// PLANTINGS
//// =====================================================================
//
//func (s *Service) PlantCrop(ctx context.Context, plotID, ownerID string, req PlantCropRequest) error {
//	plot, err := s.repo.GetPlotByID(ctx, plotID, ownerID)
//	if err != nil {
//		return err
//	}
//	if plot.OwnerID != ownerID {
//		return apperrors.ErrForbidden
//	}
//
//	// Проверяем, что грядка существует
//	bed, err := s.repo.GetBedByID(ctx, req.BedID, plotID)
//	if err != nil {
//		return err
//	}
//
//	// Получаем ячейки грядки
//	cellIDs, err := s.repo.GetCellIDsInPolygon(ctx, plotID, bed.Geom)
//	if err != nil {
//		return err
//	}
//
//	payload := CropPlantedPayload{
//		PlotID:    plotID,
//		BedID:     req.BedID,
//		CropID:    req.CropID,
//		PlantDate: req.PlantDate,
//		CellIDs:   cellIDs,
//	}
//
//	return s.appendAndProject(ctx, plotID, EventCropPlanted, payload)
//}
//
//func (s *Service) HarvestCrop(ctx context.Context, plotID, ownerID string, req HarvestCropRequest) error {
//	plot, err := s.repo.GetPlotByID(ctx, plotID, ownerID)
//	if err != nil {
//		return err
//	}
//	if plot.OwnerID != ownerID {
//		return apperrors.ErrForbidden
//	}
//
//	bed, err := s.repo.GetBedByID(ctx, req.BedID, plotID)
//	if err != nil {
//		return err
//	}
//
//	cellIDs, err := s.repo.GetCellIDsInPolygon(ctx, plotID, bed.Geom)
//	if err != nil {
//		return err
//	}
//	//
//	//harvestDate, err := time.Parse("2006-01-02", req.HarvestDate)
//	//if err != nil {
//	//	return apperrors.ErrBadRequest
//	//}
//
//	payload := CropHarvestedPayload{
//		PlotID:      plotID,
//		BedID:       req.BedID,
//		HarvestDate: req.HarvestDate,
//		YieldKg:     req.YieldKg,
//		CellIDs:     cellIDs,
//	}
//
//	return s.appendAndProject(ctx, plotID, EventCropHarvested, payload)
//}
//
//// =====================================================================
//// TIMELINE & STATE
//// =====================================================================
//
//func (s *Service) GetTimeline(ctx context.Context, plotID, ownerID string, filter TimelineFilter) ([]TimelineEvent, error) {
//	plot, err := s.repo.GetPlotByID(ctx, plotID, ownerID)
//	if err != nil {
//		return nil, err
//	}
//	if plot.OwnerID != ownerID {
//		return nil, apperrors.ErrForbidden
//	}
//	return s.repo.GetTimeline(ctx, plotID, filter)
//}
//
//func (s *Service) GetPlotState(ctx context.Context, plotID, ownerID string) (*PlotState, error) {
//	plot, err := s.repo.GetPlotByID(ctx, plotID, ownerID)
//	if err != nil {
//		return nil, err
//	}
//	if plot.OwnerID != ownerID {
//		return nil, apperrors.ErrForbidden
//	}
//
//	beds, err := s.repo.ListBedsByPlot(ctx, plotID)
//	if err != nil {
//		return nil, err
//	}
//
//	objects, err := s.repo.ListObjectsByPlot(ctx, plotID)
//	if err != nil {
//		return nil, err
//	}
//
//	grid, err := s.repo.GetGridCellsByPlot(ctx, plotID)
//	if err != nil {
//		return nil, err
//	}
//
//	return &PlotState{
//		Plot:    plot,
//		Beds:    beds,
//		Objects: objects,
//		Grid:    grid,
//	}, nil
//}
//
//// =====================================================================
//// HELPERS
//// =====================================================================
//
//// appendAndProject записывает событие и применяет его к read model.
//// В production лучше разделить: сначала append в event_store, потом async projection.
//func (s *Service) appendAndProject(ctx context.Context, plotID string, eventType int, payload interface{}) error {
//	// 1. Append event
//	if err := s.repo.AppendEvent(ctx, plotID, eventType, payload); err != nil {
//		return err
//	}
//
//	// 2. Project (синхронно для простоты)
//	payloadJSON, _ := json.Marshal(payload)
//	return s.projector.ApplyEvent(ctx, plotID, eventType, payloadJSON)
//}
