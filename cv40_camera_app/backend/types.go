package main

import (
	"path/filepath"
	"time"
)

// Camera parameter groups mirroring lt schema (subset for Phase 1)
type Colors struct {
	Hue        int   `json:"hue"`
	Saturation int   `json:"saturation"`
	Brightness int   `json:"brightness"`
	Contrast   int   `json:"contrast"`
	ColorGain  [3]float64 `json:"colorGain"`
}

type Visuals struct {
	Zoom       float64 `json:"zoom"`
	Sharpness  float64 `json:"sharpness"`
}

type White struct {
	Temperature int       `json:"temperature"`
	Balance     [3]float64 `json:"balance"`
}

type Exposure struct {
	LowLightGain float64 `json:"lowLightGain"`
}

// CameraService abstracts real SDK vs simulator
// All paths/sources are assumed camera 0 on board 0.
type CameraService interface {
	CaptureStill(destDir string) (string, error)
	StartRecording(destDir string) error
	PauseRecording() error
	ResumeRecording() error
	StopRecording() error
	WhiteBalance() error

	GetColors() (Colors, error)
	SetColors(Colors) error
	GetVisuals() (Visuals, error)
	SetVisuals(Visuals) error
	GetWhite() (White, error)
	SetWhite(White) error
	GetExposure() (Exposure, error)
	SetExposure(Exposure) error
}

// util
func makeSessionDir(root string) string {
	// Each start creates a new timestamped folder
	return filepath.Join(root, time.Now().Format("20060102_150405"))
}

