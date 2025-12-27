package events

import (
    "encoding/json"
    "net/http"
    "sync"
    "time"
    "github.com/gorilla/websocket"
)

type Hub struct {
    mu    sync.Mutex
    conns map[*websocket.Conn]bool
    up    websocket.Upgrader
}

type Event struct {
    Type string                 `json:"type"`
    Timestamp int64             `json:"timestamp"`
    Data map[string]interface{} `json:"data"`
}

func NewHub() *Hub {
    return &Hub{conns: make(map[*websocket.Conn]bool), up: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}}
}

func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request) {
    c, err := h.up.Upgrade(w, r, nil)
    if err != nil { return }
    h.mu.Lock()
    h.conns[c] = true
    h.mu.Unlock()
    for {
        if _, _, err := c.ReadMessage(); err != nil {
            h.mu.Lock()
            delete(h.conns, c)
            h.mu.Unlock()
            break
        }
    }
}

func (h *Hub) Broadcast(t string, data map[string]interface{}) {
    h.mu.Lock()
    defer h.mu.Unlock()
    e := Event{Type: t, Timestamp: time.Now().UnixMilli(), Data: data}
    b, _ := json.Marshal(e)
    for c := range h.conns {
        _ = c.WriteMessage(websocket.TextMessage, b)
    }
}
