@echo off
echo Starting CV40 Camera System - Single PC Dual Display...

REM Start backend server
cd backend
start "CV40 Backend" main.exe

REM Wait for server startup
timeout /t 3

REM Open monitor display (main monitor - Display 1)
if exist "C:\Program Files\Google\Chrome\Application\chrome.exe" (
    start "CV40 Monitor" "C:\Program Files\Google\Chrome\Application\chrome.exe" --new-window --kiosk http://localhost:8082
) else if exist "C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe" (
    start "CV40 Monitor" "C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe" --new-window --kiosk http://localhost:8082
) else (
    start "CV40 Monitor" http://localhost:8082
)

REM Wait 2 seconds
timeout /t 2

REM Open touchscreen control (touchscreen - Display 2)
cd ..
if exist "C:\Program Files\Google\Chrome\Application\chrome.exe" (
    start "CV40 Touch" "C:\Program Files\Google\Chrome\Application\chrome.exe" --new-window --app=file:///%CD%/touchscreen_config.html
) else if exist "C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe" (
    start "CV40 Touch" "C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe" --new-window --app=file:///%CD%/touchscreen_config.html
) else (
    start "CV40 Touch" touchscreen_config.html
)

echo.
echo ========================================
echo CV40 Camera System Started
echo ========================================
echo Backend API: http://localhost:8081
echo Monitor Display: http://localhost:8082 (Display 1)
echo Touch Control: touchscreen_config.html (Display 2)
echo.
echo Drag windows to correct displays if needed
echo Press any key to stop the system...
pause

REM Stop backend when user presses key
taskkill /IM main.exe /F 2>nul
echo System stopped.
