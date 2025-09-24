#!/bin/bash

# LEP System - Quick Deploy Script (Hybrid Approach)
# Resolves Terraform conflicts and deploys cleanly

set -e

# Configuration
PROJECT_ID="leps-472702"
REGION="us-central1"
ENVIRONMENT="${ENVIRONMENT:-dev}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_success() { echo -e "${PURPLE}[SUCCESS]${NC} $1"; }
log_step() { echo -e "${BLUE}[STEP]${NC} $1"; }

show_banner() {
    echo -e "${PURPLE}"
    echo "=================================================================="
    echo "              LEP System - Quick Deploy (Hybrid)                 "
    echo "=================================================================="
    echo -e "${NC}"
    echo "Environment: $ENVIRONMENT"
    echo "Project: $PROJECT_ID"
    echo ""
}

# Clean Terraform conflicts
clean_terraform() {
    log_step "Cleaning Terraform conflicts..."

    # Remove any lock files
    rm -f terraform.tfstate.lock.info 2>/dev/null || true
    rm -f .terraform.lock.hcl 2>/dev/null || true
    rm -f tfplan 2>/dev/null || true

    # Backup original main.tf if it exists
    if [[ -f "main_original.tf" ]]; then
        log_info "Original main.tf already backed up"
    elif [[ -f "main.tf" ]] && grep -q "google_artifact_registry_repository" main.tf; then
        log_info "Backing up original main.tf..."
        mv main.tf main_original.tf
    fi

    # Ensure we're using the simplified version
    if [[ ! -f "main.tf" ]] && [[ -f "main_simplified.tf" ]]; then
        log_info "Using simplified main.tf..."
        cp main_simplified.tf main.tf
    fi

    log_success "Terraform conflicts cleaned"
}

# Check if bootstrap is needed
check_bootstrap() {
    log_step "Checking if bootstrap is needed..."

    local sa_email="lep-backend-sa@${PROJECT_ID}.iam.gserviceaccount.com"

    if gcloud iam service-accounts describe $sa_email &>/dev/null; then
        log_success "Bootstrap resources already exist"
        return 0
    else
        log_warn "Bootstrap needed - Service Account not found"
        return 1
    fi
}

# Quick bootstrap
quick_bootstrap() {
    log_step "Running quick bootstrap..."

    # Enable APIs
    log_info "Enabling APIs..."
    gcloud services enable \
        secretmanager.googleapis.com \
        sqladmin.googleapis.com \
        run.googleapis.com \
        cloudbuild.googleapis.com \
        artifactregistry.googleapis.com \
        --project=$PROJECT_ID

    # Create Service Account
    log_info "Creating Service Account..."
    if ! gcloud iam service-accounts describe lep-backend-sa@${PROJECT_ID}.iam.gserviceaccount.com &>/dev/null; then
        gcloud iam service-accounts create lep-backend-sa \
            --display-name="LEP Backend Service Account" \
            --project=$PROJECT_ID
    fi

    # Create Artifact Registry
    log_info "Creating Artifact Registry..."
    if ! gcloud artifacts repositories describe lep-backend --location=$REGION &>/dev/null; then
        gcloud artifacts repositories create lep-backend \
            --repository-format=docker \
            --location=$REGION \
            --project=$PROJECT_ID
    fi

    # Create Secrets
    log_info "Creating secrets..."
    local secrets=("jwt-private-key-${ENVIRONMENT}" "jwt-public-key-${ENVIRONMENT}" "db-password-${ENVIRONMENT}")
    for secret in "${secrets[@]}"; do
        if ! gcloud secrets describe $secret &>/dev/null; then
            gcloud secrets create $secret --replication-policy="automatic" --project=$PROJECT_ID
        fi
    done

    # Configure Docker
    gcloud auth configure-docker ${REGION}-docker.pkg.dev

    log_success "Bootstrap completed"
}

# Deploy with Terraform
deploy_terraform() {
    log_step "Deploying with Terraform..."

    # Initialize
    log_info "Initializing Terraform..."
    if ! terraform init; then
        log_error "Terraform init failed"
        return 1
    fi

    # Plan
    log_info "Creating Terraform plan..."
    if ! terraform plan -var-file=environments/gcp-${ENVIRONMENT}.tfvars -out=tfplan; then
        log_error "Terraform plan failed"
        return 1
    fi

    # Apply
    log_info "Applying Terraform plan..."
    if ! terraform apply tfplan; then
        log_error "Terraform apply failed"
        return 1
    fi

    log_success "Terraform deployment completed"
}

# Build and deploy application
build_and_deploy() {
    log_step "Building and deploying application..."

    local image_tag="${REGION}-docker.pkg.dev/${PROJECT_ID}/lep-backend/lep-backend:latest"

    # Build
    log_info "Building Docker image..."
    if ! docker build -t $image_tag .; then
        log_error "Docker build failed"
        return 1
    fi

    # Push
    log_info "Pushing Docker image..."
    if ! docker push $image_tag; then
        log_error "Docker push failed"
        return 1
    fi

    log_success "Application built and deployed"
}

# Get service URL
get_service_url() {
    log_step "Getting service URL..."

    local service_name="leps-backend-${ENVIRONMENT}"
    local service_url=""

    # Wait a bit for service to be ready
    sleep 5

    service_url=$(gcloud run services describe $service_name \
        --region=$REGION \
        --format="value(status.url)" 2>/dev/null || echo "")

    if [[ -n "$service_url" ]]; then
        log_success "Service URL: $service_url"
        log_info "Health check: $service_url/health"
        log_info "API docs: $service_url/ping"
    else
        log_warn "Service URL not available yet. Check Cloud Run console."
    fi
}

# Add sample secrets (if needed)
add_sample_secrets() {
    log_step "Adding sample secrets..."

    # Add sample JWT keys (you should replace these with real ones)
    local jwt_private="jwt-private-key-${ENVIRONMENT}"
    local jwt_public="jwt-public-key-${ENVIRONMENT}"
    local db_password="db-password-${ENVIRONMENT}"

    # Check if secrets have versions
    if ! gcloud secrets versions list $jwt_private --limit=1 &>/dev/null; then
        log_info "Adding sample JWT private key..."
        echo "sample_jwt_private_key_replace_with_real_one" | gcloud secrets versions add $jwt_private --data-file=-
    fi

    if ! gcloud secrets versions list $jwt_public --limit=1 &>/dev/null; then
        log_info "Adding sample JWT public key..."
        echo "sample_jwt_public_key_replace_with_real_one" | gcloud secrets versions add $jwt_public --data-file=-
    fi

    if ! gcloud secrets versions list $db_password --limit=1 &>/dev/null; then
        log_info "Adding sample DB password..."
        echo "$(openssl rand -base64 32)" | gcloud secrets versions add $db_password --data-file=-
    fi

    log_success "Sample secrets added"
}

# Main execution
main() {
    show_banner

    # Check authentication
    if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" &>/dev/null; then
        log_error "Please authenticate with: gcloud auth login"
        exit 1
    fi

    # Set project
    gcloud config set project $PROJECT_ID

    echo ""
    clean_terraform
    echo ""

    if ! check_bootstrap; then
        echo ""
        quick_bootstrap
        echo ""
        add_sample_secrets
    fi

    echo ""
    deploy_terraform
    echo ""
    build_and_deploy
    echo ""
    get_service_url

    echo ""
    log_success "🎉 Deployment completed successfully!"
    echo ""
    echo -e "${BLUE}Next steps:${NC}"
    echo "1. Update JWT secrets with real keys:"
    echo "   gcloud secrets versions add jwt-private-key-${ENVIRONMENT} --data-file=path/to/private.pem"
    echo "   gcloud secrets versions add jwt-public-key-${ENVIRONMENT} --data-file=path/to/public.pem"
    echo ""
    echo "2. Test the deployment:"
    echo "   curl \$(gcloud run services describe leps-backend-${ENVIRONMENT} --region=${REGION} --format='value(status.url)')/health"
}

# Run main function
main "$@"