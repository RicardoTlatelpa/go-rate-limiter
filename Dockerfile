# Build stage: using official Go image
FROM golang:1.24.3-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first for better layer caching (missing this in your version)
COPY go.mod go.sum ./

# Download dependencies (only works if go.sum is present)
RUN go mod download

# Now copy the rest of the source code
COPY . .

# Build the application
RUN go build -o rate-limiter main.go

# Final image stage: a lightweight base
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy the built binary from builder stage
COPY --from=builder /app/rate-limiter .

# Set the startup command
ENTRYPOINT ["./rate-limiter"]

# Document the port the app uses
EXPOSE 8080
