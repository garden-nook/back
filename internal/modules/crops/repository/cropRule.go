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
               cr.subject_crop_id,  sub_c.name AS subject_crop_name,
               cr.subject_family_id, sub_f.name AS subject_family_name,
               cr.context_crop_id,  ctx_c.name AS context_crop_name,
               cr.context_family_id, ctx_f.name AS context_family_name
        FROM crop_rules cr
        LEFT JOIN crops sub_c ON cr.subject_crop_id = sub_c.id
        LEFT JOIN crop_families sub_f ON cr.subject_family_id = sub_f.id
        LEFT JOIN crops ctx_c ON cr.context_crop_id = ctx_c.id
        LEFT JOIN crop_families ctx_f ON cr.context_family_id = ctx_f.id
        WHERE 
            (
                (cr.subject_crop_id = $1 OR cr.subject_family_id = $2)
                AND cr.context_type IN ($3, $4)
            )
            OR (
                cr.context_type = $3
                AND (cr.context_crop_id = $1 OR cr.context_family_id = $2)
            )
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
			contextType                          enum.RuleContextType
			score                                int32
			subjectCropID, contextCropID         *int32
			subjectCropName, contextCropName     *string
			subjectFamilyID, contextFamilyID     *int32
			subjectFamilyName, contextFamilyName *string
		)
		if err := rows.Scan(&contextType, &score,
			&subjectCropID, &subjectCropName, &subjectFamilyID, &subjectFamilyName,
			&contextCropID, &contextCropName, &contextFamilyID, &contextFamilyName); err != nil {
			return nil, r.mapper.Map(err)
		}

		switch {
		case contextType == 1 && isRelated(cropID, familyID, subjectCropID, subjectFamilyID):
			addCropOrFamily(result, false, score, contextCropID, contextCropName, contextFamilyID, contextFamilyName)
		case contextType == 1 && isRelated(cropID, familyID, contextCropID, contextFamilyID):
			addCropOrFamily(result, true, score, subjectCropID, subjectCropName, subjectFamilyID, subjectFamilyName)
		case contextType == 3 && isRelated(cropID, familyID, subjectCropID, subjectFamilyID):
			addCompanionCropOrFamily(result, score, contextCropID, contextCropName, contextFamilyID, contextFamilyName)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, r.mapper.Map(err)
	}

	return result, nil
}

func isRelated(cropID, familyID int32, relCropID, relFamilyID *int32) bool {
	return (relCropID != nil && *relCropID == cropID) || (relFamilyID != nil && *relFamilyID == familyID)
}

func addCropOrFamily(result *model.CropRelations, isSuccessor bool, score int32,
	cropID *int32, cropName *string, familyID *int32, familyName *string) {
	if cropID != nil {
		rel := model.CropRelation{CropID: *cropID, CropName: *cropName, Score: score}
		if isSuccessor {
			if score > 0 {
				result.GoodSuccessors = append(result.GoodSuccessors, rel)
			} else {
				result.BadSuccessors = append(result.BadSuccessors, rel)
			}
		} else {
			if score > 0 {
				result.GoodPredecessors = append(result.GoodPredecessors, rel)
			} else {
				result.BadPredecessors = append(result.BadPredecessors, rel)
			}
		}
	} else if familyID != nil {
		rel := model.FamilyRelation{FamilyID: *familyID, FamilyName: *familyName, Score: score}
		if isSuccessor {
			if score > 0 {
				result.GoodSuccessorFamilies = append(result.GoodSuccessorFamilies, rel)
			} else {
				result.BadSuccessorFamilies = append(result.BadSuccessorFamilies, rel)
			}
		} else {
			if score > 0 {
				result.GoodPredecessorFamilies = append(result.GoodPredecessorFamilies, rel)
			} else {
				result.BadPredecessorFamilies = append(result.BadPredecessorFamilies, rel)
			}
		}
	}
}

func addCompanionCropOrFamily(result *model.CropRelations, score int32,
	cropID *int32, cropName *string, familyID *int32, familyName *string) {
	if cropID != nil {
		rel := model.CropRelation{CropID: *cropID, CropName: *cropName, Score: score}
		if score > 0 {
			result.GoodCompanions = append(result.GoodCompanions, rel)
		} else {
			result.BadCompanions = append(result.BadCompanions, rel)
		}
	} else if familyID != nil {
		rel := model.FamilyRelation{FamilyID: *familyID, FamilyName: *familyName, Score: score}
		if score > 0 {
			result.GoodCompanionFamilies = append(result.GoodCompanionFamilies, rel)
		} else {
			result.BadCompanionFamilies = append(result.BadCompanionFamilies, rel)
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
