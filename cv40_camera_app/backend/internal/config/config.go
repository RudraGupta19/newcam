package config

import (
    "encoding/json"
    "errors"
    "os"
)

type RecordingDefaults struct {
    Codec   string  `json:"codec"`
    Encoder string  `json:"encoder"`
    Bitrate int     `json:"bitrate"`
    Container string `json:"container"`
}

type SafeRanges struct {
    Brightness [2]int     `json:"brightness"`
    Contrast   [2]int     `json:"contrast"`
    Saturation [2]int     `json:"saturation"`
    Hue        [2]int     `json:"hue"`
    Sharpness  [2]float64 `json:"sharpness"`
    Zoom       [2]float64 `json:"zoom"`
    Temperature [2]int    `json:"temperature"`
    LowLightGain [2]float64 `json:"lowLightGain"`
}

type OverlaySpec struct {
    CanvasID int    `json:"canvasId"`
    Output   string `json:"output"`
}

type Config struct {
    BaseURL string `json:"baseUrl"`
    BoardID int    `json:"boardId"`
    CameraID int   `json:"cameraId"`
    OutputID string `json:"outputId"`
    StorageRoots []string `json:"storageRoots"`
    FreeSpaceGB int `json:"freeSpaceGb"`
    Recording RecordingDefaults `json:"recording"`
    Ranges SafeRanges `json:"ranges"`
    Overlay OverlaySpec `json:"overlay"`
}

func Load(path string) (Config, error) {
    var c Config
    b, err := os.ReadFile(path)
    if err != nil {
        return c, err
    }
    if err := json.Unmarshal(b, &c); err != nil {
        return c, err
    }
    if c.CameraID < 0 || c.BoardID < 0 {
        return c, errors.New("invalid board/camera id")
    }
    return c, nil
}
