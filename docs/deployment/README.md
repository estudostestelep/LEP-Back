# 🚀 Deployment & Ambientes

Documentação sobre deployment, ambientes e configurações de infraestrutura para produção.

## 📚 Arquivos

### [DEPLOYMENT.md](DEPLOYMENT.md)
Guia completo de deployment para diferentes ambientes:
- Configuração de produção
- Deploy em GCP (Cloud Run, Cloud SQL)
- Scripts de deployment
- Validação pós-deployment

### [ENVIRONMENTS.md](ENVIRONMENTS.md)
Configuração de ambientes:
- Variáveis de ambiente
- Configurações por ambiente (dev, test, staging, prod)
- Secrets management
- Database setup

---

## 🎯 Guia Rápido

### Preparar para Deploy
```bash
# 1. Revisar DEPLOYMENT.md
cat docs/deployment/DEPLOYMENT.md

# 2. Configurar variáveis em ENVIRONMENTS.md
cat docs/deployment/ENVIRONMENTS.md

# 3. Executar validações
go test ./...

# 4. Build
go build -o lep-system .
```

### Fazer Deploy
```bash
# Seguir instruções em DEPLOYMENT.md
# Variáveis definidas em ENVIRONMENTS.md
```

---

## 🔗 Links Relacionados
- [docs/QUICKSTART.md](../QUICKSTART.md) - Quick start geral
- [docs/infra/INFRASTRUCTURE_AUDIT.md](../infra/INFRASTRUCTURE_AUDIT.md) - Auditoria de infra
- [../README.md](../../README.md) - README principal

---

**Comece por**: [DEPLOYMENT.md](DEPLOYMENT.md)
