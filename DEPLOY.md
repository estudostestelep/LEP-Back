# 🚀 Guia Rápido de Deploy - LEP Backend

Este é um guia de referência rápida para deploy do LEP Backend. Para instruções detalhadas, consulte [INFRASTRUCTURE.md](./INFRASTRUCTURE.md).

## ⚡ Deploy Rápido

### 1. Pré-requisitos
```bash
# Instale as ferramentas (se necessário)
# - gcloud CLI
# - terraform
# - docker

# Autentique-se
gcloud auth login
gcloud auth application-default login
```

### 2. Configuração Inicial
```bash
# Clone o repositório
git clone <url-do-repo>
cd LEP-Back

# Configure o projeto
./scripts/setup-secrets.sh --project-id SEU-PROJECT-ID --environment dev

# Gere chaves JWT
openssl genpkey -algorithm RSA -out jwt_private_key.pem -pkcs8 -aes256
openssl rsa -pubout -in jwt_private_key.pem -out jwt_public_key.pem

# Edite terraform.tfvars com suas configurações
cp terraform.tfvars.example terraform.tfvars
# Edite o arquivo com suas chaves e configurações
```

### 3. Deploy
```bash
# Deploy completo (infraestrutura + aplicação)
./scripts/deploy.sh --environment dev

# OU deploy por etapas:

# 1. Apenas infraestrutura
terraform init
terraform apply -var="environment=dev"

# 2. Apenas aplicação
./scripts/deploy.sh --skip-terraform --environment dev
```

## 🔧 Comandos Essenciais

### Deploy por Ambiente
```bash
# Desenvolvimento
./scripts/deploy.sh --environment dev

# Staging
./scripts/deploy.sh --environment staging

# Produção
./scripts/deploy.sh --environment prod
```

### Verificação de Status
```bash
# URL do serviço
terraform output service_url

# Logs em tempo real
gcloud logs tail "resource.type=cloud_run_revision"

# Status dos recursos
gcloud run services list
gcloud sql instances list
```

### Testes Rápidos
```bash
# Health check
curl $(terraform output -raw service_url)/health

# Ping
curl $(terraform output -raw service_url)/ping
```

## 🔄 CI/CD (Opcional)

### Configurar GitHub Actions
```bash
# Conectar repositório
gcloud builds triggers create github \
  --repo-name=LEP-Back \
  --repo-owner=SEU_USERNAME \
  --branch-pattern="^main$" \
  --build-config=cloudbuild.yaml
```

## 📋 Checklist de Deploy

### Antes do Deploy
- [ ] gcloud CLI instalado e autenticado
- [ ] terraform instalado
- [ ] docker instalado e rodando
- [ ] terraform.tfvars configurado
- [ ] Chaves JWT geradas
- [ ] Credenciais Twilio/SMTP (se necessário)

### Deploy de Produção
- [ ] Usar `environment = "prod"`
- [ ] Configurar `db_tier = "db-n1-standard-1"` ou superior
- [ ] Configurar `db_availability_type = "REGIONAL"`
- [ ] Configurar `min_instances = 1` ou superior
- [ ] Ativar `enable_deletion_protection = true`
- [ ] Configurar domínio personalizado (se necessário)
- [ ] Configurar alertas de monitoramento

### Pós-Deploy
- [ ] Testar endpoints de health
- [ ] Verificar logs
- [ ] Testar autenticação
- [ ] Configurar DNS (se usando domínio personalizado)
- [ ] Configurar alertas

## 🚨 Troubleshooting Rápido

### Container não inicia
```bash
# Verificar logs
gcloud run services logs read NOME_DO_SERVICO --region=us-central1

# Testar localmente
docker run -p 8080:8080 IMAGEM_URL
```

### Erro de banco de dados
```bash
# Verificar instância
gcloud sql instances describe NOME_DA_INSTANCIA

# Verificar conectividade
gcloud sql connect NOME_DA_INSTANCIA --user=USUARIO
```

### Erro de secrets
```bash
# Listar secrets
gcloud secrets list

# Verificar conteúdo
gcloud secrets versions access latest --secret="NOME_DO_SECRET"
```

## 💡 Dicas

1. **Use `--dry-run`** para simular deploys
2. **Configure min-instances=0** para desenvolvimento (economizar)
3. **Use tags de versão** para rollbacks fáceis
4. **Monitore custos** no Console GCP
5. **Faça backup** dos secrets importantes

## 📞 Links Úteis

- [Guia Completo (INFRASTRUCTURE.md)](./INFRASTRUCTURE.md)
- [Console GCP](https://console.cloud.google.com)
- [Cloud Run Console](https://console.cloud.google.com/run)
- [Cloud SQL Console](https://console.cloud.google.com/sql)
- [Secret Manager](https://console.cloud.google.com/security/secret-manager)