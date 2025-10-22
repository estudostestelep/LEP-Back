#!/bin/bash

# LEP Remote Database Seeder Script
# This script populates a remote database via HTTP API calls

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
BASE_URL="https://lep-system-516622888070.us-central1.run.app"
ENVIRONMENT="stage"
VERBOSE=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --url)
      BASE_URL="$2"
      shift 2
      ;;
    --environment)
      ENVIRONMENT="$2"
      shift 2
      ;;
    --verbose|-v)
      VERBOSE="--verbose"
      shift
      ;;
    --help|-h)
      echo "Usage: ./scripts/run_seed_remote.sh [OPTIONS]"
      echo ""
      echo "Options:"
      echo "  --url URL              Base URL of the API (default: https://lep-system-516622888070.us-central1.run.app)"
      echo "  --environment ENV      Environment to seed (stage, prod) (default: stage)"
      echo "  --verbose, -v          Enable verbose logging"
      echo "  --help, -h             Show this help message"
      echo ""
      echo "Examples:"
      echo "  ./scripts/run_seed_remote.sh"
      echo "  ./scripts/run_seed_remote.sh --verbose"
      echo "  ./scripts/run_seed_remote.sh --url https://api.example.com --environment prod"
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      echo "Use --help to see available options"
      exit 1
      ;;
  esac
done

# Print header
echo -e "${BLUE}╔══════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   LEP Remote Database Seeder            ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════╝${NC}"
echo ""

# Print configuration
echo -e "${YELLOW}Configuration:${NC}"
echo -e "  Target URL: ${GREEN}${BASE_URL}${NC}"
echo -e "  Environment: ${GREEN}${ENVIRONMENT}${NC}"
echo -e "  Verbose: ${GREEN}${VERBOSE:-false}${NC}"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed or not in PATH${NC}"
    exit 1
fi

# Navigate to project root
cd "$(dirname "$0")/.."

# Build the seed-remote binary
echo -e "${BLUE}Building seed-remote binary...${NC}"
go build -o lep-seed-remote.exe cmd/seed-remote/main.go

if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to build seed-remote binary${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Binary built successfully${NC}"
echo ""

# Run the seeder
echo -e "${BLUE}Starting remote seeding process...${NC}"
echo ""

./lep-seed-remote.exe --url "$BASE_URL" --environment "$ENVIRONMENT" $VERBOSE

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}╔══════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║   Seeding completed successfully!       ║${NC}"
    echo -e "${GREEN}╚══════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${YELLOW}Next steps:${NC}"
    echo -e "  1. Access the application at: ${GREEN}${BASE_URL}${NC}"
    echo -e "  2. Login with credentials shown above"
    echo ""
else
    echo ""
    echo -e "${RED}╔══════════════════════════════════════════╗${NC}"
    echo -e "${RED}║   Seeding failed!                       ║${NC}"
    echo -e "${RED}╚══════════════════════════════════════════╝${NC}"
    exit 1
fi
