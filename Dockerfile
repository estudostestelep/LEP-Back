# Multi-stage build for production optimization
# This Dockerfile is optimized for security and small image size

# Build stage
FROM golang:1.24-alpine AS builder

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

# Final stage - minimal runtime image with shell support for Cloud Run
FROM alpine:latest

# Install CA certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN adduser -D -g '' lepuser

# Set working directory
WORKDIR /

# Copy the binary
COPY --from=builder /app/main /main

# Make binary executable
RUN chmod +x /main

# Use non-root user
USER lepuser

# Expose port (Cloud Run will override this with PORT env var)
EXPOSE 8080

# Run the application (no health check as Cloud Run handles this)
ENTRYPOINT ["/main"]

