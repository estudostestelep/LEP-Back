# LEP System - Deployment Guide

Guia completo para deploy da infraestrutura e aplicação LEP System no Google Cloud Platform.

---

## 📋 Índice

1. [Pré-requisitos](#pré-requisitos)
2. [Arquitetura](#arquitetura)
3. [Deploy Inicial (Bootstrap)](#deploy-inicial-bootstrap)
4. [Migrations](#migrations)
5. [Seeding de Dados](#seeding-de-dados)
6. [Deploy da Aplicação](#deploy-da-aplicação)
7. [Ambientes](#ambientes)
8. [Rollback](#rollback)
9. [Troubleshooting](#troubleshooting)

---

## Pré-requisitos

### Software Necessário

- **Google Cloud SDK** (gcloud CLI)
  ```bash
  # Verificar instalação
  gcloud --version

  # Instalar se necessário
  # https://cloud.google.com/sdk/docs/install
  ```

- **Terraform** (>= 1.0)
  ```bash
  # Verificar instalação
  terraform --version

  # Instalar se necessário
  # https://www.terraform.io/downloads
  ```

- **Go** (>= 1.24.0)
  ```bash
  # Verificar instalação
  go version

  # Instalar se necessário
  # https://go.dev/dl/
  ```

### Autenticação GCP

```bash
# 1. Login no Google Cloud
gcloud auth login

# 2. Configurar projeto
gcloud config set project leps-472702

# 3. Ativar Application Default Credentials (para Terraform)
gcloud auth application-default login
```

### Variáveis de Ambiente

Criar arquivo `.env` na raiz do projeto:

```bash
# JWT Keys (gerar com openssl)
JWT_SECRET_PRIVATE_KEY="sua_chave_privada_aqui"
JWT_SECRET_PUBLIC_KEY="sua_chave_publica_aqui"

# Twilio (opcional)
TWILIO_ACCOUNT_SID="seu_account_sid"
TWILIO_AUTH_TOKEN="seu_auth_token"
TWILIO_PHONE_NUMBER="+5511999999999"

# SMTP (opcional)
SMTP_HOST="smtp.gmail.com"
SMTP_PORT=587
SMTP_USERNAME="seu_email@gmail.com"
SMTP_PASSWORD="sua_senha_app"
```

Para gerar JWT keys:
```bash
# Private Key
openssl genrsa -out private.key 2048
cat private.key

# Public Key
openssl rsa -in private.key -pubout -out public.key
cat public.key
```

---

## Arquitetura

### Componentes GCP

```
┌─────────────────────────────────────────────────────────────┐
│                      Google Cloud Platform                   │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐      ┌──────────────┐                    │
│  │  Cloud Run   │─────▶│  Cloud SQL   │                    │
│  │  (lep-system)│      │  (PostgreSQL)│                    │
│  └──────────────┘      └──────────────┘                    │
│         │                                                    │
│         ├───────────────┐                                   │
│         │               │                                   │
│  ┌──────▼─────┐  ┌─────▼────────┐                         │
│  │   GCS      │  │ Secret       │                          │
│  │  (Images)  │  │ Manager      │                          │
│  └────────────┘  └──────────────┘                          │
│                                                              │
│  Service Account: lep-backend-sa                            │
│  Permissions: SQL Client, Storage Admin, Secret Accessor   │
└─────────────────────────────────────────────────────────────┘
```

### Recursos Gerenciados

| Recurso | Nome | Finalidade |
|---------|------|------------|
| **Cloud SQL** | `leps-postgres-{env}` | Banco de dados PostgreSQL 15 |
| **GCS Bucket** | `leps-472702-lep-images-{env}` | Armazenamento de imagens |
| **Service Account** | `leps-backend-sa` | Identidade do backend |
| **Secret Manager** | `db-password-{env}` | Senha do banco |
| **Secret Manager** | `jwt-private-key-{env}` | Chave privada JWT |
| **Secret Manager** | `jwt-public-key-{env}` | Chave pública JWT |
| **Cloud Run** | `lep-system` | Aplicação backend |

---

## Deploy Inicial (Bootstrap)

### 1. Preparar Configuração do Ambiente

Editar `environments/gcp-{env}.tfvars`:

```hcl
# environments/gcp-stage.tfvars
project_id  = "leps-472702"
environment = "stage"
region      = "us-central1"

# Database
database_name = "lep_database"
database_user = "lep_user"
db_tier       = "db-f1-micro"

# JWT Keys (copiar do .env)
jwt_private_key = "-----BEGIN RSA PRIVATE KEY-----\n..."
jwt_public_key  = "-----BEGIN PUBLIC KEY-----\n..."

# Storage
bucket_name = "leps-472702-lep-images-stage"

# Opcional: Twilio e SMTP
twilio_account_sid  = ""
twilio_auth_token   = ""
smtp_password       = ""
```

### 2. Inicializar Terraform

```bash
# Navegar para diretório raiz
cd /path/to/LEP-Back

# Inicializar Terraform
terraform init

# Validar configuração
terraform validate

# Ver plano de execução
terraform plan -var-file=environments/gcp-stage.tfvars
```

### 3. Criar Infraestrutura

```bash
# Aplicar configuração
terraform apply -var-file=environments/gcp-stage.tfvars

# Confirmar quando solicitado
# Digite: yes
```

**Tempo estimado**: 10-15 minutos (Cloud SQL leva mais tempo)

### 4. Verificar Recursos Criados

```bash
# Ver outputs
terraform output

# Listar instâncias Cloud SQL
gcloud sql instances list

# Listar buckets
gcloud storage buckets list

# Listar secrets
gcloud secrets list
```

### 5. Importar Recursos Existentes (Se Aplicável)

Se você já tem recursos criados manualmente e quer gerenciá-los pelo Terraform:

```bash
# Cloud SQL Instance
terraform import -var-file=environments/gcp-stage.tfvars \
  google_sql_database_instance.main leps-472702/leps-postgres-dev

# GCS Bucket
terraform import -var-file=environments/gcp-stage.tfvars \
  google_storage_bucket.images leps-472702-lep-images-stage

# Service Account
terraform import -var-file=environments/gcp-stage.tfvars \
  google_service_account.backend projects/leps-472702/serviceAccounts/lep-backend-sa@leps-472702.iam.gserviceaccount.com

# Verificar que não há mudanças pendentes
terraform plan -var-file=environments/gcp-stage.tfvars
# Output deve mostrar: "No changes. Your infrastructure matches the configuration."
```

---

## Migrations

### Executar Migrations Localmente

```bash
# Build do migration tool
go build -o lep-migrate.exe cmd/migrate/main.go

# Dry-run (ver o que será migrado)
ENVIRONMENT=dev ./lep-migrate.exe --dry-run

# Executar migration
ENVIRONMENT=dev ./lep-migrate.exe --verbose
```

### Executar Migrations em Staging/Produção

#### Opção A: Via Cloud SQL Proxy (Recomendado)

```bash
# 1. Baixar Cloud SQL Proxy
# Windows: https://dl.google.com/cloudsql/cloud_sql_proxy_x64.exe

# 2. Iniciar proxy em terminal separado
cloud-sql-proxy leps-472702:us-central1:leps-postgres-stage

# 3. Em outro terminal, configurar variáveis
set ENVIRONMENT=stage
set DB_HOST=127.0.0.1
set DB_PORT=5432
set DB_USER=lep_user
set DB_NAME=lep_database
set DB_PASS=<senha_do_banco>

# 4. Executar migration
./lep-migrate.exe --verbose
```

#### Opção B: Via Cloud Run Job (Produção)

```bash
# 1. Build container
gcloud builds submit --tag gcr.io/leps-472702/lep-migrate:latest \
  --project leps-472702

# 2. Criar Cloud Run Job
gcloud run jobs create lep-migrate \
  --image=gcr.io/leps-472702/lep-migrate:latest \
  --region=us-central1 \
  --add-cloudsql-instances=leps-472702:us-central1:leps-postgres-stage \
  --service-account=lep-backend-sa@leps-472702.iam.gserviceaccount.com \
  --set-env-vars="ENVIRONMENT=stage,DB_USER=lep_user,DB_NAME=lep_database,INSTANCE_UNIX_SOCKET=/cloudsql/leps-472702:us-central1:leps-postgres-stage" \
  --set-secrets="DB_PASS=db-password-stage:latest" \
  --task-timeout=10m

# 3. Executar job
gcloud run jobs execute lep-migrate --region=us-central1

# 4. Ver logs
gcloud run jobs executions logs read --region=us-central1
```

### Verificar Schema do Banco

```bash
# Conectar ao banco via proxy
psql "host=127.0.0.1 port=5432 user=lep_user dbname=lep_database"

# Listar tabelas
\dt

# Ver estrutura de tabela específica
\d organizations

# Sair
\q
```

---

## Seeding de Dados

### Seed Local (Desenvolvimento)

```bash
# Executar seed completo
bash ./scripts/run_seed.sh

# Com verbose
bash ./scripts/run_seed.sh --verbose

# Limpar dados antes
bash ./scripts/run_seed.sh --clear-first
```

### Seed Remoto (Staging/Produção)

```bash
# Build do seed remoto
go build -o lep-seed-remote.exe cmd/seed-remote/main.go

# Executar para staging
./lep-seed-remote.exe --url https://lep-system-516622888070.us-central1.run.app --verbose

# Ou usar script
bash ./scripts/run_seed_remote.sh --verbose
```

**Credenciais após seeding**:
- **Master Admins**: `pablo@lep.com / senha123`
- **Demo Users**: `teste@gmail.com / password`

### Dados Criados

- 1 Organização
- 1 Projeto
- 6 Usuários (3 admins + 3 demo)
- 12 Produtos
- 8 Mesas
- 4 Pedidos
- 6 Reservas
- 3 Waitlist entries
- 5 Clientes
- 5 Templates de notificação

---

## Deploy da Aplicação

### Deploy para Staging

```bash
# Usar script pronto
bash ./scripts/stage-deploy.sh

# Ou manualmente
ENVIRONMENT=stage gcloud run deploy lep-system \
  --source . \
  --region=us-central1 \
  --platform=managed \
  --allow-unauthenticated \
  --memory=1Gi \
  --cpu=2 \
  --min-instances=0 \
  --max-instances=10 \
  --add-cloudsql-instances=$(terraform output -raw database_connection_name) \
  --service-account=$(terraform output -raw service_account_email) \
  --set-env-vars="ENVIRONMENT=stage,STORAGE_TYPE=gcs,BUCKET_NAME=$(terraform output -raw storage_bucket_name),BASE_URL=$(terraform output -raw storage_bucket_url),DB_USER=$(terraform output -raw database_user),DB_NAME=$(terraform output -raw database_name),INSTANCE_UNIX_SOCKET=/cloudsql/$(terraform output -raw database_connection_name)" \
  --set-secrets="DB_PASS=$(terraform output -raw db_password_secret_name):latest,JWT_SECRET_PRIVATE_KEY=$(terraform output -raw jwt_private_key_secret_name):latest,JWT_SECRET_PUBLIC_KEY=$(terraform output -raw jwt_public_key_secret_name):latest"
```

### Verificar Deploy

```bash
# Ver URL do serviço
gcloud run services describe lep-system --region=us-central1 --format="value(status.url)"

# Testar health check
curl https://lep-system-516622888070.us-central1.run.app/ping
# Esperado: "pong"

curl https://lep-system-516622888070.us-central1.run.app/health
# Esperado: {"status":"healthy"}
```

### Ver Logs

```bash
# Logs em tempo real
gcloud run services logs tail lep-system --region=us-central1

# Últimas 50 linhas
gcloud run services logs read lep-system --region=us-central1 --limit=50

# Filtrar por erro
gcloud run services logs read lep-system --region=us-central1 --limit=100 | grep ERROR
```

---

## Ambientes

### Estrutura de Ambientes

| Ambiente | Descrição | URL | Banco |
|----------|-----------|-----|-------|
| **dev** | Desenvolvimento local | http://localhost:8080 | Docker PostgreSQL |
| **stage** | Staging/homologação | https://lep-system-...run.app | Cloud SQL (dev) |
| **prod** | Produção (futuro) | TBD | Cloud SQL (prod) |

### Deploy por Ambiente

#### DEV (Local)

```bash
# Usar Docker Compose
bash ./scripts/dev-local.sh

# Ou menu interativo
bash ./scripts/master-interactive.sh
# Escolher: 1. Ambiente DEV (Local)
```

#### STAGE

```bash
# Via script
bash ./scripts/stage-deploy.sh

# Ou via Terraform + gcloud
terraform apply -var-file=environments/gcp-stage.tfvars
# ... seguir passos de deploy acima
```

#### PROD (Futuro)

```bash
# Criar novo arquivo environments/gcp-prd.tfvars
# Configurar valores de produção
terraform apply -var-file=environments/gcp-prd.tfvars
```

---

## Rollback

### Rollback do Cloud Run

```bash
# Listar revisões
gcloud run revisions list --service=lep-system --region=us-central1

# Fazer rollback para revisão anterior
gcloud run services update-traffic lep-system \
  --region=us-central1 \
  --to-revisions=lep-system-00002-abc=100
```

### Rollback do Banco de Dados

```bash
# Ver backups disponíveis
gcloud sql backups list --instance=leps-postgres-stage

# Restaurar backup
gcloud sql backups restore BACKUP_ID \
  --backup-instance=leps-postgres-stage \
  --backup-instance=leps-postgres-stage

# OU restaurar para nova instância (recomendado)
gcloud sql backups restore BACKUP_ID \
  --backup-instance=leps-postgres-stage \
  --target-instance=leps-postgres-stage-restore
```

### Rollback da Infraestrutura (Terraform)

```bash
# Ver histórico de state
terraform state list

# Importar state anterior (se tiver backup)
cp terraform.tfstate.backup terraform.tfstate

# Aplicar configuração anterior
terraform apply -var-file=environments/gcp-stage.tfvars
```

---

## Troubleshooting

### Problema: Cloud Run não inicia

**Sintomas**:
- Deploy bem-sucedido mas serviço não responde
- 503 Service Unavailable

**Diagnóstico**:
```bash
# Ver logs
gcloud run services logs read lep-system --region=us-central1 --limit=50

# Ver detalhes da revisão
gcloud run revisions describe lep-system-00001-abc --region=us-central1
```

**Soluções**:
- Verificar se Cloud SQL está acessível
- Verificar se secrets existem e estão acessíveis
- Verificar env vars (DB_USER, DB_NAME, etc.)

### Problema: Erro de conexão com Cloud SQL

**Sintomas**:
- "connection refused"
- "could not connect to server"

**Soluções**:
```bash
# Verificar que Cloud SQL instance está rodando
gcloud sql instances describe leps-postgres-stage

# Verificar connection name está correto
gcloud run services describe lep-system --region=us-central1 --format="value(spec.template.spec.containers[0].env[?(@.name=='INSTANCE_UNIX_SOCKET')].value)"

# Deve ser: /cloudsql/leps-472702:us-central1:leps-postgres-stage
```

### Problema: Seed remoto falha

**Sintomas**:
- "failed to create project: status 500"
- "failed to create user: E-mail já cadastrado"

**Soluções**:
- Atualizar backend com correções de duplicata
- Rebuild seed remoto: `go build -o lep-seed-remote.exe cmd/seed-remote/main.go`
- Ou limpar banco antes: conectar via proxy e `TRUNCATE TABLE projects CASCADE;`

### Problema: Terraform não encontra recursos

**Sintomas**:
- "Error: resource not found"
- Quer criar recursos que já existem

**Solução - Importar recursos**:
```bash
terraform import -var-file=environments/gcp-stage.tfvars \
  google_sql_database_instance.main leps-472702/leps-postgres-stage
```

### Problema: Permissões insuficientes

**Sintomas**:
- "permission denied"
- "403 Forbidden"

**Soluções**:
```bash
# Verificar IAM da service account
gcloud projects get-iam-policy leps-472702 \
  --flatten="bindings[].members" \
  --filter="bindings.members:lep-backend-sa@leps-472702.iam.gserviceaccount.com"

# Adicionar permissões manualmente
gcloud projects add-iam-policy-binding leps-472702 \
  --member="serviceAccount:lep-backend-sa@leps-472702.iam.gserviceaccount.com" \
  --role="roles/cloudsql.client"
```

---

## Checklist de Deploy

### Antes do Deploy

- [ ] Terraform configurado e autenticado
- [ ] JWT keys geradas e configuradas
- [ ] Arquivo `environments/gcp-{env}.tfvars` criado
- [ ] Backup do banco atual (se aplicável)

### Durante o Deploy

- [ ] `terraform apply` executado com sucesso
- [ ] Todos os recursos criados (Cloud SQL, GCS, Secrets)
- [ ] Migrations executadas
- [ ] Cloud Run deploy bem-sucedido
- [ ] Health check respondendo

### Após o Deploy

- [ ] Seed de dados executado
- [ ] Login funcional com credenciais de teste
- [ ] Endpoints principais testados
- [ ] Logs não mostram erros críticos
- [ ] Documentar versão deployada

---

**Última atualização**: 2025-10-14
**Versão**: 1.0
**Autor**: Claude Code + Equipe LEP
