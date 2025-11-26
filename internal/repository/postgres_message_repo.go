package repository

import (
	"context"
	"time"

	"github.com/LevanPro/insider/internal/domain"
	"github.com/jmoiron/sqlx"
)

type PostgresMessageRepository struct {
	db *sqlx.DB
}

func NewPostgresMessageRepository(db *sqlx.DB) *PostgresMessageRepository {
	return &PostgresMessageRepository{db: db}
}

func (r *PostgresMessageRepository) GetNextUnsent(ctx context.Context, limit int) ([]domain.Message, error) {
	var msgs []domain.Message
	err := r.db.SelectContext(ctx, &msgs, `
      SELECT id, "to", content, status, sent_at, external_id, created_at, updated_at
      FROM messages
      WHERE status = 'pending'
      ORDER BY id
      LIMIT $1
    `, limit)
	return msgs, err
}

func (r *PostgresMessageRepository) MarkAsSent(ctx context.Context, id int64, sentAt time.Time, externalID *string) error {
	_, err := r.db.ExecContext(ctx, `
      UPDATE messages
      SET status = 'sent',
          sent_at = $2,
          external_id = $3,
          updated_at = NOW()
      WHERE id = $1
    `, id, sentAt, externalID)
	return err
}

func (r *PostgresMessageRepository) MarkAsFailed(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `
      UPDATE messages
      SET status = 'failed',
          updated_at = NOW()
      WHERE id = $1
    `, id)
	return err
}

func (r *PostgresMessageRepository) ListSent(
	ctx context.Context,
	limit, offset int,
) ([]domain.Message, error) {

	var msgs []domain.Message

	query := `
        SELECT id, "to", content, status, sent_at, external_id, created_at, updated_at
        FROM messages
        WHERE status = 'sent'
        ORDER BY sent_at DESC
        LIMIT $1 OFFSET $2
    `

	err := r.db.SelectContext(ctx, &msgs, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
