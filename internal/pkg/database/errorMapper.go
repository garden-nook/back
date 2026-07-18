package database

import (
	"errors"
	"garden-nook/internal/pkg/apperrors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ErrorMapper struct {
	pgCodeToErr map[string]error
}

func NewErrorMapper(custom map[string]error) *ErrorMapper {
	m := map[string]error{
		"23505": apperrors.ErrConflict,   // unique_violation
		"23503": apperrors.ErrBadRequest, // foreign_key_violation
	}
	for k, v := range custom {
		m[k] = v
	}
	return &ErrorMapper{pgCodeToErr: m}
}

func (m *ErrorMapper) Map(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return apperrors.ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if mapped, ok := m.pgCodeToErr[pgErr.Code]; ok {
			return mapped
		}
	}
	return err
}
