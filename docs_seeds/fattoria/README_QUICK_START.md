# 🍕 Seed Fattoria - Quick Start

Guia rápido para usar o seed da Fattoria Pizzeria.

## ⚡ Começar em 30 segundos

```bash
# 1. Acesse a raiz do backend
cd LEP-Back

# 2. Execute o seed da Fattoria
bash scripts/run_seed_fattoria.sh --clear-first

# 3. Inicie o servidor
go run main.go

# 4. Login
# Email: admin@fattoria.com.br
# Senha: password
```

## 📋 O que você vai ter

### 🍕 Menu Completo
- **5 Pizzas**: Crostini, Marguerita, Marinara, Parma, Vegana
- **4 Bebidas**: Suco de caju, Heineken s/álc, Baden Baden IPA, Coquetel Sônia e Zé

### 🏪 Estrutura Operacional
- **3 Mesas**: Prontas para receber pedidos
- **1 Ambiente**: Salão Principal (60 pessoas)
- **Tags**: Vegetariana, Vegana (para filtros)

### 🔐 Acesso
- **Email**: admin@fattoria.com.br
- **Senha**: password

## 🎯 Comandos Úteis

```bash
# Seed com limpeza de dados antigos
bash scripts/run_seed_fattoria.sh --clear-first

# Seed mantendo dados existentes
bash scripts/run_seed_fattoria.sh

# Seed com saída detalhada
bash scripts/run_seed_fattoria.sh --verbose

# Seed em ambiente de teste
bash scripts/run_seed_fattoria.sh --environment=test

# Ver ajuda completo
bash scripts/run_seed_fattoria.sh --help
```

## 🚀 Go Direto (sem script)

```bash
# Seed padrão Fattoria
go run cmd/seed/main.go --restaurant=fattoria

# Com opções
go run cmd/seed/main.go --restaurant=fattoria --clear-first --verbose
```

## 📊 Dados Inclusos

| Categoria | Item | Preço | Tags |
|-----------|------|-------|------|
| **Entradas** | Crostini | R$ 30,00 | - |
| **Pizzas** | Marguerita | R$ 80,00 | 🌱 Vegetariana |
| **Pizzas** | Marinara | R$ 58,00 | 🌿 Vegana |
| **Pizzas** | Parma | R$ 109,00 | - |
| **Pizzas** | Vegana | R$ 60,00 | 🌿 Vegana |
| **Soft** | Suco Caju | R$ 15,00 | - |
| **Cervejas** | Heineken | R$ 13,00 | - |
| **Artesanal** | Baden Baden | R$ 23,00 | - |
| **Coquetéis** | Sônia e Zé | R$ 32,00 | - |

## 🔗 Documentação Completa

Para mais detalhes, consulte [SEED_FATTORIA.md](../../SEED_FATTORIA.md)

## 💡 Tips

- Use `--clear-first` apenas na primeira execução
- Para desenvolvimento rápido, rode sem `--clear-first`
- Use `--verbose` se houver erros
- Todos os comandos suportam multiple flags

## 🆘 Troubleshooting

```bash
# PostgreSQL não está rodando?
# Inicie em Docker: docker-compose up -d db

# Erro de conexão?
# Verifique .env com credenciais corretas

# Precisa de detalhes?
bash scripts/run_seed_fattoria.sh --verbose
```

---

✅ **Pronto para começar!** Execute `bash scripts/run_seed_fattoria.sh --clear-first`
