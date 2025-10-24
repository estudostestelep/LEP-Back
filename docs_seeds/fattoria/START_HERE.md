# 🍕 Fattoria Pizzeria Seed - COMECE AQUI

Bem-vindo! Este arquivo é o ponto de entrada para o novo seed da Fattoria Pizzeria.

## ⚡ 30 Segundos para Começar

```bash
# 1. Vá para a pasta do backend
cd LEP-Back

# 2. Execute o seed com limpeza
bash scripts/run_seed_fattoria.sh --clear-first

# 3. Inicie o servidor
go run main.go

# 4. Acesse a API
curl http://localhost:8080/health
```

**Login**:
- Email: `admin@fattoria.com.br`
- Senha: `password`

## 📚 Documentação por Tipo

### 🏃 Rápido (5 minutos)
1. **[FATTORIA_MENU.txt](FATTORIA_MENU.txt)** - Menu visual em ASCII
2. **[cmd/seed/README_FATTORIA.md](cmd/seed/README_FATTORIA.md)** - Quick start

### 📖 Completo (20 minutos)
1. **[SEED_FATTORIA.md](SEED_FATTORIA.md)** - Documentação principal
2. **[cmd/seed/FATTORIA_IDS.md](cmd/seed/FATTORIA_IDS.md)** - Referência de IDs

### ✅ Implementação (30 minutos)
1. **[SEED_FATTORIA_SUMMARY.md](SEED_FATTORIA_SUMMARY.md)** - Resumo técnico
2. **[INSTALLATION_CHECKLIST.md](INSTALLATION_CHECKLIST.md)** - Validação passo-a-passo

## 🎯 O que você vai ter

### 🍕 Cardápio Completo
```
Pizzas:
  • Crostini (Entrada) - R$ 30,00
  • Marguerita 🌱 - R$ 80,00
  • Marinara 🌿 - R$ 58,00
  • Parma - R$ 109,00
  • Vegana 🌿 - R$ 60,00

Bebidas:
  • Suco de Caju - R$ 15,00
  • Heineken s/álc - R$ 13,00
  • Baden Baden IPA - R$ 23,00
  • Sônia e Zé (Coquetel) - R$ 32,00
```

### 🏪 Estrutura Operacional
- 3 Mesas (4, 2, 6 lugares)
- 1 Salão (60 pessoas)
- 2 Tags (Vegetariana, Vegana)
- 1 Admin (admin@fattoria.com.br)

## 🚀 Opções de Execução

### Opção A: Script (Recomendado)
```bash
# Seed simples
bash scripts/run_seed_fattoria.sh

# Com limpeza de dados
bash scripts/run_seed_fattoria.sh --clear-first

# Com saída detalhada
bash scripts/run_seed_fattoria.sh --verbose

# Todas as opções
bash scripts/run_seed_fattoria.sh --help
```

### Opção B: Go Direto
```bash
# Sem limpeza
go run cmd/seed/main.go --restaurant=fattoria

# Com opções
go run cmd/seed/main.go --restaurant=fattoria --clear-first --verbose
```

### Opção C: Seed Padrão (mantém compatibilidade)
```bash
# Seed do LEP Demo (não do Fattoria)
bash scripts/run_seed.sh
# ou
go run cmd/seed/main.go
```

## 🔐 Credenciais

| Campo | Valor |
|-------|-------|
| **Email** | admin@fattoria.com.br |
| **Senha** | password |
| **Organização** | Fattoria Pizzeria |

## 📁 Arquivos Criados

```
LEP-Back/
├── 📄 START_HERE.md (este arquivo)
├── 📄 SEED_FATTORIA.md (documentação completa)
├── 📄 SEED_FATTORIA_SUMMARY.md (resumo técnico)
├── 📄 FATTORIA_MENU.txt (menu visual)
├── 📄 INSTALLATION_CHECKLIST.md (validação)
│
├── utils/
│   └── 📄 seed_fattoria.go (código do seed)
│
├── scripts/
│   └── 📜 run_seed_fattoria.sh (script bash)
│
├── cmd/seed/
│   ├── 📄 README_FATTORIA.md (quick start)
│   ├── 📄 FATTORIA_IDS.md (IDs reference)
│   └── 📝 main.go (modificado - flag --restaurant)
```

## 🔧 Pré-requisitos

- [ ] Go 1.16+
- [ ] PostgreSQL rodando
- [ ] `.env` configurado
- [ ] Permissões de escrita

## ✅ Validação Rápida

Após executar o seed:

```bash
# 1. Health check
curl http://localhost:8080/health
# Esperado: {"status":"healthy"}

# 2. Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@fattoria.com.br","password":"password"}'
# Esperado: Token JWT

# 3. Produtos
curl -X GET http://localhost:8080/product \
  -H "Authorization: Bearer {TOKEN}" \
  -H "X-Lpe-Organization-Id: 223e4567-e89b-12d3-a456-426614174100" \
  -H "X-Lpe-Project-Id: 223e4567-e89b-12d3-a456-426614174101"
# Esperado: 9 produtos
```

## 📖 Próximos Passos

### Para Desenvolvedores
1. Leia [SEED_FATTORIA.md](SEED_FATTORIA.md)
2. Consulte [cmd/seed/FATTORIA_IDS.md](cmd/seed/FATTORIA_IDS.md) para IDs
3. Execute `bash scripts/run_seed_fattoria.sh --clear-first`

### Para QA/Testes
1. Leia [INSTALLATION_CHECKLIST.md](INSTALLATION_CHECKLIST.md)
2. Valide todos os itens
3. Teste via API ou Frontend

### Para Apresentações
1. Veja [FATTORIA_MENU.txt](FATTORIA_MENU.txt)
2. Execute com dados limpos
3. Mostre dados reais do Fattoria

## 🎯 Características

✅ **Completamente Implementado**
- Dados de todos os produtos
- Estrutura operacional
- Tags de dieta
- Usuário admin
- Documentação completa

✅ **Fácil de Usar**
- 3 opções diferentes
- Script bash amigável
- Mensagens em português
- Suporte a flags

✅ **Bem Documentado**
- 5 arquivos de documentação
- Exemplos de código
- Checklist de validação
- FAQ e troubleshooting

✅ **Mantém Compatibilidade**
- Seed padrão continua funcionando
- Mesma estrutura e modelos
- Sem breaking changes

## 🆘 Troubleshooting

### "Failed to connect to database"
→ Verifique PostgreSQL e `.env`

### "Unknown option: --restaurant"
→ Use a versão modificada de `cmd/seed/main.go`

### "ParentId field not found"
→ Use a versão corrigida de `utils/seed_fattoria.go`

Para mais: Veja [INSTALLATION_CHECKLIST.md](INSTALLATION_CHECKLIST.md)

## 📞 Documentação de Referência

| Documento | Quando Usar | Tempo |
|-----------|------------|-------|
| **START_HERE.md** (aqui) | Começar agora | 2 min |
| **FATTORIA_MENU.txt** | Ver menu visual | 3 min |
| **cmd/seed/README_FATTORIA.md** | Quick start | 5 min |
| **SEED_FATTORIA.md** | Documentação completa | 20 min |
| **cmd/seed/FATTORIA_IDS.md** | Referência de IDs | 15 min |
| **SEED_FATTORIA_SUMMARY.md** | Resumo técnico | 10 min |
| **INSTALLATION_CHECKLIST.md** | Validar instalação | 30 min |

## 🎉 Pronto para Começar?

```bash
# Execute um destes comandos:

# Opção 1: Rápida (mantém dados existentes)
bash scripts/run_seed_fattoria.sh

# Opção 2: Limpa (recomendado para primeira vez)
bash scripts/run_seed_fattoria.sh --clear-first

# Opção 3: Detalhada (com logs)
bash scripts/run_seed_fattoria.sh --clear-first --verbose
```

Depois inicie o servidor:
```bash
go run main.go
```

E acesse:
```bash
# API
curl http://localhost:8080/health

# Frontend (se rodando)
http://localhost:5173
```

---

## 📋 Checklist Inicial

- [ ] Leu este arquivo
- [ ] Tem Go e PostgreSQL instalados
- [ ] Arquivo `.env` configurado
- [ ] Executou: `bash scripts/run_seed_fattoria.sh --clear-first`
- [ ] Iniciou servidor: `go run main.go`
- [ ] Testou API: `curl http://localhost:8080/health`
- [ ] Fez login com `admin@fattoria.com.br`

**Parabéns! Você está pronto para começar!** 🎉

---

**Versão**: 1.0
**Status**: ✅ Completo
**Última Atualização**: 2024
**Suporte**: Veja [SEED_FATTORIA.md](SEED_FATTORIA.md)
