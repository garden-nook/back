package crops

import (
	"context"
	"errors"
	"garden-nook/internal/pkg/response"
	"log/slog"

	"garden-nook/internal/pkg/apperrors"

	"golang.org/x/sync/errgroup"
)

type Service struct {
	repo *Repository
	log  *slog.Logger
}

func NewService(repo *Repository, log *slog.Logger) *Service {
	return &Service{repo: repo, log: log}
}

// ---------- SOIL TYPES ----------

func (s *Service) ListSoilTypes(ctx context.Context) ([]SoilType, int, error) {
	return s.repo.ListSoilTypes(ctx, nil)
}

func (s *Service) GetSoilType(ctx context.Context, id int32) (*SoilType, error) {
	return s.repo.GetSoilTypeByID(ctx, id)
}

func (s *Service) CreateSoilType(ctx context.Context, req CreateSoilTypeRequest) (*response.CreateUpdateIntId, error) {
	f, err := s.repo.CreateSoilType(ctx, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrConflict) {
			return nil, err
		}
		s.log.Error("create soil type failed", "err", err)
		return nil, apperrors.ErrInternal
	}
	return &response.CreateUpdateIntId{Id: f}, nil
}

func (s *Service) UpdateSoilType(ctx context.Context, id int32, req UpdateSoilTypeRequest) (*response.CreateUpdateIntId, error) {
	f, err := s.repo.UpdateSoilType(ctx, id, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) || errors.Is(err, apperrors.ErrConflict) {
			return nil, err
		}
		s.log.Error("update soil type failed", "err", err, "id", id)
		return nil, apperrors.ErrInternal
	}
	return &response.CreateUpdateIntId{Id: f}, nil
}

func (s *Service) DeleteSoilType(ctx context.Context, id int32) error {
	err := s.repo.DeleteSoilType(ctx, id)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) && !errors.Is(err, apperrors.ErrConflict) {
		s.log.Error("delete soil type failed", "err", err, "id", id)
		return apperrors.ErrInternal
	}
	return err
}

// ---------- FAMILIES ----------

func (s *Service) ListFamilies(ctx context.Context) ([]CropFamily, int, error) {
	return s.repo.ListFamilies(ctx, nil)
}

func (s *Service) GetFamily(ctx context.Context, id int32) (*CropFamily, error) {
	return s.repo.GetFamilyByID(ctx, id)
}

func (s *Service) CreateFamily(ctx context.Context, req CreateFamilyRequest) (*response.CreateUpdateIntId, error) {
	f, err := s.repo.CreateFamily(ctx, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrConflict) {
			return nil, err
		}
		s.log.Error("create family failed", "err", err)
		return nil, apperrors.ErrInternal
	}
	return &response.CreateUpdateIntId{Id: f}, nil
}

func (s *Service) UpdateFamily(ctx context.Context, id int32, req UpdateFamilyRequest) (*response.CreateUpdateIntId, error) {
	f, err := s.repo.UpdateFamily(ctx, id, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) || errors.Is(err, apperrors.ErrConflict) {
			return nil, err
		}
		s.log.Error("update family failed", "err", err, "id", id)
		return nil, apperrors.ErrInternal
	}
	return &response.CreateUpdateIntId{Id: f}, nil
}

func (s *Service) DeleteFamily(ctx context.Context, id int32) error {
	err := s.repo.DeleteFamily(ctx, id)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) && !errors.Is(err, apperrors.ErrConflict) {
		s.log.Error("delete family failed", "err", err, "id", id)
		return apperrors.ErrInternal
	}
	return err
}

// ---------- CROPS ----------

func (s *Service) ListCrops(ctx context.Context, f ListCropsFilter) ([]Crop, int, error) {
	return s.repo.ListCrops(ctx, f, nil)
}

func (s *Service) GetCrop(ctx context.Context, id int32) (*CropExtended, error) {
	var (
		c  *Crop
		cr *CropRelations
	)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		c, err = s.repo.GetCropByID(ctx, id)
		return err
	})
	g.Go(func() error {
		var err error
		cr, err = s.repo.GetCropRelations(ctx, id)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return &CropExtended{c, cr}, nil
}

func (s *Service) CreateCrop(ctx context.Context, req CreateCropRequest) (*response.CreateUpdateIntId, error) {
	if _, err := s.repo.GetFamilyByID(ctx, req.FamilyID); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrBadRequest
		}
		return nil, err
	}
	c, err := s.repo.CreateCrop(ctx, req)
	if err != nil {
		s.log.Error("create crop failed", "err", err)
		return nil, apperrors.ErrInternal
	}
	return &response.CreateUpdateIntId{Id: c}, nil
}

func (s *Service) UpdateCrop(ctx context.Context, id int32, req UpdateCropRequest) (*response.CreateUpdateIntId, error) {
	if req.FamilyID != nil {
		if _, err := s.repo.GetFamilyByID(ctx, *req.FamilyID); err != nil {
			if errors.Is(err, apperrors.ErrNotFound) {
				return nil, apperrors.ErrBadRequest
			}
			return nil, err
		}
	}
	c, err := s.repo.UpdateCrop(ctx, id, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, err
		}
		s.log.Error("update crop failed", "err", err, "id", id)
		return nil, apperrors.ErrInternal
	}
	return &response.CreateUpdateIntId{Id: c}, nil
}

func (s *Service) DeleteCrop(ctx context.Context, id int32) error {
	err := s.repo.SoftDeleteCrop(ctx, id)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		s.log.Error("delete crop failed", "err", err, "id", id)
		return apperrors.ErrInternal
	}
	return err
}

// ---------- RULES ----------

func (s *Service) ListRules(ctx context.Context) ([]CropRule, int, error) {
	return s.repo.ListRules(ctx, nil)
}

func (s *Service) CreateRule(ctx context.Context, req CreateRuleRequest) (*response.CreateUpdateIntId, error) {
	// Бизнес-правило: должен быть указан либо crop, либо family (но не оба одновременно на одной стороне)
	if req.SubjectCropID != nil && req.SubjectFamilyID != nil {
		return nil, apperrors.ErrBadRequest
	}
	if req.ContextCropID != nil && req.ContextFamilyID != nil {
		return nil, apperrors.ErrBadRequest
	}
	if req.SubjectCropID == nil && req.SubjectFamilyID == nil {
		return nil, apperrors.ErrBadRequest
	}
	if req.ContextCropID == nil && req.ContextFamilyID == nil {
		return nil, apperrors.ErrBadRequest
	}

	r, err := s.repo.CreateRule(ctx, req)
	if err != nil {
		s.log.Error("create rule failed", "err", err)
		return nil, apperrors.ErrInternal
	}
	return &response.CreateUpdateIntId{Id: r}, nil
}

func (s *Service) DeleteRule(ctx context.Context, id int32) error {
	err := s.repo.DeleteRule(ctx, id)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		s.log.Error("delete rule failed", "err", err, "id", id)
		return apperrors.ErrInternal
	}
	return err
}
