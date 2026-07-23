package plots

//
//import (
//	"context"
//	"encoding/json"
//	"log/slog"
//	"time"
//)
//
//// Projector применяет события из event_store к read model.
//// В production лучше вынести в async worker (Kafka/RabbitMQ consumer).
//type Projector struct {
//	repo *Repository
//	log  *slog.Logger
//}
//
//func NewProjector(repo *Repository, log *slog.Logger) *Projector {
//	return &Projector{repo: repo, log: log}
//}
//
//// ApplyEvent применяет одно событие к read model.
//func (p *Projector) ApplyEvent(ctx context.Context, plotID string, eventType int, payloadJSON []byte) error {
//	switch eventType {
//	case EventPlotCreated:
//		var payload PlotCreatedPayload
//		if err := json.Unmarshal(payloadJSON, &payload); err != nil {
//			return err
//		}
//		// Инициализируем сетку
//		// Для простоты считаем, что origin = (0, 0). В reality нужно вычислять из boundary.
//		return p.repo.InitializeGrid(ctx, plotID, payload.GridCols, payload.GridRows, payload.GridCellSize, 0, 0)
//
//	case EventBedCreated:
//		var payload BedCreatedPayload
//		if err := json.Unmarshal(payloadJSON, &payload); err != nil {
//			return err
//		}
//		bed := &Bed{
//			BedID:  payload.BedID,
//			PlotID: payload.PlotID,
//			Name:   payload.Name,
//			Geom:   payload.Geom,
//		}
//		return p.repo.CreateBed(ctx, bed)
//
//	case EventBedUpdated:
//		var payload BedUpdatedPayload
//		if err := json.Unmarshal(payloadJSON, &payload); err != nil {
//			return err
//		}
//		return p.repo.UpdateBed(ctx, payload.BedID, payload.PlotID, payload.Name, payload.NewGeom)
//
//	case EventBedDeleted:
//		var payload BedDeletedPayload
//		if err := json.Unmarshal(payloadJSON, &payload); err != nil {
//			return err
//		}
//		return p.repo.SoftDeleteBed(ctx, payload.BedID, payload.PlotID)
//
//	case EventObjectCreated:
//		var payload ObjectCreatedPayload
//		if err := json.Unmarshal(payloadJSON, &payload); err != nil {
//			return err
//		}
//		obj := &UIObject{
//			ObjectID:   payload.ObjectID,
//			PlotID:     payload.PlotID,
//			Name:       payload.Name,
//			ObjectType: payload.ObjectType,
//			Geom:       payload.Geom,
//		}
//		return p.repo.CreateObject(ctx, obj)
//
//	case EventObjectDeleted:
//		// Аналогично BedDeleted
//		return nil
//
//	case EventCropPlanted:
//		var payload CropPlantedPayload
//		if err := json.Unmarshal(payloadJSON, &payload); err != nil {
//			return err
//		}
//		// Парсим plant_date
//		plantDate, err := parseDate(payload.PlantDate)
//		if err != nil {
//			return err
//		}
//		return p.repo.PlantCrop(ctx, payload.PlotID, payload.BedID, payload.CropID, plantDate, payload.CellIDs)
//
//	case EventCropHarvested:
//		var payload CropHarvestedPayload
//		if err := json.Unmarshal(payloadJSON, &payload); err != nil {
//			return err
//		}
//		harvestDate, err := parseDate(payload.HarvestDate)
//		if err != nil {
//			return err
//		}
//		return p.repo.HarvestCrop(ctx, payload.PlotID, payload.BedID, harvestDate, payload.YieldKg, payload.CellIDs)
//
//	default:
//		p.log.Warn("unknown event type", "type", eventType)
//		return nil
//	}
//}
//
//func parseDate(s string) (time.Time, error) {
//	return time.Parse("2006-01-02", s)
//}
