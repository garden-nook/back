package auth

import (
	"context"
	"errors"
	"garden-nook/internal/pkg/apperrors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// ---------- USERS ----------

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, string, error) {
	u := &User{}
	var hash string
	err := r.db.QueryRow(ctx,
		`SELECT user_id, email, display_name, password_hash, created_at
		 FROM users WHERE email=$1`, email,
	).Scan(&u.UserID, &u.Email, &u.DisplayName, &hash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, "", apperrors.ErrNotFound
	}
	if err != nil {
		return nil, "", err
	}
	return u, hash, nil
}

func (r *Repository) CreateUser(ctx context.Context, email, hash, displayName string) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx,
		`INSERT INTO users(email, password_hash, display_name)
		 VALUES ($1, $2, $3) RETURNING user_id, email, display_name, created_at`,
		email, hash, displayName,
	).Scan(&u.UserID, &u.Email, &u.DisplayName, &u.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, apperrors.ErrConflict
		}
		return nil, err
	}
	return u, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id string) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx,
		`SELECT user_id, email, display_name, created_at FROM users WHERE user_id=$1`, id,
	).Scan(&u.UserID, &u.Email, &u.DisplayName, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.ErrNotFound
	}
	return u, err
}

// ---------- ADMINS ----------

func (r *Repository) GetAdminByEmail(ctx context.Context, email string) (*Admin, string, error) {
	a := &Admin{}
	var hash string
	err := r.db.QueryRow(ctx,
		`SELECT admin_id, email, display_name, password_hash, created_at
		 FROM admins WHERE email=$1`, email,
	).Scan(&a.AdminID, &a.Email, &a.DisplayName, &hash, &a.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, "", apperrors.ErrNotFound
	}
	if err != nil {
		return nil, "", err
	}
	return a, hash, nil
}

func (r *Repository) GetAdminByID(ctx context.Context, id string) (*Admin, error) {
	a := &Admin{}
	err := r.db.QueryRow(ctx,
		`SELECT admin_id, email, display_name, created_at FROM admins WHERE admin_id=$1`, id,
	).Scan(&a.AdminID, &a.Email, &a.DisplayName, &a.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.ErrNotFound
	}
	return a, err
}
