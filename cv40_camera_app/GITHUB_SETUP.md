# 🚀 CV40 Camera System - GitHub Repository Setup

## 📋 Repository Information
- **Username**: RudraGupta19
- **Repository**: camera_enciris
- **URL**: https://github.com/RudraGupta19/camera_enciris

## 🔧 Setup Instructions

### Step 1: Create Repository on GitHub
1. Go to https://github.com
2. Sign in as **RudraGupta19**
3. Click "New repository" (green button)
4. Repository name: **camera_enciris**
5. Description: **CV40 Medical Camera Control System - Endoscopic Surgery Interface**
6. Set to **Private** (recommended for medical software)
7. ✅ Add README file
8. ✅ Add .gitignore → Choose "Go"
9. Click "Create repository"

### Step 2: Clone Repository Locally
```bash
# Navigate to your projects directory
cd C:\Users\adlife\IdeaProjects\

# Clone the repository
git clone https://github.com/RudraGupta19/camera_enciris.git

# Navigate into the repository
cd camera_enciris
```

### Step 3: Copy Your CV40 Files
```bash
# Copy all your CV40 files to the repository
xcopy /E /I "C:\Users\adlife\IdeaProjects\camera_enciris\cv40_camera_app\*" .

# Or if you're in the cv40_camera_app directory:
xcopy /E /I . "C:\Users\adlife\IdeaProjects\camera_enciris"
```

### Step 4: Initial Commit and Push
```bash
# Add all files
git add .

# Commit with descriptive message
git commit -m "Initial commit: CV40 Medical Camera Control System

- Complete Go backend with Enciris LT SDK integration
- Flutter touchscreen interface (2048x1536)
- Dual-display monitor system
- Mock mode for development
- Hardware mode for production deployment
- Session management and multi-drive recording
- White balance and preset system (Arthroscopy/Red Boost)"

# Push to GitHub
git push origin main
```

## 🔄 Alternative: Push Existing Project

If you want to push your current cv40_camera_app folder:

```bash
# Navigate to your current project
cd C:\Users\adlife\IdeaProjects\camera_enciris\cv40_camera_app

# Initialize git repository
git init

# Add remote repository
git remote add origin https://github.com/RudraGupta19/camera_enciris.git

# Add all files
git add .

# Initial commit
git commit -m "CV40 Medical Camera System - Phase 1 Complete"

# Push to GitHub
git push -u origin main
```

## 📁 Repository Structure

Your GitHub repository will contain:
```
camera_enciris/
├── README.md ⭐ (Project overview)
├── .gitignore ⭐ (Excludes binaries and build artifacts)
├── backend/ ⭐ (Go REST API server)
│   ├── main.go
│   ├── handlers.go
│   ├── mock.go
│   └── go.mod
├── lib/ ⭐ (Flutter touchscreen app)
│   ├── main.dart
│   ├── touchscreen_app.dart
│   ├── settings_page.dart
│   └── session_start_page.dart
├── go/ ⭐ (Enciris LT SDK)
│   ├── lt.go
│   ├── lt_schema.go
│   └── examples/
├── monitor_static/ ⭐ (Monitor display)
│   └── index.html
├── touchscreen_config.html ⭐ (Web-based touch controls)
├── BUILD_INSTRUCTIONS.md
├── DEPLOYMENT_GUIDE.md
└── pubspec.yaml ⭐ (Flutter dependencies)
```

## 🔄 Working on Another Device

### Step 1: Clone Repository
```bash
# On your other device
git clone https://github.com/RudraGupta19/camera_enciris.git
cd camera_enciris
```

### Step 2: Install Dependencies
```bash
# Install Go dependencies
cd backend
go mod tidy

# Install Flutter dependencies (if developing UI)
cd ..
flutter pub get
```

### Step 3: Build and Run
```bash
# Build backend
cd backend
go build -o main.exe .

# Run in mock mode
set MOCK_CAMERA=1
main.exe

# Or run Flutter app
cd ..
flutter run -d windows -t lib/touchscreen_main.dart
```

## 🔐 Repository Settings (Recommended)

### Security Settings:
1. **Repository visibility**: Private (medical software)
2. **Branch protection**: Protect main branch
3. **Required reviews**: Enable for production changes
4. **Secrets**: Add any API keys or sensitive config

### Collaboration:
1. **Add collaborators**: Invite team members
2. **Issues**: Enable for bug tracking
3. **Projects**: Use for feature planning
4. **Wiki**: Document medical compliance requirements

## 📋 Git Workflow for Development

### Daily Development:
```bash
# Pull latest changes
git pull origin main

# Create feature branch
git checkout -b feature/new-preset-system

# Make changes, then commit
git add .
git commit -m "Add new preset system for laparoscopy"

# Push feature branch
git push origin feature/new-preset-system

# Create pull request on GitHub
```

### Production Releases:
```bash
# Tag releases for deployment
git tag -a v1.0.0 -m "Phase 1 Release - Basic camera controls"
git push origin v1.0.0
```

## 🎯 Next Steps

1. **Create the repository** on GitHub (Step 1 above)
2. **Push your code** using the commands above
3. **Clone on other devices** for development
4. **Use branches** for new features
5. **Tag releases** for production deployments

Your CV40 camera system will be version-controlled and accessible from any device! 🎉
