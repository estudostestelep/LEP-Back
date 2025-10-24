# Staging Seed - Fixes Aplicadas (Sessao Final)

Data: 24 de Outubro, 2025
Status: PRONTO PARA TESTAR
Fixes: 3 correcoes criticas identificadas e corrigidas

---

## Problema Original

Ao tentar executar o seed Fattoria em staging, o erro abaixo aparecia:

./scripts/run_seed_staging.sh: line 76: export: 'max-age=7200': not a valid identifier

---

## Root Cause Analysis

Foram identificadas 2 causas raiz diferentes:

### Causa 1: Metodo de carregamento de variaveis (CRITICA)
Arquivo: scripts/run_seed_staging.sh (linha 76)

Antes (ERRADO):
export $(grep -v '^#' "$ENV_FILE" | xargs)

Por que falha: O metodo export $(... | xargs) nao consegue lidar com valores que contem espacos, virgulas ou caracteres especiais.

Depois (CORRETO):
set -a
source "$ENV_FILE"
set +a

Por que funciona: O comando source processa todo o arquivo corretamente.

### Causa 2: Valores nao escapados no .env.staging
Arquivo: .env.staging (linha 12)

Antes (ERRADO):
BUCKET_CACHE_CONTROL=public, max-age=7200

Depois (CORRETO):
BUCKET_CACHE_CONTROL="public, max-age=7200"

Regra: Sempre coloque aspas em valores que contem espacos, virgulas ou caracteres especiais!

---

## Arquivos Corrigidos

### 1. scripts/run_seed_staging.sh
Status: Corrigido
Linhas: 76-78
Mudanca: Substitui export $(grep -v '^#' "$ENV_FILE" | xargs) por set -a; source "$ENV_FILE"; set +a
Validacao: bash -n Sintaxe valida

### 2. tools/scripts/seed/run_seed_staging.sh
Status: Corrigido (arquivo espelhado)
Linhas: 76-78
Mudanca: Mesma correcao acima
Validacao: bash -n Sintaxe valida

### 3. .env.staging
Status: Corrigido
Linha: 12
Mudanca: Colocadas aspas em BUCKET_CACHE_CONTROL="public, max-age=7200"
Validacao: Variaveis carregam corretamente

### 4. .env.staging.example (NOVO)
Status: Criado
Localizacao: Raiz do projeto
Proposito: Template de referencia mostrando formato correto

---

## Testes Executados

### Teste 1: Validacao de Sintaxe Bash
bash -n scripts/run_seed_staging.sh
Sintaxe valida

bash -n tools/scripts/seed/run_seed_staging.sh
Sintaxe valida

### Teste 2: Carregamento de Variaveis
source .env.staging
DB_USER=lep_user
DB_NAME=lep_database
ENVIRONMENT=staging
BUCKET_CACHE_CONTROL=public, max-age=7200
BUCKET_TIMEOUT=60

Antes da correcao, o teste 2 retornava:
max-age=7200: command not found

---

## Como Testar Agora

### Opcao 1: Via Master Interactive Script (Recomendado)

bash ./scripts/master-interactive.sh

No menu:
- Selecione: 3 (Database & Seeding)
- Selecione: 7 (Seed Fattoria STAGE)
- Confirme: y

### Opcao 2: Executar Script Direto

bash ./scripts/run_seed_staging.sh --clear-first --verbose

### Opcao 3: Via Batch Mode

bash ./scripts/master-interactive.sh --seed-fattoria-stage

---

## Pre-requisitos para Sucesso

1. Cloud SQL Proxy esta rodando?
   ps aux | grep cloud_sql_proxy

2. GCP autenticado?
   gcloud auth list
   gcloud config set project leps-472702

3. Variaveis de ambiente carregando?
   source .env.staging
   echo $BUCKET_CACHE_CONTROL

---

## Validacao Apos Execucao

Apos rodar o seed, execute essas queries para validar:

1. Verificar Organizacao Fattoria
SELECT id, name FROM organizations 
WHERE name = 'Fattoria Pizzeria' AND deleted_at IS NULL;

2. Contar Produtos
SELECT COUNT(*) as total_produtos FROM products 
WHERE organization_id = (
  SELECT id FROM organizations WHERE name = 'Fattoria Pizzeria'
) AND deleted_at IS NULL;
Esperado: 9

3. Listar Produtos
SELECT name, price_normal FROM products 
WHERE organization_id = (
  SELECT id FROM organizations WHERE name = 'Fattoria Pizzeria'
) AND deleted_at IS NULL
ORDER BY name;

4. Verificar Admin
SELECT email, name FROM users 
WHERE email = 'admin@fattoria.com.br' AND deleted_at IS NULL;

5. Contar Mesas
SELECT COUNT(*) as total_mesas FROM tables 
WHERE organization_id = (
  SELECT id FROM organizations WHERE name = 'Fattoria Pizzeria'
) AND deleted_at IS NULL;
Esperado: 3

---

## Se Algo Ainda Nao Funcionar

### Erro: "Connection refused"
Solucao: Garantir que Cloud SQL Proxy esta rodando

### Erro: "max-age=7200: command not found"
Solucao: Verifique se .env.staging linha 12 tem aspas

### Erro: "Database does not exist"
Solucao: Criar database antes

### Sem erros, mas dados nao aparecem no banco
Verificar:
1. Qual database foi usada?
2. Qual host?
3. Cloud SQL Proxy conectado ao socket correto?

---

## Checklist de Sucesso

- [ ] Scripts rodaram sem erro bash (exit code 0)
- [ ] Nenhum erro de "max-age=7200: command not found"
- [ ] Seed foi executado
- [ ] Seed completou
- [ ] Organizacao "Fattoria Pizzeria" existe no banco
- [ ] 9 produtos foram inseridos
- [ ] Admin user admin@fattoria.com.br foi criado
- [ ] 3 mesas foram criadas

---

Documento Criado: 24 de Outubro, 2025
Versao: 1.0
Status: Todas as correcoes aplicadas e testadas
