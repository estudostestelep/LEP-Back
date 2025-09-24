#!/bin/bash

# LEP System - Health Check Script
# Verifica status de todos os ambientes e recursos

set -e

# Configuration
PROJECT_ID="leps-472702"
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
CHECK_ALL=false
ENVIRONMENT=""
DETAILED=false

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

log_check() {
    echo -e "${CYAN}[CHECK]${NC} $1"
}

# Show banner
show_banner() {
    echo -e "${PURPLE}"
    echo "=================================================================="
    echo "                LEP System - Health Check                        "
    echo "=================================================================="
    echo -e "${NC}"
}

# Show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --all               Check all environments"
    echo "  --environment ENV   Check specific environment (local-dev, gcp-dev, gcp-stage, gcp-prd)"
    echo "  --detailed          Show detailed information"
    echo "  --help              Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 --all                          # Check all environments"
    echo "  $0 --environment gcp-dev          # Check only gcp-dev"
    echo "  $0 --environment gcp-prd --detailed  # Detailed check for production"
}

# Parse arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --all)
                CHECK_ALL=true
                shift
                ;;
            --environment)
                ENVIRONMENT="$2"
                shift 2
                ;;
            --detailed)
                DETAILED=true
                shift
                ;;
            --help)
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
}

# Check prerequisites
check_prerequisites() {
    log_check "Checking prerequisites..."

    local all_good=true

    # Check gcloud
    if ! command -v gcloud &> /dev/null; then
        log_error "gcloud CLI not found"
        all_good=false
    else
        log_success "gcloud CLI found"
    fi

    # Check docker
    if ! command -v docker &> /dev/null; then
        log_warn "Docker not found (needed for local builds)"
    else
        log_success "Docker found"
    fi

    # Check go
    if ! command -v go &> /dev/null; then
        log_warn "Go not found (needed for local development)"
    else
        log_success "Go found: $(go version | cut -d' ' -f3)"
    fi

    # Check gcloud auth
    if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" &> /dev/null; then
        log_error "gcloud not authenticated"
        all_good=false
    else
        local account=$(gcloud auth list --filter=status:ACTIVE --format="value(account)")
        log_success "gcloud authenticated as: $account"
    fi

    # Check project
    local current_project=$(gcloud config get-value project 2>/dev/null)
    if [[ "$current_project" != "$PROJECT_ID" ]]; then
        log_warn "Current project: $current_project, expected: $PROJECT_ID"
    else
        log_success "Project correctly set: $PROJECT_ID"
    fi

    if [[ "$all_good" == false ]]; then
        log_error "Prerequisites check failed"
        exit 1
    fi
}

# Check local development environment
check_local_dev() {
    log_check "Checking local development environment..."

    # Check if local server is running
    if curl -s http://localhost:8080/ping &> /dev/null; then
        log_success "Local server is running on port 8080"

        # Check health endpoint
        local health_response=$(curl -s http://localhost:8080/health 2>/dev/null || echo "failed")
        if [[ "$health_response" == *"healthy"* ]]; then
            log_success "Local health check passed"
        else
            log_error "Local health check failed: $health_response"
        fi

        if [[ "$DETAILED" == true ]]; then
            echo "Health response: $health_response"
        fi
    else
        log_warn "Local server not running on port 8080"

        # Check if Docker Compose is running
        if docker-compose ps 2>/dev/null | grep -q "Up"; then
            log_info "Docker Compose services detected"
            if [[ "$DETAILED" == true ]]; then
                docker-compose ps
            fi
        else
            log_info "Local server not running. Use './run-local.sh' or 'docker-compose up'"
        fi
    fi
}

# Check GCP environment
check_gcp_environment() {
    local env=$1
    log_check "Checking GCP environment: $env..."

    local service_name=""
    case $env in
        "gcp-dev")
            service_name="leps-backend-dev"
            ;;
        "gcp-stage")
            service_name="leps-backend-stage"
            ;;
        "gcp-prd")
            service_name="leps-backend-prd"
            ;;
        *)
            log_error "Unknown GCP environment: $env"
            return 1
            ;;
    esac

    # Check if Cloud Run service exists
    if gcloud run services describe "$service_name" --region="$REGION" &> /dev/null; then
        log_success "Cloud Run service '$service_name' exists"

        # Get service URL
        local service_url=$(gcloud run services describe "$service_name" \
            --region="$REGION" \
            --format="value(status.url)" 2>/dev/null)

        if [[ -n "$service_url" ]]; then
            log_success "Service URL: $service_url"

            # Check health endpoint
            local health_response=$(curl -s "$service_url/ping" 2>/dev/null || echo "failed")
            if [[ "$health_response" == *"pong"* ]]; then
                log_success "Service health check passed"

                # Detailed health check
                local detailed_health=$(curl -s "$service_url/health" 2>/dev/null || echo "failed")
                if [[ "$detailed_health" == *"healthy"* ]]; then
                    log_success "Detailed health check passed"
                    if [[ "$DETAILED" == true ]]; then
                        echo "Health response: $detailed_health"
                    fi
                else
                    log_warn "Detailed health check failed: $detailed_health"
                fi
            else
                log_error "Service health check failed: $health_response"
            fi
        else
            log_error "Could not get service URL"
        fi

        # Check service status
        if [[ "$DETAILED" == true ]]; then
            log_info "Service details:"
            gcloud run services describe "$service_name" \
                --region="$REGION" \
                --format="table(metadata.name,status.conditions[0].type,status.conditions[0].status,status.traffic[0].percent)"
        fi

    else
        log_warn "Cloud Run service '$service_name' not found"
        log_info "Use './scripts/deploy-interactive.sh' to deploy"
    fi
}

# Check GCP resources
check_gcp_resources() {
    log_check "Checking GCP resources..."

    # Check Cloud SQL instances
    local sql_instances=$(gcloud sql instances list --format="value(name)" 2>/dev/null | grep -c "leps" || echo "0")
    if [[ "$sql_instances" -gt 0 ]]; then
        log_success "Found $sql_instances Cloud SQL instance(s)"
        if [[ "$DETAILED" == true ]]; then
            gcloud sql instances list --format="table(name,region,databaseVersion,settings.tier,status)"
        fi
    else
        log_warn "No Cloud SQL instances found"
    fi

    # Check Artifact Registry repositories
    local repos=$(gcloud artifacts repositories list --location="$REGION" --format="value(name)" 2>/dev/null | grep -c "lep" || echo "0")
    if [[ "$repos" -gt 0 ]]; then
        log_success "Found $repos Artifact Registry repositor(y/ies)"
        if [[ "$DETAILED" == true ]]; then
            gcloud artifacts repositories list --location="$REGION" --format="table(name,format,createTime)"
        fi
    else
        log_warn "No Artifact Registry repositories found"
    fi

    # Check secrets
    local secrets=$(gcloud secrets list --format="value(name)" 2>/dev/null | grep -c "lep\|jwt\|db-password" || echo "0")
    if [[ "$secrets" -gt 0 ]]; then
        log_success "Found $secrets secret(s)"
        if [[ "$DETAILED" == true ]]; then
            gcloud secrets list --format="table(name,createTime)" | grep -E "lep|jwt|db-password"
        fi
    else
        log_warn "No secrets found"
    fi

    # Check service accounts
    local service_accounts=$(gcloud iam service-accounts list --format="value(email)" 2>/dev/null | grep -c "lep" || echo "0")
    if [[ "$service_accounts" -gt 0 ]]; then
        log_success "Found $service_accounts service account(s)"
        if [[ "$DETAILED" == true ]]; then
            gcloud iam service-accounts list --format="table(email,displayName)" | grep "lep"
        fi
    else
        log_warn "No service accounts found"
    fi
}

# Check specific environment
check_environment() {
    local env=$1

    case $env in
        "local-dev")
            check_local_dev
            ;;
        "gcp-dev"|"gcp-stage"|"gcp-prd")
            check_gcp_environment "$env"
            ;;
        *)
            log_error "Unknown environment: $env"
            exit 1
            ;;
    esac
}

# Main execution
main() {
    show_banner

    # Parse arguments
    parse_args "$@"

    # Check prerequisites
    check_prerequisites

    echo ""

    # Determine what to check
    if [[ "$CHECK_ALL" == true ]]; then
        log_info "Checking all environments..."
        echo ""

        # Check GCP resources first
        check_gcp_resources
        echo ""

        # Check each environment
        check_environment "local-dev"
        echo ""
        check_environment "gcp-dev"
        echo ""
        check_environment "gcp-stage"
        echo ""
        check_environment "gcp-prd"

    elif [[ -n "$ENVIRONMENT" ]]; then
        log_info "Checking environment: $ENVIRONMENT"
        echo ""

        if [[ "$ENVIRONMENT" =~ ^gcp- ]]; then
            check_gcp_resources
            echo ""
        fi

        check_environment "$ENVIRONMENT"

    else
        log_error "Please specify --all or --environment <env>"
        show_usage
        exit 1
    fi

    echo ""
    log_success "Health check completed!"
}

# Run main function
main "$@"