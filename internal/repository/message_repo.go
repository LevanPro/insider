package repository

import (
	"context"
	"time"

	"github.com/LevanPro/insider/internal/domain"
)

type MessageRepository interface {
	GetNextUnsent(ctx context.Context, limit int) ([]domain.Message, error)
	MarkAsSent(ctx context.Context, id int64, sentAt time.Time, externalID *string) error
	MarkAsFailed(ctx context.Context, id int64) error
	ListSent(ctx context.Context, limit, offset int) ([]domain.Message, error)
}
