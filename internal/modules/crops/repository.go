package crops

import (
	"context"
	"errors"

	"garden-nook/internal/pkg/apperrors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// ---------- FAMILIES ----------

func (r *Repository) ListFamilies(ctx context.Context) ([]CropFamily, error) {
	query := `SELECT id, name, COALESCE(description,'') FROM crop_families ORDER BY name`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []CropFamily
	for rows.Next() {
		var f CropFamily
		if err := rows.Scan(&f.ID, &f.Name, &f.Description); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func (r *Repository) GetFamilyByID(ctx context.Context, id int32) (*CropFamily, error) {
	f := &CropFamily{}
	err := r.db.QueryRow(ctx,
		`SELECT id, name, COALESCE(description,'') FROM crop_families WHERE id=$1`, id,
	).Scan(&f.ID, &f.Name, &f.Description)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (r *Repository) CreateFamily(ctx context.Context, req CreateFamilyRequest) (*CropFamily, error) {
	f := &CropFamily{}
	err := r.db.QueryRow(ctx,
		`INSERT INTO crop_families(name, description) VALUES ($1, $2) RETURNING id, name, COALESCE(description,'')`,
		req.Name, req.Description,
	).Scan(&f.ID, &f.Name, &f.Description)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return nil, apperrors.ErrConflict
		}
		return nil, err
	}
	return f, nil
}

func (r *Repository) UpdateFamily(ctx context.Context, id int32, req UpdateFamilyRequest) (*CropFamily, error) {
	existing, err := r.GetFamilyByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}

	_, err = r.db.Exec(ctx,
		`UPDATE crop_families SET name=$1, description=$2 WHERE id=$3`,
		existing.Name, existing.Description, id,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, apperrors.ErrConflict
		}
		return nil, err
	}
	return existing, nil
}

func (r *Repository) DeleteFamily(ctx context.Context, id int32) error {
	// Физическое удаление семейства допустимо, если нет связанных crops.
	// FK с ON DELETE RESTRICT в БД защитит нас.
	tag, err := r.db.Exec(ctx, `DELETE FROM crop_families WHERE id=$1`, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" { // foreign_key_violation
			return apperrors.ErrConflict
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

// ---------- CROPS ----------

func (r *Repository) ListCrops(ctx context.Context, f ListCropsFilter) ([]Crop, int, error) {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 50
	}
	if f.Page <= 0 {
		f.Page = 1
	}
	offset := (f.Page - 1) * f.Limit

	where := `WHERE c.is_deleted = FALSE`
	args := []interface{}{}
	argIdx := 1

	if f.FamilyID != nil {
		where += " AND c.family_id = $" + itoa(argIdx)
		args = append(args, *f.FamilyID)
		argIdx++
	}
	if f.Search != "" {
		where += " AND LOWER(c.name) LIKE $" + itoa(argIdx)
		args = append(args, "%"+f.Search+"%")
		argIdx++
	}

	var total int
	countSQL := `SELECT COUNT(*) FROM crops c ` + where
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	dataSQL := `SELECT c.id, c.name, c.family_id, cf.name, c.vegetation_days_avg, c.sun_needs
	            FROM crops c
	            JOIN crop_families cf ON cf.id = c.family_id ` +
		where + ` ORDER BY c.name ASC LIMIT $` + itoa(argIdx) + ` OFFSET $` + itoa(argIdx+1)
	args = append(args, f.Limit, offset)

	rows, err := r.db.Query(ctx, dataSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []Crop
	for rows.Next() {
		var c Crop
		if err := rows.Scan(&c.ID, &c.Name, &c.FamilyID, &c.FamilyName, &c.VegetationDaysAvg, &c.SunNeeds); err != nil {
			return nil, 0, err
		}
		out = append(out, c)
	}
	return out, total, rows.Err()
}

func (r *Repository) GetCropByID(ctx context.Context, id int32) (*Crop, error) {
	c := &Crop{}
	err := r.db.QueryRow(ctx,
		`SELECT c.id, c.name, c.family_id, cf.name, c.vegetation_days_avg, c.sun_needs
		 FROM crops c JOIN crop_families cf ON cf.id = c.family_id
		 WHERE c.id=$1 AND c.is_deleted=FALSE`, id,
	).Scan(&c.ID, &c.Name, &c.FamilyID, &c.FamilyName, &c.VegetationDaysAvg, &c.SunNeeds)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.ErrNotFound
	}
	return c, err
}

func (r *Repository) CreateCrop(ctx context.Context, req CreateCropRequest) (*Crop, error) {
	c := &Crop{}
	err := r.db.QueryRow(ctx,
		`INSERT INTO crops(name, family_id, vegetation_days_avg, sun_needs)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, name, family_id, vegetation_days_avg, sun_needs`,
		req.Name, req.FamilyID, req.VegetationDaysAvg, req.SunNeeds,
	).Scan(&c.ID, &c.Name, &c.FamilyID, &c.VegetationDaysAvg, &c.SunNeeds)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" { // FK violation (family not exists)
			return nil, apperrors.ErrBadRequest
		}
		return nil, err
	}
	// Подтягиваем имя семейства для ответа
	fam, _ := r.GetFamilyByID(ctx, c.FamilyID)
	if fam != nil {
		c.FamilyName = fam.Name
	}
	return c, nil
}

func (r *Repository) UpdateCrop(ctx context.Context, id int32, req UpdateCropRequest) (*Crop, error) {
	existing, err := r.GetCropByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.FamilyID != nil {
		existing.FamilyID = *req.FamilyID
	}
	if req.VegetationDaysAvg != nil {
		existing.VegetationDaysAvg = *req.VegetationDaysAvg
	}
	if req.SunNeeds != nil {
		existing.SunNeeds = *req.SunNeeds
	}

	_, err = r.db.Exec(ctx,
		`UPDATE crops SET name=$1, family_id=$2, vegetation_days_avg=$3, sun_needs=$4 WHERE id=$5`,
		existing.Name, existing.FamilyID, existing.VegetationDaysAvg, existing.SunNeeds, id,
	)
	if err != nil {
		return nil, err
	}

	// Перегружаем с актуальным FamilyName
	return r.GetCropByID(ctx, id)
}

func (r *Repository) SoftDeleteCrop(ctx context.Context, id int32) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE crops SET is_deleted = TRUE WHERE id=$1 AND is_deleted=FALSE`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

// ---------- RULES ----------

func (r *Repository) ListRules(ctx context.Context) ([]CropRule, error) {
	rows, err := r.db.Query(ctx,
		`SELECT rule_id, subject_crop_id, subject_family_id, context_type,
		        context_crop_id, context_family_id, return_after_days,
		        score_modifier, explanation, priority
		 FROM crop_rules ORDER BY priority DESC, rule_id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []CropRule
	for rows.Next() {
		var ru CropRule
		if err := rows.Scan(&ru.ID, &ru.SubjectCropID, &ru.SubjectFamilyID, &ru.ContextType,
			&ru.ContextCropID, &ru.ContextFamilyID, &ru.ReturnAfterDays,
			&ru.ScoreModifier, &ru.Explanation, &ru.Priority); err != nil {
			return nil, err
		}
		out = append(out, ru)
	}
	return out, rows.Err()
}

func (r *Repository) CreateRule(ctx context.Context, req CreateRuleRequest) (*CropRule, error) {
	ru := &CropRule{}
	err := r.db.QueryRow(ctx,
		`INSERT INTO crop_rules
		  (subject_crop_id, subject_family_id, context_type, context_crop_id, context_family_id,
		   return_after_days, score_modifier, explanation, priority)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		 RETURNING rule_id, subject_crop_id, subject_family_id, context_type,
		           context_crop_id, context_family_id, return_after_days,
		           score_modifier, explanation, priority`,
		req.SubjectCropID, req.SubjectFamilyID, req.ContextType, req.ContextCropID, req.ContextFamilyID,
		req.ReturnAfterDays, req.ScoreModifier, req.Explanation, req.Priority,
	).Scan(&ru.ID, &ru.SubjectCropID, &ru.SubjectFamilyID, &ru.ContextType,
		&ru.ContextCropID, &ru.ContextFamilyID, &ru.ReturnAfterDays,
		&ru.ScoreModifier, &ru.Explanation, &ru.Priority)
	if err != nil {
		return nil, err
	}
	return ru, nil
}

func (r *Repository) DeleteRule(ctx context.Context, id int32) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM crop_rules WHERE rule_id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

// itoa — минихелпер, чтобы не тащить strconv ради одной операции.
func itoa(i int) string {
	return string(rune('0'+i%10)) + "" // упрощённо: только для i<10 достаточно
}
