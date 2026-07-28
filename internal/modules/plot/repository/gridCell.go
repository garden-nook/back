package repository

import (
	"context"
	"database/sql"
	"garden-nook/internal/modules/plot/enum"
	"garden-nook/internal/modules/plot/model"
	"garden-nook/internal/pkg/database"

	"github.com/jackc/pgx/v5"
)

type GridCellRepo struct {
	db     database.DBTX
	mapper *database.ErrorMapper
}

func NewGridCellRepo(db database.DBTX, mapper *database.ErrorMapper) *GridCellRepo {
	return &GridCellRepo{db: db, mapper: mapper}
}

func (r *GridCellRepo) WithTx(tx pgx.Tx) *GridCellRepo {
	return &GridCellRepo{db: tx, mapper: r.mapper}
}

func (r *GridCellRepo) CreateCells(ctx context.Context, plotID string, cols, rows int) error {
	batch := &pgx.Batch{}
	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			batch.Queue(
				`INSERT INTO grid_cells (plot_id, x_index, y_index) VALUES ($1, $2, $3)`,
				plotID, x, y,
			)
		}
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

func (r *GridCellRepo) GetShadedCells(ctx context.Context, plotID string) ([]model.GridCell, error) {
	rows, err := r.db.Query(ctx,
		`SELECT plot_id, x_index, y_index, shade_level
         FROM grid_cells
         WHERE plot_id = $1 AND shade_level != $2`, plotID, enum.ShadeLevelFull)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[model.GridCell])
}

func (r *GridCellRepo) GetAverageShadeForRectangle(ctx context.Context, plotID string, xStart, yStart, width, height int) (float64, error) {
	var avg sql.NullFloat64
	err := r.db.QueryRow(ctx, `
        SELECT AVG(shade_level::float)
        FROM grid_cells
        WHERE plot_id = $1
          AND x_index >= $2 AND x_index < $2 + $4
          AND y_index >= $3 AND y_index < $3 + $5
    `, plotID, xStart, yStart, width, height).Scan(&avg)
	if err != nil {
		return 0, r.mapper.Map(err)
	}
	if !avg.Valid {
		// Если по какой-то причине ячеек нет, возвращаем значение по умолчанию (полное солнце)
		return float64(enum.ShadeLevelFull), nil
	}
	return avg.Float64, nil
}
