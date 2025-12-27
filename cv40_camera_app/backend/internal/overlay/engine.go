package overlay

import (
    lt "lt/client/go"
    "cv40-camera-backend/internal/config"
    "cv40-camera-backend/internal/cv40"
)

type Engine struct {
    cli *cv40.RealClient
    cfg config.Config
}

func NewEngine(cli *cv40.RealClient, cfg config.Config) *Engine { return &Engine{cli: cli, cfg: cfg} }

func (e *Engine) InitOutput() error {
    if err := e.cli.Health(); err != nil { return err }
    out, err := e.cli.GetOutput(e.cfg.Overlay.Output)
    if err != nil { return err }
    size := out.Video.Size
    if size == [2]int{0,0} { size = [2]int{1920,1080} }
    if err := e.cli.CanvasInit(e.cfg.Overlay.CanvasID, size); err != nil { return err }
    return e.cli.ConfigureOutputOverlay(e.cfg.Overlay.Output, e.cfg.Overlay.CanvasID)
}

func (e *Engine) SetRecordingIndicator(active bool, paused bool) {
    s := "REC"
    if paused { s = "PAUSED" }
    if !active { s = "" }
    if s == "" { s = "" }
    _ = e.cli.CanvasText(e.cfg.Overlay.CanvasID, lt.CanvasText{Text: s, FontSize: 48, Color: [4]int{255,0,0,200}, Size: [2]int{1920,1080}})
}

func (e *Engine) Toast(text string, ms int) {
    _ = e.cli.CanvasText(e.cfg.Overlay.CanvasID, lt.CanvasText{Text: text, FontSize: 38, Color: [4]int{255,255,255,255}, Size: [2]int{1920,1080}})
}

func (e *Engine) Slider(name string, value string, ms int) {
    _ = e.cli.CanvasText(e.cfg.Overlay.CanvasID, lt.CanvasText{Text: name+": "+value, FontSize: 32, Color: [4]int{0,180,255,255}, Size: [2]int{1920,1080}})
}

func (e *Engine) DriveWarning(text string, ms int) {
    _ = e.cli.CanvasText(e.cfg.Overlay.CanvasID, lt.CanvasText{Text: text, FontSize: 34, Color: [4]int{255,255,0,255}, Size: [2]int{1920,1080}})
}
