#!/bin/bash

# LEP Test Runner Script
# This script runs all tests with proper formatting and reporting

set -e

echo ""
echo "🧪 LEP Test Suite Runner"
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
VERBOSE=false
COVERAGE=false
HTML_REPORT=false
SPECIFIC_TEST=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --verbose|-v)
            VERBOSE=true
            shift
            ;;
        --coverage|-c)
            COVERAGE=true
            shift
            ;;
        --html)
            HTML_REPORT=true
            COVERAGE=true
            shift
            ;;
        --test|-t)
            SPECIFIC_TEST="$2"
            shift 2
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  -v, --verbose      Enable verbose test output"
            echo "  -c, --coverage     Generate coverage report"
            echo "  --html             Generate HTML coverage report"
            echo "  -t, --test NAME    Run specific test"
            echo "  -h, --help         Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0                         # Run all tests"
            echo "  $0 --verbose --coverage    # Verbose with coverage"
            echo "  $0 --html                  # Generate HTML report"
            echo "  $0 --test TestUserRoutes   # Run specific test"
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
if [ ! -f "go.mod" ]; then
    echo -e "${RED}❌ Please run this script from the LEP-Back root directory${NC}"
    exit 1
fi

# Check if tests directory exists
if [ ! -d "tests" ]; then
    echo -e "${RED}❌ No tests directory found${NC}"
    echo "Please create tests using: mkdir tests"
    exit 1
fi

echo -e "${BLUE}📋 Test Configuration:${NC}"
echo "  Verbose: $VERBOSE"
echo "  Coverage: $COVERAGE"
echo "  HTML Report: $HTML_REPORT"
if [ -n "$SPECIFIC_TEST" ]; then
    echo "  Specific Test: $SPECIFIC_TEST"
fi
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

# Prepare test arguments
TEST_ARGS=""
if [ "$VERBOSE" = true ]; then
    TEST_ARGS="$TEST_ARGS -v"
fi

if [ "$COVERAGE" = true ]; then
    TEST_ARGS="$TEST_ARGS -cover -coverprofile=coverage.out"
fi

if [ -n "$SPECIFIC_TEST" ]; then
    TEST_ARGS="$TEST_ARGS -run $SPECIFIC_TEST"
fi

# Create reports directory if it doesn't exist
mkdir -p tests/reports

# Record start time
START_TIME=$(date +%s)

echo -e "${BLUE}🚀 Running Tests...${NC}"
echo ""

# Run the tests
echo -e "${BLUE}🧪 Running all test suites...${NC}"
echo "  • Unit tests (data generation, validation)"
echo "  • Integration tests (route availability)"
echo "  • Real integration tests (database operations)"
echo "  • Business scenario tests (complete workflows)"
echo ""

if go test ./tests $TEST_ARGS; then
    # Calculate duration
    END_TIME=$(date +%s)
    DURATION=$((END_TIME - START_TIME))

    echo ""
    echo -e "${GREEN}✅ All tests passed!${NC}"
    echo -e "${BLUE}⏱️  Duration: ${DURATION}s${NC}"

    # Generate coverage reports if requested
    if [ "$COVERAGE" = true ] && [ -f "coverage.out" ]; then
        echo ""
        echo -e "${BLUE}📊 Coverage Analysis:${NC}"

        # Show coverage summary
        go tool cover -func=coverage.out | tail -1

        if [ "$HTML_REPORT" = true ]; then
            echo "📄 Generating HTML coverage report..."
            go tool cover -html=coverage.out -o tests/reports/coverage.html
            echo -e "${GREEN}✓ HTML report: tests/reports/coverage.html${NC}"

            # Try to open in browser (if available)
            if command -v xdg-open &> /dev/null; then
                echo "🌐 Opening in browser..."
                xdg-open tests/reports/coverage.html &
            elif command -v open &> /dev/null; then
                echo "🌐 Opening in browser..."
                open tests/reports/coverage.html &
            fi
        fi
    fi

    echo ""
    echo -e "${PURPLE}🎯 Test Summary:${NC}"
    echo "================="

    # Count test files and functions
    TEST_FILES=$(find tests -name "*_test.go" | wc -l)
    TEST_FUNCS=$(grep -r "func Test" tests/ | wc -l)

    echo "  Test files: $TEST_FILES"
    echo "  Test functions: $TEST_FUNCS"

    if [ -f "coverage.out" ]; then
        COVERAGE_PERCENT=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
        echo "  Coverage: $COVERAGE_PERCENT"
    fi

    echo ""
    echo -e "${GREEN}🎉 Testing completed successfully!${NC}"

    # Cleanup coverage file if not needed
    if [ "$COVERAGE" = false ] && [ -f "coverage.out" ]; then
        rm coverage.out
    fi

else
    END_TIME=$(date +%s)
    DURATION=$((END_TIME - START_TIME))

    echo ""
    echo -e "${RED}❌ Tests failed!${NC}"
    echo -e "${BLUE}⏱️  Duration: ${DURATION}s${NC}"
    echo ""
    echo -e "${YELLOW}💡 Troubleshooting:${NC}"
    echo "  1. Check test output above for specific failures"
    echo "  2. Ensure all dependencies are installed: go mod tidy"
    echo "  3. Verify test data and mocks are correct"
    echo "  4. Run with --verbose for detailed output"
    echo ""

    # Cleanup coverage file
    if [ -f "coverage.out" ]; then
        rm coverage.out
    fi

    exit 1
fi