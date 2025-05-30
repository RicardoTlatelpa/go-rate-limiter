package limiter

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func AllowRequestRedis(rdb *redis.Client, userKey string, capacity int, refillRate float64)(bool,error) {
	keyTokens := "rate" + userKey + ":tokens"
	keyLast := "rate:" + userKey  + ":last"
	result, err := rdb.MGet(ctx, keyTokens, keyLast).Result()
	if err != nil {
		return false, err
	}

	var tokens float64 = float64(capacity)
	var lastRefillTime float64 = float64(time.Now().Unix())

	if result[0] != nil {
		tokens, _ = parseFloat(result[0])
	}
	if result[1] != nil {
		lastRefillTime, _ = parseFloat(result[1])
	}

	now := float64(time.Now().Unix())
	elapsed := now - lastRefillTime
	refilled := elapsed * refillRate
	tokens = minFloat(float64(capacity), tokens+refilled)

	allowed := false
	if tokens >= 1.0 {
		tokens -= 1.0
		allowed = true
	}
	pipe := rdb.TxPipeline()
	pipe.Set(ctx, keyTokens, tokens, 0)
	pipe.Set(ctx, keyLast, now, 0)
	_, err = pipe.Exec(ctx)

	return allowed, err
}


func parseFloat(v interface{}) (float64, error) {
	switch val := v.(type) {
	case string:
		return strconv.ParseFloat(val,64)
	case []byte:
		return strconv.ParseFloat(string(val),64)
	case float64:
		return val, nil
	default:
		return 0, nil
	}

}

func minFloat(a,b float64) float64 {
	if a < b {
		return a
	}
	return b
}