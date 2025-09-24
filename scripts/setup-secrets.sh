#!/bin/bash

# LEP Backend Secrets Setup Script
# This script helps generate and configure secrets for the LEP backend

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

show_usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Setup secrets for LEP Backend deployment

OPTIONS:
    -g, --generate-jwt      Generate new JWT key pair
    -p, --project-id ID     GCP Project ID
    -e, --environment ENV   Environment (dev, staging, prod) [default: dev]
    -h, --help              Show this help message

EXAMPLES:
    $0 --generate-jwt --project-id my-project-123
    $0 --project-id my-project-123 --environment prod

PREREQUISITES:
    - openssl installed (for JWT key generation)
    - gcloud CLI installed and authenticated

EOF
}

check_prerequisites() {
    log_info "Checking prerequisites..."

    # Check openssl
    if ! command -v openssl &> /dev/null; then
        log_error "openssl is not installed"
        exit 1
    fi

    # Check gcloud
    if ! command -v gcloud &> /dev/null; then
        log_error "gcloud CLI is not installed"
        exit 1
    fi

    # Check if user is authenticated with gcloud
    if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
        log_error "Not authenticated with gcloud"
        log_info "Run: gcloud auth login"
        exit 1
    fi

    log_success "All prerequisites checked"
}

generate_jwt_keys() {
    log_info "Generating JWT key pair..."

    local private_key_file="jwt_private_key.pem"
    local public_key_file="jwt_public_key.pem"

    # Generate private key
    openssl genpkey -algorithm RSA -out "$private_key_file" -pkcs8 -aes256

    # Generate public key
    openssl rsa -pubout -in "$private_key_file" -out "$public_key_file"

    log_success "JWT keys generated:"
    log_info "Private key: $private_key_file"
    log_info "Public key: $public_key_file"

    log_warning "IMPORTANT: Store these keys securely and add them to your terraform.tfvars file"
    log_warning "The private key is encrypted and you'll need the passphrase to use it"

    echo ""
    echo "Add these to your terraform.tfvars file:"
    echo "jwt_private_key = \"\"\""
    cat "$private_key_file"
    echo "\"\"\""
    echo ""
    echo "jwt_public_key = \"\"\""
    cat "$public_key_file"
    echo "\"\"\""
    echo ""
}

setup_terraform_vars() {
    local project_id="$1"
    local environment="$2"

    log_info "Setting up terraform.tfvars file..."

    if [[ -f "terraform.tfvars" ]]; then
        log_warning "terraform.tfvars already exists"
        read -p "Do you want to backup and recreate it? (y/N): " -r
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            mv terraform.tfvars "terraform.tfvars.backup.$(date +%Y%m%d-%H%M%S)"
            log_info "Existing terraform.tfvars backed up"
        else
            log_info "Keeping existing terraform.tfvars file"
            return
        fi
    fi

    # Copy from example
    if [[ -f "terraform.tfvars.example" ]]; then
        cp terraform.tfvars.example terraform.tfvars

        # Update project_id and environment
        sed -i.bak "s/your-gcp-project-id/$project_id/g" terraform.tfvars
        sed -i.bak "s/environment  = \"dev\"/environment  = \"$environment\"/g" terraform.tfvars
        rm -f terraform.tfvars.bak

        log_success "terraform.tfvars created from example"
        log_warning "Please edit terraform.tfvars to add your specific values:"
        log_info "- JWT keys (generate with --generate-jwt)"
        log_info "- Twilio credentials (optional)"
        log_info "- SMTP credentials (optional)"
    else
        log_error "terraform.tfvars.example not found"
        exit 1
    fi
}

validate_gcp_project() {
    local project_id="$1"

    log_info "Validating GCP project: $project_id"

    if ! gcloud projects describe "$project_id" > /dev/null 2>&1; then
        log_error "Project $project_id not found or not accessible"
        log_info "Make sure you have access to the project and it exists"
        exit 1
    fi

    # Set current project
    gcloud config set project "$project_id"

    log_success "GCP project validated and set as current"
}

enable_required_apis() {
    local project_id="$1"

    log_info "Enabling required APIs..."

    local apis=(
        "secretmanager.googleapis.com"
        "sqladmin.googleapis.com"
        "run.googleapis.com"
        "cloudbuild.googleapis.com"
        "artifactregistry.googleapis.com"
    )

    for api in "${apis[@]}"; do
        log_info "Enabling $api..."
        gcloud services enable "$api" --project="$project_id"
    done

    log_success "All required APIs enabled"
}

main() {
    local generate_jwt=false
    local project_id=""
    local environment="dev"

    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -g|--generate-jwt)
                generate_jwt=true
                shift
                ;;
            -p|--project-id)
                project_id="$2"
                shift 2
                ;;
            -e|--environment)
                environment="$2"
                shift 2
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done

    # Validate environment
    if [[ ! "$environment" =~ ^(dev|staging|prod)$ ]]; then
        log_error "Invalid environment: $environment"
        log_info "Valid environments: dev, staging, prod"
        exit 1
    fi

    check_prerequisites

    if [[ "$generate_jwt" == "true" ]]; then
        generate_jwt_keys
    fi

    if [[ -n "$project_id" ]]; then
        validate_gcp_project "$project_id"
        enable_required_apis "$project_id"
        setup_terraform_vars "$project_id" "$environment"

        log_success "Setup completed for project: $project_id"
        log_info "Next steps:"
        log_info "1. Edit terraform.tfvars with your specific values"
        log_info "2. Run: terraform init"
        log_info "3. Run: terraform plan"
        log_info "4. Run: terraform apply"
        log_info "5. Run: ./scripts/deploy.sh --environment $environment"
    else
        log_warning "No project ID provided. Use --project-id to set up GCP configuration"
    fi
}

# Execute main function with all arguments
main "$@"