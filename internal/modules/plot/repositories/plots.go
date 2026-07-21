package repositories

import (
	"context"
	"fmt"
	"garden-nook/internal/modules/plot/models"
	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db     database.DBTX
	mapper *database.ErrorMapper
}

func NewRepository(pool *pgxpool.Pool, mapper *database.ErrorMapper) *Repository {
	return &Repository{db: pool, mapper: mapper}
}

// WithTx возвращает копию репозитория, работающую в рамках переданной транзакции.
func (r *Repository) WithTx(tx pgx.Tx) *Repository {
	return &Repository{db: tx, mapper: r.mapper}
}

func (r *Repository) ListPlots(ctx context.Context, ownerID string, p *database.Pagination) ([]models.Plot, int, error) {
	where := `WHERE p.owner_id = $1 AND p.is_deleted = FALSE`
	filterArgs := []any{ownerID}
	argIdx := 2

	baseSelect := `SELECT p.plot_id, p.name, p.soil_type, st.name AS soil_name,
	       			   p.area_sq_m, p.grid_cell_size, p.grid_cols, p.grid_rows
	                FROM plots p
	                JOIN soil_types st ON st.id = p.soil_type `

	if p == nil {
		rows, err := r.db.Query(ctx, baseSelect+where+` ORDER BY p.name`, filterArgs...)
		if err != nil {
			return nil, 0, r.mapper.Map(err)
		}
		defer rows.Close()
		plots, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Plot])
		if err != nil {
			return nil, 0, r.mapper.Map(err)
		}
		return plots, len(plots), nil
	}

	batch := &pgx.Batch{}
	countSQL := `SELECT COUNT(*) FROM plots p ` + where
	batch.Queue(countSQL, filterArgs...)

	pagSQL, pagArgs := p.SQL(argIdx)
	dataSQL := baseSelect + where + ` ORDER BY p.name` + pagSQL
	dataArgs := append(filterArgs, pagArgs...)
	batch.Queue(dataSQL, dataArgs...)

	results := r.db.SendBatch(ctx, batch)
	defer results.Close()

	var total int
	if err := results.QueryRow().Scan(&total); err != nil {
		return nil, 0, r.mapper.Map(err)
	}
	rows, err := results.Query()
	if err != nil {
		return nil, 0, r.mapper.Map(err)
	}
	defer rows.Close()
	plots, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Plot])
	if err != nil {
		return nil, 0, r.mapper.Map(err)
	}
	return plots, total, nil
}

// GetPlotByID возвращает участок по его ID без проверки владельца.
func (r *Repository) GetPlotByID(ctx context.Context, plotID string) (*models.Plot, error) {
	query := `SELECT p.plot_id, p.name, p.soil_type, st.name AS soil_name,
                     p.area_sq_m, p.grid_cell_size, p.grid_cols, p.grid_rows
              FROM plots p
              JOIN soil_types st ON st.id = p.soil_type
              WHERE p.plot_id = $1 AND p.is_deleted = FALSE`

	row, err := r.db.Query(ctx, query, plotID)
	if err != nil {
		return nil, r.mapper.Map(err)
	}

	plot, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[models.Plot])
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	return plot, nil
}

// GetPlotByOwnerAndID возвращает участок, если он принадлежит указанному владельцу.
func (r *Repository) GetPlotByOwnerAndID(ctx context.Context, plotID, ownerID string) (*models.Plot, error) {
	query := `SELECT p.plot_id, p.name, p.soil_type, st.name AS soil_name,
                     p.area_sq_m, p.grid_cell_size, p.grid_cols, p.grid_rows
              FROM plots p
              JOIN soil_types st ON st.id = p.soil_type
              WHERE p.plot_id = $1 AND p.owner_id = $2 AND p.is_deleted = FALSE`

	row, err := r.db.Query(ctx, query, plotID, ownerID)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	defer row.Close()

	plot, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[models.Plot])
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	return plot, nil
}

func (r *Repository) CreatePlot(ctx context.Context, req *models.CreatePlotModel, ownerID string) (string, error) {
	var plotID string
	err := r.db.QueryRow(ctx,
		`INSERT INTO plots (name, owner_id, soil_type, boundary, area_sq_m,
			            grid_cell_size, grid_cols, grid_rows)
		 VALUES ($1, $2, $3,
		         ST_MakeEnvelope(0, 0, $4, $5, 3857),
		         $6, $7, $8, $9)
		 RETURNING plot_id`,
		req.Name, ownerID, req.SoilTypeID,
		req.BoundaryWidth, req.BoundaryHeight,
		req.AreaSqM, req.GridCellSize, req.GridCols, req.GridRows,
	).Scan(&plotID)
	if err != nil {
		return "", r.mapper.Map(err)
	}
	return plotID, nil
}

func (r *Repository) UpdatePlot(ctx context.Context, plotID, ownerID string, req models.UpdatePlotRequest) (string, error) {
	fields := []database.SetField{
		database.NewSetField("name", req.Name),
		database.NewSetField("soil_type", req.SoilTypeID),
	}
	setSQL, setArgs := database.BuildUpdateSet(1, fields...)
	if len(setArgs) == 0 {
		return plotID, nil
	}
	query := fmt.Sprintf(
		"UPDATE plots SET %s WHERE plot_id = $%d AND owner_id = $%d AND is_deleted = FALSE RETURNING plot_id",
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

func (r *Repository) ResizePlot(ctx context.Context, plotID, ownerID string, req models.ResizePlotModel) (string, error) {
	query := `
        UPDATE plots SET
            boundary = ST_MakeEnvelope(0, 0, $1, $2, 3857),
            area_sq_m = $3,
            grid_cell_size = $4,
            grid_cols = $5,
            grid_rows = $6
        WHERE plot_id = $7 AND owner_id = $8 AND is_deleted = FALSE
        RETURNING plot_id`

	var updatedID string
	err := r.db.QueryRow(ctx, query,
		req.BoundaryWidth, req.BoundaryHeight,
		req.AreaSqM, req.GridCellSize, req.GridCols, req.GridRows,
		plotID, ownerID,
	).Scan(&updatedID)
	if err != nil {
		return "", r.mapper.Map(err)
	}
	return updatedID, nil
}

func (r *Repository) SoftDeletePlot(ctx context.Context, plotID, ownerID string) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE plots SET is_deleted = TRUE WHERE plot_id = $1 AND owner_id = $2 AND is_deleted = FALSE`,
		plotID, ownerID)
	if err != nil {
		return r.mapper.Map(err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}
