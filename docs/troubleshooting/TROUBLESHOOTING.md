# 🔧 LEP System - Troubleshooting Guide

## 🚨 Problemas Comuns e Soluções

### **Erro: "unexpected character '/' in variable name"**

#### **Sintomas:**
```bash
failed to read .env: line 3: unexpected character "/" in variable name
```

#### **Causa:**
Docker Compose não suporta JWT keys multi-line no formato atual do arquivo `.env`.

#### **Solução:**
1. **Remover arquivo .env problemático:**
   ```bash
   rm .env
   ```

2. **Executar novamente o script:**
   ```bash
   ./scripts/deploy-interactive.sh
   ```

3. **Ou usar correção manual:**
   ```bash
   # Criar .env simples para desenvolvimento local
   cat > .env << EOF
   ENVIRONMENT=local-dev
   PORT=8080
   DB_HOST=postgres
   DB_PORT=5432
   DB_USER=lep_user
   DB_PASS=lep_password
   DB_NAME=lep_database
   DB_SSL_MODE=disable
   JWT_SECRET_PRIVATE_KEY=simple-local-key
   JWT_SECRET_PUBLIC_KEY=simple-local-key
   ENABLE_CRON_JOBS=false
   EOF
   ```

---

### **Docker Compose não encontrado**

#### **Sintomas:**
```bash
docker-compose: command not found
```

#### **Solução:**

**Windows (Docker Desktop):**
```bash
# Docker Desktop inclui docker-compose
# Se não funcionar, use:
docker compose up -d
```

**Linux:**
```bash
# Instalar docker-compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

**Mac:**
```bash
brew install docker-compose
```

---

### **Porta 8080 já em uso**

#### **Sintomas:**
```bash
ERROR: Port 8080 is already in use
```

#### **Solução:**
```bash
# Verificar o que está usando a porta
netstat -tulpn | grep :8080

# Parar outros serviços LEP
docker-compose down

# Ou mudar a porta no docker-compose.yml
# Alterar: "8080:8080" para "8081:8080"
```

---

### **Banco de dados não conecta**

#### **Sintomas:**
```bash
connection refused: postgres:5432
```

#### **Solução:**
```bash
# Verificar se PostgreSQL está rodando
docker-compose ps

# Reiniciar apenas o banco
docker-compose restart postgres

# Ver logs do banco
docker-compose logs postgres

# Aguardar o banco ficar pronto
docker-compose up -d postgres
sleep 30
docker-compose up -d app
```

---

### **Erro de permissão no Docker**

#### **Sintomas:**
```bash
permission denied while trying to connect to Docker daemon
```

#### **Solução:**

**Linux:**
```bash
# Adicionar usuário ao grupo docker
sudo usermod -aG docker $USER
newgrp docker

# Ou usar sudo
sudo docker-compose up -d
```

**Windows/Mac:**
```bash
# Verificar se Docker Desktop está rodando
# Restart Docker Desktop se necessário
```

---

### **Build falha por falta de espaço**

#### **Sintomas:**
```bash
no space left on device
```

#### **Solução:**
```bash
# Limpar containers e imagens não utilizadas
docker system prune -f

# Remover volumes não utilizados
docker volume prune -f

# Verificar espaço
docker system df
```

---

### **Go module download falha**

#### **Sintomas:**
```bash
go: module not found
```

#### **Solução:**
```bash
# No host (fora do container)
go mod tidy
go mod verify

# Rebuild sem cache
docker-compose build --no-cache app
```

---

## 🔄 Reset Completo

Se nada funcionar, faça um reset completo:

```bash
# Parar tudo
docker-compose down -v

# Remover imagens locais
docker rmi $(docker images lep-* -q) 2>/dev/null || true

# Limpar arquivos problemáticos
rm -f .env

# Reconstruir tudo
docker-compose build --no-cache
docker-compose up -d
```

---

## 📊 Verificação de Saúde

### **Verificar se tudo está funcionando:**

```bash
# Status dos containers
docker-compose ps

# Logs da aplicação
docker-compose logs -f app

# Teste de conectividade
curl http://localhost:8080/health

# Teste do banco
docker-compose exec postgres psql -U lep_user -d lep_database -c "SELECT version();"
```

### **URLs de verificação:**
- 🏠 **Aplicação**: http://localhost:8080
- ❤️ **Health Check**: http://localhost:8080/health
- 🏓 **Ping**: http://localhost:8080/ping
- 📧 **MailHog**: http://localhost:8025
- 🗄️ **PgAdmin**: http://localhost:5050

---

## 🆘 Ainda com problemas?

### **Coleta de informações para suporte:**

```bash
# Informações do sistema
echo "=== System Info ==="
uname -a
docker --version
docker-compose --version

# Status dos containers
echo "=== Container Status ==="
docker-compose ps

# Logs dos últimos 50 linhas
echo "=== App Logs ==="
docker-compose logs --tail=50 app

# Verificar .env
echo "=== Environment ==="
cat .env | head -10

# Verificar conectividade
echo "=== Network Test ==="
curl -v http://localhost:8080/health
```

### **Comandos de debug:**

```bash
# Entrar no container da aplicação
docker-compose exec app sh

# Entrar no banco de dados
docker-compose exec postgres psql -U lep_user -d lep_database

# Ver variáveis de ambiente do container
docker-compose exec app env | grep -E "(JWT|DB|PORT)"
```

---

## 📞 Quick Fixes

### **Desenvolvimento rápido sem Docker:**

Se Docker estiver com problemas, você pode rodar localmente:

```bash
# Instalar PostgreSQL local
# Windows: https://www.postgresql.org/download/windows/
# Mac: brew install postgresql
# Linux: sudo apt install postgresql

# Configurar banco local
createdb lep_database
psql lep_database -c "CREATE USER lep_user WITH PASSWORD 'lep_password';"

# Configurar .env local
cat > .env << EOF
ENVIRONMENT=local-dev
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=lep_user
DB_PASS=lep_password
DB_NAME=lep_database
DB_SSL_MODE=disable
JWT_SECRET_PRIVATE_KEY=simple-local-key
JWT_SECRET_PUBLIC_KEY=simple-local-key
ENABLE_CRON_JOBS=false
EOF

# Rodar aplicação
go run main.go
```