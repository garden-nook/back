package repository

import (
	"context"
	"fmt"
	"garden-nook/internal/modules/plot/model"
	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/database"
	"time"

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

func (r *BedRepo) UpdateBed(ctx context.Context, bedID string, name *string, xStart, yStart, width, height *int) error {
	fields := []database.SetField{
		database.NewSetField("name", name),
		database.NewSetField("x_start", xStart),
		database.NewSetField("y_start", yStart),
		database.NewSetField("width", width),
		database.NewSetField("height", height),
	}
	setSQL, setArgs := database.BuildUpdateSet(1, fields...)
	if len(setArgs) == 0 {
		return nil
	}
	query := fmt.Sprintf("UPDATE beds_ui SET %s WHERE bed_id=$%d", setSQL, len(setArgs)+1)
	args := append(setArgs, bedID)

	ct, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return r.mapper.Map(err)
	}
	if ct.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *BedRepo) SetCrop(ctx context.Context, bedID string, cropID int32, plantDate time.Time) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE beds_ui SET current_crop_id = $1, plant_date = $2 WHERE bed_id = $3 AND current_crop_id IS NULL`,
		cropID, plantDate, bedID)
	if err != nil {
		return r.mapper.Map(err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrConflict
	}
	return nil
}

func (r *BedRepo) ClearCrop(ctx context.Context, bedID string) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE beds_ui SET current_crop_id = NULL, plant_date = NULL 
         WHERE bed_id = $1 AND current_crop_id IS NOT NULL`, bedID)
	if err != nil {
		return r.mapper.Map(err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
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
