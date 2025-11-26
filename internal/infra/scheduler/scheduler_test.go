package scheduler_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/LevanPro/insider/internal/infra/scheduler"
)

const (
	testInterval        = 10 * time.Millisecond
	waitDuration        = 3 * testInterval
	initialWaitDuration = 5 * testInterval
)

type mockCallback struct {
	count atomic.Int32
}

func (m *mockCallback) Fn(ctx context.Context) error {
	m.count.Add(1)
	return nil
}

func (m *mockCallback) GetCount() int32 {
	return m.count.Load()
}

func TestNewScheduler(t *testing.T) {
	mock := &mockCallback{}
	s := scheduler.NewScheduler(mock.Fn, testInterval, true)

	if s == nil {
		t.Fatal("NewScheduler returned nil")
	}
	if s.IsRunning() {
		t.Error("Scheduler should not be running after creation")
	}
}

func TestStartImmediately(t *testing.T) {
	mock := &mockCallback{}
	s := scheduler.NewScheduler(mock.Fn, testInterval, true)

	if err := s.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer s.Stop()

	time.Sleep(1 * time.Millisecond)

	if mock.GetCount() != 1 {
		t.Errorf("Expected count 1 immediately after start, got %d", mock.GetCount())
	}
}

func TestNoStartImmediately(t *testing.T) {
	mock := &mockCallback{}
	s := scheduler.NewScheduler(mock.Fn, testInterval, false)

	if err := s.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer s.Stop()

	time.Sleep(1 * time.Millisecond)

	if mock.GetCount() != 0 {
		t.Errorf("Expected count 0 immediately after start, got %d", mock.GetCount())
	}
}

func TestSchedulerExecution(t *testing.T) {
	mock := &mockCallback{}
	s := scheduler.NewScheduler(mock.Fn, testInterval, false)

	if err := s.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer s.Stop()

	time.Sleep(waitDuration)

	initialCount := mock.GetCount()
	if initialCount < 2 {
		t.Errorf("Expected at least 2 executions after %v, got %d", waitDuration, initialCount)
	}

	time.Sleep(testInterval)
	finalCount := mock.GetCount()

	if finalCount <= initialCount {
		t.Errorf("Expected count to increase after another interval, initial: %d, final: %d", initialCount, finalCount)
	}
}

func TestStop(t *testing.T) {
	mock := &mockCallback{}
	s := scheduler.NewScheduler(mock.Fn, testInterval, false)

	if err := s.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// 1. Wait long enough for ticks to fire, plus a small buffer, to ensure a stable count.
	time.Sleep(initialWaitDuration)
	time.Sleep(1 * time.Millisecond) // Allow any concurrent tick to complete its callback execution

	countBeforeStop := mock.GetCount()

	// 2. Stop the scheduler
	if err := s.Stop(); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	if s.IsRunning() {
		t.Error("IsRunning should be false after Stop")
	}

	// 3. Wait for two more intervals. If the scheduler is correctly stopped, the count should not change.
	time.Sleep(2 * testInterval)
	countAfterStop := mock.GetCount()

	if countAfterStop != countBeforeStop {
		t.Errorf("Expected count to be %d after stop, got %d, meaning a tick fired after Stop()", countBeforeStop, countAfterStop)
	}
}

func TestStartAlreadyRunning(t *testing.T) {
	mock := &mockCallback{}
	s := scheduler.NewScheduler(mock.Fn, testInterval, false)

	if err := s.Start(); err != nil {
		t.Fatalf("First Start failed: %v", err)
	}
	defer s.Stop()

	err := s.Start()
	if !s.IsRunning() {
		t.Fatal("Scheduler should still be running")
	}
	if err == nil || err.Error() != scheduler.ErrAlreadyRunning.Error() {
		t.Errorf("Expected ErrAlreadyRunning, got: %v", err)
	}
}

func TestStopAlreadyStopped(t *testing.T) {
	mock := &mockCallback{}
	s := scheduler.NewScheduler(mock.Fn, testInterval, false)

	err := s.Stop()
	if err == nil || err.Error() != scheduler.ErrAlreadyStopped.Error() {
		t.Errorf("Expected ErrAlreadyStopped when stopping a non-running scheduler, got: %v", err)
	}

	if err := s.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	if err := s.Stop(); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	err = s.Stop()
	if err == nil || err.Error() != scheduler.ErrAlreadyStopped.Error() {
		t.Errorf("Expected ErrAlreadyStopped after successful stop, got: %v", err)
	}
}

func TestIsRunning(t *testing.T) {
	mock := &mockCallback{}
	s := scheduler.NewScheduler(mock.Fn, testInterval, false)

	if s.IsRunning() {
		t.Error("Should not be running initially")
	}

	if err := s.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	if !s.IsRunning() {
		t.Error("Should be running after Start")
	}

	if err := s.Stop(); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
	if s.IsRunning() {
		t.Error("Should not be running after Stop")
	}
}
