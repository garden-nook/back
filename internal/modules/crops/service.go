package crops

import (
	"context"
	"errors"
	"log/slog"

	"garden-nook/internal/pkg/apperrors"
)

type Service struct {
	repo *Repository
	log  *slog.Logger
}

func NewService(repo *Repository, log *slog.Logger) *Service {
	return &Service{repo: repo, log: log}
}

// ---------- FAMILIES ----------

func (s *Service) ListFamilies(ctx context.Context) ([]CropFamily, error) {
	return s.repo.ListFamilies(ctx)
}

func (s *Service) GetFamily(ctx context.Context, id int32) (*CropFamily, error) {
	return s.repo.GetFamilyByID(ctx, id)
}

func (s *Service) CreateFamily(ctx context.Context, req CreateFamilyRequest) (*CropFamily, error) {
	f, err := s.repo.CreateFamily(ctx, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrConflict) {
			return nil, err
		}
		s.log.Error("create family failed", "err", err)
		return nil, apperrors.ErrInternal
	}
	return f, nil
}

func (s *Service) UpdateFamily(ctx context.Context, id int32, req UpdateFamilyRequest) (*CropFamily, error) {
	f, err := s.repo.UpdateFamily(ctx, id, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) || errors.Is(err, apperrors.ErrConflict) {
			return nil, err
		}
		s.log.Error("update family failed", "err", err, "id", id)
		return nil, apperrors.ErrInternal
	}
	return f, nil
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
	return s.repo.ListCrops(ctx, f)
}

func (s *Service) GetCrop(ctx context.Context, id int32) (*Crop, error) {
	return s.repo.GetCropByID(ctx, id)
}

func (s *Service) CreateCrop(ctx context.Context, req CreateCropRequest) (*Crop, error) {
	// Проверяем существование семейства
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
	return c, nil
}

func (s *Service) UpdateCrop(ctx context.Context, id int32, req UpdateCropRequest) (*Crop, error) {
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
	return c, nil
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

func (s *Service) ListRules(ctx context.Context) ([]CropRule, error) {
	return s.repo.ListRules(ctx)
}

func (s *Service) CreateRule(ctx context.Context, req CreateRuleRequest) (*CropRule, error) {
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
	return r, nil
}

func (s *Service) DeleteRule(ctx context.Context, id int32) error {
	err := s.repo.DeleteRule(ctx, id)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		s.log.Error("delete rule failed", "err", err, "id", id)
		return apperrors.ErrInternal
	}
	return err
}
