# LEP System - Guia de Ambientes

Este documento explica como usar os ambientes padronizados do LEP System.

## Ambientes Disponíveis

### 🔧 DEV - Desenvolvimento Local
- **Objetivo**: Desenvolvimento 100% local sem dependências do GCP
- **Infraestrutura**: Docker (PostgreSQL + Redis + MailHog)
- **Storage**: localStorage (./uploads/)
- **Banco**: PostgreSQL no Docker
- **Uso**: Desenvolvimento diário, testes locais

### 🚀 STAGE - Staging (Local + GCP)
- **Objetivo**: Testes com infraestrutura GCP, execução local OU deploy
- **Infraestrutura**: Cloud SQL + Google Cloud Storage
- **Storage**: Google Cloud Storage
- **Banco**: Cloud SQL PostgreSQL
- **Uso**: Validação antes da produção, testes de integração

### 🏭 PROD - Produção
- **Objetivo**: Ambiente de produção (futuro)
- **Status**: Reservado para configurações profissionais

## Credenciais Padronizadas

Para facilitar testes e validações, **dev** e **stage** usam as mesmas credenciais:

```bash
# Database
DB_NAME=lep_database
DB_USER=lep_user
DB_PASS=lep_password

# JWT (simplificado para dev/stage)
JWT_SECRET_PRIVATE_KEY=dev-simple-private-key-for-testing-only
JWT_SECRET_PUBLIC_KEY=dev-simple-public-key-for-testing-only

# Storage - DEV (local)
STORAGE_TYPE=local
BUCKET_NAME=lep-dev-bucket
BASE_URL=http://localhost:8080
BUCKET_CACHE_CONTROL=public, max-age=3600
BUCKET_TIMEOUT=30

# Storage - STAGE (GCP)
STORAGE_TYPE=gcs
BUCKET_NAME=leps-472702-lep-images-stage
BASE_URL=https://storage.googleapis.com/leps-472702-lep-images-stage
BUCKET_CACHE_CONTROL=public, max-age=7200
BUCKET_TIMEOUT=60

# Cloud SQL - STAGE
INSTANCE_UNIX_SOCKET=/cloudsql/leps-472702:us-central1:leps-postgres-stage

# Application Settings
ENABLE_CRON_JOBS=false  # dev: false, stage: true
GIN_MODE=debug          # dev: debug, stage: release
LOG_LEVEL=debug         # dev: debug, stage: info
```

## Variáveis de Ambiente Completas

### 🔧 Ambiente DEV (.env ou docker-compose)
```bash
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
BUCKET_CACHE_CONTROL=public, max-age=3600
BUCKET_TIMEOUT=30

# SMTP (MailHog)
SMTP_HOST=mailhog
SMTP_PORT=1025
SMTP_USERNAME=
SMTP_PASSWORD=

# Application
ENABLE_CRON_JOBS=false
GIN_MODE=debug
LOG_LEVEL=debug
```

### 🚀 Ambiente STAGE (.env local ou secrets GCP)
```bash
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
BUCKET_CACHE_CONTROL=public, max-age=7200
BUCKET_TIMEOUT=60

# SMTP (opcional)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# Twilio (opcional)
TWILIO_ACCOUNT_SID=your_account_sid
TWILIO_AUTH_TOKEN=your_auth_token
TWILIO_PHONE_NUMBER=+1234567890

# Application
ENABLE_CRON_JOBS=true
GIN_MODE=release
LOG_LEVEL=info
```

## Como Usar

### 🔧 Ambiente DEV

```bash
# Iniciar ambiente local completo
./scripts/dev-local.sh

# Acessos disponíveis:
# API: http://localhost:8080
# MailHog: http://localhost:8025
# PostgreSQL: localhost:5432

# Popular com dados de exemplo
docker-compose run --rm seed

# Parar tudo
docker-compose down
```

### 🚀 Ambiente STAGE

#### Opção 1: Execução Local (conecta GCP)
```bash
# Conectar no GCP rodando localmente
./scripts/stage-local.sh
# Conecta: Cloud SQL + GCS
# Executa: go run main.go local
```

#### Opção 2: Deploy no Cloud Run
```bash
# Deploy no Google Cloud Run
./scripts/stage-deploy.sh
# Deploy: Cloud Run + Cloud SQL + GCS
# Acesso: URL pública do Cloud Run
```

## Configuração do GCP (Stage)

### Pré-requisitos
```bash
# Instalar Google Cloud SDK
# https://cloud.google.com/sdk/docs/install

# Autenticar
gcloud auth login
gcloud auth application-default login

# Configurar projeto
gcloud config set project leps-472702
```

### Criar Infraestrutura
```bash
# Terraform (primeira vez)
cd terraform
terraform init
terraform apply -var-file=../environments/gcp-stage.tfvars
```

## Arquivos de Configuração

### Docker Compose
- `docker-compose.yml` - Configurado para ambiente **dev**

### Terraform
- `environments/gcp-stage.tfvars` - Configurações do ambiente **stage**
- `environments/gcp-prd.tfvars` - Configurações do ambiente **prod** (futuro)

### Scripts Úteis
- `scripts/dev-local.sh` - Iniciar ambiente dev completo
- `scripts/stage-local.sh` - Conectar GCP localmente
- `scripts/stage-deploy.sh` - Deploy no Cloud Run

## Verificação Rápida

### DEV
```bash
curl http://localhost:8080/ping
# Resposta: "pong"
```

### STAGE Local
```bash
curl http://localhost:8080/ping
# Resposta: "pong" (conectado no GCP)
```

### STAGE Deploy
```bash
curl https://lep-system-xxx-uc.a.run.app/ping
# Resposta: "pong" (rodando no Cloud Run)
```

## Troubleshooting

### DEV - Docker
```bash
# Ver logs
docker-compose logs app
docker-compose logs postgres

# Rebuild completo
docker-compose down --volumes
docker-compose build --no-cache
./scripts/dev-local.sh
```

### STAGE - GCP
```bash
# Verificar autenticação
gcloud auth list

# Verificar projeto
gcloud config get-value project

# Ver logs do Cloud Run
gcloud logs read --limit=10

# Verificar Cloud SQL
gcloud sql instances list
```

## Fluxo de Trabalho Recomendado

1. **Desenvolvimento**: Use ambiente **dev** para desenvolvimento diário
2. **Teste Local GCP**: Use **stage local** para validar integração GCP
3. **Teste Deploy**: Use **stage deploy** para validar deploy completo
4. **Produção**: Deploy em **prod** (futuro)

## Migração de Ambientes Antigos

### Se você estava usando:
- `ENVIRONMENT=local-dev` → Use `ENVIRONMENT=dev`
- `ENVIRONMENT=staging` → Use `ENVIRONMENT=stage`
- `gcp-dev.tfvars` → **Removido**, use apenas `gcp-stage.tfvars`

### Comandos de migração:
```bash
# Parar containers antigos
docker-compose down --volumes

# Usar novo ambiente dev
./scripts/dev-local.sh
```