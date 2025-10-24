# 🍕 Guia: Seed Fattoria via Master Interactive Script

O Seed Fattoria foi integrado ao script `master-interactive.sh` para facilitar o uso!

---

## 🎯 Como Usar

### Opção 1: Menu Interativo (Recomendado)

```bash
# Iniciar o script
bash ./scripts/master-interactive.sh

# Depois:
# 1. Selecione: 3. 🌱 Database & Seeding
# 2. Escolha uma opção:
#    - 6. 🍕 Seed Fattoria (DEV)     [Popular local no Docker]
#    - 7. 🍕 Seed Fattoria (STAGE)   [Popular no Cloud SQL]
# 3. Confirme a execução
```

### Opção 2: Linha de Comando (Rápida)

```bash
# Executar Fattoria no DEV (Docker local)
bash ./scripts/master-interactive.sh --seed-fattoria-dev

# Executar Fattoria no STAGE (Cloud SQL)
bash ./scripts/master-interactive.sh --seed-fattoria-stage
```

### Opção 3: Usando Script Direto

```bash
# Se preferir o script específico do Fattoria:
bash ./scripts/run_seed_fattoria.sh --verbose
bash ./tools/scripts/seed/run_seed_fattoria.sh --verbose

# Ou para staging:
bash ./scripts/run_seed_staging.sh --verbose
```

---

## 📊 O que Será Inserido

**Fattoria Pizzeria**:
```
┌──────────────────────────────────────┐
│ 🍕 FATTORIA PIZZERIA                │
├──────────────────────────────────────┤
│ Organização: Fattoria Pizzeria       │
│ Email: contato@fattoria.com.br       │
│ Admin: admin@fattoria.com.br         │
│ Senha: password                      │
├──────────────────────────────────────┤
│ MENU (9 produtos):                   │
│ • Crostini.........................R$ 30 │
│ • Margherita........................R$ 80 │
│ • Marinara...........................R$ 58 │
│ • Parma............................R$ 109 │
│ • Vegana.............................R$ 60 │
│ • Suco de caju.......................R$ 15 │
│ • Heineken...........................R$ 13 │
│ • Baden Baden IPA....................R$ 23 │
│ • Sônia e Zé........................R$ 32 │
├──────────────────────────────────────┤
│ ESTRUTURA:                           │
│ • 1 Salão Principal (cap. 60)        │
│ • 3 Mesas (4, 2, 6 lugares)          │
└──────────────────────────────────────┘
```

---

## ✨ Exemplos Práticos

### Exemplo 1: Popular DEV com Fattoria (Uma linha)

```bash
bash ./scripts/master-interactive.sh --seed-fattoria-dev
```

**Resultado**:
- Cria organização Fattoria
- Insere 9 produtos
- Configura 3 mesas
- Cria usuário admin

### Exemplo 2: Acessar via Menu Interativo

```bash
bash ./scripts/master-interactive.sh
# Prensa ENTER no menu principal
# Selecione: 3 (Database & Seeding)
# Selecione: 6 (Fattoria DEV) ou 7 (Fattoria STAGE)
# Confirme: y
```

### Exemplo 3: Popular STAGE (Cloud SQL)

```bash
# Pré-requisito: Cloud SQL Proxy rodando
cloud_sql_proxy -instances=leps-472702:us-central1:leps-postgres-staging=unix:/cloudsql/...

# Em outro terminal:
bash ./scripts/master-interactive.sh --seed-fattoria-stage
```

---

## 🔧 Opções do Menu Database Atualizado

```
3. 🌱 Database & Seeding

  1. 🌱 Popular DEV (Docker local)        [Seed padrão]
  2. ☁️  Popular STAGE (Cloud SQL)         [Seed padrão]
  3. 🧹 Limpar e repopular DEV           [Deleta e recria]
  4. 🧹 Limpar e repopular STAGE          [Deleta e recria]
  5. 👥 Apenas usuários demo              [Só cria users]
  6. 🍕 Seed Fattoria (DEV)               [NOVO!]
  7. 🍕 Seed Fattoria (STAGE)             [NOVO!]
  8. 📊 Status das databases              [Verifica status]
  0. ⬅️  Voltar ao menu principal
```

---

## 🚀 Flags Batch Mode Disponíveis

```bash
# Ajuda
./scripts/master-interactive.sh --help

# Opções DEV
./scripts/master-interactive.sh --dev              # Iniciar ambiente
./scripts/master-interactive.sh --seed-dev         # Seed padrão
./scripts/master-interactive.sh --seed-fattoria-dev # NOVO!

# Opções STAGE
./scripts/master-interactive.sh --stage                # Menu STAGE
./scripts/master-interactive.sh --seed-fattoria-stage # NOVO!

# Outros
./scripts/master-interactive.sh --test             # Testes
./scripts/master-interactive.sh --status           # Status projeto
```

---

## ✅ Fluxo Padrão de Uso

```
1. Iniciar ambiente DEV
   $ ./scripts/master-interactive.sh --dev

2. Popular com Fattoria
   $ ./scripts/master-interactive.sh --seed-fattoria-dev

3. Testar API
   $ curl http://localhost:8080/health

4. Login Fattoria
   $ curl -X POST http://localhost:8080/login \
       -H "Content-Type: application/json" \
       -d '{"email":"admin@fattoria.com.br","password":"password"}'

5. Listar produtos Fattoria
   $ curl http://localhost:8080/product \
       -H "Authorization: Bearer <token>" \
       -H "X-Lpe-Organization-Id: <org_id>" \
       -H "X-Lpe-Project-Id: <project_id>"
```

---

## 🐛 Troubleshooting

### Problema: "Menu não mostra opção 6"

**Solução**: Script precisa estar atualizado
```bash
git pull origin dev
chmod +x ./scripts/master-interactive.sh
```

### Problema: "Seed não rodou"

**Solução**: Verificar se Docker ou GCP está disponível
```bash
# Para DEV
docker ps

# Para STAGE
gcloud auth list
gcloud config set project leps-472702
```

### Problema: "Erro ao rodar no DEV"

**Solução**: Garantir que PostgreSQL está rodando
```bash
./scripts/master-interactive.sh --dev    # Inicia containers
sleep 10                                  # Aguarda inicialização
./scripts/master-interactive.sh --seed-fattoria-dev
```

---

## 📋 Checklist de Sucesso

Após rodar `--seed-fattoria-dev`:

- [ ] Script executou sem erros (exit code 0)
- [ ] "✅ Seed execution completed successfully!" foi exibido
- [ ] Organização "Fattoria Pizzeria" foi criada
- [ ] 9 produtos foram inseridos
- [ ] Usuário admin@fattoria.com.br está ativo
- [ ] 3 mesas estão configuradas
- [ ] Consegue fazer login com admin@fattoria.com.br

---

## 🔗 Referências Relacionadas

| Recurso | Localização |
|---------|------------|
| Script Principal | `./scripts/master-interactive.sh` |
| Seed Fattoria | `./scripts/run_seed_fattoria.sh` |
| Seed Staging | `./scripts/run_seed_staging.sh` |
| Documentação Completa | `docs/seed/SEED_STAGING.md` |
| Status Seed | `STAGING_SEED_STATUS.md` |
| Guia Quick Start | `STAGING_SEED_GUIDE.md` |

---

## 💡 Dicas Rápidas

1. **Para testar rapidamente Fattoria no DEV:**
   ```bash
   bash ./scripts/master-interactive.sh --seed-fattoria-dev
   ```

2. **Para fazer backup antes de popular:**
   ```bash
   bash ./scripts/master-interactive.sh  # Menu 6 -> 3 (backup)
   ```

3. **Para limpar e repopular tudo:**
   ```bash
   bash ./scripts/master-interactive.sh  # Menu 3 -> 3 (DEV)
   # ou
   bash ./scripts/master-interactive.sh  # Menu 3 -> 4 (STAGE)
   ```

4. **Para ver todas as opções batch mode:**
   ```bash
   bash ./scripts/master-interactive.sh --help
   ```

---

**Status**: ✅ **Integração Completa**

O Seed Fattoria está totalmente integrado ao master script e pronto para uso!

Última atualização: 24 de Outubro, 2025
