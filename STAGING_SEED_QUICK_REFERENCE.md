# Quick Reference: Staging Seed Fattoria

## Execução Rápida (3 formas)

### 1️⃣ Mais Fácil - Menu Interativo
```bash
bash ./scripts/master-interactive.sh
# Menu → 3 → 7 → y
```

### 2️⃣ Rápida - Batch Mode
```bash
bash ./scripts/master-interactive.sh --seed-fattoria-stage
```

### 3️⃣ Direto - Script Puro
```bash
bash ./scripts/run_seed_staging.sh --clear-first --verbose
```

---

## Pré-requisitos (verificar antes)

```bash
# 1. Cloud SQL Proxy rodando?
ps aux | grep cloud_sql_proxy

# 2. GCP autenticado?
gcloud auth list
gcloud config set project leps-472702

# 3. Variáveis carregando?
source .env.staging && echo $BUCKET_CACHE_CONTROL
```

---

## Validação Após Execução

```bash
# Conectar ao Cloud SQL
gcloud sql connect leps-postgres-staging --user=lep_user

# Queries de validação
SELECT COUNT(*) FROM organizations WHERE name = 'Fattoria Pizzeria';
-- Esperado: 1

SELECT COUNT(*) FROM products 
WHERE organization_id = (SELECT id FROM organizations WHERE name = 'Fattoria Pizzeria');
-- Esperado: 9

SELECT COUNT(*) FROM tables 
WHERE organization_id = (SELECT id FROM organizations WHERE name = 'Fattoria Pizzeria');
-- Esperado: 3

SELECT email FROM users WHERE email = 'admin@fattoria.com.br';
-- Esperado: admin@fattoria.com.br
```

---

## Troubleshooting Rápido

| Erro | Solução |
|------|---------|
| Connection refused | Iniciar Cloud SQL Proxy |
| Database does not exist | `gcloud sql databases create lep_database --instance=leps-postgres-staging` |
| max-age=7200: command not found | Script já foi corrigido ✅ |
| Sem output / seed não roda | Verificar se Cobra output fix está em cmd/seed/main.go ✅ |
| Dados não aparecem | Verificar se `.env.staging` tem as credenciais corretas |

---

## Dados Fattoria (o que será inserido)

- **Organização**: Fattoria Pizzeria
- **Admin**: admin@fattoria.com.br / password
- **Produtos** (9): Crostini, Margherita, Marinara, Parma, Vegana, Suco de Caju, Heineken, Baden Baden IPA, Sônia e Zé
- **Mesas** (3): 3 mesas com capacidades diferentes

---

## Documentação Completa

- `STAGING_SEED_FIX_FINAL.md` - Explicação detalhada das correções
- `RUN_SEED_STAGING_INSTRUCTIONS.md` - Guia passo-a-passo
- `STAGING_SEED_DIAGNOSIS.md` - Troubleshooting avançado
- `.env.staging.example` - Template de variáveis

