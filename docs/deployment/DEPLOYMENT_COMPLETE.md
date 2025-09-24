# 🚀 LEP System - Guia Completo de Deploy

Este é o guia único e consolidado para deployment do LEP System em todos os ambientes.

## 🎯 Ambientes Suportados

### 1. 🏠 **local-dev** - Desenvolvimento Local
- Docker Compose com PostgreSQL local
- Hot reload e debug habilitado
- Configuração em `.env`

### 2. ☁️ **gcp-dev** - GCP Desenvolvimento
- Cloud Run + Cloud SQL
- Artifact Registry
- Secret Manager
- **URL**: https://leps-backend-dev-516622888070.us-central1.run.app

### 3. 🚀 **gcp-stage** - GCP Staging
- Configuração similar ao dev
- Dados de teste estáveis

### 4. 🌟 **gcp-prd** - GCP Produção
- Alta disponibilidade
- Backups automatizados

## ⚡ Deploy Rápido - GCP

### Método 1: Script Automático (Recomendado)
```bash
# Deploy completo híbrido (gcloud + Terraform)
ENVIRONMENT=dev ./scripts/quick-deploy.sh
```

### Método 2: Comandos Manuais
```bash
# 1. Bootstrap recursos via gcloud
./scripts/bootstrap-gcp.sh

# 2. Deploy aplicação
docker build -t us-central1-docker.pkg.dev/leps-472702/lep-backend/lep-backend:latest .
docker push us-central1-docker.pkg.dev/leps-472702/lep-backend/lep-backend:latest

# 3. Update Cloud Run
gcloud run deploy leps-backend-dev \
  --image=us-central1-docker.pkg.dev/leps-472702/lep-backend/lep-backend:latest \
  --region=us-central1 \
  --project=leps-472702
```

## 🔧 Configurações Atuais

### Banco de Dados (Cloud SQL)
- **Instância**: `leps-postgres-dev`
- **Database**: `lep_database`
- **Usuário**: `lep_user`
- **Senha**: `123456` (dev only)

### Credenciais de Login (após seed)
- **Admin**: `admin@lep-demo.com` / `password`
- **Garçom**: `garcom@lep-demo.com` / `password`
- **Gerente**: `gerente@lep-demo.com` / `password`

## 📋 Scripts Disponíveis

| Script | Função |
|--------|---------|
| `quick-deploy.sh` | Deploy completo automatizado |
| `bootstrap-gcp.sh` | Criar recursos base via gcloud |
| `health-check.sh` | Verificar status dos ambientes |
| `run_seed.sh` | Popular banco com dados de teste |

## 🐛 Troubleshooting

### Erro "PRI *" no Cloud Run
```bash
# Verificar logs
gcloud run services logs read leps-backend-dev --region=us-central1

# Redeploy forçado
gcloud run deploy leps-backend-dev \
  --image=us-central1-docker.pkg.dev/leps-472702/lep-backend/lep-backend:latest \
  --region=us-central1 --project=leps-472702
```

### Health Checks
```bash
# Local
curl http://localhost:8080/ping

# GCP
curl https://leps-backend-dev-516622888070.us-central1.run.app/ping
```

## 📁 Estrutura de Arquivos
- `main.tf` - Terraform híbrido (Cloud Run only)
- `scripts/` - Scripts de automação
- `docs/deployment/` - Documentação consolidada
- `environments/` - Configs por ambiente

**Notas:**
- Abordagem híbrida resolve problemas de permissão IAM
- Bootstrap via gcloud + Deploy via Terraform
- Secrets gerenciados via Secret Manager
- Logs centralizados no Cloud Logging