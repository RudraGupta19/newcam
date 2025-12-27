package api

import (
    "encoding/json"
    "bytes"
    "io"
    "log"
    "net/http"
    "path/filepath"
    "time"
    "github.com/gorilla/mux"
    lt "lt/client/go"
    "cv40-camera-backend/internal/config"
    "cv40-camera-backend/internal/cv40"
    "cv40-camera-backend/internal/events"
    "cv40-camera-backend/internal/meta"
    "cv40-camera-backend/internal/overlay"
    "cv40-camera-backend/internal/recording"
    "cv40-camera-backend/internal/state"
    "cv40-camera-backend/internal/storage"
    "cv40-camera-backend/internal/tools"
)

type Server struct {
    cfg config.Config
    cli *cv40.RealClient
    ov  *overlay.Engine
    sm  *storage.Manager
    st  *state.Store
    ev  *events.Hub
    rec *recording.Manager
    lim *tools.Limiter
    sessionID string
    sessionDirs []string
    curPreset string
    stopping bool
}

func NewServer(cfg config.Config, cli *cv40.RealClient, ov *overlay.Engine, sm *storage.Manager, st *state.Store, ev *events.Hub) *Server {
    s := &Server{cfg: cfg, cli: cli, ov: ov, sm: sm, st: st, ev: ev, rec: recording.NewManager(cli)}
    s.lim = tools.NewLimiter(cli, ov, cfg.Ranges)
    return s
}

func (s *Server) Start() error {
    r := mux.NewRouter()
    r.Use(func(next http.Handler) http.Handler { return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { w.Header().Set("Access-Control-Allow-Origin", "*"); w.Header().Set("Access-Control-Allow-Headers", "Content-Type"); w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS"); if req.Method==http.MethodOptions { w.WriteHeader(http.StatusNoContent); return }; next.ServeHTTP(w, req) }) })
    r.HandleFunc("/health", s.handleHealth).Methods("GET")
    r.HandleFunc("/state", s.handleState).Methods("GET")
    r.HandleFunc("/events", s.ev.HandleWS)
    r.HandleFunc("/tools/session/start", s.handleSessionStart).Methods("POST")
    r.HandleFunc("/tools/record/start", s.handleRecordStart).Methods("POST")
    r.HandleFunc("/tools/record/pause", s.handleRecordPause).Methods("POST")
    r.HandleFunc("/tools/record/resume", s.handleRecordResume).Methods("POST")
    r.HandleFunc("/tools/record/stop", s.handleRecordStop).Methods("POST")
    r.HandleFunc("/tools/photo/capture", s.handlePhotoCapture).Methods("POST")
    r.HandleFunc("/tools/whitebalance/run", s.handleWhiteBalance).Methods("POST")
    r.HandleFunc("/tools/settings/colors", s.handleSetColors).Methods("POST")
    r.HandleFunc("/tools/settings/visuals", s.handleSetVisuals).Methods("POST")
    r.HandleFunc("/tools/preset/apply", s.handlePresetApply).Methods("POST")
    r.HandleFunc("/controller/event", s.handleControllerEvent).Methods("POST")
    log.Println("control-service :8083")
    return http.ListenAndServe(":8083", r)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
    if err := s.cli.Health(); err != nil { w.WriteHeader(http.StatusServiceUnavailable); w.Write([]byte(err.Error())); return }
    w.WriteHeader(http.StatusOK)
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]any{"state": s.st.Get()})
}

func (s *Server) handleSessionStart(w http.ResponseWriter, r *http.Request) {
    var body struct{ Doctor, Hospital, Patient, SurgeryType string }
    _ = json.NewDecoder(r.Body).Decode(&body)
    id := time.Now().Format("20060102_150405")
    dirs := s.sm.SessionDirs(id)
    for _, d := range dirs { _ = meta.Write(d, meta.SessionMeta{SessionID: id, Doctor: body.Doctor, Hospital: body.Hospital, Patient: body.Patient, SurgeryType: body.SurgeryType}) }
    s.sessionID = id
    s.sessionDirs = dirs
    s.st.Set(state.SESSION_ACTIVE)
    s.ev.Broadcast("session_started", map[string]interface{}{"sessionId": id})
    for _, d := range dirs { _ = meta.AppendEvent(d, meta.NewEvent("session_started", map[string]any{"sessionId": id})) }
    colors, _ := s.cli.GetColors()
    visuals, _ := s.cli.GetVisuals()
    white, _ := s.cli.GetWhite()
    exposure, _ := s.cli.GetExposure()
    snap := map[string]any{"colors": colors, "visuals": visuals, "white": white, "exposure": exposure}
    for _, d := range dirs { _ = meta.AppendEvent(d, meta.NewEvent("settings_snapshot", snap)) }
    json.NewEncoder(w).Encode(map[string]any{"sessionId": id, "dirs": dirs})
}

func (s *Server) handleRecordStart(w http.ResponseWriter, r *http.Request) {
    if s.stopping { w.WriteHeader(http.StatusConflict); return }
    if s.st.Get() != state.SESSION_ACTIVE && s.st.Get() != state.READY { w.WriteHeader(http.StatusConflict); return }
    if s.sessionID == "" { w.WriteHeader(http.StatusConflict); w.Write([]byte("no active session")); return }
    dirs := s.sessionDirs
    outs := []string{}
    for _, d := range dirs { outs = append(outs, filepath.Join(d, "video")) }
    jobs, err := s.rec.Start(outs, "video/mp4")
    if err != nil { w.WriteHeader(http.StatusBadGateway); w.Write([]byte(err.Error())); return }
    s.rec.OnUpdate(func(sts []recording.JobStatus){
        active := 0; failed := 0
        for _, sjs := range sts { if sjs.Status == "ACTIVE" || sjs.Status == "recording" { active++ } ; if sjs.Status == "FAILED" { failed++ } }
        if failed > 0 {
            s.st.Set(state.DEGRADED)
            s.ov.DriveWarning("Drive failure; recording continues", 2000)
            s.ev.Broadcast("drive_failure", map[string]interface{}{"failed": failed})
            for _, d := range s.sessionDirs { _ = meta.AppendEvent(d, meta.NewEvent("drive_failure", map[string]any{"failed": failed})) }
        } else if active == 0 {
            s.st.Set(state.ERROR_BLOCKING)
            s.ev.Broadcast("recording_blocked", map[string]interface{}{})
            for _, d := range s.sessionDirs { _ = meta.AppendEvent(d, meta.NewEvent("recording_blocked", map[string]any{})) }
        } else {
            s.st.Set(state.RECORDING)
        }
    })
    s.st.Set(state.RECORDING)
    s.ov.SetRecordingIndicator(true, false)
    s.ev.Broadcast("recording_state", map[string]interface{}{"recording": true, "paused": false, "jobs": jobs})
    for _, d := range dirs { _ = meta.AppendEvent(d, meta.NewEvent("record_start", map[string]any{"jobs": jobs})) }
    json.NewEncoder(w).Encode(map[string]any{"status": "recording", "jobs": jobs})
}

func (s *Server) handleRecordPause(w http.ResponseWriter, r *http.Request) {
    if s.stopping { w.WriteHeader(http.StatusConflict); return }
    if s.st.Get() != state.RECORDING { w.WriteHeader(http.StatusConflict); return }
    if err := s.rec.Pause(); err != nil { w.WriteHeader(http.StatusBadGateway); w.Write([]byte(err.Error())); return }
    s.st.Set(state.PAUSED)
    s.ov.SetRecordingIndicator(true, true)
    s.ev.Broadcast("recording_state", map[string]interface{}{"recording": true, "paused": true})
    for _, d := range s.sessionDirs { _ = meta.AppendEvent(d, meta.NewEvent("record_pause", map[string]any{})) }
    json.NewEncoder(w).Encode(map[string]any{"status": "paused"})
}

func (s *Server) handleRecordResume(w http.ResponseWriter, r *http.Request) {
    if s.stopping { w.WriteHeader(http.StatusConflict); return }
    if s.st.Get() != state.PAUSED { w.WriteHeader(http.StatusConflict); return }
    if err := s.rec.Resume(); err != nil { w.WriteHeader(http.StatusBadGateway); w.Write([]byte(err.Error())); return }
    s.st.Set(state.RECORDING)
    s.ov.SetRecordingIndicator(true, false)
    s.ev.Broadcast("recording_state", map[string]interface{}{"recording": true, "paused": false})
    for _, d := range s.sessionDirs { _ = meta.AppendEvent(d, meta.NewEvent("record_resume", map[string]any{})) }
    json.NewEncoder(w).Encode(map[string]any{"status": "recording"})
}

func (s *Server) handleRecordStop(w http.ResponseWriter, r *http.Request) {
    s.stopping = true
    results, err := s.rec.Stop()
    if err != nil { s.stopping = false; w.WriteHeader(http.StatusBadGateway); w.Write([]byte(err.Error())); return }
    s.st.Set(state.SESSION_ACTIVE)
    s.ov.SetRecordingIndicator(false, false)
    s.ev.Broadcast("recording_state", map[string]interface{}{"recording": false, "paused": false})
    for _, d := range s.sessionDirs { _ = meta.AppendEvent(d, meta.NewEvent("record_stop", map[string]any{"results": results})) }
    json.NewEncoder(w).Encode(map[string]any{"status": "stopped", "results": results})
    s.stopping = false
}

func (s *Server) handlePhotoCapture(w http.ResponseWriter, r *http.Request) {
    dirs := s.sessionDirs
    for _, d := range dirs { _ = s.cli.CaptureStill(filepath.Join(d, "photos")) }
    s.ev.Broadcast("photo_captured", map[string]interface{}{"timestamp": time.Now().UnixMilli()})
    s.ov.Toast("Photo captured", 2000)
    for _, d := range dirs { _ = meta.AppendEvent(d, meta.NewEvent("photo_captured", map[string]any{})) }
    json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
}

func (s *Server) handleWhiteBalance(w http.ResponseWriter, r *http.Request) {
    wb := lt.CameraWhite{Temperature: 6500}
    if err := s.cli.SetWhite(wb); err != nil { w.WriteHeader(http.StatusBadGateway); w.Write([]byte(err.Error())); return }
    s.ev.Broadcast("white_balance", map[string]interface{}{"complete": true})
    s.ov.Toast("White balance complete", 2000)
    for _, d := range s.sessionDirs { _ = meta.AppendEvent(d, meta.NewEvent("white_balance", map[string]any{"complete": true})) }
    json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
}

func (s *Server) handleSetColors(w http.ResponseWriter, r *http.Request) {
    var v lt.CameraColors
    _ = json.NewDecoder(r.Body).Decode(&v)
    s.lim.Colors(v)
    s.ov.Slider("Colors", "pending", 800)
    s.ev.Broadcast("parameter_change", map[string]interface{}{"parameter": "colors", "value": v})
    json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
}

func (s *Server) handleSetVisuals(w http.ResponseWriter, r *http.Request) {
    var v lt.CameraVisuals
    _ = json.NewDecoder(r.Body).Decode(&v)
    s.lim.Visuals(v)
    s.ov.Slider("Visuals", "pending", 800)
    s.ev.Broadcast("parameter_change", map[string]interface{}{"parameter": "visuals", "value": v})
    json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
}

func (s *Server) handlePresetApply(w http.ResponseWriter, r *http.Request) {
    var req struct{ Preset string `json:"preset"` }
    _ = json.NewDecoder(r.Body).Decode(&req)
    switch req.Preset {
    case "arthroscopy":
        colors := lt.CameraColors{Brightness: 10, Contrast: 15, Saturation: 5, Hue: 0}
        visuals := lt.CameraVisuals{Zoom: 1.0, Sharpness: 0.7}
        white := lt.CameraWhite{Temperature: 6500}
        if err := s.cli.SetColors(colors); err != nil { w.WriteHeader(http.StatusBadGateway); w.Write([]byte(err.Error())); return }
        if err := s.cli.SetVisuals(visuals); err != nil { w.WriteHeader(http.StatusBadGateway); w.Write([]byte(err.Error())); return }
        if err := s.cli.SetWhite(white); err != nil { w.WriteHeader(http.StatusBadGateway); w.Write([]byte(err.Error())); return }
    case "red_boost":
        colors := lt.CameraColors{Brightness: 20, Contrast: 25, Saturation: 15, Hue: -5}
        visuals := lt.CameraVisuals{Zoom: 1.0, Sharpness: 0.8}
        white := lt.CameraWhite{Temperature: 5800}
        if err := s.cli.SetColors(colors); err != nil { w.WriteHeader(http.StatusBadGateway); w.Write([]byte(err.Error())); return }
        if err := s.cli.SetVisuals(visuals); err != nil { w.WriteHeader(http.StatusBadGateway); w.Write([]byte(err.Error())); return }
        if err := s.cli.SetWhite(white); err != nil { w.WriteHeader(http.StatusBadGateway); w.Write([]byte(err.Error())); return }
    default:
        w.WriteHeader(http.StatusBadRequest); w.Write([]byte("unknown preset")); return
    }
    s.ov.Toast("Preset applied: "+req.Preset, 1500)
    s.ev.Broadcast("preset_applied", map[string]interface{}{"preset": req.Preset})
    for _, d := range s.sessionDirs { _ = meta.AppendEvent(d, meta.NewEvent("preset_applied", map[string]any{"preset": req.Preset})) }
    json.NewEncoder(w).Encode(map[string]any{"status": "applied", "preset": req.Preset})
}

func (s *Server) handleControllerEvent(w http.ResponseWriter, r *http.Request) {
    var req struct{ DeviceId string `json:"deviceId"`; Btn int `json:"btn"`; Press string `json:"press"` }
    _ = json.NewDecoder(r.Body).Decode(&req)
    switch req.Btn {
    case 1:
        if req.Press == "long" { s.handleRecordStop(w, r) } else { s.handlePhotoCapture(w, r) }
        return
    case 2:
        if req.Press == "long" { s.handleWhiteBalance(w, r); return }
        if s.st.Get() == state.RECORDING { s.handleRecordPause(w, r); return }
        if s.st.Get() == state.PAUSED { s.handleRecordResume(w, r); return }
        s.handleRecordStart(w, r)
        return
    case 3:
        next := "arthroscopy"
        if s.curPreset == "arthroscopy" { next = "red_boost" }
        b, _ := json.Marshal(map[string]string{"preset": next})
        req2 := &http.Request{Method: "POST", Body: io.NopCloser(bytes.NewReader(b))}
        s.handlePresetApply(w, req2)
        s.curPreset = next
        return
    case 4:
        var v lt.CameraVisuals
        v.Zoom = 1.1
        _ = s.cli.SetVisuals(v)
        s.ov.Slider("Zoom", "1.1x", 1000)
        json.NewEncoder(w).Encode(map[string]any{"status": "zoom"})
        return
    }
    w.WriteHeader(http.StatusBadRequest)
}
