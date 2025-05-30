package middleware

import (
	"net"
	"net/http"
	"os"

	"github.com/RicardoTlatelpa/go-rate-limiter/limiter"
	"github.com/redis/go-redis/v9"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr: os.Getenv("REDIS_ADDR"),
	Password: "",
	DB: 0,
})

func RedisRateLimitMiddleware(capacity int, refillRate float64, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = "global"
		}
		allowed, err := limiter.AllowRequestRedis(redisClient, ip, capacity, refillRate)
		if err != nil {
			http.Error(w, "Redis error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if allowed {
			next.ServeHTTP(w,r)
		} else {
			http.Error(w, "429 - Rate limit exceeded", http.StatusTooManyRequests)
		}
	})
}