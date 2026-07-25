package repository

import (
	"context"
	"fmt"
	"garden-nook/internal/modules/crops/dto"
	"garden-nook/internal/modules/crops/model"
	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SoilTypeRepo struct {
	db     database.DBTX
	mapper *database.ErrorMapper
}

func NewSoilTypeRepo(pool *pgxpool.Pool, mapper *database.ErrorMapper) *SoilTypeRepo {
	return &SoilTypeRepo{db: pool, mapper: mapper}
}

func (r *SoilTypeRepo) WithTx(tx pgx.Tx) *SoilTypeRepo {
	return &SoilTypeRepo{db: tx, mapper: r.mapper}
}

func (r *SoilTypeRepo) ListSoilTypes(ctx context.Context, p *database.Pagination) ([]model.SoilType, int, error) {
	baseQuery := `SELECT id, name, description FROM soil_types`
	orderBy := "name"

	plots, total, err := database.ListQuery[model.SoilType](ctx, r.db, baseQuery, "", []any{}, orderBy, p)
	if err != nil {
		return nil, 0, r.mapper.Map(err)
	}
	return plots, total, nil
}

func (r *SoilTypeRepo) GetSoilTypeByID(ctx context.Context, id int32) (*model.SoilType, error) {
	row, err := r.db.Query(ctx,
		`SELECT id, name, description FROM soil_types WHERE id = $1`, id,
	)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	soilType, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[model.SoilType])
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	return soilType, nil
}

func (r *SoilTypeRepo) CreateSoilType(ctx context.Context, req dto.CreateSoilTypeRequest) (int32, error) {
	var id int32
	err := r.db.QueryRow(ctx,
		`INSERT INTO soil_types (name, description) VALUES ($1, $2) RETURNING id`,
		req.Name, req.Description,
	).Scan(&id)
	if err != nil {
		return 0, r.mapper.Map(err)
	}
	return id, nil
}

func (r *SoilTypeRepo) UpdateSoilType(ctx context.Context, id int32, req dto.UpdateSoilTypeRequest) (int32, error) {
	fields := []database.SetField{
		database.NewSetField("name", req.Name),
		database.NewSetField("description", req.Description),
	}
	setSQL, setArgs := database.BuildUpdateSet(1, fields...)
	if len(setArgs) == 0 {
		return id, nil
	}
	query := fmt.Sprintf(
		"UPDATE soil_types SET %s WHERE id = $%d RETURNING id",
		setSQL, len(setArgs)+1,
	)
	args := append(setArgs, id)

	var updatedID int32
	err := r.db.QueryRow(ctx, query, args...).Scan(&updatedID)
	if err != nil {
		return 0, r.mapper.Map(err)
	}
	return updatedID, nil
}

func (r *SoilTypeRepo) DeleteSoilType(ctx context.Context, id int32) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM soil_types WHERE id = $1`, id)
	if err != nil {
		return r.mapper.Map(err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}
