package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/LevanPro/insider/internal/domain"
	"github.com/LevanPro/insider/internal/repository"
	"go.uber.org/zap"
)

var (
	ErrGetMessageFail = errors.New("failed to get unsent messages")
)

type MessageService struct {
	repo       repository.MessageRepository
	sender     Sender
	batchSize  int
	numWorkers int
	log        *zap.SugaredLogger
}

func NewMessageService(repo repository.MessageRepository, sender Sender, batchSize int, numWorkers int, log *zap.SugaredLogger) *MessageService {
	if numWorkers <= 0 {
		numWorkers = 1
	}
	return &MessageService{
		repo:       repo,
		sender:     sender,
		batchSize:  batchSize,
		numWorkers: numWorkers,
		log:        log,
	}
}

func (s *MessageService) ProcessNextUnsent(ctx context.Context) error {
	s.log.Infow("Starting processing")
	defer s.log.Infow("End processing")

	msgs, err := s.repo.GetNextUnsent(ctx, s.batchSize)
	if err != nil {
		return ErrGetMessageFail
	}

	if len(msgs) == 0 {
		s.log.Infow("No messages to process")
		return nil
	}

	s.log.Infow("Processing messages", "count", len(msgs))

	msgChan := make(chan domain.Message, len(msgs))

	var wg sync.WaitGroup

	workerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := 0; i < s.numWorkers; i++ {
		wg.Add(1)
		go s.worker(workerCtx, i, msgChan, &wg)
	}

	for _, msg := range msgs {
		select {
		case msgChan <- msg:
		case <-ctx.Done():
			s.log.Warnw("Ctx cancelled while sending data to workers")
			close(msgChan)
			cancel()
			wg.Wait()
			return ctx.Err()
		}
	}

	close(msgChan)

	wg.Wait()

	return nil
}

func (s *MessageService) worker(ctx context.Context, workerID int, msgChan <-chan domain.Message, wg *sync.WaitGroup) {
	defer wg.Done()

	s.log.Infow("Worker started", "workerID", workerID)
	defer s.log.Infow("Worker stopped", "workerID", workerID)

	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return
			}

			s.processMessage(ctx, workerID, msg)

		case <-ctx.Done():
			s.log.Warnw("Worker ctx cancelled", "workerID", workerID)
			return
		}
	}
}

func (s *MessageService) processMessage(ctx context.Context, workerID int, msg domain.Message) {
	s.log.Infow("Processing message",
		"workerID", workerID,
		"messageID", msg.ID,
		"to", msg.To)

	// TODO:: need to think of what do to with such messages
	if len(msg.Content) > 160 {
		s.log.Warnw("message content exceeds 160 characters", "workerID", workerID, "messageID", msg.ID, "length", len(msg.Content))
		return
	}

	// TODO:: implement retry logic
	resp, err := s.sender.Send(ctx, msg.To, msg.Content)
	if err != nil {
		s.log.Errorw("Failed to send message", "workerID", workerID, "messageID", msg.ID, "error", err)
		return
	}

	now := time.Now().UTC()
	extID := resp.MessageID
	if err := s.repo.MarkAsSent(ctx, msg.ID, now, &extID); err != nil {
		s.log.Errorw("Unable to mark message as sent", "workerID", workerID, "messageID", msg.ID, "externalID", extID, "error", err)
		return
	}

	s.log.Infow("Message has been sent successfully", "workerID", workerID, "messageID", msg.ID, "externalID", extID)
}

func (s *MessageService) ListSent(ctx context.Context, limit, offset int) ([]domain.Message, error) {
	return s.repo.ListSent(ctx, limit, offset)
}
