package repository

import (
	"context"
	"garden-nook/internal/modules/plot/model"
	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/database"

	"github.com/jackc/pgx/v5"
)

type BedRepo struct {
	db     database.DBTX
	mapper *database.ErrorMapper
}

func NewBedRepo(db database.DBTX, mapper *database.ErrorMapper) *BedRepo {
	return &BedRepo{db: db, mapper: mapper}
}

func (r *BedRepo) WithTx(tx pgx.Tx) *BedRepo {
	return &BedRepo{db: tx, mapper: r.mapper}
}

func (r *BedRepo) GetBedsByPlot(ctx context.Context, plotID string) ([]model.Bed, error) {
	rows, err := r.db.Query(ctx,
		`SELECT bed_id, plot_id, name, x_start, y_start, width, height, current_crop_id, plant_date
         FROM beds_ui
         WHERE plot_id = $1
         ORDER BY name`, plotID)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Bed])
}

func (r *BedRepo) CreateBedWithID(ctx context.Context, bedID, plotID, name string, xStart, yStart, width, height int) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO beds_ui (bed_id, plot_id, name, x_start, y_start, width, height)
         VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		bedID, plotID, name, xStart, yStart, width, height)
	if err != nil {
		return r.mapper.Map(err)
	}
	return nil
}

func (r *BedRepo) GetBedByID(ctx context.Context, bedID string) (*model.Bed, error) {
	query := `SELECT bed_id, plot_id, name, x_start, y_start, width, height, current_crop_id, plant_date
              FROM beds_ui
              WHERE bed_id = $1`
	row, err := r.db.Query(ctx, query, bedID)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	defer row.Close()

	bed, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[model.Bed])
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	return bed, nil
}

func (r *BedRepo) DeleteBed(ctx context.Context, bedID string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM beds_ui WHERE bed_id = $1`, bedID)
	if err != nil {
		return r.mapper.Map(err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}
