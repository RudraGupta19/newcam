package tools

import (
    "time"
    lt "lt/client/go"
    "cv40-camera-backend/internal/config"
    "cv40-camera-backend/internal/cv40"
    "cv40-camera-backend/internal/overlay"
)

type Limiter struct {
    ranges config.SafeRanges
    cli *cv40.RealClient
    ov *overlay.Engine
    colorsCh chan lt.CameraColors
    visualsCh chan lt.CameraVisuals
}

func NewLimiter(cli *cv40.RealClient, ov *overlay.Engine, ranges config.SafeRanges) *Limiter {
    l := &Limiter{cli: cli, ov: ov, ranges: ranges, colorsCh: make(chan lt.CameraColors, 8), visualsCh: make(chan lt.CameraVisuals, 8)}
    l.start()
    return l
}

func (l *Limiter) start() {
    go func() {
        var last lt.CameraColors
        ticker := time.NewTicker(75 * time.Millisecond)
        defer ticker.Stop()
        pending := false
        for {
            select {
            case v := <-l.colorsCh:
                v = l.clampColors(v)
                last = v
                pending = true
            case <-ticker.C:
                if pending {
                    _ = l.cli.SetColors(last)
                    l.ov.Slider("Colors", "applied", 800)
                    pending = false
                }
            }
        }
    }()
    go func() {
        var last lt.CameraVisuals
        ticker := time.NewTicker(75 * time.Millisecond)
        defer ticker.Stop()
        pending := false
        for {
            select {
            case v := <-l.visualsCh:
                v = l.clampVisuals(v)
                last = v
                pending = true
            case <-ticker.C:
                if pending {
                    _ = l.cli.SetVisuals(last)
                    l.ov.Slider("Visuals", "applied", 800)
                    pending = false
                }
            }
        }
    }()
}

func (l *Limiter) clampColors(v lt.CameraColors) lt.CameraColors {
    if v.Brightness < l.ranges.Brightness[0] { v.Brightness = l.ranges.Brightness[0] }
    if v.Brightness > l.ranges.Brightness[1] { v.Brightness = l.ranges.Brightness[1] }
    if v.Contrast < l.ranges.Contrast[0] { v.Contrast = l.ranges.Contrast[0] }
    if v.Contrast > l.ranges.Contrast[1] { v.Contrast = l.ranges.Contrast[1] }
    if v.Saturation < l.ranges.Saturation[0] { v.Saturation = l.ranges.Saturation[0] }
    if v.Saturation > l.ranges.Saturation[1] { v.Saturation = l.ranges.Saturation[1] }
    if v.Hue < l.ranges.Hue[0] { v.Hue = l.ranges.Hue[0] }
    if v.Hue > l.ranges.Hue[1] { v.Hue = l.ranges.Hue[1] }
    return v
}

func (l *Limiter) clampVisuals(v lt.CameraVisuals) lt.CameraVisuals {
    if v.Zoom < l.ranges.Zoom[0] { v.Zoom = l.ranges.Zoom[0] }
    if v.Zoom > l.ranges.Zoom[1] { v.Zoom = l.ranges.Zoom[1] }
    if v.Sharpness < l.ranges.Sharpness[0] { v.Sharpness = l.ranges.Sharpness[0] }
    if v.Sharpness > l.ranges.Sharpness[1] { v.Sharpness = l.ranges.Sharpness[1] }
    return v
}

func (l *Limiter) Colors(v lt.CameraColors) { select { case l.colorsCh <- v: default: } }
func (l *Limiter) Visuals(v lt.CameraVisuals) { select { case l.visualsCh <- v: default: } }
