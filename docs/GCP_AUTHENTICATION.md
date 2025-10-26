# Autenticação Google Cloud Platform (GCP)

Este documento explica como configurar as credenciais do Google Cloud para usar o Google Cloud Storage (GCS) no LEP System.

## Problema Identificado

Ao tentar usar `STORAGE_TYPE=gcs`, a aplicação retorna:
```
google: could not find default credentials. See https://cloud.google.com/docs/authentication/external/set-up-adc for more information
```

Isso ocorre porque o cliente Go do Google Cloud não conseguiu encontrar as credenciais de autenticação.

## Soluções

### Solução 1: Application Default Credentials (ADC) - Recomendado para Desenvolvimento Local

Este é o método mais simples e recomendado para desenvolvimento local.

#### Passos:

1. **Faça login com sua conta Google:**
```bash
gcloud auth application-default login
```

2. Este comando abrirá um navegador para autenticação
3. As credenciais serão salvas automaticamente em:
   - **Windows**: `C:\Users\<seu-usuario>\AppData\Roaming\gcloud\application_default_credentials.json`
   - **macOS/Linux**: `~/.config/gcloud/application_default_credentials.json`

4. A aplicação encontrará essas credenciais automaticamente

#### Vantagens:
✅ Simples de usar
✅ Funciona imediatamente após o login
✅ Sem necessidade de arquivos de configuração

#### Desvantagens:
❌ Vinculado à sua conta pessoal (não adequado para produção)

---

### Solução 2: Service Account Key - Para Staging/Produção

Use uma Service Account para ambientes de staging e produção.

#### Passos:

1. **Crie uma Service Account no GCP Console:**
   - Vá para: https://console.cloud.google.com/iam-admin/serviceaccounts
   - Clique em "Create Service Account"
   - Nome: `lep-backend-staging`
   - Clique em "Create and Continue"

2. **Conceda permissões necessárias:**
   - Role: `Storage Admin` (para gerenciar GCS)
   - Role: `Cloud SQL Client` (se usar Cloud SQL)
   - Clique em "Continue"

3. **Crie uma chave JSON:**
   - Vá para a aba "Keys"
   - Clique em "Add Key" → "Create new key"
   - Selecione "JSON"
   - Salve o arquivo como `service-account-key.json` no raiz do projeto

4. **Configure a variável de ambiente:**

**No `.env.staging`:**
```bash
GOOGLE_APPLICATION_CREDENTIALS=./service-account-key.json
```

**Ou diretamente no terminal:**
```bash
export GOOGLE_APPLICATION_CREDENTIALS="./service-account-key.json"
go run main.go
```

#### Segurança:
⚠️ **IMPORTANTE**: Não commit o arquivo `service-account-key.json` no Git!
```bash
# Adicione ao .gitignore se não estiver lá
echo "service-account-key.json" >> .gitignore
```

---

### Solução 3: Cloud Run (Produção)

Quando deployar em Cloud Run, a autenticação é automática usando a Service Account associada ao Cloud Run.

#### Passos:

1. **Crie uma Service Account para Cloud Run:**
```bash
gcloud iam service-accounts create lep-backend-prod
```

2. **Conceda as permissões necessárias:**
```bash
gcloud projects add-iam-policy-binding leps-472702 \
  --member="serviceAccount:lep-backend-prod@leps-472702.iam.gserviceaccount.com" \
  --role="roles/storage.admin"

gcloud projects add-iam-policy-binding leps-472702 \
  --member="serviceAccount:lep-backend-prod@leps-472702.iam.gserviceaccount.com" \
  --role="roles/cloudsql.client"
```

3. **Deploy no Cloud Run com a Service Account:**
```bash
gcloud run deploy lep-backend \
  --service-account lep-backend-prod@leps-472702.iam.gserviceaccount.com \
  --image gcr.io/leps-472702/lep-backend:latest
```

#### Vantagens:
✅ Seguro - credenciais não salvas em arquivos
✅ Isolado - cada serviço tem sua própria conta
✅ Auditável - todas as ações são rastreadas

---

## Fluxo de Desenvolvimento Recomendado

### 1. Desenvolvimento Local
```bash
# Use ADC (Application Default Credentials)
gcloud auth application-default login

# Configure .env para usar local storage
STORAGE_TYPE=local
ENVIRONMENT=dev

# Execute a aplicação
go run main.go
```

### 2. Testes com GCS Local
```bash
# Configure para usar GCS
STORAGE_TYPE=gcs
BUCKET_NAME=leps-472702-lep-images-dev
ENVIRONMENT=dev

# Autentique com ADC ou Service Account
gcloud auth application-default login
# OU
export GOOGLE_APPLICATION_CREDENTIALS="./service-account-key.json"

# Execute a aplicação
go run main.go
```

### 3. Staging/Produção
```bash
# Use Service Account ou Cloud Run
# Confira as credenciais estão configuradas
echo $GOOGLE_APPLICATION_CREDENTIALS

# Build e deploy
docker build -t lep-backend:stage .
docker tag lep-backend:stage gcr.io/leps-472702/lep-backend:stage
docker push gcr.io/leps-472702/lep-backend:stage
```

---

## Verificação

Para verificar se as credenciais estão configuradas corretamente:

```bash
# Teste a autenticação
gcloud auth list

# Veja qual conta está ativa
gcloud config get-value account

# Teste o acesso ao bucket
gsutil ls gs://leps-472702-lep-images-stage/
```

Ou teste direto na aplicação - o log deve mostrar:
```
🔑 Usando credenciais do arquivo: ./service-account-key.json
✅ Conectado ao GCS bucket: leps-472702-lep-images-stage
```

---

## Troubleshooting

### Erro: "could not find default credentials"

**Verificação:**
1. Confirme que executou `gcloud auth application-default login`
2. Verifique o arquivo de credenciais existe:
   - Windows: `%APPDATA%\gcloud\application_default_credentials.json`
   - Linux/macOS: `~/.config/gcloud/application_default_credentials.json`

**Solução:**
```bash
gcloud auth application-default login
```

### Erro: "permission denied" ao fazer upload

**Verificação:**
1. A Service Account tem permissão `roles/storage.admin`?
2. O bucket existe e é acessível?

**Solução:**
```bash
# Verifique as permissões
gcloud projects get-iam-policy leps-472702 \
  --flatten="bindings[].members" \
  --filter="bindings.members:serviceAccount:*"

# Se necessário, adicione a role
gcloud projects add-iam-policy-binding leps-472702 \
  --member="serviceAccount:lep-backend-staging@leps-472702.iam.gserviceaccount.com" \
  --role="roles/storage.admin"
```

### Erro: "bucket not found"

**Solução:**
1. Confirme o nome do bucket está correto em `BUCKET_NAME`
2. Confirme o bucket foi criado no projeto correto
3. Verifique a localização do bucket

```bash
# Liste todos os buckets
gsutil ls

# Crie o bucket se não existir
gsutil mb gs://leps-472702-lep-images-stage
```

---

## Referências

- [Google Cloud Authentication](https://cloud.google.com/docs/authentication/external/set-up-adc)
- [Service Accounts](https://cloud.google.com/iam/docs/service-accounts)
- [Google Cloud Storage Go Client](https://pkg.go.dev/cloud.google.com/go/storage)
- [Cloud Run Authentication](https://cloud.google.com/run/docs/authenticating/service-to-service)
