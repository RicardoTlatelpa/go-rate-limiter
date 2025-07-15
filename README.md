# Go Rate Limiter middleware

A reusable, Redis-Backed **token bucket rate limiter middleware** for Go HTTP servers.
Provides a simple way to protect your APIs from abuse and control request rates on a per-client basis.

## Features
- Token bucket algorithm
- Redis-based distributed store
- Middleware for 'net/http'
- Support configurable capacity and refill rates
- Includes '/status' endpoint for debugging stats

## Monitoring

Gathering analytical data to measure how effective the rate limiter is.

- The rate limiting algorithm is effective
- The rate limiting rules are effective

## Installation
```bash
go get github.com/RicardoTlatelpa/go-rate-limiter/middleware
```
### import in your go app
```go
import "github.com/RicardoTlatelpa/go-rate-limiter/middleware"
```
## Executing
### locally
```bash
brew install redis
brew services start redis
go run main.go
```
### With kubernetes
```bash
kubectl apply -f k8s/
minikube service rate-limiter-service --url
```

