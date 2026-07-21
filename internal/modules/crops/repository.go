package crops

import (
	"context"
	"fmt"
	"garden-nook/internal/pkg/database"

	"garden-nook/internal/pkg/apperrors"

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

// WithTx возвращает новый экземпляр репозитория, привязанный к транзакции.
func (r *Repository) WithTx(tx pgx.Tx) *Repository {
	return &Repository{db: tx, mapper: r.mapper}
}

//  ---------- SOIL TYPES ----------

func (r *Repository) ListSoilTypes(ctx context.Context, p *database.Pagination) ([]SoilType, int, error) {
	base := `SELECT id, name, description FROM soil_types`
	order := ` ORDER BY name`

	if p == nil {
		rows, err := r.db.Query(ctx, base+order)
		if err != nil {
			return nil, 0, r.mapper.Map(err)
		}
		defer rows.Close()
		soilTypes, err := pgx.CollectRows(rows, pgx.RowToStructByName[SoilType])
		if err != nil {
			return nil, 0, r.mapper.Map(err)
		}
		return soilTypes, len(soilTypes), nil
	}

	batch := &pgx.Batch{}
	batch.Queue(`SELECT COUNT(*) FROM soil_types`)

	pagSQL, pagArgs := p.SQL(1)
	batch.Queue(base+order+pagSQL, pagArgs...)

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
	soilTypes, err := pgx.CollectRows(rows, pgx.RowToStructByName[SoilType])
	if err != nil {
		return nil, 0, r.mapper.Map(err)
	}
	return soilTypes, total, nil
}

// GetSoilTypeByID возвращает тип почвы по ID.
func (r *Repository) GetSoilTypeByID(ctx context.Context, id int32) (*SoilType, error) {
	row, err := r.db.Query(ctx,
		`SELECT id, name, description FROM soil_types WHERE id = $1`, id,
	)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	soilType, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[SoilType])
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	return soilType, nil
}

// CreateSoilType создаёт новый тип почвы и возвращает его ID.
func (r *Repository) CreateSoilType(ctx context.Context, req CreateSoilTypeRequest) (int32, error) {
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

// UpdateSoilType обновляет тип почвы и возвращает его ID.
func (r *Repository) UpdateSoilType(ctx context.Context, id int32, req UpdateSoilTypeRequest) (int32, error) {
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

// DeleteSoilType удаляет тип почвы по ID.
func (r *Repository) DeleteSoilType(ctx context.Context, id int32) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM soil_types WHERE id = $1`, id)
	if err != nil {
		return r.mapper.Map(err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

// ---------- FAMILIES ----------

func (r *Repository) ListFamilies(ctx context.Context, p *database.Pagination) ([]CropFamily, int, error) {
	base := `SELECT id, name, description FROM crop_families`
	order := ` ORDER BY name`

	if p == nil {
		rows, err := r.db.Query(ctx, base+order)
		if err != nil {
			return nil, 0, r.mapper.Map(err)
		}
		defer rows.Close()
		families, err := pgx.CollectRows(rows, pgx.RowToStructByName[CropFamily])
		if err != nil {
			return nil, 0, r.mapper.Map(err)
		}
		return families, len(families), nil
	}

	batch := &pgx.Batch{}
	batch.Queue(`SELECT COUNT(*) FROM crop_families`)

	pagSQL, pagArgs := p.SQL(1)
	batch.Queue(base+order+pagSQL, pagArgs...)

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
	families, err := pgx.CollectRows(rows, pgx.RowToStructByName[CropFamily])
	if err != nil {
		return nil, 0, r.mapper.Map(err)
	}
	return families, total, nil
}

func (r *Repository) GetFamilyByID(ctx context.Context, id int32) (*CropFamily, error) {
	row, err := r.db.Query(ctx, `SELECT id, name, description FROM crop_families WHERE id=$1`, id)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	family, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[CropFamily])
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	return family, nil
}

func (r *Repository) CreateFamily(ctx context.Context, req CreateFamilyRequest) (int32, error) {
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

func (r *Repository) UpdateFamily(ctx context.Context, id int32, req UpdateFamilyRequest) (int32, error) {
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

func (r *Repository) DeleteFamily(ctx context.Context, id int32) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM crop_families WHERE id=$1`, id)
	if err != nil {
		return r.mapper.Map(err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

// ---------- CROPS ----------

func (r *Repository) ListCrops(ctx context.Context, f ListCropsFilter, p *database.Pagination) ([]Crop, int, error) {
	where := `WHERE c.is_deleted = FALSE`
	var filterArgs []any
	argIdx := 1

	if f.FamilyID != nil {
		where += fmt.Sprintf(" AND c.family_id = $%d", argIdx)
		filterArgs = append(filterArgs, *f.FamilyID)
		argIdx++
	}
	if f.SoilTypeID != nil {
		where += fmt.Sprintf(" AND c.soil_type_id = $%d", argIdx)
		filterArgs = append(filterArgs, *f.SoilTypeID)
		argIdx++
	}
	if f.Search != "" {
		where += fmt.Sprintf(" AND LOWER(c.name) LIKE $%d", argIdx)
		filterArgs = append(filterArgs, "%"+f.Search+"%")
		argIdx++
	}

	baseSelect := `SELECT c.id, c.name, c.description, c.family_id, cf.name as family_name,
	               c.vegetation_days_avg, c.sun_needs, c.soil_type_id, st.name as soil_name
	        FROM crops c
	        JOIN crop_families cf ON cf.id = c.family_id
	        JOIN soil_types st ON st.id = c.soil_type_id `

	if p == nil {
		rows, err := r.db.Query(ctx, baseSelect+where+` ORDER BY c.name ASC`, filterArgs)
		if err != nil {
			return nil, 0, r.mapper.Map(err)
		}
		defer rows.Close()
		crops, err := pgx.CollectRows(rows, pgx.RowToStructByName[Crop])
		if err != nil {
			return nil, 0, r.mapper.Map(err)
		}
		return crops, len(crops), nil
	}

	batch := &pgx.Batch{}
	countSQL := `SELECT COUNT(*) FROM crops c ` + where
	batch.Queue(countSQL, filterArgs...)

	pagSQL, pagArgs := p.SQL(argIdx)
	dataSQL := baseSelect + where + ` ORDER BY c.name ASC` + pagSQL
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
	crops, err := pgx.CollectRows(rows, pgx.RowToStructByName[Crop])
	if err != nil {
		return nil, 0, r.mapper.Map(err)
	}
	return crops, total, nil
}

func (r *Repository) GetCropByID(ctx context.Context, id int32) (*Crop, error) {
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
	crop, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[Crop])
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	return crop, nil
}

func (r *Repository) CreateCrop(ctx context.Context, req CreateCropRequest) (int32, error) {
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

func (r *Repository) UpdateCrop(ctx context.Context, id int32, req UpdateCropRequest) (int32, error) {
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

func (r *Repository) SoftDeleteCrop(ctx context.Context, id int32) error {
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

// ---------- RULES ----------

func (r *Repository) ListRules(ctx context.Context, p *database.Pagination) ([]CropRule, int, error) {
	base := `SELECT rule_id, subject_crop_id, subject_family_id, context_type,
		        context_crop_id, context_family_id, return_after_days,
		        score_modifier, explanation, priority
		 FROM crop_rules`
	order := ` ORDER BY priority DESC, rule_id ASC`

	if p == nil {
		rows, err := r.db.Query(ctx, base+order)
		if err != nil {
			return nil, 0, r.mapper.Map(err)
		}
		defer rows.Close()
		rules, err := pgx.CollectRows(rows, pgx.RowToStructByName[CropRule])
		if err != nil {
			return nil, 0, r.mapper.Map(err)
		}
		return rules, len(rules), nil
	}

	batch := &pgx.Batch{}
	batch.Queue(`SELECT COUNT(*) FROM crop_rules`)

	pagSQL, pagArgs := p.SQL(1)
	batch.Queue(base+order+pagSQL, pagArgs...)

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
	rules, err := pgx.CollectRows(rows, pgx.RowToStructByName[CropRule])
	if err != nil {
		return nil, 0, r.mapper.Map(err)
	}
	return rules, total, nil
}

func (r *Repository) CreateRule(ctx context.Context, req CreateRuleRequest) (int32, error) {
	var id int32
	err := r.db.QueryRow(ctx,
		`INSERT INTO crop_rules
		  (subject_crop_id, subject_family_id, context_type, context_crop_id, context_family_id,
		   return_after_days, score_modifier, explanation, priority)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		 RETURNING rule_id`,
		req.SubjectCropID, req.SubjectFamilyID, req.ContextType, req.ContextCropID, req.ContextFamilyID,
		req.ReturnAfterDays, req.ScoreModifier, req.Explanation, req.Priority,
	).Scan(&id)
	if err != nil {
		return 0, r.mapper.Map(err)
	}
	return id, nil
}

func (r *Repository) DeleteRule(ctx context.Context, id int32) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM crop_rules WHERE rule_id=$1`, id)
	if err != nil {
		return r.mapper.Map(err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

// GetCropRelations возвращает все связи для культуры с заданным ID.
func (r *Repository) GetCropRelations(ctx context.Context, cropID int32) (*CropRelations, error) {
	// Получаем family_id культуры
	var familyID int32
	err := r.db.QueryRow(ctx, `SELECT family_id FROM crops WHERE id=$1`, cropID).Scan(&familyID)
	if err != nil {
		return nil, r.mapper.Map(err)
	}

	// Запрос, возвращающий все культуры, связанные с заданной через правила.
	query := `
		SELECT cr.context_type, cr.score_modifier,
               cr.context_crop_id, c.name AS crop_name,
               cr.context_family_id, cf.name AS family_name
        FROM crop_rules cr
        LEFT JOIN crops c ON cr.context_crop_id = c.id
        LEFT JOIN crop_families cf ON cr.context_family_id = cf.id
        WHERE (cr.subject_crop_id = $1 OR cr.subject_family_id = $2)
          AND cr.context_type IN ($3, $4, $5)
	`

	rows, err := r.db.Query(ctx, query,
		cropID, familyID,
		RuleContextPredecessor, RuleContextSuccessor, RuleContextCompanion,
	)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	defer rows.Close()

	result := &CropRelations{}
	for rows.Next() {
		var (
			contextType int32
			score       int32
			//explanation string
			cropIDPtr   *int32
			cropName    *string
			familyIDPtr *int32
			familyName  *string
		)
		if err = rows.Scan(&contextType, &score, /* &explanation,*/
			&cropIDPtr, &cropName, &familyIDPtr, &familyName); err != nil {
			return nil, r.mapper.Map(err)
		}

		switch {
		case cropIDPtr != nil:
			rel := CropRelation{
				CropID:   *cropIDPtr,
				CropName: *cropName,
				Score:    score,
				//Explanation: explanation,
			}
			appendCropRelation(result, contextType, rel)
		case familyIDPtr != nil:
			rel := FamilyRelation{
				FamilyID:   *familyIDPtr,
				FamilyName: *familyName,
				Score:      score,
				//Explanation: explanation,
			}
			appendFamilyRelation(result, contextType, rel)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, r.mapper.Map(err)
	}

	return result, nil
}

func appendCropRelation(r *CropRelations, ctxType int32, rel CropRelation) {
	switch ctxType {
	case RuleContextPredecessor:
		if rel.Score > 0 {
			r.GoodPredecessors = append(r.GoodPredecessors, rel)
		} else if rel.Score < 0 {
			r.BadPredecessors = append(r.BadPredecessors, rel)
		}
	case RuleContextSuccessor:
		if rel.Score > 0 {
			r.GoodSuccessors = append(r.GoodSuccessors, rel)
		} else if rel.Score < 0 {
			r.BadSuccessors = append(r.BadSuccessors, rel)
		}
	case RuleContextCompanion:
		if rel.Score > 0 {
			r.GoodCompanions = append(r.GoodCompanions, rel)
		} else if rel.Score < 0 {
			r.BadCompanions = append(r.BadCompanions, rel)
		}
	}
}

func appendFamilyRelation(r *CropRelations, ctxType int32, rel FamilyRelation) {
	switch ctxType {
	case RuleContextPredecessor:
		if rel.Score > 0 {
			r.GoodPredecessorFamilies = append(r.GoodPredecessorFamilies, rel)
		} else if rel.Score < 0 {
			r.BadPredecessorFamilies = append(r.BadPredecessorFamilies, rel)
		}
	case RuleContextSuccessor:
		if rel.Score > 0 {
			r.GoodSuccessorFamilies = append(r.GoodSuccessorFamilies, rel)
		} else if rel.Score < 0 {
			r.BadSuccessorFamilies = append(r.BadSuccessorFamilies, rel)
		}
	case RuleContextCompanion:
		if rel.Score > 0 {
			r.GoodCompanionFamilies = append(r.GoodCompanionFamilies, rel)
		} else if rel.Score < 0 {
			r.BadCompanionFamilies = append(r.BadCompanionFamilies, rel)
		}
	}
}
