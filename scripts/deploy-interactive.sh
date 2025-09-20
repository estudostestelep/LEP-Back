#!/bin/bash

# LEP System - Interactive Multi-Environment Deployment Script
# Supports: local-dev, gcp-dev, gcp-stage, gcp-prd

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
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Global variables
SELECTED_ENV=""
FAILED_COMMANDS=()
REMAINING_COMMANDS=()
CURRENT_STEP=0
TOTAL_STEPS=0

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

log_step() {
    echo -e "${BLUE}[STEP $1/$2]${NC} $3"
}

log_success() {
    echo -e "${PURPLE}[SUCCESS]${NC} $1"
}

log_command() {
    echo -e "${CYAN}[CMD]${NC} $1"
}

# Show banner
show_banner() {
    echo -e "${PURPLE}"
    echo "=================================================================="
    echo "           LEP System - Interactive Deployment                    "
    echo "=================================================================="
    echo -e "${NC}"
    echo "Choose your deployment environment:"
    echo ""
    echo "1. 🏠 local-dev     - Local development with Docker Compose"
    echo "2. ☁️  gcp-dev      - GCP minimal setup (testing only)"
    echo "3. 🚀 gcp-stage     - GCP staging (production-like, no Twilio)"
    echo "4. 🌟 gcp-prd       - GCP production (full features)"
    echo ""
}

# Environment selection
select_environment() {
    while true; do
        read -p "Select environment (1-4): " choice
        case $choice in
            1)
                SELECTED_ENV="local-dev"
                break
                ;;
            2)
                SELECTED_ENV="gcp-dev"
                break
                ;;
            3)
                SELECTED_ENV="gcp-stage"
                break
                ;;
            4)
                SELECTED_ENV="gcp-prd"
                break
                ;;
            *)
                log_error "Invalid choice. Please select 1-4."
                ;;
        esac
    done

    log_info "Selected environment: ${SELECTED_ENV}"
}

# Dependency checks
check_dependencies() {
    log_step $((++CURRENT_STEP)) $TOTAL_STEPS "Checking dependencies for ${SELECTED_ENV}"

    local missing_deps=()

    case $SELECTED_ENV in
        "local-dev")
            command -v docker >/dev/null 2>&1 || missing_deps+=("docker")
            command -v docker-compose >/dev/null 2>&1 || missing_deps+=("docker-compose")
            command -v go >/dev/null 2>&1 || missing_deps+=("go")
            ;;
        "gcp-"*)
            command -v gcloud >/dev/null 2>&1 || missing_deps+=("gcloud")
            command -v terraform >/dev/null 2>&1 || missing_deps+=("terraform")
            command -v docker >/dev/null 2>&1 || missing_deps+=("docker")
            command -v go >/dev/null 2>&1 || missing_deps+=("go")
            ;;
    esac

    if [ ${#missing_deps[@]} -ne 0 ]; then
        log_error "Missing required dependencies: ${missing_deps[*]}"
        echo ""
        echo "Installation instructions:"
        for dep in "${missing_deps[@]}"; do
            case $dep in
                "docker")
                    echo "  - Docker: https://docs.docker.com/get-docker/"
                    ;;
                "docker-compose")
                    echo "  - Docker Compose: https://docs.docker.com/compose/install/"
                    ;;
                "gcloud")
                    echo "  - Google Cloud CLI: https://cloud.google.com/sdk/docs/install"
                    ;;
                "terraform")
                    echo "  - Terraform: https://learn.hashicorp.com/tutorials/terraform/install-cli"
                    ;;
                "go")
                    echo "  - Go: https://golang.org/doc/install"
                    ;;
            esac
        done
        exit 1
    fi

    log_success "All dependencies are installed"
}

# Environment variable validation
validate_environment_vars() {
    log_step $((++CURRENT_STEP)) $TOTAL_STEPS "Validating environment variables"

    local missing_vars=()
    local config_file=""

    case $SELECTED_ENV in
        "local-dev")
            config_file="environments/local-dev.env"
            if [ ! -f "$config_file" ]; then
                missing_vars+=("environments/local-dev.env file")
            fi
            ;;
        "gcp-dev")
            config_file="environments/gcp-dev.tfvars"
            ;;
        "gcp-stage")
            config_file="environments/gcp-stage.tfvars"
            # Check for SMTP configuration
            if ! grep -q "smtp_username.*@" "environments/gcp-stage.tfvars" 2>/dev/null; then
                missing_vars+=("SMTP_USERNAME in gcp-stage.tfvars")
            fi
            ;;
        "gcp-prd")
            config_file="environments/gcp-prd.tfvars"
            # Check for Twilio configuration
            if ! grep -q "twilio_account_sid.*AC" "environments/gcp-prd.tfvars" 2>/dev/null; then
                missing_vars+=("TWILIO_ACCOUNT_SID in gcp-prd.tfvars")
            fi
            # Check for SMTP configuration
            if ! grep -q "smtp_username.*@" "environments/gcp-prd.tfvars" 2>/dev/null; then
                missing_vars+=("SMTP_USERNAME in gcp-prd.tfvars")
            fi
            ;;
    esac

    # Check config file exists
    if [ ! -f "$config_file" ]; then
        missing_vars+=("$config_file")
    fi

    # For GCP environments, check JWT keys
    if [[ $SELECTED_ENV == gcp-* ]]; then
        if [ ! -f "jwt_private_key.pem" ] || [ ! -f "jwt_public_key.pem" ]; then
            missing_vars+=("JWT key files (jwt_private_key.pem, jwt_public_key.pem)")
        fi
    fi

    if [ ${#missing_vars[@]} -ne 0 ]; then
        log_error "Missing required configuration: ${missing_vars[*]}"
        echo ""
        echo "Required configurations:"
        for var in "${missing_vars[@]}"; do
            echo "  - $var"
        done
        exit 1
    fi

    log_success "Environment validation passed"
}

# Execute command with error handling
execute_command() {
    local cmd="$1"
    local description="$2"

    log_command "$cmd"

    if eval "$cmd"; then
        log_success "$description completed"
        return 0
    else
        local exit_code=$?
        log_error "$description failed with exit code $exit_code"
        FAILED_COMMANDS+=("$cmd")

        # Show remaining commands
        if [ ${#REMAINING_COMMANDS[@]} -gt 0 ]; then
            echo ""
            log_warn "Remaining commands to complete deployment:"
            for remaining_cmd in "${REMAINING_COMMANDS[@]}"; do
                echo "  - $remaining_cmd"
            done
        fi

        echo ""
        log_error "Deployment failed at: $description"
        log_error "Failed command: $cmd"
        exit $exit_code
    fi
}

# Local development deployment
deploy_local_dev() {
    log_info "🏠 Deploying LOCAL DEVELOPMENT environment"

    REMAINING_COMMANDS=(
        "docker-compose down"
        "docker-compose build"
        "docker-compose up -d"
        "docker-compose logs -f app"
    )

    TOTAL_STEPS=6

    execute_command "docker-compose down -v" "Cleaning existing containers"
    REMAINING_COMMANDS=("${REMAINING_COMMANDS[@]:1}")

    execute_command "cp environments/local-dev.env .env" "Setting up environment file"

    execute_command "docker-compose build --no-cache app" "Building application container"
    REMAINING_COMMANDS=("${REMAINING_COMMANDS[@]:1}")

    execute_command "docker-compose up -d postgres redis mailhog" "Starting infrastructure services"

    log_info "Waiting for database to be ready..."
    sleep 10

    execute_command "docker-compose up -d app" "Starting application"
    REMAINING_COMMANDS=("${REMAINING_COMMANDS[@]:1}")

    log_success "🎉 Local development environment is ready!"
    echo ""
    echo "📍 Available services:"
    echo "  - Application: http://localhost:8080"
    echo "  - Health check: http://localhost:8080/health"
    echo "  - Database (PgAdmin): http://localhost:5050 (admin@lep.local / admin123)"
    echo "  - Email testing (MailHog): http://localhost:8025"
    echo ""
    echo "📋 Useful commands:"
    echo "  - View logs: docker-compose logs -f app"
    echo "  - Stop services: docker-compose down"
    echo "  - Database shell: docker-compose exec postgres psql -U lep_user -d lep_database"
    echo ""

    read -p "Show application logs? (y/n): " show_logs
    if [[ $show_logs =~ ^[Yy]$ ]]; then
        docker-compose logs -f app
    fi
}

# GCP deployment
deploy_gcp() {
    local env_suffix=""
    case $SELECTED_ENV in
        "gcp-dev") env_suffix="dev" ;;
        "gcp-stage") env_suffix="staging" ;;
        "gcp-prd") env_suffix="prod" ;;
    esac

    log_info "☁️ Deploying to GCP environment: ${SELECTED_ENV}"

    REMAINING_COMMANDS=(
        "gcloud auth login"
        "gcloud config set project ${PROJECT_ID}"
        "terraform init"
        "terraform plan -var-file=environments/${SELECTED_ENV}.tfvars"
        "terraform apply -var-file=environments/${SELECTED_ENV}.tfvars"
        "docker build"
        "docker push"
        "gcloud run deploy"
    )

    TOTAL_STEPS=10

    # Authenticate with GCP
    log_step $((++CURRENT_STEP)) $TOTAL_STEPS "Authenticating with Google Cloud"
    if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q "@"; then
        execute_command "gcloud auth login" "Google Cloud authentication"
        execute_command "gcloud auth application-default login" "Application default credentials"
        execute_command "gcloud auth configure-docker ${REGION}-docker.pkg.dev" "Docker authentication"
    else
        log_success "Already authenticated with Google Cloud"
    fi
    REMAINING_COMMANDS=("${REMAINING_COMMANDS[@]:1}")

    # Set project
    execute_command "gcloud config set project ${PROJECT_ID}" "Setting GCP project"
    REMAINING_COMMANDS=("${REMAINING_COMMANDS[@]:1}")

    # Initialize Terraform
    log_step $((++CURRENT_STEP)) $TOTAL_STEPS "Initializing Terraform"
    execute_command "terraform init" "Terraform initialization"
    REMAINING_COMMANDS=("${REMAINING_COMMANDS[@]:1}")

    # Plan Terraform
    log_step $((++CURRENT_STEP)) $TOTAL_STEPS "Planning infrastructure"
    execute_command "terraform plan -var-file=environments/${SELECTED_ENV}.tfvars -out=tfplan" "Terraform planning"
    REMAINING_COMMANDS=("${REMAINING_COMMANDS[@]:1}")

    # Confirm deployment
    echo ""
    log_warn "Review the Terraform plan above."
    read -p "Continue with deployment? (y/n): " confirm
    if [[ ! $confirm =~ ^[Yy]$ ]]; then
        log_warn "Deployment cancelled by user."
        exit 0
    fi

    # Apply Terraform
    log_step $((++CURRENT_STEP)) $TOTAL_STEPS "Creating infrastructure"
    execute_command "terraform apply tfplan" "Infrastructure deployment"
    REMAINING_COMMANDS=("${REMAINING_COMMANDS[@]:1}")

    # Build and push Docker image
    log_step $((++CURRENT_STEP)) $TOTAL_STEPS "Building Docker image"
    local image_url="${REGION}-docker.pkg.dev/${PROJECT_ID}/lep-backend/lep-backend:latest"
    execute_command "docker build -t ${image_url} ." "Docker image build"
    REMAINING_COMMANDS=("${REMAINING_COMMANDS[@]:1}")

    log_step $((++CURRENT_STEP)) $TOTAL_STEPS "Pushing Docker image"
    execute_command "docker push ${image_url}" "Docker image push"
    REMAINING_COMMANDS=("${REMAINING_COMMANDS[@]:1}")

    # Deploy to Cloud Run
    log_step $((++CURRENT_STEP)) $TOTAL_STEPS "Deploying to Cloud Run"
    local service_name="${PROJECT_NAME}-backend-${env_suffix}"
    execute_command "gcloud run deploy ${service_name} --image=${image_url} --region=${REGION} --platform=managed --allow-unauthenticated" "Cloud Run deployment"
    REMAINING_COMMANDS=("${REMAINING_COMMANDS[@]:1}")

    # Get service URL
    local service_url=$(gcloud run services describe ${service_name} --region=${REGION} --format="value(status.url)" 2>/dev/null || echo "")

    log_success "🎉 GCP deployment completed!"
    echo ""
    echo "📍 Deployment information:"
    echo "  - Environment: ${SELECTED_ENV}"
    echo "  - Service: ${service_name}"
    echo "  - Region: ${REGION}"
    if [ -n "$service_url" ]; then
        echo "  - URL: ${service_url}"
        echo "  - Health check: ${service_url}/health"
    fi
    echo ""

    # Health check
    if [ -n "$service_url" ]; then
        log_info "Performing health check..."
        sleep 10
        if curl -s -f "${service_url}/health" > /dev/null; then
            log_success "Application is healthy!"
        else
            log_warn "Health check failed. Service might still be starting up."
        fi
    fi
}

# Main execution
main() {
    show_banner
    select_environment

    # Set total steps based on environment
    case $SELECTED_ENV in
        "local-dev") TOTAL_STEPS=6 ;;
        "gcp-"*) TOTAL_STEPS=10 ;;
    esac

    check_dependencies
    validate_environment_vars

    echo ""
    log_info "🚀 Starting deployment to ${SELECTED_ENV}"
    echo ""

    case $SELECTED_ENV in
        "local-dev")
            deploy_local_dev
            ;;
        "gcp-"*)
            deploy_gcp
            ;;
    esac
}

# Execute main function
main "$@"