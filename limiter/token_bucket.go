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

/* 
	Why use float for tokens
		smooths out burstiness
		prevents under-utilization
		useufll when refill rate is less than 1/sec
*/

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
	/* 
		Locks the mutex to make this section thread safe
		concurrent requests don't corrupt shared state (tokens, lastRefillTime)		
	*/
	defer tb.mu.Unlock()
	/*
		mutex is automatically unlocked when the function exits, even if
		early return occurs
	*/
	now := time.Now()
	elapsed := now.Sub(tb.lastRefillTime).Seconds()
	refilledTokens := elapsed * tb.refillRate
	/*
		Refilling tokens every second
	*/
	tb.tokens = min(float64(tb.capacity), tb.tokens+refilledTokens)
	tb.lastRefillTime = now

	/*

	*/
	if tb.tokens >= 1 {
		tb.tokens -= 1
		return true
	}
	/*
		Consume the token if there's at least 1 token in the TokenBucket
	*/
	return false
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}