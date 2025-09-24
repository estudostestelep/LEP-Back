# 🚑 LEP System - Guia de Deploy de Correções (Hotfix)

*Data: 23/09/2024*

## 📋 Visão Geral

Este guia fornece comandos rápidos e seguros para deploy de correções urgentes no LEP System, com opções de rollback e validação.

## ⚡ Deploy Rápido (TL;DR)

### **Opção 1: Script Automatizado (Recomendado)**
```bash
# Deploy completo com validações
./scripts/deploy-interactive.sh

# Deploy direto para ambiente específico
ENVIRONMENT=gcp-stage ./scripts/deploy-interactive.sh
```

### **Opção 2: Comandos Manuais (Urgência)**
```bash
# Build e push
docker build -t us-central1-docker.pkg.dev/leps-472702/lep-backend/lep-backend:latest .
docker push us-central1-docker.pkg.dev/leps-472702/lep-backend/lep-backend:latest

# Deploy (escolher ambiente)
gcloud run deploy leps-backend-dev --image=us-central1-docker.pkg.dev/leps-472702/lep-backend/lep-backend:latest --region=us-central1
```

## 🎯 Ambientes de Deploy

### **1. 🧪 gcp-dev (Testes)**
- **Finalidade**: Validar correção antes de staging/produção
- **Comando**:
  ```bash
  ENVIRONMENT=gcp-dev ./scripts/deploy-interactive.sh
  ```
- **URL**: Será exibida no final do deploy
- **Rollback**: Rápido, pode quebrar sem impacto

### **2. 🚀 gcp-stage (Staging)**
- **Finalidade**: Validação final em ambiente similar à produção
- **Comando**:
  ```bash
  ENVIRONMENT=gcp-stage ./scripts/deploy-interactive.sh
  ```
- **Validação**: ✅ OBRIGATÓRIA antes de produção
- **Notificações**: Apenas SMTP (sem Twilio)

### **3. 🌟 gcp-prd (Produção)**
- **Finalidade**: Ambiente de produção com usuários reais
- **Comando**:
  ```bash
  ENVIRONMENT=gcp-prd ./scripts/deploy-interactive.sh
  ```
- **⚠️ CUIDADO**: Sempre testar em dev/stage primeiro
- **Notificações**: Twilio + SMTP ativas

## 🔧 Comandos por Situação

### **Correção Simples (Sem Dependências)**

#### **Desenvolvimento → Staging → Produção**
```bash
# 1. Deploy em dev para teste
ENVIRONMENT=gcp-dev ./scripts/deploy-interactive.sh

# 2. Validar funcionamento
curl https://DEV_URL/health
curl https://DEV_URL/ping

# 3. Deploy em staging
ENVIRONMENT=gcp-stage ./scripts/deploy-interactive.sh

# 4. Deploy em produção (após validação)
ENVIRONMENT=gcp-prd ./scripts/deploy-interactive.sh
```

### **Correção Urgente (Pular Dev)**

#### **Direto para Staging → Produção**
```bash
# 1. Deploy direto em staging
ENVIRONMENT=gcp-stage ./scripts/deploy-interactive.sh

# 2. Validação rápida
curl https://STAGE_URL/health

# 3. Deploy em produção
ENVIRONMENT=gcp-prd ./scripts/deploy-interactive.sh
```

### **Hotfix Crítico (Deploy Manual Rápido)**

```bash
# 1. Build local
docker build -t lep-backend-hotfix .

# 2. Tag para registry
docker tag lep-backend-hotfix us-central1-docker.pkg.dev/leps-472702/lep-backend/lep-backend:hotfix-$(date +%Y%m%d-%H%M)

# 3. Push
docker push us-central1-docker.pkg.dev/leps-472702/lep-backend/lep-backend:hotfix-$(date +%Y%m%d-%H%M)

# 4. Deploy imediato
gcloud run deploy leps-backend-prd \
    --image=us-central1-docker.pkg.dev/leps-472702/lep-backend/lep-backend:hotfix-$(date +%Y%m%d-%H%M) \
    --region=us-central1 \
    --no-traffic  # Deploy sem tráfego para validar

# 5. Validar e ativar tráfego
gcloud run services update-traffic leps-backend-prd \
    --to-latest \
    --region=us-central1
```

## 🔍 Validação Pós-Deploy

### **Health Checks Obrigatórios**

```bash
# Definir URL do ambiente
export APP_URL="https://sua-url-cloud-run"

# 1. Health check básico
curl -f $APP_URL/health
# Esperado: {"status":"healthy","environment":"..."}

# 2. Ping test
curl -f $APP_URL/ping
# Esperado: "pong"

# 3. Verificar logs (últimos 10min)
gcloud run services logs read leps-backend-prd \
    --region=us-central1 \
    --limit=50 \
    --since="10m"
```

### **Testes Funcionais Críticos**

```bash
# Headers necessários para APIs protegidas
export ORG_ID="your-org-uuid"
export PROJECT_ID="your-project-uuid"

# 1. Teste de autenticação (se aplicável)
curl -H "X-Lpe-Organization-Id: $ORG_ID" \
     -H "X-Lpe-Project-Id: $PROJECT_ID" \
     $APP_URL/health

# 2. Teste de endpoint crítico (ex: reservas)
curl -H "X-Lpe-Organization-Id: $ORG_ID" \
     -H "X-Lpe-Project-Id: $PROJECT_ID" \
     $APP_URL/reservation

# 3. Verificar notificações (se alteradas)
# Fazer reserva de teste e verificar se SMS/Email chegam
```

## 🔄 Rollback Rápido

### **Opção 1: Rollback via Cloud Run Console**
1. Acesse [Cloud Run Console](https://console.cloud.google.com/run)
2. Selecione o serviço (`leps-backend-xxx`)
3. Aba "REVISIONS"
4. Clique na revisão anterior estável
5. "MANAGE TRAFFIC" → 100% para revisão anterior

### **Opção 2: Rollback via CLI**

```bash
# 1. Listar revisões disponíveis
gcloud run revisions list \
    --service=leps-backend-prd \
    --region=us-central1

# 2. Identificar revisão anterior (ex: leps-backend-prd-00123-abc)
export PREVIOUS_REVISION="leps-backend-prd-00123-abc"

# 3. Rollback imediato
gcloud run services update-traffic leps-backend-prd \
    --to-revisions=$PREVIOUS_REVISION=100 \
    --region=us-central1

# 4. Verificar rollback
curl https://your-service-url/health
```

### **Opção 3: Rollback via Image Tag**

```bash
# 1. Deploy image anterior conhecida
gcloud run deploy leps-backend-prd \
    --image=us-central1-docker.pkg.dev/leps-472702/lep-backend/lep-backend:stable \
    --region=us-central1

# 2. Verificar funcionamento
curl https://your-service-url/health
```

## 📊 Monitoramento Pós-Deploy

### **Verificar Métricas (Primeiros 10 minutos)**

```bash
# 1. Verificar se há erros 5xx
gcloud logging read "resource.type=\"cloud_run_revision\"
    AND resource.labels.service_name=\"leps-backend-prd\"
    AND httpRequest.status>=500" \
    --freshness=10m \
    --limit=20

# 2. Verificar latência de resposta
gcloud logging read "resource.type=\"cloud_run_revision\"
    AND resource.labels.service_name=\"leps-backend-prd\"
    AND httpRequest.latency>\"1s\"" \
    --freshness=10m \
    --limit=10

# 3. Verificar memory/CPU usage se disponível
gcloud run services describe leps-backend-prd \
    --region=us-central1 \
    --format="value(status.traffic[0].latestRevision)"
```

## 🚨 Troubleshooting de Deploy

### **Problema: Build falha**

```bash
# Verificar logs detalhados
docker build --progress=plain --no-cache -t lep-backend .

# Verificar dependências Go
go mod tidy
go mod verify

# Testar build local
go build -o lep-test .
./lep-test --version
```

### **Problema: Push para registry falha**

```bash
# Reautenticar Docker
gcloud auth configure-docker us-central1-docker.pkg.dev

# Verificar acesso ao repositório
gcloud artifacts repositories describe lep-backend \
    --location=us-central1

# Verificar permissões
gcloud projects get-iam-policy leps-472702 | grep "$(gcloud config get-value account)"
```

### **Problema: Deploy Cloud Run falha**

```bash
# Verificar se service existe
gcloud run services describe leps-backend-prd --region=us-central1

# Verificar últimos logs de deploy
gcloud run services logs read leps-backend-prd \
    --region=us-central1 \
    --since="5m"

# Verificar se image existe no registry
gcloud artifacts docker images list us-central1-docker.pkg.dev/leps-472702/lep-backend
```

### **Problema: Service não responde**

```bash
# Verificar status do serviço
gcloud run services describe leps-backend-prd \
    --region=us-central1 \
    --format="value(status.conditions[0].message)"

# Verificar se está recebendo tráfego
gcloud run services describe leps-backend-prd \
    --region=us-central1 \
    --format="value(status.traffic[0].percent)"

# Forçar nova revisão
gcloud run deploy leps-backend-prd \
    --image=us-central1-docker.pkg.dev/leps-472702/lep-backend/lep-backend:latest \
    --region=us-central1 \
    --set-env-vars="DEPLOY_TIME=$(date)"
```

## 📋 Checklist de Deploy

### **Pré-Deploy**
- [ ] Código revisado e testado localmente
- [ ] Backup/tag da versão atual em produção
- [ ] Ambiente de staging testado (se não for urgência)
- [ ] Variáveis de ambiente conferidas
- [ ] Health checks funcionando localmente

### **Durante Deploy**
- [ ] Build completed sem warnings críticos
- [ ] Push para registry bem-sucedido
- [ ] Deploy Cloud Run sem erros
- [ ] Health check passing imediatamente
- [ ] Logs não mostram erros críticos

### **Pós-Deploy**
- [ ] URL de produção respondendo
- [ ] Health check returning "healthy"
- [ ] Funcionalidades críticas testadas
- [ ] Logs limpos por 5-10 minutos
- [ ] Notificações funcionando (se alteradas)
- [ ] Métricas estáveis

## 📞 Suporte e Escalação

### **Comandos de Emergência**
```bash
# Status completo de todos os serviços
./scripts/health-check.sh --all-environments

# Rollback imediato
gcloud run services update-traffic leps-backend-prd \
    --to-revisions=REVISION_ANTERIOR=100 \
    --region=us-central1

# Logs em tempo real
gcloud run services logs tail leps-backend-prd --region=us-central1
```

### **Contatos de Escalação**
- **Infra GCP**: Verificar console GCP
- **Logs centralizados**: Cloud Logging
- **Monitoramento**: Cloud Monitoring (se configurado)

---

*Use este guia para deployments rápidos e seguros. Sempre priorize a validação e tenha um plano de rollback pronto.*