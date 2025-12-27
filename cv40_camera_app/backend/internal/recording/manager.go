package recording

import (
    "path/filepath"
    "time"
    "os"
    "cv40-camera-backend/internal/cv40"
)

type Job struct { URL string; Target string }

type JobStatus struct { Job Job; Status string; Error string }

type Manager struct {
    cli *cv40.RealClient
    jobs []Job
    pollStop chan struct{}
    onUpdate func([]JobStatus)
}

func NewManager(cli *cv40.RealClient) *Manager { return &Manager{cli: cli} }

func (m *Manager) Start(destDirs []string, media string) ([]Job, error) {
    m.jobs = nil
    for _, d := range destDirs {
        u, err := m.cli.CreateVideoWorker(d, media)
        if err != nil { return nil, err }
        if err := m.cli.StartWorker(u); err != nil { return nil, err }
        m.jobs = append(m.jobs, Job{URL: u, Target: d})
    }
    m.startPolling()
    return m.jobs, nil
}

func (m *Manager) startPolling() {
    if m.pollStop != nil { close(m.pollStop) }
    m.pollStop = make(chan struct{})
    go func(stop <-chan struct{}) {
        ticker := time.NewTicker(300 * time.Millisecond)
        defer ticker.Stop()
        for {
            select {
            case <-stop:
                return
            case <-ticker.C:
                statuses := []JobStatus{}
                for _, j := range m.jobs {
                    w, err := m.cli.GetWorker(j.URL)
                    if err != nil { statuses = append(statuses, JobStatus{Job: j, Status: "FAILED", Error: err.Error()}); continue }
                    statuses = append(statuses, JobStatus{Job: j, Status: w.Status})
                }
                if m.onUpdate != nil { m.onUpdate(statuses) }
            }
        }
    }(m.pollStop)
}

func (m *Manager) OnUpdate(fn func([]JobStatus)) { m.onUpdate = fn }

func (m *Manager) Pause() error {
    for _, j := range m.jobs { if err := m.cli.PauseWorker(j.URL); err != nil { return err } }
    return nil
}

func (m *Manager) Resume() error {
    for _, j := range m.jobs { if err := m.cli.StartWorker(j.URL); err != nil { return err } }
    return nil
}

type RecordingResult struct { Target string; File string; Size int64 }

func (m *Manager) Stop() ([]RecordingResult, error) {
    for _, j := range m.jobs { if err := m.cli.StopWorker(j.URL); err != nil { return nil, err } }
    done := make(chan struct{})
    go func() {
        timeout := time.After(5 * time.Second)
        for {
            select {
            case <-timeout:
                close(done)
                return
            default:
                allStopped := true
                for _, j := range m.jobs {
                    w, err := m.cli.GetWorker(j.URL)
                    if err != nil { continue }
                    if w.Status != "stopped" && w.Status != "finalized" { allStopped = false }
                }
                if allStopped { close(done); return }
                time.Sleep(200 * time.Millisecond)
            }
        }
    }()
    <-done
    results := m.verifyFiles()
    if m.pollStop != nil { close(m.pollStop) }
    m.jobs = nil
    return results, nil
}

func (m *Manager) verifyFiles() []RecordingResult {
    results := []RecordingResult{}
    for _, j := range m.jobs {
        dir := j.Target
        matches, _ := filepath.Glob(filepath.Join(dir, "*.mp4"))
        var best string
        var bestSize int64
        for _, p := range matches {
            fi, err := os.Stat(p)
            if err == nil && fi.Size() > bestSize { best = p; bestSize = fi.Size() }
        }
        results = append(results, RecordingResult{Target: dir, File: best, Size: bestSize})
    }
    return results
}
