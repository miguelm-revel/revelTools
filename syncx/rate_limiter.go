package syncx

import (
	"context"
	"time"
)

// RateLimiter is a token-bucket rate limiter backed by a buffered channel.
// Tokens are replenished at a steady interval derived from the configured rate and duration.
type RateLimiter struct {
	tokens chan struct{}
	stop   chan struct{}
}

// NewRateLimiter creates a RateLimiter that allows rate events per the given duration.
// The bucket is pre-filled to rate tokens and a background goroutine replenishes
// one token per sub-interval until Stop is called.
func NewRateLimiter(rate int, per time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, rate),
		stop:   make(chan struct{}),
	}

	for i := 0; i < rate; i++ {
		rl.tokens <- struct{}{}
	}

	interval := per / time.Duration(rate)
	t := time.NewTicker(interval)

	go func() {
		defer t.Stop()
		for {
			select {
			case <-rl.stop:
				return
			case <-t.C:
				select {
				case rl.tokens <- struct{}{}:
				default:
				}
			}
		}
	}()

	return rl
}

// Allow reports whether a token is immediately available and, if so, consumes it.
func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

// Wait blocks until a token is available or ctx is cancelled.
// Returns ctx.Err() if the context is done before a token becomes available.
func (rl *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-rl.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stop shuts down the background token-replenishment goroutine.
func (rl *RateLimiter) Stop() { close(rl.stop) }
