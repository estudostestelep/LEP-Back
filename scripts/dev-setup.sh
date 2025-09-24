#!/bin/bash

# LEP Development Setup Script
# This script sets up the complete development environment

set -e

echo ""
echo "🚀 LEP Development Setup"
echo "========================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Configuration
SKIP_TESTS=false
SKIP_SEED=false
START_SERVER=true
CLEAR_DATA=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-tests)
            SKIP_TESTS=true
            shift
            ;;
        --skip-seed)
            SKIP_SEED=true
            shift
            ;;
        --no-server)
            START_SERVER=false
            shift
            ;;
        --clear-data)
            CLEAR_DATA=true
            shift
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "This script sets up the complete LEP development environment:"
            echo "  1. Checks dependencies"
            echo "  2. Builds the application"
            echo "  3. Runs tests"
            echo "  4. Seeds the database"
            echo "  5. Starts the server"
            echo ""
            echo "Options:"
            echo "  --skip-tests     Skip running tests"
            echo "  --skip-seed      Skip database seeding"
            echo "  --no-server      Don't start the server"
            echo "  --clear-data     Clear existing data before seeding"
            echo "  --help           Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0                    # Full setup"
            echo "  $0 --skip-tests       # Setup without tests"
            echo "  $0 --clear-data       # Fresh setup with clean data"
            exit 0
            ;;
        *)
            echo -e "${RED}❌ Unknown option: $1${NC}"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

echo -e "${BLUE}📋 Setup Configuration:${NC}"
echo "  Skip tests: $SKIP_TESTS"
echo "  Skip seed: $SKIP_SEED"
echo "  Start server: $START_SERVER"
echo "  Clear data: $CLEAR_DATA"
echo ""

# Step 1: Check dependencies
echo -e "${BLUE}1. 📦 Checking Dependencies${NC}"
echo "==============================="

# Check Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed${NC}"
    echo "Please install Go from: https://golang.org/dl/"
    exit 1
fi
echo -e "${GREEN}✓ Go $(go version | cut -d' ' -f3)${NC}"

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo -e "${RED}❌ Please run this script from the LEP-Back root directory${NC}"
    exit 1
fi
echo -e "${GREEN}✓ LEP-Back directory${NC}"

# Check Docker (optional)
if command -v docker &> /dev/null; then
    echo -e "${GREEN}✓ Docker available${NC}"
    DOCKER_AVAILABLE=true
else
    echo -e "${YELLOW}⚠️  Docker not available (optional)${NC}"
    DOCKER_AVAILABLE=false
fi

echo ""

# Step 2: Environment setup
echo -e "${BLUE}2. ⚙️  Environment Setup${NC}"
echo "=========================="

# Check .env file
if [ ! -f ".env" ]; then
    echo -e "${YELLOW}⚠️  No .env file found. Creating from example...${NC}"
    if [ -f ".env.example" ]; then
        cp .env.example .env
        echo -e "${BLUE}ℹ️  Please update .env with your database credentials${NC}"
    else
        echo -e "${YELLOW}⚠️  No .env.example found, continuing...${NC}"
    fi
fi

# Download dependencies
echo "📥 Downloading Go dependencies..."
if go mod tidy; then
    echo -e "${GREEN}✓ Dependencies updated${NC}"
else
    echo -e "${YELLOW}⚠️  Warning: Some dependencies failed to download${NC}"
fi

echo ""

# Step 3: Build application
echo -e "${BLUE}3. 🔨 Building Application${NC}"
echo "==========================="

if go build -o lep-system .; then
    echo -e "${GREEN}✓ Application built successfully${NC}"
    echo "  Binary: ./lep-system"
else
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi

echo ""

# Step 4: Run tests
if [ "$SKIP_TESTS" = false ]; then
    echo -e "${BLUE}4. 🧪 Running Tests${NC}"
    echo "===================="

    if [ -d "tests" ]; then
        echo "Running API test suite..."
        if go test ./tests -v; then
            echo -e "${GREEN}✓ All tests passed${NC}"
        else
            echo -e "${YELLOW}⚠️  Some tests failed, continuing...${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  No tests directory found, skipping...${NC}"
    fi

    echo ""
else
    echo -e "${YELLOW}⏭️  Skipping tests${NC}"
    echo ""
fi

# Step 5: Database setup
if [ "$SKIP_SEED" = false ]; then
    echo -e "${BLUE}5. 🌱 Database Setup${NC}"
    echo "===================="

    # Check if seeder exists
    if [ -f "cmd/seed/main.go" ]; then
        SEED_ARGS="--environment=dev"
        if [ "$CLEAR_DATA" = true ]; then
            SEED_ARGS="$SEED_ARGS --clear-first"
        fi

        echo "Seeding database with sample data..."
        if go run cmd/seed/main.go $SEED_ARGS; then
            echo -e "${GREEN}✓ Database seeded successfully${NC}"
        else
            echo -e "${YELLOW}⚠️  Database seeding failed, continuing...${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  No database seeder found, skipping...${NC}"
    fi

    echo ""
else
    echo -e "${YELLOW}⏭️  Skipping database seeding${NC}"
    echo ""
fi

# Step 6: Start server
if [ "$START_SERVER" = true ]; then
    echo -e "${BLUE}6. 🌐 Starting Server${NC}"
    echo "======================"

    echo -e "${GREEN}🎉 LEP Development Environment Ready!${NC}"
    echo ""
    echo -e "${PURPLE}📊 Available Endpoints:${NC}"
    echo "  Health: http://localhost:8080/health"
    echo "  Ping:   http://localhost:8080/ping"
    echo "  API:    http://localhost:8080/*"
    echo ""
    echo -e "${PURPLE}👤 Login Credentials:${NC}"
    echo "  Admin:   admin@lep-demo.com / password"
    echo "  Waiter:  garcom@lep-demo.com / password"
    echo "  Manager: gerente@lep-demo.com / password"
    echo ""
    echo -e "${PURPLE}🔧 Development Commands:${NC}"
    echo "  Run tests:     go test ./tests -v"
    echo "  Reseed DB:     ./scripts/run_seed.sh --clear-first"
    echo "  Check health:  curl http://localhost:8080/health"
    echo ""
    echo -e "${BLUE}🚀 Starting LEP Backend Server...${NC}"
    echo "Press Ctrl+C to stop"
    echo ""

    # Start the server
    if [ -f "./lep-system" ]; then
        ./lep-system
    else
        go run main.go
    fi
else
    echo -e "${GREEN}🎉 LEP Development Environment Ready!${NC}"
    echo ""
    echo -e "${BLUE}Manual Start Commands:${NC}"
    echo "  ./lep-system"
    echo "  # or"
    echo "  go run main.go"
    echo ""
    echo -e "${PURPLE}📊 Available Endpoints (when server is running):${NC}"
    echo "  Health: http://localhost:8080/health"
    echo "  Ping:   http://localhost:8080/ping"
    echo "  API:    http://localhost:8080/*"
    echo ""
    echo -e "${PURPLE}👤 Login Credentials:${NC}"
    echo "  Admin:   admin@lep-demo.com / password"
    echo "  Waiter:  garcom@lep-demo.com / password"
    echo "  Manager: gerente@lep-demo.com / password"
fi