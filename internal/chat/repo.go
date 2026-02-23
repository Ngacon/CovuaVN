package chat

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, input CreateMessageInput) (Message, error)
	List(ctx context.Context) ([]Message, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, input CreateMessageInput) (Message, error) {
	const q = `
		INSERT INTO chat_messages (sender_username, message)
		VALUES ($1, $2)
		RETURNING id, sender_username, message, created_at
	`

	var m Message
	err := r.db.QueryRow(ctx, q, input.SenderUsername, input.Message).Scan(
		&m.ID,
		&m.SenderUsername,
		&m.Message,
		&m.CreatedAt,
	)
	return m, err
}

func (r *repository) List(ctx context.Context) ([]Message, error) {
	const q = `
		SELECT id, sender_username, message, created_at
		FROM chat_messages
		ORDER BY id DESC
		LIMIT 100
	`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]Message, 0)
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.SenderUsername, &m.Message, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
