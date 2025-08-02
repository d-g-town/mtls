# Multi-stage build for Service2 (PaaS-ready)
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service2 .

# Final stage
FROM alpine:latest

# Install ca-certificates for SSL/TLS
RUN apk --no-cache add ca-certificates

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/service2 .

# Create directory for certificates (will be mounted as secrets)
RUN mkdir -p /certs

# Expose port (configurable via PORT env var)
EXPOSE 8444

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-check-certificate --spider https://localhost:${PORT:-8444}/health || exit 1

# Run the application
CMD ["./service2"]