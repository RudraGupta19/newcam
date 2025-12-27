package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	// Start monitor server in a separate goroutine
	go startMonitorServer()

	router := mux.NewRouter()

	// API routes
	router.HandleFunc("/api/session/start", handleSessionStart).Methods("POST")
	router.HandleFunc("/api/capture", handleCapture).Methods("POST")
	router.HandleFunc("/api/recording/start", handleRecordingStart).Methods("POST")
	router.HandleFunc("/api/recording/stop", handleRecordingStop).Methods("POST")
	router.HandleFunc("/api/recording/pause", handleRecordingPause).Methods("POST")
	router.HandleFunc("/api/recording/resume", handleRecordingResume).Methods("POST")
	router.HandleFunc("/api/white-balance", handleWhiteBalance).Methods("POST")
	router.HandleFunc("/api/settings", handleGetSettings).Methods("GET")
	router.HandleFunc("/api/settings", handlePostSettings).Methods("POST")
	router.HandleFunc("/api/presets/apply", handleApplyPreset).Methods("POST")
	router.HandleFunc("/api/controller/button", handleControllerButton).Methods("POST")

	// Enable CORS
	router.Use(corsMiddleware)

	log.Println("Server starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
