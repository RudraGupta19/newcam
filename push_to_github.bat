@echo off
echo Pushing CV40 Camera System to GitHub...
echo.

echo Step 1: Setting up remote repository...
git remote remove origin 2>nul
git remote add origin https://github.com/RudraGupta19/newcam.git

echo Step 2: Pushing to GitHub...
git push -u origin main

echo.
if %errorlevel% == 0 (
    echo ========================================
    echo SUCCESS! Repository pushed to GitHub
    echo ========================================
    echo.
    echo Repository URL: https://github.com/RudraGupta19/newcam
    echo.
    echo You can now clone this repository on any device:
    echo git clone https://github.com/RudraGupta19/newcam.git
) else (
    echo ========================================
    echo FAILED! Please check:
    echo ========================================
    echo 1. Repository exists on GitHub
    echo 2. You're signed in as RudraGupta19
    echo 3. Repository name is exactly: camera_enciris
    echo 4. You have push permissions
)

echo.
pause
