package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type MonitorServer struct {
	mu          sync.Mutex
	connections map[*websocket.Conn]bool
	upgrader    websocket.Upgrader
}

var monitorServer = &MonitorServer{
	connections: make(map[*websocket.Conn]bool),
	upgrader: websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	},
}

type OSDEvent struct {
	Type      string                 `json:"type"`
	Timestamp int64                  `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

func (ms *MonitorServer) broadcast(event OSDEvent) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	event.Timestamp = time.Now().UnixMilli()
	data, _ := json.Marshal(event)
	
	for conn := range ms.connections {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			conn.Close()
			delete(ms.connections, conn)
		}
	}
}

func handleMonitorWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := monitorServer.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	monitorServer.mu.Lock()
	monitorServer.connections[conn] = true
	monitorServer.mu.Unlock()

	// Send initial state
	monitorServer.broadcast(OSDEvent{
		Type: "connected",
		Data: map[string]interface{}{
			"recording": app.recording,
			"paused":    app.paused,
		},
	})

	// Keep connection alive and handle incoming messages
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			monitorServer.mu.Lock()
			delete(monitorServer.connections, conn)
			monitorServer.mu.Unlock()
			break
		}
	}
}

// Broadcast recording state changes
func broadcastRecordingState(recording, paused bool) {
	monitorServer.broadcast(OSDEvent{
		Type: "recording_state",
		Data: map[string]interface{}{
			"recording": recording,
			"paused":    paused,
		},
	})
}

// Broadcast parameter changes
func broadcastParameterChange(parameter string, value float64) {
	monitorServer.broadcast(OSDEvent{
		Type: "parameter_change",
		Data: map[string]interface{}{
			"parameter": parameter,
			"value":     value,
		},
	})
}

// Broadcast photo capture
func broadcastPhotoCapture() {
	monitorServer.broadcast(OSDEvent{
		Type: "photo_captured",
		Data: map[string]interface{}{
			"timestamp": time.Now().UnixMilli(),
		},
	})
}

// Broadcast white balance progress
func broadcastWhiteBalance(complete bool) {
	monitorServer.broadcast(OSDEvent{
		Type: "white_balance",
		Data: map[string]interface{}{
			"complete": complete,
		},
	})
}

func broadcastPresetApplied(preset string) {
    monitorServer.broadcast(OSDEvent{
        Type: "preset_applied",
        Data: map[string]interface{}{
            "preset": preset,
        },
    })
}

func startMonitorServer() {
	router := mux.NewRouter()
	
	// WebSocket endpoint for OSD events
	router.HandleFunc("/monitor/ws", handleMonitorWebSocket)
	
	// Static file serving for monitor display
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../monitor_static/")))
	
	// Enable CORS
	router.Use(corsMiddleware)
	
	log.Println("Monitor server starting on :8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}
