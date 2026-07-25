package repository

import (
	"context"
	"garden-nook/internal/modules/plot/model"
	"garden-nook/internal/pkg/database"

	"github.com/jackc/pgx/v5"
)

type HistoryRepo struct {
	db     database.DBTX
	mapper *database.ErrorMapper
}

func NewHistoryRepo(db database.DBTX, mapper *database.ErrorMapper) *HistoryRepo {
	return &HistoryRepo{db: db, mapper: mapper}
}

func (r *HistoryRepo) WithTx(tx pgx.Tx) *HistoryRepo {
	return &HistoryRepo{db: tx, mapper: r.mapper}
}

func (r *HistoryRepo) AddHistoryRecords(ctx context.Context, records []model.CellHistory) error {
	batch := &pgx.Batch{}
	for _, rec := range records {
		batch.Queue(
			`INSERT INTO cell_crop_history (plot_id, x_index, y_index, crop_id, plant_date, harvest_date)
             VALUES ($1, $2, $3, $4, $5, $6)`,
			rec.PlotID, rec.XIndex, rec.YIndex, rec.CropID, rec.PlantDate, rec.HarvestDate,
		)
	}
	br := r.db.SendBatch(ctx, batch)
	defer br.Close()
	for i := 0; i < batch.Len(); i++ {
		if _, err := br.Exec(); err != nil {
			return r.mapper.Map(err)
		}
	}
	return nil
}

func (r *HistoryRepo) GetPredecessorsForBed(ctx context.Context, plotID string, xStart, yStart, width, height int) ([]model.PredecessorInfo, error) {
	query := `
        SELECT h.crop_id, c.family_id, MAX(h.harvest_date)
        FROM cell_crop_history h
        JOIN crops c ON c.id = h.crop_id
        WHERE h.plot_id = $1
          AND h.x_index >= $2 AND h.x_index < $2 + $4
          AND h.y_index >= $3 AND h.y_index < $3 + $5
        GROUP BY h.crop_id, c.family_id
    `
	rows, err := r.db.Query(ctx, query, plotID, xStart, yStart, width, height)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	defer rows.Close()

	var predecessors []model.PredecessorInfo
	for rows.Next() {
		var p model.PredecessorInfo
		if err = rows.Scan(&p.CropID, &p.FamilyID, &p.LastHarvestDate); err != nil {
			return nil, r.mapper.Map(err)
		}
		predecessors = append(predecessors, p)
	}
	return predecessors, rows.Err()
}
