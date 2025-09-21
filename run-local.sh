#!/bin/bash

echo ""
echo "=============================================="
echo " LEP System Backend - Local Development"
echo "=============================================="
echo ""

echo "Verificando se o Go está instalado..."
if ! command -v go &> /dev/null; then
    echo "ERRO: Go não está instalado ou não está no PATH"
    echo "Baixe o Go em: https://golang.org/dl/"
    exit 1
fi

go version

echo ""
echo "Baixando dependências..."
go mod tidy
if [ $? -ne 0 ]; then
    echo "AVISO: Falha ao baixar dependências, continuando..."
fi

echo ""
echo "Verificando build..."
go build .
if [ $? -ne 0 ]; then
    echo "ERRO: Falha no build da aplicação"
    exit 1
fi

echo ""
echo "=============================================="
echo " Iniciando LEP Backend na porta 8080"
echo "=============================================="
echo ""
echo "Endpoints disponíveis:"
echo "  http://localhost:8080/ping"
echo "  http://localhost:8080/health"
echo ""
echo "Pressione Ctrl+C para parar"
echo ""

go run main.go