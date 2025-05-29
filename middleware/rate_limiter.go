package middleware // intercepting and processing HTTP requests and responses

import (
	"net"
	"net/http"
	"sync"

	"github.com/RicardoTlatelpa/go-rate-limiter/limiter"
)

type RateLimiterMiddleware struct {
	buckets map[string]*limiter.TokenBucket
	mu sync.Mutex // prevent inconsistencies when reading/writing data
	cap int
	refill float64
}

func NewRateLimiterMiddleware(capacity int, refillRate float64) *RateLimiterMiddleware{
	return &RateLimiterMiddleware{
		buckets: make(map[string]*limiter.TokenBucket),
		cap: capacity,
		refill: refillRate,
	}
}

// getBucket fetches or creates a TokenBucket for a given IP address

func (rl *RateLimiterMiddleware) getBucket(ip string) *limiter.TokenBucket {
	rl.mu.Lock() // lock the function execution
	defer rl.mu.Unlock() // after function runs => unlock the mutex

	if bucket, exists := rl.buckets[ip]; exists {
		return bucket
	}

	bucket := limiter.NewTokenBucket(rl.cap, rl.refill)
	rl.buckets[ip] = bucket
	return bucket
}

func (rl *RateLimiterMiddleware) MiddlewareFunc(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = "global"
		}

		bucket := rl.getBucket(ip)
		if bucket.Allow(){
			next.ServeHTTP(w,r)
		}else {
			http.Error(w, "429 - Rate limit exceeded", http.StatusTooManyRequests)
		}
	}) 
}