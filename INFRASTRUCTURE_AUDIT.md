# LEP System - Infrastructure Audit Report
**Data**: 2025-10-14
**Objetivo**: Validar alinhamento entre Models, Terraform, Database e Scripts

---

## 📊 Executive Summary

### Status Geral
- ✅ **Models do PostgreSQL**: 30 entidades bem definidas
- ⚠️ **Terraform**: Configuração minimalista (apenas senha do banco)
- ✅ **Scripts de Deploy**: Completos e funcionais
- ⚠️ **Infraestrutura GCP**: Parcialmente gerenciada via Terraform

### Principais Achados
1. **Infraestrutura existente não está no Terraform** - Cloud SQL, buckets e service accounts foram criados manualmente
2. **Models estão completos** - 30 entidades com todos os campos necessários
3. **Auto-migration do GORM** - O sistema usa GORM AutoMigrate para criar tabelas
4. **Seed remoto criado** - Nova ferramenta para popular banco via HTTP

---

## 🗂️ Database Models (PostgresLEP.go)

### Entidades Core (10)
1. **Organization** - Organizações (multi-tenant principal)
2. **Project** - Projetos dentro de organizações
3. **User** - Usuários/funcionários
4. **UserOrganization** - Relacionamento N:N User-Organization
5. **UserProject** - Relacionamento N:N User-Project
6. **Customer** - Clientes do restaurante
7. **Table** - Mesas do restaurante
8. **Product** - Produtos do cardápio (refatorado)
9. **Order** - Pedidos
10. **Reservation** - Reservas de mesa

### Entidades de Suporte (8)
11. **Waitlist** - Lista de espera
12. **Settings** - Configurações paramet rizáveis
13. **Environment** - Ambientes físicos do restaurante
14. **AuditLog** - Log de auditoria
15. **BannedLists** - Tokens banidos (JWT blacklist)
16. **LoggedLists** - Tokens ativos (JWT whitelist)
17. **Tag** - Tags/categorias genéricas
18. **ProductTag** - Relacionamento Product-Tag

### Sistema de Cardápio (4)
19. **Menu** - Cardápios principais
20. **Category** - Categorias do cardápio
21. **Subcategory** - Subcategorias
22. **SubcategoryCategory** - Relacionamento N:N

### Sistema de Notificações (5)
23. **NotificationConfig** - Configuração de eventos
24. **NotificationTemplate** - Templates de mensagens
25. **NotificationLog** - Log de envios
26. **NotificationEvent** - Fila de eventos
27. **NotificationInbound** - Mensagens recebidas (2-way)

### Features Avançadas (3)
28. **BlockedPeriod** - Períodos bloqueados para reservas
29. **Lead** - Sistema básico de CRM
30. **ReportMetric** - Métricas para relatórios

### Campos Especiais Importantes
- **Product**: Refatorado com suporte para Pratos, Bebidas e Vinhos
  - Campos específicos de vinho: `Vintage`, `Country`, `Region`, `Winery`, `WineType`, `Grapes[]`
  - Preços múltiplos: `PriceBottle`, `PriceHalfBottle`, `PriceGlass`
  - Tipo e organização: `Type`, `Order`, `Active`, `PDVCode`
  - Relacionamentos: `CategoryId`, `SubcategoryId`

- **Order**: Sistema completo de tracking
  - `OrderItems` (JSONB) - Array de itens com suporte a objeto único (backward compatibility)
  - Timestamps de workflow: `StartedAt`, `ReadyAt`, `DeliveredAt`
  - Tempo estimado: `EstimatedPrepTime`, `EstimatedDeliveryTime`

---

## 🏗️ Terraform Configuration

### Arquivos Identificados
1. **main.tf** - Configuração principal (MUITO BÁSICO)
2. **variables.tf** - Variáveis completas e validações
3. **environments/gcp-stage.tfvars** - Config para staging
4. **environments/gcp-prd.tfvars** - Config para produção

### Estado Atual do main.tf
```terraform
# Apenas 1 recurso definido:
- random_password.db_password

# Comentado/não gerenciado:
- Cloud Run service (gerenciado via gcloud)
- Cloud SQL instance (não definido)
- GCS buckets (não definido)
- Service Account (não definido)
- Secrets Manager (não definido)
```

### ⚠️ Problema Identificado
**A infraestrutura GCP existente NÃO está sendo gerenciada pelo Terraform!**

#### Recursos GCP Existentes (criados manualmente):
- ✅ Cloud SQL: `leps-postgres-dev`
- ✅ Bucket GCS: `leps-472702-lep-images-dev`
- ✅ Bucket GCS: `leps-472702-lep-images-stage`
- ✅ Service Account: `lep-backend-sa@leps-472702.iam.gserviceaccount.com`
- ✅ Secrets Manager: `db-password-dev`, `jwt-private-key-dev`, `jwt-public-key-dev`
- ✅ Cloud Run Service: `lep-system` (https://lep-system-516622888070.us-central1.run.app)

#### Recursos NÃO definidos no Terraform:
- ❌ Cloud SQL PostgreSQL instance
- ❌ Google Cloud Storage buckets
- ❌ Service Account + IAM bindings
- ❌ Secret Manager secrets
- ❌ API enablement
- ❌ Network configuration

---

## 📜 Scripts de Deployment

### Scripts Principais
1. **`scripts/stage-deploy.sh`** ✅
   - Deploy via `gcloud run deploy`
   - Configura env vars e secrets
   - Conecta Cloud SQL
   - **Status**: Funcional e completo

2. **`scripts/master-interactive.sh`** ✅
   - Menu interativo unificado
   - Suporta DEV (local) e STAGE (GCP)
   - Bootstrap de infraestrutura (linha 332-365)
   - **Status**: Completo

3. **`scripts/run_seed_remote.sh`** ✅ (NOVO)
   - Seed via HTTP API
   - Não requer acesso direto ao PostgreSQL
   - **Status**: Criado nesta sessão

### Fluxo de Bootstrap (master-interactive.sh)
```bash
stage_bootstrap_infrastructure() {
    # Linha 347: Executa terraform init
    # Linha 360: Executa terraform apply -var-file=environments/gcp-stage.tfvars
    # ⚠️ PROBLEMA: main.tf não tem recursos para criar!
}
```

---

## 🔄 Auto-Migration (GORM)

### Como Funciona
O sistema usa **GORM AutoMigrate** para criar e atualizar tabelas automaticamente:

```go
// Localização: cmd/seed/main.go (linha 100-126)
func runMigrations(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.Organization{},
        &models.Project{},
        &models.User{},
        // ... todas as 30 entidades
    )
}
```

### ✅ Vantagens
- Schema sempre alinhado com os models
- Não requer migrations manuais
- Cria/atualiza tabelas automaticamente

### ⚠️ Desvantagens
- Não remove colunas antigas
- Não altera tipos incompatíveis
- Não gerencia índices complexos

### 📍 Quando é Executado
1. **Seed local**: `go run cmd/seed/main.go` (linha 59-63)
2. **Seed remote**: Apenas via endpoints HTTP (não executa migrations)

---

## 🌱 Sistema de Seeding

### Seeding Local (cmd/seed/main.go)
```bash
# Execução
./scripts/run_seed.sh
# OU
go run cmd/seed/main.go

# Características:
- Conecta diretamente ao PostgreSQL
- Executa AutoMigrate ANTES de inserir dados
- Usa httptest para testar endpoints localmente
- Suporta --clear-first para limpar dados
```

### Seeding Remoto (cmd/seed-remote/main.go) 🆕
```bash
# Execução
./lep-seed-remote.exe --url https://lep-system-516622888070.us-central1.run.app

# Características:
- Faz requisições HTTP reais para a API
- NÃO executa migrations (assume que tabelas existem)
- Trata duplicatas graciosamente
- Não requer credenciais do banco
```

### Dados Criados pelo Seed
- 1 Organização: "LEP Demo Organization"
- 1 Projeto: "LEP Demo Project"
- 6 Usuários (3 master admins + 3 demo users)
- 12 Produtos em 3 categorias
- 8 Mesas com status variados
- 4 Pedidos ativos
- 6 Reservas (passado, presente, futuro)
- 3 Entradas na lista de espera
- 5 Clientes
- 5 Templates de notificação

---

## 🚨 Problemas Identificados

### 1. Erro no Seed Remoto - Projeto já existe (500)
**Sintoma**:
```
failed to create project: status 500 - {"error":"Error creating project"}
```

**Causa Provável**:
- Banco em staging já tem um projeto com mesmo ID
- Endpoint `/project` retorna 500 em vez de 409 para duplicatas
- Falta tratamento de erro no backend

**Solução Recomendada**:
1. Verificar handler de criação de projeto
2. Retornar 409 Conflict para duplicatas (como nas outras entidades)
3. Ou: Limpar banco staging antes de rodar seed

### 2. Infraestrutura não gerenciada pelo Terraform
**Problema**:
- Resources GCP existem mas não estão no Terraform
- `terraform apply` não faz nada útil
- Drift entre estado real e código

**Impacto**:
- Dificulta replicar ambiente
- Sem controle de versão da infra
- Risco de mudanças manuais não documentadas

**Solução Recomendada**:
Opção A: **Import existente para Terraform**
```bash
terraform import google_sql_database_instance.main leps-472702/leps-postgres-dev
terraform import google_storage_bucket.images leps-472702-lep-images-dev
# ... etc
```

Opção B: **Criar novo main.tf completo** (recomendado)
```terraform
# Ver seção "Terraform Completo Recomendado" abaixo
```

### 3. AutoMigrate não executa em produção
**Problema**:
- Seed remoto não executa migrations
- Tabelas precisam existir antes do seed
- Sem processo claro de migration em prod

**Solução**:
1. Executar seed local UMA VEZ para criar schema
2. OU: Criar job de migration separado que apenas roda AutoMigrate
3. OU: Usar migrations manuais (golang-migrate)

---

## ✅ Recomendações

### Prioridade ALTA

#### 1. Corrigir Handler de Projeto
**Arquivo**: `server/project.go` ou `handler/project.go`

Garantir que retorna 409 para duplicatas:
```go
// Se projeto já existe com mesmo ID/nome
if err == gorm.ErrDuplicatedKey || strings.Contains(err.Error(), "duplicate key") {
    utils.SendConflictError(c, "Project")
    return
}
```

#### 2. Importar Infraestrutura Existente para Terraform

**Criar arquivo**: `main.tf` completo

```terraform
# Cloud SQL Instance
resource "google_sql_database_instance" "main" {
  name             = "${var.project_name}-postgres-${var.environment}"
  database_version = "POSTGRES_15"
  region           = var.region
  project          = var.project_id

  settings {
    tier              = var.db_tier
    availability_type = var.db_availability_type
    disk_size         = var.db_disk_size

    backup_configuration {
      enabled    = true
      start_time = "03:00"
    }

    ip_configuration {
      ipv4_enabled = false
      require_ssl  = false
    }

    database_flags {
      name  = "max_connections"
      value = "100"
    }
  }

  deletion_protection = var.enable_deletion_protection
}

resource "google_sql_database" "main" {
  name     = var.database_name
  instance = google_sql_database_instance.main.name
  project  = var.project_id
}

resource "google_sql_user" "main" {
  name     = var.database_user
  instance = google_sql_database_instance.main.name
  password = random_password.db_password.result
  project  = var.project_id
}

# GCS Buckets
resource "google_storage_bucket" "images" {
  name          = "${var.project_id}-${var.project_name}-images-${var.environment}"
  location      = var.region
  project       = var.project_id
  force_destroy = var.environment != "prod"

  uniform_bucket_level_access = true

  cors {
    origin          = ["*"]
    method          = ["GET", "POST", "PUT", "DELETE"]
    response_header = ["*"]
    max_age_seconds = 3600
  }
}

resource "google_storage_bucket_iam_member" "images_public_read" {
  bucket = google_storage_bucket.images.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

# Service Account
resource "google_service_account" "backend" {
  account_id   = "${var.project_name}-backend-sa"
  display_name = "LEP Backend Service Account"
  project      = var.project_id
}

resource "google_project_iam_member" "backend_sql_client" {
  project = var.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.backend.email}"
}

resource "google_project_iam_member" "backend_storage_admin" {
  project = var.project_id
  role    = "roles/storage.admin"
  member  = "serviceAccount:${google_service_account.backend.email}"
}

resource "google_project_iam_member" "backend_secret_accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.backend.email}"
}

# Secret Manager
resource "google_secret_manager_secret" "db_password" {
  secret_id = "db-password-${var.environment}"
  project   = var.project_id

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "db_password" {
  secret      = google_secret_manager_secret.db_password.id
  secret_data = random_password.db_password.result
}

resource "google_secret_manager_secret" "jwt_private_key" {
  secret_id = "jwt-private-key-${var.environment}"
  project   = var.project_id

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "jwt_private_key" {
  secret      = google_secret_manager_secret.jwt_private_key.id
  secret_data = var.jwt_private_key
}

resource "google_secret_manager_secret" "jwt_public_key" {
  secret_id = "jwt-public-key-${var.environment}"
  project   = var.project_id

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "jwt_public_key" {
  secret      = google_secret_manager_secret.jwt_public_key.id
  secret_data = var.jwt_public_key
}

# Enable Required APIs
resource "google_project_service" "services" {
  for_each = toset([
    "run.googleapis.com",
    "sql-component.googleapis.com",
    "sqladmin.googleapis.com",
    "storage.googleapis.com",
    "secretmanager.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "iam.googleapis.com",
    "cloudbuild.googleapis.com"
  ])

  project = var.project_id
  service = each.key

  disable_on_destroy = false
}
```

**Depois de criar, importar recursos existentes**:
```bash
# 1. Inicializar
terraform init

# 2. Importar recursos existentes
terraform import -var-file=environments/gcp-stage.tfvars \
  google_sql_database_instance.main leps-472702/leps-postgres-dev

terraform import -var-file=environments/gcp-stage.tfvars \
  google_storage_bucket.images leps-472702-lep-images-stage

terraform import -var-file=environments/gcp-stage.tfvars \
  google_service_account.backend projects/leps-472702/serviceAccounts/lep-backend-sa@leps-472702.iam.gserviceaccount.com

# 3. Verificar plano (não deve ter mudanças)
terraform plan -var-file=environments/gcp-stage.tfvars
```

#### 3. Criar Job de Migration Dedicado

**Novo arquivo**: `cmd/migrate/main.go`

```go
package main

import (
	"fmt"
	"lep/config"
	"lep/repositories/models"
	"lep/resource"
	"log"

	"gorm.io/gorm"
)

func main() {
	fmt.Println("🔄 LEP Database Migration Tool")
	fmt.Println("================================")

	// Connect to database
	db, err := resource.OpenConnDBPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run auto-migration
	fmt.Println("📦 Running auto-migration...")
	err = runMigrations(db)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("✅ Migration completed successfully!")
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Organization{},
		&models.Project{},
		&models.User{},
		&models.UserOrganization{},
		&models.UserProject{},
		&models.Customer{},
		&models.Menu{},
		&models.Category{},
		&models.Subcategory{},
		&models.SubcategoryCategory{},
		&models.Tag{},
		&models.Product{},
		&models.ProductTag{},
		&models.Table{},
		&models.Order{},
		&models.Reservation{},
		&models.Waitlist{},
		&models.Environment{},
		&models.Settings{},
		&models.NotificationTemplate{},
		&models.NotificationConfig{},
		&models.NotificationLog{},
		&models.NotificationEvent{},
		&models.NotificationInbound{},
		&models.BlockedPeriod{},
		&models.Lead{},
		&models.ReportMetric{},
		&models.BannedLists{},
		&models.LoggedLists{},
		&models.AuditLog{},
	)
}
```

**Uso**:
```bash
# Build
go build -o lep-migrate.exe cmd/migrate/main.go

# Executar localmente
ENVIRONMENT=dev ./lep-migrate.exe

# Executar em staging via Cloud SQL Proxy
ENVIRONMENT=stage ./lep-migrate.exe

# OU via Cloud Run job (recomendado para prod)
gcloud run jobs create lep-migrate \
  --image=gcr.io/leps-472702/lep-migrate:latest \
  --add-cloudsql-instances=leps-472702:us-central1:leps-postgres-dev \
  --set-env-vars="ENVIRONMENT=stage" \
  --set-secrets="DB_PASS=db-password-dev:latest"
```

### Prioridade MÉDIA

#### 4. Adicionar Validação no Seed Remoto
```go
// Antes de começar seed, verificar que as tabelas existem
func validateTablesExist(baseURL string) error {
    // Fazer GET /health ou criar endpoint /schema/validate
    // Retornar erro se schema não está pronto
}
```

#### 5. Documentar Processo de Bootstrap
Criar `DEPLOYMENT.md` com:
- [ ] Passos para criar ambiente do zero
- [ ] Ordem de execução (Terraform → Migration → Seed)
- [ ] Comandos exatos para cada ambiente
- [ ] Rollback procedures

### Prioridade BAIXA

#### 6. Considerar Migrations Manuais (golang-migrate)
Se AutoMigrate se tornar limitante:
```bash
# Instalar
go get -u github.com/golang-migrate/migrate

# Criar migrations
migrate create -ext sql -dir migrations -seq add_product_wine_fields

# Executar
migrate -path migrations -database "postgresql://..." up
```

---

## 📋 Checklist de Ações Imediatas

### Para Resolver Erro do Seed Remoto

- [ ] 1. Verificar handler de projeto em `server/project.go`
- [ ] 2. Adicionar tratamento de duplicatas (retornar 409)
- [ ] 3. Rebuild e redeploy do backend
- [ ] 4. Testar seed remoto novamente

### Para Alinhar Infraestrutura

- [ ] 1. Criar `main.tf` completo com todos os recursos
- [ ] 2. Executar `terraform import` para cada recurso existente
- [ ] 3. Validar com `terraform plan` (deve mostrar "no changes")
- [ ] 4. Commitar `main.tf` atualizado
- [ ] 5. Documentar processo no README.md

### Para Garantir Schema Consistente

- [ ] 1. Criar `cmd/migrate/main.go`
- [ ] 2. Executar migration em staging UMA VEZ
- [ ] 3. Verificar que todas as tabelas existem corretamente
- [ ] 4. Executar seed remoto
- [ ] 5. Validar dados no banco

---

## 🎯 Próximos Passos Recomendados

### Opção A: Resolver Apenas o Seed (Rápido - 30min)
1. Corrigir handler de projeto (5min)
2. Rebuild e deploy (10min)
3. Executar seed remoto (5min)
4. Validar dados (10min)

### Opção B: Alinhar Tudo (Completo - 2-3h)
1. Criar `main.tf` completo (1h)
2. Importar recursos existentes (30min)
3. Criar tool de migration (30min)
4. Executar migration em staging (10min)
5. Corrigir handler de projeto (10min)
6. Deploy e seed (20min)
7. Documentar processo (30min)

### Opção C: Começar do Zero (Mais Seguro - 1 dia)
1. Criar novo ambiente staging-v2 com Terraform
2. Aplicar migrations
3. Seed com dados
4. Testar completamente
5. Migrar tráfego de staging antigo
6. Deletar staging antigo

---

## 📊 Resumo Técnico

| Componente | Status | Observações |
|------------|--------|-------------|
| **Models (30)** | ✅ Completo | Todas as entidades bem definidas |
| **GORM AutoMigrate** | ✅ Funcional | Usado no seed local |
| **Terraform** | ⚠️ Incompleto | Apenas senha do banco |
| **GCP Infra** | ✅ Existente | Criado manualmente |
| **Seed Local** | ✅ Funcional | Com migrations |
| **Seed Remoto** | ⚠️ Erro | Projeto duplicado (500) |
| **Deploy Scripts** | ✅ Completo | `stage-deploy.sh` OK |
| **Master Script** | ✅ Completo | Menu interativo |

---

## 🔗 Arquivos Relevantes

- **Models**: `repositories/models/PostgresLEP.go`
- **Terraform**: `main.tf`, `variables.tf`, `environments/gcp-stage.tfvars`
- **Seed Local**: `cmd/seed/main.go`, `utils/seed_data.go`
- **Seed Remoto**: `cmd/seed-remote/main.go`
- **Deploy**: `scripts/stage-deploy.sh`, `scripts/master-interactive.sh`
- **Handlers**: `handler/project.go`, `server/project.go`

---

**Gerado por**: Claude Code
**Data**: 2025-10-14
**Versão**: 1.0
