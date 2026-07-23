package service

import (
	"context"
	"encoding/json"
	"errors"
	"garden-nook/internal/modules/plot/dto"
	"garden-nook/internal/modules/plot/enum"
	"garden-nook/internal/modules/plot/model"
	"garden-nook/internal/modules/plot/payload"
	"garden-nook/internal/modules/plot/repository"
	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/helpers"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventService struct {
	pool        *pgxpool.Pool
	plotRepo    *repository.PlotRepo
	bedRepo     *repository.BedRepo
	objectRepo  *repository.ObjectRepo
	eventRepo   *repository.EventStoreRepo
	historyRepo *repository.HistoryRepo
	seh         *helpers.ServiceErrorHandler
}

func NewEventService(
	pool *pgxpool.Pool,
	plotRepo *repository.PlotRepo,
	bedRepo *repository.BedRepo,
	objectRepo *repository.ObjectRepo,
	eventRepo *repository.EventStoreRepo,
	historyRepo *repository.HistoryRepo,
	seh *helpers.ServiceErrorHandler,
) *EventService {
	return &EventService{
		pool:        pool,
		plotRepo:    plotRepo,
		bedRepo:     bedRepo,
		objectRepo:  objectRepo,
		eventRepo:   eventRepo,
		historyRepo: historyRepo,
		seh:         seh,
	}
}

func (s *EventService) HandleEvents(ctx context.Context, plotID, ownerID string, events []dto.Event) error {
	if len(events) == 0 {
		return nil // нет событий – нечего делать
	}

	for _, ev := range events {
		if !enum.IsValidEventType(ev.Type) {
			return apperrors.ErrBadRequest
		}
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return s.seh.HandleError(err, "begin tx")
	}
	defer tx.Rollback(ctx)

	if _, err = s.plotRepo.WithTx(tx).GetPlotByOwnerAndID(ctx, plotID, ownerID); err != nil {
		return s.seh.HandleError(err, "check plot ownership")
	}

	for _, ev := range events {
		if err = s.handleEvent(ctx, plotID, ev, tx); err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return s.seh.HandleError(err, "commit events")
	}
	return nil
}

func (s *EventService) handleEvent(ctx context.Context, plotID string, ev dto.Event, tx pgx.Tx) error {
	switch ev.Type {
	case enum.EventTypeBedCreated:
		return s.handleBedCreated(ctx, plotID, ev.Payload, tx)
	case enum.EventTypeBedUpdated:
		return s.handleBedUpdated(ctx, plotID, ev.Payload, tx)
	case enum.EventTypeBedDeleted:
		return s.handleBedDeleted(ctx, plotID, ev.Payload, tx)
	case enum.EventTypeObjectCreated:
		return s.handleObjectCreated(ctx, plotID, ev.Payload, tx)
	case enum.EventTypeObjectUpdated:
		return s.handleObjectUpdated(ctx, plotID, ev.Payload, tx)
	case enum.EventTypeObjectDeleted:
		return s.handleObjectDeleted(ctx, plotID, ev.Payload, tx)
	case enum.EventTypeCropPlanted:
		return s.handleCropPlanted(ctx, plotID, ev.Payload, tx)
	case enum.EventTypeCropRemoved:
		return s.handleCropRemoved(ctx, plotID, ev.Payload, tx)
	case enum.EventTypeCellShadeUpdated:
		return s.handleCellShadeUpdated(ctx, plotID, ev.Payload, tx)
	case enum.EventTypePlotResized:
		return s.handlePlotResized(ctx, plotID, ev.Payload, tx)
	default:
		return apperrors.ErrBadRequest
	}
}

func rectanglesOverlap(x1, y1, w1, h1, x2, y2, w2, h2 int) bool {
	if x1 >= x2+w2 || x2 >= x1+w1 || y1 >= y2+h2 || y2 >= y1+h1 {
		return false
	}
	return true
}

func (s *EventService) handleBedCreated(ctx context.Context, plotID string, payloadStr json.RawMessage, tx pgx.Tx) error {
	var req dto.BedCreatedRequest
	if err := json.Unmarshal(payloadStr, &req); err != nil {
		return apperrors.ErrBadRequest
	}

	plot, err := s.plotRepo.WithTx(tx).GetPlotByID(ctx, plotID)
	if err != nil {
		return s.seh.HandleError(err, "get plot for bed creation")
	}

	if req.XStart < 0 || req.YStart < 0 ||
		req.XStart+req.Width > plot.GridCols ||
		req.YStart+req.Height > plot.GridRows {
		return apperrors.ErrBadRequest
	}

	beds, err := s.bedRepo.WithTx(tx).GetBedsByPlot(ctx, plotID)
	if err != nil {
		return s.seh.HandleError(err, "get beds for overlap check")
	}
	for _, bed := range beds {
		if rectanglesOverlap(req.XStart, req.YStart, req.Width, req.Height,
			bed.XStart, bed.YStart, bed.Width, bed.Height) {
			return apperrors.ErrConflict
		}
	}

	objects, err := s.objectRepo.WithTx(tx).GetObjectsByPlot(ctx, plotID)
	if err != nil {
		return s.seh.HandleError(err, "get objects for overlap check")
	}
	for _, obj := range objects {
		if rectanglesOverlap(req.XStart, req.YStart, req.Width, req.Height,
			obj.XStart, obj.YStart, obj.Width, obj.Height) {
			return apperrors.ErrConflict
		}
	}

	bedID := uuid.New().String()

	if err = s.bedRepo.WithTx(tx).CreateBedWithID(ctx, bedID, plotID, req.Name,
		req.XStart, req.YStart, req.Width, req.Height); err != nil {
		return s.seh.HandleError(err, "insert bed")
	}

	eventPayload := payload.BedCreated{
		BedID:  bedID,
		Name:   req.Name,
		XStart: req.XStart,
		YStart: req.YStart,
		Width:  req.Width,
		Height: req.Height,
	}
	rawPayload, err := json.Marshal(eventPayload)
	if err != nil {
		return s.seh.HandleError(err, "marshal event")
	}
	if _, err = s.eventRepo.WithTx(tx).AppendEvent(ctx, plotID, enum.EventTypeBedCreated, rawPayload); err != nil {
		return s.seh.HandleError(err, "append event")
	}

	return nil
}

func (s *EventService) handleBedUpdated(ctx context.Context, plotID string, payloadStr json.RawMessage, tx pgx.Tx) error {
	var req dto.BedUpdatedRequest
	if err := json.Unmarshal(payloadStr, &req); err != nil {
		return apperrors.ErrBadRequest
	}

	bed, err := s.bedRepo.WithTx(tx).GetBedByID(ctx, req.BedID)
	if err != nil {
		return s.seh.HandleError(err, "get bed for update")
	}
	if bed.PlotID != plotID {
		return apperrors.ErrNotFound
	}

	if (req.XStart != nil || req.YStart != nil || req.Width != nil || req.Height != nil) && bed.CurrentCropID != nil {
		return apperrors.ErrConflict
	}

	newX := bed.XStart
	newY := bed.YStart
	newW := bed.Width
	newH := bed.Height
	newName := bed.Name
	if req.XStart != nil {
		newX = *req.XStart
	}
	if req.YStart != nil {
		newY = *req.YStart
	}
	if req.Width != nil {
		newW = *req.Width
	}
	if req.Height != nil {
		newH = *req.Height
	}
	if req.Name != nil {
		newName = *req.Name
	}

	plot, err := s.plotRepo.WithTx(tx).GetPlotByID(ctx, plotID)
	if err != nil {
		return s.seh.HandleError(err, "get plot for bed update")
	}
	if newX < 0 || newY < 0 || newX+newW > plot.GridCols || newY+newH > plot.GridRows {
		return apperrors.ErrBadRequest
	}

	beds, err := s.bedRepo.WithTx(tx).GetBedsByPlot(ctx, plotID)
	if err != nil {
		return s.seh.HandleError(err, "get beds for overlap check")
	}
	for _, b := range beds {
		if b.BedID == req.BedID {
			continue
		}
		if rectanglesOverlap(newX, newY, newW, newH, b.XStart, b.YStart, b.Width, b.Height) {
			return apperrors.ErrConflict
		}
	}

	objects, err := s.objectRepo.WithTx(tx).GetObjectsByPlot(ctx, plotID)
	if err != nil {
		return s.seh.HandleError(err, "get objects for overlap check")
	}
	for _, obj := range objects {
		if rectanglesOverlap(newX, newY, newW, newH, obj.XStart, obj.YStart, obj.Width, obj.Height) {
			return apperrors.ErrConflict
		}
	}

	if err = s.bedRepo.WithTx(tx).UpdateBed(ctx, req.BedID, req.Name, req.XStart, req.YStart, req.Width, req.Height); err != nil {
		return s.seh.HandleError(err, "update bed")
	}

	eventPayload := payload.BedUpdated{
		BedID:  req.BedID,
		Name:   newName,
		XStart: newX,
		YStart: newY,
		Width:  newW,
		Height: newH,
	}
	rawPayload, err := json.Marshal(eventPayload)
	if err != nil {
		return s.seh.HandleError(err, "marshal event")
	}
	if _, err = s.eventRepo.WithTx(tx).AppendEvent(ctx, plotID, enum.EventTypeBedUpdated, rawPayload); err != nil {
		return s.seh.HandleError(err, "append event")
	}

	return nil
}

func (s *EventService) handleBedDeleted(ctx context.Context, plotID string, payloadStr json.RawMessage, tx pgx.Tx) error {
	var req dto.BedDeletedRequest
	if err := json.Unmarshal(payloadStr, &req); err != nil {
		return apperrors.ErrBadRequest
	}

	bed, err := s.bedRepo.WithTx(tx).GetBedByID(ctx, req.BedID)
	if err != nil {
		return s.seh.HandleError(err, "get bed for deletion")
	}
	if bed.PlotID != plotID {
		return apperrors.ErrNotFound
	}

	if bed.CurrentCropID != nil {
		return apperrors.ErrConflict
	}

	if err = s.bedRepo.WithTx(tx).DeleteBed(ctx, req.BedID); err != nil {
		return s.seh.HandleError(err, "delete bed")
	}

	eventPayload := payload.BedDeleted{BedID: req.BedID}
	rawPayload, err := json.Marshal(eventPayload)
	if err != nil {
		return s.seh.HandleError(err, "marshal event")
	}
	if _, err = s.eventRepo.WithTx(tx).AppendEvent(ctx, plotID, enum.EventTypeBedDeleted, rawPayload); err != nil {
		return s.seh.HandleError(err, "append event")
	}

	return nil
}

func (s *EventService) handleObjectCreated(ctx context.Context, plotID string, payloadStr json.RawMessage, tx pgx.Tx) error {
	var req dto.ObjectCreatedRequest
	if err := json.Unmarshal(payloadStr, &req); err != nil {
		return apperrors.ErrBadRequest
	}
	return errors.New("not implemented")
}

func (s *EventService) handleObjectUpdated(ctx context.Context, plotID string, payloadStr json.RawMessage, tx pgx.Tx) error {
	var req dto.ObjectUpdatedRequest
	if err := json.Unmarshal(payloadStr, &req); err != nil {
		return apperrors.ErrBadRequest
	}
	return errors.New("not implemented")
}

func (s *EventService) handleObjectDeleted(ctx context.Context, plotID string, payloadStr json.RawMessage, tx pgx.Tx) error {
	var req dto.ObjectDeletedRequest
	if err := json.Unmarshal(payloadStr, &req); err != nil {
		return apperrors.ErrBadRequest
	}
	return errors.New("not implemented")
}

func (s *EventService) handleCropPlanted(ctx context.Context, plotID string, payloadStr json.RawMessage, tx pgx.Tx) error {
	var req dto.CropPlantedRequest
	if err := json.Unmarshal(payloadStr, &req); err != nil {
		return apperrors.ErrBadRequest
	}

	bed, err := s.bedRepo.WithTx(tx).GetBedByID(ctx, req.BedID)
	if err != nil {
		return s.seh.HandleError(err, "get bed for planting")
	}
	if bed.PlotID != plotID {
		return apperrors.ErrNotFound
	}

	plantDate := time.Now().UTC()
	if req.PlantDate != nil && *req.PlantDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.PlantDate)
		if err != nil {
			return apperrors.ErrBadRequest
		}
		plantDate = parsed
	}

	if err = s.bedRepo.WithTx(tx).SetCrop(ctx, req.BedID, req.CropID, plantDate); err != nil {
		return s.seh.HandleError(err, "set crop")
	}

	eventPayload := payload.CropPlanted{
		BedID:     req.BedID,
		CropID:    req.CropID,
		PlantDate: plantDate.Format("2006-01-02"),
	}
	rawPayload, err := json.Marshal(eventPayload)
	if err != nil {
		return s.seh.HandleError(err, "marshal event")
	}
	if _, err = s.eventRepo.WithTx(tx).AppendEvent(ctx, plotID, enum.EventTypeCropPlanted, rawPayload); err != nil {
		return s.seh.HandleError(err, "append event")
	}

	return nil
}

func (s *EventService) handleCropRemoved(ctx context.Context, plotID string, payloadStr json.RawMessage, tx pgx.Tx) error {
	var req dto.CropRemovedRequest
	if err := json.Unmarshal(payloadStr, &req); err != nil {
		return apperrors.ErrBadRequest
	}

	bed, err := s.bedRepo.WithTx(tx).GetBedByID(ctx, req.BedID)
	if err != nil {
		return s.seh.HandleError(err, "get bed for crop removal")
	}
	if bed.PlotID != plotID {
		return apperrors.ErrNotFound
	}
	if bed.CurrentCropID == nil {
		return apperrors.ErrConflict
	}
	if bed.PlantDate == nil {
		return apperrors.ErrInternal
	}

	harvestDate := time.Now().UTC()
	if req.Date != nil && *req.Date != "" {
		parsed, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			return apperrors.ErrBadRequest
		}
		harvestDate = parsed
	}

	var cells []model.CellCoord
	for x := bed.XStart; x < bed.XStart+bed.Width; x++ {
		for y := bed.YStart; y < bed.YStart+bed.Height; y++ {
			cells = append(cells, model.CellCoord{X: x, Y: y})
		}
	}

	if req.Harvested {
		var historyRecords []model.CellHistory
		for _, cell := range cells {
			historyRecords = append(historyRecords, model.CellHistory{
				PlotID:      plotID,
				XIndex:      cell.X,
				YIndex:      cell.Y,
				CropID:      *bed.CurrentCropID,
				PlantDate:   *bed.PlantDate,
				HarvestDate: harvestDate,
			})
		}
		if err = s.historyRepo.WithTx(tx).AddHistoryRecords(ctx, historyRecords); err != nil {
			return s.seh.HandleError(err, "insert history")
		}
	}

	if err = s.bedRepo.WithTx(tx).ClearCrop(ctx, req.BedID); err != nil {
		return s.seh.HandleError(err, "clear crop")
	}

	eventPayload := payload.CropRemoved{
		BedID:     req.BedID,
		CropID:    *bed.CurrentCropID,
		Harvested: req.Harvested,
		Date:      harvestDate.Format("2006-01-02"),
	}
	rawPayload, _ := json.Marshal(eventPayload)
	if _, err = s.eventRepo.WithTx(tx).AppendEvent(ctx, plotID, enum.EventTypeCropRemoved, rawPayload); err != nil {
		return s.seh.HandleError(err, "append event")
	}

	return nil
}

func (s *EventService) handleCellShadeUpdated(ctx context.Context, plotID string, payloadStr json.RawMessage, tx pgx.Tx) error {
	var req dto.CellShadeUpdatedRequest
	if err := json.Unmarshal(payloadStr, &req); err != nil {
		return apperrors.ErrBadRequest
	}
	return errors.New("not implemented")
}

func (s *EventService) handlePlotResized(ctx context.Context, plotID string, payloadStr json.RawMessage, tx pgx.Tx) error {
	var req dto.PlotResizedRequest
	if err := json.Unmarshal(payloadStr, &req); err != nil {
		return apperrors.ErrBadRequest
	}
	return errors.New("not implemented")
}
