# 🚀 LEP Backend - Quick Start

Guia rápido para começar a trabalhar com o backend LEP.

## ⚡ 5 Minutos para Começar

### 1. Instalar Dependências
```bash
go mod tidy
```

### 2. Configurar Ambiente
```bash
# Copiar .env.example para .env
cp .env.example .env

# Editar .env com suas credenciais PostgreSQL
# DB_USER=seu_usuario
# DB_PASS=sua_senha
# DB_NAME=seu_banco
```

### 3. Escolher Seed
```bash
# OPÇÃO A: Seed Padrão LEP Demo
bash scripts/run_seed.sh

# OPÇÃO B: Seed Fattoria Pizzeria (novo!)
bash scripts/run_seed_fattoria.sh --clear-first

# OPÇÃO C: Nenhum seed (banco vazio)
# Pule este passo
```

### 4. Iniciar Servidor
```bash
go run main.go
```

Servidor rodando em: `http://localhost:8080`

### 5. Testar API
```bash
# Health check
curl http://localhost:8080/health

# Login (com seed padrão)
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@lep-demo.com","password":"password"}'
```

## 📚 Documentação Completa

| Tópico | Arquivo | Tempo |
|--------|---------|-------|
| **Setup Inicial** | [docs/SETUP.md](SETUP.md) | 10 min |
| **Seeds** | [docs_seeds/README.md](../docs_seeds/README.md) | 5 min |
| **Deployment** | [docs/deployment/DEPLOYMENT.md](deployment/DEPLOYMENT.md) | 15 min |
| **Infraestrutura** | [docs/infra/INFRASTRUCTURE_AUDIT.md](infra/INFRASTRUCTURE_AUDIT.md) | 20 min |
| **Frontend** | [docs/frontend/FRONTEND_INSTRUCTIONS.md](frontend/FRONTEND_INSTRUCTIONS.md) | 10 min |

## 🔐 Credenciais Padrão (com seed padrão)

| Email | Senha | Rol |
|-------|-------|-----|
| admin@lep-demo.com | password | Admin |
| garcom@lep-demo.com | password | Waiter |
| gerente@lep-demo.com | password | Manager |

## 🍕 Credenciais Fattoria (com seed Fattoria)

| Email | Senha | Rol |
|-------|-------|-----|
| admin@fattoria.com.br | password | Admin |

## 🛠️ Troubleshooting

### PostgreSQL não está rodando
```bash
# Docker
docker-compose up -d db

# ou manualmente
brew services start postgresql  # macOS
systemctl start postgresql      # Linux
```

### Erro de porta ocupada (8080)
```bash
# Verificar processo
lsof -i :8080

# Ou usar porta diferente
PORT=8081 go run main.go
```

### Erro de conexão com BD
Verificar:
- PostgreSQL está rodando: `psql --version`
- Credenciais em `.env`: `DB_USER`, `DB_PASS`, `DB_NAME`
- Banco existe: `createdb seu_banco_name`

## 📂 Estrutura Principais Pastas

- `config/` - Configurações da aplicação
- `handler/` - Lógica de negócio
- `repositories/` - Acesso ao banco de dados
- `routes/` - Definição de rotas
- `utils/` - Funções utilitárias
- `cmd/seed/` - Seeds de banco de dados
- `scripts/` - Scripts de automação

## 🚀 Próximos Passos

1. Leia [docs/SETUP.md](SETUP.md) para configuração detalhada
2. Escolha e execute um seed em [docs_seeds/README.md](../docs_seeds/README.md)
3. Faça chamadas à API:
   - `GET /health` - Status do servidor
   - `POST /login` - Autenticação
   - `GET /product` - Listar produtos
   - `GET /table` - Listar mesas

## 🔗 Links Úteis

- [README Principal](../README.md)
- [Documentação de Seeds](../docs_seeds/README.md)
- [Instruções Frontend](frontend/FRONTEND_INSTRUCTIONS.md)
- [Deployment](deployment/DEPLOYMENT.md)

---

**Pronto para começar?** Execute: `go run main.go`
