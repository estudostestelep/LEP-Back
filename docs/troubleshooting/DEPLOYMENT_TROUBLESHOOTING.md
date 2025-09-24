# 🔧 LEP Backend - Troubleshooting e Deployment

*Atualizado: 23/09/2024*

## 🚨 Problema Identificado

O deployment via Docker estava falhando devido a problemas de conectividade de rede que impedem:
1. Download de pacotes Alpine Linux (`apk add`)
2. Download de dependências Go (`go mod download`)
3. Acesso ao Go module proxy

## ✅ Soluções Implementadas

### 1. **Execução Local (Recomendada)**

Para contornar os problemas de rede, use os scripts locais:

#### Windows:
```batch
run-local.bat
```

#### Linux/Mac:
```bash
./run-local.sh
```

Estes scripts:
- ✅ Verificam se Go está instalado
- ✅ Tentam baixar dependências (com fallback)
- ✅ Fazem build da aplicação
- ✅ Iniciam o servidor na porta 8080

### 2. **Docker Fixes Implementados**

Múltiplas versões do Dockerfile.dev foram criadas:

#### `Dockerfile.dev` (Atual - Minimal)
```dockerfile
FROM golang:1.23-alpine
WORKDIR /app
COPY . .
RUN mkdir -p /app/logs
EXPOSE 8080
ENV GO_ENV=development PORT=8080 GIN_MODE=debug
CMD ["go", "run", "main.go"]
```

#### Versões Alternativas Criadas:
- `Dockerfile.dev.backup` - Versão original
- `Dockerfile.dev.robust` - Com retry logic e múltiplos mirrors
- `Dockerfile.dev.minimal` - Apenas essenciais
- `Dockerfile.dev.vendor` - Usando vendor directory
- `Dockerfile.dev.offline` - Para ambientes sem internet

### 3. **Twilio Security Fix Aplicado**

```go
// ANTES (VULNERABILIDADE):
func (t *TwilioService) ValidateWebhookSignature(signature, url, body string) bool {
    return true // ❌ INSEGURO
}

// DEPOIS (SEGURO):
func (t *TwilioService) ValidateWebhookSignature(signature, requestUrl, body string) bool {
    authToken := t.AuthToken
    if authToken == "" {
        authToken = os.Getenv("TWILIO_AUTH_TOKEN")
    }

    mac := hmac.New(sha1.New, []byte(authToken))
    mac.Write([]byte(requestUrl + body))
    expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

    return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
```

### 4. **Reports Routes Implementados**

Rotas de relatórios agora funcionais:
- `GET /reports/occupancy` - Métricas de ocupação
- `GET /reports/reservations` - Estatísticas de reservas
- `GET /reports/waitlist` - Métricas de fila de espera
- `GET /reports/leads` - Relatório de leads
- `GET /reports/export/:type` - Export CSV

## 🏃‍♂️ Como Executar Agora

### Opção 1: Local (Recomendada)
```bash
# Windows
run-local.bat

# Linux/Mac
./run-local.sh
```

### Opção 2: Docker (Se rede funcionar)
```bash
docker-compose build --no-cache app
docker-compose up -d
```

### Opção 3: Go Direto
```bash
go mod tidy
go run main.go
```

## 🔍 Verificação de Funcionamento

Após iniciar, teste:

```bash
# Health check
curl http://localhost:8080/health
# Esperado: {"status":"healthy"}

# Ping test
curl http://localhost:8080/ping
# Esperado: "pong"

# Reports test (requer headers)
curl -H "X-Lpe-Organization-Id: uuid" -H "X-Lpe-Project-Id: uuid" \
     http://localhost:8080/reports/occupancy
```

## 📋 Status de Correções

### ✅ Completado
- [x] Implementação do setupReportsRoutes
- [x] Correção da validação Twilio signature
- [x] Fix do campo permissions do User
- [x] Criação de scripts de deployment local
- [x] Múltiplas opções de Dockerfile

### ⚠️ Problemas de Ambiente
- [ ] Conectividade de rede (problema do ambiente, não do código)
- [ ] Access to Alpine repositories (problema de proxy/firewall)
- [ ] Go module proxy access (problema de DNS/proxy)

## 🛠️ Próximos Passos

1. **Teste Local**: Use `run-local.bat` para validar funcionamento
2. **Configurar Proxy**: Se em ambiente corporativo, configurar proxy HTTP
3. **Deploy em Produção**: Usar ambiente com conectividade estável

## 🔧 Configuração de Proxy (Se Necessário)

Se estiver atrás de proxy corporativo:

```bash
# Configure Go proxy
export GOPROXY=https://proxy.company.com
export GONOPROXY="github.com/company/*"

# Configure Docker
export HTTP_PROXY=http://proxy.company.com:8080
export HTTPS_PROXY=http://proxy.company.com:8080
```

## 🚨 Problemas GCP e Terraform

### **Problema: Permissões Terraform**

Durante o setup inicial com Terraform, encontramos problemas de permissões mesmo com role de `owner`:

```bash
# Erro comum
Error: Permission 'iam.serviceAccounts.create' denied on resource
Error: Permission 'secretmanager.secrets.create' denied for resource
```

### **Solução: Abordagem Híbrida**

1. **Resources criados manualmente via gcloud**:
   ```bash
   # APIs
   gcloud services enable secretmanager.googleapis.com sqladmin.googleapis.com run.googleapis.com

   # Service Account
   gcloud iam service-accounts create lep-backend-sa

   # Artifact Registry
   gcloud artifacts repositories create lep-backend --repository-format=docker --location=us-central1

   # Secrets
   gcloud secrets create jwt-private-key-dev --replication-policy="automatic"
   ```

2. **Terraform usa data sources** para referenciar recursos existentes:
   ```hcl
   data "google_service_account" "lep_backend_sa" {
     account_id = "lep-backend-sa"
     project    = var.project_id
   }
   ```

### **Problema: Cloud SQL "No IP enabled"**

```bash
# Erro
Invalid request: At least one of Public IP or Private IP or PSC connectivity must be enabled
```

**Solução**:
```bash
# Criar com IP público temporariamente
gcloud sql instances create leps-postgres-dev \
    --database-version=POSTGRES_15 \
    --tier=db-f1-micro \
    --region=us-central1 \
    --assign-ip
```

### **Problema: Propagação de Permissões**

Permissões IAM podem demorar até 5 minutos para propagar.

**Solução**:
```bash
# Aguardar e testar
sleep 300
gcloud projects test-iam-permissions leps-472702 \
    --permissions="iam.serviceAccounts.create,secretmanager.secrets.create"
```

## 🔧 Troubleshooting de Conta GCP

### **Problema: Multiple Accounts**

```bash
# Verificar conta ativa
gcloud auth list

# Trocar conta se necessário
gcloud config set account novo_email@gmail.com
```

### **Problema: Project não encontrado**

```bash
# Listar projetos disponíveis
gcloud projects list

# Definir projeto correto
gcloud config set project leps-472702
```

### **Problema: Docker Registry Access**

```bash
# Erro: permission denied on registry
# Solução:
gcloud auth configure-docker us-central1-docker.pkg.dev

# Verificar acesso
gcloud artifacts repositories list --location=us-central1
```

## 🚀 Deploy em Ambientes GCP

### **Checklist Pré-Deploy**

- [ ] **Autenticação**: `gcloud auth list` mostra conta correta
- [ ] **Projeto**: `gcloud config get-value project` = `leps-472702`
- [ ] **Docker**: `gcloud auth configure-docker` configurado
- [ ] **Permissões**: Roles necessárias ativas
- [ ] **APIs**: Todas as APIs habilitadas

### **Comando de Deploy Completo**

```bash
# Deploy interativo (recomendado)
./scripts/deploy-interactive.sh

# Deploy direto por ambiente
ENVIRONMENT=gcp-dev ./scripts/deploy-interactive.sh
ENVIRONMENT=gcp-stage ./scripts/deploy-interactive.sh
ENVIRONMENT=gcp-prd ./scripts/deploy-interactive.sh
```

### **Deploy Manual de Emergência**

```bash
# 1. Build e push
docker build -t us-central1-docker.pkg.dev/leps-472702/lep-backend/lep-backend:$(date +%Y%m%d-%H%M) .
docker push us-central1-docker.pkg.dev/leps-472702/lep-backend/lep-backend:$(date +%Y%m%d-%H%M)

# 2. Deploy Cloud Run
gcloud run deploy leps-backend-dev \
    --image=us-central1-docker.pkg.dev/leps-472702/lep-backend/lep-backend:$(date +%Y%m%d-%H%M) \
    --region=us-central1 \
    --allow-unauthenticated

# 3. Verificar
curl https://$(gcloud run services describe leps-backend-dev --region=us-central1 --format="value(status.url)")/health
```

## 📋 Recursos Adicionais

### **Novos Guias Criados**

- 📖 **ACCOUNT_MIGRATION.md**: Guia completo para trocar conta GCP
- 🚑 **HOTFIX_DEPLOYMENT.md**: Deploy rápido de correções
- 🏗️ **DEPLOYMENT-GUIDE.md**: Deploy interativo multi-ambiente

### **Scripts Importantes**

- 🔧 **scripts/deploy-interactive.sh**: Deploy automatizado com validações
- 🏥 **scripts/health-check.sh**: Verificação de status (a ser criado)
- 🛠️ **scripts/troubleshoot.sh**: Diagnóstico automatizado (a ser criado)

## 📞 Suporte Expandido

### **Logs da Aplicação**
```bash
# Local
./run-local.sh  # Logs detalhados
docker-compose logs -f app  # Docker logs

# GCP
gcloud run services logs read SERVICE_NAME --region=us-central1
gcloud run services logs tail SERVICE_NAME --region=us-central1  # Real-time
```

### **Build Errors**
```bash
# Diagnóstico local
go build .
go mod verify
go mod tidy

# Docker build debug
docker build --progress=plain --no-cache .
```

### **Comandos de Diagnóstico**
```bash
# Status completo do sistema
echo "=== GCLOUD CONFIG ==="
gcloud config list
echo "=== AUTH STATUS ==="
gcloud auth list
echo "=== PROJECT PERMISSIONS ==="
gcloud projects get-iam-policy leps-472702
echo "=== SERVICES STATUS ==="
gcloud run services list --region=us-central1
```

### **Links de Apoio**
- 📖 [GCP IAM Troubleshooting](https://cloud.google.com/iam/docs/troubleshooting-access)
- 🐳 [Cloud Run Troubleshooting](https://cloud.google.com/run/docs/troubleshooting)
- 🔧 [Artifact Registry Auth](https://cloud.google.com/artifact-registry/docs/docker/authentication)

---

*Sistema LEP Backend agora suporta desenvolvimento local E deploy completo em GCP*