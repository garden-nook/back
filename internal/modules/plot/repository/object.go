package repository

import (
	"context"
	"fmt"
	"garden-nook/internal/modules/plot/model"
	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/database"

	"github.com/jackc/pgx/v5"
)

type ObjectRepo struct {
	db     database.DBTX
	mapper *database.ErrorMapper
}

func NewObjectRepo(db database.DBTX, mapper *database.ErrorMapper) *ObjectRepo {
	return &ObjectRepo{db: db, mapper: mapper}
}

func (r *ObjectRepo) WithTx(tx pgx.Tx) *ObjectRepo {
	return &ObjectRepo{db: tx, mapper: r.mapper}
}

func (r *ObjectRepo) GetObjectsByPlot(ctx context.Context, plotID string) ([]model.Object, error) {
	rows, err := r.db.Query(ctx,
		`SELECT object_id, plot_id, name, object_type, x_start, y_start, width, height
         FROM objects_ui
         WHERE plot_id = $1
         ORDER BY name`, plotID)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Object])
}

func (r *ObjectRepo) CreateObjectWithID(ctx context.Context, objectID, plotID, name string, objectType int32, xStart, yStart, width, height int) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO objects_ui (object_id, plot_id, name, object_type, x_start, y_start, width, height)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		objectID, plotID, name, objectType, xStart, yStart, width, height)
	if err != nil {
		return r.mapper.Map(err)
	}
	return nil
}

func (r *ObjectRepo) UpdateObject(ctx context.Context, objectID string, name *string, objectType *int32, xStart, yStart, width, height *int) error {
	fields := []database.SetField{
		database.NewSetField("name", name),
		database.NewSetField("object_type", objectType),
		database.NewSetField("x_start", xStart),
		database.NewSetField("y_start", yStart),
		database.NewSetField("width", width),
		database.NewSetField("height", height),
	}
	setSQL, setArgs := database.BuildUpdateSet(1, fields...)
	if len(setArgs) == 0 {
		return nil
	}
	query := fmt.Sprintf("UPDATE objects_ui SET %s WHERE object_id=$%d", setSQL, len(setArgs)+1)
	args := append(setArgs, objectID)

	ct, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return r.mapper.Map(err)
	}
	if ct.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *ObjectRepo) GetObjectByID(ctx context.Context, objectID string) (*model.Object, error) {
	query := `SELECT object_id, plot_id, name, object_type, x_start, y_start, width, height
              FROM objects_ui
              WHERE object_id = $1`
	row, err := r.db.Query(ctx, query, objectID)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	defer row.Close()

	obj, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[model.Object])
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	return obj, nil
}

func (r *ObjectRepo) DeleteObject(ctx context.Context, objectID string) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM objects_ui WHERE object_id = $1`, objectID)
	if err != nil {
		return r.mapper.Map(err)
	}
	if ct.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}
