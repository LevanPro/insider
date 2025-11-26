package service

import "context"

type SendResponse struct {
	MessageID string
}

type Sender interface {
	Send(ctx context.Context, to, content string) (*SendResponse, error)
}
