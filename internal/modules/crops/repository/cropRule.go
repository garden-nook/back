package repository

import (
	"context"
	"garden-nook/internal/modules/crops/dto"
	"garden-nook/internal/modules/crops/enum"
	"garden-nook/internal/modules/crops/model"
	"garden-nook/internal/pkg/database"

	"garden-nook/internal/pkg/apperrors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CropRuleRepo struct {
	db     database.DBTX
	mapper *database.ErrorMapper
}

func NewCropRuleRepo(pool *pgxpool.Pool, mapper *database.ErrorMapper) *CropRuleRepo {
	return &CropRuleRepo{db: pool, mapper: mapper}
}

func (r *CropRuleRepo) WithTx(tx pgx.Tx) *CropRuleRepo {
	return &CropRuleRepo{db: tx, mapper: r.mapper}
}

func (r *CropRuleRepo) ListRules(ctx context.Context, p *database.Pagination) ([]model.CropRule, int, error) {
	baseQuery := `SELECT rule_id, subject_crop_id, subject_family_id, context_type,
		        context_crop_id, context_family_id, return_after_days,
		        score_modifier, explanation, priority
		 FROM crop_rules`
	orderBy := "priority DESC, rule_id ASC"

	plots, total, err := database.ListQuery[model.CropRule](ctx, r.db, baseQuery, "", []any{}, orderBy, p)
	if err != nil {
		return nil, 0, r.mapper.Map(err)
	}
	return plots, total, nil
}

func (r *CropRuleRepo) CreateRule(ctx context.Context, req dto.CreateRuleRequest) (int32, error) {
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

func (r *CropRuleRepo) DeleteRule(ctx context.Context, id int32) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM crop_rules WHERE rule_id=$1`, id)
	if err != nil {
		return r.mapper.Map(err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *CropRuleRepo) GetCropRelations(ctx context.Context, cropID int32) (*model.CropRelations, error) {
	var familyID int32
	err := r.db.QueryRow(ctx, `SELECT family_id FROM crops WHERE id=$1`, cropID).Scan(&familyID)
	if err != nil {
		return nil, r.mapper.Map(err)
	}

	query := `
		SELECT cr.context_type, cr.score_modifier,
               cr.context_crop_id, c.name AS crop_name,
               cr.context_family_id, cf.name AS family_name
        FROM crop_rules cr
        LEFT JOIN crops c ON cr.context_crop_id = c.id
        LEFT JOIN crop_families cf ON cr.context_family_id = cf.id
        WHERE (cr.subject_crop_id = $1 OR cr.subject_family_id = $2)
          AND cr.context_type IN ($3, $4)
	`

	rows, err := r.db.Query(ctx, query,
		cropID, familyID,
		enum.RuleContextPredecessor, enum.RuleContextCompanion,
	)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	defer rows.Close()

	result := &model.CropRelations{}
	for rows.Next() {
		var (
			contextType enum.RuleContextType
			score       int32
			cropIDPtr   *int32
			cropName    *string
			familyIDPtr *int32
			familyName  *string
		)
		if err = rows.Scan(&contextType, &score,
			&cropIDPtr, &cropName, &familyIDPtr, &familyName); err != nil {
			return nil, r.mapper.Map(err)
		}

		switch {
		case cropIDPtr != nil:
			rel := model.CropRelation{
				CropID:   *cropIDPtr,
				CropName: *cropName,
				Score:    score,
			}
			appendCropRelation(result, contextType, rel)
		case familyIDPtr != nil:
			rel := model.FamilyRelation{
				FamilyID:   *familyIDPtr,
				FamilyName: *familyName,
				Score:      score,
			}
			appendFamilyRelation(result, contextType, rel)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, r.mapper.Map(err)
	}

	return result, nil
}

func appendCropRelation(r *model.CropRelations, ctxType enum.RuleContextType, rel model.CropRelation) {
	switch ctxType {
	case enum.RuleContextPredecessor:
		if rel.Score > 0 {
			r.GoodPredecessors = append(r.GoodPredecessors, rel)
		} else if rel.Score < 0 {
			r.BadPredecessors = append(r.BadPredecessors, rel)
		}
	//case RuleContextSuccessor:
	//	if rel.Score > 0 {
	//		r.GoodSuccessors = append(r.GoodSuccessors, rel)
	//	} else if rel.Score < 0 {
	//		r.BadSuccessors = append(r.BadSuccessors, rel)
	//	}
	case enum.RuleContextCompanion:
		if rel.Score > 0 {
			r.GoodCompanions = append(r.GoodCompanions, rel)
		} else if rel.Score < 0 {
			r.BadCompanions = append(r.BadCompanions, rel)
		}
	}
}

func appendFamilyRelation(r *model.CropRelations, ctxType enum.RuleContextType, rel model.FamilyRelation) {
	switch ctxType {
	case enum.RuleContextPredecessor:
		if rel.Score > 0 {
			r.GoodPredecessorFamilies = append(r.GoodPredecessorFamilies, rel)
		} else if rel.Score < 0 {
			r.BadPredecessorFamilies = append(r.BadPredecessorFamilies, rel)
		}
	//case RuleContextSuccessor:
	//	if rel.Score > 0 {
	//		r.GoodSuccessorFamilies = append(r.GoodSuccessorFamilies, rel)
	//	} else if rel.Score < 0 {
	//		r.BadSuccessorFamilies = append(r.BadSuccessorFamilies, rel)
	//	}
	case enum.RuleContextCompanion:
		if rel.Score > 0 {
			r.GoodCompanionFamilies = append(r.GoodCompanionFamilies, rel)
		} else if rel.Score < 0 {
			r.BadCompanionFamilies = append(r.BadCompanionFamilies, rel)
		}
	}
}

func (r *CropRuleRepo) ListRulesFull(ctx context.Context) ([]model.RuleInfo, error) {
	query := `SELECT subject_crop_id, subject_family_id, context_crop_id, context_family_id,
                     return_after_days, score_modifier, explanation
              FROM crop_rules WHERE context_type = $1`
	rows, err := r.db.Query(ctx, query, enum.RuleContextPredecessor)
	if err != nil {
		return nil, r.mapper.Map(err)
	}
	defer rows.Close()

	var out []model.RuleInfo
	for rows.Next() {
		var ri model.RuleInfo
		if err := rows.Scan(&ri.SubjectCropID, &ri.SubjectFamilyID, &ri.ContextCropID, &ri.ContextFamilyID,
			&ri.ReturnAfterDays, &ri.ScoreModifier, &ri.Explanation); err != nil {
			return nil, r.mapper.Map(err)
		}
		out = append(out, ri)
	}
	return out, rows.Err()
}
