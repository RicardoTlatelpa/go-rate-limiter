package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)


var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: os.Getenv("REDIS_ADDR"),
	Password: "",
	DB: 0,
})

func RedisStatusHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = "global"			
		}

		keys := []string{
			"rate:" + ip + ":requests",
			"rate:" + ip + ":allowed",
			"rate:" + ip + ":blocked",
			"rate:" + ip + ":first_seen",
			"rate:" + ip + ":last_seen",
		}

		vals, err := rdb.MGet(ctx, keys...).Result()

		if err != nil {
			http.Error(w, "Redis error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, "Stats for IP", ip)
		for i, key := range keys {
			value := "0"
			if vals[i] != nil {
				value = fmt.Sprintf("%v", vals[i])				
			}

			if keyHasSuffix(key, "first_seen") || keyHasSuffix(key, "last_seen") {
				if ts, err := parseUnix(value); err == nil {
					value = ts.Format(time.RFC1123)
				}
			}
			fmt.Fprintf(w, "%-25s: %s\n",key,value)
		}
	}) 
}


func parseUnix(raw string)(time.Time, error) {
	sec, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(sec, 0), nil
}

func keyHasSuffix(key, suffix string) bool {
	n := len(suffix)
	return len(key) >= n && key[len(key)-n:] == suffix
}