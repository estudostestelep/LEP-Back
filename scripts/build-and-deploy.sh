#!/bin/bash

# LEP System - Build and Deploy Script
# This script builds the Docker image and deploys to Google Cloud Run

set -e

# Configuration
PROJECT_ID="leps-472702"
REGION="us-central1"
SERVICE_NAME="lep-system"
REPOSITORY="lep-backend"
IMAGE_TAG="latest"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
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

# Check if required tools are installed
check_dependencies() {
    log_info "Checking dependencies..."

    if ! command -v gcloud &> /dev/null; then
        log_error "gcloud CLI is not installed. Please install it first."
    fi

    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed. Please install it first."
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
        gcloud auth configure-docker ${REGION}-docker.pkg.dev
    fi

    # Set the project
    gcloud config set project ${PROJECT_ID}
    log_info "Using project: ${PROJECT_ID}"
}

# Build Docker image
build_image() {
    log_info "Building Docker image..."

    # Use main Dockerfile (production-optimized)
    IMAGE_URL="${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/lep-backend:${IMAGE_TAG}"

    docker build -t ${IMAGE_URL} .

    log_info "Image built successfully: ${IMAGE_URL}"
}

# Push image to Artifact Registry
push_image() {
    log_info "Pushing image to Artifact Registry..."

    IMAGE_URL="${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/lep-backend:${IMAGE_TAG}"

    docker push ${IMAGE_URL}

    log_info "Image pushed successfully to Artifact Registry."
}

# Deploy to Cloud Run
deploy_service() {
    log_info "Deploying to Cloud Run..."

    IMAGE_URL="${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/lep-backend:${IMAGE_TAG}"

    # The service should already be created by Terraform
    # This command updates the existing service with the new image
    gcloud run deploy ${SERVICE_NAME} \
        --image=${IMAGE_URL} \
        --region=${REGION} \
        --platform=managed \
        --allow-unauthenticated \
        --port=8080 \
        --timeout=300 \
        --memory=512Mi \
        --cpu=1 \
        --max-instances=10 \
        --min-instances=0

    log_info "Service deployed successfully!"

    # Get service URL
    SERVICE_URL=$(gcloud run services describe ${SERVICE_NAME} --region=${REGION} --format="value(status.url)")
    log_info "Service URL: ${SERVICE_URL}"
}

# Health check
health_check() {
    log_info "Performing health check..."

    SERVICE_URL=$(gcloud run services describe ${SERVICE_NAME} --region=${REGION} --format="value(status.url)" 2>/dev/null)

    if [ -z "$SERVICE_URL" ]; then
        log_warn "Could not get service URL. Skipping health check."
        return
    fi

    # Wait a moment for the service to be ready
    sleep 10

    # Check health endpoint
    if curl -s -f "${SERVICE_URL}/health" > /dev/null; then
        log_info "Health check passed!"
    else
        log_warn "Health check failed. Service might still be starting up."
    fi
}

# Main execution
main() {
    log_info "Starting LEP System deployment..."

    check_dependencies
    authenticate
    build_image
    push_image
    deploy_service
    health_check

    log_info "Deployment completed successfully!"

    # Show service information
    echo ""
    log_info "Service Information:"
    gcloud run services describe ${SERVICE_NAME} --region=${REGION} --format="table(metadata.name,status.url,status.conditions[0].type,status.conditions[0].status)"
}

# Execute main function
main "$@"