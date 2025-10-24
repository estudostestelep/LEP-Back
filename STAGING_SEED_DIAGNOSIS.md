# 🔍 Diagnóstico: Seed Fattoria no Staging (GCP)

**Situação**: Você rodou o seed em staging, mas não encontrou dados no banco GCP.

---

## 📋 Checklist de Diagnóstico

### 1. Verificar Conectividade ao Cloud SQL

```bash
# Teste 1: Cloud SQL Proxy rodando?
ps aux | grep cloud_sql_proxy
# Esperado: processo ativo com -instances=leps-472702:us-central1:leps-postgres-staging

# Teste 2: Tentar conectar ao banco
gcloud sql connect leps-postgres-staging --user=lep_user

# Teste 3: Verificar se banco existe
gcloud sql databases list --instance=leps-postgres-staging
# Esperado: lep_database deve estar listado
```

### 2. Verificar Autenticação GCP

```bash
# Teste 1: Estar autenticado?
gcloud auth list
# Esperado: conta ativa com @ (email)

# Teste 2: Projeto correto?
gcloud config get-value project
# Esperado: leps-472702

# Teste 3: Ter permissões?
gcloud projects get-iam-policy leps-472702 --flatten="bindings[].members" --filter="bindings.members:$(gcloud config get-value account)"
# Esperado: deve ter roles como roles/cloudsql.client
```

### 3. Verificar Banco de Dados

```bash
# Conectar ao Cloud SQL
gcloud sql connect leps-postgres-staging --user=lep_user

-- Dentro do psql:

-- Teste 1: Banco existe?
\l
-- Esperado: lep_database na lista

-- Teste 2: Tabelas existem?
\dt
-- Esperado: tabelas como organizations, products, users

-- Teste 3: Dados Fattoria?
SELECT COUNT(*) FROM organizations WHERE name = 'Fattoria Pizzeria';
-- Se retornar 0: Seed não rodou

-- Teste 4: Qualquer produto Fattoria?
SELECT COUNT(*) FROM products WHERE name LIKE '%Pizza%' OR name LIKE '%Caju%';
-- Se retornar 0: Sem dados

-- Teste 5: Ver todo conteúdo
SELECT
  (SELECT COUNT(*) FROM organizations) as orgs,
  (SELECT COUNT(*) FROM products) as products,
  (SELECT COUNT(*) FROM users) as users,
  (SELECT COUNT(*) FROM tables) as tables;
```

---

## 🔧 Possíveis Causas e Soluções

### Causa 1: Cloud SQL Proxy Não Está Rodando

**Sintoma**: `ERROR: could not connect to server`

**Solução**:
```bash
# Iniciar Cloud SQL Proxy
cloud_sql_proxy -instances=leps-472702:us-central1:leps-postgres-staging=unix:/cloudsql/leps-472702:us-central1:leps-postgres-staging &

# Ou para TCP (mais simples):
cloud_sql_proxy -instances=leps-472702:us-central1:leps-postgres-staging=tcp:5432 &

# Aguardar inicialização
sleep 5

# Tentar conectar novamente
gcloud sql connect leps-postgres-staging --user=lep_user
```

### Causa 2: Variáveis de Ambiente Não Foram Carregadas

**Sintoma**: Seed tentou conectar mas usou valores padrão/errados

**Solução**:
```bash
# Verificar se variáveis estão carregadas
echo "DB_USER: $DB_USER"
echo "DB_PASS: $DB_PASS"
echo "INSTANCE_UNIX_SOCKET: $INSTANCE_UNIX_SOCKET"

# Se estiverem vazias, carregar manualmente:
export $(grep -v '^#' .env.staging | xargs)

# Verificar novamente
echo "DB_USER: $DB_USER"  # Deve mostrar: lep_user
```

### Causa 3: Seed Rodou Mas Em Banco Diferente

**Sintoma**: Dados em DEV (Docker) mas não em Staging

**Solução**:
```bash
# Verificar qual banco foi populado
# DEV (Docker local)
docker exec lep-postgres psql -U lep_user -d lep_database -c \
  "SELECT COUNT(*) FROM products WHERE name LIKE '%Pizza%';"

# STAGE (Cloud SQL)
gcloud sql connect leps-postgres-staging --user=lep_user --quiet -c \
  "SELECT COUNT(*) FROM products WHERE name LIKE '%Pizza%';"
```

### Causa 4: Seed Falhou Silenciosamente

**Sintoma**: Nenhuma mensagem de erro, mas sem dados

**Solução**:
```bash
# Rodar seed com verbose máximo
ENVIRONMENT=staging go run cmd/seed/main.go \
  --restaurant=fattoria \
  --environment=staging \
  --verbose 2>&1 | tee seed_debug.log

# Verificar arquivo de log
cat seed_debug.log | grep -E "(ERROR|FAIL|error|fail)"
```

### Causa 5: Permissões Insuficientes do Usuário

**Sintoma**: "Permission denied" ou similar

**Solução**:
```bash
# Conectar como admin (postgres)
gcloud sql connect leps-postgres-staging --user=postgres

-- Dar permissões completas ao lep_user:
ALTER USER lep_user WITH SUPERUSER;

-- Ou permissões específicas:
GRANT ALL ON DATABASE lep_database TO lep_user;
GRANT ALL ON SCHEMA public TO lep_user;
GRANT ALL ON ALL TABLES IN SCHEMA public TO lep_user;
```

---

## 🔄 Passos de Recuperação

Se o seed falhou, siga estes passos:

### Passo 1: Limpar Dados Antigos (se houver)

```bash
# Conectar ao banco
gcloud sql connect leps-postgres-staging --user=lep_user

-- Limpar dados (manter schema)
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
```

### Passo 2: Garantir Cloud SQL Proxy Rodando

```bash
# Verificar
ps aux | grep cloud_sql_proxy

# Se não estiver, iniciar:
cloud_sql_proxy -instances=leps-472702:us-central1:leps-postgres-staging=tcp:5432 &
sleep 5
```

### Passo 3: Carregar Variáveis Corretamente

```bash
# Carregar .env.staging
export $(grep -v '^#' .env.staging | xargs)

# Verificar
echo "ENVIRONMENT=$ENVIRONMENT"  # Deve ser: staging
echo "DB_USER=$DB_USER"          # Deve ser: lep_user
```

### Passo 4: Rodar Seed com Debug

```bash
# Opção A: Via master script (mais fácil)
bash ./scripts/master-interactive.sh --seed-fattoria-stage

# Opção B: Via script direto
bash ./scripts/run_seed_staging.sh --verbose

# Opção C: Via Go direto (mais controle)
ENVIRONMENT=staging go run cmd/seed/main.go \
  --restaurant=fattoria \
  --environment=staging \
  --verbose 2>&1 | tee seed_output.log
```

### Passo 5: Validar Dados

```bash
# Conectar
gcloud sql connect leps-postgres-staging --user=lep_user

-- Contar produtos Fattoria
SELECT COUNT(*) as produto_count FROM products
WHERE organization_id IN (
  SELECT id FROM organizations WHERE name = 'Fattoria Pizzeria'
);

-- Esperado: 9 produtos
```

---

## 📊 Queries de Diagnóstico Rápidas

Salve estas queries para testar rapidamente:

```sql
-- Query 1: Verificar se Fattoria existe
SELECT id, name, email FROM organizations WHERE name = 'Fattoria Pizzeria' LIMIT 1;

-- Query 2: Contar todos os dados
SELECT
  'Organizations' as entity, COUNT(*) as count FROM organizations
UNION ALL
SELECT 'Projects', COUNT(*) FROM projects
UNION ALL
SELECT 'Users', COUNT(*) FROM users
UNION ALL
SELECT 'Products', COUNT(*) FROM products
UNION ALL
SELECT 'Tables', COUNT(*) FROM tables
UNION ALL
SELECT 'Categories', COUNT(*) FROM categories
UNION ALL
SELECT 'Customers', COUNT(*) FROM customers;

-- Query 3: Listar produtos Fattoria
SELECT name, price_normal
FROM products
WHERE organization_id = (SELECT id FROM organizations WHERE name = 'Fattoria Pizzeria' LIMIT 1)
ORDER BY name;

-- Query 4: Verificar audit logs (últimas operações)
SELECT entity, action, created_at
FROM audit_logs
ORDER BY created_at DESC LIMIT 20;

-- Query 5: Status geral do banco
SELECT
  schemaname,
  tablename,
  pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables
WHERE schemaname != 'pg_catalog'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

---

## 🚀 Próximas Ações

1. **Imediatamente**:
   - [ ] Verificar se Cloud SQL Proxy está rodando
   - [ ] Conectar ao banco e validar estrutura
   - [ ] Rodar as queries de diagnóstico

2. **Se Fattoria Não Existe**:
   - [ ] Garantir .env.staging está correto (sem `;`)
   - [ ] Carregar variáveis: `export $(grep -v '^#' .env.staging | xargs)`
   - [ ] Rodar seed novamente: `bash ./scripts/master-interactive.sh --seed-fattoria-stage`

3. **Se Seed Falhar**:
   - [ ] Verificar logs: `tail -100 seed_output.log`
   - [ ] Procurar por erros de conexão
   - [ ] Validar permissões no Cloud SQL

4. **Se Dados Aparecerem**:
   - [ ] Validar 9 produtos foram inseridos
   - [ ] Testar login: admin@fattoria.com.br / password
   - [ ] Deploy da API e teste

---

## 📞 Contato & Suporte

Se o diagnóstico não resolver:

1. Envie os arquivos de log:
   - `seed_output.log` (com --verbose)
   - Saída do `gcloud sql connect`

2. Compartilhe o resultado das queries de diagnóstico

3. Confirme:
   - [ ] Cloud SQL Proxy rodando?
   - [ ] Autenticado no GCP (`gcloud auth list`)?
   - [ ] Projeto correto (`gcloud config get-value project`)?
   - [ ] `.env.staging` sem pontos-e-vírgula?

---

**Data**: 24 de Outubro, 2025
**Status**: Guia de Diagnóstico Completo
**Próximo**: Execute o diagnóstico e relate os resultados
