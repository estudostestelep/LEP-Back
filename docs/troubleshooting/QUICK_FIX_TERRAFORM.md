# 🚑 LEP System - Quick Fix para Terraform

*Data: 23/09/2024*

## ❌ Problema Encontrado

O erro de "Duplicate resource configuration" acontece porque existem dois arquivos Terraform principais no mesmo diretório:
- `main.tf` (versão original com problemas de permissão)
- `main_simplified.tf` (versão híbrida que funciona)

## ✅ Solução Imediata (FEITO)

```bash
# 1. Mover arquivo conflitante
mkdir -p backups
mv main_original.tf backups/

# 2. Limpar locks
rm -f terraform.tfstate.lock.info .terraform.lock.hcl tfplan

# 3. Testar init
terraform init
# ✅ SUCESSO: Terraform has been successfully initialized!
```

## 🚀 Como Usar Agora

### **Opção 1: Script Quick Deploy (Recomendado)**
```bash
# Deploy completo automatizado
ENVIRONMENT=dev ./scripts/quick-deploy.sh
```

### **Opção 2: Comandos Manuais**
```bash
# 1. Executar bootstrap se necessário (primeira vez)
./scripts/bootstrap-gcp.sh

# 2. Deploy com Terraform
terraform init
terraform plan -var-file=environments/gcp-dev.tfvars -out=tfplan
terraform apply tfplan

# 3. Build e deploy da aplicação
docker build -t us-central1-docker.pkg.dev/leps-472702/lep-backend/lep-backend:latest .
docker push us-central1-docker.pkg.dev/leps-472702/lep-backend/lep-backend:latest
```

### **Opção 3: Deploy Interativo Original**
```bash
# Usar script original (agora que Terraform funciona)
ENVIRONMENT=gcp-dev ./scripts/deploy-interactive.sh
```

## 📋 Status Atual

### ✅ **Resolvido**
- [x] Conflitos de arquivos Terraform eliminados
- [x] `terraform init` funcionando
- [x] Abordagem híbrida implementada
- [x] Scripts de bootstrap criados
- [x] Script de quick deploy criado

### 🎯 **Próximos Passos**
1. **Testar deploy**: Use `ENVIRONMENT=dev ./scripts/quick-deploy.sh`
2. **Verificar recursos**: Use `./scripts/health-check.sh --environment gcp-dev`
3. **Atualizar secrets**: Adicionar chaves JWT reais

## 🏗️ Estrutura Atual de Arquivos

```
LEP-Back/
├── main.tf                    # ✅ Versão híbrida (usando data sources)
├── backups/
│   └── main_original.tf       # 📦 Versão original (backup)
├── scripts/
│   ├── bootstrap-gcp.sh       # 🚀 Bootstrap manual via gcloud
│   ├── quick-deploy.sh        # ⚡ Deploy rápido híbrido
│   ├── deploy-interactive.sh  # 🎯 Deploy interativo original
│   └── health-check.sh        # 🏥 Verificação de status
└── environments/
    ├── gcp-dev.tfvars        # 🧪 Config desenvolvimento
    ├── gcp-stage.tfvars      # 🚀 Config staging
    └── gcp-prd.tfvars        # 🌟 Config produção
```

## 🔧 Troubleshooting

### **Se ainda der erro de duplicata:**
```bash
# Verificar se há outros arquivos .tf
ls -la *.tf

# Mover qualquer arquivo extra para backup
mv problematic_file.tf backups/
```

### **Se terraform init falhar:**
```bash
# Limpar completamente
rm -rf .terraform .terraform.lock.hcl terraform.tfstate*
terraform init
```

### **Se resources não existirem:**
```bash
# Executar bootstrap primeiro
./scripts/bootstrap-gcp.sh

# Depois rodar terraform
terraform init
terraform plan -var-file=environments/gcp-dev.tfvars
```

## 📞 Comandos de Verificação

```bash
# Status de todos os ambientes
./scripts/health-check.sh --all

# Status específico
./scripts/health-check.sh --environment gcp-dev --detailed

# Verificar recursos GCP
gcloud run services list --region=us-central1
gcloud sql instances list
gcloud secrets list
```

---

## 🎉 Resumo da Solução

**PROBLEMA**: Conflito entre `main.tf` e `main_original.tf`
**SOLUÇÃO**: Mover `main_original.tf` para `backup/`
**RESULTADO**: ✅ `terraform init` funcionando
**PRÓXIMO PASSO**: Execute `ENVIRONMENT=dev ./scripts/quick-deploy.sh`

*Agora o LEP System está pronto para deploy sem problemas de Terraform!*