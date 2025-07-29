@echo off

REM Build script for Claude Config Manager (Windows)

REM Check if wails is installed
where wails >nul 2>nul
if %errorlevel% neq 0 (
    echo Wails is not installed. Installing...
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
)

REM Install frontend dependencies
echo Installing frontend dependencies...
cd frontend
npm install
cd ..

REM Build for current platform
echo Building for current platform...
wails build

echo Build completed successfully!
pause