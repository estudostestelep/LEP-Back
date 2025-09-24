# 🏗️ LEP System - Abordagem Híbrida Terraform + gcloud

*Data: 23/09/2024*

## 📋 Visão Geral

Este documento explica por que adotamos uma abordagem híbrida (Terraform + gcloud manual) para a infraestrutura do LEP System e como isso resolve problemas de permissões do GCP.

## ❌ Problemas Enfrentados com Terraform Puro

### **1. Permissões IAM Complexas**

Mesmo sendo `owner` do projeto, certas operações Terraform falhavam:

```bash
# Erros comuns encontrados
Error: Permission 'iam.serviceAccounts.create' denied on resource
Error: Permission 'secretmanager.secrets.create' denied for resource
Error: Permission 'artifactregistry.repositories.create' denied
Error: Request `List Project Services` returned error: Permission denied
```

### **2. Propagação de Permissões**

- Roles concedidas demoram até 5 minutos para propagar
- Terraform não aguarda propagação automaticamente
- Service account creation precede outras operações que dependem dele

### **3. Ordem de Dependências**

```hcl
# Terraform tentava criar resources antes das APIs estarem ativas
resource "google_project_service" "required_apis" {
  # APIs sendo habilitadas...
}

resource "google_service_account" "lep_backend_sa" {
  # Erro: ServiceAccounts API não propagou ainda
  depends_on = [google_project_service.required_apis]
}
```

### **4. Limitações do Provider Google**

- Algumas operações são melhor suportadas via gcloud CLI
- Provider Terraform pode estar desatualizado vs. gcloud latest
- Erros de timeout em operações longas (Cloud SQL creation)

## ✅ Solução: Abordagem Híbrida

### **Princípio da Divisão**

1. **gcloud manual**: Resources sensíveis, bootstrap, e operações que exigem permissões especiais
2. **Terraform**: Configuração, relações entre resources, e deploy de aplicação

### **Phase 1: Bootstrap Manual (gcloud)**

```bash
# 1. Habilitar APIs
gcloud services enable \
    secretmanager.googleapis.com \
    sqladmin.googleapis.com \
    run.googleapis.com \
    cloudbuild.googleapis.com \
    artifactregistry.googleapis.com

# 2. Criar Service Account
gcloud iam service-accounts create lep-backend-sa \
    --display-name="LEP Backend Service Account"

# 3. Criar Artifact Registry
gcloud artifacts repositories create lep-backend \
    --repository-format=docker \
    --location=us-central1

# 4. Criar Secrets
gcloud secrets create jwt-private-key-dev --replication-policy="automatic"
gcloud secrets create jwt-public-key-dev --replication-policy="automatic"
gcloud secrets create db-password-dev --replication-policy="automatic"

# 5. Criar Cloud SQL (operação longa)
gcloud sql instances create leps-postgres-dev \
    --database-version=POSTGRES_15 \
    --tier=db-f1-micro \
    --region=us-central1 \
    --assign-ip
```

### **Phase 2: Terraform com Data Sources**

```hcl
# Referenciar recursos existentes
data "google_service_account" "lep_backend_sa" {
  account_id = "lep-backend-sa"
  project    = var.project_id
}

data "google_artifact_registry_repository" "lep_repo" {
  repository_id = "lep-backend"
  location      = var.region
  project       = var.project_id
}

data "google_secret_manager_secret" "jwt_private_key" {
  secret_id = "jwt-private-key-dev"
  project   = var.project_id
}

# Criar resources que dependem dos existentes
resource "google_secret_manager_secret_version" "jwt_private_key_version" {
  secret      = data.google_secret_manager_secret.jwt_private_key.id
  secret_data = var.jwt_private_key
}

resource "google_cloud_run_v2_service" "lep_backend" {
  name         = "${var.project_name}-backend-${var.environment}"
  location     = var.region

  template {
    service_account = data.google_service_account.lep_backend_sa.email
    # ... resto da configuração
  }
}
```

## 📁 Estrutura de Arquivos

### **Arquivos Terraform**

```
LEP-Back/
├── main.tf                 # ❌ Versão original (problemas de permissão)
├── main_simplified.tf      # ✅ Versão híbrida com data sources
├── variables.tf            # Variáveis compartilhadas
├── outputs.tf              # Outputs necessários
├── terraform.tfvars        # Valores para variáveis
└── environments/
    ├── gcp-dev.tfvars     # Config ambiente dev
    ├── gcp-stage.tfvars   # Config ambiente staging
    └── gcp-prd.tfvars     # Config ambiente produção
```

### **Scripts de Automação**

```
scripts/
├── bootstrap-gcp.sh        # ⚡ Cria resources via gcloud (Phase 1)
├── deploy-terraform.sh     # 🏗️ Executa terraform (Phase 2)
├── deploy-interactive.sh   # 🎯 Script completo híbrido
└── cleanup-resources.sh    # 🧹 Limpeza para recomeçar
```

## 🔧 Implementação Prática

### **1. Script de Bootstrap**

```bash
#!/bin/bash
# bootstrap-gcp.sh - Cria resources base via gcloud

PROJECT_ID="leps-472702"
REGION="us-central1"

echo "🚀 Bootstrapping GCP resources for LEP System..."

# Enable APIs
echo "📡 Enabling APIs..."
gcloud services enable \
    secretmanager.googleapis.com \
    sqladmin.googleapis.com \
    run.googleapis.com \
    cloudbuild.googleapis.com \
    artifactregistry.googleapis.com \
    --project=$PROJECT_ID

# Create Service Account
echo "👤 Creating Service Account..."
gcloud iam service-accounts create lep-backend-sa \
    --display-name="LEP Backend Service Account" \
    --description="Service account for LEP Backend Cloud Run service" \
    --project=$PROJECT_ID

# Create Artifact Registry
echo "📦 Creating Artifact Registry..."
gcloud artifacts repositories create lep-backend \
    --repository-format=docker \
    --location=$REGION \
    --description="LEP Backend Docker repository" \
    --project=$PROJECT_ID

# Create Secrets
echo "🔐 Creating Secret Manager secrets..."
gcloud secrets create jwt-private-key-dev \
    --replication-policy="automatic" \
    --project=$PROJECT_ID

gcloud secrets create jwt-public-key-dev \
    --replication-policy="automatic" \
    --project=$PROJECT_ID

gcloud secrets create db-password-dev \
    --replication-policy="automatic" \
    --project=$PROJECT_ID

echo "✅ Bootstrap completed! Now run Terraform..."
```

### **2. Terraform Simplificado**

```hcl
# main_simplified.tf
terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.26.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# Import existing resources
data "google_service_account" "lep_backend_sa" {
  account_id = "lep-backend-sa"
  project    = var.project_id
}

data "google_artifact_registry_repository" "lep_repo" {
  repository_id = "lep-backend"
  location      = var.region
  project       = var.project_id
}

# Create only what Terraform manages well
resource "google_cloud_run_v2_service" "lep_backend" {
  name         = "${var.project_name}-backend-${var.environment}"
  location     = var.region

  template {
    service_account = data.google_service_account.lep_backend_sa.email

    containers {
      image = "${var.region}-docker.pkg.dev/${var.project_id}/${data.google_artifact_registry_repository.lep_repo.repository_id}/lep-backend:latest"
      # ... configuração do container
    }
  }
}
```

### **3. Deploy Integrado**

```bash
#!/bin/bash
# deploy-interactive.sh (versão simplificada)

echo "🎯 LEP System Hybrid Deployment"

# Phase 1: Bootstrap (se necessário)
if ! gcloud iam service-accounts describe lep-backend-sa@$PROJECT_ID.iam.gserviceaccount.com &>/dev/null; then
    echo "🚀 Running bootstrap..."
    ./scripts/bootstrap-gcp.sh
fi

# Phase 2: Terraform
echo "🏗️ Running Terraform..."
terraform init
terraform plan -var-file=environments/$ENVIRONMENT.tfvars -out=tfplan
terraform apply tfplan

# Phase 3: Build and Deploy
echo "🐳 Building and deploying application..."
docker build -t $REGION-docker.pkg.dev/$PROJECT_ID/lep-backend/lep-backend:latest .
docker push $REGION-docker.pkg.dev/$PROJECT_ID/lep-backend/lep-backend:latest

echo "✅ Deployment completed!"
```

## 📊 Comparação: Antes vs. Depois

### **Terraform Puro (❌ Problemático)**

| Aspecto | Status | Problemas |
|---------|--------|-----------|
| APIs | ❌ | Permission denied |
| Service Account | ❌ | IAM create denied |
| Secrets | ❌ | Secret Manager denied |
| Cloud SQL | ❌ | Timeout issues |
| Cloud Run | ❌ | Dependências quebradas |
| **Deploy Time** | **N/A** | **Falha** |

### **Híbrido (✅ Funcional)**

| Aspecto | Método | Status | Tempo |
|---------|---------|--------|-------|
| APIs | gcloud | ✅ | 2min |
| Service Account | gcloud | ✅ | 30s |
| Secrets | gcloud | ✅ | 1min |
| Cloud SQL | gcloud | ✅ | 10-15min |
| Cloud Run | Terraform | ✅ | 2min |
| **Deploy Time** | **Híbrido** | **✅** | **15-20min** |

## 🎯 Vantagens da Abordagem Híbrida

### **1. Confiabilidade**
- ✅ Bootstrap sempre funciona (gcloud é mais estável)
- ✅ Terraform só gerencia o que faz melhor
- ✅ Rollback mais previsível

### **2. Flexibilidade**
- ✅ Pode usar Terraform puro no futuro (quando permissões estiverem corretas)
- ✅ gcloud para operações de emergência
- ✅ Diferentes métodos para diferentes ambientes

### **3. Manutenibilidade**
- ✅ Scripts claros e específicos
- ✅ Separação de responsabilidades
- ✅ Debug mais fácil (falha em bootstrap vs. falha em terraform)

### **4. Velocidade**
- ✅ Bootstrap paralelo (APIs + Service Accounts)
- ✅ Terraform só para Cloud Run (rápido)
- ✅ Menos retry logic necessário

## 🔮 Migração Futura para Terraform Puro

Quando as permissões estiverem estáveis, podemos migrar:

### **1. Import Resources**
```bash
# Importar resources criados manualmente
terraform import google_service_account.lep_backend_sa projects/$PROJECT_ID/serviceAccounts/lep-backend-sa@$PROJECT_ID.iam.gserviceaccount.com

terraform import google_artifact_registry_repository.lep_repo projects/$PROJECT_ID/locations/$REGION/repositories/lep-backend
```

### **2. Gradual Migration**
- Migrar secrets primeiro
- Depois artifact registry
- Por último service accounts (mais crítico)

### **3. Validation**
```bash
# Verificar state após import
terraform plan  # Deve mostrar "No changes"
```

## 🚨 Troubleshooting da Abordagem Híbrida

### **Problema: Resource já existe**
```bash
# Erro gcloud
ERROR: The repository [lep-backend] already exists

# Solução: Verificar antes de criar
if ! gcloud artifacts repositories describe lep-backend --location=$REGION &>/dev/null; then
    gcloud artifacts repositories create lep-backend
fi
```

### **Problema: Terraform não encontra data source**
```bash
# Erro terraform
Error: Error reading service account: Service account not found

# Solução: Verificar naming e projeto
data "google_service_account" "lep_backend_sa" {
  account_id = "lep-backend-sa"  # ✅ Sem @project.iam
  project    = var.project_id    # ✅ Projeto explícito
}
```

### **Problema: Permissões para Terraform**
```bash
# Terraform ainda precisa de algumas permissões
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="user:$(gcloud config get-value account)" \
    --role="roles/run.admin"
```

## 📚 Recursos e Referências

### **Documentação Oficial**
- [Terraform Google Provider](https://registry.terraform.io/providers/hashicorp/google/latest/docs)
- [GCP IAM Best Practices](https://cloud.google.com/iam/docs/using-iam-securely)
- [gcloud CLI Reference](https://cloud.google.com/sdk/gcloud/reference)

### **Arquivos no Projeto**
- `main_simplified.tf` - Implementação híbrida
- `scripts/deploy-interactive.sh` - Deploy automatizado
- `DEPLOYMENT_TROUBLESHOOTING.md` - Troubleshooting detalhado

### **Commands Quick Reference**
```bash
# Bootstrap manual
./scripts/bootstrap-gcp.sh

# Deploy híbrido
./scripts/deploy-interactive.sh

# Health check
./scripts/health-check.sh --environment gcp-dev

# Cleanup (para recomeçar)
./scripts/cleanup-resources.sh
```

---

*A abordagem híbrida fornece o melhor dos dois mundos: confiabilidade do gcloud CLI com a flexibilidade do Terraform para configuração de aplicação.*