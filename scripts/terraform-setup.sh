#!/bin/bash

# LEP System - Terraform Setup Script
# This script initializes and applies the Terraform configuration

set -e

# Configuration
PROJECT_ID="leps-472702"
REGION="us-central1"

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

# Check if required tools are installed
check_dependencies() {
    log_info "Checking dependencies..."

    if ! command -v gcloud &> /dev/null; then
        log_error "gcloud CLI is not installed. Please install it first."
    fi

    if ! command -v terraform &> /dev/null; then
        log_error "Terraform is not installed. Please install it first."
    fi

    log_info "All dependencies are installed."
}

# Authenticate with Google Cloud
authenticate() {
    log_info "Checking Google Cloud authentication..."

    if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q "@"; then
        log_warn "Not authenticated with Google Cloud. Please login."
        gcloud auth login
        gcloud auth application-default login
    fi

    # Set the project
    gcloud config set project ${PROJECT_ID}
    log_info "Using project: ${PROJECT_ID}"
}

# Enable required APIs
enable_apis() {
    log_step "Enabling required Google Cloud APIs..."

    APIs=(
        "secretmanager.googleapis.com"
        "sqladmin.googleapis.com"
        "run.googleapis.com"
        "cloudbuild.googleapis.com"
        "artifactregistry.googleapis.com"
        "compute.googleapis.com"
        "servicenetworking.googleapis.com"
    )

    for api in "${APIs[@]}"; do
        log_info "Enabling ${api}..."
        gcloud services enable ${api} --project=${PROJECT_ID}
    done

    log_info "All APIs enabled successfully."
}

# Initialize Terraform
terraform_init() {
    log_step "Initializing Terraform..."

    if [ ! -f "terraform.tfvars" ]; then
        log_error "terraform.tfvars file not found. Please create it from terraform.tfvars.example"
    fi

    terraform init

    log_info "Terraform initialized successfully."
}

# Plan Terraform changes
terraform_plan() {
    log_step "Planning Terraform changes..."

    terraform plan -out=tfplan

    log_info "Terraform plan completed. Review the changes above."

    read -p "Do you want to apply these changes? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_warn "Deployment cancelled."
        exit 0
    fi
}

# Apply Terraform changes
terraform_apply() {
    log_step "Applying Terraform changes..."

    terraform apply tfplan

    log_info "Terraform applied successfully!"
}

# Show outputs
show_outputs() {
    log_step "Deployment Information:"

    echo ""
    log_info "Service URL:"
    terraform output -raw service_url

    echo ""
    log_info "Database Connection:"
    terraform output -raw database_connection_name

    echo ""
    log_info "Docker Repository:"
    terraform output -raw docker_repository_url

    echo ""
    log_info "Available commands:"
    echo "  Build and deploy: $(terraform output -raw docker_build_command)"
    echo "  Cloud Run deploy: $(terraform output -raw cloud_run_deploy_command)"
}

# Validate configuration
validate_config() {
    log_step "Validating configuration..."

    # Check if JWT keys are properly set
    if grep -q "YOUR_PRIVATE_KEY_HERE" terraform.tfvars; then
        log_error "JWT private key not set in terraform.tfvars. Please update it with your actual key."
    fi

    if grep -q "YOUR_PUBLIC_KEY_HERE" terraform.tfvars; then
        log_error "JWT public key not set in terraform.tfvars. Please update it with your actual key."
    fi

    log_info "Configuration validation passed."
}

# Main execution
main() {
    log_info "Starting LEP System Terraform setup..."

    check_dependencies
    authenticate
    validate_config
    enable_apis
    terraform_init
    terraform_plan
    terraform_apply
    show_outputs

    echo ""
    log_info "Setup completed successfully!"
    log_info "You can now run './scripts/build-and-deploy.sh' to deploy your application."
}

# Execute main function
main "$@"