package scheduler

import (
	"context"
	"errors"
	"sync"
	"time"
)

type CallbackFn func(context.Context) error

var (
	ErrAlreadyRunning = errors.New("scheduler is already running")
	ErrAlreadyStopped = errors.New("scheduler is already stopped")
)

type Scheduler struct {
	callBackFn       CallbackFn
	interval         time.Duration
	startImmediately bool

	mu      sync.Mutex
	ticker  *time.Ticker
	quit    chan struct{}
	running bool
}

func NewScheduler(callBackFn CallbackFn, interval time.Duration, startImmediately bool) *Scheduler {
	return &Scheduler{
		callBackFn:       callBackFn,
		interval:         interval,
		startImmediately: startImmediately,
	}
}

func (s *Scheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return ErrAlreadyRunning
	}

	s.ticker = time.NewTicker(s.interval)
	s.quit = make(chan struct{})
	s.running = true

	go s.run()

	return nil
}

func (s *Scheduler) run() {

	if s.startImmediately {
		_ = s.callBackFn(context.Background())
	}

	for {
		select {
		case <-s.ticker.C:
			_ = s.callBackFn(context.Background())
		case <-s.quit:
			s.ticker.Stop()
			return
		}
	}
}

func (s *Scheduler) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return ErrAlreadyStopped
	}

	close(s.quit)
	s.running = false

	return nil
}

func (s *Scheduler) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}
