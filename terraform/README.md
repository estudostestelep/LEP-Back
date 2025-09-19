# LEP System - Terraform Infrastructure

Este diretório contém a configuração do Terraform para provisionamento da infraestrutura do LEP System no Google Cloud Platform (GCP).

## 📋 Pré-requisitos

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install)
- Conta GCP com projeto configurado
- Permissões necessárias no projeto GCP

## 🔧 Configuração Inicial

### 1. Autenticação com GCP

```bash
# Login na sua conta Google
gcloud auth login

# Configurar projeto padrão
gcloud config set project YOUR_PROJECT_ID

# Configurar credenciais para Terraform
gcloud auth application-default login
```

### 2. Configuração de Variáveis

```bash
# Copie o arquivo de exemplo
cp terraform.tfvars.example terraform.tfvars

# Edite com seus valores
nano terraform.tfvars
```

### Variáveis Obrigatórias

- `project_id`: ID do seu projeto GCP
- `database_password`: Senha segura para o banco PostgreSQL
- `jwt_secret_private_key`: Chave secreta para JWT
- `jwt_secret_public_key`: Chave pública para JWT (pode ser a mesma para HS256)
- `container_image`: Imagem Docker da aplicação

## 🚀 Deploy

### 1. Inicialização

```bash
# Inicializar Terraform
terraform init
```

### 2. Planejar Mudanças

```bash
# Verificar o que será criado
terraform plan
```

### 3. Aplicar Infraestrutura

```bash
# Aplicar mudanças
terraform apply

# Confirmar com 'yes' quando solicitado
```

### 4. Verificar Deploy

```bash
# Ver outputs importantes
terraform output cloud_run_url
terraform output database_connection_name
```

## 🏗️ Recursos Criados

### Core Infrastructure
- **Cloud Run Service**: Aplicação containerizada com auto-scaling
- **Cloud SQL (PostgreSQL)**: Banco de dados gerenciado
- **Service Account**: Para comunicação segura entre serviços

### Security
- **Secret Manager**: Armazenamento seguro de senhas e chaves JWT
- **IAM Roles**: Permissões mínimas necessárias
- **Private IP**: Banco de dados isolado na rede interna

### Monitoring & Backup
- **Cloud SQL Backups**: Backup automático com point-in-time recovery
- **Health Checks**: Probes de startup e liveness
- **Logging**: Logs automáticos de conexões e operações

## 🔒 Segurança

### Melhores Práticas Implementadas

- ✅ Banco de dados sem IP público (apenas private IP)
- ✅ SSL obrigatório para conexões com banco
- ✅ Secrets gerenciados pelo Secret Manager
- ✅ Service Account com permissões mínimas
- ✅ Backup automático habilitado
- ✅ Logs de auditoria habilitados

### Configurações de Produção

Para ambiente de produção, ajuste:

```hcl
# terraform.tfvars
environment = "prod"
database_tier = "db-custom-2-4096"  # Maior performance
max_instances = 10                   # Maior escala
deletion_protection = true           # Proteção contra exclusão
```

## 🌍 Ambientes

### Desenvolvimento
```bash
terraform workspace new dev
terraform apply -var-file="dev.tfvars"
```

### Produção
```bash
terraform workspace new prod
terraform apply -var-file="prod.tfvars"
```

## 📊 Monitoramento

### URLs Importantes

Após o deploy, acesse:

- **Aplicação**: `terraform output cloud_run_url`
- **Cloud Run Console**: [Console GCP - Cloud Run](https://console.cloud.google.com/run)
- **Cloud SQL Console**: [Console GCP - Cloud SQL](https://console.cloud.google.com/sql)
- **Secret Manager**: [Console GCP - Secret Manager](https://console.cloud.google.com/security/secret-manager)

### Health Check

```bash
# URL da aplicação
APP_URL=$(terraform output -raw cloud_run_url)

# Teste de health check
curl $APP_URL/health

# Teste de login (se endpoint público)
curl -X POST $APP_URL/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'
```

## 🧹 Limpeza

### Destruir Infraestrutura

```bash
# CUIDADO: Isso remove TODOS os recursos
terraform destroy
```

### Remoção Segura

Para produção, considere:

1. Fazer backup manual do banco
2. Exportar secrets importantes
3. Remover `deletion_protection = true`
4. Executar `terraform destroy`

## 🔧 Troubleshooting

### Problemas Comuns

**1. Erro de API não habilitada:**
```bash
# Habilitar APIs manualmente
gcloud services enable run.googleapis.com sqladmin.googleapis.com secretmanager.googleapis.com
```

**2. Erro de permissões:**
```bash
# Verificar permissões do usuário
gcloud projects get-iam-policy YOUR_PROJECT_ID
```

**3. Falha na conexão com banco:**
- Verificar se a aplicação está usando Unix socket correto
- Confirmar se service account tem role `cloudsql.client`

### Logs

```bash
# Ver logs do Cloud Run
gcloud logging read "resource.type=cloud_run_revision" --limit=50

# Ver logs do Cloud SQL
gcloud logging read "resource.type=cloudsql_database" --limit=50
```

## 📚 Recursos Adicionais

- [Terraform Google Provider](https://registry.terraform.io/providers/hashicorp/google/latest/docs)
- [Cloud Run Documentation](https://cloud.google.com/run/docs)
- [Cloud SQL for PostgreSQL](https://cloud.google.com/sql/docs/postgres)
- [Secret Manager](https://cloud.google.com/secret-manager/docs)

## 🤝 Contribuição

1. Faça alterações em uma branch separada
2. Teste em ambiente de desenvolvimento primeiro
3. Valide com `terraform plan`
4. Abra PR com documentação das mudanças