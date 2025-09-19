# Guia de Infraestrutura - LEP Backend

Este guia fornece instruções detalhadas para configurar e implantar a infraestrutura do LEP Backend no Google Cloud Platform (GCP).

## 📋 Índice

1. [Pré-requisitos](#pré-requisitos)
2. [Configuração Inicial](#configuração-inicial)
3. [Configuração de Secrets](#configuração-de-secrets)
4. [Deploy da Infraestrutura](#deploy-da-infraestrutura)
5. [Deploy da Aplicação](#deploy-da-aplicação)
6. [CI/CD com Cloud Build](#cicd-com-cloud-build)
7. [Monitoramento e Logs](#monitoramento-e-logs)
8. [Troubleshooting](#troubleshooting)
9. [Custos Estimados](#custos-estimados)

## 🔧 Pré-requisitos

### Ferramentas Necessárias

1. **Google Cloud CLI (gcloud)**
   ```bash
   # Windows (via Chocolatey)
   choco install gcloudsdk

   # macOS (via Homebrew)
   brew install google-cloud-sdk

   # Linux (via package manager)
   curl https://sdk.cloud.google.com | bash
   ```

2. **Terraform**
   ```bash
   # Windows (via Chocolatey)
   choco install terraform

   # macOS (via Homebrew)
   brew install terraform

   # Linux (via package manager)
   wget https://releases.hashicorp.com/terraform/1.6.0/terraform_1.6.0_linux_amd64.zip
   unzip terraform_1.6.0_linux_amd64.zip
   sudo mv terraform /usr/local/bin/
   ```

3. **Docker**
   - Instale o [Docker Desktop](https://www.docker.com/products/docker-desktop/)

4. **OpenSSL** (para geração de chaves JWT)
   - Já incluído no Git Bash (Windows)
   - Pré-instalado no macOS/Linux

### Conta GCP

1. **Crie um projeto GCP**
   ```bash
   gcloud projects create SEU-PROJECT-ID --name="LEP Backend"
   gcloud config set project SEU-PROJECT-ID
   ```

2. **Ative o faturamento**
   - Acesse o [Console GCP](https://console.cloud.google.com)
   - Vá para "Billing" e associe uma conta de faturamento

3. **Autentique-se**
   ```bash
   gcloud auth login
   gcloud auth application-default login
   ```

## ⚙️ Configuração Inicial

### 1. Clone e Configure o Repositório

```bash
git clone <url-do-repositorio>
cd LEP-Back
```

### 2. Execute o Script de Setup

**Linux/macOS:**
```bash
chmod +x scripts/setup-secrets.sh
./scripts/setup-secrets.sh --project-id SEU-PROJECT-ID --environment dev
```

**Windows PowerShell:**
```powershell
# Execute como Administrador se necessário
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
.\scripts\setup-secrets.ps1 -ProjectId "SEU-PROJECT-ID" -Environment "dev"
```

Este script irá:
- Validar o projeto GCP
- Habilitar APIs necessárias
- Criar o arquivo `terraform.tfvars` baseado no exemplo

## 🔐 Configuração de Secrets

### 1. Gerar Chaves JWT

**Importante**: Use chaves fortes para produção!

```bash
# Gerar chave privada (será solicitada uma senha)
openssl genpkey -algorithm RSA -out jwt_private_key.pem -pkcs8 -aes256

# Gerar chave pública
openssl rsa -pubout -in jwt_private_key.pem -out jwt_public_key.pem
```

### 2. Configurar terraform.tfvars

Edite o arquivo `terraform.tfvars` com suas configurações:

```hcl
# Configurações obrigatórias
project_id = "seu-project-id"
environment = "dev"  # ou staging, prod

# Chaves JWT (copie o conteúdo dos arquivos .pem)
jwt_private_key = """-----BEGIN PRIVATE KEY-----
CONTEUDO_DA_CHAVE_PRIVADA
-----END PRIVATE KEY-----"""

jwt_public_key = """-----BEGIN PUBLIC KEY-----
CONTEUDO_DA_CHAVE_PUBLICA
-----END PUBLIC KEY-----"""

# Configurações opcionais de produção
db_tier = "db-n1-standard-1"  # Para produção
db_availability_type = "REGIONAL"  # Para alta disponibilidade
min_instances = 1  # Para evitar cold starts
max_instances = 20  # Baseado na demanda esperada

# Configurações Twilio (opcional)
twilio_account_sid = "ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
twilio_auth_token = "seu_auth_token"
twilio_phone_number = "+5511999999999"

# Configurações SMTP (opcional)
smtp_username = "seu.email@gmail.com"
smtp_password = "sua_senha_de_app"  # Use App Password do Gmail
```

### 3. Configurações de Terceiros

#### Twilio (SMS/WhatsApp)
1. Crie conta em [twilio.com](https://www.twilio.com)
2. Obtenha Account SID e Auth Token
3. Configure um número Twilio
4. Para WhatsApp Business, siga o [guia oficial](https://www.twilio.com/docs/whatsapp)

#### SMTP (Email)
1. **Gmail**: Ative autenticação de 2 fatores e gere uma [senha de app](https://support.google.com/accounts/answer/185833)
2. **Outros provedores**: Configure SMTP host, porta, usuário e senha

## 🚀 Deploy da Infraestrutura

### 1. Validação e Deploy

```bash
# Inicializar Terraform
terraform init

# Planejar deploy (verificar o que será criado)
terraform plan -var="environment=dev"

# Aplicar configurações (confirme com 'yes')
terraform apply -var="environment=dev"
```

### 2. Verificar Recursos Criados

```bash
# Listar outputs do Terraform
terraform output

# Verificar Cloud Run service
gcloud run services list

# Verificar Cloud SQL instance
gcloud sql instances list

# Verificar secrets
gcloud secrets list
```

## 📦 Deploy da Aplicação

### 1. Deploy Manual

**Linux/macOS:**
```bash
chmod +x scripts/deploy.sh
./scripts/deploy.sh --environment dev
```

**Windows PowerShell:**
```powershell
.\scripts\deploy.ps1 -Environment dev
```

### 2. Deploy por Etapas

```bash
# Apenas infraestrutura
./scripts/deploy.sh --skip-build --environment dev

# Apenas build e deploy da aplicação
./scripts/deploy.sh --skip-terraform --environment dev

# Dry run (simular sem executar)
./scripts/deploy.sh --dry-run --environment dev
```

### 3. Teste Manual

```bash
# Obter URL do serviço
SERVICE_URL=$(terraform output -raw service_url)

# Testar endpoints
curl $SERVICE_URL/health
curl $SERVICE_URL/ping

# Testar autenticação (substitua dados reais)
curl -X POST $SERVICE_URL/login \
  -H "Content-Type: application/json" \
  -H "X-Lpe-Organization-Id: ORG_ID" \
  -H "X-Lpe-Project-Id: PROJECT_ID" \
  -d '{"username":"admin","password":"password"}'
```

## 🔄 CI/CD com Cloud Build

### 1. Configurar GitHub Integration

```bash
# Conectar repositório GitHub
gcloud source repos create lep-backend
gcloud builds triggers create github \
  --repo-name=LEP-Back \
  --repo-owner=SEU_GITHUB_USERNAME \
  --branch-pattern="^dev$" \
  --build-config=cloudbuild.yaml \
  --name=lep-backend-dev
```

### 2. Configurar IAM para Cloud Build

```bash
# Obter email do Cloud Build service account
PROJECT_ID=$(gcloud config get-value project)
BUILD_SA="${PROJECT_ID}@cloudbuild.gserviceaccount.com"

# Adicionar permissões necessárias
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$BUILD_SA" \
  --role="roles/run.admin"

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$BUILD_SA" \
  --role="roles/artifactregistry.writer"

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$BUILD_SA" \
  --role="roles/secretmanager.secretAccessor"

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$BUILD_SA" \
  --role="roles/iam.serviceAccountUser"
```

### 3. Triggers por Ambiente

```bash
# Trigger para staging (branch staging)
gcloud builds triggers create github \
  --repo-name=LEP-Back \
  --repo-owner=SEU_GITHUB_USERNAME \
  --branch-pattern="^staging$" \
  --build-config=cloudbuild.yaml \
  --name=lep-backend-staging

# Trigger para produção (branch main)
gcloud builds triggers create github \
  --repo-name=LEP-Back \
  --repo-owner=SEU_GITHUB_USERNAME \
  --branch-pattern="^main$" \
  --build-config=cloudbuild.yaml \
  --name=lep-backend-prod
```

## 📊 Monitoramento e Logs

### 1. Cloud Logging

```bash
# Ver logs do Cloud Run
gcloud logs read "resource.type=cloud_run_revision" --limit=50

# Filtrar logs por severidade
gcloud logs read "resource.type=cloud_run_revision AND severity>=ERROR" --limit=20

# Seguir logs em tempo real
gcloud logs tail "resource.type=cloud_run_revision"
```

### 2. Cloud Monitoring

1. Acesse [Cloud Monitoring](https://console.cloud.google.com/monitoring)
2. Configure alertas para:
   - CPU utilization > 80%
   - Memory utilization > 80%
   - Error rate > 5%
   - Response time > 2s

### 3. Métricas Importantes

- **Latência**: P50, P95, P99 response times
- **Throughput**: Requests per second
- **Error Rate**: Percentage of 4xx/5xx responses
- **Disponibilidade**: Uptime percentage

## 🔧 Troubleshooting

### Problemas Comuns

#### 1. Erro de Autenticação JWT
```bash
# Verificar se as chaves estão corretas no Secret Manager
gcloud secrets versions access latest --secret="jwt-private-key-dev"
```

#### 2. Erro de Conexão com Banco
```bash
# Verificar status da instância Cloud SQL
gcloud sql instances describe lep-postgres-dev

# Verificar logs do Cloud SQL
gcloud logs read "resource.type=cloudsql_database" --limit=20
```

#### 3. Container não Inicia
```bash
# Verificar logs do Cloud Run
gcloud run services logs read lep-backend-dev --region=us-central1

# Testar imagem localmente
docker run -p 8080:8080 IMAGEM_URL
```

#### 4. Problemas de Deploy
```bash
# Verificar build logs
gcloud builds list --limit=10
gcloud builds log BUILD_ID
```

### Comandos Úteis

```bash
# Redeploy forçado
gcloud run services replace service.yaml --region=us-central1

# Rollback para versão anterior
gcloud run services update-traffic lep-backend-dev --to-revisions=REVISION=100 --region=us-central1

# Escalar instâncias manualmente
gcloud run services update lep-backend-dev --min-instances=2 --region=us-central1

# Verificar recursos
gcloud run services describe lep-backend-dev --region=us-central1
```

## 💰 Custos Estimados

### Desenvolvimento (configuração mínima)
- **Cloud Run**: ~$5-15/mês
- **Cloud SQL (db-f1-micro)**: ~$7/mês
- **Secret Manager**: ~$1/mês
- **Artifact Registry**: ~$0.10/GB
- **Cloud Build**: 120 builds gratuitos/dia
- **Total estimado**: ~$15-25/mês

### Produção (configuração recomendada)
- **Cloud Run (regional)**: ~$30-100/mês
- **Cloud SQL (db-n1-standard-1, regional)**: ~$50-70/mês
- **Secret Manager**: ~$2/mês
- **Load Balancer**: ~$18/mês
- **Monitoring e Logging**: ~$10-20/mês
- **Total estimado**: ~$110-210/mês

### Otimização de Custos

1. **Use preemptible instances** para desenvolvimento
2. **Configure min-instances=0** para desenvolvimento
3. **Monitore logs** para evitar excessive logging costs
4. **Use Cloud SQL read replicas** apenas se necessário
5. **Configure lifecycle policies** para Artifact Registry

## 🔗 Links Úteis

- [Cloud Run Documentation](https://cloud.google.com/run/docs)
- [Cloud SQL Documentation](https://cloud.google.com/sql/docs)
- [Secret Manager Documentation](https://cloud.google.com/secret-manager/docs)
- [Cloud Build Documentation](https://cloud.google.com/build/docs)
- [Terraform GCP Provider](https://registry.terraform.io/providers/hashicorp/google/latest/docs)
- [GCP Pricing Calculator](https://cloud.google.com/products/calculator)

## 📞 Suporte

Para problemas específicos da infraestrutura:

1. Verifique os logs primeiro
2. Consulte este guia
3. Procure no [Stack Overflow](https://stackoverflow.com/questions/tagged/google-cloud-platform)
4. Abra um issue no repositório do projeto