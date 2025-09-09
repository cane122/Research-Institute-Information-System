@echo off
echo ==================================================
echo Research Institute Information System - Setup
echo ==================================================
echo.

echo Checking prerequisites...
echo.

:: Check if Go is installed
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Go is not installed or not in PATH
    echo Please install Go from: https://golang.org/dl/
    pause
    exit /b 1
) else (
    echo [OK] Go is installed
)

:: Check if Wails is installed
wails version >nul 2>&1
if %errorlevel% neq 0 (
    echo [INFO] Installing Wails CLI...
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    if %errorlevel% neq 0 (
        echo [ERROR] Failed to install Wails CLI
        pause
        exit /b 1
    )
    echo [OK] Wails CLI installed
) else (
    echo [OK] Wails CLI is installed
)

echo.
echo ==================================================
echo Setup completed successfully!
echo ==================================================
echo.
echo Next steps:
echo 1. Install PostgreSQL if not already installed
echo 2. Create database: CREATE DATABASE research_institute;
echo 3. Run database schema: psql -U postgres -d research_institute -f database/schema.sql
echo 4. Build the application: wails build
echo 5. Run: .\build\bin\research-institute-system.exe
echo.
echo For detailed instructions, see README.md
echo.
pause
