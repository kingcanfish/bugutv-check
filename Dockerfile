# Build stage
FROM golang:1.24-alpine  AS builder

WORKDIR /app
COPY . .
# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o bugutv-signin main.go

# Runtime stage
FROM alpine:latest
WORKDIR /app
# Copy the compiled binary from builder
COPY --from=builder /app/bugutv-signin .

# Environment variables should be passed at runtime
# e.g., docker run -e BUGUTV_USERNAME=xxx -e BUGUTV_PASSWORD=yyy bugutv-signin

ENTRYPOINT ["./bugutv-signin"]