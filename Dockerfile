# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy web assets
COPY --from=builder /app/web ./web

# Create posts directory
RUN mkdir -p /data/posts

# Copy posts from builder stage
COPY --from=builder /app/posts /data/posts

# Expose port
EXPOSE 3000

# Set environment variables
ENV POSTS_DIR=/data/posts

# Run the application
CMD ["./main"]
