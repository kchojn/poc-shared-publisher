# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=$(git describe --tags --always --dirty) \
    -X main.BuildTime=$(date -u +%Y-%m-%d_%H:%M:%S) \
    -X main.GitCommit=$(git rev-parse HEAD)" \
    -o poc-shared-publisher \
    cmd/publisher/main.go

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 publisher && \
    adduser -u 1000 -G publisher -D publisher

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/poc-shared-publisher /app/
COPY --from=builder /build/configs/config.example.yaml /app/configs/

# Set ownership
RUN chown -R publisher:publisher /app

# Switch to non-root user
USER publisher

# Expose ports
EXPOSE 8080 8081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# Set entrypoint
ENTRYPOINT ["/app/poc-shared-publisher"]
CMD ["-config", "/app/configs/config.yaml"]
