package limiter

/*
We will build a struct to store
 Max tokens
 Current tokens
 Refill rate
 Last refill timestamp
 Mutex for concurrency
Methods
 Allow(userID string) bool - the core method
 refill() - update tokens based on elapsed time

*/

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity int
	tokens float64
	refillRate float64
	lastRefillTime time.Time
	mu sync.Mutex // thread safety
}

func NewTokenBucket(capacity int, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity: capacity,
		tokens: float64(capacity),
		refillRate: refillRate,
		lastRefillTime: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	
	now := time.Now()
	elapsed := now.Sub(tb.lastRefillTime).Seconds()
	refilledTokens := elapsed * tb.refillRate

	tb.tokens = min(float64(tb.capacity), tb.tokens+refilledTokens)
	tb.lastRefillTime = now

	if tb.tokens >= 1 {
		tb.tokens -= 1
		return true
	}

	return false
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}