#!/bin/bash

# LEP System - GCP Bootstrap Script
# Creates base resources manually via gcloud to avoid Terraform permission issues

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
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_success() {
    echo -e "${PURPLE}[SUCCESS]${NC} $1"
}

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# Show banner
show_banner() {
    echo -e "${PURPLE}"
    echo "=================================================================="
    echo "                LEP System - GCP Bootstrap                       "
    echo "=================================================================="
    echo -e "${NC}"
    echo "Environment: $ENVIRONMENT"
    echo "Project: $PROJECT_ID"
    echo "Region: $REGION"
    echo ""
}

# Check prerequisites
check_prerequisites() {
    log_step "Checking prerequisites..."

    # Check gcloud
    if ! command -v gcloud &> /dev/null; then
        log_error "gcloud CLI not found. Install from: https://cloud.google.com/sdk/docs/install"
        exit 1
    fi

    # Check authentication
    if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" &> /dev/null; then
        log_error "gcloud not authenticated. Run: gcloud auth login"
        exit 1
    fi

    # Check project
    local current_project=$(gcloud config get-value project 2>/dev/null)
    if [[ "$current_project" != "$PROJECT_ID" ]]; then
        log_warn "Setting project to $PROJECT_ID"
        gcloud config set project $PROJECT_ID
    fi

    log_success "Prerequisites check passed"
}

# Enable APIs
enable_apis() {
    log_step "Enabling required APIs..."

    local apis=(
        "secretmanager.googleapis.com"
        "sqladmin.googleapis.com"
        "run.googleapis.com"
        "cloudbuild.googleapis.com"
        "artifactregistry.googleapis.com"
        "storage.googleapis.com"
    )

    gcloud services enable \
        secretmanager.googleapis.com \
        sqladmin.googleapis.com \
        run.googleapis.com \
        cloudbuild.googleapis.com \
        artifactregistry.googleapis.com \
        storage.googleapis.com \
        --project=$PROJECT_ID

    log_success "APIs enabled successfully"
}

# Create Service Account
create_service_account() {
    log_step "Creating Service Account..."

    local sa_email="lep-backend-sa@${PROJECT_ID}.iam.gserviceaccount.com"

    if gcloud iam service-accounts describe $sa_email &>/dev/null; then
        log_warn "Service Account already exists: $sa_email"
    else
        gcloud iam service-accounts create lep-backend-sa \
            --display-name="LEP Backend Service Account" \
            --description="Service account for LEP Backend Cloud Run service" \
            --project=$PROJECT_ID

        log_success "Service Account created: $sa_email"
    fi
}

# Create Artifact Registry
create_artifact_registry() {
    log_step "Creating Artifact Registry..."

    if gcloud artifacts repositories describe lep-backend --location=$REGION &>/dev/null; then
        log_warn "Artifact Registry already exists: lep-backend"
    else
        gcloud artifacts repositories create lep-backend \
            --repository-format=docker \
            --location=$REGION \
            --description="LEP Backend Docker repository" \
            --project=$PROJECT_ID

        log_success "Artifact Registry created: lep-backend"
    fi
}

# Create secrets
create_secrets() {
    log_step "Creating Secret Manager secrets..."

    local secrets=(
        "jwt-private-key-${ENVIRONMENT}"
        "jwt-public-key-${ENVIRONMENT}"
        "db-password-${ENVIRONMENT}"
    )

    for secret in "${secrets[@]}"; do
        if gcloud secrets describe $secret &>/dev/null; then
            log_warn "Secret already exists: $secret"
        else
            gcloud secrets create $secret \
                --replication-policy="automatic" \
                --project=$PROJECT_ID

            log_success "Secret created: $secret"
        fi
    done
}

# Create Cloud SQL instance
create_cloud_sql() {
    log_step "Creating Cloud SQL instance..."

    local instance_name="leps-postgres-${ENVIRONMENT}"

    if gcloud sql instances describe $instance_name &>/dev/null; then
        log_warn "Cloud SQL instance already exists: $instance_name"
    else
        log_info "Creating Cloud SQL instance (this may take 10-15 minutes)..."

        gcloud sql instances create $instance_name \
            --database-version=POSTGRES_15 \
            --tier=db-f1-micro \
            --region=$REGION \
            --storage-type=SSD \
            --storage-size=20GB \
            --assign-ip \
            --project=$PROJECT_ID

        log_success "Cloud SQL instance created: $instance_name"
    fi
}

# Create storage bucket for images
create_storage_bucket() {
    log_step "Creating Cloud Storage bucket for images..."

    local bucket_name="${PROJECT_ID}-lep-images-${ENVIRONMENT}"

    if gsutil ls -b gs://$bucket_name &>/dev/null; then
        log_warn "Storage bucket already exists: $bucket_name"
    else
        log_info "Creating storage bucket: $bucket_name"

        # Create bucket with regional storage in same region as other resources
        gsutil mb -p $PROJECT_ID -c STANDARD -l $REGION gs://$bucket_name

        # Set public access for images (read-only)
        gsutil iam ch allUsers:objectViewer gs://$bucket_name

        # Set CORS policy for web access
        cat > /tmp/cors.json << 'EOF'
[
    {
        "origin": ["*"],
        "method": ["GET", "HEAD"],
        "responseHeader": ["Content-Type", "Access-Control-Allow-Origin"],
        "maxAgeSeconds": 3600
    }
]
EOF

        gsutil cors set /tmp/cors.json gs://$bucket_name
        rm -f /tmp/cors.json

        log_info "Bucket created with public read access and CORS configured"
    fi

    # Store bucket name in Secret Manager for application use
    local bucket_secret_name="lep-storage-bucket-name"

    if gcloud secrets describe $bucket_secret_name &>/dev/null; then
        log_info "Updating existing bucket name secret..."
        echo -n "$bucket_name" | gcloud secrets versions add $bucket_secret_name --data-file=-
    else
        log_info "Creating bucket name secret..."
        echo -n "$bucket_name" | gcloud secrets create $bucket_secret_name --data-file=-
    fi

    log_success "Storage bucket created and configured: gs://$bucket_name"
}

# Configure Docker authentication
configure_docker_auth() {
    log_step "Configuring Docker authentication..."

    gcloud auth configure-docker ${REGION}-docker.pkg.dev

    log_success "Docker authentication configured"
}

# Show next steps
show_next_steps() {
    echo ""
    log_success "Bootstrap completed successfully!"
    echo ""
    echo -e "${CYAN}Next steps:${NC}"
    echo "1. Add values to secrets:"
    echo "   gcloud secrets versions add jwt-private-key-${ENVIRONMENT} --data-file=path/to/private.pem"
    echo "   gcloud secrets versions add jwt-public-key-${ENVIRONMENT} --data-file=path/to/public.pem"
    echo ""
    echo "2. Run Terraform:"
    echo "   terraform init"
    echo "   terraform plan -var-file=environments/gcp-${ENVIRONMENT}.tfvars -out=tfplan"
    echo "   terraform apply tfplan"
    echo ""
    echo "3. Or use the interactive deployment script:"
    echo "   ./scripts/deploy-interactive.sh"
}

# Main execution
main() {
    show_banner
    check_prerequisites
    echo ""

    enable_apis
    echo ""

    create_service_account
    echo ""

    create_artifact_registry
    echo ""

    create_secrets
    echo ""

    create_cloud_sql
    echo ""

    create_storage_bucket
    echo ""

    configure_docker_auth

    show_next_steps
}

# Run main function
main "$@"