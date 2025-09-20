# 🚀 LEP System - Guia de Deploy Interativo

## 📋 Resumo

Este guia documenta o sistema de deploy interativo multi-ambiente do LEP System, que suporta 4 ambientes distintos com validação completa de dependências e variáveis.

## 🎯 Ambientes Suportados

### 1. 🏠 **local-dev** (Desenvolvimento Local)
- **Descrição**: Ambiente completo local usando Docker Compose
- **Dependências**: Docker, Docker Compose, Go
- **Infraestrutura**: PostgreSQL, Redis, MailHog, PgAdmin
- **Características**:
  - Zero custo
  - Hot reloading
  - Interface de email local
  - Banco de dados visual

### 2. ☁️ **gcp-dev** (GCP Mínimo)
- **Descrição**: Setup mínimo no GCP para testes rápidos
- **Dependências**: gcloud, terraform, docker, go
- **Infraestrutura**: Cloud SQL (f1-micro), Cloud Run (mínimo)
- **Características**:
  - Recursos mínimos (custo baixo)
  - Scale-to-zero
  - Sem notificações
  - Sem alta disponibilidade

### 3. 🚀 **gcp-stage** (GCP Staging)
- **Descrição**: Ambiente similar à produção, mas sem Twilio
- **Dependências**: gcloud, terraform, docker, go + SMTP
- **Infraestrutura**: Cloud SQL (regional), Cloud Run (production-like)
- **Características**:
  - Alta disponibilidade
  - SMTP notifications apenas
  - Recursos de produção
  - Deletion protection

### 4. 🌟 **gcp-prd** (GCP Produção)
- **Descrição**: Ambiente de produção completo
- **Dependências**: gcloud, terraform, docker, go + Twilio + SMTP
- **Infraestrutura**: Cloud SQL (high-performance), Cloud Run (full capacity)
- **Características**:
  - Todas as funcionalidades
  - Twilio + SMTP
  - Custom domain
  - Monitoring completo

## 🛠️ Scripts de Deploy

### **Linux/Mac: `./scripts/deploy-interactive.sh`**
```bash
# Deploy interativo
./scripts/deploy-interactive.sh

# Deploy direto para um ambiente específico
ENVIRONMENT=local-dev ./scripts/deploy-interactive.sh
```

### **Windows: `./scripts/Deploy-Interactive.ps1`**
```powershell
# Deploy interativo
.\scripts\Deploy-Interactive.ps1

# Deploy direto para um ambiente específico
.\scripts\Deploy-Interactive.ps1 -Environment "local-dev"
```

## 📁 Estrutura de Configuração

```
LEP-Back/
├── environments/                    # Configurações por ambiente
│   ├── local-dev.env               # Variáveis locais
│   ├── gcp-dev.tfvars             # Terraform dev
│   ├── gcp-stage.tfvars           # Terraform staging
│   └── gcp-prd.tfvars             # Terraform produção
├── scripts/
│   ├── deploy-interactive.sh       # Script principal (Bash)
│   ├── Deploy-Interactive.ps1      # Script principal (PowerShell)
│   └── init-db.sql                # Inicialização do banco local
├── docker-compose.yml              # Ambiente local completo
├── Dockerfile                      # Production-ready container
├── Dockerfile.dev                  # Development container
└── main.go                         # Aplicação com detecção de ambiente
```

## 🔧 Pré-requisitos por Ambiente

### **local-dev**
- ✅ Docker & Docker Compose
- ✅ Go 1.21+
- ✅ Arquivo `environments/local-dev.env`

### **gcp-dev**
- ✅ Google Cloud CLI
- ✅ Terraform
- ✅ Docker
- ✅ Go 1.21+
- ✅ Arquivo `environments/gcp-dev.tfvars`
- ✅ JWT keys (`.pem` files)
- ✅ GCP Project: `leps-472702`

### **gcp-stage**
- ✅ Todos os requisitos do gcp-dev
- ✅ SMTP credentials configuradas
- ✅ Arquivo `environments/gcp-stage.tfvars`

### **gcp-prd**
- ✅ Todos os requisitos do gcp-stage
- ✅ Twilio credentials configuradas
- ✅ Domain configurado (opcional)
- ✅ Arquivo `environments/gcp-prd.tfvars`

## 🚦 Fluxo de Validação

O script executa validações automáticas antes do deploy:

### 1. **Verificação de Dependências**
```bash
✅ Docker instalado e funcionando
✅ Docker Compose disponível
✅ Google Cloud CLI autenticado
✅ Terraform instalado
✅ Go environment configurado
```

### 2. **Validação de Variáveis Obrigatórias**
```bash
✅ Arquivos de configuração existem
✅ JWT keys estão presentes
✅ SMTP configurado (stage/prod)
✅ Twilio configurado (prod)
✅ Project ID correto
```

### 3. **Verificação de Acesso**
```bash
✅ GCP authentication ativa
✅ Permissões no projeto
✅ Docker registry access
✅ Terraform state access
```

## ⚡ Execução e Monitoramento

### **Durante o Deploy**
- 📊 **Progress tracking**: Mostra etapa atual (X/Y)
- 🎨 **Logs coloridos**: INFO, WARN, ERROR, SUCCESS
- 📝 **Comandos visíveis**: Todos os comandos executados são mostrados
- ⏱️ **Error handling**: Em caso de falha, mostra comandos restantes

### **Exemplo de Output**
```bash
[STEP 3/6] Building Docker image
[CMD] docker build -t us-central1-docker.pkg.dev/leps-472702/lep-backend/lep-backend:latest .
[SUCCESS] Docker image build completed

[STEP 4/6] Pushing Docker image
[CMD] docker push us-central1-docker.pkg.dev/leps-472702/lep-backend/lep-backend:latest
[SUCCESS] Docker image push completed
```

### **Em Caso de Erro**
```bash
[ERROR] Infrastructure deployment failed with exit code 1
[ERROR] Failed command: terraform apply tfplan

Remaining commands to complete deployment:
  - docker build
  - docker push
  - gcloud run deploy
```

## 🏥 Health Checks e Verificação

### **Endpoints Automáticos**
- 🔍 `/health` - Status detalhado da aplicação
- 🏓 `/ping` - Verificação básica de conectividade

### **Informações de Health Check**
```json
{
  "status": "healthy",
  "environment": "gcp-stage",
  "version": "1.0.0",
  "mode": "gcp",
  "platform": "cloud-run"
}
```

## 📊 Recursos por Ambiente

| Recurso | local-dev | gcp-dev | gcp-stage | gcp-prd |
|---------|-----------|---------|-----------|---------|
| **CPU** | Host | 0.5 vCPU | 2 vCPU | 4 vCPU |
| **Memory** | Host | 256Mi | 1Gi | 2Gi |
| **Database** | PostgreSQL local | f1-micro | n1-standard-1 | n1-standard-2 |
| **Availability** | Single | ZONAL | REGIONAL | REGIONAL |
| **Instances** | 1 | 0-3 | 1-20 | 2-50 |
| **Storage** | Local | 10GB | 50GB | 100GB |
| **Notifications** | ❌ | ❌ | SMTP only | Twilio + SMTP |

## 🔒 Segurança

### **Secrets Management**
- 🔐 JWT keys via Secret Manager (GCP)
- 🗝️ Database passwords auto-generated
- 📧 SMTP credentials em Secret Manager
- 📱 Twilio credentials em Secret Manager

### **Network Security**
- 🛡️ Private database (no public IP)
- 🔒 HTTPS automatic (Cloud Run)
- 🚪 IAM roles com least privilege
- 🧱 VPC connector para conectividade segura

## 🚨 Troubleshooting

### **Problemas Comuns**

#### **Docker não encontrado**
```bash
# Windows
winget install Docker.DockerDesktop

# Mac
brew install --cask docker

# Linux
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
```

#### **gcloud não autenticado**
```bash
gcloud auth login
gcloud auth application-default login
gcloud config set project leps-472702
```

#### **Terraform state locked**
```bash
terraform force-unlock <LOCK_ID>
```

#### **Cloud Run deploy falha**
```bash
# Verificar image registry
gcloud auth configure-docker us-central1-docker.pkg.dev

# Verificar permissões
gcloud projects add-iam-policy-binding leps-472702 \
    --member="user:$(gcloud config get-value account)" \
    --role="roles/run.admin"
```

### **Logs e Debugging**

#### **Local Development**
```bash
# Ver logs da aplicação
docker-compose logs -f app

# Acessar container
docker-compose exec app sh

# Verificar banco
docker-compose exec postgres psql -U lep_user -d lep_database
```

#### **GCP Environments**
```bash
# Ver logs do Cloud Run
gcloud run services logs read leps-backend-dev --region=us-central1

# Verificar status
gcloud run services describe leps-backend-dev --region=us-central1
```

## 📞 Suporte

### **Comandos Úteis**
```bash
# Verificar status de todos os ambientes
./scripts/deploy-interactive.sh --check-all

# Cleanup local
docker-compose down -v
docker system prune -f

# Verificar custos GCP
gcloud billing budgets list
```

### **Links Úteis**
- 📖 [Documentação GCP Cloud Run](https://cloud.google.com/run/docs)
- 🏗️ [Terraform GCP Provider](https://registry.terraform.io/providers/hashicorp/google/latest/docs)
- 🐳 [Docker Compose Reference](https://docs.docker.com/compose/compose-file/)
- 🔧 [Go Environment Setup](https://golang.org/doc/install)

---

## 🎉 Quick Start

1. **Clone e setup inicial:**
   ```bash
   git clone <repo>
   cd LEP-Back
   ```

2. **Para desenvolvimento local:**
   ```bash
   ./scripts/deploy-interactive.sh
   # Escolha opção 1 (local-dev)
   ```

3. **Para deploy em GCP:**
   ```bash
   ./scripts/deploy-interactive.sh
   # Escolha opção 2-4 baseado no ambiente desejado
   ```

O script irá guiá-lo através de todo o processo com validações automáticas e feedback detalhado! 🚀