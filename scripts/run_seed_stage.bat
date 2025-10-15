@echo off
REM Script para executar seed no ambiente stage

echo ========================================
echo LEP Database Seeder - Stage Environment
echo ========================================
echo.

REM Compilar o seeder
echo [1/2] Compilando seeder...
"C:\Go\bin\go.exe" build -o lep-seed.exe cmd/seed/*.go
if %errorlevel% neq 0 (
    echo Erro ao compilar o seeder!
    exit /b %errorlevel%
)

echo [2/2] Executando seed no ambiente stage...
echo.

REM Carregar variáveis do .env.stage
for /f "usebackq tokens=1,* delims==" %%a in (".env.stage") do (
    set "%%a=%%b"
)

REM Executar o seed
lep-seed.exe --clear-first --verbose

echo.
echo ========================================
echo Seed concluído!
echo ========================================
