package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RicardoTlatelpa/go-rate-limiter/middleware"
)

func main() {
	// Custom router for precise control over routes
	mux := http.NewServeMux()

	// Route: /status
	mux.Handle("/status", middleware.RedisStatusHandler())

	// Route: /
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Request allowed!")
	})
	mux.Handle("/", middleware.RedisRateLimitMiddleware(5, 1.0, testHandler))

	// Start the server
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
