package notification

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, input CreateNotificationInput) (Notification, error)
	List(ctx context.Context) ([]Notification, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, input CreateNotificationInput) (Notification, error) {
	const q = `
		INSERT INTO notifications (title, body)
		VALUES ($1, $2)
		RETURNING id, title, body, created_at
	`

	var n Notification
	err := r.db.QueryRow(ctx, q, input.Title, input.Body).Scan(
		&n.ID,
		&n.Title,
		&n.Body,
		&n.CreatedAt,
	)
	return n, err
}

func (r *repository) List(ctx context.Context) ([]Notification, error) {
	const q = `
		SELECT id, title, body, created_at
		FROM notifications
		ORDER BY id DESC
		LIMIT 100
	`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := make([]Notification, 0)
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.Title, &n.Body, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}
