@echo off
setlocal enabledelayedexpansion

REM LEP Database Seeder Script for Windows
REM This script populates the database with realistic sample data

echo.
echo 🌱 LEP Database Seeder
echo ======================
echo.

REM Default values
set CLEAR_FIRST=false
set ENVIRONMENT=dev
set VERBOSE=false
set ARGS=

REM Parse command line arguments
:parse_args
if "%~1"=="" goto end_parse_args
if "%~1"=="--clear-first" (
    set CLEAR_FIRST=true
    shift
    goto parse_args
)
if "%~1"=="--environment" (
    set ENVIRONMENT=%~2
    shift
    shift
    goto parse_args
)
if "%~1"=="--verbose" (
    set VERBOSE=true
    shift
    goto parse_args
)
if "%~1"=="--help" (
    echo Usage: %0 [OPTIONS]
    echo.
    echo Options:
    echo   --clear-first    Clear existing data before seeding
    echo   --environment    Environment to seed (dev, test, staging^) [default: dev]
    echo   --verbose        Enable verbose logging
    echo   --help           Show this help message
    echo.
    echo Examples:
    echo   %0                                    # Basic seeding
    echo   %0 --clear-first --verbose            # Clear and verbose seed
    echo   %0 --environment test                 # Seed test environment
    exit /b 0
)
echo ❌ Unknown option: %~1
echo Use --help for usage information
exit /b 1

:end_parse_args

REM Check if Go is installed
go version >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo ❌ Go is not installed or not in PATH
    echo Please install Go from: https://golang.org/dl/
    exit /b 1
)

REM Check if we're in the right directory
if not exist "go.mod" (
    echo ❌ Please run this script from the LEP-Back root directory
    exit /b 1
)
if not exist "cmd\seed" (
    echo ❌ Please run this script from the LEP-Back root directory
    exit /b 1
)

REM Check if .env file exists
if not exist ".env" (
    echo ⚠️  No .env file found. Creating from example...
    if exist ".env.example" (
        copy ".env.example" ".env" >nul
        echo ℹ️  Please update .env with your database credentials
    ) else (
        echo ❌ No .env.example file found
        exit /b 1
    )
)

REM Build arguments for the seeder
if "%CLEAR_FIRST%"=="true" (
    set ARGS=!ARGS! --clear-first
)
if "%VERBOSE%"=="true" (
    set ARGS=!ARGS! --verbose
)
set ARGS=!ARGS! --environment=%ENVIRONMENT%

echo 📋 Configuration:
echo   Environment: %ENVIRONMENT%
echo   Clear first: %CLEAR_FIRST%
echo   Verbose: %VERBOSE%
echo.

REM Download dependencies if needed
echo 📦 Checking dependencies...
go mod verify >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo   Downloading dependencies...
    go mod tidy
    if %ERRORLEVEL% neq 0 (
        echo ⚠️  Warning: Failed to download dependencies, continuing...
    )
)

REM Run the seeder
echo 🚀 Running database seeder...
echo.

REM Execute the seeder
go run cmd/seed/main.go %ARGS%

if %ERRORLEVEL% equ 0 (
    echo.
    echo ✅ Database seeding completed successfully!
    echo.
    echo 🎯 Next Steps:
    echo   1. Start the backend server:
    echo      go run main.go
    echo.
    echo   2. Test the API:
    echo      curl http://localhost:8080/health
    echo.
    echo   3. Login credentials:
    echo      admin@lep-demo.com / password (Admin^)
    echo      garcom@lep-demo.com / password (Waiter^)
    echo      gerente@lep-demo.com / password (Manager^)
    echo.
    echo   4. Run tests:
    echo      go test ./tests -v
    echo.
) else (
    echo.
    echo ❌ Database seeding failed!
    echo.
    echo 💡 Troubleshooting:
    echo   1. Check if PostgreSQL is running
    echo   2. Verify database credentials in .env
    echo   3. Ensure database exists and is accessible
    echo   4. Check network connectivity
    echo.
    echo   For more details, run with --verbose flag
    exit /b 1
)

pause