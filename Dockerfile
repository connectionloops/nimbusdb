# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build arguments
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X 'NimbusDb/version.Version=${VERSION}' -X 'NimbusDb/version.Commit=${COMMIT}' -X 'NimbusDb/version.BuildDate=${BUILD_DATE}'" \
    -o nimbusdb \
    .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS, tzdata for timezone support, and wget for health checks
RUN apk --no-cache add ca-certificates tzdata wget

# Create non-root user
RUN addgroup -g 1000 nimbusdb && \
    adduser -D -u 1000 -G nimbusdb nimbusdb

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/nimbusdb /app/nimbusdb

# Create health check script
RUN echo '#!/bin/sh' > /app/healthcheck.sh && \
    echo 'PORT=${HEALTH_PORT:-8080}' >> /app/healthcheck.sh && \
    echo 'wget --no-verbose --tries=1 --spider "http://localhost:${PORT}/health" || exit 1' >> /app/healthcheck.sh && \
    chmod +x /app/healthcheck.sh

# Change ownership to non-root user
RUN chown -R nimbusdb:nimbusdb /app

# Switch to non-root user
USER nimbusdb

# Expose health check port (default 8080, can be overridden via HEALTH_PORT env var)
EXPOSE 8080

# Health check using HTTP endpoint (uses HEALTH_PORT env var, defaults to 8080)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD /app/healthcheck.sh

# Run the application
CMD ["./nimbusdb"]

