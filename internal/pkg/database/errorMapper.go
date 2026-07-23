package database

import (
	"errors"
	"garden-nook/internal/pkg/apperrors"
	"garden-nook/internal/pkg/helpers"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ErrorMapper struct {
	pgCodeToErr map[string]error
	log         *slog.Logger
}

func NewErrorMapper(custom map[string]error, log *slog.Logger) *ErrorMapper {
	m := map[string]error{
		"23505": apperrors.ErrConflict,   // unique_violation
		"23503": apperrors.ErrBadRequest, // foreign_key_violation
	}
	for k, v := range custom {
		m[k] = v
	}
	return &ErrorMapper{pgCodeToErr: m, log: log}
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
	m.log.Error("DB error", slog.Any("err", err), slog.String("source", helpers.GetLogSource()))
	return err
}
