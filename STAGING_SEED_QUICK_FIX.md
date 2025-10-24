# ⚡ Quick Fix: Dados não aparecem no Cloud SQL (GCP)

Você rodou o seed em staging, mas não encontrou dados no banco GCP. Aqui está como resolver:

---

## 🚀 Solução Rápida (5 Passos)

### Passo 1: Verificar Cloud SQL Proxy

```bash
# Está rodando?
ps aux | grep cloud_sql_proxy

# Se NÃO estiver, inicie:
cloud_sql_proxy -instances=leps-472702:us-central1:leps-postgres-staging=tcp:5432 &

# Aguarde
sleep 5
```

### Passo 2: Verificar Autenticação GCP

```bash
# Autenticado?
gcloud auth list

# Projeto correto?
gcloud config get-value project
# Esperado: leps-472702

# Se não for, configure:
gcloud config set project leps-472702
```

### Passo 3: Testar Conectividade ao Banco

```bash
# Conectar (a senha é: 123456)
gcloud sql connect leps-postgres-staging --user=lep_user

-- Verificar se banco existe:
\l
-- Esperado: lep_database na lista

-- Sair:
\q
```

### Passo 4: Verificar se Fattoria Existe

```bash
# Conectar
gcloud sql connect leps-postgres-staging --user=lep_user

-- Contar Fattoria:
SELECT COUNT(*) FROM organizations WHERE name = 'Fattoria Pizzeria';

-- Se retornar 0: Seed não rodou (veja Passo 5)
-- Se retornar 1: Seed rodou, mas produtos não foram inseridos (investigue)

\q
```

### Passo 5: Se Fattoria Não Existe, Rodar Seed Novamente

```bash
# Opção A: Via master script (RECOMENDADO)
bash ./scripts/master-interactive.sh --seed-fattoria-stage

# Opção B: Via script direto
bash ./scripts/run_seed_staging.sh --verbose

# Opção C: Via Go (com debug)
export $(grep -v '^#' .env.staging | xargs)
ENVIRONMENT=staging go run cmd/seed/main.go \
  --restaurant=fattoria \
  --environment=staging \
  --verbose
```

---

## ✅ Validar Sucesso

Após rodar seed, execute:

```bash
# Conectar ao banco
gcloud sql connect leps-postgres-staging --user=lep_user

-- Contar produtos Fattoria (esperado: 9)
SELECT COUNT(*) FROM products
WHERE organization_id IN (
  SELECT id FROM organizations WHERE name = 'Fattoria Pizzeria'
);

-- Listar produtos
SELECT name, price_normal FROM products
WHERE organization_id IN (
  SELECT id FROM organizations WHERE name = 'Fattoria Pizzeria'
)
ORDER BY name;

-- Verificar usuário admin
SELECT email FROM users WHERE email = 'admin@fattoria.com.br';

\q
```

**Esperado**:
- 9 produtos listados
- Usuário admin@fattoria.com.br existe

---

## 🔧 Se Ainda Não Funcionar

### Problema: Conexão Recusada

```bash
# Iniciar Cloud SQL Proxy em TCP (mais fácil)
cloud_sql_proxy -instances=leps-472702:us-central1:leps-postgres-staging=tcp:5432 &

# Aguardar
sleep 3

# Testar
psql -h localhost -U lep_user -d lep_database -c "SELECT 1;"
```

### Problema: Auth Failed

```bash
# Re-autenticar
gcloud auth login

# Aplicar credenciais padrão
gcloud auth application-default login
```

### Problema: Banco Não Existe

```bash
# Criar banco
gcloud sql databases create lep_database --instance=leps-postgres-staging

# Verificar
gcloud sql databases list --instance=leps-postgres-staging
```

### Problema: Permissão Negada

```bash
# Conectar como postgres (usuário admin)
gcloud sql connect leps-postgres-staging --user=postgres

-- Dar permissões
ALTER USER lep_user WITH SUPERUSER;
GRANT ALL ON DATABASE lep_database TO lep_user;

\q
```

---

## 📋 Checklist de Status

- [ ] Cloud SQL Proxy rodando (`ps aux | grep cloud_sql_proxy`)
- [ ] Autenticado no GCP (`gcloud auth list`)
- [ ] Projeto correto (`gcloud config get-value project`)
- [ ] `.env.staging` sem pontos-e-vírgula
- [ ] Conecta ao banco (`gcloud sql connect ...`)
- [ ] Fattoria existe no banco (`SELECT COUNT(*)...`)
- [ ] 9 produtos foram inseridos
- [ ] Usuário admin existe

---

## 🚨 Erro Mais Comum

**Erro**: `ERROR: could not connect to server`

**Causa**: Cloud SQL Proxy não está rodando

**Solução**:
```bash
# Verificar
ps aux | grep cloud_sql_proxy

# Se não houver resultado, inicie:
cloud_sql_proxy -instances=leps-472702:us-central1:leps-postgres-staging=tcp:5432 &

# Aguarde 5 segundos
sleep 5

# Tente novamente
gcloud sql connect leps-postgres-staging --user=lep_user
```

---

## 📞 Se Nada Funcionar

Procure no `STAGING_SEED_DIAGNOSIS.md` para diagnóstico detalhado, que contém:
- Causa 1: Cloud SQL Proxy não rodando
- Causa 2: Variáveis não carregadas
- Causa 3: Seed rodou em banco errado
- Causa 4: Seed falhou silenciosamente
- Causa 5: Permissões insuficientes

---

**Status**: 🎯 Quick Fix fornecido
**Próximo**: Execute os passos acima e reporte se funcionou
