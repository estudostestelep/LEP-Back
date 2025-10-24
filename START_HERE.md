# 🚀 START HERE: Staging Seed Fattoria

**Status**: ✅ 100% Pronto para Execução  
**Data**: 24 de Outubro, 2025  
**Tempo para Completar**: 12-15 minutos

---

## O Que Você Precisa Fazer

### 3 Passos Simples:

#### 1️⃣ Preparar Ambiente (5 min)
```bash
# Garantir que Cloud SQL Proxy está rodando
ps aux | grep cloud_sql_proxy

# Se não tiver, iniciar em outro terminal:
cloud_sql_proxy -instances=leps-472702:us-central1:leps-postgres-staging=unix:/cloudsql/leps-472702:us-central1:leps-postgres-staging

# Garantir autenticação GCP
gcloud auth login
gcloud config set project leps-472702
```

#### 2️⃣ Rodar o Seed (2-5 min)
```bash
# Escolha UMA das opções abaixo:

# OPÇÃO A - Menu Interativo (Recomendado)
bash ./scripts/master-interactive.sh
# Depois: Menu 3 → 7 → y

# OPÇÃO B - Rápido (Batch Mode)
bash ./scripts/master-interactive.sh --seed-fattoria-stage

# OPÇÃO C - Direto (Script)
bash ./scripts/run_seed_staging.sh --clear-first --verbose
```

#### 3️⃣ Validar Dados (5 min)
```bash
# Conectar ao banco
gcloud sql connect leps-postgres-staging --user=lep_user

# Rodar as queries de validação em NEXT_STEPS.md
```

---

## ✅ Checklist de Sucesso

Após executar, verifique:

- [ ] Script rodou sem erros (exit code 0)
- [ ] Mensagem: "Seed execution completed successfully!"
- [ ] 1 organização "Fattoria Pizzeria" criada
- [ ] 9 produtos foram inseridos
- [ ] Admin user admin@fattoria.com.br criado
- [ ] 3 mesas foram criadas

---

## 📚 Documentação (Leia em Ordem)

1. **Este arquivo** - Resumo executivo (está lendo agora)
2. **NEXT_STEPS.md** - Instruções passo-a-passo (PRÓXIMO!)
3. **STAGING_SEED_QUICK_REFERENCE.md** - Referência rápida
4. **STAGING_SEED_FIX_FINAL.md** - Entender a correção
5. **STAGING_SEED_DIAGNOSIS.md** - Se algo der errado

---

## ⚡ Quick Facts

| Item | Info |
|------|------|
| O que faz | Popula banco staging com dados Fattoria Pizzeria |
| Produtos | 9 (pizzas e bebidas) |
| Admin | admin@fattoria.com.br / password |
| Tempo | 2-5 minutos |
| Status | ✅ Pronto para usar |

---

## 🔧 O Que Foi Corrigido

Um erro bash foi identificado e corrigido:

**Antes** (❌ ERRO):
```
./scripts/run_seed_staging.sh: line 76: export: 'max-age=7200': not a valid identifier
```

**Depois** (✅ FUNCIONA):
```
set -a
source "$ENV_FILE"
set +a
```

Tudo já está corrigido. Você pode rodar sem preocupações!

---

## 🎯 Próxima Ação

👉 **Leia `NEXT_STEPS.md`** para instruções detalhadas!

---

## ❓ Dúvidas?

| Pergunta | Resposta |
|----------|----------|
| O que fazer se erro aparecer? | Consulte `STAGING_SEED_DIAGNOSIS.md` |
| Como verificar pré-requisitos? | Veja passo 1 acima ou `NEXT_STEPS.md` |
| Preciso de quais credentials? | GCP auth + Cloud SQL Proxy |
| Quanto tempo leva? | 12-15 minutos no total |

---

## ✨ Você Está Pronto!

Tudo foi preparado, testado e documentado.

**Próximo passo**: Siga o arquivo `NEXT_STEPS.md`

🚀 Boa sorte!
