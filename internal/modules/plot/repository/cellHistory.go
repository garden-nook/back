package repository

//
//import (
//	"context"
//	"fmt"
//	"garden-nook/internal/modules/plot/dto"
//	"garden-nook/internal/modules/plot/model"
//	"garden-nook/internal/pkg/database"
//	"time"
//
//	"github.com/jackc/pgx/v5"
//)
//
//type HistoryRepo struct {
//	db     database.DBTX
//	mapper *database.ErrorMapper
//}
//
//func NewHistoryRepo(db database.DBTX, mapper *database.ErrorMapper) *HistoryRepo {
//	return &HistoryRepo{db: db, mapper: mapper}
//}
//
//func (r *HistoryRepo) WithTx(tx pgx.Tx) *HistoryRepo {
//	return &HistoryRepo{db: tx, mapper: r.mapper}
//}
//
//// AddHistoryRecord вставляет новую запись в историю для одной ячейки.
//func (r *HistoryRepo) AddHistoryRecord(ctx context.Context, req dto.AddHistoryRequest) (string, error) {
//	var historyID string
//	err := r.db.QueryRow(ctx,
//		`INSERT INTO cell_crop_history
//		 (plot_id, x_index, y_index, crop_id, family_id, bed_id, plant_date)
//		 VALUES ($1, $2, $3, $4, $5, $6, $7)
//		 RETURNING history_id`,
//		req.PlotID, req.XIndex, req.YIndex, req.CropID, req.FamilyID,
//		req.BedID, req.PlantDate,
//	).Scan(&historyID)
//	if err != nil {
//		return "", r.mapper.Map(err)
//	}
//	return historyID, nil
//}
//
//// UpdateHarvest обновляет дату сбора и метаданные для конкретной ячейки.
//func (r *HistoryRepo) UpdateHarvest(ctx context.Context, plotID string, xIndex, yIndex int, harvestDate time.Time, metadata *string) error {
//	query := `UPDATE cell_crop_history SET harvest_date = $1, metadata = $2
//		WHERE plot_id = $3 AND x_index = $4 AND y_index = $5
//		  AND harvest_date IS NULL
//		RETURNING history_id`
//	var historyID string
//	err := r.db.QueryRow(ctx, query,
//		harvestDate, metadata, plotID, xIndex, yIndex,
//	).Scan(&historyID)
//	if err != nil {
//		return r.mapper.Map(err)
//	}
//	return nil
//}
//
//// DeleteHistory удаляет записи истории для указанных ячеек, где harvest_date IS NULL.
//func (r *HistoryRepo) DeleteHistory(ctx context.Context, plotID string, cells []model.CellCoord) error {
//	batch := &pgx.Batch{}
//	for _, cell := range cells {
//		sql := `DELETE FROM cell_crop_history
//			WHERE plot_id = $1 AND x_index = $2 AND y_index = $3
//			  AND harvest_date IS NULL`
//		batch.Queue(sql, plotID, cell.XIndex, cell.YIndex)
//	}
//	results := r.db.SendBatch(ctx, batch)
//	defer results.Close()
//	for i := 0; i < batch.Len(); i++ {
//		if _, err := results.Exec(); err != nil {
//			return r.mapper.Map(err)
//		}
//	}
//	return nil
//}
//
//// GetHistoryForCell возвращает историю посадок для конкретной ячейки (сортировка по дате посадки DESC).
//func (r *HistoryRepo) GetHistoryForCell(ctx context.Context, plotID string, xIndex, yIndex int, p *database.Pagination) ([]model.CellHistory, int, error) {
//	where := `WHERE plot_id = $1 AND x_index = $2 AND y_index = $3`
//	filterArgs := []any{plotID, xIndex, yIndex}
//	argIdx := 4
//
//	return r.queryHistory(ctx, where, filterArgs, argIdx, p)
//}
//
//// GetHistoryForPlot возвращает историю посадок для всего участка с возможностью фильтрации.
//func (r *HistoryRepo) GetHistoryForPlot(ctx context.Context, plotID string, filter dto.HistoryFilter, p *database.Pagination) ([]model.CellHistory, int, error) {
//	where := `WHERE plot_id = $1`
//	filterArgs := []any{plotID}
//	argIdx := 2
//
//	if filter.BedID != nil {
//		where += fmt.Sprintf(" AND bed_id = $%d", argIdx)
//		filterArgs = append(filterArgs, *filter.BedID)
//		argIdx++
//	}
//	if filter.CropID != nil {
//		where += fmt.Sprintf(" AND crop_id = $%d", argIdx)
//		filterArgs = append(filterArgs, *filter.CropID)
//		argIdx++
//	}
//	if filter.FamilyID != nil {
//		where += fmt.Sprintf(" AND family_id = $%d", argIdx)
//		filterArgs = append(filterArgs, *filter.FamilyID)
//		argIdx++
//	}
//	if filter.FromDate != nil {
//		where += fmt.Sprintf(" AND plant_date >= $%d", argIdx)
//		filterArgs = append(filterArgs, *filter.FromDate)
//		argIdx++
//	}
//	if filter.ToDate != nil {
//		where += fmt.Sprintf(" AND plant_date <= $%d", argIdx)
//		filterArgs = append(filterArgs, *filter.ToDate)
//		argIdx++
//	}
//
//	return r.queryHistory(ctx, where, filterArgs, argIdx, p)
//}
//
//// queryHistory – внутренняя функция для выполнения запроса истории с пагинацией.
//func (r *HistoryRepo) queryHistory(ctx context.Context, where string, filterArgs []any, argIdx int, p *database.Pagination) ([]model.CellHistory, int, error) {
//	baseSQL := `SELECT history_id, plot_id, x_index, y_index, crop_id, family_id, bed_id, plant_date, harvest_date, metadata
//		FROM cell_crop_history `
//
//	if p == nil {
//		rows, err := r.db.Query(ctx, baseSQL+where+` ORDER BY plant_date DESC`)
//		if err != nil {
//			return nil, 0, r.mapper.Map(err)
//		}
//		defer rows.Close()
//		entries, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.CellHistory])
//		if err != nil {
//			return nil, 0, r.mapper.Map(err)
//		}
//		return entries, len(entries), nil
//	}
//
//	batch := &pgx.Batch{}
//	countSQL := `SELECT COUNT(*) FROM cell_crop_history ` + where
//	batch.Queue(countSQL, filterArgs...)
//
//	pagSQL, pagArgs := p.SQL(argIdx)
//	dataSQL := baseSQL + where + ` ORDER BY plant_date DESC` + pagSQL
//	dataArgs := append(filterArgs, pagArgs...)
//	batch.Queue(dataSQL, dataArgs...)
//
//	results := r.db.SendBatch(ctx, batch)
//	defer results.Close()
//
//	var total int
//	if err := results.QueryRow().Scan(&total); err != nil {
//		return nil, 0, r.mapper.Map(err)
//	}
//	rows, err := results.Query()
//	if err != nil {
//		return nil, 0, r.mapper.Map(err)
//	}
//	defer rows.Close()
//
//	entries, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.CellHistory])
//	if err != nil {
//		return nil, 0, r.mapper.Map(err)
//	}
//	return entries, total, nil
//}
