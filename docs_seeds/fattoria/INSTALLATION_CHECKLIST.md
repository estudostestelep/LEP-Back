# 🍕 Fattoria Pizzeria Seed - Checklist de Instalação

Checklist passo-a-passo para instalar e validar o seed da Fattoria.

## ✅ Pré-requisitos

- [ ] Go 1.16+ instalado (`go version`)
- [ ] PostgreSQL rodando e acessível
- [ ] Arquivo `.env` configurado com credenciais do banco
- [ ] Permissões de escrita na pasta do projeto

## ✅ Instalação dos Arquivos

### Novos Arquivos Criados
- [ ] `utils/seed_fattoria.go` (412 linhas) ✅ Criado
- [ ] `scripts/run_seed_fattoria.sh` (150+ linhas) ✅ Criado
- [ ] `SEED_FATTORIA.md` (500+ linhas) ✅ Criado
- [ ] `SEED_FATTORIA_SUMMARY.md` (200+ linhas) ✅ Criado
- [ ] `cmd/seed/README_FATTORIA.md` (150+ linhas) ✅ Criado
- [ ] `cmd/seed/FATTORIA_IDS.md` (300+ linhas) ✅ Criado
- [ ] `FATTORIA_MENU.txt` (visual) ✅ Criado
- [ ] `INSTALLATION_CHECKLIST.md` (este arquivo) ✅ Criado

### Arquivos Modificados
- [ ] `cmd/seed/main.go`
  - [ ] Adicionado flag `--restaurant`
  - [ ] Adicionado switch statement
  - [ ] Mensagens atualizadas

## ✅ Configuração do Ambiente

### Verificar Go
```bash
# Verificar versão
go version
# Esperado: go version go1.16+ ...
```
- [ ] Go version >= 1.16

### Verificar Dependências
```bash
# Na raiz do LEP-Back
go mod tidy
```
- [ ] Sem erros de dependências

### Compilação de Teste
```bash
# Compilar o seed
go build -o /dev/null cmd/seed/main.go
```
- [ ] Compila sem erros
- [ ] Sem warnings

### Verificar PostgreSQL
```bash
# Testar conexão (substitute com suas credenciais)
psql -U postgres -h localhost -c "SELECT version();"
```
- [ ] Conexão bem-sucedida

### Verificar .env
```bash
# Verificar arquivo .env
cat .env | grep -E "DB_USER|DB_PASS|DB_NAME"
```
- [ ] `DB_USER` configurado
- [ ] `DB_PASS` configurado
- [ ] `DB_NAME` configurado
- [ ] `INSTANCE_UNIX_SOCKET` (se usando GCP)

## ✅ Teste de Execução

### Opção 1: Script Bash

```bash
# 1. Verificar se script é executável
ls -l scripts/run_seed_fattoria.sh
# Esperado: -rwx------

# 2. Dar permissões se necessário
chmod +x scripts/run_seed_fattoria.sh

# 3. Executar seed (sem limpeza)
bash scripts/run_seed_fattoria.sh
```
- [ ] Script é executável
- [ ] Script executa sem erros
- [ ] Mensagens de sucesso aparecem

### Opção 2: Comando Go Direto

```bash
# Executar seed diretamente
go run cmd/seed/main.go --restaurant=fattoria
```
- [ ] Comando executa
- [ ] Sem erros de compilação
- [ ] Saída em português

### Opção 3: Com Limpeza (CUIDADO!)

```bash
# AVISO: Isso vai limpar TODOS os dados!
bash scripts/run_seed_fattoria.sh --clear-first
```
- [ ] Confirmar que deseja limpar dados
- [ ] Executa sem erros
- [ ] Database limpo e repopulado

## ✅ Validação de Dados

### Verificar Seed via SQL

```bash
# Conectar ao banco
psql -U $DB_USER -d $DB_NAME -h localhost
```

```sql
-- 1. Verificar organizações
SELECT id, name FROM organizations WHERE name LIKE '%Fattoria%';
-- Esperado: 1 linha com "Fattoria Pizzeria"

-- 2. Verificar projetos
SELECT id, name FROM projects WHERE name LIKE '%Fattoria%';
-- Esperado: 1 linha com "Fattoria Pizzeria - Projeto Principal"

-- 3. Verificar categorias
SELECT id, name FROM categories
WHERE organization_id = '223e4567-e89b-12d3-a456-426614174100'
ORDER BY "order";
-- Esperado: 8 categorias (Pizzas, Bebidas, + 6 subcategorias)

-- 4. Verificar produtos
SELECT id, name, price_normal FROM products
WHERE organization_id = '223e4567-e89b-12d3-a456-426614174100'
ORDER BY "order";
-- Esperado: 9 produtos
-- Margherita (80.00), Marinara (58.00), Parma (109.00), etc.

-- 5. Verificar tags
SELECT id, name FROM tags
WHERE organization_id = '223e4567-e89b-12d3-a456-426614174100';
-- Esperado: 2 tags (Vegetariana, Vegana)

-- 6. Verificar mesas
SELECT number, capacity, status FROM tables
WHERE organization_id = '223e4567-e89b-12d3-a456-426614174100'
ORDER BY number;
-- Esperado: 3 mesas

-- 7. Verificar usuário admin
SELECT id, email, name FROM users
WHERE email = 'admin@fattoria.com.br';
-- Esperado: 1 usuário admin
```

### Checklist de Validação SQL
- [ ] 1 organização Fattoria
- [ ] 1 projeto Fattoria
- [ ] 8 categorias (corretas)
- [ ] 9 produtos (corretos)
- [ ] 2 tags (Vegetariana, Vegana)
- [ ] 3 mesas
- [ ] 1 usuário admin
- [ ] 1 ambiente (Salão Principal)

## ✅ Teste de API

### 1. Health Check
```bash
curl -X GET http://localhost:8080/health
```
- [ ] Status: 200 OK
- [ ] Resposta: `{"status":"healthy"}`

### 2. Ping
```bash
curl -X GET http://localhost:8080/ping
```
- [ ] Status: 200 OK
- [ ] Resposta: `pong`

### 3. Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@fattoria.com.br",
    "password": "password"
  }'
```
- [ ] Status: 200 OK
- [ ] Resposta contém `"token"`
- [ ] Token é uma string válida

### 4. Obter Produtos
```bash
# Usar token do login anterior
TOKEN="seu_token_aqui"
ORG_ID="223e4567-e89b-12d3-a456-426614174100"
PROJ_ID="223e4567-e89b-12d3-a456-426614174101"

curl -X GET http://localhost:8080/product \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Lpe-Organization-Id: $ORG_ID" \
  -H "X-Lpe-Project-Id: $PROJ_ID"
```
- [ ] Status: 200 OK
- [ ] Resposta contém array de produtos
- [ ] Todos os 9 produtos presentes

### 5. Obter Mesas
```bash
curl -X GET http://localhost:8080/table \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Lpe-Organization-Id: $ORG_ID" \
  -H "X-Lpe-Project-Id: $PROJ_ID"
```
- [ ] Status: 200 OK
- [ ] 3 mesas retornadas
- [ ] Mesas têm números corretos (1, 2, 3)

## ✅ Teste Manual no Frontend

Se tiver o frontend rodando:

1. Acessar `http://localhost:5173`
2. Login com:
   - Email: `admin@fattoria.com.br`
   - Senha: `password`
3. Verificar:
   - [ ] Menu carrega corretamente
   - [ ] Produtos aparecem na lista
   - [ ] Categorias são corretas
   - [ ] Preços estão corretos
   - [ ] Tags aparecem (Vegetariana, Vegana)
   - [ ] Mesas aparecem na interface

## ✅ Troubleshooting

### Erro: "Failed to connect to database"
```bash
# Verificar se PostgreSQL está rodando
docker ps | grep postgres
# ou
psql --version
```
- [ ] PostgreSQL está rodando
- [ ] Credenciais em `.env` estão corretas

### Erro: "Unknown option: --restaurant"
```bash
# Verificar se main.go foi modificado corretamente
grep -n "restaurant" cmd/seed/main.go
```
- [ ] Linha com `StringVar` para `--restaurant`
- [ ] Switch statement presente

### Erro: "ParentId field not found"
✅ **JÁ CORRIGIDO** - Use a versão corrigida de `seed_fattoria.go`

### Erro: "Table already exists"
```bash
# Execute com --clear-first para limpar primeiro
bash scripts/run_seed_fattoria.sh --clear-first
```
- [ ] Dados limpados
- [ ] Novo seed executado

### Erro: "Permission denied" no script
```bash
# Dar permissões de execução
chmod +x scripts/run_seed_fattoria.sh
```
- [ ] Script tem permissão de execução

## ✅ Performance e Validação Final

### Tempo de Execução Esperado
- [ ] Seed sem limpeza: < 30 segundos
- [ ] Seed com limpeza: < 60 segundos

### Consumo de Recursos
- [ ] CPU < 50% durante execução
- [ ] Memória < 200MB
- [ ] Disco: ~1-2MB adicionados

### Logs Esperados
- [ ] "✅ Database seeding completed successfully!"
- [ ] "📊 Seeding Summary" com contagens
- [ ] Credenciais de login exibidas

## ✅ Setup Completo

Complete todos os itens acima, então:

### Opção 1: Desenvolvimento Local
```bash
# Terminal 1 - Backend
cd LEP-Back
go run main.go

# Terminal 2 - Frontend (se aplicável)
cd LEP-Front
npm run dev

# Browser
http://localhost:5173
```

### Opção 2: Apenas Backend
```bash
cd LEP-Back
bash scripts/run_seed_fattoria.sh --clear-first
go run main.go

# Teste via curl
curl http://localhost:8080/health
```

### Opção 3: Teste Automatizado
```bash
bash scripts/run_seed_fattoria.sh --clear-first --verbose
bash scripts/run_tests.sh
```

## 📋 Documentação de Referência

Após validar tudo, consulte:

- [SEED_FATTORIA.md](SEED_FATTORIA.md) - Documentação completa
- [cmd/seed/README_FATTORIA.md](cmd/seed/README_FATTORIA.md) - Quick start
- [cmd/seed/FATTORIA_IDS.md](cmd/seed/FATTORIA_IDS.md) - IDs reference
- [FATTORIA_MENU.txt](FATTORIA_MENU.txt) - Visual do menu

## 🎉 Sucesso!

Se todos os itens acima estão marcados, o seed da Fattoria está **totalmente funcional** e pronto para:

- ✅ Desenvolvimento
- ✅ Testes e QA
- ✅ Demonstrações
- ✅ Validação de integrações

## 📞 Suporte

Se encontrar problemas:

1. Consulte o arquivo `SEED_FATTORIA.md`
2. Verifique este checklist
3. Valide com SQL commands
4. Teste com curl commands
5. Execute com `--verbose` para mais detalhes

---

**Status**: ✅ Checklist Completo
**Versão**: 1.0
**Última Atualização**: 2024
