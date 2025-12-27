package storage

import (
    "errors"
    "os"
    "path/filepath"
    "cv40-camera-backend/internal/config"
)

type Target struct {
    Root string
}

type Manager struct {
    cfg config.Config
    targets []Target
}

func NewManager(cfg config.Config) *Manager { return &Manager{cfg: cfg} }

func (m *Manager) InitTargets() error {
    roots := m.cfg.StorageRoots
    if len(roots) == 0 {
        wd, _ := os.Getwd()
        roots = []string{filepath.Join(wd, "images"), filepath.Join(wd, "videos")}
    }
    m.targets = nil
    for _, r := range roots {
        if err := os.MkdirAll(r, 0o755); err != nil { return err }
        m.targets = append(m.targets, Target{Root: r})
    }
    if len(m.targets) == 0 { return errors.New("no storage targets") }
    return nil
}

func (m *Manager) SessionDirs(session string) []string {
    out := []string{}
    for _, t := range m.targets {
        d := filepath.Join(t.Root, "Sessions", session)
        _ = os.MkdirAll(filepath.Join(d, "video"), 0o755)
        _ = os.MkdirAll(filepath.Join(d, "photos"), 0o755)
        _ = os.MkdirAll(filepath.Join(d, "logs"), 0o755)
        out = append(out, d)
    }
    return out
}
