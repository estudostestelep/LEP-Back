# 📖 Documentação LEP Backend

Documentação completa do backend LEP System.

## 📋 Índice de Documentação

### 🚀 Começar
1. **[QUICKSTART.md](QUICKSTART.md)** - Start em 5 minutos
2. **[SETUP.md](SETUP.md)** - Setup detalhado do ambiente

### 📦 Seeds & Banco de Dados
- **[docs_seeds/README.md](../docs_seeds/README.md)** - Seeds disponíveis
  - LEP Demo (padrão)
  - Fattoria Pizzeria (novo!)

### 🚀 Deployment & Produção
- **[deployment/README.md](deployment/README.md)**
  - [DEPLOYMENT.md](deployment/DEPLOYMENT.md) - Guia de deployment
  - [ENVIRONMENTS.md](deployment/ENVIRONMENTS.md) - Configuração de ambientes

### 🏗️ Infraestrutura
- **[infra/README.md](infra/README.md)**
  - [INFRASTRUCTURE_AUDIT.md](infra/INFRASTRUCTURE_AUDIT.md) - Auditoria de infra
  - [CHANGELOG_INFRA.md](infra/CHANGELOG_INFRA.md) - Histórico de mudanças

### 🎨 Frontend & Integração
- **[frontend/README.md](frontend/README.md)**
  - [FRONTEND_INSTRUCTIONS.md](frontend/FRONTEND_INSTRUCTIONS.md) - Instruções frontend
  - [FRONTEND_MIGRATION_GUIDE.md](frontend/FRONTEND_MIGRATION_GUIDE.md) - Guia de migração
  - [FRONTEND_CARDAPIO_PROMPT.md](frontend/FRONTEND_CARDAPIO_PROMPT.md) - Cardápio digital

---

## 🎯 Fluxo de Trabalho Recomendado

### 1️⃣ Primeira Vez (15 minutos)
```
QUICKSTART.md → SETUP.md → go run main.go
```

### 2️⃣ Usar Seeds (5 minutos)
```
docs_seeds/README.md → escolher seed → bash script
```

### 3️⃣ Desenvolvimento (Contínuo)
```
Referência: FRONTEND_INSTRUCTIONS.md + QUICKSTART.md
```

### 4️⃣ Deployment (20 minutos)
```
deployment/DEPLOYMENT.md → ENVIRONMENTS.md → Deploy!
```

---

## 📁 Estrutura de Pastas

```
docs/
├── README.md (este arquivo)
├── QUICKSTART.md (⭐ comece aqui)
├── SETUP.md
├── deployment/
│   ├── README.md
│   ├── DEPLOYMENT.md
│   └── ENVIRONMENTS.md
├── infra/
│   ├── README.md
│   ├── INFRASTRUCTURE_AUDIT.md
│   └── CHANGELOG_INFRA.md
└── frontend/
    ├── README.md
    ├── FRONTEND_INSTRUCTIONS.md
    ├── FRONTEND_MIGRATION_GUIDE.md
    └── FRONTEND_CARDAPIO_PROMPT.md

docs_seeds/
├── README.md (índice de seeds)
└── fattoria/
    ├── START_HERE.md
    ├── SEED_FATTORIA.md
    ├── ... (6 arquivos)
```

---

## 🔍 Procurando por...?

| Preciso de... | Arquivo | Tempo |
|---------------|---------|-------|
| Iniciar rápido | [QUICKSTART.md](QUICKSTART.md) | 5 min |
| Setup completo | [SETUP.md](SETUP.md) | 10 min |
| Dados de demo | [docs_seeds/README.md](../docs_seeds/README.md) | 5 min |
| Fazer deploy | [deployment/DEPLOYMENT.md](deployment/DEPLOYMENT.md) | 20 min |
| Entender infra | [infra/INFRASTRUCTURE_AUDIT.md](infra/INFRASTRUCTURE_AUDIT.md) | 15 min |
| Desenvolver frontend | [frontend/FRONTEND_INSTRUCTIONS.md](frontend/FRONTEND_INSTRUCTIONS.md) | 10 min |

---

## ⚡ Comandos Rápidos

```bash
# Setup básico
go mod tidy
cp .env.example .env
# (editar .env com credenciais)

# Executar com seed
bash tools/scripts/seed/run_seed_fattoria.sh --clear-first
go run main.go

# Testar API
curl http://localhost:8080/health

# Frontend (em outra pasta)
cd ../LEP-Front
npm install
npm run dev
```

---

## 🤔 FAQ Rápido

### Não consigo iniciar o servidor
→ Ver [SETUP.md](SETUP.md) - Seção PostgreSQL

### Qual seed devo usar?
→ Ver [docs_seeds/README.md](../docs_seeds/README.md)

### Como faço deploy?
→ Ver [deployment/DEPLOYMENT.md](deployment/DEPLOYMENT.md)

### Onde configuro variáveis?
→ Ver [deployment/ENVIRONMENTS.md](deployment/ENVIRONMENTS.md)

### Como integro com frontend?
→ Ver [frontend/FRONTEND_INSTRUCTIONS.md](frontend/FRONTEND_INSTRUCTIONS.md)

---

## 🔗 Links Úteis

- [README Principal](../README.md) - Visão geral do projeto
- [cmd/seed/main.go](../cmd/seed/main.go) - Código do seeder
- [CLAUDE.md](../CLAUDE.md) - Instruções para Claude AI

---

## 📞 Estrutura de Suporte

Cada pasta tem seu próprio README com informações específicas:
- **docs/** - Documentação geral
- **docs_seeds/** - Seeds específicos
- **docs/deployment/** - Deploy e produção
- **docs/infra/** - Infraestrutura
- **docs/frontend/** - Frontend e integração

---

## ✅ Checklist de Primeira Vez

- [ ] Li QUICKSTART.md
- [ ] Executei setup (SETUP.md)
- [ ] Instalei dependências (`go mod tidy`)
- [ ] Configurei .env
- [ ] Executei um seed
- [ ] Iniciei servidor (`go run main.go`)
- [ ] Testei API (`curl http://localhost:8080/health`)
- [ ] Explorei documentação específica

---

**Pronto para começar?** → Abra [QUICKSTART.md](QUICKSTART.md)
