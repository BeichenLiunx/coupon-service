package ratelimiter

import (
	"context"
	"sync"
	"time"
)

type Limiter struct {
	Attempts   int
	WindowEnds time.Time
}

type RateLimiter struct {
	limiters          map[string]*Limiter
	mu                sync.Mutex
	requestsPerSecond int
}

const (
	RATE_LIMIT_WINDOW = time.Minute
)

func NewRateLimiter(requestsPerSecond int, ctx context.Context) *RateLimiter {
	r := &RateLimiter{
		limiters:          make(map[string]*Limiter),
		requestsPerSecond: requestsPerSecond,
	}

	go r.startRateLimitCleanup(ctx)

	return r
}

func (r *RateLimiter) IsRateLimited(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	info, exists := r.limiters[ip]
	now := time.Now()

	if !exists || now.After(info.WindowEnds) {
		// reset counter
		r.limiters[ip] = &Limiter{
			Attempts:   1,
			WindowEnds: now.Add(RATE_LIMIT_WINDOW),
		}
		return false
	}

	if info.Attempts >= r.requestsPerSecond {
		return true
	}

	info.Attempts++
	return false
}

func (r *RateLimiter) startRateLimitCleanup(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			r.mu.Lock()
			now := time.Now()
			for ip, info := range r.limiters {
				if now.After(info.WindowEnds) {
					delete(r.limiters, ip)
				}
			}
			r.mu.Unlock()
		}
	}
}
