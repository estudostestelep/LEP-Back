# 🌱 Seeds - Banco de Dados

Documentação sobre como o sistema de seeds funciona e como as organizações são criadas.

## 📚 Arquivos

### [SEED_ARCHITECTURE.md](SEED_ARCHITECTURE.md) ⭐ **LEIA PRIMEIRO**
Explicação completa da arquitetura de seeds:
- ✅ Como seed cria nova organização
- ✅ Acesso de usuários entre orgs
- ✅ Sistema de IDs dinâmicos
- ✅ Fluxo completo de bootstrap
- ✅ Casos de uso práticos

---

## 🚀 Quick Answer

### "Seed de Fattoria cria nova organização?"
**SIM** ✅

### "Master admin anterior tem acesso?"
**NÃO automático** ⚠️
- Cada seed cria seu próprio usuário admin
- Master admin precisa de relacionamento adicional (UserOrganization)

### "Posso ter múltiplas orgs?"
**SIM** ✅
- Pode rodar múltiplos seeds
- Cada um cria nova org
- Acessa via diferentes headers

---

## 🔗 Links Relacionados

- [docs_seeds/README.md](../../docs_seeds/README.md) - Seeds disponíveis
- [docs_seeds/fattoria/START_HERE.md](../../docs_seeds/fattoria/START_HERE.md) - Seed Fattoria
- [docs/SETUP.md](../SETUP.md) - Setup inicial

---

**Para entender melhor**: Leia [SEED_ARCHITECTURE.md](SEED_ARCHITECTURE.md)
