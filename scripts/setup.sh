#!/bin/bash

# LEP System - Complete Environment Setup Script
# This script sets up the entire development and deployment environment

set -e

# Configuration
PROJECT_ID="leps-472702"
PROJECT_NAME="leps"
REGION="us-central1"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
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

log_success() {
    echo -e "${PURPLE}[SUCCESS]${NC} $1"
}

# Show banner
show_banner() {
    echo -e "${PURPLE}"
    echo "=================================="
    echo "   LEP System Environment Setup   "
    echo "=================================="
    echo -e "${NC}"
    echo "This script will set up your complete LEP System environment."
    echo "Project ID: ${PROJECT_ID}"
    echo "Region: ${REGION}"
    echo ""
}

# Check if required tools are installed
check_dependencies() {
    log_step "Checking system dependencies..."

    local missing_tools=()

    # Check for required tools
    if ! command -v gcloud &> /dev/null; then
        missing_tools+=("gcloud")
    fi

    if ! command -v terraform &> /dev/null; then
        missing_tools+=("terraform")
    fi

    if ! command -v docker &> /dev/null; then
        missing_tools+=("docker")
    fi

    if ! command -v go &> /dev/null; then
        missing_tools+=("go")
    fi

    if ! command -v git &> /dev/null; then
        missing_tools+=("git")
    fi

    if [ ${#missing_tools[@]} -ne 0 ]; then
        log_error "Missing required tools: ${missing_tools[*]}"
        echo ""
        echo "Please install the missing tools:"
        echo "- gcloud: https://cloud.google.com/sdk/docs/install"
        echo "- terraform: https://learn.hashicorp.com/tutorials/terraform/install-cli"
        echo "- docker: https://docs.docker.com/get-docker/"
        echo "- go: https://golang.org/doc/install"
        echo "- git: https://git-scm.com/downloads"
        exit 1
    fi

    log_success "All required tools are installed."
}

# Setup Google Cloud
setup_gcloud() {
    log_step "Setting up Google Cloud..."

    # Check authentication
    if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q "@"; then
        log_warn "Not authenticated with Google Cloud. Please login."
        gcloud auth login
        gcloud auth application-default login
    fi

    # Set project
    gcloud config set project ${PROJECT_ID}

    # Verify project access
    if ! gcloud projects describe ${PROJECT_ID} &> /dev/null; then
        log_error "Cannot access project ${PROJECT_ID}. Please check your permissions."
    fi

    log_success "Google Cloud setup completed."
}

# Create directory structure
setup_directories() {
    log_step "Creating project directory structure..."

    # Create necessary directories
    mkdir -p bin
    mkdir -p logs
    mkdir -p scripts
    mkdir -p keys
    mkdir -p docs

    log_success "Directory structure created."
}

# Setup Go dependencies
setup_go() {
    log_step "Setting up Go dependencies..."

    # Ensure we're in the right directory
    if [ ! -f "go.mod" ]; then
        log_error "go.mod not found. Please run this script from the project root."
    fi

    # Download and verify dependencies
    go mod tidy
    go mod verify

    # Build the application to verify everything works
    log_info "Building application to verify setup..."
    go build -o bin/lep-system .

    log_success "Go setup completed successfully."
}

# Setup environment files
setup_environment() {
    log_step "Setting up environment configuration..."

    # Create .env file if it doesn't exist
    if [ ! -f ".env" ]; then
        log_info "Creating .env file..."
        cat > .env << EOF
# Database configuration
DB_USER=postgres
DB_PASS=your_database_password
DB_NAME=lep_database

# JWT configuration (will be populated after key generation)
JWT_SECRET_PRIVATE_KEY=
JWT_SECRET_PUBLIC_KEY=

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

    log_success "Environment configuration setup completed."
}

# Generate JWT keys if they don't exist
setup_jwt_keys() {
    log_step "Setting up JWT keys..."

    if [ ! -f "jwt_private_key.pem" ] || [ ! -f "jwt_public_key.pem" ]; then
        log_info "JWT keys not found. They should already exist in your project."
        log_info "Checking for existing keys..."

        if [ -f "jwt_private_key.pem" ] && [ -f "jwt_public_key.pem" ]; then
            log_success "JWT keys found and ready to use."
        else
            log_warn "JWT keys not found. Please ensure you have the jwt_private_key.pem and jwt_public_key.pem files."
        fi
    else
        log_success "JWT keys are already configured."
    fi
}

# Setup Docker
setup_docker() {
    log_step "Setting up Docker..."

    # Check if Docker is running
    if ! docker info &> /dev/null; then
        log_error "Docker is not running. Please start Docker first."
    fi

    # Configure Docker for Google Cloud
    log_info "Configuring Docker for Google Cloud..."
    gcloud auth configure-docker ${REGION}-docker.pkg.dev

    log_success "Docker setup completed."
}

# Validate configuration
validate_setup() {
    log_step "Validating setup..."

    local errors=0

    # Check if terraform.tfvars exists and is configured
    if [ ! -f "terraform.tfvars" ]; then
        log_warn "terraform.tfvars not found."
        errors=$((errors + 1))
    else
        if grep -q "YOUR_PRIVATE_KEY_HERE" terraform.tfvars; then
            log_warn "JWT keys not configured in terraform.tfvars."
            errors=$((errors + 1))
        fi
    fi

    # Check if Go application builds
    if [ ! -f "bin/lep-system" ]; then
        log_warn "Application binary not found. Build may have failed."
        errors=$((errors + 1))
    fi

    if [ $errors -eq 0 ]; then
        log_success "All validations passed!"
    else
        log_warn "Found $errors configuration issues. Please address them before deploying."
    fi
}

# Show next steps
show_next_steps() {
    echo ""
    log_step "Setup completed! Next steps:"
    echo ""
    echo "1. Update configuration files:"
    echo "   - .env: Add your database and service credentials"
    echo "   - terraform.tfvars: Verify JWT keys and other settings"
    echo ""
    echo "2. Initialize infrastructure:"
    echo "   ./scripts/terraform-setup.sh"
    echo ""
    echo "3. Deploy your application:"
    echo "   ./scripts/build-and-deploy.sh"
    echo ""
    echo "4. For local development:"
    echo "   ./scripts/local-dev.sh run"
    echo ""
    echo "5. Useful commands:"
    echo "   - Health check: curl http://localhost:8080/health"
    echo "   - View logs: ./scripts/local-dev.sh logs"
    echo "   - Clean build: ./scripts/local-dev.sh clean"
    echo ""
    log_success "LEP System is ready to use!"
}

# Main execution
main() {
    show_banner
    check_dependencies
    setup_gcloud
    setup_directories
    setup_go
    setup_environment
    setup_jwt_keys
    setup_docker
    validate_setup
    show_next_steps
}

# Handle script interruption
trap 'log_error "Setup interrupted by user."' INT TERM

# Execute main function
main "$@"