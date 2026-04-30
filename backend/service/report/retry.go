package report

import (
	"context"
	"time"
)

const maxQueues = 10

// QueueWait returns the backoff before queue n (1-indexed). Linear: 3 + 6*(n-1) seconds.
func QueueWait(n int) time.Duration {
	return time.Duration(3+6*(n-1)) * time.Second
}

// RetryScheduler runs attempt over an initial list, then retries failures
// across up to 10 queues with linear backoff.
type RetryScheduler struct {
	initial     []string
	sleep       func(ctx context.Context, d time.Duration)
	onFinalFail func(ctx context.Context, id string)
}

func NewRetryScheduler(initial []string) *RetryScheduler {
	return &RetryScheduler{
		initial:     initial,
		sleep:       defaultSleep,
		onFinalFail: func(context.Context, string) {},
	}
}

// SetFinalFailureHook sets the callback for cows that exhaust all 10 queues.
func (s *RetryScheduler) SetFinalFailureHook(fn func(ctx context.Context, id string)) {
	s.onFinalFail = fn
}

// Run drains the initial list and retry queues. Returns ctx.Err() on cancel.
func (s *RetryScheduler) Run(ctx context.Context, attempt func(ctx context.Context, id string) error) error {
	queue := drain(ctx, s.initial, attempt)
	for n := 1; n <= maxQueues; n++ {
		if len(queue) == 0 {
			break
		}
		s.sleep(ctx, QueueWait(n))
		if ctx.Err() != nil {
			return ctx.Err()
		}
		queue = drain(ctx, queue, attempt)
	}
	for _, id := range queue {
		s.onFinalFail(ctx, id)
	}
	return nil
}

// drain runs attempt for each id and returns the ones that failed.
func drain(ctx context.Context, ids []string, attempt func(context.Context, string) error) []string {
	var fails []string
	for _, id := range ids {
		if err := attempt(ctx, id); err != nil {
			fails = append(fails, id)
		}
	}
	return fails
}

func defaultSleep(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
	case <-ctx.Done():
	}
}
