@echo off
echo Creating CV40 Transfer Package...
echo.

REM Clean up any existing transfer folder
if exist transfer rmdir /s /q transfer

REM Create transfer structure
mkdir transfer
mkdir transfer\backend
mkdir transfer\monitor_static
mkdir transfer\go

echo [1/6] Copying backend files...
copy backend\main.exe transfer\backend\ >nul
copy backend\*.go transfer\backend\ >nul
copy backend\go.mod transfer\backend\ >nul

echo [2/6] Copying monitor display...
copy monitor_static\index.html transfer\monitor_static\ >nul

echo [3/6] Copying touchscreen interface...
copy touchscreen_config.html transfer\ >nul

echo [4/6] Copying Enciris SDK...
xcopy /E /I /Q ..\go transfer\go >nul

echo [5/6] Copying startup script...
copy start_system.bat transfer\startup_script.bat >nul

echo [6/6] Creating documentation...
copy transfer\README_FIRST.txt transfer\ >nul 2>&1
copy transfer\INSTALLATION_GUIDE.md transfer\ >nul 2>&1
copy transfer\QUICK_START.txt transfer\ >nul 2>&1

echo.
echo Creating ZIP archive for easier transfer...
powershell -command "Compress-Archive -Path 'transfer\*' -DestinationPath 'CV40_System.zip' -Force"

echo.
echo ========================================
echo Transfer Package Created Successfully!
echo ========================================
echo.
echo Option 1: Copy "transfer" folder to Shuttle PC
echo Option 2: Copy "CV40_System.zip" and extract on Shuttle PC
echo.
echo Files ready for transfer:
dir transfer /b
echo.
echo ZIP file size:
dir CV40_System.zip
echo.
pause
