package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RicardoTlatelpa/go-rate-limiter/middleware"
)


func main() {
	// 5 max tokens per user, 1 token added per second
	
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintln(w, "Request allowed!")
	})	

	http.Handle("/", middleware.RedisRateLimitMiddleware(5, 1.0, testHandler))
	http.Handle("/status", middleware.RedisStatusHandler())
	
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}