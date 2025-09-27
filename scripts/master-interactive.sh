#!/bin/bash

# LEP System - Master Interactive Script
# Unified interface for all LEP System operations and deployments
# Consolidates all existing scripts into a single, user-friendly menu system

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
    echo "                    LEP System - Master Control                  "
    echo "=================================================================="
    echo -e "${NC}"
    echo -e "${WHITE}Project:${NC} ${PROJECT_NAME} (${PROJECT_ID})"
    echo -e "${WHITE}Region:${NC} ${REGION}"
    echo -e "${WHITE}Version:${NC} 1.0.0"
    echo ""
}

# Main menu display
show_main_menu() {
    echo -e "${CYAN}🎛️  Selecione uma categoria:${NC}"
    echo ""
    echo "  1. 🏠 Desenvolvimento Local"
    echo "  2. ⚙️  Setup & Configuração"
    echo "  3. 🌱 Database & Seeding"
    echo "  4. 🧪 Testes"
    echo "  5. ☁️  Deploy GCP"
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
# 1. DESENVOLVIMENTO LOCAL
# ==============================================================================

show_dev_menu() {
    show_main_banner
    log_header "dev" "Desenvolvimento Local"
    echo ""
    echo "  1. 🚀 Iniciar servidor (go run main.go)"
    echo "  2. 🔨 Build da aplicação"
    echo "  3. 🐳 Docker local"
    echo "  4. 💊 Health check"
    echo "  5. 📊 Status do servidor"
    echo "  6. 🧹 Limpeza de artifacts"
    echo "  0. ⬅️  Voltar ao menu principal"
    echo ""
}

handle_dev_menu() {
    while true; do
        show_dev_menu
        read -p "Selecione uma opção: " choice

        case $choice in
            1) dev_start_server ;;
            2) dev_build_app ;;
            3) dev_docker_local ;;
            4) dev_health_check ;;
            5) dev_server_status ;;
            6) dev_clean_artifacts ;;
            0) return ;;
            *) log_error "Opção inválida. Tente novamente." ; press_enter ;;
        esac
    done
}

dev_start_server() {
    log_step "Iniciando servidor de desenvolvimento..."
    cd "$ROOT_DIR"

    if [ ! -f ".env" ]; then
        log_warn ".env não encontrado. Criando arquivo exemplo..."
        cat > .env << 'EOF'
# Database configuration
DB_USER=postgres
DB_PASS=password
DB_NAME=lep_database

# JWT configuration
JWT_SECRET_PRIVATE_KEY=your_private_key_here
JWT_SECRET_PUBLIC_KEY=your_public_key_here

# Application configuration
PORT=8080
ENABLE_CRON_JOBS=true
EOF
        log_warn "Por favor, configure o arquivo .env com suas credenciais reais."
        press_enter
        return
    fi

    log_info "Servidor será iniciado em http://localhost:8080"
    log_info "Pressione Ctrl+C para parar o servidor"
    echo ""
    go run main.go
}

dev_build_app() {
    log_step "Construindo aplicação..."
    cd "$ROOT_DIR"

    if ! check_dependencies "go"; then
        press_enter
        return
    fi

    mkdir -p bin
    if go build -o bin/lep-system .; then
        log_success "Build concluído com sucesso!"
        log_info "Binário criado em: ./bin/lep-system"
    else
        log_error "Falha no build da aplicação"
    fi
    press_enter
}

dev_docker_local() {
    log_step "Construindo e executando container Docker local..."
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

    log_info "Construindo imagem Docker..."
    if docker build -t lep-system:local .; then
        log_success "Imagem construída com sucesso!"

        if [ -f ".env" ]; then
            log_info "Executando container na porta 8080..."
            docker run -p 8080:8080 --env-file .env lep-system:local
        else
            log_warn "Arquivo .env não encontrado. Executando sem variáveis de ambiente..."
            docker run -p 8080:8080 lep-system:local
        fi
    else
        log_error "Falha ao construir a imagem Docker"
        press_enter
    fi
}

dev_health_check() {
    log_step "Verificando saúde da aplicação..."

    if curl -s -f "http://localhost:8080/health" > /dev/null; then
        log_success "✅ Aplicação está saudável!"

        # Mostrar informações detalhadas se disponíveis
        response=$(curl -s "http://localhost:8080/health" 2>/dev/null || echo "{}")
        if [ "$response" != "{}" ]; then
            echo ""
            log_info "Resposta do health check:"
            echo "$response" | python3 -m json.tool 2>/dev/null || echo "$response"
        fi
    else
        log_error "❌ Aplicação não está respondendo ou não está rodando"
        log_info "Verifique se o servidor está rodando em http://localhost:8080"
    fi
    press_enter
}

dev_server_status() {
    log_step "Verificando status do servidor..."

    # Verificar se a porta 8080 está em uso
    if lsof -i :8080 &> /dev/null || netstat -an 2>/dev/null | grep :8080 &> /dev/null; then
        log_success "✅ Servidor está rodando na porta 8080"

        # Tentar obter informações do processo
        local pid=$(lsof -t -i :8080 2>/dev/null || echo "")
        if [ -n "$pid" ]; then
            log_info "PID do processo: $pid"
            log_info "Comando: $(ps -p $pid -o command= 2>/dev/null || echo 'N/A')"
        fi

        # Verificar conectividade
        if curl -s --connect-timeout 3 "http://localhost:8080/ping" &> /dev/null; then
            log_success "✅ Endpoint /ping respondendo"
        else
            log_warn "⚠️  Endpoint /ping não está respondendo"
        fi

    else
        log_warn "❌ Nenhum servidor detectado na porta 8080"
    fi
    press_enter
}

dev_clean_artifacts() {
    log_step "Limpando artifacts de build..."
    cd "$ROOT_DIR"

    # Remover binários
    if [ -d "bin" ]; then
        rm -rf bin/
        log_info "Diretório bin/ removido"
    fi

    # Remover arquivos temporários
    rm -f tfplan 2>/dev/null || true
    rm -f .terraform.lock.hcl 2>/dev/null || true
    rm -f terraform.tfstate.lock.info 2>/dev/null || true

    # Limpar cache do Go
    go clean -cache 2>/dev/null || true
    go clean -modcache 2>/dev/null || log_info "Cache do Go mantido (requer permissões)"

    log_success "Limpeza concluída!"
    press_enter
}

# ==============================================================================
# 2. SETUP & CONFIGURAÇÃO
# ==============================================================================

show_setup_menu() {
    show_main_banner
    log_header "setup" "Setup & Configuração"
    echo ""
    echo "  1. 🏗️  Setup completo do ambiente"
    echo "  2. 📦 Verificar e instalar dependências"
    echo "  3. 🔐 Gerar chaves JWT"
    echo "  4. ⚙️  Configurar Google Cloud"
    echo "  5. 🐳 Configurar Docker"
    echo "  6. ✅ Validar configuração completa"
    echo "  7. 📄 Criar arquivos de configuração"
    echo "  0. ⬅️  Voltar ao menu principal"
    echo ""
}

handle_setup_menu() {
    while true; do
        show_setup_menu
        read -p "Selecione uma opção: " choice

        case $choice in
            1) setup_complete_environment ;;
            2) setup_check_dependencies ;;
            3) setup_generate_jwt_keys ;;
            4) setup_gcloud_config ;;
            5) setup_docker_config ;;
            6) setup_validate_config ;;
            7) setup_create_config_files ;;
            0) return ;;
            *) log_error "Opção inválida. Tente novamente." ; press_enter ;;
        esac
    done
}

setup_complete_environment() {
    log_step "Executando setup completo do ambiente..."

    if confirm_action "Isso executará o setup completo incluindo verificação de dependências, configuração do Google Cloud, Docker e criação de arquivos de configuração."; then
        cd "$ROOT_DIR"

        if [ -f "scripts/setup.sh" ]; then
            chmod +x scripts/setup.sh
            ./scripts/setup.sh
        else
            log_error "Script setup.sh não encontrado"
        fi
    fi
    press_enter
}

setup_check_dependencies() {
    log_step "Verificando dependências do sistema..."

    local missing_tools=()
    local required_tools=("go" "git" "curl" "docker" "gcloud" "terraform")

    echo ""
    echo -e "${CYAN}Dependências obrigatórias:${NC}"
    for tool in "${required_tools[@]}"; do
        if command -v "$tool" &> /dev/null; then
            local version=$(${tool} version 2>/dev/null | head -n1 || echo "versão não detectada")
            echo -e "  ✅ ${tool}: ${version}"
        else
            echo -e "  ❌ ${tool}: não instalado"
            missing_tools+=("$tool")
        fi
    done

    echo ""
    if [ ${#missing_tools[@]} -eq 0 ]; then
        log_success "Todas as dependências estão instaladas!"

        # Verificar dependências opcionais
        echo ""
        echo -e "${CYAN}Dependências opcionais:${NC}"
        local optional_tools=("openssl" "jq" "python3")
        for tool in "${optional_tools[@]}"; do
            if command -v "$tool" &> /dev/null; then
                echo -e "  ✅ ${tool}: disponível"
            else
                echo -e "  ⚠️  ${tool}: não instalado (opcional)"
            fi
        done
    else
        log_error "Dependências em falta: ${missing_tools[*]}"
        echo ""
        echo "Links para instalação:"
        for tool in "${missing_tools[@]}"; do
            case $tool in
                "gcloud") echo "  - Google Cloud CLI: https://cloud.google.com/sdk/docs/install" ;;
                "terraform") echo "  - Terraform: https://learn.hashicorp.com/tutorials/terraform/install-cli" ;;
                "docker") echo "  - Docker: https://docs.docker.com/get-docker/" ;;
                "go") echo "  - Go: https://golang.org/doc/install" ;;
                "git") echo "  - Git: https://git-scm.com/downloads" ;;
            esac
        done
    fi
    press_enter
}

setup_generate_jwt_keys() {
    log_step "Gerando chaves JWT..."

    if ! check_dependencies "openssl"; then
        press_enter
        return
    fi

    cd "$ROOT_DIR"

    # Verificar se as chaves já existem
    if [ -f "jwt_private_key.pem" ] && [ -f "jwt_public_key.pem" ]; then
        if ! confirm_action "Chaves JWT já existem. Sobrescrever?"; then
            return
        fi
    fi

    log_info "Gerando chave privada RSA..."
    if openssl genpkey -algorithm RSA -out jwt_private_key.pem -pkcs8 -aes256; then
        log_info "Gerando chave pública..."
        if openssl rsa -pubout -in jwt_private_key.pem -out jwt_public_key.pem; then
            log_success "Chaves JWT geradas com sucesso!"
            log_info "Arquivos criados:"
            log_info "  - jwt_private_key.pem (chave privada)"
            log_info "  - jwt_public_key.pem (chave pública)"
            log_warn "Atualize seus arquivos .env e terraform.tfvars com as novas chaves."
        else
            log_error "Falha ao gerar chave pública"
        fi
    else
        log_error "Falha ao gerar chave privada"
    fi
    press_enter
}

setup_gcloud_config() {
    log_step "Configurando Google Cloud..."

    if ! check_dependencies "gcloud"; then
        press_enter
        return
    fi

    # Verificar autenticação
    if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q "@"; then
        log_warn "Não autenticado no Google Cloud."
        if confirm_action "Fazer login no Google Cloud?"; then
            gcloud auth login
            gcloud auth application-default login
        else
            press_enter
            return
        fi
    fi

    # Configurar projeto
    log_info "Configurando projeto..."
    gcloud config set project "$PROJECT_ID"

    # Verificar acesso ao projeto
    if gcloud projects describe "$PROJECT_ID" &> /dev/null; then
        log_success "Projeto $PROJECT_ID configurado com sucesso!"

        # Configurar Docker
        log_info "Configurando Docker para Artifact Registry..."
        gcloud auth configure-docker "${REGION}-docker.pkg.dev"

        log_success "Configuração do Google Cloud concluída!"
    else
        log_error "Não foi possível acessar o projeto $PROJECT_ID"
        log_info "Verifique suas permissões no projeto"
    fi
    press_enter
}

setup_docker_config() {
    log_step "Configurando Docker..."

    if ! check_dependencies "docker"; then
        press_enter
        return
    fi

    if ! docker info &> /dev/null; then
        log_error "Docker não está rodando. Inicie o Docker primeiro."
        press_enter
        return
    fi

    log_success "Docker está rodando!"

    # Mostrar informações do Docker
    echo ""
    log_info "Informações do Docker:"
    docker version --format "  Version: {{.Server.Version}}"
    docker system df --format "  Disk Usage: {{.Size}}" 2>/dev/null || true

    # Configurar para Google Cloud se gcloud estiver disponível
    if command -v gcloud &> /dev/null; then
        log_info "Configurando autenticação do Docker para Google Cloud..."
        gcloud auth configure-docker "${REGION}-docker.pkg.dev" --quiet
        log_success "Docker configurado para Google Cloud!"
    fi

    press_enter
}

setup_validate_config() {
    log_step "Validando configuração completa..."

    cd "$ROOT_DIR"
    local errors=0

    echo ""
    echo -e "${CYAN}Verificando arquivos de configuração:${NC}"

    # Verificar arquivos essenciais
    local files=(".env" "go.mod" "main.go" "Dockerfile")
    for file in "${files[@]}"; do
        if [ -f "$file" ]; then
            echo -e "  ✅ $file"
        else
            echo -e "  ❌ $file (não encontrado)"
            errors=$((errors + 1))
        fi
    done

    # Verificar chaves JWT
    if [ -f "jwt_private_key.pem" ] && [ -f "jwt_public_key.pem" ]; then
        echo -e "  ✅ Chaves JWT"
    else
        echo -e "  ❌ Chaves JWT (não encontradas)"
        errors=$((errors + 1))
    fi

    echo ""
    echo -e "${CYAN}Verificando dependências:${NC}"

    local required_tools=("go" "git" "docker" "gcloud")
    for tool in "${required_tools[@]}"; do
        if command -v "$tool" &> /dev/null; then
            echo -e "  ✅ $tool"
        else
            echo -e "  ❌ $tool"
            errors=$((errors + 1))
        fi
    done

    echo ""
    echo -e "${CYAN}Verificando conectividade:${NC}"

    # Testar build do Go
    if go build -o /tmp/test-build . &> /dev/null; then
        echo -e "  ✅ Build do Go"
        rm -f /tmp/test-build
    else
        echo -e "  ❌ Build do Go (falhou)"
        errors=$((errors + 1))
    fi

    # Testar Docker
    if docker info &> /dev/null; then
        echo -e "  ✅ Docker"
    else
        echo -e "  ❌ Docker (não está rodando)"
        errors=$((errors + 1))
    fi

    # Testar Google Cloud
    if gcloud projects describe "$PROJECT_ID" &> /dev/null; then
        echo -e "  ✅ Google Cloud (projeto $PROJECT_ID)"
    else
        echo -e "  ⚠️  Google Cloud (projeto $PROJECT_ID não acessível)"
    fi

    echo ""
    if [ $errors -eq 0 ]; then
        log_success "🎉 Configuração completamente válida!"
        log_info "Seu ambiente está pronto para desenvolvimento e deploy."
    else
        log_warn "⚠️  Encontrados $errors problemas de configuração."
        log_info "Resolva os problemas acima antes de prosseguir."
    fi

    press_enter
}

setup_create_config_files() {
    log_step "Criando arquivos de configuração..."
    cd "$ROOT_DIR"

    # Criar .env se não existir
    if [ ! -f ".env" ]; then
        log_info "Criando arquivo .env..."
        cat > .env << 'EOF'
# LEP System - Environment Configuration

# Database Configuration
DB_USER=postgres
DB_PASS=your_database_password
DB_NAME=lep_database
DB_HOST=localhost
DB_PORT=5432

# JWT Configuration
JWT_SECRET_PRIVATE_KEY=your_jwt_private_key_content_here
JWT_SECRET_PUBLIC_KEY=your_jwt_public_key_content_here

# Twilio Configuration (SMS/WhatsApp)
TWILIO_ACCOUNT_SID=your_twilio_account_sid
TWILIO_AUTH_TOKEN=your_twilio_auth_token
TWILIO_PHONE_NUMBER=+1234567890

# SMTP Configuration (Email)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=your_app_password

# Storage Configuration
STORAGE_TYPE=local
STORAGE_BUCKET_NAME=your_gcs_bucket_name
BASE_URL=http://localhost:8080

# Application Configuration
PORT=8080
ENABLE_CRON_JOBS=true
EOF
        log_success "Arquivo .env criado"
    else
        log_info "Arquivo .env já existe"
    fi

    # Criar .env.example se não existir
    if [ ! -f ".env.example" ]; then
        cp .env .env.example 2>/dev/null || true
        log_info "Arquivo .env.example criado"
    fi

    # Criar .gitignore se não existir ou atualizar
    if [ ! -f ".gitignore" ]; then
        log_info "Criando arquivo .gitignore..."
        cat > .gitignore << 'EOF'
# Binaries
bin/
*.exe

# Environment files
.env
.env.local

# JWT Keys
jwt_private_key.pem
jwt_public_key.pem

# Terraform
terraform.tfvars
*.tfstate
*.tfstate.*
.terraform/
tfplan

# IDE
.vscode/
.idea/

# OS
.DS_Store
Thumbs.db

# Logs
logs/
*.log

# Temporary files
*.tmp
*.temp
EOF
        log_success "Arquivo .gitignore criado"
    fi

    log_info "Arquivos criados com sucesso!"
    log_warn "Não se esqueça de configurar o arquivo .env com suas credenciais reais."
    press_enter
}

# ==============================================================================
# 3. DATABASE & SEEDING
# ==============================================================================

show_database_menu() {
    show_main_banner
    log_header "database" "Database & Seeding"
    echo ""
    echo "  1. 🌱 Popular database com dados demo"
    echo "  2. 🧹 Limpar e repopular database"
    echo "  3. 🔧 Popular ambiente específico"
    echo "  4. 📊 Status da database"
    echo "  5. 🗂️  Popular apenas estruturas básicas"
    echo "  6. 👥 Popular apenas usuários demo"
    echo "  0. ⬅️  Voltar ao menu principal"
    echo ""
}

handle_database_menu() {
    while true; do
        show_database_menu
        read -p "Selecione uma opção: " choice

        case $choice in
            1) database_seed_demo_data ;;
            2) database_clear_and_seed ;;
            3) database_seed_environment ;;
            4) database_status ;;
            5) database_seed_basic ;;
            6) database_seed_users ;;
            0) return ;;
            *) log_error "Opção inválida. Tente novamente." ; press_enter ;;
        esac
    done
}

database_seed_demo_data() {
    log_step "Populando database com dados demo..."
    cd "$ROOT_DIR"

    if [ -f "scripts/run_seed.sh" ]; then
        chmod +x scripts/run_seed.sh
        ./scripts/run_seed.sh --verbose
    else
        # Fallback para execução direta
        log_info "Script de seeding não encontrado. Executando diretamente..."
        if [ -d "cmd/seed" ]; then
            go run cmd/seed/main.go --verbose
        else
            log_error "Sistema de seeding não encontrado"
        fi
    fi
    press_enter
}

database_clear_and_seed() {
    log_step "Limpando e repopulando database..."

    if confirm_action "⚠️  ATENÇÃO: Isso apagará TODOS os dados existentes na database!"; then
        cd "$ROOT_DIR"

        if [ -f "scripts/run_seed.sh" ]; then
            chmod +x scripts/run_seed.sh
            ./scripts/run_seed.sh --clear-first --verbose
        else
            if [ -d "cmd/seed" ]; then
                go run cmd/seed/main.go --clear-first --verbose
            else
                log_error "Sistema de seeding não encontrado"
            fi
        fi
    fi
    press_enter
}

database_seed_environment() {
    log_step "Populando ambiente específico..."
    echo ""
    echo "Ambientes disponíveis:"
    echo "  1. dev (desenvolvimento)"
    echo "  2. test (testes)"
    echo "  3. staging (homologação)"
    echo ""

    read -p "Selecione o ambiente (1-3): " env_choice

    local environment=""
    case $env_choice in
        1) environment="dev" ;;
        2) environment="test" ;;
        3) environment="staging" ;;
        *) log_error "Ambiente inválido" ; press_enter ; return ;;
    esac

    log_info "Populando ambiente: $environment"
    cd "$ROOT_DIR"

    if [ -f "scripts/run_seed.sh" ]; then
        chmod +x scripts/run_seed.sh
        ./scripts/run_seed.sh --environment="$environment" --verbose
    else
        if [ -d "cmd/seed" ]; then
            go run cmd/seed/main.go --environment="$environment" --verbose
        else
            log_error "Sistema de seeding não encontrado"
        fi
    fi
    press_enter
}

database_status() {
    log_step "Verificando status da database..."
    cd "$ROOT_DIR"

    # Verificar se conseguimos conectar na database
    echo ""
    log_info "Testando conectividade com a database..."

    # Tentar usar go run para testar a conexão
    if go run -c 'package main; import "fmt"; func main() { fmt.Println("Database connection test") }' &> /dev/null; then
        log_success "Go está funcionando"

        # Executar um teste de conexão básico se possível
        if curl -s --connect-timeout 5 "http://localhost:8080/health" &> /dev/null; then
            log_success "API está respondendo - database provavelmente OK"
        else
            log_warn "API não está respondendo - verifique se o servidor está rodando"
        fi
    else
        log_error "Problemas com Go ou dependências"
    fi

    # Mostrar configuração de database do .env
    if [ -f ".env" ]; then
        log_info "Configuração da database (.env):"
        grep -E "^DB_" .env 2>/dev/null || log_warn "Configurações de DB não encontradas no .env"
    else
        log_warn "Arquivo .env não encontrado"
    fi

    press_enter
}

database_seed_basic() {
    log_step "Populando apenas estruturas básicas..."
    cd "$ROOT_DIR"

    log_info "Criando organizações, projetos e configurações básicas..."

    if [ -d "cmd/seed" ]; then
        # Se houver parâmetros específicos para seed básico, usar aqui
        go run cmd/seed/main.go --basic-only --verbose 2>/dev/null || \
        go run cmd/seed/main.go --verbose
    else
        log_error "Sistema de seeding não encontrado"
    fi
    press_enter
}

database_seed_users() {
    log_step "Populando apenas usuários demo..."
    cd "$ROOT_DIR"

    log_info "Criando usuários de demonstração..."
    log_info "Credenciais que serão criadas:"
    echo "  - admin@lep-demo.com / password (Admin)"
    echo "  - garcom@lep-demo.com / password (Garçom)"
    echo "  - gerente@lep-demo.com / password (Gerente)"

    if confirm_action "Criar estes usuários?"; then
        if [ -d "cmd/seed" ]; then
            go run cmd/seed/main.go --users-only --verbose 2>/dev/null || \
            go run cmd/seed/main.go --verbose
        else
            log_error "Sistema de seeding não encontrado"
        fi
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
    echo "  2. 📊 Testes com cobertura de código"
    echo "  3. 📄 Relatório HTML de cobertura"
    echo "  4. 🎯 Executar teste específico"
    echo "  5. 🔍 Testes verbosos"
    echo "  6. ⚡ Testes rápidos (sem cache)"
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
            3) tests_html_coverage_report ;;
            4) tests_run_specific ;;
            5) tests_run_verbose ;;
            6) tests_run_fast ;;
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
        # Fallback para execução direta
        if ! check_dependencies "go"; then
            press_enter
            return
        fi

        log_info "Executando: go test ./..."
        if go test ./... -v; then
            log_success "Todos os testes passaram!"
        else
            log_error "Alguns testes falharam"
        fi
    fi
    press_enter
}

tests_run_with_coverage() {
    log_step "Executando testes com cobertura de código..."
    cd "$ROOT_DIR"

    if [ -f "scripts/run_tests.sh" ]; then
        chmod +x scripts/run_tests.sh
        ./scripts/run_tests.sh --coverage
    else
        if ! check_dependencies "go"; then
            press_enter
            return
        fi

        log_info "Executando testes com cobertura..."
        if go test ./... -coverprofile=coverage.out -v; then
            log_success "Testes concluídos!"

            if [ -f "coverage.out" ]; then
                log_info "Cobertura de código:"
                go tool cover -func=coverage.out | tail -n 1
            fi
        else
            log_error "Alguns testes falharam"
        fi
    fi
    press_enter
}

tests_html_coverage_report() {
    log_step "Gerando relatório HTML de cobertura..."
    cd "$ROOT_DIR"

    if [ -f "scripts/run_tests.sh" ]; then
        chmod +x scripts/run_tests.sh
        ./scripts/run_tests.sh --html
    else
        if ! check_dependencies "go"; then
            press_enter
            return
        fi

        log_info "Executando testes com cobertura..."
        if go test ./... -coverprofile=coverage.out; then
            log_info "Gerando relatório HTML..."
            go tool cover -html=coverage.out -o coverage.html
            log_success "Relatório HTML gerado: coverage.html"

            # Tentar abrir o relatório no navegador
            if command -v xdg-open &> /dev/null; then
                xdg-open coverage.html
            elif command -v open &> /dev/null; then
                open coverage.html
            else
                log_info "Abra o arquivo coverage.html no seu navegador"
            fi
        else
            log_error "Falha ao executar testes"
        fi
    fi
    press_enter
}

tests_run_specific() {
    log_step "Executar teste específico..."
    echo ""
    read -p "Digite o nome do teste ou padrão (ex: TestUserRoutes): " test_pattern

    if [ -z "$test_pattern" ]; then
        log_error "Nome do teste não pode estar vazio"
        press_enter
        return
    fi

    cd "$ROOT_DIR"

    if [ -f "scripts/run_tests.sh" ]; then
        chmod +x scripts/run_tests.sh
        ./scripts/run_tests.sh --test "$test_pattern"
    else
        if ! check_dependencies "go"; then
            press_enter
            return
        fi

        log_info "Executando teste: $test_pattern"
        if go test ./... -run "$test_pattern" -v; then
            log_success "Teste executado com sucesso!"
        else
            log_error "Teste falhou"
        fi
    fi
    press_enter
}

tests_run_verbose() {
    log_step "Executando testes com saída detalhada..."
    cd "$ROOT_DIR"

    if [ -f "scripts/run_tests.sh" ]; then
        chmod +x scripts/run_tests.sh
        ./scripts/run_tests.sh --verbose
    else
        if ! check_dependencies "go"; then
            press_enter
            return
        fi

        log_info "Executando testes verbosos..."
        go test ./... -v -count=1
    fi
    press_enter
}

tests_run_fast() {
    log_step "Executando testes rápidos (sem cache)..."
    cd "$ROOT_DIR"

    if ! check_dependencies "go"; then
        press_enter
        return
    fi

    log_info "Limpando cache de testes..."
    go clean -testcache

    log_info "Executando testes sem cache..."
    if go test ./... -count=1; then
        log_success "Testes rápidos concluídos!"
    else
        log_error "Alguns testes falharam"
    fi
    press_enter
}

# ==============================================================================
# 5. DEPLOY GCP
# ==============================================================================

show_deploy_menu() {
    show_main_banner
    log_header "deploy" "Deploy GCP"
    echo ""
    echo "  1. 🚀 Deploy interativo completo"
    echo "  2. ⚡ Deploy rápido (quick-deploy)"
    echo "  3. 🏗️  Bootstrap inicial do GCP"
    echo "  4. 🏗️  Deploy apenas infraestrutura"
    echo "  5. 📦 Deploy apenas aplicação"
    echo "  6. 🔄 Atualizar segredos (secrets)"
    echo "  7. 📊 Status dos serviços"
    echo "  0. ⬅️  Voltar ao menu principal"
    echo ""
}

handle_deploy_menu() {
    while true; do
        show_deploy_menu
        read -p "Selecione uma opção: " choice

        case $choice in
            1) deploy_interactive_complete ;;
            2) deploy_quick ;;
            3) deploy_bootstrap_gcp ;;
            4) deploy_infrastructure_only ;;
            5) deploy_application_only ;;
            6) deploy_update_secrets ;;
            7) deploy_services_status ;;
            0) return ;;
            *) log_error "Opção inválida. Tente novamente." ; press_enter ;;
        esac
    done
}

deploy_interactive_complete() {
    log_step "Iniciando deploy interativo completo..."
    cd "$ROOT_DIR"

    if [ -f "scripts/deploy-interactive.sh" ]; then
        chmod +x scripts/deploy-interactive.sh
        ./scripts/deploy-interactive.sh
    else
        log_error "Script deploy-interactive.sh não encontrado"
        press_enter
    fi
}

deploy_quick() {
    log_step "Executando deploy rápido..."
    cd "$ROOT_DIR"

    echo ""
    echo "Ambientes disponíveis:"
    echo "  1. dev (desenvolvimento)"
    echo "  2. staging (homologação)"
    echo "  3. prod (produção)"
    echo ""

    read -p "Selecione o ambiente (1-3): " env_choice

    local environment=""
    case $env_choice in
        1) environment="dev" ;;
        2) environment="staging" ;;
        3) environment="prod" ;;
        *) log_error "Ambiente inválido" ; press_enter ; return ;;
    esac

    if confirm_action "Deploy para ambiente '$environment'?"; then
        if [ -f "scripts/quick-deploy.sh" ]; then
            chmod +x scripts/quick-deploy.sh
            ENVIRONMENT="$environment" ./scripts/quick-deploy.sh
        else
            log_error "Script quick-deploy.sh não encontrado"
        fi
    fi
    press_enter
}

deploy_bootstrap_gcp() {
    log_step "Executando bootstrap inicial do GCP..."
    cd "$ROOT_DIR"

    log_warn "O bootstrap criará recursos básicos no GCP como:"
    log_info "  - Service Account"
    log_info "  - Artifact Registry"
    log_info "  - Secrets Manager"
    log_info "  - Habilitação de APIs necessárias"

    if confirm_action "Continuar com o bootstrap?"; then
        if [ -f "scripts/bootstrap-gcp.sh" ]; then
            chmod +x scripts/bootstrap-gcp.sh
            ./scripts/bootstrap-gcp.sh
        else
            log_error "Script bootstrap-gcp.sh não encontrado"
        fi
    fi
    press_enter
}

deploy_infrastructure_only() {
    log_step "Deploy apenas da infraestrutura (Terraform)..."
    cd "$ROOT_DIR"

    if ! check_dependencies "terraform"; then
        press_enter
        return
    fi

    if ! check_dependencies "gcloud"; then
        press_enter
        return
    fi

    # Verificar se terraform.tfvars existe
    if [ ! -f "terraform.tfvars" ]; then
        log_error "Arquivo terraform.tfvars não encontrado"
        log_info "Crie o arquivo com as configurações necessárias"
        press_enter
        return
    fi

    if confirm_action "Deploy apenas da infraestrutura?"; then
        log_info "Inicializando Terraform..."
        terraform init

        log_info "Criando plano de execução..."
        terraform plan -out=tfplan

        log_info "Aplicando infraestrutura..."
        terraform apply tfplan

        log_success "Deploy da infraestrutura concluído!"
    fi
    press_enter
}

deploy_application_only() {
    log_step "Deploy apenas da aplicação (Cloud Run)..."
    cd "$ROOT_DIR"

    if ! check_dependencies "docker"; then
        press_enter
        return
    fi

    if ! check_dependencies "gcloud"; then
        press_enter
        return
    fi

    local image_tag="${REGION}-docker.pkg.dev/${PROJECT_ID}/lep-backend/lep-backend:latest"

    if confirm_action "Build e deploy da aplicação?"; then
        log_info "Construindo imagem Docker..."
        docker build -t "$image_tag" .

        log_info "Enviando imagem para Artifact Registry..."
        docker push "$image_tag"

        log_info "Fazendo deploy para Cloud Run..."
        gcloud run deploy "leps-backend-dev" \
            --image="$image_tag" \
            --region="$REGION" \
            --platform=managed \
            --allow-unauthenticated

        log_success "Deploy da aplicação concluído!"

        # Obter URL do serviço
        local service_url=$(gcloud run services describe "leps-backend-dev" \
            --region="$REGION" \
            --format="value(status.url)" 2>/dev/null)

        if [ -n "$service_url" ]; then
            log_info "URL do serviço: $service_url"
            log_info "Health check: $service_url/health"
        fi
    fi
    press_enter
}

deploy_update_secrets() {
    log_step "Atualizando segredos no Google Cloud..."
    cd "$ROOT_DIR"

    if ! check_dependencies "gcloud"; then
        press_enter
        return
    fi

    echo ""
    echo "Segredos disponíveis para atualização:"
    echo "  1. JWT Private Key"
    echo "  2. JWT Public Key"
    echo "  3. Database Password"
    echo "  4. Twilio Auth Token"
    echo "  5. SMTP Password"
    echo ""

    read -p "Qual segredo atualizar (1-5)? " secret_choice

    local secret_name=""
    local secret_file=""

    case $secret_choice in
        1)
            secret_name="jwt-private-key-dev"
            secret_file="jwt_private_key.pem"
            ;;
        2)
            secret_name="jwt-public-key-dev"
            secret_file="jwt_public_key.pem"
            ;;
        3)
            secret_name="db-password-dev"
            log_warn "Digite a senha do database:"
            read -s db_password
            echo "$db_password" | gcloud secrets versions add "$secret_name" --data-file=-
            log_success "Senha do database atualizada!"
            press_enter
            return
            ;;
        4|5)
            log_error "Atualização manual necessária via Console GCP"
            press_enter
            return
            ;;
        *)
            log_error "Opção inválida"
            press_enter
            return
            ;;
    esac

    if [ -n "$secret_file" ] && [ -f "$secret_file" ]; then
        if confirm_action "Atualizar $secret_name com o arquivo $secret_file?"; then
            gcloud secrets versions add "$secret_name" --data-file="$secret_file"
            log_success "Segredo $secret_name atualizado!"
        fi
    else
        log_error "Arquivo $secret_file não encontrado"
    fi

    press_enter
}

deploy_services_status() {
    log_step "Verificando status dos serviços..."

    if ! check_dependencies "gcloud"; then
        press_enter
        return
    fi

    echo ""
    log_info "Serviços Cloud Run:"
    gcloud run services list --region="$REGION" --format="table(metadata.name,status.url,status.traffic[0].percent)"

    echo ""
    log_info "Últimas revisões:"
    gcloud run revisions list --region="$REGION" --limit=5 --format="table(metadata.name,status.conditions[0].status,metadata.creationTimestamp)"

    echo ""
    log_info "Health check dos serviços ativos:"
    local services=$(gcloud run services list --region="$REGION" --format="value(status.url)" 2>/dev/null)

    for service_url in $services; do
        if [ -n "$service_url" ]; then
            local service_name=$(echo "$service_url" | sed 's|https://||' | cut -d'.' -f1)
            if curl -s --connect-timeout 5 "$service_url/health" &> /dev/null; then
                echo -e "  ✅ $service_name: saudável"
            else
                echo -e "  ❌ $service_name: não responsivo"
            fi
        fi
    done

    press_enter
}

# ==============================================================================
# 6. UTILITÁRIOS
# ==============================================================================

show_utilities_menu() {
    show_main_banner
    log_header "utilities" "Utilitários"
    echo ""
    echo "  1. 🔍 Verificar dependências do sistema"
    echo "  2. 🧹 Limpeza completa do projeto"
    echo "  3. 💾 Backup de configurações"
    echo "  4. 📊 Status geral do projeto"
    echo "  5. 🩺 Troubleshooting automático"
    echo "  6. 📝 Gerar relatório do ambiente"
    echo "  7. 🔄 Reset completo para desenvolvimento"
    echo "  0. ⬅️  Voltar ao menu principal"
    echo ""
}

handle_utilities_menu() {
    while true; do
        show_utilities_menu
        read -p "Selecione uma opção: " choice

        case $choice in
            1) utilities_check_dependencies ;;
            2) utilities_complete_cleanup ;;
            3) utilities_backup_configs ;;
            4) utilities_project_status ;;
            5) utilities_troubleshooting ;;
            6) utilities_environment_report ;;
            7) utilities_reset_development ;;
            0) return ;;
            *) log_error "Opção inválida. Tente novamente." ; press_enter ;;
        esac
    done
}

utilities_check_dependencies() {
    log_step "Verificação completa de dependências..."

    echo ""
    echo -e "${CYAN}=== DEPENDÊNCIAS OBRIGATÓRIAS ===${NC}"
    local required_tools=("go" "git" "curl" "docker" "gcloud" "terraform")
    local missing_required=0

    for tool in "${required_tools[@]}"; do
        if command -v "$tool" &> /dev/null; then
            local version=$($tool version 2>/dev/null | head -n1 || echo "versão não detectada")
            echo -e "  ✅ ${tool}: ${version}"
        else
            echo -e "  ❌ ${tool}: NÃO INSTALADO"
            missing_required=$((missing_required + 1))
        fi
    done

    echo ""
    echo -e "${CYAN}=== DEPENDÊNCIAS OPCIONAIS ===${NC}"
    local optional_tools=("openssl" "jq" "python3" "make" "wget")

    for tool in "${optional_tools[@]}"; do
        if command -v "$tool" &> /dev/null; then
            echo -e "  ✅ ${tool}: disponível"
        else
            echo -e "  ⚠️  ${tool}: não instalado"
        fi
    done

    echo ""
    echo -e "${CYAN}=== VERIFICAÇÕES DE SISTEMA ===${NC}"

    # Verificar Docker
    if docker info &> /dev/null; then
        echo -e "  ✅ Docker: rodando"
    else
        echo -e "  ❌ Docker: não está rodando"
        missing_required=$((missing_required + 1))
    fi

    # Verificar Google Cloud
    if gcloud projects describe "$PROJECT_ID" &> /dev/null; then
        echo -e "  ✅ Google Cloud: autenticado (projeto: $PROJECT_ID)"
    else
        echo -e "  ⚠️  Google Cloud: não autenticado ou sem acesso ao projeto"
    fi

    # Verificar Go modules
    cd "$ROOT_DIR"
    if go mod verify &> /dev/null; then
        echo -e "  ✅ Go modules: válidos"
    else
        echo -e "  ⚠️  Go modules: podem precisar de atualização"
    fi

    echo ""
    if [ $missing_required -eq 0 ]; then
        log_success "🎉 Todas as dependências obrigatórias estão instaladas!"
    else
        log_error "❌ Faltam $missing_required dependências obrigatórias"
        echo ""
        echo "Para instalar as dependências em falta:"
        echo "  - Docker: https://docs.docker.com/get-docker/"
        echo "  - Google Cloud CLI: https://cloud.google.com/sdk/docs/install"
        echo "  - Terraform: https://learn.hashicorp.com/tutorials/terraform/install-cli"
        echo "  - Go: https://golang.org/doc/install"
    fi

    press_enter
}

utilities_complete_cleanup() {
    log_step "Limpeza completa do projeto..."

    log_warn "⚠️  Esta operação irá:"
    log_info "  - Remover todos os binários compilados"
    log_info "  - Limpar cache do Go"
    log_info "  - Remover arquivos temporários do Terraform"
    log_info "  - Limpar imagens Docker locais"
    log_info "  - Remover logs antigos"

    if confirm_action "Continuar com a limpeza completa?"; then
        cd "$ROOT_DIR"

        # Remover binários
        log_info "Removendo binários..."
        rm -rf bin/ 2>/dev/null || true
        find . -name "*.exe" -delete 2>/dev/null || true

        # Limpar Go
        log_info "Limpando cache do Go..."
        go clean -cache 2>/dev/null || true
        go clean -modcache 2>/dev/null || log_warn "Cache de módulos mantido (requer permissões)"

        # Limpar Terraform
        log_info "Limpando arquivos do Terraform..."
        rm -f tfplan terraform.tfstate.backup 2>/dev/null || true
        rm -f .terraform.lock.hcl terraform.tfstate.lock.info 2>/dev/null || true

        # Limpar Docker (cuidadosamente)
        if docker info &> /dev/null; then
            log_info "Removendo imagens Docker não utilizadas..."
            docker system prune -f &> /dev/null || true
        fi

        # Limpar logs
        log_info "Removendo logs antigos..."
        rm -rf logs/ 2>/dev/null || true
        find . -name "*.log" -delete 2>/dev/null || true

        # Limpar arquivos temporários
        log_info "Removendo arquivos temporários..."
        find . -name "*.tmp" -delete 2>/dev/null || true
        find . -name "*.temp" -delete 2>/dev/null || true
        rm -f coverage.out coverage.html 2>/dev/null || true

        log_success "🎉 Limpeza completa concluída!"

        # Mostrar espaço liberado se possível
        if command -v du &> /dev/null; then
            log_info "Espaço total do projeto: $(du -sh . 2>/dev/null | cut -f1)"
        fi
    fi

    press_enter
}

utilities_backup_configs() {
    log_step "Backup de configurações..."
    cd "$ROOT_DIR"

    local backup_dir="backup-$(date +%Y%m%d-%H%M%S)"
    mkdir -p "$backup_dir"

    log_info "Criando backup em: $backup_dir"

    # Arquivos de configuração importantes
    local config_files=(
        ".env"
        "terraform.tfvars"
        "jwt_private_key.pem"
        "jwt_public_key.pem"
        "docker-compose.yml"
        "Dockerfile"
        "go.mod"
        "go.sum"
    )

    local backed_up=0
    for file in "${config_files[@]}"; do
        if [ -f "$file" ]; then
            cp "$file" "$backup_dir/" 2>/dev/null && {
                echo -e "  ✅ $file"
                backed_up=$((backed_up + 1))
            }
        else
            echo -e "  ⚠️  $file (não encontrado)"
        fi
    done

    # Backup de diretórios importantes
    local config_dirs=("environments" "scripts")
    for dir in "${config_dirs[@]}"; do
        if [ -d "$dir" ]; then
            cp -r "$dir" "$backup_dir/" 2>/dev/null && {
                echo -e "  ✅ $dir/ (diretório)"
                backed_up=$((backed_up + 1))
            }
        fi
    done

    if [ $backed_up -gt 0 ]; then
        log_success "Backup criado com sucesso!"
        log_info "Arquivos salvos em: $backup_dir"
        log_info "Total de itens: $backed_up"
    else
        log_warn "Nenhum arquivo foi feito backup"
        rmdir "$backup_dir" 2>/dev/null || true
    fi

    press_enter
}

utilities_project_status() {
    log_step "Status geral do projeto..."
    cd "$ROOT_DIR"

    echo ""
    echo -e "${CYAN}=== STATUS DO PROJETO ===${NC}"

    # Informações básicas
    echo -e "${WHITE}Diretório:${NC} $(pwd)"
    echo -e "${WHITE}Projeto:${NC} $PROJECT_NAME ($PROJECT_ID)"

    # Status do Git
    if [ -d ".git" ]; then
        local branch=$(git branch --show-current 2>/dev/null || echo "unknown")
        local status=$(git status --porcelain 2>/dev/null | wc -l || echo "0")
        echo -e "${WHITE}Git:${NC} branch '$branch', $status arquivos modificados"
    else
        echo -e "${WHITE}Git:${NC} não é um repositório Git"
    fi

    echo ""
    echo -e "${CYAN}=== ARQUIVOS DE CONFIGURAÇÃO ===${NC}"

    local config_files=(".env" "terraform.tfvars" "jwt_private_key.pem" "docker-compose.yml")
    for file in "${config_files[@]}"; do
        if [ -f "$file" ]; then
            local size=$(ls -lh "$file" | awk '{print $5}')
            echo -e "  ✅ $file ($size)"
        else
            echo -e "  ❌ $file (não encontrado)"
        fi
    done

    echo ""
    echo -e "${CYAN}=== SERVIÇOS E CONECTIVIDADE ===${NC}"

    # Verificar se a aplicação está rodando
    if curl -s --connect-timeout 3 "http://localhost:8080/ping" &> /dev/null; then
        echo -e "  ✅ Aplicação local: rodando (porta 8080)"
    else
        echo -e "  ❌ Aplicação local: não está rodando"
    fi

    # Verificar Docker
    if docker info &> /dev/null; then
        local containers=$(docker ps -q | wc -l)
        echo -e "  ✅ Docker: $containers containers rodando"
    else
        echo -e "  ❌ Docker: não está rodando"
    fi

    # Verificar Google Cloud
    if gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q "@" 2>/dev/null; then
        echo -e "  ✅ Google Cloud: autenticado"
    else
        echo -e "  ❌ Google Cloud: não autenticado"
    fi

    echo ""
    echo -e "${CYAN}=== ESTATÍSTICAS DO CÓDIGO ===${NC}"

    if [ -f "go.mod" ]; then
        local go_files=$(find . -name "*.go" -not -path "./vendor/*" | wc -l)
        local test_files=$(find . -name "*_test.go" | wc -l)
        echo -e "  📄 Arquivos Go: $go_files"
        echo -e "  🧪 Arquivos de teste: $test_files"

        # Dependências Go
        local dependencies=$(go list -m all 2>/dev/null | wc -l || echo "0")
        echo -e "  📦 Dependências: $dependencies"
    fi

    # Tamanho do projeto
    if command -v du &> /dev/null; then
        local project_size=$(du -sh . 2>/dev/null | cut -f1)
        echo -e "  💾 Tamanho total: $project_size"
    fi

    press_enter
}

utilities_troubleshooting() {
    log_step "Executando troubleshooting automático..."

    echo ""
    echo -e "${CYAN}=== DIAGNÓSTICO AUTOMÁTICO ===${NC}"

    local issues_found=0

    # 1. Verificar Go
    if ! command -v go &> /dev/null; then
        echo -e "  ❌ Go não instalado"
        issues_found=$((issues_found + 1))
    else
        if ! go version | grep -q "go1\." &> /dev/null; then
            echo -e "  ⚠️  Versão do Go pode estar desatualizada"
        else
            echo -e "  ✅ Go: OK"
        fi
    fi

    # 2. Verificar arquivos essenciais
    cd "$ROOT_DIR"
    local essential_files=("go.mod" "main.go")
    for file in "${essential_files[@]}"; do
        if [ ! -f "$file" ]; then
            echo -e "  ❌ Arquivo essencial não encontrado: $file"
            issues_found=$((issues_found + 1))
        fi
    done

    # 3. Verificar dependências Go
    if [ -f "go.mod" ]; then
        if ! go mod verify &> /dev/null; then
            echo -e "  ⚠️  Módulos Go precisam ser atualizados"
            log_info "Executando 'go mod tidy'..."
            go mod tidy
        else
            echo -e "  ✅ Módulos Go: OK"
        fi
    fi

    # 4. Verificar build
    if ! go build -o /tmp/test-build . &> /dev/null; then
        echo -e "  ❌ Falha no build do Go"
        issues_found=$((issues_found + 1))
        log_info "Tentando identificar o erro..."
        go build . 2>&1 | head -5
    else
        echo -e "  ✅ Build do Go: OK"
        rm -f /tmp/test-build
    fi

    # 5. Verificar portas em uso
    if lsof -i :8080 &> /dev/null || netstat -an 2>/dev/null | grep :8080 &> /dev/null; then
        echo -e "  ⚠️  Porta 8080 já está em uso"
    else
        echo -e "  ✅ Porta 8080: disponível"
    fi

    # 6. Verificar Docker
    if command -v docker &> /dev/null; then
        if ! docker info &> /dev/null; then
            echo -e "  ⚠️  Docker instalado mas não está rodando"
        else
            echo -e "  ✅ Docker: OK"
        fi
    fi

    # 7. Verificar permissões de arquivos
    if [ ! -x "scripts/run_seed.sh" ] 2>/dev/null; then
        log_info "Corrigindo permissões de scripts..."
        chmod +x scripts/*.sh 2>/dev/null || true
    fi

    echo ""
    if [ $issues_found -eq 0 ]; then
        log_success "🎉 Nenhum problema crítico encontrado!"
        log_info "Seu ambiente parece estar configurado corretamente."
    else
        log_warn "⚠️  Encontrados $issues_found problemas."
        echo ""
        echo -e "${CYAN}Soluções sugeridas:${NC}"
        echo "  1. Execute o setup completo: Menu 2 -> Opção 1"
        echo "  2. Verifique as dependências: Menu 6 -> Opção 1"
        echo "  3. Reconfigure o ambiente: Menu 2 -> Opção 6"
    fi

    press_enter
}

utilities_environment_report() {
    log_step "Gerando relatório do ambiente..."
    cd "$ROOT_DIR"

    local report_file="environment-report-$(date +%Y%m%d-%H%M%S).txt"

    log_info "Criando relatório detalhado..."

    {
        echo "LEP System - Environment Report"
        echo "Generated on: $(date)"
        echo "======================================"
        echo ""

        echo "=== SYSTEM INFO ==="
        echo "OS: $(uname -a)"
        echo "User: $(whoami)"
        echo "PWD: $(pwd)"
        echo "Shell: $SHELL"
        echo ""

        echo "=== PROJECT INFO ==="
        echo "Project ID: $PROJECT_ID"
        echo "Project Name: $PROJECT_NAME"
        echo "Region: $REGION"
        echo ""

        echo "=== DEPENDENCIES ==="
        for tool in go git docker gcloud terraform; do
            if command -v "$tool" &> /dev/null; then
                echo "$tool: $($tool version 2>/dev/null | head -1)"
            else
                echo "$tool: NOT INSTALLED"
            fi
        done
        echo ""

        echo "=== GO INFO ==="
        if command -v go &> /dev/null; then
            go env
        else
            echo "Go not installed"
        fi
        echo ""

        echo "=== PROJECT FILES ==="
        ls -la
        echo ""

        echo "=== GO MODULES ==="
        if [ -f "go.mod" ]; then
            cat go.mod
            echo ""
            echo "Dependencies:"
            go list -m all 2>/dev/null || echo "Error listing dependencies"
        fi
        echo ""

        echo "=== DOCKER INFO ==="
        if command -v docker &> /dev/null && docker info &> /dev/null; then
            docker --version
            docker images | head -10
        else
            echo "Docker not available"
        fi
        echo ""

        echo "=== GCLOUD INFO ==="
        if command -v gcloud &> /dev/null; then
            gcloud config list
            echo ""
            gcloud auth list
        else
            echo "gcloud not available"
        fi
        echo ""

        echo "=== ENVIRONMENT VARIABLES ==="
        env | grep -E "(GO|GCLOUD|DOCKER|PATH)" | sort
        echo ""

        echo "=== NETWORK ==="
        echo "Listening ports:"
        netstat -an 2>/dev/null | grep LISTEN | head -10 || lsof -i -P -n | grep LISTEN | head -10 || echo "Cannot determine listening ports"

    } > "$report_file"

    log_success "Relatório gerado: $report_file"

    # Mostrar resumo na tela
    echo ""
    log_info "Resumo do ambiente:"
    grep -E "(OS:|Dependencies:|Project)" "$report_file" | head -10

    press_enter
}

utilities_reset_development() {
    log_step "Reset completo para desenvolvimento..."

    log_warn "⚠️  ATENÇÃO: Esta operação irá:"
    log_info "  - Remover TODOS os arquivos de build"
    log_info "  - Limpar TODOS os caches"
    log_info "  - Resetar configurações de desenvolvimento"
    log_info "  - Recriar arquivos de configuração padrão"
    log_info "  - Executar setup inicial"

    echo ""
    log_error "ARQUIVOS QUE SERÃO PRESERVADOS:"
    log_info "  - Código fonte (.go files)"
    log_info "  - Scripts"
    log_info "  - Documentação"
    log_info "  - Chaves JWT existentes"

    if confirm_action "CONFIRMA o reset completo para desenvolvimento?"; then
        cd "$ROOT_DIR"

        # 1. Limpeza completa
        log_info "1. Executando limpeza completa..."
        rm -rf bin/ logs/ coverage.out coverage.html 2>/dev/null || true
        go clean -cache -modcache 2>/dev/null || true

        # 2. Reset Docker
        if docker info &> /dev/null; then
            log_info "2. Limpando containers e imagens não utilizadas..."
            docker system prune -f &> /dev/null || true
        fi

        # 3. Reset Terraform
        log_info "3. Limpando estado do Terraform..."
        rm -f tfplan terraform.tfstate.backup .terraform.lock.hcl 2>/dev/null || true

        # 4. Recriar .env
        if [ -f ".env" ]; then
            log_info "4. Fazendo backup do .env atual..."
            mv .env .env.backup-$(date +%H%M%S) 2>/dev/null || true
        fi

        log_info "5. Criando novo .env padrão..."
        cat > .env << 'EOF'
# LEP System - Development Environment

# Database Configuration
DB_USER=postgres
DB_PASS=password
DB_NAME=lep_database
DB_HOST=localhost
DB_PORT=5432

# JWT Configuration (update after generating keys)
JWT_SECRET_PRIVATE_KEY=your_jwt_private_key_here
JWT_SECRET_PUBLIC_KEY=your_jwt_public_key_here

# Application Configuration
PORT=8080
ENABLE_CRON_JOBS=true

# Storage Configuration
STORAGE_TYPE=local
BASE_URL=http://localhost:8080

# Optional: Twilio (for SMS/WhatsApp)
TWILIO_ACCOUNT_SID=
TWILIO_AUTH_TOKEN=
TWILIO_PHONE_NUMBER=

# Optional: SMTP (for Email)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=
SMTP_PASSWORD=
EOF

        # 5. Atualizar dependências
        log_info "6. Atualizando dependências Go..."
        go mod tidy

        # 6. Teste de build
        log_info "7. Testando build..."
        if go build -o bin/lep-test . &> /dev/null; then
            log_success "Build de teste: OK"
            rm -f bin/lep-test
        else
            log_error "Build de teste: FALHOU"
        fi

        log_success "🎉 Reset para desenvolvimento concluído!"
        echo ""
        log_info "Próximos passos:"
        log_info "  1. Configure o arquivo .env com suas credenciais"
        log_info "  2. Gere novas chaves JWT se necessário (Menu 2 -> Opção 3)"
        log_info "  3. Execute o seeding da database (Menu 3 -> Opção 1)"
        log_info "  4. Inicie o servidor (Menu 1 -> Opção 1)"
    fi

    press_enter
}

# ==============================================================================
# 7. AJUDA E INFORMAÇÕES
# ==============================================================================

show_help_menu() {
    show_main_banner
    log_header "help" "Ajuda e Informações"
    echo ""
    echo "  1. 📖 Guia de primeiros passos"
    echo "  2. 🔧 Comandos úteis"
    echo "  3. 🚨 Solução de problemas comuns"
    echo "  4. 📝 Informações do sistema"
    echo "  5. 🌐 Links úteis"
    echo "  6. 📚 Documentação dos scripts"
    echo "  0. ⬅️  Voltar ao menu principal"
    echo ""
}

handle_help_menu() {
    while true; do
        show_help_menu
        read -p "Selecione uma opção: " choice

        case $choice in
            1) help_getting_started ;;
            2) help_useful_commands ;;
            3) help_troubleshooting ;;
            4) help_system_info ;;
            5) help_useful_links ;;
            6) help_scripts_documentation ;;
            0) return ;;
            *) log_error "Opção inválida. Tente novamente." ; press_enter ;;
        esac
    done
}

help_getting_started() {
    show_main_banner
    log_header "help" "Guia de Primeiros Passos"
    echo ""

    echo -e "${CYAN}🚀 Configuração Inicial (primeira vez):${NC}"
    echo "  1. Execute o setup completo: Menu Principal -> 2 -> 1"
    echo "  2. Gere chaves JWT: Menu Principal -> 2 -> 3"
    echo "  3. Configure Google Cloud: Menu Principal -> 2 -> 4"
    echo "  4. Valide a configuração: Menu Principal -> 2 -> 6"
    echo ""

    echo -e "${CYAN}🏠 Desenvolvimento Local:${NC}"
    echo "  1. Popule o database: Menu Principal -> 3 -> 1"
    echo "  2. Inicie o servidor: Menu Principal -> 1 -> 1"
    echo "  3. Teste a API: curl http://localhost:8080/health"
    echo "  4. Execute testes: Menu Principal -> 4 -> 1"
    echo ""

    echo -e "${CYAN}☁️ Deploy para Produção:${NC}"
    echo "  1. Bootstrap GCP: Menu Principal -> 5 -> 3"
    echo "  2. Deploy interativo: Menu Principal -> 5 -> 1"
    echo "  3. Ou deploy rápido: Menu Principal -> 5 -> 2"
    echo ""

    echo -e "${CYAN}💡 Dicas Importantes:${NC}"
    echo "  - Sempre configure o arquivo .env antes de iniciar"
    echo "  - Use 'Menu 6 -> 5' para troubleshooting automático"
    echo "  - Execute 'Menu 6 -> 4' para verificar status geral"
    echo "  - Mantenha backup das configurações: Menu 6 -> 3"
    echo ""

    echo -e "${CYAN}🔑 Credenciais Padrão (após seeding):${NC}"
    echo "  - admin@lep-demo.com / password (Admin)"
    echo "  - garcom@lep-demo.com / password (Garçom)"
    echo "  - gerente@lep-demo.com / password (Gerente)"
    echo ""

    press_enter
}

help_useful_commands() {
    show_main_banner
    log_header "help" "Comandos Úteis"
    echo ""

    echo -e "${CYAN}📦 Comandos Go:${NC}"
    echo "  go run main.go                 # Executar aplicação"
    echo "  go build -o bin/lep .          # Build da aplicação"
    echo "  go test ./...                  # Executar todos os testes"
    echo "  go mod tidy                    # Limpar dependências"
    echo "  go mod verify                  # Verificar dependências"
    echo ""

    echo -e "${CYAN}🐳 Comandos Docker:${NC}"
    echo "  docker build -t lep .          # Build da imagem"
    echo "  docker run -p 8080:8080 lep    # Executar container"
    echo "  docker ps                      # Listar containers rodando"
    echo "  docker system prune -f         # Limpar sistema Docker"
    echo ""

    echo -e "${CYAN}☁️ Comandos Google Cloud:${NC}"
    echo "  gcloud auth login              # Login no GCP"
    echo "  gcloud config set project ID   # Definir projeto"
    echo "  gcloud run services list       # Listar serviços"
    echo "  gcloud logs read SERVICE       # Ver logs do serviço"
    echo ""

    echo -e "${CYAN}🏗️ Comandos Terraform:${NC}"
    echo "  terraform init                 # Inicializar Terraform"
    echo "  terraform plan                 # Planejar mudanças"
    echo "  terraform apply                # Aplicar mudanças"
    echo "  terraform destroy              # Destruir recursos"
    echo ""

    echo -e "${CYAN}🔍 Comandos de Debug:${NC}"
    echo "  curl http://localhost:8080/health    # Health check"
    echo "  curl http://localhost:8080/ping      # Conectividade"
    echo "  lsof -i :8080                       # Ver processo na porta"
    echo "  netstat -tulpn | grep 8080          # Status da porta"
    echo ""

    echo -e "${CYAN}📊 Comandos de Database:${NC}"
    echo "  go run cmd/seed/main.go              # Popular database"
    echo "  go run cmd/seed/main.go --clear-first # Limpar e popular"
    echo ""

    press_enter
}

help_troubleshooting() {
    show_main_banner
    log_header "help" "Solução de Problemas Comuns"
    echo ""

    echo -e "${CYAN}🔧 Problema: Porta 8080 já está em uso${NC}"
    echo "  Solução:"
    echo "    1. lsof -i :8080              # Ver processo na porta"
    echo "    2. kill -9 PID               # Matar processo (substitua PID)"
    echo "    3. Ou mude a porta no .env   # PORT=8081"
    echo ""

    echo -e "${CYAN}🔧 Problema: 'go: command not found'${NC}"
    echo "  Solução:"
    echo "    1. Instale Go: https://golang.org/doc/install"
    echo "    2. Adicione Go ao PATH"
    echo "    3. Reinicie o terminal"
    echo ""

    echo -e "${CYAN}🔧 Problema: Docker não responde${NC}"
    echo "  Solução:"
    echo "    1. Verifique se Docker está rodando"
    echo "    2. Reinicie Docker Desktop"
    echo "    3. Execute: docker info"
    echo ""

    echo -e "${CYAN}🔧 Problema: 'gcloud: command not found'${NC}"
    echo "  Solução:"
    echo "    1. Instale gcloud CLI"
    echo "    2. Execute: gcloud auth login"
    echo "    3. Configure projeto: gcloud config set project PROJECT_ID"
    echo ""

    echo -e "${CYAN}🔧 Problema: Build falha${NC}"
    echo "  Solução:"
    echo "    1. go mod tidy               # Atualizar dependências"
    echo "    2. go mod verify             # Verificar dependências"
    echo "    3. go clean -cache           # Limpar cache"
    echo ""

    echo -e "${CYAN}🔧 Problema: Testes falham${NC}"
    echo "  Solução:"
    echo "    1. Verifique database está rodando"
    echo "    2. Execute seeding: Menu 3 -> 1"
    echo "    3. Verifique .env está configurado"
    echo ""

    echo -e "${CYAN}🔧 Problema: Deploy falha${NC}"
    echo "  Solução:"
    echo "    1. Verifique autenticação: gcloud auth list"
    echo "    2. Verifique permissões do projeto"
    echo "    3. Execute bootstrap: Menu 5 -> 3"
    echo ""

    echo -e "${CYAN}🔧 Problema: JWT errors${NC}"
    echo "  Solução:"
    echo "    1. Gere novas chaves: Menu 2 -> 3"
    echo "    2. Atualize .env com as chaves"
    echo "    3. Reinicie a aplicação"
    echo ""

    press_enter
}

help_system_info() {
    show_main_banner
    log_header "help" "Informações do Sistema"
    echo ""

    echo -e "${CYAN}📋 LEP System Master Script${NC}"
    echo "  Versão: 1.0.0"
    echo "  Autor: LEP Development Team"
    echo "  Última atualização: $(date +%Y-%m-%d)"
    echo ""

    echo -e "${CYAN}🏗️ Projeto:${NC}"
    echo "  Nome: $PROJECT_NAME"
    echo "  ID: $PROJECT_ID"
    echo "  Região: $REGION"
    echo ""

    echo -e "${CYAN}📁 Estrutura:${NC}"
    if [ -d "$ROOT_DIR" ]; then
        echo "  Root: $ROOT_DIR"
        echo "  Scripts: $SCRIPT_DIR"

        # Contar arquivos
        local go_files=$(find "$ROOT_DIR" -name "*.go" -not -path "*/vendor/*" 2>/dev/null | wc -l)
        local scripts=$(find "$SCRIPT_DIR" -name "*.sh" 2>/dev/null | wc -l)

        echo "  Arquivos Go: $go_files"
        echo "  Scripts: $scripts"
    fi
    echo ""

    echo -e "${CYAN}🔧 Funcionalidades:${NC}"
    echo "  ✅ Desenvolvimento local"
    echo "  ✅ Setup automático"
    echo "  ✅ Database seeding"
    echo "  ✅ Sistema de testes"
    echo "  ✅ Deploy GCP"
    echo "  ✅ Utilitários"
    echo "  ✅ Troubleshooting"
    echo ""

    echo -e "${CYAN}🌟 Características:${NC}"
    echo "  - Interface interativa unificada"
    echo "  - Consolidação de todos os scripts existentes"
    echo "  - Validações automáticas"
    echo "  - Tratamento robusto de erros"
    echo "  - Suporte multi-ambiente"
    echo "  - Backup automático de configurações"
    echo ""

    press_enter
}

help_useful_links() {
    show_main_banner
    log_header "help" "Links Úteis"
    echo ""

    echo -e "${CYAN}📚 Documentação Oficial:${NC}"
    echo "  Go Language: https://golang.org/doc/"
    echo "  Docker: https://docs.docker.com/"
    echo "  Google Cloud: https://cloud.google.com/docs"
    echo "  Terraform: https://www.terraform.io/docs"
    echo ""

    echo -e "${CYAN}⚙️ Instalação de Ferramentas:${NC}"
    echo "  Go: https://golang.org/doc/install"
    echo "  Docker: https://docs.docker.com/get-docker/"
    echo "  gcloud CLI: https://cloud.google.com/sdk/docs/install"
    echo "  Terraform: https://learn.hashicorp.com/tutorials/terraform/install-cli"
    echo ""

    echo -e "${CYAN}🎓 Tutoriais:${NC}"
    echo "  Go Tutorial: https://tour.golang.org/"
    echo "  Docker Getting Started: https://docs.docker.com/get-started/"
    echo "  GCP Quickstart: https://cloud.google.com/docs/get-started"
    echo "  Terraform Learn: https://learn.hashicorp.com/terraform"
    echo ""

    echo -e "${CYAN}🔧 Referências da API:${NC}"
    echo "  Gin Framework: https://gin-gonic.com/docs/"
    echo "  GORM: https://gorm.io/docs/"
    echo "  Google Cloud APIs: https://cloud.google.com/apis/docs/overview"
    echo ""

    echo -e "${CYAN}🆘 Suporte:${NC}"
    echo "  Go Community: https://golang.org/help/"
    echo "  Docker Community: https://forums.docker.com/"
    echo "  Stack Overflow: https://stackoverflow.com/"
    echo ""

    press_enter
}

help_scripts_documentation() {
    show_main_banner
    log_header "help" "Documentação dos Scripts"
    echo ""

    echo -e "${CYAN}📜 Scripts Consolidados no Master:${NC}"
    echo ""

    echo -e "${WHITE}1. setup.sh${NC}"
    echo "   - Setup completo do ambiente de desenvolvimento"
    echo "   - Verificação de dependências"
    echo "   - Configuração inicial de arquivos"
    echo ""

    echo -e "${WHITE}2. local-dev.sh${NC}"
    echo "   - Comandos para desenvolvimento local"
    echo "   - Build, run, test, docker, clean"
    echo "   - Geração de chaves JWT"
    echo ""

    echo -e "${WHITE}3. run_seed.sh${NC}"
    echo "   - População da database com dados demo"
    echo "   - Suporte a ambientes (dev/test/staging)"
    echo "   - Opções para limpar e repopular"
    echo ""

    echo -e "${WHITE}4. run_tests.sh${NC}"
    echo "   - Execução de testes com opções avançadas"
    echo "   - Cobertura de código e relatórios HTML"
    echo "   - Testes específicos e verbosos"
    echo ""

    echo -e "${WHITE}5. deploy-interactive.sh${NC}"
    echo "   - Deploy interativo multi-ambiente"
    echo "   - local-dev, gcp-dev, gcp-stage, gcp-prd"
    echo "   - Validações automáticas"
    echo ""

    echo -e "${WHITE}6. quick-deploy.sh${NC}"
    echo "   - Deploy rápido para GCP"
    echo "   - Resolução de conflitos Terraform"
    echo "   - Abordagem híbrida"
    echo ""

    echo -e "${WHITE}7. bootstrap-gcp.sh${NC}"
    echo "   - Criação inicial de recursos GCP"
    echo "   - Service Account, Artifact Registry"
    echo "   - Configuração de APIs"
    echo ""

    echo ""
    echo -e "${CYAN}✨ Vantagens do Master Script:${NC}"
    echo "  - Todas as funcionalidades em um só lugar"
    echo "  - Interface consistente e intuitiva"
    echo "  - Validações e verificações automáticas"
    echo "  - Tratamento robusto de erros"
    echo "  - Suporte a troubleshooting"
    echo ""

    press_enter
}

# ==============================================================================
# MAIN EXECUTION LOGIC
# ==============================================================================

# Handle direct command line arguments for batch execution
handle_batch_mode() {
    case "$1" in
        "--help"|"-h")
            show_main_banner
            echo "LEP System Master Script"
            echo ""
            echo "Uso: $0 [OPÇÃO]"
            echo ""
            echo "Opções:"
            echo "  --help, -h           Mostrar esta ajuda"
            echo "  --quick-deploy       Deploy rápido interativo"
            echo "  --setup              Setup completo do ambiente"
            echo "  --seed               Popular database com dados demo"
            echo "  --test               Executar todos os testes"
            echo "  --status             Mostrar status do projeto"
            echo "  --clean              Limpeza completa"
            echo ""
            echo "Sem argumentos: Iniciar modo interativo"
            exit 0
            ;;
        "--quick-deploy")
            deploy_quick
            exit 0
            ;;
        "--setup")
            setup_complete_environment
            exit 0
            ;;
        "--seed")
            database_seed_demo_data
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
        "--clean")
            utilities_complete_cleanup
            exit 0
            ;;
        "")
            # No arguments - start interactive mode
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
            2) handle_setup_menu ;;
            3) handle_database_menu ;;
            4) handle_tests_menu ;;
            5) handle_deploy_menu ;;
            6) handle_utilities_menu ;;
            7) handle_help_menu ;;
            0)
                echo ""
                log_success "👋 Obrigado por usar o LEP System Master!"
                log_info "Tenha um ótimo desenvolvimento!"
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

# Script initialization and cleanup
init_script() {
    # Ensure we're in the correct directory
    cd "$ROOT_DIR" 2>/dev/null || {
        log_error "Não foi possível acessar o diretório do projeto: $ROOT_DIR"
        log_info "Execute este script a partir do diretório raiz do projeto LEP-Back"
        exit 1
    }

    # Create necessary directories
    mkdir -p bin logs 2>/dev/null || true

    # Set proper permissions for scripts
    chmod +x scripts/*.sh 2>/dev/null || true

    # Handle script interruption gracefully
    trap 'echo ""; log_warn "Script interrompido pelo usuário."; exit 130' INT TERM
}

cleanup_script() {
    # Perform any necessary cleanup before exit
    log_info "Executando limpeza final..."
    # Add cleanup logic here if needed
}

# ==============================================================================
# SCRIPT ENTRY POINT
# ==============================================================================

# Main execution
main() {
    init_script

    # Handle command line arguments
    handle_batch_mode "$1"

    # Start interactive mode
    main_loop

    cleanup_script
}

# Execute main function with all arguments
main "$@"