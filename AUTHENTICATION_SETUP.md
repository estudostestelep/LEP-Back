# Setup de Autenticação Google Cloud (GCP)

## Problema Atual

Você está recebendo dois erros de autenticação:

1. **Terraform**: `oauth2: "invalid_grant" "reauth related error (invalid_rapt)"`
2. **Docker Push**: `failed to authorize: 403 Forbidden`

Ambos indicam que a sessão GCP expirou e precisa ser reauthenticada.

---

## Solução Rápida (5 minutos)

### Passo 1: Reauthentique com GCP

```bash
# Limpe a autenticação anterior
gcloud auth revoke

# Faça login novamente
gcloud auth login
```

Este comando abrirá seu navegador para autenticação. Confirme o login.

### Passo 2: Configure o Application Default Credentials (ADC)

```bash
# Isso permitirá que Go acesse GCP automaticamente
gcloud auth application-default login
```

Novamente, seu navegador abrirá para autenticação.

### Passo 3: Configure Docker para acessar GCR

```bash
# Configure gcloud como helper de autenticação do Docker
gcloud auth configure-docker
```

Este comando configura Docker para usar `gcloud` como provider de autenticação.

### Passo 4: Tente novamente

```bash
# Reexecute o terraform apply
cd LEP-Back
terraform apply tfplan-stage

# Ou tente fazer push novamente
docker push gcr.io/leps-472702/lep-backend:stage
```

---

## Se Ainda Não Funcionar

### Para Terraform:

```bash
# Verifique sua autenticação atual
gcloud auth list

# Confirme que seu projeto está correto
gcloud config get-value project

# Se necessário, defina o projeto
gcloud config set project leps-472702

# Tente novamente
terraform apply tfplan-stage
```

### Para Docker Push:

```bash
# Verifique se Docker consegue falar com GCR
docker run hello-world

# Se GCR falhar, tente:
gcloud auth configure-docker

# Teste com uma imagem simples primeiro
docker tag hello-world gcr.io/leps-472702/hello-world:test
docker push gcr.io/leps-472702/hello-world:test
```

---

## Entender o Que Aconteceu

### RAPT Token
- GCP usa "Re-Authentication Challenge Token (RAPT)" para segurança
- Esse token expira após algum tempo de inatividade
- Solução: Fazer login novamente

### Docker Push 403
- Docker não tem permissão para acessar GCR
- Causa: Credentials não configuradas ou expiradas
- Solução: Executar `gcloud auth configure-docker`

---

## Referência de Comandos

```bash
# Verificar status de autenticação
gcloud auth list
gcloud config list

# Configurar credenciais
gcloud auth login                              # Login interativo
gcloud auth application-default login          # ADC para Go/aplicações
gcloud auth configure-docker                   # Docker para GCR

# Limpar credenciais
gcloud auth revoke
gcloud auth revoke --all

# Projeto
gcloud config set project leps-472702
gcloud config get-value project

# Permissões
gcloud projects get-iam-policy leps-472702
```

---

## Próximas Etapas

Após reauthenticar:

1. **Reexecute Terraform**:
   ```bash
   cd LEP-Back
   terraform apply tfplan-stage
   ```

2. **Faça Push da Imagem Docker**:
   ```bash
   docker push gcr.io/leps-472702/lep-backend:stage
   ```

3. **Deploy no Cloud Run** (opcional, se usando):
   ```bash
   gcloud run deploy lep-backend \
     --image gcr.io/leps-472702/lep-backend:stage \
     --region us-central1
   ```

---

## Dicas

- ✅ Mantenha sua sessão gcloud ativa executando `gcloud auth login` regularmente
- ✅ Use `gcloud auth application-default login` para Go/aplicações
- ✅ Use `gcloud auth configure-docker` para Docker
- ✅ Sempre confirme o projeto correto: `gcloud config get-value project`
- ⚠️ Não compartilhe chaves de serviço em públicos
- ⚠️ Use `gcloud auth revoke` quando terminar em computadores compartilhados

