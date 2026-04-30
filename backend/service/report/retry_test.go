package report

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryScheduler_AllSucceedFirstTry(t *testing.T) {
	sched := newSchedulerWithFakeSleep([]string{"a", "b", "c"})
	var got []string
	err := sched.Run(context.Background(), func(_ context.Context, id string) error {
		got = append(got, id)
		return nil
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(got) != 3 {
		t.Errorf("calls = %d, want 3", len(got))
	}
	if len(sched.failed) != 0 {
		t.Errorf("failed = %v, want []", sched.failed)
	}
}

func TestRetryScheduler_RetriesUntilSuccess(t *testing.T) {
	sched := newSchedulerWithFakeSleep([]string{"x"})
	attempts := 0
	err := sched.Run(context.Background(), func(_ context.Context, id string) error {
		attempts++
		if attempts < 3 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if attempts != 3 {
		t.Errorf("attempts = %d, want 3", attempts)
	}
	if len(sched.failed) != 0 {
		t.Errorf("failed = %v", sched.failed)
	}
}

func TestRetryScheduler_FailsAfterAllQueues(t *testing.T) {
	sched := newSchedulerWithFakeSleep([]string{"y"})
	var attempts int
	err := sched.Run(context.Background(), func(_ context.Context, id string) error {
		attempts++
		return errors.New("nope")
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	// 1 main loop + 10 retry queues = 11 attempts
	if attempts != 11 {
		t.Errorf("attempts = %d, want 11", attempts)
	}
	if len(sched.failed) != 1 || sched.failed[0] != "y" {
		t.Errorf("failed = %v, want [y]", sched.failed)
	}
}

func TestRetryScheduler_EmptyQueueBreaksEarly(t *testing.T) {
	sched := newSchedulerWithFakeSleep([]string{"a", "b"})
	attempts := 0
	err := sched.Run(context.Background(), func(_ context.Context, id string) error {
		attempts++
		if id == "a" && attempts == 1 {
			return errors.New("once")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	// a fails main loop (1) + a succeeds Q1 (1) + b succeeds main loop (1) = 3
	if attempts != 3 {
		t.Errorf("attempts = %d, want 3", attempts)
	}
	// After Q1 the queue is empty → no Q2 sleep should happen.
	if len(sched.sleeps) != 1 {
		t.Errorf("sleeps = %v, want 1 (Q1 only)", sched.sleeps)
	}
	if sched.sleeps[0] != 3*time.Second {
		t.Errorf("sleep[0] = %v, want 3s", sched.sleeps[0])
	}
}

func TestRetryScheduler_WaitSequence(t *testing.T) {
	want := []time.Duration{
		3 * time.Second, 9 * time.Second, 15 * time.Second, 21 * time.Second, 27 * time.Second,
		33 * time.Second, 39 * time.Second, 45 * time.Second, 51 * time.Second, 57 * time.Second,
	}
	for i := 0; i < 10; i++ {
		if got := QueueWait(i + 1); got != want[i] {
			t.Errorf("Q%d = %v, want %v", i+1, got, want[i])
		}
	}
}

// helpers

func newSchedulerWithFakeSleep(initial []string) *fakeScheduler {
	fs := &fakeScheduler{}
	s := NewRetryScheduler(initial)
	s.sleep = fs.fakeSleep
	s.onFinalFail = fs.markFailed
	fs.s = s
	return fs
}

type fakeScheduler struct {
	s      *RetryScheduler
	sleeps []time.Duration
	failed []string
}

func (f *fakeScheduler) fakeSleep(_ context.Context, d time.Duration) {
	f.sleeps = append(f.sleeps, d)
}

func (f *fakeScheduler) markFailed(_ context.Context, id string) {
	f.failed = append(f.failed, id)
}

func (f *fakeScheduler) Run(ctx context.Context, attempt func(ctx context.Context, id string) error) error {
	return f.s.Run(ctx, attempt)
}
