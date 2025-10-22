@echo off
REM LEP Remote Database Seeder Script (Windows)
REM This script populates a remote database via HTTP API calls

setlocal enabledelayedexpansion

REM Default values
set "BASE_URL=https://lep-system-516622888070.us-central1.run.app"
set "ENVIRONMENT=stage"
set "VERBOSE="

REM Parse command line arguments
:parse_args
if "%1"=="" goto end_parse
if "%1"=="--url" (
    set "BASE_URL=%2"
    shift
    shift
    goto parse_args
)
if "%1"=="--environment" (
    set "ENVIRONMENT=%2"
    shift
    shift
    goto parse_args
)
if "%1"=="--verbose" (
    set "VERBOSE=--verbose"
    shift
    goto parse_args
)
if "%1"=="-v" (
    set "VERBOSE=--verbose"
    shift
    goto parse_args
)
if "%1"=="--help" goto show_help
if "%1"=="-h" goto show_help

echo Unknown option: %1
echo Use --help to see available options
exit /b 1

:show_help
echo Usage: scripts\run_seed_remote.bat [OPTIONS]
echo.
echo Options:
echo   --url URL              Base URL of the API (default: https://lep-system-516622888070.us-central1.run.app)
echo   --environment ENV      Environment to seed (stage, prod) (default: stage)
echo   --verbose, -v          Enable verbose logging
echo   --help, -h             Show this help message
echo.
echo Examples:
echo   scripts\run_seed_remote.bat
echo   scripts\run_seed_remote.bat --verbose
echo   scripts\run_seed_remote.bat --url https://api.example.com --environment prod
exit /b 0

:end_parse

REM Print header
echo ========================================
echo    LEP Remote Database Seeder
echo ========================================
echo.

REM Print configuration
echo Configuration:
echo   Target URL: %BASE_URL%
echo   Environment: %ENVIRONMENT%
echo   Verbose: %VERBOSE%
echo.

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo Error: Go is not installed or not in PATH
    exit /b 1
)

REM Navigate to project root
cd /d "%~dp0\.."

REM Build the seed-remote binary
echo Building seed-remote binary...
"C:\Go\bin\go.exe" build -o lep-seed-remote.exe cmd/seed-remote/main.go

if %ERRORLEVEL% neq 0 (
    echo Failed to build seed-remote binary
    exit /b 1
)

echo [OK] Binary built successfully
echo.

REM Run the seeder
echo Starting remote seeding process...
echo.

lep-seed-remote.exe --url "%BASE_URL%" --environment "%ENVIRONMENT%" %VERBOSE%

if %ERRORLEVEL% equ 0 (
    echo.
    echo ========================================
    echo    Seeding completed successfully!
    echo ========================================
    echo.
    echo Next steps:
    echo   1. Access the application at: %BASE_URL%
    echo   2. Login with credentials shown above
    echo.
) else (
    echo.
    echo ========================================
    echo    Seeding failed!
    echo ========================================
    exit /b 1
)

endlocal
