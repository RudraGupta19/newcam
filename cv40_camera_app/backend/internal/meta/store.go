package meta

import (
    "encoding/json"
    "os"
    "path/filepath"
)

type SessionMeta struct {
    SessionID string `json:"sessionId"`
    Doctor string `json:"doctor"`
    Hospital string `json:"hospital"`
    Patient string `json:"patient"`
    SurgeryType string `json:"surgeryType"`
}

func Write(dir string, m SessionMeta) error {
    p := filepath.Join(dir, "meta.json")
    b, _ := json.MarshalIndent(m, "", "  ")
    return os.WriteFile(p, b, 0o644)
}
