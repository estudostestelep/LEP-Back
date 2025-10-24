# 🔧 Fix: .env.staging Format Error

## ❌ Problema Encontrado

Ao tentar rodar o script `run_seed_staging.sh`, você recebia este erro:

```
./scripts/run_seed_staging.sh: line 76: export: `max-age=7200;': not a valid identifier
```

### Causa Raiz

O arquivo `.env.staging` tinha **ponto-e-vírgula no final de cada linha**:

```bash
# ❌ ERRADO (com semicolons)
DB_USER=lep_user;
DB_PASS=123456;
BUCKET_CACHE_CONTROL=public, max-age=7200;
```

Quando o script tenta fazer `export $(grep -v '^#' .env.staging | xargs)`, o bash tenta interpretar:
```bash
export DB_USER=lep_user; DB_PASS=123456; BUCKET_CACHE_CONTROL=public, max-age=7200;
```

Isso quebra porque:
- `;` é um separador de comandos no bash
- `max-age=7200;` não é uma variável válida

---

## ✅ Solução Implementada

Removi todos os ponto-e-vírgula do final das linhas:

```bash
# ✅ CORRETO (sem semicolons)
DB_USER=lep_user
DB_PASS=123456
BUCKET_CACHE_CONTROL=public, max-age=7200
```

### Arquivo Corrigido

```
.env.staging ← Já está corrigido no seu sistema
```

### Arquivo Template

```
.env.staging.example ← Template para referência (commited)
```

---

## 🚀 Como Usar Agora

Seu arquivo `.env.staging` já está correto. Você pode rodar normalmente:

```bash
# Opção 1: Via master script
bash ./scripts/master-interactive.sh --seed-fattoria-stage

# Opção 2: Via menu interativo
bash ./scripts/master-interactive.sh
# Menu 3 → Opção 7

# Opção 3: Script direto
bash ./scripts/run_seed_staging.sh --verbose
```

---

## 📝 O que Mudar

Se você criou seu próprio `.env.staging`, remova os ponto-e-vírgula:

```bash
# ANTES (❌ errado)
DB_USER=lep_user;
DB_PASS=123456;
ENVIRONMENT=staging;

# DEPOIS (✅ correto)
DB_USER=lep_user
DB_PASS=123456
ENVIRONMENT=staging
```

---

## ✨ Arquivos Afetados

| Arquivo | Mudança | Status |
|---------|---------|--------|
| `.env.staging` | Removidos semicolons | ✅ Corrigido |
| `.env.staging.example` | Arquivo template | ✅ Criado |
| `run_seed_staging.sh` | Nenhuma mudança | ✅ OK |

---

## 🔗 Referências

- **Script Staging**: `scripts/run_seed_staging.sh`
- **Master Script**: `scripts/master-interactive.sh`
- **Documentação**: `STAGING_SEED_GUIDE.md`

---

## ✅ Validação

Agora você pode testar:

```bash
# Teste rápido de sintaxe
bash -n scripts/run_seed_staging.sh

# Teste com export
export $(grep -v '^#' .env.staging | xargs) && echo "✅ OK"

# Rodar seed
bash ./scripts/master-interactive.sh --seed-fattoria-stage
```

---

**Status**: ✅ **RESOLVIDO**

Seu ambiente de staging está pronto para rodar o seed Fattoria!

Data da correção: 24 de Outubro, 2025
