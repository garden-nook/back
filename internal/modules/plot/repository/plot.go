package repository

import (
	"context"
	"fmt"
	"garden-nook/internal/modules/plot/dto"
	"garden-nook/internal/modules/plot/model"
	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlotRepo struct {
	db     database.DBTX
	mapper *database.ErrorMapper
}

func NewPlotRepo(pool *pgxpool.Pool, mapper *database.ErrorMapper) *PlotRepo {
	return &PlotRepo{db: pool, mapper: mapper}
}

func (r *PlotRepo) WithTx(tx pgx.Tx) *PlotRepo {
	return &PlotRepo{db: tx, mapper: r.mapper}
}

func (r *PlotRepo) ListPlots(ctx context.Context, ownerID string, p *database.Pagination) ([]model.Plot, int, error) {
	baseQuery := `SELECT p.plot_id, p.name, p.soil_type, st.name AS soil_name,
	              p.grid_cell_size, p.grid_cols, p.grid_rows
	              FROM plots p
	              JOIN soil_types st ON st.id = p.soil_type`
	whereClause := "WHERE p.owner_id = $1"
	whereArgs := []any{ownerID}
	orderBy := "p.name"

	plots, total, err := database.ListQuery[model.Plot](ctx, r.db, baseQuery, whereClause, whereArgs, orderBy, p)
	if err != nil {
		return nil, 0, r.mapper.Map(err)
	}
	return plots, total, nil
}

func (r *PlotRepo) CreatePlot(ctx context.Context, req *model.CreatePlotModel, ownerID string) (string, error) {
	var plotID string
	err := r.db.QueryRow(ctx,
		`INSERT INTO plots (name, owner_id, soil_type, boundary,
			            grid_cell_size, grid_cols, grid_rows)
		 VALUES ($1, $2, $3,
		         ST_MakeEnvelope(0, 0, $4, $5, 3857),
		         $6, $7, $8)
		 RETURNING plot_id`,
		req.Name, ownerID, req.SoilTypeID,
		req.BoundaryWidth, req.BoundaryHeight,
		req.GridCellSize, req.GridCols, req.GridRows,
	).Scan(&plotID)
	if err != nil {
		return "", r.mapper.Map(err)
	}
	return plotID, nil
}

func (r *PlotRepo) UpdatePlot(ctx context.Context, plotID, ownerID string, req dto.UpdatePlotRequest) (string, error) {
	fields := []database.SetField{
		database.NewSetField("name", req.Name),
		database.NewSetField("soil_type", req.SoilTypeID),
	}
	setSQL, setArgs := database.BuildUpdateSet(1, fields...)
	if len(setArgs) == 0 {
		return plotID, nil
	}
	query := fmt.Sprintf(
		"UPDATE plots SET %s WHERE plot_id = $%d AND owner_id = $%d RETURNING plot_id",
		setSQL, len(setArgs)+1, len(setArgs)+2,
	)
	args := append(setArgs, plotID, ownerID)

	var updatedID string
	err := r.db.QueryRow(ctx, query, args...).Scan(&updatedID)
	if err != nil {
		return "", r.mapper.Map(err)
	}
	return updatedID, nil
}

func (r *PlotRepo) DeletePlot(ctx context.Context, plotID, ownerID string) error {
	ct, err := r.db.Exec(ctx,
		`DELETE FROM plots WHERE plot_id = $1 AND owner_id = $2`,
		plotID, ownerID)
	if err != nil {
		return r.mapper.Map(err)
	}
	if ct.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *PlotRepo) GetPlotByOwnerAndID(ctx context.Context, plotID, ownerID string) (*model.Plot, error) {
	query := `SELECT p.plot_id, p.name, p.soil_type, st.name AS soil_name,
                     p.grid_cell_size, p.grid_cols, p.grid_rows
              FROM plots p
              JOIN soil_types st ON st.id = p.soil_type
              WHERE p.plot_id = $1 AND p.owner_id = $2`

	row, err := r.db.Query(ctx, query, plotID, ownerID)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	defer row.Close()

	plot, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[model.Plot])
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	return plot, nil
}

func (r *PlotRepo) GetPlotByID(ctx context.Context, plotID string) (*model.Plot, error) {
	query := `SELECT p.plot_id, p.name, p.soil_type, st.name AS soil_name,
                     p.grid_cell_size, p.grid_cols, p.grid_rows
              FROM plots p
              JOIN soil_types st ON st.id = p.soil_type
              WHERE p.plot_id = $1`

	row, err := r.db.Query(ctx, query, plotID)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	defer row.Close()

	plot, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[model.Plot])
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	return plot, nil
}
