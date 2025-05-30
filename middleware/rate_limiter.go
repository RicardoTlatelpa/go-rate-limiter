package middleware // intercepting and processing HTTP requests and responses

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/RicardoTlatelpa/go-rate-limiter/limiter"
)
type ClientStats struct {
	Requests int
	Allowed int
	Blocked int
	FirstSeen time.Time
	LastSeen time.Time

}
type RateLimiterMiddleware struct {
	buckets map[string]*limiter.TokenBucket
	Stats map[string]*ClientStats
	Mu sync.Mutex // prevent inconsistencies when reading/writing data
	cap int
	refill float64
}

func NewRateLimiterMiddleware(capacity int, refillRate float64) *RateLimiterMiddleware{
	return &RateLimiterMiddleware{
		buckets: make(map[string]*limiter.TokenBucket),
		Stats: make(map[string]*ClientStats),
		cap: capacity,
		refill: refillRate,
	}
}

// getBucket fetches or creates a TokenBucket for a given IP address

func (rl *RateLimiterMiddleware) getBucket(ip string) *limiter.TokenBucket {
	rl.Mu.Lock() // lock the function execution
	defer rl.Mu.Unlock() // after function runs => unlock the mutex

	if _, exists := rl.buckets[ip]; !exists {
		rl.buckets[ip] = limiter.NewTokenBucket(rl.cap, rl.refill)
		rl.Stats[ip] = &ClientStats{
			FirstSeen: time.Now(),
		}
	}
	rl.Stats[ip].LastSeen = time.Now()
	return rl.buckets[ip]
}

func (rl *RateLimiterMiddleware) MiddlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = "global"
		}

		bucket := rl.getBucket(ip)
		allowed := bucket.Allow()

		rl.Mu.Lock()
		Stats := rl.Stats[ip]
		Stats.Requests++
		if allowed {
			Stats.Allowed++			
		} else {
			Stats.Blocked++
		}
		rl.Mu.Unlock()
		if allowed {
			next.ServeHTTP(w,r)
		} else {
			http.Error(w, "429 - Rate limit exceeded", http.StatusTooManyRequests)
		}
	}) 
}