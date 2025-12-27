# CV40 Control-Service Runbook

## How to Run the Service

- Prerequisites
  - Windows device connected to Enciris LT board (real hardware only)
  - Go toolchain installed (Go 1.22+), Enciris LT SDK available (provided in `cv40_camera_app/go`)
  - Confirm `outputId` and `canvasId` for your surgical monitor

- Build and run
  - Navigate to `cv40_camera_app/backend`
  - Create `config.json` (see Required Config Fields)
  - Run the control-service on port `8083`:
    - Development: `go run ./cmd/control-service`
    - Build binary: `go build ./cmd/control-service` then run `./control-service` (Windows: `control-service.exe`)
  - Verify health: `curl -i http://localhost:8083/health`

- Notes
  - The legacy server on `8081` can remain during migration; new clients should use `:8083` tool endpoints
  - Overlays render via CV40 canvas; ensure `overlay.output` and `overlay.canvasId` match your device

## Required Config Fields

- File: `cv40_camera_app/backend/config.json`
- Template (edit values to match your hardware):

```json
{
  "baseUrl": "cv40:/",
  "boardId": 0,
  "cameraId": 0,
  "outputId": "hdmi-out/0",
  "storageRoots": [
    "D:\\CV40\\recordings",
    "E:\\CV40\\recordings"
  ],
  "freeSpaceGb": 10,
  "recording": {
    "codec": "h264",
    "encoder": "hw",
    "bitrate": 12000,
    "container": "mp4"
  },
  "ranges": {
    "brightness": [-100, 100],
    "contrast": [-100, 100],
    "saturation": [-100, 100],
    "hue": [-180, 180],
    "sharpness": [0.0, 2.0],
    "zoom": [1.0, 4.0],
    "temperature": [2000, 10000],
    "lowLightGain": [0.0, 1.0]
  },
  "overlay": {
    "canvasId": 0,
    "output": "hdmi-out/0"
  }
}
```

## Endpoints and Curl Examples (port 8083)

- Health
  - `GET /health`
  - `curl -i http://localhost:8083/health`

- State
  - `GET /state`
  - `curl -s http://localhost:8083/state`

- Events (WebSocket)
  - `WS /events`
  - Example: `websocat ws://localhost:8083/events` (or any WS client)

- Session
  - `POST /tools/session/start`
  - `curl -s -X POST http://localhost:8083/tools/session/start -H "Content-Type: application/json" -d '{"Doctor":"Dr. Smith","Hospital":"Test Hospital","Patient":"John Doe","SurgeryType":"Arthroscopy"}'`

- Recording
  - `POST /tools/record/start`
  - `curl -s -X POST http://localhost:8083/tools/record/start`
  - `POST /tools/record/pause`
  - `curl -s -X POST http://localhost:8083/tools/record/pause`
  - `POST /tools/record/resume`
  - `curl -s -X POST http://localhost:8083/tools/record/resume`
  - `POST /tools/record/stop`
  - `curl -s -X POST http://localhost:8083/tools/record/stop`

- Photo
  - `POST /tools/photo/capture`
  - `curl -s -X POST http://localhost:8083/tools/photo/capture`

- White Balance
  - `POST /tools/whitebalance/run`
  - `curl -s -X POST http://localhost:8083/tools/whitebalance/run`

- Settings
  - `POST /tools/settings/colors`
  - `curl -s -X POST http://localhost:8083/tools/settings/colors -H "Content-Type: application/json" -d '{"brightness":10,"contrast":15,"saturation":5,"hue":0}'`
  - `POST /tools/settings/visuals`
  - `curl -s -X POST http://localhost:8083/tools/settings/visuals -H "Content-Type: application/json" -d '{"zoom":1.25,"sharpness":0.8}'`

- Presets
  - `POST /tools/preset/apply`
  - Arthroscopy: `curl -s -X POST http://localhost:8083/tools/preset/apply -H "Content-Type: application/json" -d '{"preset":"arthroscopy"}'`
  - Red Boost: `curl -s -X POST http://localhost:8083/tools/preset/apply -H "Content-Type: application/json" -d '{"preset":"red_boost"}'`

- Controller (ESP32 handheld)
  - `POST /controller/event`
  - BTN1 short → photo: `curl -s -X POST http://localhost:8083/controller/event -H "Content-Type: application/json" -d '{"deviceId":"handheld-01","btn":1,"press":"short"}'`
  - BTN1 long → stop: `curl -s -X POST http://localhost:8083/controller/event -H "Content-Type: application/json" -d '{"deviceId":"handheld-01","btn":1,"press":"long"}'`
  - BTN2 short → start/pause/resume toggle: `curl -s -X POST http://localhost:8083/controller/event -H "Content-Type: application/json" -d '{"deviceId":"handheld-01","btn":2,"press":"short"}'`
  - BTN2 long → white balance: `curl -s -X POST http://localhost:8083/controller/event -H "Content-Type: application/json" -d '{"deviceId":"handheld-01","btn":2,"press":"long"}'`
  - BTN3 short → preset toggle: `curl -s -X POST http://localhost:8083/controller/event -H "Content-Type: application/json" -d '{"deviceId":"handheld-01","btn":3,"press":"short"}'`
  - BTN4 short → zoom step: `curl -s -X POST http://localhost:8083/controller/event -H "Content-Type: application/json" -d '{"deviceId":"handheld-01","btn":4,"press":"short"}'`

## Acceptance Test Script/Checklist

- Preflight
  - Confirm `config.json` is correct for `outputId` and `canvasId`
  - Start control-service: `go run ./cmd/control-service`
  - Connect WebSocket client to `/events` to observe state changes

- Tests
  - Health
    - `curl -i http://localhost:8083/health`
    - Expect 200 OK; camera detected
  - Start Session
    - `curl -s -X POST http://localhost:8083/tools/session/start -H "Content-Type: application/json" -d '{"Doctor":"Dr. Smith","Hospital":"Test","Patient":"John","SurgeryType":"Arthroscopy"}'`
    - Verify per-target folders: `Sessions/<sessionId>/{video,photos,logs}`
    - Verify `meta.json` exists
  - Start Recording
    - `curl -s -X POST http://localhost:8083/tools/record/start`
    - Verify MP4 grows in each target `video` folder
    - Overlay shows REC on monitor
  - Pause
    - `curl -s -X POST http://localhost:8083/tools/record/pause`
    - Verify MP4 stops growing; overlay shows PAUSED
  - Resume
    - `curl -s -X POST http://localhost:8083/tools/record/resume`
    - Verify MP4 continues; overlay shows REC
  - Photo During Recording
    - `curl -s -X POST http://localhost:8083/tools/photo/capture`
    - Verify JPEG saved in each target `photos` folder; toast overlay
  - White Balance
    - `curl -s -X POST http://localhost:8083/tools/whitebalance/run`
    - Verify toast overlay; event logged
  - Stop
    - `curl -s -X POST http://localhost:8083/tools/record/stop`
    - Verify final MP4 exists and is playable; overlay cleared
    - Check `logs/events.jsonl` for `record_stop` with per-target results
  - Multi-drive Degrade
    - Unplug one drive during recording; observe `/events` and monitor overlay warning
    - Ensure other drives continue recording; state transitions to `DEGRADED`
    - Check `logs/events.jsonl` for `drive_failure`
  - ESP32 Controller
    - Trigger BTN mappings via `/controller/event` examples above; behavior matches touch UI
  - Kiosk Behavior (optional)
    - Reboot device and confirm system auto-starts without desktop steps; monitor shows video with overlays

- Troubleshooting
  - If overlays not visible: check `outputId` and `overlay.canvasId` in config, verify device routes in hardware
  - If recording not created: verify storage paths exist and are writable; inspect `events.jsonl`

---

- Primary service files
  - Entry: `cv40_camera_app/backend/cmd/control-service/main.go:1`
  - API: `cv40_camera_app/backend/internal/api/server.go:35`
  - CV40 client: `cv40_camera_app/backend/internal/cv40/real_client.go:22`
  - Overlay engine: `cv40_camera_app/backend/internal/overlay/engine.go:16`
  - Recording manager: `cv40_camera_app/backend/internal/recording/manager.go:16`
  - Storage manager: `cv40_camera_app/backend/internal/storage/manager.go:14`
  - Meta events: `cv40_camera_app/backend/internal/meta/events.go:1`
```
