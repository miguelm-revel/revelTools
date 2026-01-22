package syncx

import (
	"context"
	"time"
)

type RateLimiter struct {
	tokens chan struct{}
	stop   chan struct{}
}

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

func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-rl.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (rl *RateLimiter) Stop() { close(rl.stop) }
