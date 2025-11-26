package domain

import "time"

type MessageStatus string

const (
	StatusPending MessageStatus = "pending"
	StatusSent    MessageStatus = "sent"
	StatusFailed  MessageStatus = "failed"
)

type Message struct {
	ID         int64         `db:"id"`
	To         string        `db:"to"`
	Content    string        `db:"content"`
	Status     MessageStatus `db:"status"`
	SentAt     *time.Time    `db:"sent_at"`
	ExternalID *string       `db:"external_id"`
	CreatedAt  time.Time     `db:"created_at"`
	UpdatedAt  time.Time     `db:"updated_at"`
}
