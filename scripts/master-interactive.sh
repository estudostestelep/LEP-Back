#!/bin/bash

# LEP System - Master Interactive Script
# Unified interface for all LEP System operations and deployments
# Updated for new environment structure: dev (local) | stage (GCP) | prod (future)

set -e

# Global Configuration
PROJECT_ID="leps-472702"
PROJECT_NAME="leps"
REGION="us-central1"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors and formatting
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
NC='\033[0m' # No Color

# Logging functions
log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_success() { echo -e "${PURPLE}[SUCCESS]${NC} $1"; }
log_step() { echo -e "${BLUE}[STEP]${NC} $1"; }
log_header() { echo -e "${CYAN}[$(echo $1 | tr '[:lower:]' '[:upper:]')]${NC} $2"; }

# Banner and branding
show_main_banner() {
    clear
    echo -e "${PURPLE}"
    echo "=================================================================="
    echo "           LEP System - Master Control (Refatorado)              "
    echo "=================================================================="
    echo -e "${NC}"
    echo -e "${WHITE}Project:${NC} ${PROJECT_NAME} (${PROJECT_ID})"
    echo -e "${WHITE}Region:${NC} ${REGION}"
    echo -e "${WHITE}Ambientes:${NC} dev (local) | stage (GCP) | prod (futuro)"
    echo ""
}

# Main menu display
show_main_menu() {
    echo -e "${CYAN}🎛️  Selecione uma categoria:${NC}"
    echo ""
    echo "  1. 🔧 Ambiente DEV (Local)"
    echo "  2. 🚀 Ambiente STAGE (GCP)"
    echo "  3. 🌱 Database & Seeding"
    echo "  4. 🧪 Testes"
    echo "  5. ⚙️  Setup & Configuração"
    echo "  6. 🛠️  Utilitários"
    echo "  7. ❓ Ajuda"
    echo "  0. 🚪 Sair"
    echo ""
}

# Utility functions
press_enter() {
    echo ""
    echo -e "${YELLOW}Pressione ENTER para continuar...${NC}"
    read
}

confirm_action() {
    local message="$1"
    echo ""
    echo -e "${YELLOW}$message${NC}"
    read -p "Continuar? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_warn "Operação cancelada pelo usuário."
        return 1
    fi
    return 0
}

check_dependencies() {
    local tool="$1"
    if ! command -v "$tool" &> /dev/null; then
        log_error "$tool não está instalado ou não está no PATH"
        return 1
    fi
    return 0
}

# ==============================================================================
# 1. AMBIENTE DEV (LOCAL)
# ==============================================================================

show_dev_menu() {
    show_main_banner
    log_header "dev" "Ambiente DEV - 100% Local"
    echo ""
    echo -e "${CYAN}🏠 Docker + localStorage + credenciais padronizadas${NC}"
    echo ""
    echo "  1. 🚀 Iniciar ambiente completo (Docker)"
    echo "  2. 🔨 Build aplicação Go"
    echo "  3. 💊 Health check local"
    echo "  4. 📊 Status do ambiente dev"
    echo "  5. 🧹 Parar e limpar ambiente"
    echo "  6. 🌱 Popular dados demo"
    echo "  0. ⬅️  Voltar ao menu principal"
    echo ""
}

handle_dev_menu() {
    while true; do
        show_dev_menu
        read -p "Selecione uma opção: " choice

        case $choice in
            1) dev_start_environment ;;
            2) dev_build_app ;;
            3) dev_health_check ;;
            4) dev_status ;;
            5) dev_stop_and_clean ;;
            6) dev_seed_data ;;
            0) return ;;
            *) log_error "Opção inválida. Tente novamente." ; press_enter ;;
        esac
    done
}

dev_start_environment() {
    log_step "Iniciando ambiente DEV completo..."
    cd "$ROOT_DIR"

    if ! check_dependencies "docker"; then
        press_enter
        return
    fi

    if ! docker info &> /dev/null; then
        log_error "Docker não está rodando. Inicie o Docker primeiro."
        press_enter
        return
    fi

    log_info "Executando script dev-local.sh..."
    if [ -f "scripts/dev-local.sh" ]; then
        chmod +x scripts/dev-local.sh
        ./scripts/dev-local.sh
    else
        log_error "Script dev-local.sh não encontrado"
    fi
    press_enter
}

dev_build_app() {
    log_step "Build da aplicação Go..."
    cd "$ROOT_DIR"

    if ! check_dependencies "go"; then
        press_enter
        return
    fi

    mkdir -p bin
    if go build -o bin/lep-system .; then
        log_success "Build concluído: ./bin/lep-system"
    else
        log_error "Falha no build"
    fi
    press_enter
}

dev_health_check() {
    log_step "Health check do ambiente DEV..."

    if curl -s -f "http://localhost:8080/health" > /dev/null; then
        log_success "✅ Aplicação DEV está saudável!"
        response=$(curl -s "http://localhost:8080/health" 2>/dev/null || echo "{}")
        if [ "$response" != "{}" ]; then
            echo ""
            log_info "Resposta do health check:"
            echo "$response" | python3 -m json.tool 2>/dev/null || echo "$response"
        fi
    else
        log_error "❌ Aplicação DEV não está respondendo"
        log_info "Verifique se o ambiente está rodando (opção 1)"
    fi
    press_enter
}

dev_status() {
    log_step "Status do ambiente DEV..."
    cd "$ROOT_DIR"

    echo ""
    log_info "Docker containers:"
    docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep -E "(lep-|postgres|mailhog|redis)" || echo "Nenhum container LEP rodando"

    echo ""
    log_info "Conectividade:"
    if curl -s --connect-timeout 3 "http://localhost:8080/ping" &> /dev/null; then
        echo -e "  ✅ API: http://localhost:8080"
    else
        echo -e "  ❌ API: não está respondendo"
    fi

    if curl -s --connect-timeout 3 "http://localhost:8025" &> /dev/null; then
        echo -e "  ✅ MailHog: http://localhost:8025"
    else
        echo -e "  ❌ MailHog: não está rodando"
    fi

    press_enter
}

dev_stop_and_clean() {
    log_step "Parando e limpando ambiente DEV..."
    cd "$ROOT_DIR"

    if confirm_action "Parar todos os containers Docker do LEP?"; then
        log_info "Parando containers..."
        docker-compose down --remove-orphans 2>/dev/null || true

        log_info "Removendo volumes órfãos..."
        docker volume prune -f 2>/dev/null || true

        log_success "Ambiente DEV parado e limpo!"
    fi
    press_enter
}

dev_seed_data() {
    log_step "Populando dados demo no ambiente DEV..."
    cd "$ROOT_DIR"

    if [ -f "scripts/run_seed.sh" ]; then
        chmod +x scripts/run_seed.sh
        ENVIRONMENT=dev ./scripts/run_seed.sh --verbose
    else
        go run cmd/seed/main.go --environment=dev --verbose
    fi
    press_enter
}

# ==============================================================================
# 2. AMBIENTE STAGE (GCP)
# ==============================================================================

show_stage_menu() {
    show_main_banner
    log_header "stage" "Ambiente STAGE - GCP"
    echo ""
    echo -e "${CYAN}☁️ Cloud SQL + GCS + credenciais padronizadas${NC}"
    echo ""
    echo "  1. 🖥️  Executar local (conecta GCP)"
    echo "  2. 🚀 Deploy no Cloud Run"
    echo "  3. 🏗️  Bootstrap infraestrutura"
    echo "  4. 💊 Health check STAGE"
    echo "  5. 📊 Status serviços GCP"
    echo "  6. 🌱 Popular dados demo STAGE"
    echo "  0. ⬅️  Voltar ao menu principal"
    echo ""
}

handle_stage_menu() {
    while true; do
        show_stage_menu
        read -p "Selecione uma opção: " choice

        case $choice in
            1) stage_run_local ;;
            2) stage_deploy_cloud_run ;;
            3) stage_bootstrap_infrastructure ;;
            4) stage_health_check ;;
            5) stage_services_status ;;
            6) stage_seed_data ;;
            0) return ;;
            *) log_error "Opção inválida. Tente novamente." ; press_enter ;;
        esac
    done
}

stage_run_local() {
    log_step "Executando STAGE local (conecta GCP)..."
    cd "$ROOT_DIR"

    if ! check_dependencies "gcloud"; then
        press_enter
        return
    fi

    # Verificar autenticação
    if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q "@"; then
        log_error "Não autenticado no Google Cloud. Execute:"
        log_info "  gcloud auth login"
        log_info "  gcloud auth application-default login"
        press_enter
        return
    fi

    log_info "Executando script stage-local.sh..."
    if [ -f "scripts/stage-local.sh" ]; then
        chmod +x scripts/stage-local.sh
        ./scripts/stage-local.sh
    else
        log_error "Script stage-local.sh não encontrado"
    fi
    press_enter
}

stage_deploy_cloud_run() {
    log_step "Deploy STAGE no Cloud Run..."
    cd "$ROOT_DIR"

    if ! check_dependencies "gcloud"; then
        press_enter
        return
    fi

    if confirm_action "Deploy no Cloud Run para acesso online?"; then
        log_info "Executando script stage-deploy.sh..."
        if [ -f "scripts/stage-deploy.sh" ]; then
            chmod +x scripts/stage-deploy.sh
            ./scripts/stage-deploy.sh
        else
            log_error "Script stage-deploy.sh não encontrado"
        fi
    fi
    press_enter
}

stage_bootstrap_infrastructure() {
    log_step "Bootstrap da infraestrutura STAGE..."
    cd "$ROOT_DIR"

    if ! check_dependencies "terraform"; then
        press_enter
        return
    fi

    log_warn "O bootstrap criará:"
    log_info "  - Cloud SQL PostgreSQL"
    log_info "  - Google Cloud Storage bucket"
    log_info "  - Service Account e secrets"
    log_info "  - Habilitação de APIs necessárias"

    if confirm_action "Continuar com o bootstrap?"; then
        log_info "Inicializando Terraform..."
        # Os arquivos .tf estão na raiz do projeto
        if [ ! -f "main.tf" ]; then
            log_error "Arquivos Terraform não encontrados na raiz do projeto"
            press_enter
            return
        fi

        log_info "Executando terraform init..."
        terraform init

        log_info "Executando terraform apply..."
        terraform apply -var-file=environments/gcp-stage.tfvars

        log_success "Bootstrap concluído!"
    fi
    press_enter
}

stage_health_check() {
    log_step "Health check do ambiente STAGE..."

    # Verificar se há serviços rodando no Cloud Run
    if ! check_dependencies "gcloud"; then
        press_enter
        return
    fi

    local services=$(gcloud run services list --region="$REGION" --format="value(status.url)" 2>/dev/null)

    if [ -z "$services" ]; then
        log_warn "Nenhum serviço Cloud Run encontrado"
        log_info "Execute o deploy primeiro (opção 2)"
    else
        for service_url in $services; do
            if [ -n "$service_url" ]; then
                local service_name=$(echo "$service_url" | sed 's|https://||' | cut -d'.' -f1)
                log_info "Testando $service_name..."
                if curl -s --connect-timeout 5 "$service_url/health" &> /dev/null; then
                    log_success "✅ $service_name: saudável"
                    log_info "URL: $service_url"
                else
                    log_error "❌ $service_name: não responsivo"
                fi
            fi
        done
    fi
    press_enter
}

stage_services_status() {
    log_step "Status dos serviços STAGE no GCP..."

    if ! check_dependencies "gcloud"; then
        press_enter
        return
    fi

    echo ""
    log_info "Serviços Cloud Run:"
    gcloud run services list --region="$REGION" --format="table(metadata.name,status.url,status.traffic[0].percent)" 2>/dev/null || log_warn "Nenhum serviço encontrado"

    echo ""
    log_info "Cloud SQL instâncias:"
    gcloud sql instances list --format="table(name,region,databaseVersion,settings.tier,status)" 2>/dev/null || log_warn "Nenhuma instância SQL encontrada"

    echo ""
    log_info "Storage buckets:"
    gcloud storage ls 2>/dev/null | grep lep || log_warn "Nenhum bucket LEP encontrado"

    press_enter
}

stage_seed_data() {
    log_step "Populando dados demo no ambiente STAGE..."
    cd "$ROOT_DIR"

    log_warn "Isso populará o banco Cloud SQL com dados demo"
    if confirm_action "Continuar?"; then
        if [ -f "scripts/run_seed.sh" ]; then
            chmod +x scripts/run_seed.sh
            ENVIRONMENT=stage ./scripts/run_seed.sh --verbose
        else
            ENVIRONMENT=stage go run cmd/seed/main.go --environment=stage --verbose
        fi
    fi
    press_enter
}

# ==============================================================================
# 3. DATABASE & SEEDING
# ==============================================================================

show_database_menu() {
    show_main_banner
    log_header "database" "Database & Seeding"
    echo ""
    echo "  1. 🌱 Popular DEV (Docker local)"
    echo "  2. ☁️ Popular STAGE (Cloud SQL)"
    echo "  3. 🧹 Limpar e repopular DEV"
    echo "  4. 🧹 Limpar e repopular STAGE"
    echo "  5. 👥 Apenas usuários demo"
    echo "  6. 📊 Status das databases"
    echo "  0. ⬅️  Voltar ao menu principal"
    echo ""
}

handle_database_menu() {
    while true; do
        show_database_menu
        read -p "Selecione uma opção: " choice

        case $choice in
            1) database_seed_dev ;;
            2) database_seed_stage ;;
            3) database_clear_and_seed_dev ;;
            4) database_clear_and_seed_stage ;;
            5) database_seed_users_only ;;
            6) database_status ;;
            0) return ;;
            *) log_error "Opção inválida. Tente novamente." ; press_enter ;;
        esac
    done
}

database_seed_dev() {
    log_step "Populando database DEV (Docker local)..."
    cd "$ROOT_DIR"
    ENVIRONMENT=dev ./scripts/run_seed.sh --verbose 2>/dev/null || go run cmd/seed/main.go --environment=dev --verbose
    press_enter
}

database_seed_stage() {
    log_step "Populando database STAGE (Cloud SQL)..."
    cd "$ROOT_DIR"
    log_warn "Conectará no Cloud SQL para popular dados"
    if confirm_action "Continuar?"; then
        ENVIRONMENT=stage ./scripts/run_seed.sh --verbose 2>/dev/null || ENVIRONMENT=stage go run cmd/seed/main.go --environment=stage --verbose
    fi
    press_enter
}

database_clear_and_seed_dev() {
    log_step "Limpando e repopulando DEV..."
    if confirm_action "⚠️ Apagar TODOS os dados do DEV?"; then
        cd "$ROOT_DIR"
        ENVIRONMENT=dev ./scripts/run_seed.sh --clear-first --verbose 2>/dev/null || go run cmd/seed/main.go --environment=dev --clear-first --verbose
    fi
    press_enter
}

database_clear_and_seed_stage() {
    log_step "Limpando e repopulando STAGE..."
    if confirm_action "⚠️ Apagar TODOS os dados do Cloud SQL STAGE?"; then
        cd "$ROOT_DIR"
        ENVIRONMENT=stage ./scripts/run_seed.sh --clear-first --verbose 2>/dev/null || ENVIRONMENT=stage go run cmd/seed/main.go --environment=stage --clear-first --verbose
    fi
    press_enter
}

database_seed_users_only() {
    log_step "Criando apenas usuários demo..."
    echo ""
    echo "Credenciais que serão criadas:"
    echo "  - admin@lep-demo.com / password (Admin)"
    echo "  - garcom@lep-demo.com / password (Garçom)"
    echo "  - gerente@lep-demo.com / password (Gerente)"

    if confirm_action "Criar estes usuários?"; then
        cd "$ROOT_DIR"
        go run cmd/seed/main.go --users-only --verbose 2>/dev/null || go run cmd/seed/main.go --verbose
    fi
    press_enter
}

database_status() {
    log_step "Status das databases..."

    echo ""
    log_info "DEV (Docker local):"
    if docker ps | grep postgres &> /dev/null; then
        echo -e "  ✅ PostgreSQL container rodando"
    else
        echo -e "  ❌ PostgreSQL container não encontrado"
    fi

    echo ""
    log_info "STAGE (Cloud SQL):"
    if gcloud sql instances list --format="value(name)" 2>/dev/null | grep -q "leps-postgres-stage"; then
        echo -e "  ✅ Cloud SQL instância encontrada"
    else
        echo -e "  ❌ Cloud SQL instância não encontrada"
    fi

    press_enter
}

# ==============================================================================
# 4. TESTES
# ==============================================================================

show_tests_menu() {
    show_main_banner
    log_header "tests" "Testes"
    echo ""
    echo "  1. 🧪 Executar todos os testes"
    echo "  2. 📊 Testes com cobertura"
    echo "  3. 📄 Relatório HTML de cobertura"
    echo "  4. 🎯 Teste específico"
    echo "  5. ⚡ Testes rápidos (sem cache)"
    echo "  0. ⬅️  Voltar ao menu principal"
    echo ""
}

handle_tests_menu() {
    while true; do
        show_tests_menu
        read -p "Selecione uma opção: " choice

        case $choice in
            1) tests_run_all ;;
            2) tests_run_with_coverage ;;
            3) tests_html_coverage ;;
            4) tests_run_specific ;;
            5) tests_run_fast ;;
            0) return ;;
            *) log_error "Opção inválida. Tente novamente." ; press_enter ;;
        esac
    done
}

tests_run_all() {
    log_step "Executando todos os testes..."
    cd "$ROOT_DIR"

    if [ -f "scripts/run_tests.sh" ]; then
        chmod +x scripts/run_tests.sh
        ./scripts/run_tests.sh
    else
        go test ./... -v
    fi
    press_enter
}

tests_run_with_coverage() {
    log_step "Testes com cobertura..."
    cd "$ROOT_DIR"

    if [ -f "scripts/run_tests.sh" ]; then
        chmod +x scripts/run_tests.sh
        ./scripts/run_tests.sh --coverage
    else
        go test ./... -coverprofile=coverage.out -v
        if [ -f "coverage.out" ]; then
            go tool cover -func=coverage.out | tail -n 1
        fi
    fi
    press_enter
}

tests_html_coverage() {
    log_step "Relatório HTML de cobertura..."
    cd "$ROOT_DIR"

    if [ -f "scripts/run_tests.sh" ]; then
        chmod +x scripts/run_tests.sh
        ./scripts/run_tests.sh --html
    else
        go test ./... -coverprofile=coverage.out
        go tool cover -html=coverage.out -o coverage.html
        log_success "Relatório HTML: coverage.html"
    fi
    press_enter
}

tests_run_specific() {
    echo ""
    read -p "Digite o nome do teste: " test_pattern
    if [ -n "$test_pattern" ]; then
        cd "$ROOT_DIR"
        go test ./... -run "$test_pattern" -v
    fi
    press_enter
}

tests_run_fast() {
    log_step "Testes rápidos (sem cache)..."
    cd "$ROOT_DIR"
    go clean -testcache
    go test ./... -count=1
    press_enter
}

# ==============================================================================
# 5. SETUP & CONFIGURAÇÃO
# ==============================================================================

show_setup_menu() {
    show_main_banner
    log_header "setup" "Setup & Configuração"
    echo ""
    echo "  1. 🔧 Verificar dependências"
    echo "  2. 🔐 Gerar chaves JWT"
    echo "  3. ⚙️  Configurar Google Cloud"
    echo "  4. 📄 Criar arquivos .env"
    echo "  5. ✅ Validar configuração"
    echo "  6. 🏗️  Setup completo inicial"
    echo "  0. ⬅️  Voltar ao menu principal"
    echo ""
}

handle_setup_menu() {
    while true; do
        show_setup_menu
        read -p "Selecione uma opção: " choice

        case $choice in
            1) setup_check_dependencies ;;
            2) setup_generate_jwt ;;
            3) setup_gcloud ;;
            4) setup_create_env_files ;;
            5) setup_validate_config ;;
            6) setup_complete_initial ;;
            0) return ;;
            *) log_error "Opção inválida. Tente novamente." ; press_enter ;;
        esac
    done
}

setup_check_dependencies() {
    log_step "Verificando dependências..."

    local required_tools=("go" "git" "curl" "docker" "gcloud")
    local missing_tools=()

    echo ""
    for tool in "${required_tools[@]}"; do
        if command -v "$tool" &> /dev/null; then
            echo -e "  ✅ ${tool}: $(${tool} version 2>/dev/null | head -n1 || echo 'instalado')"
        else
            echo -e "  ❌ ${tool}: NÃO INSTALADO"
            missing_tools+=("$tool")
        fi
    done

    if [ ${#missing_tools[@]} -eq 0 ]; then
        log_success "Todas as dependências estão instaladas!"
    else
        log_error "Dependências em falta: ${missing_tools[*]}"
    fi
    press_enter
}

setup_generate_jwt() {
    log_step "Gerando chaves JWT..."
    cd "$ROOT_DIR"

    if [ -f "jwt_private_key.pem" ]; then
        if ! confirm_action "Chaves JWT já existem. Sobrescrever?"; then
            return
        fi
    fi

    if openssl genpkey -algorithm RSA -out jwt_private_key.pem -pkcs8 2>/dev/null; then
        openssl rsa -pubout -in jwt_private_key.pem -out jwt_public_key.pem 2>/dev/null
        log_success "Chaves JWT geradas!"
        log_info "Atualize seus arquivos .env com as novas chaves"
    else
        log_error "Falha ao gerar chaves JWT"
    fi
    press_enter
}

setup_gcloud() {
    log_step "Configurando Google Cloud..."

    if ! check_dependencies "gcloud"; then
        press_enter
        return
    fi

    if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q "@"; then
        if confirm_action "Fazer login no Google Cloud?"; then
            gcloud auth login
            gcloud auth application-default login
        fi
    fi

    gcloud config set project "$PROJECT_ID"
    log_success "Google Cloud configurado para projeto $PROJECT_ID"
    press_enter
}

setup_create_env_files() {
    log_step "Criando arquivos de configuração..."
    cd "$ROOT_DIR"

    # Criar .env para DEV se não existir
    if [ ! -f ".env" ]; then
        cat > .env << 'EOF'
# LEP System - DEV Environment (Local)
ENVIRONMENT=dev
PORT=8080

# Database (Docker)
DB_HOST=postgres
DB_PORT=5432
DB_USER=lep_user
DB_PASS=lep_password
DB_NAME=lep_database
DB_SSL_MODE=disable

# JWT
JWT_SECRET_PRIVATE_KEY=dev-simple-private-key-for-testing-only
JWT_SECRET_PUBLIC_KEY=dev-simple-public-key-for-testing-only

# Storage (Local)
STORAGE_TYPE=local
BUCKET_NAME=lep-dev-bucket
BASE_URL=http://localhost:8080

# SMTP (MailHog)
SMTP_HOST=mailhog
SMTP_PORT=1025

# Application
ENABLE_CRON_JOBS=false
GIN_MODE=debug
LOG_LEVEL=debug
EOF
        log_success "Arquivo .env criado para DEV"
    fi

    # Criar .env.stage para STAGE se não existir
    if [ ! -f ".env.stage" ]; then
        cat > .env.stage << 'EOF'
# LEP System - STAGE Environment (GCP)
ENVIRONMENT=stage
PORT=8080

# Database (Cloud SQL)
DB_USER=lep_user
DB_PASS=lep_password
DB_NAME=lep_database
INSTANCE_UNIX_SOCKET=/cloudsql/leps-472702:us-central1:leps-postgres-stage

# JWT
JWT_SECRET_PRIVATE_KEY=dev-simple-private-key-for-testing-only
JWT_SECRET_PUBLIC_KEY=dev-simple-public-key-for-testing-only

# Storage (GCS)
STORAGE_TYPE=gcs
BUCKET_NAME=leps-472702-lep-images-stage
BASE_URL=https://storage.googleapis.com/leps-472702-lep-images-stage

# Application
ENABLE_CRON_JOBS=true
GIN_MODE=release
LOG_LEVEL=info
EOF
        log_success "Arquivo .env.stage criado para STAGE"
    fi

    press_enter
}

setup_validate_config() {
    log_step "Validando configuração..."
    cd "$ROOT_DIR"

    local errors=0

    # Verificar arquivos
    local files=(".env" "go.mod" "main.go" "docker-compose.yml")
    for file in "${files[@]}"; do
        if [ -f "$file" ]; then
            echo -e "  ✅ $file"
        else
            echo -e "  ❌ $file"
            errors=$((errors + 1))
        fi
    done

    # Testar build
    if go build -o /tmp/test-build . &> /dev/null; then
        echo -e "  ✅ Build Go"
        rm -f /tmp/test-build
    else
        echo -e "  ❌ Build Go falhou"
        errors=$((errors + 1))
    fi

    if [ $errors -eq 0 ]; then
        log_success "Configuração válida!"
    else
        log_warn "$errors problemas encontrados"
    fi
    press_enter
}

setup_complete_initial() {
    log_step "Setup completo inicial..."

    if confirm_action "Executar setup completo (dependências + arquivos + validação)?"; then
        setup_check_dependencies
        setup_create_env_files
        setup_validate_config
        log_success "Setup inicial concluído!"
    fi
    press_enter
}

# ==============================================================================
# 6. UTILITÁRIOS
# ==============================================================================

show_utilities_menu() {
    show_main_banner
    log_header "utilities" "Utilitários"
    echo ""
    echo "  1. 🔍 Status geral do projeto"
    echo "  2. 🧹 Limpeza completa"
    echo "  3. 💾 Backup de configurações"
    echo "  4. 🩺 Troubleshooting automático"
    echo "  5. 📊 Relatório do ambiente"
    echo "  6. 🔄 Reset completo DEV"
    echo "  0. ⬅️  Voltar ao menu principal"
    echo ""
}

handle_utilities_menu() {
    while true; do
        show_utilities_menu
        read -p "Selecione uma opção: " choice

        case $choice in
            1) utilities_project_status ;;
            2) utilities_complete_cleanup ;;
            3) utilities_backup_configs ;;
            4) utilities_troubleshooting ;;
            5) utilities_environment_report ;;
            6) utilities_reset_dev ;;
            0) return ;;
            *) log_error "Opção inválida. Tente novamente." ; press_enter ;;
        esac
    done
}

utilities_project_status() {
    log_step "Status geral do projeto..."
    cd "$ROOT_DIR"

    echo ""
    echo -e "${CYAN}=== AMBIENTE DEV ===${NC}"
    if docker ps | grep -E "(lep-|postgres)" &> /dev/null; then
        echo -e "  ✅ Containers Docker rodando"
    else
        echo -e "  ❌ Containers Docker parados"
    fi

    echo ""
    echo -e "${CYAN}=== AMBIENTE STAGE ===${NC}"
    if gcloud run services list --region="$REGION" 2>/dev/null | grep -q "lep-system"; then
        echo -e "  ✅ Serviços Cloud Run ativos"
    else
        echo -e "  ❌ Nenhum serviço Cloud Run encontrado"
    fi

    echo ""
    echo -e "${CYAN}=== ARQUIVOS ===${NC}"
    local files=(".env" ".env.stage" "docker-compose.yml" "scripts/dev-local.sh" "scripts/stage-local.sh" "scripts/stage-deploy.sh")
    for file in "${files[@]}"; do
        if [ -f "$file" ]; then
            echo -e "  ✅ $file"
        else
            echo -e "  ❌ $file"
        fi
    done

    press_enter
}

utilities_complete_cleanup() {
    log_step "Limpeza completa..."

    if confirm_action "Limpar binários, cache Go, containers Docker?"; then
        cd "$ROOT_DIR"

        rm -rf bin/ logs/ coverage.* 2>/dev/null || true
        go clean -cache -modcache 2>/dev/null || true
        docker system prune -f 2>/dev/null || true

        log_success "Limpeza concluída!"
    fi
    press_enter
}

utilities_backup_configs() {
    log_step "Backup de configurações..."
    cd "$ROOT_DIR"

    local backup_dir="backup-$(date +%Y%m%d-%H%M%S)"
    mkdir -p "$backup_dir"

    local files=(".env" ".env.stage" "docker-compose.yml" "environments/gcp-stage.tfvars")
    for file in "${files[@]}"; do
        if [ -f "$file" ]; then
            cp "$file" "$backup_dir/"
            echo -e "  ✅ $file"
        fi
    done

    log_success "Backup criado em: $backup_dir"
    press_enter
}

utilities_troubleshooting() {
    log_step "Troubleshooting automático..."

    local issues=0

    # Verificar Go
    if ! command -v go &> /dev/null; then
        echo -e "  ❌ Go não instalado"
        issues=$((issues + 1))
    else
        echo -e "  ✅ Go instalado"
    fi

    # Verificar Docker
    if ! docker info &> /dev/null; then
        echo -e "  ❌ Docker não está rodando"
        issues=$((issues + 1))
    else
        echo -e "  ✅ Docker funcionando"
    fi

    # Verificar build
    cd "$ROOT_DIR"
    if ! go build -o /tmp/test . &> /dev/null; then
        echo -e "  ❌ Build falha"
        issues=$((issues + 1))
    else
        echo -e "  ✅ Build OK"
        rm -f /tmp/test
    fi

    if [ $issues -eq 0 ]; then
        log_success "Nenhum problema encontrado!"
    else
        log_warn "$issues problemas encontrados"
    fi
    press_enter
}

utilities_environment_report() {
    log_step "Gerando relatório do ambiente..."
    cd "$ROOT_DIR"

    local report="environment-report-$(date +%Y%m%d-%H%M%S).txt"
    {
        echo "LEP System Environment Report"
        echo "Generated: $(date)"
        echo "=========================="
        echo ""
        echo "Go: $(go version 2>/dev/null || echo 'not installed')"
        echo "Docker: $(docker --version 2>/dev/null || echo 'not installed')"
        echo "gcloud: $(gcloud version --format='value(VERSION)' 2>/dev/null || echo 'not installed')"
        echo ""
        echo "Project files:"
        ls -la
    } > "$report"

    log_success "Relatório gerado: $report"
    press_enter
}

utilities_reset_dev() {
    log_step "Reset completo do ambiente DEV..."

    log_warn "Isso irá:"
    log_info "  - Parar todos os containers"
    log_info "  - Limpar volumes Docker"
    log_info "  - Recriar .env padrão"
    log_info "  - Limpar cache Go"

    if confirm_action "Confirma o reset do DEV?"; then
        cd "$ROOT_DIR"

        docker-compose down --volumes 2>/dev/null || true
        docker system prune -f 2>/dev/null || true
        go clean -cache 2>/dev/null || true

        # Recriar .env
        mv .env .env.backup-$(date +%H%M%S) 2>/dev/null || true
        setup_create_env_files

        log_success "Reset DEV concluído!"
    fi
    press_enter
}

# ==============================================================================
# 7. AJUDA
# ==============================================================================

show_help_menu() {
    show_main_banner
    log_header "help" "Ajuda"
    echo ""
    echo "  1. 📖 Guia de primeiros passos"
    echo "  2. 🔧 Novos ambientes (dev/stage)"
    echo "  3. 🚨 Solução de problemas"
    echo "  4. 📝 Comandos úteis"
    echo "  5. 🌐 Links úteis"
    echo "  0. ⬅️  Voltar ao menu principal"
    echo ""
}

handle_help_menu() {
    while true; do
        show_help_menu
        read -p "Selecione uma opção: " choice

        case $choice in
            1) help_getting_started ;;
            2) help_new_environments ;;
            3) help_troubleshooting ;;
            4) help_useful_commands ;;
            5) help_useful_links ;;
            0) return ;;
            *) log_error "Opção inválida. Tente novamente." ; press_enter ;;
        esac
    done
}

help_getting_started() {
    show_main_banner
    log_header "help" "Guia de Primeiros Passos"
    echo ""

    echo -e "${CYAN}🚀 Para começar (primeira vez):${NC}"
    echo "  1. Execute setup: Menu 5 -> 6"
    echo "  2. Inicie DEV: Menu 1 -> 1"
    echo "  3. Popule dados: Menu 1 -> 6"
    echo "  4. Teste: curl http://localhost:8080/health"
    echo ""

    echo -e "${CYAN}🔑 Credenciais padrão:${NC}"
    echo "  - admin@lep-demo.com / password"
    echo "  - garcom@lep-demo.com / password"
    echo "  - gerente@lep-demo.com / password"
    echo ""

    press_enter
}

help_new_environments() {
    show_main_banner
    log_header "help" "Novos Ambientes"
    echo ""

    echo -e "${CYAN}🔧 DEV - Desenvolvimento Local:${NC}"
    echo "  - 100% local com Docker"
    echo "  - PostgreSQL + Redis + MailHog"
    echo "  - localStorage para uploads"
    echo "  - Credenciais padronizadas"
    echo ""

    echo -e "${CYAN}🚀 STAGE - GCP:${NC}"
    echo "  - Cloud SQL + Google Cloud Storage"
    echo "  - Execução local OU Cloud Run"
    echo "  - Mesmas credenciais do DEV"
    echo "  - Para testes de integração"
    echo ""

    echo -e "${CYAN}🏭 PROD - Produção (futuro):${NC}"
    echo "  - Configurações profissionais"
    echo "  - Credenciais únicas"
    echo "  - Alta disponibilidade"
    echo ""

    press_enter
}

help_troubleshooting() {
    show_main_banner
    log_header "help" "Solução de Problemas"
    echo ""

    echo -e "${CYAN}🔧 Problema: Containers não sobem${NC}"
    echo "  1. docker system prune -f"
    echo "  2. Reiniciar Docker Desktop"
    echo "  3. Menu 1 -> 1 (dev-local.sh)"
    echo ""

    echo -e "${CYAN}🔧 Problema: Build falha${NC}"
    echo "  1. go mod tidy"
    echo "  2. go clean -cache"
    echo "  3. Menu 6 -> 4 (troubleshooting)"
    echo ""

    echo -e "${CYAN}🔧 Problema: GCP não conecta${NC}"
    echo "  1. gcloud auth login"
    echo "  2. gcloud config set project leps-472702"
    echo "  3. Menu 5 -> 3 (configurar gcloud)"
    echo ""

    press_enter
}

help_useful_commands() {
    show_main_banner
    log_header "help" "Comandos Úteis"
    echo ""

    echo -e "${CYAN}📦 Go:${NC}"
    echo "  go run main.go       # Executar"
    echo "  go build .           # Build"
    echo "  go test ./...        # Testes"
    echo ""

    echo -e "${CYAN}🐳 Docker:${NC}"
    echo "  docker-compose up    # Subir containers"
    echo "  docker-compose down  # Parar containers"
    echo "  docker ps            # Ver containers"
    echo ""

    echo -e "${CYAN}☁️ GCP:${NC}"
    echo "  gcloud auth login    # Login"
    echo "  gcloud run services list  # Listar serviços"
    echo ""

    press_enter
}

help_useful_links() {
    show_main_banner
    log_header "help" "Links Úteis"
    echo ""

    echo -e "${CYAN}📚 Documentação:${NC}"
    echo "  Go: https://golang.org/doc/"
    echo "  Docker: https://docs.docker.com/"
    echo "  Google Cloud: https://cloud.google.com/docs"
    echo ""

    echo -e "${CYAN}⚙️ Downloads:${NC}"
    echo "  Go: https://golang.org/dl/"
    echo "  Docker: https://docker.com/get-started"
    echo "  gcloud: https://cloud.google.com/sdk/docs/install"
    echo ""

    press_enter
}

# ==============================================================================
# MAIN EXECUTION
# ==============================================================================

# Handle command line arguments
handle_batch_mode() {
    case "$1" in
        "--help"|"-h")
            show_main_banner
            echo "LEP System Master Script (Refatorado)"
            echo ""
            echo "Uso: $0 [OPÇÃO]"
            echo ""
            echo "Opções:"
            echo "  --help, -h       Mostrar ajuda"
            echo "  --dev            Iniciar ambiente DEV"
            echo "  --stage          Menu ambiente STAGE"
            echo "  --seed-dev       Popular dados DEV"
            echo "  --test           Executar testes"
            echo "  --status         Status do projeto"
            echo ""
            exit 0
            ;;
        "--dev")
            dev_start_environment
            exit 0
            ;;
        "--stage")
            handle_stage_menu
            exit 0
            ;;
        "--seed-dev")
            database_seed_dev
            exit 0
            ;;
        "--test")
            tests_run_all
            exit 0
            ;;
        "--status")
            utilities_project_status
            exit 0
            ;;
        "")
            return 0
            ;;
        *)
            log_error "Opção desconhecida: $1"
            echo "Use --help para ver opções disponíveis"
            exit 1
            ;;
    esac
}

# Main interactive loop
main_loop() {
    while true; do
        show_main_banner
        show_main_menu

        read -p "Selecione uma opção: " choice

        case $choice in
            1) handle_dev_menu ;;
            2) handle_stage_menu ;;
            3) handle_database_menu ;;
            4) handle_tests_menu ;;
            5) handle_setup_menu ;;
            6) handle_utilities_menu ;;
            7) handle_help_menu ;;
            0)
                echo ""
                log_success "👋 Obrigado por usar o LEP System!"
                echo ""
                exit 0
                ;;
            *)
                log_error "Opção inválida. Tente novamente."
                press_enter
                ;;
        esac
    done
}

# Script initialization
init_script() {
    cd "$ROOT_DIR" 2>/dev/null || {
        log_error "Diretório do projeto não encontrado: $ROOT_DIR"
        exit 1
    }

    mkdir -p bin logs 2>/dev/null || true
    chmod +x scripts/*.sh 2>/dev/null || true
    trap 'echo ""; log_warn "Script interrompido."; exit 130' INT TERM
}

# Main execution
main() {
    init_script
    handle_batch_mode "$1"
    main_loop
}

# Execute
main "$@"