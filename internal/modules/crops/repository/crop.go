package repository

import (
	"context"
	"fmt"
	"garden-nook/internal/modules/crops/dto"
	"garden-nook/internal/modules/crops/model"
	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/database"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CropRepo struct {
	db     database.DBTX
	mapper *database.ErrorMapper
}

func NewCropRepo(pool *pgxpool.Pool, mapper *database.ErrorMapper) *CropRepo {
	return &CropRepo{db: pool, mapper: mapper}
}

func (r *CropRepo) WithTx(tx pgx.Tx) *CropRepo {
	return &CropRepo{db: tx, mapper: r.mapper}
}

func (r *CropRepo) ListCrops(ctx context.Context, f dto.ListCropsFilter, p *database.Pagination) ([]model.Crop, int, error) {
	baseQuery := `SELECT c.id, c.name, c.description, c.family_id, cf.name as family_name,
	               c.vegetation_days_avg, c.sun_needs, c.soil_type_id, st.name as soil_name
	        FROM crops c
	        JOIN crop_families cf ON cf.id = c.family_id
	        JOIN soil_types st ON st.id = c.soil_type_id`
	whereClause := "WHERE c.is_deleted = FALSE"
	whereArgs := []any{}
	argIdx := 1
	if f.FamilyID != nil {
		whereClause += fmt.Sprintf(" AND c.family_id = $%d", argIdx)
		whereArgs = append(whereArgs, *f.FamilyID)
		argIdx++
	}
	if f.SoilTypeID != nil {
		whereClause += fmt.Sprintf(" AND c.soil_type_id = $%d", argIdx)
		whereArgs = append(whereArgs, *f.SoilTypeID)
		argIdx++
	}
	if f.Search != "" {
		whereClause += fmt.Sprintf(" AND c.name ILIKE $%d", argIdx)
		whereArgs = append(whereArgs, "%"+f.Search+"%")
		argIdx++
	}

	orderBy := "c.name ASC"

	plots, total, err := database.ListQuery[model.Crop](ctx, r.db, baseQuery, whereClause, whereArgs, orderBy, p)
	if err != nil {
		return nil, 0, r.mapper.Map(err)
	}
	return plots, total, nil
}

func (r *CropRepo) GetCropByID(ctx context.Context, id int32) (*model.Crop, error) {
	row, err := r.db.Query(ctx,
		`SELECT c.id, c.name, c.description, c.family_id, cf.name as family_name, 
			 c.vegetation_days_avg, c.sun_needs, c.soil_type_id, st.name as soil_name
		 FROM crops c 
		 JOIN crop_families cf ON cf.id = c.family_id
	     JOIN soil_types st ON st.id = c.soil_type_id
		 WHERE c.id=$1 AND c.is_deleted=FALSE`,
		id,
	)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	crop, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[model.Crop])
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	return crop, nil
}

func (r *CropRepo) CreateCrop(ctx context.Context, req dto.CreateCropRequest) (int32, error) {
	var id int32
	err := r.db.QueryRow(ctx,
		`INSERT INTO crops(name, family_id, soil_type_id, vegetation_days_avg, sun_needs) 
		 VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		req.Name, req.FamilyID, req.SoilTypeID, req.VegetationDaysAvg, req.SunNeeds,
	).Scan(&id)
	if err != nil {
		return 0, r.mapper.Map(err)
	}
	return id, nil
}

func (r *CropRepo) UpdateCrop(ctx context.Context, id int32, req dto.UpdateCropRequest) (int32, error) {
	fields := []database.SetField{
		database.NewSetField("name", req.Name),
		database.NewSetField("family_id", req.FamilyID),
		database.NewSetField("vegetation_days_avg", req.VegetationDaysAvg),
		database.NewSetField("sun_needs", req.SunNeeds),
	}
	setSQL, setArgs := database.BuildUpdateSet(1, fields...)
	if len(setArgs) == 0 {
		return id, nil
	}
	query := fmt.Sprintf(
		"UPDATE crops SET %s WHERE id=$%d AND is_deleted=FALSE RETURNING id",
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

func (r *CropRepo) SoftDeleteCrop(ctx context.Context, id int32) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE crops SET is_deleted = TRUE WHERE id=$1 AND is_deleted=FALSE`, id)
	if err != nil {
		return r.mapper.Map(err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *CropRepo) ListCropsFiltered(ctx context.Context, soilTypeID *int32, sunNeeds []int32, search string, limit int) ([]model.Crop, error) {
	where := `WHERE c.is_deleted = FALSE`
	var args []any
	argIdx := 1

	if soilTypeID != nil {
		where += fmt.Sprintf(" AND c.soil_type_id = $%d", argIdx)
		args = append(args, *soilTypeID)
		argIdx++
	}
	if len(sunNeeds) > 0 {
		placeholders := make([]string, len(sunNeeds))
		for i := range sunNeeds {
			placeholders[i] = fmt.Sprintf("$%d", argIdx)
			args = append(args, sunNeeds[i])
			argIdx++
		}
		where += fmt.Sprintf(" AND c.sun_needs IN (%s)", strings.Join(placeholders, ","))
	}
	if search != "" {
		where += fmt.Sprintf(" AND c.name ILIKE $%d", argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query := fmt.Sprintf(`
        SELECT c.id, c.name, c.family_id, cf.name AS family_name, c.sun_needs, c.soil_type_id
        FROM crops c
        JOIN crop_families cf ON cf.id = c.family_id
        %s
        ORDER BY c.name
        LIMIT %d
    `, where, limit)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	defer rows.Close()

	var crops []model.Crop
	for rows.Next() {
		var c model.Crop
		if err := rows.Scan(&c.ID, &c.Name, &c.FamilyID, &c.FamilyName, &c.SunNeeds, &c.SoilTypeID); err != nil {
			return nil, r.mapper.Map(err)
		}
		crops = append(crops, c)
	}
	return crops, rows.Err()
}
