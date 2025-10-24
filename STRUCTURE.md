# 📁 Estrutura do Projeto LEP-Back

Mapa completo de pastas e arquivos do backend LEP.

## 🏠 Raiz do Projeto

```
LEP-Back/
├── 📄 README.md ⭐ (Documentação principal do projeto)
├── 📄 CLAUDE.md (Instruções para Claude AI)
├── 📄 STRUCTURE.md (este arquivo)
│
├── 📁 docs/ (Documentação)
│   ├── 📄 README.md (Índice principal de docs)
│   ├── 📄 QUICKSTART.md (Start em 5 minutos)
│   ├── 📄 SETUP.md (Setup detalhado)
│   │
│   ├── 📁 deployment/ (Deploy & Ambientes)
│   │   ├── 📄 README.md
│   │   ├── 📄 DEPLOYMENT.md
│   │   └── 📄 ENVIRONMENTS.md
│   │
│   ├── 📁 infra/ (Infraestrutura)
│   │   ├── 📄 README.md
│   │   ├── 📄 INFRASTRUCTURE_AUDIT.md
│   │   └── 📄 CHANGELOG_INFRA.md
│   │
│   ├── 📁 frontend/ (Frontend & Integração)
│   │   ├── 📄 README.md
│   │   ├── 📄 FRONTEND_INSTRUCTIONS.md
│   │   ├── 📄 FRONTEND_MIGRATION_GUIDE.md
│   │   └── 📄 FRONTEND_CARDAPIO_PROMPT.md
│   │
│   └── 📁 guides/ (Guias & Referências)
│       ├── 📄 README.md
│       ├── 📄 CHANGELOG.md
│       └── 📄 IMPLEMENTACAO_CARDAPIO.md
│
├── 📁 docs_seeds/ (Seeds de Banco de Dados)
│   ├── 📄 README.md (Índice de seeds)
│   │
│   └── 📁 fattoria/ (Seed Fattoria Pizzeria)
│       ├── 📄 START_HERE.md ⭐
│       ├── 📄 SEED_FATTORIA.md
│       ├── 📄 SEED_FATTORIA_SUMMARY.md
│       ├── 📄 INSTALLATION_CHECKLIST.md
│       ├── 📄 README_QUICK_START.md
│       ├── 📄 FATTORIA_IDS.md
│       ├── 📄 FATTORIA_MENU.txt
│       └── 📄 FILES_MANIFEST.md
│
├── 📁 tools/ (Ferramentas & Scripts)
│   ├── 📁 scripts/ (Scripts de Automação)
│   │   ├── 📁 seed/ (Scripts de Seed)
│   │   │   └── 📜 run_seed_fattoria.sh
│   │   ├── 📁 deploy/ (Scripts de Deploy)
│   │   └── 📁 dev/ (Scripts de Desenvolvimento)
│   │
│   └── (futuro: mais ferramentas)
│
├── 📁 cmd/ (Comandos & Binários)
│   ├── 📁 seed/ (Seeder da aplicação)
│   │   ├── main.go
│   │   ├── bootstrap_helpers.go
│   │   └── README_FATTORIA.md (documentação específica)
│   │
│   └── (outros comandos)
│
├── 📁 config/ (Configurações)
├── 📁 handler/ (Lógica de Negócio)
├── 📁 middleware/ (Middlewares HTTP)
├── 📁 repositories/ (Acesso a Dados)
│   └── 📁 models/ (Modelos de Dados)
├── 📁 routes/ (Definição de Rotas)
├── 📁 utils/ (Funções Utilitárias)
│   ├── seed_data.go (Seed padrão LEP)
│   ├── seed_fattoria.go (Seed Fattoria)
│   ├── event_service.go
│   ├── cron_service.go
│   ├── notification_service.go
│   └── (outras utilidades)
│
├── 📁 scripts/ (Scripts Originais)
│   ├── run_seed.sh (Seed padrão LEP)
│   ├── run_tests.sh
│   ├── dev-local.sh
│   ├── stage-local.sh
│   ├── stage-deploy.sh
│   └── (outros scripts originais)
│
├── 📁 resource/ (Gerenciamento de Recursos)
├── 📁 tests/ (Testes Automatizados)
├── 📁 migrations/ (Migrações de BD)
│
├── 📄 go.mod (Dependências Go)
├── 📄 go.sum
├── 📄 .env.example
├── 📄 .gitignore
├── 📄 main.go (Entrada da aplicação)
├── 📄 Dockerfile
└── 📄 docker-compose.yml
```

---

## 📚 Documentação por Localização

### Na Raiz (⭐ Principais)
| Arquivo | Propósito |
|---------|-----------|
| **README.md** | Documentação principal do projeto |
| **CLAUDE.md** | Instruções para Claude AI |
| **STRUCTURE.md** | Este arquivo - estrutura de pastas |

### Em `docs/`
| Arquivo | Propósito |
|---------|-----------|
| **README.md** | Índice de toda documentação |
| **QUICKSTART.md** | Começar em 5 minutos |
| **SETUP.md** | Setup detalhado |

### Em `docs/deployment/`
| Arquivo | Propósito |
|---------|-----------|
| **README.md** | Índice de deploy |
| **DEPLOYMENT.md** | Guia de deployment |
| **ENVIRONMENTS.md** | Configuração de ambientes |

### Em `docs/infra/`
| Arquivo | Propósito |
|---------|-----------|
| **README.md** | Índice de infraestrutura |
| **INFRASTRUCTURE_AUDIT.md** | Auditoria de infra |
| **CHANGELOG_INFRA.md** | Histórico de mudanças |

### Em `docs/frontend/`
| Arquivo | Propósito |
|---------|-----------|
| **README.md** | Índice de frontend |
| **FRONTEND_INSTRUCTIONS.md** | Setup frontend |
| **FRONTEND_MIGRATION_GUIDE.md** | Guia de migração |
| **FRONTEND_CARDAPIO_PROMPT.md** | Cardápio digital |

### Em `docs/guides/`
| Arquivo | Propósito |
|---------|-----------|
| **README.md** | Índice de guias |
| **CHANGELOG.md** | Histórico do projeto |
| **IMPLEMENTACAO_CARDAPIO.md** | Implementação de cardápio |

### Em `docs_seeds/`
| Arquivo | Propósito |
|---------|-----------|
| **README.md** | Índice de seeds |

### Em `docs_seeds/fattoria/`
| Arquivo | Propósito |
|---------|-----------|
| **START_HERE.md** | Comece aqui com Fattoria |
| **SEED_FATTORIA.md** | Documentação completa |
| **INSTALLATION_CHECKLIST.md** | Validação passo-a-passo |
| **FATTORIA_IDS.md** | Referência de IDs |
| **FATTORIA_MENU.txt** | Menu visual |
| **FILES_MANIFEST.md** | Lista de arquivos |

---

## 🗂️ Organização de Scripts

```
tools/scripts/
├── seed/
│   └── run_seed_fattoria.sh (Seed da Fattoria)
├── deploy/
│   └── (scripts de deployment)
└── dev/
    └── (scripts de desenvolvimento)

scripts/ (Scripts originais - mantidos para compatibilidade)
├── run_seed.sh (Seed padrão LEP)
├── run_tests.sh
├── dev-local.sh
└── stage-deploy.sh
```

---

## 🎯 Fluxo de Navegação

### Primeira Vez?
```
README.md → docs/QUICKSTART.md → docs/SETUP.md → go run main.go
```

### Quer Dados?
```
docs_seeds/README.md → docs_seeds/fattoria/ → bash tools/scripts/seed/run_seed_fattoria.sh
```

### Quer Deploy?
```
docs/deployment/DEPLOYMENT.md → docs/deployment/ENVIRONMENTS.md → Deploy!
```

### Desenvolver Frontend?
```
docs/frontend/FRONTEND_INSTRUCTIONS.md → LEP-Front/README.md → npm run dev
```

### Entender Infra?
```
docs/infra/INFRASTRUCTURE_AUDIT.md → docs/infra/CHANGELOG_INFRA.md
```

---

## 📊 Estatísticas

- **Arquivos .md**: ~20+ (bem organizados)
- **Scripts**: 2+ ativos (tools/scripts/ e scripts/)
- **Documentação**: 7 pastas principais
- **Seeds**: 2 (LEP Demo + Fattoria)

---

## 🔍 Buscar Rapidamente

| Preciso... | Procure em... |
|-----------|---------------|
| Começar rápido | `docs/QUICKSTART.md` |
| Setup completo | `docs/SETUP.md` |
| Seed Fattoria | `docs_seeds/fattoria/START_HERE.md` |
| Deploy | `docs/deployment/DEPLOYMENT.md` |
| Frontend | `docs/frontend/FRONTEND_INSTRUCTIONS.md` |
| Infra | `docs/infra/INFRASTRUCTURE_AUDIT.md` |
| Histórico | `docs/guides/CHANGELOG.md` |
| Cardápio | `docs/frontend/FRONTEND_CARDAPIO_PROMPT.md` |

---

## ✅ Organização Completa

- ✅ Documentação centralizada em `docs/`
- ✅ Seeds específicos em `docs_seeds/`
- ✅ Scripts em `tools/scripts/`
- ✅ Raiz limpa (apenas README.md e CLAUDE.md)
- ✅ READMEs em cada pasta para navegação
- ✅ Estrutura lógica e intuitiva

---

**Última Atualização**: 2024
**Status**: ✅ Reorganização Completa
