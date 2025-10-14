@echo off
echo ========================================
echo Script de Rebuild do LEP System
echo ========================================
echo.

echo Passo 1: Fechando processos Go em execucao...
taskkill /F /IM go.exe 2>nul
taskkill /F /IM lep-system.exe 2>nul
timeout /t 2 >nul

echo Passo 2: Limpando cache de build...
go clean -cache 2>nul

echo Passo 3: Removendo binario antigo...
if exist lep-system.exe del /F lep-system.exe

echo Passo 4: Atualizando modulos...
go mod tidy

echo Passo 5: Compilando projeto...
go build -o lep-system.exe .

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ========================================
    echo BUILD CONCLUIDO COM SUCESSO!
    echo ========================================
    echo Binario criado: lep-system.exe
) else (
    echo.
    echo ========================================
    echo ERRO NO BUILD
    echo ========================================
    echo Se o erro persistir, reinicie o computador
    echo para liberar o cache do Go.
)

echo.
pause
