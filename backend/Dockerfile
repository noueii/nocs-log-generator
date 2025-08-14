# Multi-stage Dockerfile for Go backend
# Stage 1: Development stage with hot reload
FROM golang:1.24-alpine AS dev
LABEL stage=development

# Install necessary tools
RUN apk add --no-cache git make curl

# Install air for hot reload in development
RUN go install github.com/air-verse/air@latest

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Default command for development (will be overridden by docker-compose)
CMD ["air", "-c", ".air.toml"]

# Stage 2: Builder stage for production build
FROM golang:1.24-alpine AS builder
LABEL stage=build

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary, GOOS=linux for Linux target
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Stage 3: Production runtime stage
FROM alpine:3.19 AS production
LABEL stage=production

# Install ca-certificates for HTTPS requests and curl for health checks
RUN apk --no-cache add ca-certificates curl

# Create non-root user for security
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/server .

# Change ownership of the binary
RUN chown appuser:appuser server

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the binary
CMD ["./server"]