package repository

import (
	"context"
	"garden-nook/internal/modules/plot/model"
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
