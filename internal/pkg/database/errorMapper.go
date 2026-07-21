package database

import (
	"errors"
	"fmt"
	"garden-nook/internal/pkg/apperrors"
	"log/slog"
	"runtime"

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
	m.log.Error("DB error handled", "err", err, "source", formatSource())
	return err
}

func formatSource() string {
	pc, file, line, ok := runtime.Caller(2)
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		return fmt.Sprintf("%s (%s:%d)", funcName, file, line)
	}
	return "undefined"
}
