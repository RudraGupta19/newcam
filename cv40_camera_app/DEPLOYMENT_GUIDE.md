# CV40 Camera System - Deployment Guide

## 📦 Package Contents

### For Main Device (Shuttle PC + Enciris LT Board)
```
cv40_camera_system/
├── backend/                 # Go backend server
│   ├── main.exe            # Compiled backend (Windows)
│   ├── handlers.go
│   ├── mock.go
│   ├── monitor_server.go
│   └── utils.go
├── monitor_static/          # Monitor display files
│   └── index.html          # Full-screen monitor display
├── go/                     # Enciris LT SDK
│   ├── lt.go
│   ├── lt_schema.go
│   └── examples/
└── start_system.bat        # Startup script
```

### For Touchscreen Device (Separate PC/Tablet)
```
cv40_touchscreen/
├── cv40_touch.exe          # Flutter touchscreen app (Windows)
├── data/                   # Flutter app data
└── start_touchscreen.bat   # Startup script
```

## 🚀 Deployment Instructions

### Main Device Setup (Shuttle PC)

1. **Copy Files**
   ```
   Copy entire cv40_camera_system/ folder to C:\CV40\
   ```

2. **Configure Network**
   - Set static IP: `192.168.1.100`
   - Subnet: `255.255.255.0`
   - Ensure both devices on same network

3. **Run System**
   ```batch
   cd C:\CV40
   start_system.bat
   ```
   
   This starts:
   - Backend API server on port 8081
   - Monitor display server on port 8082
   - Opens monitor display in full-screen browser

### Touchscreen Device Setup

1. **Copy Files**
   ```
   Copy cv40_touchscreen/ folder to C:\CV40Touch\
   ```

2. **Configure Connection**
   - Edit `start_touchscreen.bat`
   - Change IP to main device: `192.168.1.100`

3. **Run Touchscreen**
   ```batch
   cd C:\CV40Touch
   start_touchscreen.bat
   ```

## 📋 Startup Scripts

### start_system.bat (Main Device)
```batch
@echo off
echo Starting CV40 Camera System...

REM Start backend server
cd backend
start "CV40 Backend" main.exe

REM Wait for server startup
timeout /t 3

REM Open monitor display in full-screen
start "CV40 Monitor" "C:\Program Files\Google\Chrome\Application\chrome.exe" --kiosk --app=http://localhost:8082

echo System started. Backend on :8081, Monitor on :8082
pause
```

### start_touchscreen.bat (Touch Device)
```batch
@echo off
echo Starting CV40 Touchscreen...

REM Set main device IP
set MAIN_DEVICE_IP=192.168.1.100

REM Start touchscreen app
cv40_touch.exe --dart-define=API_BASE=http://%MAIN_DEVICE_IP%:8081

pause
```

## 🔧 Configuration

### Network Settings
- **Main Device**: `192.168.1.100`
- **Touch Device**: `192.168.1.101` (or DHCP)
- **API Endpoint**: `http://192.168.1.100:8081`
- **Monitor Display**: `http://192.168.1.100:8082`

### Hardware Mode vs Mock Mode
- **Hardware Mode**: Remove `MOCK_CAMERA=1` from environment
- **Mock Mode**: Set `MOCK_CAMERA=1` for testing without Enciris board

## 🖥️ Display Configuration

### Main Monitor (Surgery Display)
1. Set as primary display
2. Resolution: 1920x1080 or higher
3. Browser opens in kiosk mode automatically

### Touchscreen Panel
1. Set as secondary display
2. Touch calibration recommended
3. Flutter app auto-sizes to screen

## 🔒 Production Settings

### Windows Configuration
```batch
REM Disable Windows updates during surgery
net stop wuauserv

REM Set high performance power plan
powercfg /setactive 8c5e7fda-e8bf-4a96-9a85-a6e23a8c635c

REM Disable sleep/hibernate
powercfg /change standby-timeout-ac 0
powercfg /change hibernate-timeout-ac 0
```

### Firewall Rules
```batch
REM Allow CV40 ports
netsh advfirewall firewall add rule name="CV40 API" dir=in action=allow protocol=TCP localport=8081
netsh advfirewall firewall add rule name="CV40 Monitor" dir=in action=allow protocol=TCP localport=8082
```

## 🧪 Testing Deployment

### 1. Network Connectivity
```batch
REM From touchscreen device
ping 192.168.1.100
telnet 192.168.1.100 8081
```

### 2. API Testing
- Open `http://192.168.1.100:8081/api/settings` in browser
- Should return JSON camera settings

### 3. Monitor Display
- Open `http://192.168.1.100:8082` on main monitor
- Should show simulated camera feed with test controls

### 4. Touch Integration
- Start session on touchscreen
- Press record → Red dot should appear on monitor
- Adjust settings → Sliders should show on monitor

## 🚨 Troubleshooting

### Backend Won't Start
```batch
REM Check if ports are in use
netstat -ano | findstr :8081
netstat -ano | findstr :8082

REM Kill conflicting processes
taskkill /PID [PID_NUMBER] /F
```

### Touchscreen Can't Connect
1. Check network connectivity
2. Verify main device IP address
3. Check Windows Firewall settings
4. Ensure backend is running

### Monitor Display Issues
1. Try different browser (Chrome recommended)
2. Check WebSocket connection in browser console
3. Verify port 8082 is accessible

## 📁 File Locations

### Main Device
- **Backend**: `C:\CV40\backend\main.exe`
- **Monitor**: `C:\CV40\monitor_static\index.html`
- **Logs**: `C:\CV40\logs\`
- **Recordings**: `C:\CV40\recordings\` (or USB drives)

### Touch Device
- **App**: `C:\CV40Touch\cv40_touch.exe`
- **Config**: `C:\CV40Touch\config.json`

## 🔄 Updates

### Backend Updates
1. Stop system: `taskkill /IM main.exe /F`
2. Replace `main.exe`
3. Restart: `start_system.bat`

### Touchscreen Updates
1. Stop app: Close Flutter window
2. Replace `cv40_touch.exe`
3. Restart: `start_touchscreen.bat`

---

**Ready for Production Deployment** ✅
