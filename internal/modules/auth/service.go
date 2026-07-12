package auth

import (
	"context"
	"errors"
	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/jwt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo   *Repository
	jwtMgr *jwt.Manager
	log    *slog.Logger
}

func NewService(repo *Repository, jwtMgr *jwt.Manager, log *slog.Logger) *Service {
	return &Service{repo: repo, jwtMgr: jwtMgr, log: log}
}

// ---------- USER AUTH ----------

func (s *Service) RegisterUser(ctx context.Context, req RegisterRequest) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("bcrypt hash failed", "err", err)
		return nil, apperrors.ErrInternal
	}

	u, err := s.repo.CreateUser(ctx, req.Email, string(hash), req.DisplayName)
	if err != nil {
		if errors.Is(err, apperrors.ErrConflict) {
			return nil, err
		}
		s.log.Error("create user failed", "err", err)
		return nil, apperrors.ErrInternal
	}
	return u, nil
}

func (s *Service) LoginUser(ctx context.Context, req LoginRequest) (*TokenResponse, error) {
	user, hash, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			// Не раскрываем, что email не существует
			return nil, apperrors.ErrUnauthorized
		}
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrUnauthorized
	}

	token, err := s.jwtMgr.Generate(user.UserID, jwt.SubjectUser)
	if err != nil {
		s.log.Error("generate user token failed", "err", err)
		return nil, apperrors.ErrInternal
	}

	return &TokenResponse{
		AccessToken: token,
		TokenType:   string(jwt.SubjectUser),
		ExpiresIn:   int(s.jwtMgr.UserTTL().Seconds()),
	}, nil
}

func (s *Service) GetUserProfile(ctx context.Context, userID string) (*MeResponse, error) {
	u, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &MeResponse{
		ID:          u.UserID,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		Type:        string(jwt.SubjectUser),
	}, nil
}

// ---------- ADMIN AUTH ----------

func (s *Service) LoginAdmin(ctx context.Context, req LoginRequest) (*TokenResponse, error) {
	admin, hash, err := s.repo.GetAdminByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrUnauthorized
		}
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrUnauthorized
	}

	token, err := s.jwtMgr.Generate(admin.AdminID, jwt.SubjectAdmin)
	if err != nil {
		s.log.Error("generate admin token failed", "err", err)
		return nil, apperrors.ErrInternal
	}

	return &TokenResponse{
		AccessToken: token,
		TokenType:   string(jwt.SubjectAdmin),
		ExpiresIn:   int(s.jwtMgr.AdminTTL().Seconds()),
	}, nil
}

func (s *Service) GetAdminProfile(ctx context.Context, adminID string) (*MeResponse, error) {
	a, err := s.repo.GetAdminByID(ctx, adminID)
	if err != nil {
		return nil, err
	}
	return &MeResponse{
		ID:          a.AdminID,
		Email:       a.Email,
		DisplayName: a.DisplayName,
		Type:        string(jwt.SubjectAdmin),
	}, nil
}
