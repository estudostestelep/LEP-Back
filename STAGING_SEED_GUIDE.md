# 🚀 Guia Rápido: Rodar Seed no Staging

**TL;DR** - Comando rápido para rodar o seed Fattoria em staging:

```bash
bash scripts/run_seed_staging.sh --clear-first --verbose
```

---

## 📋 Checklist de Pré-Requisitos

- [ ] Cloud SQL Proxy rodando: `cloud_sql_proxy -instances=leps-472702:us-central1:leps-postgres-staging=unix:/cloudsql/...`
- [ ] Autenticado no GCP: `gcloud auth list`
- [ ] Projeto correto: `gcloud config set project leps-472702`
- [ ] Arquivo `.env.staging` existe na raiz
- [ ] Go 1.23.0+ instalado: `go version`

---

## 🎯 Executar Seed

### Opção 1: Script (Recomendado - Mais Fácil)
```bash
# Sem limpeza de dados
bash scripts/run_seed_staging.sh

# Limpar dados antigos primeiro
bash scripts/run_seed_staging.sh --clear-first

# Com output verboso
bash scripts/run_seed_staging.sh --clear-first --verbose
```

### Opção 2: Comando Direto (Mais Controle)
```bash
export $(grep -v '^#' .env.staging | xargs)
go run ./cmd/seed/ --restaurant=fattoria --environment=staging --verbose
```

### Opção 3: Docker (Se Preferir Isolado)
```bash
docker build -t lep-seed:staging -f Dockerfile.dev .
docker run --rm --env-file .env.staging lep-seed:staging \
  go run cmd/seed/main.go --restaurant=fattoria --verbose
```

---

## ✅ Validar Dados

### Conectar ao Banco de Staging
```bash
gcloud sql connect leps-postgres-staging --user=lep_user --quiet
```

### Queries de Validação
```sql
-- 1. Verificar organização Fattoria
SELECT id, name FROM organizations WHERE name = 'Fattoria Pizzeria';

-- 2. Contar produtos (esperado: 9)
SELECT COUNT(*) as total_products FROM products
WHERE organization_id = '<ID_DA_FATTORIA>';

-- 3. Listar produtos Fattoria
SELECT name, price_normal FROM products
WHERE organization_id = '<ID_DA_FATTORIA>'
ORDER BY name;

-- 4. Verificar credenciais
SELECT email, name FROM users
WHERE email = 'admin@fattoria.com.br';
```

---

## 🔍 Monitorar Execução

### Em Tempo Real
```bash
# Ver operações do Cloud SQL
gcloud sql operations list --instance=leps-postgres-staging --limit=5

# Ver logs da API (se já deployada)
gcloud logging read "resource.type=cloud_run_revision" --limit 20
```

### Após Conclusão
```bash
# Ver dados inseridos
gcloud sql connect leps-postgres-staging --user=lep_user --quiet

-- Contar registros por tabela
SELECT 'organizations' as table_name, COUNT(*) FROM organizations
UNION ALL
SELECT 'products', COUNT(*) FROM products
UNION ALL
SELECT 'users', COUNT(*) FROM users
UNION ALL
SELECT 'tables', COUNT(*) FROM tables;
```

---

## 🐛 Problemas Comuns

| Problema | Solução |
|----------|---------|
| **Connection refused** | Iniciar Cloud SQL Proxy: `cloud_sql_proxy -instances=...` |
| **Permission denied** | `sudo chown -R $(whoami) /cloudsql && chmod 755 /cloudsql` |
| **Auth failed** | Verificar credenciais em `.env.staging` |
| **Silent output** | Adicionar flag: `--verbose` |
| **Sem feedback** | Usar: `2>&1 \| tee seed.log` para capturar output |

---

## 📊 O que será Inserido

### Organização
- **Nome**: Fattoria Pizzeria
- **Email**: contato@fattoria.com.br
- **Endereço**: Rua dos Italianos, 456 - São Paulo, SP

### Usuário Admin
- **Email**: admin@fattoria.com.br
- **Senha**: password
- **Permissões**: admin, products, orders, reservations, customers, tables, reports

### Menu (9 Produtos)
1. Crostini (R$ 30.00)
2. Margherita (R$ 80.00)
3. Marinara (R$ 58.00)
4. Parma (R$ 109.00)
5. Vegana (R$ 60.00)
6. Suco de caju (R$ 15.00)
7. Heineken (R$ 13.00)
8. Baden Baden IPA (R$ 23.00)
9. Sônia e Zé (R$ 32.00)

### Estrutura Física
- 1 Salão Principal (capacidade: 60)
- 3 Mesas (4, 2, 6 lugares)

---

## 🔗 Documentação Completa

Para detalhes adicionais, consulte:
- **[docs/seed/SEED_STAGING.md](docs/seed/SEED_STAGING.md)** - Guia completo e detalhado
- **[docs/seed/SEED_ARCHITECTURE.md](docs/seed/SEED_ARCHITECTURE.md)** - Arquitetura do sistema
- **[SEED_EXECUTION_REPORT.md](SEED_EXECUTION_REPORT.md)** - Troubleshooting avançado
- **[SESSION_SUMMARY.md](SESSION_SUMMARY.md)** - Histórico de mudanças

---

## ✨ Próximas Ações

```bash
# 1. Após seed bem-sucedido, fazer deploy
gcloud run deploy lep-api-staging \
  --source . \
  --region us-central1

# 2. Testar API
curl https://staging-api.lep.example.com/health

# 3. Testar Login
curl -X POST https://staging-api.lep.example.com/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@fattoria.com.br","password":"password"}'
```

---

**Last Updated**: 24 de Outubro, 2025
**Status**: ✅ Pronto para Staging
**Ambiente**: GCP Cloud SQL - us-central1
