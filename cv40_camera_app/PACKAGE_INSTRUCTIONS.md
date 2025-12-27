# 📦 CV40 Camera System - Packaging Instructions

## 🎯 Two-Device Setup

### Device 1: Main System (Shuttle PC + Enciris LT)
**Copy these files to your main device:**

```
📁 cv40_camera_system/
├── 📁 backend/
│   ├── main.exe ⭐ (Compiled Go backend)
│   ├── handlers.go
│   ├── mock.go
│   ├── monitor_server.go
│   ├── utils.go
│   └── go.mod
├── 📁 monitor_static/
│   └── index.html ⭐ (Monitor display)
├── 📁 go/ ⭐ (Enciris SDK - REQUIRED)
│   ├── lt.go
│   ├── lt_schema.go
│   ├── lt_windows.go
│   └── examples/
└── start_system.bat ⭐ (Startup script)
```

### Device 2: Touchscreen Control
**Copy these files to your touchscreen device:**

```
📁 cv40_touchscreen/
├── touchscreen_config.html ⭐ (Touch control panel)
└── start_touchscreen.bat ⭐ (Startup script)
```

## 🚀 Quick Start Instructions

### Step 1: Setup Main Device (Shuttle PC)

1. **Copy Files**
   ```
   Copy entire cv40_camera_system folder to: C:\CV40\
   ```

2. **Set Network IP** (Important!)
   - Set static IP: `192.168.1.100`
   - Subnet mask: `255.255.255.0`

3. **Start System**
   ```
   Double-click: C:\CV40\start_system.bat
   ```
   
   ✅ This will:
   - Start backend server (port 8081)
   - Start monitor server (port 8082)  
   - Open full-screen monitor display

### Step 2: Setup Touchscreen Device

1. **Copy Files**
   ```
   Copy cv40_touchscreen folder to: C:\CV40Touch\
   ```

2. **Connect to Network**
   - Same network as main device
   - Can use DHCP (automatic IP)

3. **Start Touchscreen**
   ```
   Double-click: C:\CV40Touch\start_touchscreen.bat
   ```
   
   ✅ This opens the touch control panel in browser

## 🔧 Configuration

### Network Setup
- **Main Device**: `192.168.1.100` (static)
- **Touch Device**: Any IP on same network
- **Connection**: Touch panel connects to main device

### First Time Setup
1. On touchscreen, click the IP config (top-left)
2. Enter main device IP: `192.168.1.100`
3. Click "Update" - should show "Connected"

## 📋 Startup Scripts

### start_system.bat (Main Device)
```batch
@echo off
echo Starting CV40 Camera System...

cd backend
start "CV40 Backend" main.exe

timeout /t 3

REM Open monitor in full-screen
start "CV40 Monitor" chrome.exe --kiosk --app=http://localhost:8082

echo System started on ports 8081 and 8082
pause
```

### start_touchscreen.bat (Touch Device)
```batch
@echo off
echo Starting CV40 Touchscreen Control...

REM Open touchscreen control panel
start "CV40 Touch" chrome.exe --app=touchscreen_config.html

echo Touchscreen started
pause
```

## 🧪 Testing the Setup

### 1. Test Main Device
- Backend should start on port 8081
- Monitor display should open full-screen
- Test buttons in monitor display work

### 2. Test Touchscreen
- Control panel should open
- Top-left should show "Connected" status
- All buttons should be responsive

### 3. Test Integration
- Press **Record** on touch → Red dot appears on monitor
- Press **Photo** on touch → Photo popup on monitor
- Adjust settings → Parameter sliders on monitor

## 🔒 Hardware vs Mock Mode

### Mock Mode (Testing)
- Set environment variable: `MOCK_CAMERA=1`
- Simulates all camera functions
- Perfect for testing without Enciris hardware

### Hardware Mode (Production)
- Remove `MOCK_CAMERA=1` environment variable
- Connects to real Enciris LT board via named pipes
- Requires Enciris SDK and hardware

## 🚨 Troubleshooting

### "Backend won't start"
```batch
REM Check if port is in use
netstat -ano | findstr :8081

REM Kill conflicting process
taskkill /PID [PID] /F
```

### "Touchscreen can't connect"
1. Check main device IP address
2. Ping test: `ping 192.168.1.100`
3. Check Windows Firewall
4. Verify backend is running

### "Monitor display blank"
1. Try different browser (Chrome recommended)
2. Check URL: `http://localhost:8082`
3. Check browser console for errors

## 📁 File Locations After Setup

### Main Device (C:\CV40\)
- **Backend**: `backend\main.exe`
- **Monitor**: `monitor_static\index.html`
- **SDK**: `go\` folder (required for hardware)
- **Recordings**: Auto-created in `recordings\`

### Touch Device (C:\CV40Touch\)
- **Control Panel**: `touchscreen_config.html`
- **Config**: Saved in browser localStorage

## 🔄 Updates

### Update Backend
1. Stop system (close command window)
2. Replace `main.exe` with new version
3. Restart `start_system.bat`

### Update Touchscreen
1. Replace `touchscreen_config.html`
2. Refresh browser or restart

---

## ⚡ Ready to Deploy!

1. **Package the files** as described above
2. **Copy to your devices** 
3. **Set network IPs**
4. **Run the startup scripts**
5. **Test the integration**

Your dual-device CV40 camera system is ready! 🎉
