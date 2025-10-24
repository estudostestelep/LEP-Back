# 🛠️ Setup Inicial - Guia Completo

Guia passo-a-passo para configurar o ambiente de desenvolvimento do LEP Backend.

## 📋 Pré-requisitos

### Sistema Operacional
- Linux, macOS ou Windows (com WSL2 ou Git Bash)

### Ferramentas Obrigatórias
- [ ] Go 1.16+
- [ ] PostgreSQL 12+
- [ ] Git
- [ ] Node.js (apenas para frontend)

### Ferramentas Opcionais
- [ ] Docker & Docker Compose (recomendado)
- [ ] Postman ou Insomnia (para testar API)
- [ ] DBeaver ou pgAdmin (para gerenciar BD)

---

## 🔧 Instalação de Dependências

### 1. Go
```bash
# Verificar se Go está instalado
go version

# Se não estiver, baixar em:
# https://golang.org/dl/

# Linux/macOS (via brew):
brew install go

# Windows: Baixar instalador
```

### 2. PostgreSQL

#### Via Docker (Recomendado)
```bash
# Iniciar PostgreSQL em Docker
docker run --name postgres \
  -e POSTGRES_PASSWORD=senha123 \
  -e POSTGRES_DB=lep_database \
  -p 5432:5432 \
  -d postgres:latest

# Verificar se está rodando
docker ps | grep postgres
```

#### macOS (Homebrew)
```bash
brew install postgresql
brew services start postgresql
```

#### Linux (Debian/Ubuntu)
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib

# Iniciar serviço
sudo systemctl start postgresql
```

#### Windows
- Baixar instalador em: https://www.postgresql.org/download/windows/
- Instalar seguindo o wizard
- Iniciar PostgreSQL via Services

### 3. Git
```bash
# Verificar se Git está instalado
git --version

# Se não, instalar de: https://git-scm.com/
```

---

## 📁 Setup do Projeto

### 1. Clonar Repositório
```bash
git clone <url-do-repositorio>
cd LEP/LEP-Back
```

### 2. Instalar Dependências Go
```bash
go mod tidy
```

### 3. Configurar Variáveis de Ambiente
```bash
# Copiar arquivo de exemplo
cp .env.example .env

# Editar .env com suas credenciais
nano .env
# ou
vim .env
# ou abrir em editor (VSCode, etc)
```

### 4. Conteúdo do .env
```env
# ===== DATABASE =====
DB_USER=postgres
DB_PASS=sua_senha_aqui
DB_NAME=lep_database
INSTANCE_UNIX_SOCKET=        # Deixar vazio para desenvolvimento local

# ===== AUTHENTICATION =====
JWT_SECRET_PRIVATE_KEY=sua_chave_privada_aqui
JWT_SECRET_PUBLIC_KEY=sua_chave_publica_aqui

# ===== NOTIFICATIONS (Opcional) =====
TWILIO_ACCOUNT_SID=seu_account_sid
TWILIO_AUTH_TOKEN=seu_auth_token
TWILIO_PHONE_NUMBER=+5511999999999

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=seu_email@gmail.com
SMTP_PASSWORD=sua_senha_app

# ===== FEATURE FLAGS =====
ENABLE_CRON_JOBS=true

# ===== STORAGE =====
STORAGE_TYPE=local
BASE_URL=http://localhost:8080
```

### 5. Criar Banco de Dados
```bash
# Via psql (recomendado)
psql -U postgres -c "CREATE DATABASE lep_database;"

# Verificar
psql -U postgres -l | grep lep_database
```

### 6. Testar Conexão
```bash
# Verificar se consegue conectar
psql -U postgres -d lep_database -c "SELECT version();"
```

---

## 🌱 Executar Seeds

### Opção A: LEP Demo (Recomendado para Começar)
```bash
bash scripts/run_seed.sh
```

### Opção B: Fattoria Pizzeria (Dados Reais)
```bash
bash tools/scripts/seed/run_seed_fattoria.sh --clear-first
```

### Opção C: Sem Seed (Banco Vazio)
```bash
# Pular os seeds acima
```

---

## 🚀 Iniciar Servidor

```bash
# Iniciar o servidor
go run main.go

# Esperado:
# [GIN-debug] Listening and serving HTTP on :8080
```

### Alterar Porta (opcional)
```bash
PORT=3000 go run main.go
```

---

## ✅ Validações

### 1. Server Running
```bash
# Testar saúde da API
curl http://localhost:8080/health

# Resposta esperada:
# {"status":"healthy"}
```

### 2. Database Connection
```bash
# Testar ping
curl http://localhost:8080/ping

# Resposta esperada:
# pong
```

### 3. Authentication
```bash
# Fazer login (se executou seed)
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@lep-demo.com","password":"password"}'

# Resposta esperada:
# {"token":"eyJhbGciOiJIUzI1NiI...","user":{...}}
```

### 4. Database Query
```bash
# Verificar se tabelas foram criadas
psql -U postgres -d lep_database \
  -c "\dt"

# Deve listar tabelas: organizations, projects, users, etc.
```

---

## 🎨 Setup Frontend (Opcional)

Se quer rodar o frontend junto:

```bash
# 1. Navegar para frontend
cd ../LEP-Front

# 2. Instalar dependências
npm install

# 3. Iniciar dev server
npm run dev

# Frontend rodará em: http://localhost:5173
```

---

## 🔧 Troubleshooting

### "Connection refused" (PostgreSQL)
```bash
# Verificar se PostgreSQL está rodando
# Docker:
docker ps | grep postgres

# macOS:
brew services list | grep postgresql

# Linux:
systemctl status postgresql

# Windows: Verificar Services (services.msc)
```

### "Erro de credenciais"
```bash
# Verificar credenciais no .env
cat .env | grep DB_

# Testar conexão manualmente
psql -U postgres -d lep_database
```

### "Porta 8080 ocupada"
```bash
# Usar porta diferente
PORT=8081 go run main.go

# Ou verificar qual processo está usando
lsof -i :8080
kill -9 <PID>
```

### "Tabelas não existem"
```bash
# Re-executar seed
bash scripts/run_seed.sh --clear-first

# Ou esperar auto-migrate da próxima inicialização
```

### "GORM auto-migration error"
```bash
# Deletar banco e criar novamente
dropdb -U postgres lep_database
createdb -U postgres lep_database

# Re-executar seed
bash scripts/run_seed.sh
```

---

## 📚 Documentação Relacionada

- [QUICKSTART.md](QUICKSTART.md) - Setup rápido (5 min)
- [docs_seeds/README.md](../docs_seeds/README.md) - Seeds disponíveis
- [frontend/FRONTEND_INSTRUCTIONS.md](frontend/FRONTEND_INSTRUCTIONS.md) - Frontend setup
- [deployment/DEPLOYMENT.md](deployment/DEPLOYMENT.md) - Deploy para produção

---

## 🎯 Próximas Etapas

### Após Setup Completo
1. [ ] Explorar API endpoints
2. [ ] Ler documentação de seeds ([docs_seeds/README.md](../docs_seeds/README.md))
3. [ ] Setup frontend (se necessário)
4. [ ] Revisar estrutura de código
5. [ ] Fazer seu primeiro commit

### Antes de Deploy
1. [ ] Revisar [deployment/DEPLOYMENT.md](deployment/DEPLOYMENT.md)
2. [ ] Configurar variáveis em [deployment/ENVIRONMENTS.md](deployment/ENVIRONMENTS.md)
3. [ ] Executar testes
4. [ ] Validar em staging

---

## ✅ Checklist Final

- [ ] Go 1.16+ instalado
- [ ] PostgreSQL 12+ instalado e rodando
- [ ] Dependências Go instaladas (`go mod tidy`)
- [ ] Arquivo .env criado e configurado
- [ ] Banco de dados criado
- [ ] Seed executado
- [ ] Server iniciado (`go run main.go`)
- [ ] API respondendo (`curl http://localhost:8080/health`)
- [ ] Login funcionando (se seed executado)

---

## 🚀 Você Está Pronto!

Setup concluído com sucesso. Agora você pode:

- Explorar a API
- Executar seeds
- Desenvolver features
- Consultar documentação

**Próximo passo?** Leia [QUICKSTART.md](QUICKSTART.md) para um overview rápido.

---

**Data**: 2024
**Status**: ✅ Guia Completo
**Última Atualização**: 2024
