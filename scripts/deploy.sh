#!/bin/bash

# LEP Backend Deployment Script
# This script automates the deployment process to Google Cloud Platform

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
ENVIRONMENT="dev"
SKIP_BUILD=false
SKIP_TERRAFORM=false
DRY_RUN=false

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

Deploy LEP Backend to Google Cloud Platform

OPTIONS:
    -e, --environment ENV    Environment to deploy (dev, staging, prod) [default: dev]
    -s, --skip-build        Skip Docker image build
    -t, --skip-terraform    Skip Terraform infrastructure deployment
    -d, --dry-run           Show what would be deployed without executing
    -h, --help              Show this help message

EXAMPLES:
    $0 --environment prod
    $0 --skip-build --environment staging
    $0 --dry-run --environment prod

PREREQUISITES:
    - gcloud CLI installed and authenticated
    - Docker installed
    - Terraform installed
    - terraform.tfvars file configured
    - JWT keys generated and configured

EOF
}

check_prerequisites() {
    log_info "Checking prerequisites..."

    # Check gcloud
    if ! command -v gcloud &> /dev/null; then
        log_error "gcloud CLI is not installed"
        exit 1
    fi

    # Check Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        exit 1
    fi

    # Check Terraform
    if ! command -v terraform &> /dev/null; then
        log_error "Terraform is not installed"
        exit 1
    fi

    # Check terraform.tfvars
    if [[ ! -f "terraform.tfvars" ]]; then
        log_error "terraform.tfvars file not found"
        log_info "Copy terraform.tfvars.example to terraform.tfvars and configure it"
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

get_terraform_output() {
    local output_name="$1"
    terraform output -raw "$output_name" 2>/dev/null || echo ""
}

deploy_infrastructure() {
    if [[ "$SKIP_TERRAFORM" == "true" ]]; then
        log_info "Skipping Terraform infrastructure deployment"
        return
    fi

    log_info "Deploying infrastructure with Terraform..."

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "DRY RUN: Would run terraform plan and apply"
        return
    fi

    # Initialize Terraform
    terraform init

    # Plan deployment
    terraform plan -var="environment=$ENVIRONMENT" -out=tfplan

    # Apply deployment
    log_info "Applying Terraform configuration..."
    terraform apply tfplan

    # Clean up plan file
    rm -f tfplan

    log_success "Infrastructure deployed successfully"
}

build_and_push_image() {
    if [[ "$SKIP_BUILD" == "true" ]]; then
        log_info "Skipping Docker image build"
        return
    fi

    log_info "Building and pushing Docker image..."

    # Get repository URL from Terraform output
    local docker_repo_url
    docker_repo_url=$(get_terraform_output "docker_repository_url")

    if [[ -z "$docker_repo_url" ]]; then
        log_error "Could not get Docker repository URL from Terraform output"
        log_info "Make sure infrastructure is deployed first"
        exit 1
    fi

    local image_tag="$docker_repo_url/lep-backend:latest"
    local image_tag_env="$docker_repo_url/lep-backend:$ENVIRONMENT-$(date +%Y%m%d-%H%M%S)"

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "DRY RUN: Would build and push:"
        log_info "  - $image_tag"
        log_info "  - $image_tag_env"
        return
    fi

    # Configure Docker to authenticate with GCP
    log_info "Configuring Docker authentication..."
    gcloud auth configure-docker "${docker_repo_url%%/*}" --quiet

    # Build image
    log_info "Building Docker image..."
    docker build -f Dockerfile.prod -t "$image_tag" -t "$image_tag_env" .

    # Push images
    log_info "Pushing Docker images..."
    docker push "$image_tag"
    docker push "$image_tag_env"

    log_success "Docker images built and pushed successfully"
    log_info "Latest image: $image_tag"
    log_info "Tagged image: $image_tag_env"
}

deploy_service() {
    log_info "Deploying Cloud Run service..."

    local service_name
    service_name=$(get_terraform_output "service_name")

    local region
    region=$(get_terraform_output "region")

    local image_url
    image_url=$(get_terraform_output "docker_image_url")

    if [[ -z "$service_name" || -z "$region" || -z "$image_url" ]]; then
        log_error "Could not get service information from Terraform output"
        exit 1
    fi

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "DRY RUN: Would deploy Cloud Run service:"
        log_info "  Service: $service_name"
        log_info "  Region: $region"
        log_info "  Image: $image_url"
        return
    fi

    # Deploy to Cloud Run
    gcloud run deploy "$service_name" \
        --image="$image_url" \
        --region="$region" \
        --platform=managed \
        --quiet

    log_success "Cloud Run service deployed successfully"

    # Get service URL
    local service_url
    service_url=$(get_terraform_output "service_url")
    log_success "Service available at: $service_url"
}

test_deployment() {
    log_info "Testing deployment..."

    local service_url
    service_url=$(get_terraform_output "service_url")

    if [[ -z "$service_url" ]]; then
        log_warning "Could not get service URL for testing"
        return
    fi

    if [[ "$DRY_RUN" == "true" ]]; then
        log_info "DRY RUN: Would test endpoints:"
        log_info "  - $service_url/health"
        log_info "  - $service_url/ping"
        return
    fi

    # Test health endpoint
    log_info "Testing health endpoint..."
    if curl -sf "$service_url/health" > /dev/null; then
        log_success "Health check passed"
    else
        log_warning "Health check failed"
    fi

    # Test ping endpoint
    log_info "Testing ping endpoint..."
    if curl -sf "$service_url/ping" > /dev/null; then
        log_success "Ping check passed"
    else
        log_warning "Ping check failed"
    fi
}

main() {
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -e|--environment)
                ENVIRONMENT="$2"
                shift 2
                ;;
            -s|--skip-build)
                SKIP_BUILD=true
                shift
                ;;
            -t|--skip-terraform)
                SKIP_TERRAFORM=true
                shift
                ;;
            -d|--dry-run)
                DRY_RUN=true
                shift
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
    if [[ ! "$ENVIRONMENT" =~ ^(dev|staging|prod)$ ]]; then
        log_error "Invalid environment: $ENVIRONMENT"
        log_info "Valid environments: dev, staging, prod"
        exit 1
    fi

    log_info "Starting deployment for environment: $ENVIRONMENT"

    if [[ "$DRY_RUN" == "true" ]]; then
        log_warning "DRY RUN MODE - No actual changes will be made"
    fi

    check_prerequisites
    deploy_infrastructure
    build_and_push_image
    deploy_service
    test_deployment

    log_success "Deployment completed successfully!"

    if [[ "$DRY_RUN" != "true" ]]; then
        local service_url
        service_url=$(get_terraform_output "service_url")
        log_info "Service URL: $service_url"
    fi
}

# Execute main function with all arguments
main "$@"