# 🌱 Rodando Seed no Ambiente de Staging

Guia completo para executar o seed Fattoria no ambiente de staging (Cloud SQL no GCP).

---

## 📋 Pré-requisitos

### 1. Ambiente GCP Configurado
```bash
# Verificar se você está autenticado no GCP
gcloud auth list
gcloud config list

# Certificar-se de estar no projeto certo
gcloud config set project leps-472702
```

### 2. Cloud SQL Proxy em Execução
```bash
# Se usando Unix socket (Linux/Mac)
cloud_sql_proxy -instances=leps-472702:us-central1:leps-postgres-staging=unix:/cloudsql/leps-472702:us-central1:leps-postgres-staging

# Se usando TCP (Windows ou sem Unix socket)
cloud_sql_proxy -instances=leps-472702:us-central1:leps-postgres-staging=tcp:5432
```

### 3. Variáveis de Ambiente
- Arquivo `.env.staging` deve existir na raiz do projeto
- Contém credenciais de staging (não commitar!)

### 4. Go instalado
```bash
go version  # Deve ser >= 1.23.0
```

---

## 🚀 Execução

### Opção 1: Usando o Script Bash (Recomendado)

```bash
# Seed padrão Fattoria
bash scripts/run_seed_staging.sh

# Com limpeza de dados anterior
bash scripts/run_seed_staging.sh --clear-first

# Com saída verbosa
bash scripts/run_seed_staging.sh --verbose

# Combinar flags
bash scripts/run_seed_staging.sh --clear-first --verbose
```

**Outra localização**:
```bash
bash tools/scripts/seed/run_seed_staging.sh
```

### Opção 2: Comando Direto Go

```bash
# Carregar variáveis de staging
export $(grep -v '^#' .env.staging | xargs)

# Rodar seed
go run ./cmd/seed/ \
  --restaurant=fattoria \
  --environment=staging \
  --clear-first \
  --verbose
```

### Opção 3: Via Docker (Se disponível)

```bash
# Build e execute via Docker
docker build -t lep-seed:staging -f Dockerfile.dev .

docker run --rm \
  --env-file .env.staging \
  -e DB_HOST=host.docker.internal \
  lep-seed:staging \
  go run cmd/seed/main.go --restaurant=fattoria --verbose
```

---

## 📊 Monitoramento da Execução

### Durante a Execução

```bash
# Ver logs em tempo real
gcloud sql operations describe OPERATION_ID \
  --instance=leps-postgres-staging

# Ver conexões ativas
gcloud sql connect leps-postgres-staging \
  --user=lep_user \
  --quiet

-- Depois na conexão:
SELECT * FROM pg_stat_activity;
```

### Após a Execução

```bash
# Conectar ao banco
gcloud sql connect leps-postgres-staging \
  --user=lep_user \
  --quiet

-- Validar dados inseridos:
SELECT COUNT(*) as product_count FROM products;
SELECT COUNT(*) as org_count FROM organizations WHERE name = 'Fattoria Pizzeria';
SELECT name, price_normal FROM products WHERE name LIKE '%Pizza%' LIMIT 5;
```

---

## ✅ Validação de Sucesso

### Checklist Pós-Seed

```sql
-- 1. Organização Fattoria criada
SELECT id, name, email FROM organizations
WHERE name = 'Fattoria Pizzeria';

-- 2. Produtos inseridos (esperado: 9)
SELECT COUNT(*) FROM products
WHERE organization_id = <fattoria_org_id>;

-- 3. Usuário admin Fattoria
SELECT email, name FROM users
WHERE email = 'admin@fattoria.com.br';

-- 4. Mesas e ambientes configurados
SELECT COUNT(*) FROM tables
WHERE organization_id = <fattoria_org_id>;

-- 5. Categorias carregadas
SELECT COUNT(*) FROM categories
WHERE organization_id = <fattoria_org_id>;
```

---

## 🐛 Troubleshooting

### Problema: "Connection refused"
```
Solução: Garantir que Cloud SQL Proxy está rodando
- Verifique: ps aux | grep cloud_sql_proxy
- Inicie se necessário: cloud_sql_proxy -instances=...
```

### Problema: "Permission denied" (Unix socket)
```
Solução: Verificar permissões do diretório
sudo chown -R $(whoami):$(whoami) /cloudsql
chmod 755 /cloudsql
```

### Problema: "Authentication failed"
```
Solução: Verificar credenciais em .env.staging
- Confirmar DB_USER e DB_PASS estão corretos
- Verificar se usuário existe no Cloud SQL
```

### Problema: "Database does not exist"
```
Solução: Criar banco antes de rodar seed
gcloud sql databases create lep_database \
  --instance=leps-postgres-staging
```

### Problema: "Silent execution" (sem output)
```
Solução: Usar flag --verbose
bash scripts/run_seed_staging.sh --verbose

Ou adicionar debug:
go run ./cmd/seed/ --restaurant=fattoria --verbose 2>&1 | tee seed.log
```

---

## 📋 Dados Inseridos

### Organização
| Campo | Valor |
|-------|-------|
| Nome | Fattoria Pizzeria |
| Email | contato@fattoria.com.br |
| Endereço | Rua dos Italianos, 456 - São Paulo, SP |

### Usuário Admin
| Campo | Valor |
|-------|-------|
| Email | admin@fattoria.com.br |
| Senha | password |
| Permissões | admin, products, orders, reservations, customers, tables, reports |

### Menu
**9 Produtos**:
1. **Crostini** - R$ 30.00 (Entrada)
2. **Margherita** - R$ 80.00 (Pizza)
3. **Marinara** - R$ 58.00 (Pizza)
4. **Parma** - R$ 109.00 (Pizza)
5. **Vegana** - R$ 60.00 (Pizza)
6. **Suco de caju integral** - R$ 15.00 (Bebida)
7. **Heineken s/ álcool** - R$ 13.00 (Cerveja)
8. **Baden Baden IPA** - R$ 23.00 (Cerveja Artesanal)
9. **Sônia e Zé** - R$ 32.00 (Coquetel)

### Ambientes
- 1 salão principal com capacidade de 60 pessoas
- 3 mesas (4, 2, 6 lugares)

---

## 🔗 Próximas Ações

Após rodar o seed com sucesso:

1. **Deploy da API**
   ```bash
   gcloud run deploy lep-api-staging \
     --source . \
     --region us-central1 \
     --platform managed
   ```

2. **Testar Conectividade**
   ```bash
   curl https://staging-api.lep.example.com/health
   ```

3. **Testar Login**
   ```bash
   curl -X POST https://staging-api.lep.example.com/login \
     -H "Content-Type: application/json" \
     -d '{"email":"admin@fattoria.com.br","password":"password"}'
   ```

4. **Monitorar em Produção**
   ```bash
   gcloud logging read "resource.type=cloud_run_revision" \
     --limit 50 \
     --format json
   ```

---

## 📞 Referências

- **Documentação Seed**: [docs/seed/SEED_ARCHITECTURE.md](SEED_ARCHITECTURE.md)
- **Troubleshooting Geral**: [SEED_EXECUTION_REPORT.md](../SEED_EXECUTION_REPORT.md)
- **Cloud SQL Proxy**: https://cloud.google.com/sql/docs/postgres/sql-proxy
- **GCP IAM**: https://cloud.google.com/iam/docs

---

**Última Atualização**: 24 de Outubro, 2025
**Ambiente**: Staging (Cloud SQL - GCP)
**Responsável**: Claude Code

