# 🖥️ CV40 Camera System - Single PC Dual Display Setup

## ✅ Correct Configuration

You have:
- **1x Shuttle PC** (main computer with Enciris LT board)
- **2x HDMI outputs** from the same PC
- **Display 1**: Main monitor (surgery view - full screen)
- **Display 2**: Touchscreen panel (controls - 2048x1536)

This is **much simpler** than separate devices!

## 🚀 Single PC Deployment

### Step 1: Copy Files to ONE Location
```
Copy entire cv40_camera_app folder to your Shuttle PC: C:\CV40\
```

### Step 2: Windows Display Setup
1. **Connect both displays** via HDMI to your Shuttle PC
2. **Windows Settings** → Display
3. **Set up as "Extend these displays"**
4. **Identify displays**:
   - Display 1: Main monitor (surgery view)
   - Display 2: Touchscreen panel (controls)
5. **Set touchscreen as touch input** for Display 2

### Step 3: Start the System
```powershell
# Navigate to CV40 folder
cd C:\CV40

# Run the startup script
start_system.bat
```

## 🎯 What Happens When You Start

1. **Backend starts** on port 8081
2. **Monitor display opens** (http://localhost:8082) - drag to Display 1
3. **Touchscreen control opens** (touchscreen_config.html) - drag to Display 2
4. **Both apps connect** to localhost (same PC)

## 📋 Updated Files Structure

```
C:\CV40\
├── backend/
│   ├── main.exe ⭐ (Go backend server)
│   └── *.go files
├── monitor_static/
│   └── index.html ⭐ (Monitor display)
├── go/ ⭐ (Enciris SDK)
├── touchscreen_config.html ⭐ (Touch controls)
└── start_system.bat ⭐ (Starts everything)
```

## 🔧 Display Positioning

### Automatic Window Positioning (Optional)
If you want windows to open on specific displays automatically, you can use:

```powershell
# For advanced users - position windows on specific displays
# This requires additional tools like DisplayFusion or custom scripts
```

### Manual Positioning (Recommended)
1. **Start the system** with `start_system.bat`
2. **Two browser windows open**
3. **Drag monitor display** to your main monitor (Display 1)
4. **Drag touchscreen control** to your touchscreen (Display 2)
5. **Make monitor display full-screen** (F11)

## 🧪 Testing Your Setup

### 1. Start System
```powershell
cd C:\CV40
start_system.bat
```

### 2. Verify Both Windows Open
- Monitor display: Full-screen surgical view
- Touchscreen: 2048x1536 control panel with debug info

### 3. Test Integration
- Press buttons on touchscreen
- See immediate response on monitor display
- Check debug messages show "Connected" status

## 🎯 Hardware vs Mock Mode

### Mock Mode (Testing)
```powershell
# Set environment variable for testing
set MOCK_CAMERA=1
start_system.bat
```

### Hardware Mode (Production)
```powershell
# Remove mock mode to use real Enciris LT board
start_system.bat
```

## 🔧 Troubleshooting

### Windows Won't Open on Correct Display
1. Start the system
2. Manually drag windows to correct displays
3. Windows will remember positions for next time

### Touchscreen Not Responding
1. **Windows Settings** → Devices → Pen & Windows Ink
2. **Calibrate touch** for your touchscreen display
3. **Set touch target** to correct display

### Backend Connection Issues
- Both apps connect to `localhost:8081` (same PC)
- No network configuration needed
- Check Windows Firewall if issues persist

## 📐 Display Specifications

### Main Monitor (Display 1)
- **Content**: Full-screen surgical view with OSD
- **Resolution**: Any (1920x1080+ recommended)
- **Mode**: Full-screen browser (kiosk mode)

### Touchscreen (Display 2)
- **Content**: Control panel interface
- **Size**: 2048x1536 (exact match to your specs)
- **Mode**: App window sized to touchscreen
- **Touch**: Calibrated for precise button presses

## ✅ Advantages of Single PC Setup

- **Simpler**: No network configuration needed
- **Faster**: No network latency between displays
- **Reliable**: No network connection issues
- **Cheaper**: One computer instead of two
- **Easier**: Single point of control and updates

---

## 🎉 Ready to Deploy!

Your setup is actually **much simpler** than I initially thought. Just:

1. **Copy files** to `C:\CV40\` on your Shuttle PC
2. **Connect both displays** via HDMI
3. **Run** `start_system.bat`
4. **Drag windows** to correct displays

**Single PC, dual display, perfect integration!** 🚀
