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

type CropFamilyRepo struct {
	db     database.DBTX
	mapper *database.ErrorMapper
}

func NewCropFamilyRepo(pool *pgxpool.Pool, mapper *database.ErrorMapper) *CropFamilyRepo {
	return &CropFamilyRepo{db: pool, mapper: mapper}
}

func (r *CropFamilyRepo) WithTx(tx pgx.Tx) *CropFamilyRepo {
	return &CropFamilyRepo{db: tx, mapper: r.mapper}
}

func (r *CropFamilyRepo) ListFamilies(ctx context.Context, p *database.Pagination) ([]model.CropFamily, int, error) {
	baseQuery := `SELECT id, name, description FROM crop_families`
	orderBy := "name"

	plots, total, err := database.ListQuery[model.CropFamily](ctx, r.db, baseQuery, "", []any{}, orderBy, p)
	if err != nil {
		return nil, 0, r.mapper.Map(err)
	}
	return plots, total, nil
}

func (r *CropFamilyRepo) GetFamilyByID(ctx context.Context, id int32) (*model.CropFamily, error) {
	row, err := r.db.Query(ctx, `SELECT id, name, description FROM crop_families WHERE id=$1`, id)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	family, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[model.CropFamily])
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	return family, nil
}

func (r *CropFamilyRepo) CreateFamily(ctx context.Context, req dto.CreateFamilyRequest) (int32, error) {
	var id int32
	err := r.db.QueryRow(ctx,
		`INSERT INTO crop_families(name, description) VALUES ($1,$2) RETURNING id`,
		req.Name, req.Description,
	).Scan(&id)
	if err != nil {
		return 0, r.mapper.Map(err)
	}
	return id, nil
}

func (r *CropFamilyRepo) UpdateFamily(ctx context.Context, id int32, req dto.UpdateFamilyRequest) (int32, error) {
	fields := []database.SetField{
		database.NewSetField("name", req.Name),
		database.NewSetField("description", req.Description),
	}
	setSQL, setArgs := database.BuildUpdateSet(1, fields...)
	if len(setArgs) == 0 {
		return id, nil
	}
	query := fmt.Sprintf("UPDATE crop_families SET %s WHERE id=$%d RETURNING id",
		setSQL, len(setArgs)+1)
	args := append(setArgs, id)

	var updatedID int32
	err := r.db.QueryRow(ctx, query, args...).Scan(&updatedID)
	if err != nil {
		return 0, r.mapper.Map(err)
	}
	return updatedID, nil
}

func (r *CropFamilyRepo) DeleteFamily(ctx context.Context, id int32) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM crop_families WHERE id=$1`, id)
	if err != nil {
		return r.mapper.Map(err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}
