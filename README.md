# CV40 Camera System - Phase 1

A complete medical camera control system built with Go backend and Flutter UI, designed for endoscopic surgery applications.

## 🚀 Quick Start (Local Development)

### Prerequisites
- Go 1.22+ 
- Flutter SDK (for UI development)
- Windows 10/11 (target platform)

### Running the System

1. **Start the Backend (Mock Mode)**
   ```powershell
   cd cv40_camera_app/backend
   $env:MOCK_CAMERA="1"
   go run .
   ```
   Server starts on http://localhost:8081

2. **Test the API**
   Open `cv40_camera_app/test_api.html` in your browser to test all endpoints interactively.

3. **Run Flutter UI** (if Flutter is installed)
   ```powershell
   cd cv40_camera_app
   flutter run -d windows
   ```

## 📋 Features Implemented

### ✅ Phase 1 Complete
- **Session Management**: Surgery metadata capture and folder creation
- **Recording Controls**: Start/pause/resume/stop with multi-drive mirroring
- **Still Capture**: JPEG capture during or outside recording
- **Settings System**: Brightness, zoom, contrast, sharpness, color controls
- **Preset System**: Arthroscopy (default) and Red Boost presets
- **White Balance**: Automated color temperature adjustment
- **Mock Mode**: Full local development without hardware

### 🎯 Core Endpoints
- `POST /api/session/start` - Initialize surgery session
- `POST /api/capture` - Capture still image
- `POST /api/recording/{start,pause,resume,stop}` - Recording control
- `POST /api/white-balance` - Execute white balance
- `GET/POST /api/settings` - Camera parameter control
- `POST /api/presets/apply` - Apply Arthroscopy/Red Boost presets

### 🎨 UI Components
- **Main Interface**: Touch-optimized camera controls
- **Settings Page**: 3-tab interface (Primary/Colour/Advanced)
- **Session Dialog**: Surgery metadata entry
- **Red Boost Toggle**: Real-time preset switching

## 🏗️ Architecture

```
cv40_camera_app/
├── backend/           # Go REST API server
│   ├── main.go       # HTTP router and server
│   ├── handlers.go   # API endpoint implementations
│   ├── mock.go       # Hardware simulation for development
│   └── utils.go      # CORS and JSON helpers
├── lib/              # Flutter UI application
│   ├── main.dart     # Main camera interface
│   ├── settings_page.dart    # Settings controls
│   └── session_start_page.dart # Surgery metadata entry
└── go/               # Enciris LT SDK (provided)
    ├── lt.go         # Core SDK client
    ├── lt_schema.go  # Camera/device data structures
    └── examples/     # SDK usage examples
```

## 🔧 Hardware Integration

### Enciris LT SDK Integration
- **Transport**: Named pipes (Windows) / Unix sockets (Linux)
- **Camera Control**: Full parameter access (exposure, colors, visuals)
- **Recording**: Hardware-accelerated H.264/HEVC encoding
- **Still Capture**: JPEG with metadata
- **Button Input**: Head-mounted button support

### Multi-Drive Recording
- Automatic USB drive detection
- Simultaneous recording to all connected drives
- Per-surgery folder creation with metadata

## 📱 Development Workflow

### Local Testing (No Hardware)
1. Set `MOCK_CAMERA=1` environment variable
2. Backend simulates all camera responses
3. UI fully functional for development/demo

### Hardware Deployment
1. Remove `MOCK_CAMERA` environment variable
2. Ensure Enciris LT board is connected
3. Backend automatically connects to cv40:/ endpoints

## 🎯 Next Steps (Phase 2 & 3)

### Phase 2 - Enhanced Presets
- [ ] Multiple surgery type presets
- [ ] Custom preset creation and saving
- [ ] Preset import/export functionality
- [ ] Real-time parameter overlay on monitor

### Phase 3 - Mobile & Cloud
- [ ] iOS/Android companion app
- [ ] AWS S3 cloud storage integration
- [ ] Remote camera control via mobile
- [ ] Subscription-based storage plans

## 🔒 Security Features
- Code signing for production binaries
- Settings validation and sanitization
- Secure preset import/export with signatures
- Kiosk mode for production deployment

## 📊 File Organization
- **Session Folders**: `YYYYMMDD_HHMMSS_Surgery_Patient/`
- **Metadata**: `metadata.json` with surgery details
- **Media Files**: `video_001.mp4`, `photo_001.jpg` etc.
- **Multi-Drive**: Identical copies on all connected storage

## 🛠️ Production Deployment
- Windows 11 kiosk mode with Shell Launcher v2
- BitLocker encryption and Secure Boot
- Automatic startup and crash recovery
- Medical-grade hardware compliance ready

---

**Status**: Phase 1 Complete ✅  
**Next**: Hardware integration testing and Phase 2 development
