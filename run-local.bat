@echo off
echo.
echo ==============================================
echo  LEP System Backend - Local Development
echo ==============================================
echo.

echo Verificando se o Go está instalado...
go version
if %ERRORLEVEL% neq 0 (
    echo ERRO: Go não está instalado ou não está no PATH
    echo Baixe o Go em: https://golang.org/dl/
    pause
    exit /b 1
)

echo.
echo Baixando dependências...
go mod tidy
if %ERRORLEVEL% neq 0 (
    echo AVISO: Falha ao baixar dependências, continuando...
)

echo.
echo Verificando build...
go build .
if %ERRORLEVEL% neq 0 (
    echo ERRO: Falha no build da aplicação
    pause
    exit /b 1
)

echo.
echo ==============================================
echo  Iniciando LEP Backend na porta 8080
echo ==============================================
echo.
echo Endpoints disponíveis:
echo   http://localhost:8080/ping
echo   http://localhost:8080/health
echo.
echo Pressione Ctrl+C para parar
echo.

go run main.go