# Próximos Passos: Staging Seed Fattoria

**Status Atual**: ✅ Scripts corrigidos e prontos para execução

---

## O que foi feito ✅

- [x] Identificado erro bash: `export $(... | xargs)` com valores complexos
- [x] Corrigido método de carregamento: `set -a; source; set +a`
- [x] Testado carregamento de variáveis
- [x] Sincronizados ambos os scripts
- [x] Criada documentação completa
- [x] Realizados 3 git commits

---

## Próximas Ações (Você faz isso)

### 1. Verificar Pré-requisitos

```bash
# Cloud SQL Proxy rodando?
ps aux | grep cloud_sql_proxy
# Se não tiver, iniciar em outro terminal:
cloud_sql_proxy -instances=leps-472702:us-central1:leps-postgres-staging=unix:/cloudsql/leps-472702:us-central1:leps-postgres-staging

# GCP autenticado?
gcloud auth list
gcloud config set project leps-472702

# Variáveis carregando?
cd LEP-Back
source .env.staging
echo $BUCKET_CACHE_CONTROL
# Deve mostrar: public, max-age=7200
```

### 2. Executar o Seed (escolha uma opção)

#### Opção A: Menu Interativo (Recomendado)
```bash
bash ./scripts/master-interactive.sh
# Selecione: 3 (Database & Seeding)
# Selecione: 7 (🍕 Seed Fattoria STAGE)
# Confirme: y
```

#### Opção B: Batch Mode (Rápido)
```bash
bash ./scripts/master-interactive.sh --seed-fattoria-stage
```

#### Opção C: Script Direto
```bash
bash ./scripts/run_seed_staging.sh --clear-first --verbose
```

### 3. Validar Dados no Banco

Após o seed completar (esperado: 2-5 minutos):

```bash
# Conectar ao Cloud SQL
gcloud sql connect leps-postgres-staging --user=lep_user

# Executar as queries abaixo
```

**Query 1: Verificar Organização**
```sql
SELECT id, name FROM organizations 
WHERE name = 'Fattoria Pizzeria' AND deleted_at IS NULL;
-- Esperado: 1 registro
```

**Query 2: Contar Produtos**
```sql
SELECT COUNT(*) as produto_count FROM products 
WHERE organization_id = (
  SELECT id FROM organizations WHERE name = 'Fattoria Pizzeria'
) AND deleted_at IS NULL;
-- Esperado: 9
```

**Query 3: Listar Produtos**
```sql
SELECT name, price_normal FROM products 
WHERE organization_id = (
  SELECT id FROM organizations WHERE name = 'Fattoria Pizzeria'
) AND deleted_at IS NULL
ORDER BY name;
-- Esperado: 9 produtos em ordem alfabética
```

**Query 4: Verificar Admin**
```sql
SELECT email, name FROM users 
WHERE email = 'admin@fattoria.com.br' AND deleted_at IS NULL;
-- Esperado: 1 usuário com esse email
```

**Query 5: Contar Mesas**
```sql
SELECT COUNT(*) as mesa_count FROM tables 
WHERE organization_id = (
  SELECT id FROM organizations WHERE name = 'Fattoria Pizzeria'
) AND deleted_at IS NULL;
-- Esperado: 3
```

### 4. Verificar Sucesso

Após as validações acima:

- [ ] Exit code foi 0 (sem erros)
- [ ] Mensagem "Seed execution completed successfully!" apareceu
- [ ] 1 organização "Fattoria Pizzeria" criada
- [ ] 9 produtos inseridos
- [ ] Admin user admin@fattoria.com.br criado
- [ ] 3 mesas criadas
- [ ] Todas as queries retornaram resultados esperados

---

## Se Algo Não Funcionar

### Erro: "Connection refused"
```bash
# Cloud SQL Proxy não está rodando. Solução:
cloud_sql_proxy -instances=leps-472702:us-central1:leps-postgres-staging=unix:/cloudsql/leps-472702:us-central1:leps-postgres-staging
```

### Erro: "max-age=7200: command not found"
Este erro foi corrigido. Se ainda aparecer:
```bash
# Verificar se script foi atualizado
git pull origin dev
bash ./scripts/run_seed_staging.sh --verbose
```

### Erro: "Database does not exist"
```bash
# Criar database antes
gcloud sql databases create lep_database \
  --instance=leps-postgres-staging
```

### Sem erros, mas sem dados no banco
Consulte: `STAGING_SEED_DIAGNOSIS.md` (solução de problemas completa)

---

## Referência Rápida de Arquivos

| Arquivo | Leia Se... |
|---------|-----------|
| `STAGING_SEED_QUICK_REFERENCE.md` | Quer uma referência rápida (2 min) |
| `STAGING_SEED_FIX_FINAL.md` | Quer entender o que foi corrigido (5 min) |
| `RUN_SEED_STAGING_INSTRUCTIONS.md` | Quer guia passo-a-passo completo |
| `STAGING_SEED_DIAGNOSIS.md` | Algo não funciona e precisa troubleshooting |
| `.env.staging.example` | Quer ver formato correto das variáveis |

---

## Cronograma Esperado

| Atividade | Tempo |
|-----------|-------|
| Verificar pré-requisitos | 5 min |
| Rodar seed | 2-5 min |
| Validar dados (queries) | 5 min |
| **Total** | **12-15 min** |

---

## Sucesso = Quando você conseguir:

✅ Rodar seed sem erros (exit code 0)
✅ Ver mensagem "completed successfully!"
✅ Validar 9 produtos na base de dados
✅ Encontrar admin user criado
✅ Confirmar 3 mesas criadas

---

## Documentação Disponível

```
Leia em ordem:
1. STAGING_SEED_QUICK_REFERENCE.md     (este arquivo + início rápido)
2. STAGING_SEED_FIX_FINAL.md           (entender a correção)
3. RUN_SEED_STAGING_INSTRUCTIONS.md    (referência completa)
4. STAGING_SEED_DIAGNOSIS.md           (se der problema)
```

---

## Dúvidas Frequentes

**P: Posso rodar sem Cloud SQL Proxy?**
R: Não. O seed precisa conectar ao Cloud SQL. Cloud SQL Proxy é obrigatório.

**P: Preciso de GCP credentials?**
R: Sim. `gcloud auth login` é necessário para conectar ao Cloud SQL.

**P: O seed vai sobrescrever dados existentes?**
R: Apenas se usar `--clear-first`. Sem a flag, dados são adicionados.

**P: Quanto tempo leva?**
R: Tipicamente 2-5 minutos, dependendo de latência de rede.

---

**Status**: ✅ Pronto para você executar!

Próximo: Siga as ações acima na ordem indicada.
