package service

import (
	"context"
	"garden-nook/internal/modules/plot/dto"
	"garden-nook/internal/modules/plot/enum"
	"garden-nook/internal/modules/plot/model"
	"garden-nook/internal/modules/plot/provider"
	"garden-nook/internal/modules/plot/repository"
	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/helpers"
	"math"
	"sort"
	"time"
)

const (
	// ShadeTolerance – допустимое отклонение среднего уровня затенения для включения соседних уровней.
	ShadeTolerance = 0.7
)

type RecommendationService struct {
	bedRepo     *repository.BedRepo
	plotRepo    *repository.PlotRepo
	gridRepo    *repository.GridCellRepo
	historyRepo *repository.HistoryRepo
	cropProv    provider.CropProvider
	ruleProv    provider.RuleProvider
	seh         *helpers.ServiceErrorHandler
}

func NewRecommendationService(
	bedRepo *repository.BedRepo,
	plotRepo *repository.PlotRepo,
	gridRepo *repository.GridCellRepo,
	historyRepo *repository.HistoryRepo,
	cropProv provider.CropProvider,
	ruleProv provider.RuleProvider,
	seh *helpers.ServiceErrorHandler,
) *RecommendationService {
	return &RecommendationService{
		bedRepo:     bedRepo,
		plotRepo:    plotRepo,
		gridRepo:    gridRepo,
		historyRepo: historyRepo,
		cropProv:    cropProv,
		ruleProv:    ruleProv,
		seh:         seh,
	}
}

// GetBedRecommendations возвращает рекомендации культур для указанной грядки.
// Если search не пуст, дополнительно возвращает результаты поиска (search_results).
func (s *RecommendationService) GetBedRecommendations(
	ctx context.Context,
	bedID string,
	ownerID string,
	search string,
	limit int,
	searchLimit int,
	disableFilters bool,
) (*dto.BedRecommendationsResponse, error) {
	bed, err := s.bedRepo.GetBedByID(ctx, bedID)
	if err != nil {
		return nil, s.seh.HandleError(err, "get bed")
	}
	if bed.CurrentCropID != nil {
		return nil, apperrors.ErrConflict
	}

	plot, err := s.plotRepo.GetPlotByOwnerAndID(ctx, bed.PlotID, ownerID)
	if err != nil {
		return nil, s.seh.HandleError(err, "check plot ownership")
	}

	allowedSunNeeds, err := s.calculateAllowedSunNeeds(ctx, bed)
	if err != nil {
		return nil, err
	}

	predecessors, err := s.historyRepo.GetPredecessorsForBed(ctx, bed.PlotID, bed.XStart, bed.YStart, bed.Width, bed.Height)
	if err != nil {
		return nil, s.seh.HandleError(err, "get predecessors")
	}

	cropFilter := model.CropFilter{
		Search:     search,
		Limit:      limit,
		SoilTypeID: &plot.SoilTypeID,
		SunNeeds:   allowedSunNeeds,
	}

	candidates, err := s.cropProv.ListCropsFiltered(ctx, cropFilter)
	if err != nil {
		return nil, s.seh.HandleError(err, "list crops")
	}

	var searchDTO []dto.CropSearchResult
	if search != "" {
		searchFilter := model.CropFilter{
			Search: search,
			Limit:  searchLimit,
		}
		if !disableFilters {
			searchFilter.SoilTypeID = &plot.SoilTypeID
			searchFilter.SunNeeds = allowedSunNeeds
		}
		searchResults, err := s.cropProv.ListCropsFiltered(ctx, searchFilter)
		if err != nil {
			return nil, s.seh.HandleError(err, "search crops")
		}
		searchDTO = make([]dto.CropSearchResult, len(searchResults))
		for i, c := range searchResults {
			searchDTO[i] = dto.CropSearchResult{
				CropID:     c.ID,
				Name:       c.Name,
				FamilyName: c.FamilyName,
			}
		}
	}

	now := time.Now().UTC()
	recommendations := make([]dto.CropRecommendation, 0, len(candidates))
	for _, crop := range candidates {
		score := int32(0)
		var reasons []dto.ReasonDetail

		cropRules, err := s.ruleProv.GetRulesBySubjectCropID(ctx, crop.ID)
		if err != nil {
			return nil, s.seh.HandleError(err, "get crop rules")
		}
		familyRules, err := s.ruleProv.GetRulesBySubjectFamilyID(ctx, crop.FamilyID)
		if err != nil {
			return nil, s.seh.HandleError(err, "get family rules")
		}
		allRules := append(cropRules, familyRules...)

		// Проверяем каждое правило против каждого предшественника
		for _, pre := range predecessors {
			daysSinceHarvest := int32(now.Sub(pre.LastHarvestDate).Hours() / 24)
			for _, rule := range allRules {
				// Проверяем совпадение контекста
				matches := false
				if rule.ContextCropID != nil {
					if *rule.ContextCropID == pre.CropID {
						matches = true
					}
				} else if rule.ContextFamilyID != nil {
					if *rule.ContextFamilyID == pre.FamilyID {
						matches = true
					}
				}
				if !matches {
					continue
				}

				// Условия применения правила
				apply := false
				if rule.ScoreModifier > 0 {
					if rule.ReturnAfterDays == 0 || daysSinceHarvest >= rule.ReturnAfterDays {
						apply = true
					}
				} else { // < 0
					if rule.ReturnAfterDays == 0 || daysSinceHarvest < rule.ReturnAfterDays {
						apply = true
					}
				}

				if apply {
					score += rule.ScoreModifier
					reasons = append(reasons, dto.ReasonDetail{
						Explanation: rule.Explanation,
						Score:       rule.ScoreModifier,
						IsPositive:  rule.ScoreModifier >= 0,
					})
				}
			}
		}

		sort.Slice(reasons, func(i, j int) bool {
			return reasons[i].Score > reasons[j].Score
		})

		if len(reasons) != 0 {
			recommendations = append(recommendations, dto.CropRecommendation{
				CropID:     crop.ID,
				Name:       crop.Name,
				FamilyName: crop.FamilyName,
				Score:      score,
				Reasons:    reasons,
				IsPositive: score >= 0,
			})
		}
	}

	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return &dto.BedRecommendationsResponse{
		Recommendations: recommendations,
		SearchResults:   searchDTO,
	}, nil
}

func (s *RecommendationService) calculateAllowedSunNeeds(ctx context.Context, bed *model.Bed) ([]int32, error) {
	avg, err := s.gridRepo.GetAverageShadeForRectangle(ctx, bed.PlotID, bed.XStart, bed.YStart, bed.Width, bed.Height)
	if err != nil {
		return nil, s.seh.HandleError(err, "get average shade")
	}

	allowedSet := make(map[int32]struct{})
	for _, level := range []int32{int32(enum.ShadeLevelFull), int32(enum.ShadeLevelPartial), int32(enum.ShadeLevelShade)} {
		if math.Abs(avg-float64(level)) <= ShadeTolerance {
			allowedSet[level] = struct{}{}
		}
	}

	if len(allowedSet) == 0 {
		allowedSet[int32(enum.ShadeLevelFull)] = struct{}{}
	}

	result := make([]int32, 0, len(allowedSet))
	for level := range allowedSet {
		result = append(result, level)
	}
	return result, nil
}
