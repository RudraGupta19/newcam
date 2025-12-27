package main

import (
    "log"
    "os"
    "path/filepath"
    "cv40-camera-backend/internal/api"
    "cv40-camera-backend/internal/config"
    "cv40-camera-backend/internal/cv40"
    "cv40-camera-backend/internal/events"
    "cv40-camera-backend/internal/overlay"
    "cv40-camera-backend/internal/state"
    "cv40-camera-backend/internal/storage"
)

func main() {
    cfgPath := os.Getenv("CV40_CONFIG")
    if cfgPath == "" {
        cfgPath = filepath.Join(".", "config.json")
    }
    cfg, err := config.Load(cfgPath)
    if err != nil {
        log.Fatal(err)
    }

    ev := events.NewHub()
    st := state.NewStore()
    st.Set(state.BOOTING)

    client := cv40.NewRealClient(cfg)
    if err := client.Health(); err != nil {
        st.Set(state.ERROR_BLOCKING)
        log.Fatal(err)
    }

    ov := overlay.NewEngine(client, cfg)
    if err := ov.InitOutput(); err != nil {
        log.Println("overlay init:", err)
    }

    sm := storage.NewManager(cfg)
    if err := sm.InitTargets(); err != nil {
        st.Set(state.DEGRADED)
        log.Println("storage init:", err)
    }

    st.Set(state.READY)

    srv := api.NewServer(cfg, client, ov, sm, st, ev)
    log.Fatal(srv.Start())
}
