#!/bin/bash
# LEP System - Local Development Environment (100% local)
# Este script inicia o ambiente de desenvolvimento completamente local com Docker

set -e

echo "=== LEP System - Ambiente DEV Local ==="
echo "• Docker PostgreSQL + Redis + MailHog"
echo "• localStorage para imagens"
echo "• Credenciais padronizadas dev/stage"
echo

# Verificar se Docker está rodando
if ! docker info >/dev/null 2>&1; then
    echo "❌ Erro: Docker não está rodando. Por favor, inicie o Docker primeiro."
    exit 1
fi

# Verificar se docker-compose está disponível
if ! command -v docker-compose >/dev/null 2>&1; then
    echo "❌ Erro: docker-compose não encontrado. Por favor, instale o Docker Compose."
    exit 1
fi

# Parar containers existentes (se houver)
echo "🧹 Parando containers existentes..."
docker-compose down --remove-orphans >/dev/null 2>&1 || true

# Construir e iniciar serviços
echo "🔨 Construindo containers..."
docker-compose build --no-cache

echo "🚀 Iniciando ambiente DEV local..."
docker-compose up -d postgres redis mailhog

# Aguardar PostgreSQL ficar pronto
echo "⏳ Aguardando PostgreSQL ficar pronto..."
timeout=60
while ! docker-compose exec -T postgres pg_isready -U lep_user -d lep_database >/dev/null 2>&1; do
    sleep 2
    timeout=$((timeout - 2))
    if [ $timeout -le 0 ]; then
        echo "❌ Timeout aguardando PostgreSQL. Verifique os logs: docker-compose logs postgres"
        exit 1
    fi
done

echo "✅ PostgreSQL pronto!"

# Iniciar aplicação principal
echo "🚀 Iniciando aplicação LEP..."
docker-compose up -d app

# Aguardar aplicação ficar pronta
echo "⏳ Aguardando aplicação ficar pronta..."
sleep 10

# Verificar se aplicação está respondendo
if curl -f -s http://localhost:8080/ping >/dev/null; then
    echo "✅ Aplicação rodando em http://localhost:8080"
else
    echo "❌ Aplicação não está respondendo. Verifique os logs: docker-compose logs app"
    exit 1
fi

echo
echo "=== Ambiente DEV Local Iniciado com Sucesso! ==="
echo "🌐 API: http://localhost:8080"
echo "📧 MailHog: http://localhost:8025 (SMTP interface)"
echo "🗄️  PostgreSQL: localhost:5432 (lep_user/lep_password/lep_database)"
echo "🔴 Redis: localhost:6379"
echo
echo "Comandos úteis:"
echo "  docker-compose logs app          # Ver logs da aplicação"
echo "  docker-compose logs postgres     # Ver logs do PostgreSQL"
echo "  docker-compose exec app bash     # Acessar container da aplicação"
echo "  docker-compose down              # Parar tudo"
echo
echo "Para popular com dados de exemplo:"
echo "  docker-compose run --rm seed"
echo