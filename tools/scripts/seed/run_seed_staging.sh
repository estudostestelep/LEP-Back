#!/bin/bash

# 🌱 LEP Seed Runner - Staging Environment
# Este script executa o seed para o ambiente de staging

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
ENV_FILE="${PROJECT_DIR}/.env.staging"
RESTAURANT="fattoria"
CLEAR_FIRST=false
VERBOSE=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --clear-first)
            CLEAR_FIRST=true
            shift
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        --restaurant)
            RESTAURANT="$2"
            shift 2
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Function to print colored output
print_header() {
    echo -e "${BLUE}================================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}================================================${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Main execution
print_header "🌱 LEP Database Seeder - STAGING ENVIRONMENT"

# Check if .env.staging exists
if [ ! -f "$ENV_FILE" ]; then
    print_error ".env.staging file not found at $ENV_FILE"
    exit 1
fi

print_success "Found .env.staging configuration"

# Load environment variables from .env.staging
export $(grep -v '^#' "$ENV_FILE" | xargs)

# Validate required variables
required_vars=("DB_USER" "DB_PASS" "DB_NAME" "ENVIRONMENT")
for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        print_error "Missing required variable: $var"
        exit 1
    fi
done

print_success "Environment variables loaded"
echo ""
echo -e "${BLUE}Configuration:${NC}"
echo "  Environment: $ENVIRONMENT"
echo "  Database: $DB_NAME"
echo "  Database User: $DB_USER"
echo "  Restaurant: $RESTAURANT"
echo "  Clear First: $CLEAR_FIRST"
echo "  Verbose: $VERBOSE"
echo ""

# Build seed command
cd "$PROJECT_DIR"

SEED_CMD="go run ./cmd/seed/ --restaurant=$RESTAURANT --environment=$ENVIRONMENT"

if [ "$CLEAR_FIRST" = true ]; then
    SEED_CMD="$SEED_CMD --clear-first"
fi

if [ "$VERBOSE" = true ]; then
    SEED_CMD="$SEED_CMD --verbose"
fi

# Check if we can connect to the database
print_header "🔍 Validating Database Connection"

if [ -n "$INSTANCE_UNIX_SOCKET" ]; then
    print_warning "Using Cloud SQL Unix Socket: $INSTANCE_UNIX_SOCKET"
    print_warning "Ensure Cloud SQL Proxy is running or you have direct access"
else
    print_success "Using TCP connection to $DB_HOST:$DB_PORT"
fi

# Run the seed
print_header "🌱 Running Seed Execution"
echo ""
echo -e "${BLUE}Command:${NC} $SEED_CMD"
echo ""

if eval "$SEED_CMD"; then
    print_success "Seed execution completed successfully!"
    echo ""
    print_header "📊 Next Steps"
    echo "1. Verify data in staging database"
    echo "2. Run: gcloud sql connect <instance-name> --user=lep_user"
    echo "3. Test API endpoint: https://staging-api.lep.example.com/health"
    echo ""
    print_success "Staging seed ready!"
else
    print_error "Seed execution failed!"
    exit 1
fi
