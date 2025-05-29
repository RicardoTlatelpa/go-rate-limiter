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
	lastRefillRate time.Time
	mu sync.Mutex // thread safety
}