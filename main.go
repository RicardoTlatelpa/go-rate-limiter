package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RicardoTlatelpa/go-rate-limiter/middleware"
)


func main() {
	// 5 max tokens per user, 1 token added per second
	rateLimiter := middleware.NewRateLimiterMiddleware(5,1.0)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintln(w, "Request allowed!")
	})

	http.Handle("/", rateLimiter.MiddlewareFunc(testHandler))
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request){
		rl := rateLimiter
		rl.Mu.Lock()
		defer rl.Mu.Unlock()
		fmt.Fprint(w, "Rate limiter Report")
		for ip, stats := range rl.Stats {
			fmt.Fprintf(w, "IP: %s\n", ip)
			fmt.Fprintf(w, "  Total Requests:   %d\n", stats.Requests)
			fmt.Fprintf(w, "  Allowed:          %d\n", stats.Allowed)
			fmt.Fprintf(w, "  Blocked:          %d\n", stats.Blocked)
			fmt.Fprintf(w, "  First Seen:       %s\n", stats.FirstSeen.Format(time.RFC3339))
			fmt.Fprintf(w, "  Last Seen:        %s\n", stats.LastSeen.Format(time.RFC3339))
			fmt.Fprintf(w, "\n")
		}
	})

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}