# 📋 Instruções para Executar Seed em Staging

## ⚠️ Situação Atual

O script `run_seed_staging.sh` foi criado, mas **não foi executado em ambiente real** porque:

1. **Cloud SQL Proxy não está rodando** neste ambiente local
2. **Sem acesso direto ao GCP** (leps-472702:us-central1:leps-postgres-staging)
3. **Unix socket não disponível** (`/cloudsql/...`)

---

## 🚀 Como Executar em Staging (Instruções para Você)

### Pré-requisito 1: Cloud SQL Proxy em Execução

No seu ambiente **com acesso GCP** (seu laptop/servidor), execute:

```bash
# Para Unix socket (Linux/Mac)
cloud_sql_proxy \
  -instances=leps-472702:us-central1:leps-postgres-staging=unix:/cloudsql/leps-472702:us-central1:leps-postgres-staging

# OU para TCP (qualquer SO)
cloud_sql_proxy \
  -instances=leps-472702:us-central1:leps-postgres-staging=tcp:5432
```

**Deixar este processo rodando em um terminal**

### Pré-requisito 2: Autenticação GCP

```bash
# Verificar se está autenticado
gcloud auth list

# Se não estiver:
gcloud auth login
gcloud config set project leps-472702
```

### Pré-requisito 3: Clonar/Atualizar Repositório

```bash
# Atualizar código
cd /caminho/para/LEP-Back
git pull origin dev
```

---

## ✅ Executar o Seed

### Opção 1: Usar o Script Bash (Recomendado)

```bash
# Navegar para o projeto
cd /caminho/para/LEP-Back

# Executar com limpeza de dados antigos
bash scripts/run_seed_staging.sh --clear-first --verbose

# OU sem limpeza (adiciona aos dados existentes)
bash scripts/run_seed_staging.sh --verbose
```

**Saída esperada**:
```
================================================
🌱 LEP Database Seeder - STAGING ENVIRONMENT
================================================
✅ Found .env.staging configuration
✅ Environment variables loaded

Configuration:
  Environment: staging
  Database: lep_database
  Database User: lep_user
  Restaurant: fattoria
  Clear First: true
  Verbose: true

🔍 Validating Database Connection
⚠️  Using Cloud SQL Unix Socket: /cloudsql/leps-472702:us-central1:leps-postgres-staging
⚠️  Ensure Cloud SQL Proxy is running or you have direct access

🌱 Running Seed Execution
================================================

... [progresso do seed] ...

✅ Seed execution completed successfully!

📊 Next Steps
1. Verify data in staging database
2. Run: gcloud sql connect <instance-name> --user=lep_user
3. Test API endpoint: https://staging-api.lep.example.com/health

✅ Staging seed ready!
================================================
```

### Opção 2: Comando Direto Go

```bash
# Carregar variáveis e rodar
cd /caminho/para/LEP-Back
export $(grep -v '^#' .env.staging | xargs)
go run ./cmd/seed/ --restaurant=fattoria --environment=staging --clear-first --verbose
```

### Opção 3: Docker (Se Preferir)

```bash
cd /caminho/para/LEP-Back

# Build
docker build -t lep-seed:staging -f Dockerfile.dev .

# Run
docker run --rm \
  --env-file .env.staging \
  -e DB_HOST=host.docker.internal \
  lep-seed:staging \
  go run cmd/seed/main.go --restaurant=fattoria --clear-first --verbose
```

---

## 📊 Validar Dados Inseridos

Após o seed completar com sucesso:

### Conectar ao Cloud SQL

```bash
gcloud sql connect leps-postgres-staging \
  --user=lep_user \
  --quiet
```

### Executar Validações

```sql
-- 1. Verificar organização foi criada
SELECT COUNT(*) as fattoria_count FROM organizations
WHERE name = 'Fattoria Pizzeria';
-- Esperado: 1

-- 2. Contar produtos (esperado: 9)
SELECT COUNT(*) as product_count FROM products
WHERE organization_id IN (
  SELECT id FROM organizations WHERE name = 'Fattoria Pizzeria'
);
-- Esperado: 9

-- 3. Listar produtos
SELECT name, price_normal FROM products
WHERE organization_id IN (
  SELECT id FROM organizations WHERE name = 'Fattoria Pizzeria'
)
ORDER BY name;

-- 4. Verificar admin user
SELECT email, name FROM users
WHERE email = 'admin@fattoria.com.br';

-- 5. Verificar mesas
SELECT COUNT(*) as table_count FROM tables
WHERE organization_id IN (
  SELECT id FROM organizations WHERE name = 'Fattoria Pizzeria'
);
-- Esperado: 3

-- 6. Validar categorias
SELECT COUNT(*) as category_count FROM categories
WHERE organization_id IN (
  SELECT id FROM organizations WHERE name = 'Fattoria Pizzeria'
);
```

---

## 🐛 Se Algo der Errado

### Problema: "Connection refused"

```bash
# Verificar se Cloud SQL Proxy está rodando
ps aux | grep cloud_sql_proxy

# Se não estiver, inicie novamente:
cloud_sql_proxy -instances=leps-472702:us-central1:leps-postgres-staging=unix:/cloudsql/...
```

### Problema: "Permission denied" para Unix socket

```bash
# Verificar permissões
ls -la /cloudsql/

# Corrigir se necessário:
sudo chown -R $(whoami):$(whoami) /cloudsql
chmod 755 /cloudsql
```

### Problema: "Authentication failed"

```bash
# Verificar credenciais em .env.staging
cat .env.staging | grep "DB_"

# Testar conexão manualmente:
psql -h /var/run/cloudsql/leps-472702:us-central1:leps-postgres-staging \
     -U lep_user \
     -d lep_database \
     -c "SELECT 1;"
```

### Problema: Silent output (sem feedback)

```bash
# Adicionar redirecionamento de output:
bash scripts/run_seed_staging.sh --clear-first --verbose 2>&1 | tee seed_staging.log

# Depois verificar o arquivo:
cat seed_staging.log
```

### Problema: "database does not exist"

```bash
# Criar banco antes:
gcloud sql databases create lep_database \
  --instance=leps-postgres-staging
```

---

## 📋 Checklist de Sucesso

Após rodar o seed, validar:

- [ ] Script executou sem erros (exit code 0)
- [ ] 1 organização "Fattoria Pizzeria" foi criada
- [ ] 9 produtos foram inseridos
- [ ] 1 usuário admin criado (admin@fattoria.com.br)
- [ ] 3 mesas criadas
- [ ] Categorias carregadas
- [ ] Audit logs mostram operações

---

## 🔄 Próximas Ações Após Seed

```bash
# 1. Deploy da API para Cloud Run
gcloud run deploy lep-api-staging \
  --source . \
  --region us-central1 \
  --platform managed \
  --env-vars-file .env.staging

# 2. Testar health endpoint
curl https://staging-api.lep.example.com/health

# 3. Testar login
curl -X POST https://staging-api.lep.example.com/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@fattoria.com.br",
    "password": "password"
  }'

# 4. Listar produtos via API
curl https://staging-api.lep.example.com/product \
  -H "Authorization: Bearer <seu_token>" \
  -H "X-Lpe-Organization-Id: <fattoria_org_id>" \
  -H "X-Lpe-Project-Id: <fattoria_project_id>"
```

---

## 📞 Referências

- **Script Principal**: `scripts/run_seed_staging.sh`
- **Documentação Completa**: `docs/seed/SEED_STAGING.md`
- **Quick Guide**: `STAGING_SEED_GUIDE.md`
- **Troubleshooting**: `SEED_EXECUTION_REPORT.md`
- **Arquitetura**: `docs/seed/SEED_ARCHITECTURE.md`

---

## 📝 Status Final

✅ **Scripts Criados e Documentados**
- `scripts/run_seed_staging.sh` - Pronto para usar
- `tools/scripts/seed/run_seed_staging.sh` - Cópia espelhada

✅ **Documentação Completa**
- Instruções passo a passo
- Troubleshooting detalhado
- Validação de dados

⏳ **Aguardando**
- Sua execução em ambiente com acesso GCP

---

**Instruções Criadas em**: 24 de Outubro, 2025
**Versão**: 1.0
**Status**: Pronto para Execução
