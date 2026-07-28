package service

import (
	"context"
	"encoding/json"
	"garden-nook/internal/modules/plot/dto"
	"garden-nook/internal/modules/plot/enum"
	"garden-nook/internal/modules/plot/model"
	"garden-nook/internal/modules/plot/payload"
	"garden-nook/internal/modules/plot/repository"
	"garden-nook/internal/pkg/database"
	"garden-nook/internal/pkg/helpers"
	"garden-nook/internal/pkg/response"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PlotService struct {
	plotRepo   *repository.PlotRepo
	gridRepo   *repository.GridCellRepo
	bedRepo    *repository.BedRepo
	objectRepo *repository.ObjectRepo
	eventRepo  *repository.EventStoreRepo
	seh        *helpers.ServiceErrorHandler
	pool       *pgxpool.Pool
}

func NewPlotService(
	pool *pgxpool.Pool,
	plotRepo *repository.PlotRepo,
	gridRepo *repository.GridCellRepo,
	bedRepo *repository.BedRepo,
	objectRepo *repository.ObjectRepo,
	eventRepo *repository.EventStoreRepo,
	seh *helpers.ServiceErrorHandler,
) *PlotService {
	return &PlotService{
		pool:       pool,
		plotRepo:   plotRepo,
		gridRepo:   gridRepo,
		bedRepo:    bedRepo,
		objectRepo: objectRepo,
		eventRepo:  eventRepo,
		seh:        seh,
	}
}

func (s *PlotService) ListPlots(ctx context.Context, ownerID string, p *database.Pagination) ([]model.Plot, int, error) {
	return s.plotRepo.ListPlots(ctx, ownerID, p)
}

func (s *PlotService) GetPlotStructure(ctx context.Context, plotID, ownerID string) (*model.PlotStructure, error) {
	plot, err := s.plotRepo.GetPlotByOwnerAndID(ctx, plotID, ownerID)
	if err != nil {
		return nil, s.seh.HandleError(err, "get plot")
	}

	beds, err := s.bedRepo.GetBedsByPlot(ctx, plotID)
	if err != nil {
		return nil, s.seh.HandleError(err, "get beds")
	}

	objects, err := s.objectRepo.GetObjectsByPlot(ctx, plotID)
	if err != nil {
		return nil, s.seh.HandleError(err, "get objects")
	}

	cells, err := s.gridRepo.GetShadedCells(ctx, plotID)
	if err != nil {
		return nil, s.seh.HandleError(err, "get shaded cells")
	}
	shadeGroups := groupShadeCells(cells)

	return &model.PlotStructure{
		Plot:        *plot,
		Beds:        beds,
		Objects:     objects,
		ShadeGroups: shadeGroups,
	}, nil
}

func groupShadeCells(cells []model.GridCell) []model.ShadeGroup {
	groupsMap := make(map[enum.ShadeLevel][]model.CellCoord)
	for _, c := range cells {
		groupsMap[c.ShadeLevel] = append(groupsMap[c.ShadeLevel], model.CellCoord{X: c.XIndex, Y: c.YIndex})
	}
	var groups []model.ShadeGroup
	for level, coords := range groupsMap {
		groups = append(groups, model.ShadeGroup{ShadeLevel: level, Cells: coords})
	}
	return groups
}

func (s *PlotService) CreatePlot(ctx context.Context, ownerID string, req dto.CreatePlotRequest) (*response.CreateUpdateUuidId, error) {
	cellSize := model.DefaultGridCellSize
	cols := int(math.Ceil(req.WidthMeters / cellSize))
	rows := int(math.Ceil(req.HeightMeters / cellSize))

	plot := &model.CreatePlotModel{
		Name:           req.Name,
		SoilTypeID:     req.SoilTypeID,
		BoundaryWidth:  float64(cols) * cellSize,
		BoundaryHeight: float64(rows) * cellSize,
		GridCellSize:   model.DefaultGridCellSize,
		GridCols:       cols,
		GridRows:       rows,
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, s.seh.HandleError(err, "begin tx")
	}
	defer tx.Rollback(ctx)

	plotID, err := s.plotRepo.WithTx(tx).CreatePlot(ctx, plot, ownerID)
	if err != nil {
		return nil, s.seh.HandleError(err, "create plot")
	}

	if err = s.gridRepo.WithTx(tx).CreateCells(ctx, plotID, cols, rows); err != nil {
		return nil, s.seh.HandleError(err, "create cells")
	}

	eventPayload := payload.PlotCreated{
		Name:         req.Name,
		SoilTypeID:   req.SoilTypeID,
		WidthMeters:  req.WidthMeters,
		HeightMeters: req.HeightMeters,
		GridCellSize: cellSize,
		GridCols:     cols,
		GridRows:     rows,
	}
	rawPayload, err := json.Marshal(eventPayload)
	if err != nil {
		return nil, s.seh.HandleError(err, "marshal event")
	}
	if _, err = s.eventRepo.WithTx(tx).AppendEvent(ctx, plotID, enum.EventTypePlotCreated, rawPayload); err != nil {
		return nil, s.seh.HandleError(err, "append event")
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, s.seh.HandleError(err, "commit")
	}

	return &response.CreateUpdateUuidId{Id: plotID}, nil
}

func (s *PlotService) UpdatePlot(ctx context.Context, plotID, ownerID string, req dto.UpdatePlotRequest) (*response.CreateUpdateUuidId, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, s.seh.HandleError(err, "begin tx update plot")
	}
	defer tx.Rollback(ctx)

	_, err = s.plotRepo.WithTx(tx).UpdatePlot(ctx, plotID, ownerID, req)
	if err != nil {
		return nil, s.seh.HandleError(err, "update plot")
	}

	eventPayload := payload.PlotUpdated{
		Name:       req.Name,
		SoilTypeID: req.SoilTypeID,
	}
	rawPayload, err := json.Marshal(eventPayload)
	if err != nil {
		return nil, s.seh.HandleError(err, "marshal event")
	}

	if _, err = s.eventRepo.WithTx(tx).AppendEvent(ctx, plotID, enum.EventTypePlotUpdated, rawPayload); err != nil {
		return nil, s.seh.HandleError(err, "append event")
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, s.seh.HandleError(err, "commit")
	}

	return &response.CreateUpdateUuidId{Id: plotID}, nil
}

func (s *PlotService) DeletePlot(ctx context.Context, plotID, ownerID string) error {
	err := s.plotRepo.DeletePlot(ctx, plotID, ownerID)
	if err != nil {
		return s.seh.HandleError(err, "delete plot")
	}
	return nil
}
