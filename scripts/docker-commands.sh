#!/bin/bash

# LEP Docker Commands Helper
# This script provides easy commands for Docker-based development

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

show_help() {
    echo ""
    echo "🐳 LEP Docker Commands Helper"
    echo "============================="
    echo ""
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Available commands:"
    echo ""
    echo "🚀 Development:"
    echo "  start              Start all development services (postgres, redis, mailhog, app)"
    echo "  stop               Stop all services"
    echo "  restart            Restart all services"
    echo "  logs               Show logs from all services"
    echo "  logs-app           Show logs from app only"
    echo ""
    echo "🌱 Database:"
    echo "  seed               Populate database with sample data"
    echo "  seed-fresh         Clear database and seed with fresh data"
    echo "  db-shell           Connect to PostgreSQL shell"
    echo ""
    echo "🧪 Testing:"
    echo "  test               Run tests in container"
    echo "  test-coverage      Run tests with coverage report"
    echo ""
    echo "🔧 Tools:"
    echo "  pgadmin            Start PgAdmin (database management)"
    echo "  build              Build application container"
    echo "  clean              Remove all containers and volumes"
    echo "  status             Show status of all services"
    echo ""
    echo "Examples:"
    echo "  $0 start                    # Start development environment"
    echo "  $0 seed                     # Populate database"
    echo "  $0 test                     # Run tests"
    echo "  $0 logs-app                 # Watch app logs"
    echo ""
}

# Check if docker and docker-compose are available
check_docker() {
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}❌ Docker is not installed or not in PATH${NC}"
        echo "Please install Docker from: https://docs.docker.com/get-docker/"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        echo -e "${RED}❌ Docker Compose is not installed or not in PATH${NC}"
        echo "Please install Docker Compose from: https://docs.docker.com/compose/install/"
        exit 1
    fi
}

# Check if we're in the right directory
check_directory() {
    if [ ! -f "docker-compose.yml" ]; then
        echo -e "${RED}❌ Please run this script from the LEP-Back root directory${NC}"
        exit 1
    fi
}

# Main command handler
case "$1" in
    "start")
        check_docker
        check_directory
        echo -e "${BLUE}🚀 Starting LEP development environment...${NC}"
        docker-compose up -d postgres redis mailhog
        echo -e "${YELLOW}⏳ Waiting for services to be ready...${NC}"
        sleep 10
        docker-compose up -d app
        echo -e "${GREEN}✅ Development environment started!${NC}"
        echo ""
        echo -e "${PURPLE}📊 Available services:${NC}"
        echo "  App:      http://localhost:8080"
        echo "  MailHog:  http://localhost:8025"
        echo ""
        echo -e "${BLUE}📝 Next steps:${NC}"
        echo "  • Run: $0 seed              # Populate database"
        echo "  • Run: $0 logs-app          # Watch app logs"
        echo "  • Run: $0 pgadmin           # Database management"
        ;;

    "stop")
        check_docker
        check_directory
        echo -e "${YELLOW}🛑 Stopping LEP services...${NC}"
        docker-compose down
        echo -e "${GREEN}✅ All services stopped!${NC}"
        ;;

    "restart")
        check_docker
        check_directory
        echo -e "${YELLOW}🔄 Restarting LEP services...${NC}"
        docker-compose restart
        echo -e "${GREEN}✅ Services restarted!${NC}"
        ;;

    "logs")
        check_docker
        check_directory
        echo -e "${BLUE}📋 Showing logs from all services (Ctrl+C to exit):${NC}"
        docker-compose logs -f
        ;;

    "logs-app")
        check_docker
        check_directory
        echo -e "${BLUE}📋 Showing app logs (Ctrl+C to exit):${NC}"
        docker-compose logs -f app
        ;;

    "seed")
        check_docker
        check_directory
        echo -e "${BLUE}🌱 Seeding database with sample data...${NC}"
        docker-compose run --rm seed
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✅ Database seeded successfully!${NC}"
        else
            echo -e "${RED}❌ Database seeding failed!${NC}"
        fi
        ;;

    "seed-fresh")
        check_docker
        check_directory
        echo -e "${BLUE}🌱 Clearing database and seeding with fresh data...${NC}"
        docker-compose run --rm seed go run cmd/seed/main.go --environment=dev --clear-first --verbose
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✅ Database seeded successfully with fresh data!${NC}"
        else
            echo -e "${RED}❌ Database seeding failed!${NC}"
        fi
        ;;

    "db-shell")
        check_docker
        check_directory
        echo -e "${BLUE}🐘 Connecting to PostgreSQL shell...${NC}"
        docker-compose exec postgres psql -U lep_user -d lep_database
        ;;

    "test")
        check_docker
        check_directory
        echo -e "${BLUE}🧪 Running tests in container...${NC}"
        docker-compose run --rm test
        ;;

    "test-coverage")
        check_docker
        check_directory
        echo -e "${BLUE}🧪 Running tests with coverage...${NC}"
        docker-compose run --rm test go test ./tests -v -cover -coverprofile=coverage.out
        ;;

    "pgadmin")
        check_docker
        check_directory
        echo -e "${BLUE}🔧 Starting PgAdmin...${NC}"
        docker-compose --profile tools up -d pgadmin
        echo -e "${GREEN}✅ PgAdmin started!${NC}"
        echo ""
        echo -e "${PURPLE}📊 PgAdmin access:${NC}"
        echo "  URL:      http://localhost:5050"
        echo "  Email:    admin@lep.local"
        echo "  Password: admin123"
        echo ""
        echo -e "${BLUE}📝 Database connection:${NC}"
        echo "  Host:     postgres (or localhost)"
        echo "  Port:     5432"
        echo "  Database: lep_database"
        echo "  Username: lep_user"
        echo "  Password: lep_password"
        ;;

    "build")
        check_docker
        check_directory
        echo -e "${BLUE}🔨 Building application container...${NC}"
        docker-compose build app
        echo -e "${GREEN}✅ Application container built!${NC}"
        ;;

    "clean")
        check_docker
        check_directory
        echo -e "${YELLOW}🧹 Cleaning up containers and volumes...${NC}"
        read -p "This will remove all containers and volumes. Are you sure? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            docker-compose down -v --remove-orphans
            docker system prune -f
            echo -e "${GREEN}✅ Cleanup completed!${NC}"
        else
            echo -e "${BLUE}ℹ️  Cleanup cancelled${NC}"
        fi
        ;;

    "status")
        check_docker
        check_directory
        echo -e "${BLUE}📊 LEP Services Status:${NC}"
        echo "========================"
        docker-compose ps
        echo ""
        echo -e "${BLUE}🐳 Docker System Info:${NC}"
        echo "======================"
        docker system df
        ;;

    "help"|"--help"|"-h"|"")
        show_help
        ;;

    *)
        echo -e "${RED}❌ Unknown command: $1${NC}"
        echo "Use '$0 help' to see available commands"
        exit 1
        ;;
esac