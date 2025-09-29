# 🔧 Environment Setup Guide

Este guia explica como configurar as variáveis de ambiente para diferentes ambientes do LEP System.

## 📁 Arquivos de Ambiente

### `.env.example`
Template com todas as variáveis necessárias e valores de exemplo.

### `.env`
Configuração para desenvolvimento local.

### `.env.staging`
Configuração para ambiente de staging.

## 🚀 Quick Setup

### Desenvolvimento Local

1. **Copiar template**:
   ```bash
   cp .env.example .env
   ```

2. **O arquivo `.env` já está pré-configurado para desenvolvimento local com**:
   - PostgreSQL local (Docker Compose)
   - Storage local
   - SMTP via MailHog
   - Chaves JWT simples para teste
   - Bucket configurações para desenvolvimento

3. **Iniciar ambiente**:
   ```bash
   docker-compose up -d postgres mailhog
   go run main.go
   ```

### Staging

1. **Copiar template**:
   ```bash
   cp .env.staging .env
   ```

2. **Atualizar variáveis sensíveis**:
   ```bash
   # Database
   DB_PASS=sua_senha_staging

   # JWT Keys
   JWT_SECRET_PRIVATE_KEY=sua_chave_privada
   JWT_SECRET_PUBLIC_KEY=sua_chave_publica

   # SMTP
   SMTP_USERNAME=seu_email@gmail.com
   SMTP_PASSWORD=sua_senha_app
   ```

## 🔧 Novas Variáveis BUCKET_*

### `BUCKET_NAME`
- **Função**: Nome do bucket GCS (substitui `STORAGE_BUCKET_NAME`)
- **Local**: `lep-local-bucket`
- **Dev**: `leps-472702-lep-images-dev`
- **Staging**: `leps-472702-lep-images-staging`
- **Prod**: `leps-472702-lep-images-prod`

### `BUCKET_CACHE_CONTROL`
- **Função**: Header de cache para arquivos uploadados
- **Local**: `public, max-age=3600` (1 hora)
- **Dev**: `public, max-age=3600` (1 hora)
- **Staging**: `public, max-age=7200` (2 horas)
- **Prod**: `public, max-age=86400` (24 horas)

### `BUCKET_TIMEOUT`
- **Função**: Timeout em segundos para operações GCS
- **Local**: `30` segundos
- **Dev**: `30` segundos
- **Staging**: `60` segundos
- **Prod**: `120` segundos

## 🔒 Segurança

### Arquivos Protegidos
O `.gitignore` protege automaticamente:
- `.env`
- `.env.*`
- `*.pem`
- `*.key`
- `secrets/`

### Variáveis Sensíveis
Nunca commite:
- Senhas de banco
- Chaves JWT reais
- Credenciais SMTP/Twilio
- Tokens de API

## 📋 Checklist de Setup

### Desenvolvimento Local
- [ ] Arquivo `.env` criado
- [ ] PostgreSQL rodando (Docker Compose)
- [ ] MailHog rodando para SMTP
- [ ] Aplicação conecta ao banco
- [ ] Upload de imagens funciona

### Staging
- [ ] Arquivo `.env.staging` configurado
- [ ] Chaves JWT reais configuradas
- [ ] Credenciais SMTP configuradas
- [ ] Bucket GCS criado
- [ ] Deploy no Cloud Run funcionando

### Produção
- [ ] Secrets no Secret Manager
- [ ] Variáveis de ambiente no Cloud Run
- [ ] Backup do banco configurado
- [ ] Monitoramento ativo

## 🛠️ Troubleshooting

### Storage Local não funciona
```bash
# Verificar permissões
ls -la ./uploads/
mkdir -p ./uploads/products/
```

### GCS não conecta
```bash
# Verificar credenciais
gcloud auth list
gcloud config list project

# Verificar bucket existe
gsutil ls gs://leps-472702-lep-images-dev/
```

### Banco não conecta
```bash
# Local
docker-compose ps postgres

# GCP
gcloud sql instances list
```

## 📚 Documentação Adicional

- [Docker Compose Setup](../docker-compose.yml)
- [GCP Deploy Guide](./deployment/DEPLOYMENT_COMPLETE.md)
- [Storage Configuration](../CLAUDE.md#image-storage-system)