#!/bin/bash

# LEP Database Seeder for Fattoria Pizzeria
# This script populates the database with Fattoria Pizzeria sample data

set -e

echo ""
echo "🍕 LEP Database Seeder - Fattoria Pizzeria"
echo "==========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
CLEAR_FIRST=false
ENVIRONMENT="dev"
VERBOSE=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --clear-first)
            CLEAR_FIRST=true
            shift
            ;;
        --environment)
            ENVIRONMENT="$2"
            shift 2
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --clear-first    Clear existing data before seeding"
            echo "  --environment    Environment to seed (dev, test, staging) [default: dev]"
            echo "  --verbose        Enable verbose logging"
            echo "  --help           Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0                                    # Basic seeding with Fattoria data"
            echo "  $0 --clear-first --verbose            # Clear and verbose seed"
            echo "  $0 --environment=test                 # Seed test environment"
            exit 0
            ;;
        *)
            echo -e "${RED}❌ Unknown option: $1${NC}"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed or not in PATH${NC}"
    echo "Please install Go from: https://golang.org/dl/"
    exit 1
fi

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -d "cmd/seed" ]; then
    echo -e "${RED}❌ Please run this script from the LEP-Back root directory${NC}"
    exit 1
fi

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo -e "${YELLOW}⚠️  No .env file found. Creating from example...${NC}"
    if [ -f ".env.example" ]; then
        cp .env.example .env
        echo -e "${BLUE}ℹ️  Please update .env with your database credentials${NC}"
    else
        echo -e "${RED}❌ No .env.example file found${NC}"
        exit 1
    fi
fi

# Build arguments for the seeder
ARGS="--restaurant=fattoria"
if [ "$CLEAR_FIRST" = true ]; then
    ARGS="$ARGS --clear-first"
fi
if [ "$VERBOSE" = true ]; then
    ARGS="$ARGS --verbose"
fi
ARGS="$ARGS --environment=$ENVIRONMENT"

echo -e "${BLUE}📋 Configuration:${NC}"
echo "  Restaurant: Fattoria Pizzeria 🍕"
echo "  Environment: $ENVIRONMENT"
echo "  Clear first: $CLEAR_FIRST"
echo "  Verbose: $VERBOSE"
echo ""

# Download dependencies if needed
echo -e "${BLUE}📦 Checking dependencies...${NC}"
if ! go mod verify &> /dev/null; then
    echo "  Downloading dependencies..."
    go mod tidy
    if [ $? -ne 0 ]; then
        echo -e "${YELLOW}⚠️  Warning: Failed to download dependencies, continuing...${NC}"
    fi
fi

# Run the seeder
echo -e "${BLUE}🚀 Running Fattoria database seeder...${NC}"
echo ""

# Execute the seeder
if go run cmd/seed/main.go $ARGS; then
    echo ""
    echo -e "${GREEN}✅ Fattoria database seeding completed successfully!${NC}"
    echo ""
    echo -e "${BLUE}🎯 Next Steps:${NC}"
    echo "  1. Start the backend server:"
    echo -e "     ${YELLOW}go run main.go${NC}"
    echo ""
    echo "  2. Test the API:"
    echo -e "     ${YELLOW}curl http://localhost:8080/health${NC}"
    echo ""
    echo "  3. Login credentials:"
    echo -e "     ${YELLOW}admin@fattoria.com.br / password (Admin)${NC}"
    echo ""
    echo "  4. Fattoria Menu Data:"
    echo -e "     ${YELLOW}🍕 Pizzas:${NC}"
    echo -e "        • Crostini (Entradas) - R$ 30,00"
    echo -e "        • Marguerita (Vegetariana) - R$ 80,00"
    echo -e "        • Marinara (Vegana) - R$ 58,00"
    echo -e "        • Parma - R$ 109,00"
    echo -e "        • Vegana (Vegana) - R$ 60,00"
    echo -e "     ${YELLOW}🥤 Bebidas:${NC}"
    echo -e "        • Suco de caju integral - R$ 15,00"
    echo -e "        • Heineken s/ álcool - R$ 13,00"
    echo -e "        • Baden Baden IPA - R$ 23,00"
    echo -e "        • Sônia e Zé (Coquetel) - R$ 32,00"
    echo ""
    echo "  5. Run tests:"
    echo -e "     ${YELLOW}go test ./tests -v${NC}"
    echo ""
else
    echo ""
    echo -e "${RED}❌ Fattoria database seeding failed!${NC}"
    echo ""
    echo -e "${YELLOW}💡 Troubleshooting:${NC}"
    echo "  1. Check if PostgreSQL is running"
    echo "  2. Verify database credentials in .env"
    echo "  3. Ensure database exists and is accessible"
    echo "  4. Check network connectivity"
    echo ""
    echo "  For more details, run with --verbose flag"
    exit 1
fi
