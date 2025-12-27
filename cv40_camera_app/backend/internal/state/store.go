package state

import "sync"

type Status string

const (
    BOOTING        Status = "BOOTING"
    READY          Status = "READY"
    SESSION_ACTIVE Status = "SESSION_ACTIVE"
    RECORDING      Status = "RECORDING"
    PAUSED         Status = "PAUSED"
    DEGRADED       Status = "DEGRADED"
    ERROR_BLOCKING Status = "ERROR_BLOCKING"
)

type Store struct {
    mu    sync.RWMutex
    state Status
}

func NewStore() *Store { return &Store{state: BOOTING} }

func (s *Store) Get() Status {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.state
}

func (s *Store) Set(st Status) {
    s.mu.Lock()
    s.state = st
    s.mu.Unlock()
}
