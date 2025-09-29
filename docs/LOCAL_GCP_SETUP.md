# 🌐 Executar Local com Serviços GCP

Guia para rodar a aplicação **localmente** usando serviços **GCP** (banco e storage) do ambiente staging.

## 🔧 Configuração Necessária

### 1. **Autenticação GCP**

```bash
# Fazer login no GCP
gcloud auth login

# Configurar projeto
gcloud config set project leps-472702

# Configurar credenciais para aplicação
gcloud auth application-default login
```

### 2. **Instalar Cloud SQL Proxy**

```bash
# Windows
curl -o cloud-sql-proxy.exe https://storage.googleapis.com/cloud-sql-connectors/cloud-sql-proxy/v2.8.0/cloud-sql-proxy.x64.exe

# Ou via gcloud
gcloud components install cloud-sql-proxy
```

### 3. **Iniciar Cloud SQL Proxy**

```bash
# Em um terminal separado, manter rodando:
cloud-sql-proxy leps-472702:us-central1:leps-postgres-staging --port=5432

# Ou se já tiver PostgreSQL local na 5432, usar outra porta:
cloud-sql-proxy leps-472702:us-central1:leps-postgres-staging --port=5433
```

### 4. **Configurar .env para Proxy**

Se usar proxy na porta 5433, ajustar `.env`:
```bash
# Database Configuration (via Cloud SQL Proxy)
DB_HOST=localhost;
DB_PORT=5433;
DB_USER=lep_user_staging;
DB_PASS=sua_senha_staging;
DB_NAME=lep_database_staging;
DB_SSL_MODE=disable;
# Não precisa SSL via proxy
# Remover: INSTANCE_UNIX_SOCKET
```

### 5. **Verificar Credenciais de Banco**

```bash
# Verificar qual senha está configurada no Secret Manager
gcloud secrets versions access latest --secret="db-password-staging"
```

## 🚀 Comandos de Inicialização

### Terminal 1: Cloud SQL Proxy
```bash
cloud-sql-proxy leps-472702:us-central1:leps-postgres-staging --port=5432
```

### Terminal 2: Aplicação Go
```bash
go run main.go
```

## ✅ Verificações

### Banco de Dados
```bash
# Testar conexão diretamente
psql -h localhost -p 5432 -U lep_user_staging -d lep_database_staging
```

### Storage GCS
```bash
# Verificar acesso ao bucket
gsutil ls gs://leps-472702-lep-images-staging/

# Testar upload
curl -X POST http://localhost:8080/upload/product/image \
  -H "X-Lpe-Organization-Id: seu-org-id" \
  -H "X-Lpe-Project-Id: seu-project-id" \
  -F "image=@test-image.jpg"
```

### Health Check
```bash
curl http://localhost:8080/health
curl http://localhost:8080/ping
```

## 🔍 Troubleshooting

### Erro de Autenticação GCS
```bash
# Reconfigurar credenciais
gcloud auth application-default login
gcloud config set project leps-472702
```

### Erro de Conexão com Banco
```bash
# Verificar se proxy está rodando
netstat -an | findstr :5432

# Verificar instância Cloud SQL
gcloud sql instances describe leps-postgres-staging
```

### Erro de Permissões
```bash
# Verificar IAM da sua conta
gcloud projects get-iam-policy leps-472702 --flatten="bindings[].members" --filter="bindings.members:user:seu-email@gmail.com"
```

## 📊 Vantagens desta Configuração

✅ **Dados Reais**: Usa dados reais do staging
✅ **Debug Local**: Hot reload e debug facilitados
✅ **Storage Real**: Upload direto para GCS
✅ **Consistência**: Mesmo ambiente que staging
✅ **Performance**: Rede local + serviços cloud

## ⚠️ Cuidados

- ⚠️ **Não fazer alterações destrutivas** no banco staging
- ⚠️ **Não commitar credenciais** reais no .env
- ⚠️ **Monitorar custos** do Cloud SQL (sempre rodando)
- ⚠️ **Usar dados de teste** para uploads

## 🔄 Voltar para Local

Para voltar ao desenvolvimento totalmente local:

```bash
# Parar Cloud SQL Proxy
# Ctrl+C no terminal do proxy

# Restaurar .env original
cp .env.local .env  # se tiver backup

# Iniciar serviços locais
docker-compose up -d postgres mailhog
```