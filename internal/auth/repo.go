package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, username string, passwordHash string) (Account, error)
	GetByUsername(ctx context.Context, username string) (Account, error)
	DeleteByID(ctx context.Context, id int64) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, username string, passwordHash string) (Account, error) {
	const q = `
		INSERT INTO accounts (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, username, password_hash, created_at, updated_at
	`

	var a Account
	err := r.db.QueryRow(ctx, q, username, passwordHash).Scan(
		&a.ID,
		&a.Username,
		&a.PasswordHash,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	return a, err
}

func (r *repository) GetByUsername(ctx context.Context, username string) (Account, error) {
	const q = `
		SELECT id, username, password_hash, created_at, updated_at
		FROM accounts
		WHERE username = $1
	`

	var a Account
	err := r.db.QueryRow(ctx, q, username).Scan(
		&a.ID,
		&a.Username,
		&a.PasswordHash,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	return a, err
}

func (r *repository) DeleteByID(ctx context.Context, id int64) error {
	const q = `DELETE FROM accounts WHERE id = $1`
	_, err := r.db.Exec(ctx, q, id)
	return err
}
