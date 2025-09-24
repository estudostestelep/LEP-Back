# Multi-stage build for production optimization
# This Dockerfile is optimized for security and small image size

# Build stage
FROM golang:1.23-alpine AS builder

# Install necessary packages for building
RUN apk add --no-cache git ca-certificates tzdata

# Create non-root user for security
RUN adduser -D -g '' lepuser

# Set working directory
WORKDIR /app

# Copy go modules files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Verify dependencies
RUN go mod verify

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main .

# Final stage - minimal runtime image
FROM scratch

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy SSL certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy user information
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary
COPY --from=builder /app/main /main

# Use non-root user
USER lepuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/main", "-health-check"] || exit 1

# Run the application
ENTRYPOINT ["/main"]

