# 🏗️ CV40 Camera System - Build & Deploy Instructions

## 📱 Building for Your Actual Device

### Prerequisites
- Go 1.22+ installed
- Flutter SDK installed (for touchscreen app)
- Windows 10/11 target devices
- Network connectivity between devices

## 🖥️ Main Device Build (Shuttle PC + Enciris LT)

### Step 1: Compile Backend
```powershell
# Navigate to backend directory
cd cv40_camera_app/backend

# Build Windows executable
go build -o main.exe .

# Verify build
dir main.exe
```

### Step 2: Prepare Main Device Package
```powershell
# Create deployment folder
mkdir C:\CV40_Deploy\cv40_camera_system

# Copy backend
copy backend\main.exe C:\CV40_Deploy\cv40_camera_system\
copy backend\*.go C:\CV40_Deploy\cv40_camera_system\

# Copy monitor display
mkdir C:\CV40_Deploy\cv40_camera_system\monitor_static
copy monitor_static\index.html C:\CV40_Deploy\cv40_camera_system\monitor_static\

# Copy Enciris SDK (CRITICAL!)
xcopy /E /I go C:\CV40_Deploy\cv40_camera_system\go

# Copy startup script
copy start_system.bat C:\CV40_Deploy\cv40_camera_system\
```

## 📱 Touchscreen Device Build

### Option 1: Flutter Native App (Recommended)
```powershell
# Build Flutter Windows app for touchscreen
cd cv40_camera_app
flutter build windows --release -t lib/touchscreen_main.dart

# Package will be in: build\windows\x64\runner\Release\
```

### Option 2: Web-based Touchscreen (Simpler)
```powershell
# Just copy the HTML file
mkdir C:\CV40_Deploy\cv40_touchscreen
copy touchscreen_config.html C:\CV40_Deploy\cv40_touchscreen\
copy start_touchscreen.bat C:\CV40_Deploy\cv40_touchscreen\
```

## 🚀 Deployment to Actual Devices

### Main Device (Shuttle PC) Setup

1. **Copy Files**
   ```
   Copy C:\CV40_Deploy\cv40_camera_system\ to target device: C:\CV40\
   ```

2. **Network Configuration**
   ```powershell
   # Set static IP
   netsh interface ip set address "Ethernet" static 192.168.1.100 255.255.255.0
   
   # Verify
   ipconfig
   ```

3. **Firewall Configuration**
   ```powershell
   # Allow CV40 ports
   netsh advfirewall firewall add rule name="CV40 API" dir=in action=allow protocol=TCP localport=8081
   netsh advfirewall firewall add rule name="CV40 Monitor" dir=in action=allow protocol=TCP localport=8082
   ```

4. **Hardware vs Mock Mode**
   ```powershell
   # For hardware mode (production)
   cd C:\CV40
   start_system.bat
   
   # For mock mode (testing)
   set MOCK_CAMERA=1
   start_system.bat
   ```

### Touchscreen Device Setup

1. **Copy Files**
   ```
   Copy touchscreen files to target device: C:\CV40Touch\
   ```

2. **Network Configuration**
   ```
   Connect to same network as main device (can use DHCP)
   ```

3. **Start Touchscreen**
   ```powershell
   cd C:\CV40Touch
   start_touchscreen.bat
   ```

## 🔧 Display Configuration

### Main Device (Monitor Display)
- **Resolution**: Set to 1920x1080 or higher
- **Display**: Set as primary display
- **Browser**: Chrome/Edge will open in kiosk mode automatically

### Touchscreen Device
- **Resolution**: 2048x1536 (as per your specs)
- **Touch Calibration**: Calibrate touch input in Windows settings
- **App Size**: Flutter app will render at exact 2048x1536 dimensions

## 🧪 Testing the Build

### 1. Test Main Device
```powershell
# Start backend
cd C:\CV40
start_system.bat

# Verify services
netstat -ano | findstr :8081
netstat -ano | findstr :8082

# Test API
curl http://localhost:8081/api/settings
```

### 2. Test Touchscreen
```powershell
# Start touchscreen
cd C:\CV40Touch
start_touchscreen.bat

# Verify connection to main device
ping 192.168.1.100
```

### 3. Test Integration
- Press buttons on touchscreen
- Verify actions appear on main monitor
- Check debug messages in touchscreen app

## 📐 App Dimensions (As Per Your Specs)

### Touchscreen App: 2048x1536
- **Top Bar**: 2048x267 (pencil icon 150x150, text 878x191, settings 120x120)
- **Main Area**: 2048x1249 with 235px top offset
- **Padding**: top:171px, right:95px, bottom:171px, left:95px, gap:79px
- **Control Buttons**: 600x600 each (camera, record, white balance)

### Monitor Display: Full Screen
- **Camera Feed**: Simulated tissue background
- **OSD Elements**: Recording dot, parameter sliders, photo popups
- **Real-time Sync**: WebSocket connection to touchscreen

## 🔄 Production Deployment

### Auto-Start Configuration
```powershell
# Main device - Windows service or startup script
schtasks /create /tn "CV40System" /tr "C:\CV40\start_system.bat" /sc onlogon

# Touchscreen - Auto-start browser
schtasks /create /tn "CV40Touch" /tr "C:\CV40Touch\start_touchscreen.bat" /sc onlogon
```

### Kiosk Mode (Optional)
```powershell
# Enable Shell Launcher v2 for true kiosk mode
# This replaces Windows Explorer with your app
```

## 🚨 Troubleshooting

### Backend Won't Start
```powershell
# Check port conflicts
netstat -ano | findstr :8081
taskkill /PID [PID] /F

# Check Enciris SDK
dir C:\CV40\go\lt.go
```

### Touchscreen Connection Issues
```powershell
# Test network connectivity
ping 192.168.1.100
telnet 192.168.1.100 8081

# Check firewall
netsh advfirewall show allprofiles
```

### Flutter Build Issues
```powershell
# Clean and rebuild
flutter clean
flutter pub get
flutter build windows --release -t lib/touchscreen_main.dart
```

## 📋 Final Checklist

- [ ] Backend compiled and tested
- [ ] Monitor display opens full-screen
- [ ] Touchscreen app renders at 2048x1536
- [ ] Network connectivity between devices
- [ ] API calls working (check debug messages)
- [ ] Recording/capture functions tested
- [ ] White balance and presets working
- [ ] Hardware mode tested (if Enciris board available)

---

**Your CV40 system is ready for production deployment!** 🎉
