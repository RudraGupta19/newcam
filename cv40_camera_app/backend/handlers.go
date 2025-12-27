package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	lt "lt/client/go"
)

type Server struct {
    mu           sync.Mutex
    recording    bool
    paused       bool
    workers      []string // active recording worker URLs (one per destination)
    sessionRoot  string   // per-surgery folder name (base only)
    destinations []string // absolute directories where we mirror outputs
    preset       string
}

var app = &Server{preset: "arthroscopy"}

// POST /api/session/start
// Body: { "doctor":..., "hospital":..., "surgery":..., "patient":..., "technician":... }
func handleSessionStart(w http.ResponseWriter, r *http.Request) {
	var meta map[string]any
	if err := json.NewDecoder(r.Body).Decode(&meta); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	// Create a base folder name. Example: 20250101_153012_Surgery_Patient
	when := time.Now().Format("20060102_150405")
	surgery, _ := meta["surgery"].(string)
	patient, _ := meta["patient"].(string)
	base := fmt.Sprintf("%s_%s_%s", when, safeName(surgery), safeName(patient))

	dests := findDestinations()
	if len(dests) == 0 {
		// Fallback: use local images/ and videos/ dirs (development)
		dests = []string{filepath.Join("images"), filepath.Join("videos")}
	}

	app.mu.Lock()
	app.sessionRoot = base
	app.destinations = dests
	app.mu.Unlock()

	// Create per-destination directories and write metadata.json
	meta["_folder"] = base
	for _, root := range dests {
		folder := filepath.Join(root, base)
		_ = os.MkdirAll(folder, 0o755)
		f, err := os.Create(filepath.Join(folder, "metadata.json"))
		if err == nil {
			_ = json.NewEncoder(f).Encode(meta)
			_ = f.Close()
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{"folder": base, "destinations": dests})
}

// POST /api/capture
func handleCapture(w http.ResponseWriter, r *http.Request) {
	app.mu.Lock()
	base := app.sessionRoot
	dests := append([]string(nil), app.destinations...)
	app.mu.Unlock()
	if base == "" {
		base = time.Now().Format("20060102_150405")
	}

	client := createClient()
	defer client.Close()

	// Validate camera present
	var cam lt.Camera
	if err := client.Get("cv40:/0/camera/0", &cam); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
		return
	}
	if cam.Video.Signal == "" { // proceed even if not locked during local dev
		log.Println("warning: camera signal not locked")
	}

	results := []map[string]string{}
	for _, root := range dests {
		folder := filepath.Join(root, base)
		_ = os.MkdirAll(folder, 0o755)
		// Create still capture worker to that folder
		err := client.Post("cv40:/0/camera/0/file", lt.ImageFileWorker{Media: "image/jpeg", Location: folder}, nil)
		if !errors.Is(err, lt.ErrRedirect) {
			results = append(results, map[string]string{"destination": folder, "status": "error", "error": err.Error()})
			continue
		}
		results = append(results, map[string]string{"destination": folder, "status": "ok"})
	}
	broadcastPhotoCapture()
	writeJSON(w, http.StatusOK, map[string]any{"results": results})
}

// POST /api/recording/start
func handleRecordingStart(w http.ResponseWriter, r *http.Request) {
	app.mu.Lock()
	if app.recording {
		app.mu.Unlock()
		writeJSON(w, http.StatusConflict, map[string]string{"error": "already recording"})
		return
	}
	base := app.sessionRoot
	dests := append([]string(nil), app.destinations...)
	app.mu.Unlock()
	if base == "" {
		base = time.Now().Format("20060102_150405")
	}

	client := createClient()
	defer client.Close()

	workers := []string{}
	for _, root := range dests {
		folder := filepath.Join(root, base)
		_ = os.MkdirAll(folder, 0o755)
		err := client.Post("cv40:/0/camera/0/file", lt.VideoFileWorker{
			Media:    "video/mp4",
			Location: folder,
			// Extra: defaults; you can tune HW/codec later on the device
		}, nil)
		if !errors.Is(err, lt.ErrRedirect) {
			writeJSON(w, http.StatusBadGateway, map[string]string{"error": fmt.Sprintf("worker create on %s: %v", folder, err)})
			return
		}
		workerURL := lt.RedirectLocation(err)
		if err := client.Post(workerURL+"/start", nil, nil); err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{"error": fmt.Sprintf("start on %s: %v", folder, err)})
			return
		}
		workers = append(workers, workerURL)
	}

	app.mu.Lock()
	app.recording = true
	app.paused = false
	app.workers = workers
	app.mu.Unlock()

	broadcastRecordingState(true, false)
	writeJSON(w, http.StatusOK, map[string]any{"workers": workers})
}

// POST /api/recording/pause
func handleRecordingPause(w http.ResponseWriter, r *http.Request) {
	app.mu.Lock()
	workers := append([]string(nil), app.workers...)
	app.mu.Unlock()

	client := createClient()
	defer client.Close()
	for _, u := range workers {
		if err := client.Post(u+"/pause", nil, nil); err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
			return
		}
	}
	app.mu.Lock()
	app.paused = true
	app.mu.Unlock()
	broadcastRecordingState(app.recording, true)
	writeJSON(w, http.StatusOK, map[string]string{"status": "paused"})
}

// POST /api/recording/resume
func handleRecordingResume(w http.ResponseWriter, r *http.Request) {
	app.mu.Lock()
	workers := append([]string(nil), app.workers...)
	app.mu.Unlock()

	client := createClient()
	defer client.Close()
	for _, u := range workers {
		if err := client.Post(u+"/start", nil, nil); err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
			return
		}
	}
	app.mu.Lock()
	app.paused = false
	app.mu.Unlock()
	broadcastRecordingState(app.recording, false)
	writeJSON(w, http.StatusOK, map[string]string{"status": "recording"})
}

// POST /api/recording/stop
func handleRecordingStop(w http.ResponseWriter, r *http.Request) {
	app.mu.Lock()
	workers := append([]string(nil), app.workers...)
	app.mu.Unlock()

	client := createClient()
	defer client.Close()
	for _, u := range workers {
		if err := client.Post(u+"/stop", nil, nil); err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
			return
		}
	}
	app.mu.Lock()
	app.recording = false
	app.paused = false
	app.workers = nil
	app.mu.Unlock()
	broadcastRecordingState(false, false)
	writeJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
}

// POST /api/white-balance
func handleWhiteBalance(w http.ResponseWriter, r *http.Request) {
	// Simple implementation: set temperature and balance to defaults/neutral.
	// On the actual device we can compute balance from scene or call an SDK WB op if present.
	client := createClient()
	defer client.Close()
	wb := lt.CameraWhite{Temperature: 6500}
	if err := client.Post("cv40:/0/camera/0/white", &wb, nil); err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}
	broadcastWhiteBalance(true)
	writeJSON(w, http.StatusOK, map[string]string{"status": "white-balance-set"})
}

// GET /api/settings -> aggregate of visuals/colors/white/exposure
func handleGetSettings(w http.ResponseWriter, r *http.Request) {
	client := createClient()
	defer client.Close()
	var visuals lt.CameraVisuals
	var colors lt.CameraColors
	var white lt.CameraWhite
	var exposure lt.CameraExposure
	_ = client.Get("cv40:/0/camera/0/visuals", &visuals)
	_ = client.Get("cv40:/0/camera/0/colors", &colors)
	_ = client.Get("cv40:/0/camera/0/white", &white)
	_ = client.Get("cv40:/0/camera/0/exposure", &exposure)
	writeJSON(w, http.StatusOK, map[string]any{
		"visuals":  visuals,
		"colors":   colors,
		"white":    white,
		"exposure": exposure,
	})
}

// POST /api/settings -> body may contain partial updates
func handlePostSettings(w http.ResponseWriter, r *http.Request) {
	var req map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	client := createClient()
	defer client.Close()
	if v, ok := req["visuals"]; ok {
		var visuals lt.CameraVisuals
		if err := json.Unmarshal(v, &visuals); err == nil {
			_ = client.Post("cv40:/0/camera/0/visuals", &visuals, nil)
			if visuals.Zoom > 0 {
				broadcastParameterChange("Zoom", visuals.Zoom)
			}
			if visuals.Sharpness > 0 {
				broadcastParameterChange("Sharpness", visuals.Sharpness)
			}
		}
	}
	if v, ok := req["colors"]; ok {
		var colors lt.CameraColors
		if err := json.Unmarshal(v, &colors); err == nil {
			_ = client.Post("cv40:/0/camera/0/colors", &colors, nil)
			broadcastParameterChange("Brightness", float64(colors.Brightness))
			broadcastParameterChange("Contrast", float64(colors.Contrast))
			broadcastParameterChange("Saturation", float64(colors.Saturation))
			broadcastParameterChange("Hue", float64(colors.Hue))
		}
	}
	if v, ok := req["white"]; ok {
		var white lt.CameraWhite
		if err := json.Unmarshal(v, &white); err == nil {
			_ = client.Post("cv40:/0/camera/0/white", &white, nil)
			broadcastParameterChange("Temperature", float64(white.Temperature))
		}
	}
	if v, ok := req["exposure"]; ok {
		var exposure lt.CameraExposure
		if err := json.Unmarshal(v, &exposure); err == nil {
			_ = client.Post("cv40:/0/camera/0/exposure", &exposure, nil)
			broadcastParameterChange("LowLightGain", exposure.LowLightGain)
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// POST /api/presets/apply
func handleApplyPreset(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Preset string `json:"preset"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	client := createClient()
	defer client.Close()

    switch req.Preset {
	case "arthroscopy":
		// Default arthroscopy settings
		colors := lt.CameraColors{
			Brightness: 10,
			Contrast:   15,
			Saturation: 5,
			Hue:        0,
			Gamma:      1.0,
			ColorGain:  [3]float64{1.0, 1.0, 1.0},
		}
		visuals := lt.CameraVisuals{
			Zoom:      1.0,
			Sharpness: 0.7,
		}
		white := lt.CameraWhite{
			Temperature: 6500,
		}
		_ = client.Post("cv40:/0/camera/0/colors", &colors, nil)
		_ = client.Post("cv40:/0/camera/0/visuals", &visuals, nil)
		_ = client.Post("cv40:/0/camera/0/white", &white, nil)

	case "red_boost":
		// Red boost preset - enhanced reds and yellows for better tissue contrast
		colors := lt.CameraColors{
			Brightness: 20,
			Contrast:   25,
			Saturation: 15,
			Hue:        -5,
			Gamma:      0.9,
			ColorGain:  [3]float64{1.2, 0.95, 0.85}, // Boost red, slightly reduce green/blue
		}
		visuals := lt.CameraVisuals{
			Zoom:      1.0,
			Sharpness: 0.8,
		}
		white := lt.CameraWhite{
			Temperature: 5800, // Warmer temperature
		}
		_ = client.Post("cv40:/0/camera/0/colors", &colors, nil)
		_ = client.Post("cv40:/0/camera/0/visuals", &visuals, nil)
		_ = client.Post("cv40:/0/camera/0/white", &white, nil)

    default:
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unknown preset"})
        return
    }
    app.mu.Lock()
    app.preset = req.Preset
    app.mu.Unlock()
    broadcastPresetApplied(req.Preset)
    writeJSON(w, http.StatusOK, map[string]string{"preset": req.Preset, "status": "applied"})
}

// POST /api/controller/button
func handleControllerButton(w http.ResponseWriter, r *http.Request) {
    var req struct{
        ID string `json:"id"`
        Press string `json:"press"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        return
    }
    switch req.ID {
    case "photo":
        app.mu.Lock()
        base := app.sessionRoot
        dests := append([]string(nil), app.destinations...)
        app.mu.Unlock()
        if base == "" { base = time.Now().Format("20060102_150405") }
        client := createClient()
        defer client.Close()
        var cam lt.Camera
        if err := client.Get("cv40:/0/camera/0", &cam); err != nil {
            writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
            return
        }
        for _, root := range dests {
            folder := filepath.Join(root, base)
            _ = os.MkdirAll(folder, 0o755)
            err := client.Post("cv40:/0/camera/0/file", lt.ImageFileWorker{Media: "image/jpeg", Location: folder}, nil)
            if !errors.Is(err, lt.ErrRedirect) {
                writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
                return
            }
        }
        broadcastPhotoCapture()
        writeJSON(w, http.StatusOK, map[string]string{"status": "photo"})
    case "record":
        app.mu.Lock()
        recording := app.recording
        paused := app.paused
        base := app.sessionRoot
        dests := append([]string(nil), app.destinations...)
        workers := append([]string(nil), app.workers...)
        app.mu.Unlock()
        client := createClient()
        defer client.Close()
        if req.Press == "long" {
            for _, u := range workers { _ = client.Post(u+"/stop", nil, nil) }
            app.mu.Lock()
            app.recording = false
            app.paused = false
            app.workers = nil
            app.mu.Unlock()
            broadcastRecordingState(false, false)
            writeJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
            return
        }
        if !recording {
            if base == "" { base = time.Now().Format("20060102_150405") }
            newWorkers := []string{}
            for _, root := range dests {
                folder := filepath.Join(root, base)
                _ = os.MkdirAll(folder, 0o755)
                err := client.Post("cv40:/0/camera/0/file", lt.VideoFileWorker{Media: "video/mp4", Location: folder}, nil)
                if !errors.Is(err, lt.ErrRedirect) {
                    writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
                    return
                }
                workerURL := lt.RedirectLocation(err)
                if err := client.Post(workerURL+"/start", nil, nil); err != nil {
                    writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
                    return
                }
                newWorkers = append(newWorkers, workerURL)
            }
            app.mu.Lock()
            app.recording = true
            app.paused = false
            app.workers = newWorkers
            app.mu.Unlock()
            broadcastRecordingState(true, false)
            writeJSON(w, http.StatusOK, map[string]string{"status": "recording"})
            return
        }
        if paused {
            for _, u := range workers { _ = client.Post(u+"/start", nil, nil) }
            app.mu.Lock(); app.paused = false; app.mu.Unlock()
            broadcastRecordingState(true, false)
            writeJSON(w, http.StatusOK, map[string]string{"status": "recording"})
            return
        }
        for _, u := range workers { _ = client.Post(u+"/pause", nil, nil) }
        app.mu.Lock(); app.paused = true; app.mu.Unlock()
        broadcastRecordingState(true, true)
        writeJSON(w, http.StatusOK, map[string]string{"status": "paused"})
    case "wb":
        client := createClient()
        defer client.Close()
        wb := lt.CameraWhite{Temperature: 6500}
        if err := client.Post("cv40:/0/camera/0/white", &wb, nil); err != nil {
            writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
            return
        }
        broadcastWhiteBalance(true)
        writeJSON(w, http.StatusOK, map[string]string{"status": "white-balance-set"})
    case "preset":
        app.mu.Lock(); current := app.preset; app.mu.Unlock()
        next := "arthroscopy"
        if current == "arthroscopy" { next = "red_boost" }
        client := createClient()
        defer client.Close()
        switch next {
        case "arthroscopy":
            colors := lt.CameraColors{Brightness: 10, Contrast: 15, Saturation: 5, Hue: 0, Gamma: 1.0, ColorGain: [3]float64{1.0,1.0,1.0}}
            visuals := lt.CameraVisuals{Zoom: 1.0, Sharpness: 0.7}
            white := lt.CameraWhite{Temperature: 6500}
            _ = client.Post("cv40:/0/camera/0/colors", &colors, nil)
            _ = client.Post("cv40:/0/camera/0/visuals", &visuals, nil)
            _ = client.Post("cv40:/0/camera/0/white", &white, nil)
        case "red_boost":
            colors := lt.CameraColors{Brightness: 20, Contrast: 25, Saturation: 15, Hue: -5, Gamma: 0.9, ColorGain: [3]float64{1.2,0.95,0.85}}
            visuals := lt.CameraVisuals{Zoom: 1.0, Sharpness: 0.8}
            white := lt.CameraWhite{Temperature: 5800}
            _ = client.Post("cv40:/0/camera/0/colors", &colors, nil)
            _ = client.Post("cv40:/0/camera/0/visuals", &visuals, nil)
            _ = client.Post("cv40:/0/camera/0/white", &white, nil)
        }
        app.mu.Lock(); app.preset = next; app.mu.Unlock()
        broadcastPresetApplied(next)
        writeJSON(w, http.StatusOK, map[string]string{"preset": next, "status": "applied"})
    default:
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unknown controller id"})
    }
}

func safeName(s string) string {
	if s == "" { return "NA" }
	clean := s
	forbidden := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, f := range forbidden {
		clean = strings.ReplaceAll(clean, f, "_")
	}
	return clean
}

// findDestinations returns absolute output directories for mirroring.
// DEV: If none found, we return cv40_camera_app/backend/images and videos as fallbacks.
func findDestinations() []string {
	// TODO: Implement WMI removable-drive discovery for Windows. For now, dev-only defaults.
	wd, _ := os.Getwd()
	return []string{ filepath.Join(wd, "images"), filepath.Join(wd, "videos") }
}
