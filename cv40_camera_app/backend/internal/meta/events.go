package meta

import (
    "encoding/json"
    "os"
    "path/filepath"
    "time"
)

type Event struct {
    Type string `json:"type"`
    Timestamp int64 `json:"timestamp"`
    Data map[string]any `json:"data"`
}

func AppendEvent(sessionDir string, e Event) error {
    p := filepath.Join(sessionDir, "logs", "events.jsonl")
    b, _ := json.Marshal(e)
    b = append(b, '\n')
    f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
    if err != nil { return err }
    defer f.Close()
    _, err = f.Write(b)
    return err
}

func NewEvent(t string, data map[string]any) Event {
    return Event{Type: t, Timestamp: time.Now().UnixMilli(), Data: data}
}
