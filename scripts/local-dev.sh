#!/bin/bash

# LEP System - Local Development Script
# This script helps with local development tasks

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# Show help
show_help() {
    echo "LEP System - Local Development Helper"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  run          Start the Go application locally"
    echo "  build        Build the Go application"
    echo "  test         Run tests"
    echo "  docker       Build and run Docker container locally"
    echo "  clean        Clean build artifacts"
    echo "  deps         Download and verify dependencies"
    echo "  generate     Generate JWT keys"
    echo "  help         Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 run       # Start the application on port 8080"
    echo "  $0 docker    # Build and run in Docker"
    echo "  $0 generate  # Generate new JWT keys"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed. Please install Go first."
    fi
}

# Download dependencies
deps() {
    log_step "Downloading and verifying dependencies..."

    go mod tidy
    go mod verify

    log_info "Dependencies updated successfully."
}

# Build the application
build() {
    log_step "Building the application..."

    check_go

    go build -o bin/lep-system .

    log_info "Build completed successfully. Binary: ./bin/lep-system"
}

# Run the application
run() {
    log_step "Starting the application..."

    check_go

    # Check if .env file exists
    if [ ! -f ".env" ]; then
        log_warn ".env file not found. Creating example .env file..."
        cat > .env << EOF
# Database configuration
DB_USER=postgres
DB_PASS=password
DB_NAME=lep_database

# JWT configuration
JWT_SECRET_PRIVATE_KEY=your_private_key_here
JWT_SECRET_PUBLIC_KEY=your_public_key_here

# Twilio configuration (optional)
TWILIO_ACCOUNT_SID=
TWILIO_AUTH_TOKEN=
TWILIO_PHONE_NUMBER=

# SMTP configuration (optional)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=
SMTP_PASSWORD=

# Application configuration
ENABLE_CRON_JOBS=true
PORT=8080
EOF
        log_warn "Please update the .env file with your actual configuration."
    fi

    log_info "Starting server on http://localhost:8080"
    go run main.go
}

# Run tests
test() {
    log_step "Running tests..."

    check_go

    go test -v ./...

    log_info "Tests completed."
}

# Build and run Docker container
docker() {
    log_step "Building and running Docker container..."

    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed. Please install Docker first."
    fi

    # Build image
    log_info "Building Docker image..."
    docker build -t lep-system:local .

    # Run container
    log_info "Running Docker container..."
    docker run -p 8080:8080 --env-file .env lep-system:local
}

# Clean build artifacts
clean() {
    log_step "Cleaning build artifacts..."

    rm -rf bin/
    rm -f tfplan

    log_info "Clean completed."
}

# Generate JWT keys
generate_keys() {
    log_step "Generating JWT keys..."

    if ! command -v openssl &> /dev/null; then
        log_error "OpenSSL is not installed. Please install OpenSSL first."
    fi

    # Create keys directory if it doesn't exist
    mkdir -p keys

    # Generate private key
    log_info "Generating private key..."
    openssl genpkey -algorithm RSA -out keys/jwt_private_key.pem -pkcs8 -aes256

    # Generate public key
    log_info "Generating public key..."
    openssl rsa -pubout -in keys/jwt_private_key.pem -out keys/jwt_public_key.pem

    log_info "JWT keys generated successfully in ./keys/ directory"
    log_warn "Please update your .env and terraform.tfvars files with the new keys."
}

# Health check
health_check() {
    log_step "Performing health check..."

    if curl -s -f "http://localhost:8080/health" > /dev/null; then
        log_info "Application is healthy!"
    else
        log_warn "Application is not responding or not running."
    fi
}

# Main execution
main() {
    case "${1:-help}" in
        "run")
            run
            ;;
        "build")
            build
            ;;
        "test")
            test
            ;;
        "docker")
            docker
            ;;
        "clean")
            clean
            ;;
        "deps")
            deps
            ;;
        "generate")
            generate_keys
            ;;
        "health")
            health_check
            ;;
        "help"|*)
            show_help
            ;;
    esac
}

# Execute main function
main "$@"