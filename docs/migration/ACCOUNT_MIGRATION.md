# 🔄 LEP System - Guia de Migração de Conta GCP

*Data: 23/09/2024*

## 📋 Visão Geral

Este guia fornece instruções passo-a-passo para migrar o LEP System de uma conta GCP para outra, incluindo validações, permissões e troubleshooting.

## ⚠️ Pré-requisitos

### **Informações Necessárias:**
- ✅ **Email da nova conta GCP**: `novo_email@gmail.com`
- ✅ **Project ID**: `leps-472702` (se mantido)
- ✅ **Permissões**: Owner ou roles específicas
- ✅ **Access às chaves JWT**: Arquivos `.pem`
- ✅ **Credenciais Twilio/SMTP**: Se aplicável

### **Ferramentas Instaladas:**
- ✅ Google Cloud CLI (`gcloud`)
- ✅ Docker
- ✅ Terraform (se usando)
- ✅ Go 1.22+

## 🚀 Passo-a-Passo

### **1. Backup da Configuração Atual**

```bash
# Backup das configurações atuais
mkdir -p ~/lep-backup-$(date +%Y%m%d)
cd ~/lep-backup-$(date +%Y%m%d)

# Backup configuração gcloud
gcloud config configurations list > gcloud-configs.txt
gcloud config list > current-config.txt

# Backup terraform state (se aplicável)
cp -r /path/to/lep-backend/terraform.tfstate* .

# Backup variáveis de ambiente
cp -r /path/to/lep-backend/environments/ .
```

### **2. Logout da Conta Atual**

```bash
# Ver contas autenticadas atualmente
gcloud auth list

# Logout de todas as contas
gcloud auth revoke --all

# Limpar cache de credenciais
gcloud auth application-default revoke
```

### **3. Autenticar Nova Conta**

```bash
# Login interativo
gcloud auth login
# Escolher a nova conta no browser

# Configurar Application Default Credentials
gcloud auth application-default login
# Escolher a mesma conta

# Verificar autenticação
gcloud auth list
```

### **4. Configurar Projeto e Região**

```bash
# Definir projeto padrão
gcloud config set project leps-472702

# Definir região padrão
gcloud config set compute/region us-central1
gcloud config set compute/zone us-central1-a

# Verificar configuração
gcloud config list
```

### **5. Verificar Permissões da Nova Conta**

```bash
# Listar permissões da conta no projeto
gcloud projects get-iam-policy leps-472702

# Verificar se tem role necessária
gcloud projects get-iam-policy leps-472702 \
    --flatten="bindings[].members" \
    --format="table(bindings.role)" \
    --filter="bindings.members:novo_email@gmail.com"
```

### **6. Conceder Permissões (Se Necessário)**

#### **Via Console Web (Recomendado para Owner):**
1. Acesse [GCP Console IAM](https://console.cloud.google.com/iam-admin/iam)
2. Selecione projeto `leps-472702`
3. Clique "ADD" ou "GRANT ACCESS"
4. Email: `novo_email@gmail.com`
5. Role: `Owner` ou roles específicas abaixo

#### **Via CLI (Se você já tem acesso):**
```bash
# Role de Owner (mais simples)
gcloud projects add-iam-policy-binding leps-472702 \
    --member="user:novo_email@gmail.com" \
    --role="roles/owner"

# OU roles específicas (mais seguro)
gcloud projects add-iam-policy-binding leps-472702 \
    --member="user:novo_email@gmail.com" \
    --role="roles/run.admin"

gcloud projects add-iam-policy-binding leps-472702 \
    --member="user:novo_email@gmail.com" \
    --role="roles/cloudsql.admin"

gcloud projects add-iam-policy-binding leps-472702 \
    --member="user:novo_email@gmail.com" \
    --role="roles/secretmanager.admin"

gcloud projects add-iam-policy-binding leps-472702 \
    --member="user:novo_email@gmail.com" \
    --role="roles/artifactregistry.admin"

gcloud projects add-iam-policy-binding leps-472702 \
    --member="user:novo_email@gmail.com" \
    --role="roles/iam.serviceAccountAdmin"
```

### **7. Configurar Docker Registry**

```bash
# Configurar autenticação Docker
gcloud auth configure-docker us-central1-docker.pkg.dev

# Testar acesso ao registry
gcloud artifacts repositories list --location=us-central1
```

### **8. Validar Acesso aos Recursos**

```bash
# Verificar Cloud Run services
gcloud run services list --region=us-central1

# Verificar Cloud SQL instances
gcloud sql instances list

# Verificar Artifact Registry
gcloud artifacts repositories list --location=us-central1

# Verificar Secret Manager
gcloud secrets list

# Verificar Service Accounts
gcloud iam service-accounts list
```

### **9. Testar Deploy**

```bash
# Navegar para o projeto
cd /path/to/LEP-Back

# Testar script de deploy
./scripts/deploy-interactive.sh

# Escolher ambiente de teste (gcp-dev recomendado)
# Verificar se tudo funciona sem erros
```

## 🔧 Troubleshooting

### **Problema: "Permission denied"**

```bash
# Verificar conta ativa
gcloud config get-value account

# Verificar permissões específicas
gcloud projects test-iam-permissions leps-472702 \
    --permissions="cloudsql.instances.create,secretmanager.secrets.create,run.services.create"

# Se faltar permissões, repetir passo 6
```

### **Problema: "Docker push permission denied"**

```bash
# Reconfigurar Docker auth
gcloud auth configure-docker us-central1-docker.pkg.dev

# Verificar se pode listar repositórios
gcloud artifacts repositories list --location=us-central1

# Se não funcionar, verificar role artifactregistry.admin
```

### **Problema: "Terraform state locked"**

```bash
# Listar locks ativos
terraform force-unlock <LOCK_ID>

# Se não souber o LOCK_ID, verificar erro anterior
```

### **Problema: "Application Default Credentials"**

```bash
# Reconfigurar ADC
gcloud auth application-default login

# Verificar arquivo de credenciais
ls -la ~/.config/gcloud/application_default_credentials.json

# Definir projeto quota
gcloud auth application-default set-quota-project leps-472702
```

## ✅ Checklist de Validação

Após completar a migração, verificar:

- [ ] **Autenticação**: `gcloud auth list` mostra nova conta
- [ ] **Projeto**: `gcloud config get-value project` = `leps-472702`
- [ ] **Permissões**: `gcloud projects get-iam-policy leps-472702` inclui nova conta
- [ ] **Docker**: `gcloud artifacts repositories list` funciona
- [ ] **Cloud Run**: `gcloud run services list` funciona
- [ ] **Cloud SQL**: `gcloud sql instances list` funciona
- [ ] **Secrets**: `gcloud secrets list` funciona
- [ ] **Deploy**: `./scripts/deploy-interactive.sh` funciona em gcp-dev

## 🚨 Rollback (Em Caso de Problemas)

### **1. Restaurar Conta Original**
```bash
# Logout nova conta
gcloud auth revoke novo_email@gmail.com

# Login conta original
gcloud auth login conta_original@gmail.com
gcloud auth application-default login

# Restaurar configuração
gcloud config set project leps-472702
```

### **2. Restaurar Terraform State**
```bash
# Se state foi corrompido
cp ~/lep-backup-YYYYMMDD/terraform.tfstate* .
```

## 📞 Suporte

### **Comandos Úteis de Diagnóstico**
```bash
# Status completo da conta
echo "=== GCLOUD CONFIG ==="
gcloud config list
echo ""
echo "=== AUTH LIST ==="
gcloud auth list
echo ""
echo "=== PROJECT PERMISSIONS ==="
gcloud projects get-iam-policy leps-472702
echo ""
echo "=== SERVICES STATUS ==="
gcloud services list --enabled
```

### **Logs de Debug**
```bash
# Habilitar logs detalhados
export CLOUDSDK_CORE_VERBOSITY=debug
gcloud <comando-que-falha>
```

## 📋 Documento de Migração Concluída

Após sucesso da migração, documentar:

```text
Data: ___________
Conta anterior: ___________
Conta nova: ___________
Projeto: leps-472702
Status: ✅ Migração concluída com sucesso
Validações: ✅ Todas passaram
Deploy test: ✅ gcp-dev funcionando
```

---

*Este documento deve ser usado sempre que houver necessidade de trocar a conta GCP responsável pelo projeto LEP System.*