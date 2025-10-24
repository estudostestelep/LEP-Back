# 📦 LEP Database Seeds

Documentação de todos os seeds disponíveis para popular o banco de dados.

## 🚀 Seeds Disponíveis

### 1. 🏢 Seed Padrão LEP Demo
**Localização**: `scripts/run_seed.sh`

Seed padrão com dados de demonstração do LEP:
- 3 usuários de teste
- 6 produtos de exemplo
- 2 mesas
- 1 reserva e 1 waitlist
- 1 pedido em preparo

**Executar**:
```bash
bash scripts/run_seed.sh
```

**Login padrão**:
- Email: `admin@lep-demo.com`
- Senha: `password`

---

### 2. 🍕 Seed Fattoria Pizzeria ⭐ (NOVO!)
**Localização**: `docs_seeds/fattoria/`

Seed completo da Fattoria Pizzeria com:
- **9 Produtos**: 5 pizzas + 4 bebidas
- **8 Categorias**: Organizado por tipo
- **3 Mesas**: Pronto para operação
- **2 Tags**: Vegetariana e Vegana
- **1 Admin**: Usuário específico da Fattoria

**Documentação Fattoria**:
- [START_HERE.md](fattoria/START_HERE.md) - Comece aqui (5 min)
- [SEED_FATTORIA.md](fattoria/SEED_FATTORIA.md) - Documentação completa (20 min)
- [FATTORIA_MENU.txt](fattoria/FATTORIA_MENU.txt) - Menu visual em ASCII
- [INSTALLATION_CHECKLIST.md](fattoria/INSTALLATION_CHECKLIST.md) - Validação passo-a-passo
- [FATTORIA_IDS.md](fattoria/FATTORIA_IDS.md) - Referência de IDs
- [FILES_MANIFEST.md](fattoria/FILES_MANIFEST.md) - Lista de arquivos

**Executar**:
```bash
# Opção 1: Via script (recomendado)
bash tools/scripts/seed/run_seed_fattoria.sh --clear-first

# Opção 2: Via Go direto
go run cmd/seed/main.go --restaurant=fattoria --clear-first

# Opção 3: Via script original (compatibilidade)
bash scripts/run_seed_fattoria.sh --clear-first
```

**Login Fattoria**:
- Email: `admin@fattoria.com.br`
- Senha: `password`

---

## 📋 Comparação de Seeds

| Aspecto | LEP Demo | Fattoria |
|---------|----------|----------|
| **Produtos** | 6 itens genéricos | 9 itens reais |
| **Pizzas** | - | 5 pizzas |
| **Bebidas** | 1 bebida | 4 bebidas |
| **Mesas** | 2 mesas | 3 mesas |
| **Tags** | - | Vegetariana, Vegana |
| **Docs** | README padrão | 8 arquivos |
| **Admin** | admin@lep-demo.com | admin@fattoria.com.br |

---

## 🎯 Como Escolher?

### Use **LEP Demo** se:
- [ ] Quer dados genéricos para testes
- [ ] Precisa de estrutura simples
- [ ] Está começando a explorar o sistema
- [ ] Quer exemplo padrão

### Use **Fattoria** se:
- [ ] Quer dados realistas de um restaurante
- [ ] Precisa testar funcionalidades completas
- [ ] Vai fazer uma apresentação
- [ ] Quer entender como estruturar um restaurante
- [ ] Precisa de documentação detalhada

---

## 🔄 Alternando Entre Seeds

```bash
# 1. Limpar dados atuais
bash scripts/run_seed.sh --clear-first   # LEP Demo
# ou
bash tools/scripts/seed/run_seed_fattoria.sh --clear-first  # Fattoria

# 2. Executar novo seed
go run main.go
```

---

## 📂 Estrutura de Seeds

```
docs_seeds/
├── README.md (este arquivo)
└── fattoria/
    ├── START_HERE.md ⭐ (comece aqui)
    ├── SEED_FATTORIA.md (completo)
    ├── SEED_FATTORIA_SUMMARY.md
    ├── INSTALLATION_CHECKLIST.md
    ├── FATTORIA_MENU.txt
    ├── FATTORIA_IDS.md
    ├── README_QUICK_START.md
    └── FILES_MANIFEST.md
```

---

## 🛠️ Scripts de Seed

### LEP Demo (padrão)
```bash
# Localização: scripts/run_seed.sh
bash scripts/run_seed.sh
bash scripts/run_seed.sh --clear-first
bash scripts/run_seed.sh --verbose
bash scripts/run_seed.sh --environment=test
```

### Fattoria (novo)
```bash
# Localização: tools/scripts/seed/run_seed_fattoria.sh
bash tools/scripts/seed/run_seed_fattoria.sh
bash tools/scripts/seed/run_seed_fattoria.sh --clear-first
bash tools/scripts/seed/run_seed_fattoria.sh --verbose
bash tools/scripts/seed/run_seed_fattoria.sh --environment=test
```

---

## 💡 Flags Comuns

| Flag | Efeito | Exemplo |
|------|--------|---------|
| `--clear-first` | Limpa todos os dados antes de popular | `--clear-first` |
| `--verbose` | Saída detalhada de debug | `--verbose` |
| `--environment` | Seleciona ambiente | `--environment=test` |
| `--help` | Mostra ajuda | `--help` |

---

## ✅ Validar Seed Executado

```bash
# 1. Testar saúde da API
curl http://localhost:8080/health

# 2. Fazer login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@lep-demo.com","password":"password"}'

# 3. Listar produtos
curl -X GET http://localhost:8080/product \
  -H "Authorization: Bearer {TOKEN}" \
  -H "X-Lpe-Organization-Id: {ORG_ID}" \
  -H "X-Lpe-Project-Id: {PROJ_ID}"
```

---

## 🔗 Documentação Relacionada

- [docs/QUICKSTART.md](../docs/QUICKSTART.md) - Quick start geral
- [docs/SETUP.md](../docs/SETUP.md) - Setup detalhado
- [docs_seeds/fattoria/START_HERE.md](fattoria/START_HERE.md) - Guia Fattoria
- [cmd/seed/main.go](../cmd/seed/main.go) - Código do seeder

---

## 🚀 Começar Agora

### Opção 1: LEP Demo (rápido)
```bash
bash scripts/run_seed.sh
go run main.go
```

### Opção 2: Fattoria (completo)
```bash
bash tools/scripts/seed/run_seed_fattoria.sh --clear-first
go run main.go
```

### Opção 3: Nenhum seed (vazio)
```bash
# Pule os seeds e inicie com banco vazio
go run main.go
```

---

**Dúvida?** Consulte a documentação específica do seed que está usando.
