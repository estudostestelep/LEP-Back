#!/bin/bash
# LEP System - Stage Local (conecta GCP localmente)
# Este script configura variáveis para conectar no GCP rodando go run main.go local

set -e

echo "=== LEP System - Ambiente STAGE Local ==="
echo "• Cloud SQL PostgreSQL (GCP)"
echo "• Google Cloud Storage"
echo "• Execução local via go run main.go"
echo "• Credenciais padronizadas dev/stage"
echo

# Verificar se Go está instalado
if ! command -v go >/dev/null 2>&1; then
    echo "❌ Erro: Go não está instalado. Por favor, instale Go primeiro."
    exit 1
fi

# Verificar se gcloud está configurado
if ! command -v gcloud >/dev/null 2>&1; then
    echo "❌ Erro: gcloud CLI não encontrado. Por favor, instale e configure o Google Cloud SDK."
    echo "   https://cloud.google.com/sdk/docs/install"
    exit 1
fi

# Verificar autenticação
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" >/dev/null 2>&1; then
    echo "❌ Erro: Não autenticado no Google Cloud. Execute:"
    echo "   gcloud auth login"
    echo "   gcloud auth application-default login"
    exit 1
fi

# Verificar projeto ativo
PROJECT_ID=$(gcloud config get-value project 2>/dev/null)
if [ "$PROJECT_ID" != "leps-472702" ]; then
    echo "❌ Erro: Projeto Google Cloud incorreto. Execute:"
    echo "   gcloud config set project leps-472702"
    exit 1
fi

echo "✅ Google Cloud configurado: $PROJECT_ID"

# Configurar variáveis de ambiente para stage
export ENVIRONMENT=stage
export PORT=8080

# Database - Cloud SQL
export DB_USER=lep_user
export DB_PASS=lep_password
export DB_NAME=lep_database
export INSTANCE_UNIX_SOCKET=/cloudsql/leps-472702:us-central1:leps-postgres-stage

# JWT - chaves padronizadas
export JWT_SECRET_PRIVATE_KEY=staging-jwt-private-key-for-testing-only
export JWT_SECRET_PUBLIC_KEY=staging-jwt-public-key-for-testing-only

# Storage - Google Cloud Storage
export STORAGE_TYPE=gcs
export BUCKET_NAME=leps-472702-lep-images-stage
export BASE_URL=https://storage.googleapis.com/leps-472702-lep-images-stage
export BUCKET_CACHE_CONTROL="public, max-age=7200"
export BUCKET_TIMEOUT=60

# Outras configurações
export ENABLE_CRON_JOBS=true
export GIN_MODE=release
export LOG_LEVEL=info

echo "🔧 Variáveis de ambiente configuradas para STAGE"

# Instalar dependências se necessário
echo "📦 Verificando dependências Go..."
go mod tidy

echo "🚀 Iniciando aplicação em modo STAGE local..."
echo "   Conectando em: Cloud SQL + GCS"
echo "   Executando: go run main.go"
echo

# Executar aplicação
go run main.go