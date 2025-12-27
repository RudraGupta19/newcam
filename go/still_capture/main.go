package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	lt "lt/client/go"
)

func main() {
	var client lt.Client
	defer client.Close()

	// Find first active input video source
	var sourceURL string
	sources := []string{
		"cv40:/0/camera/0",
	}
	for _, source := range sources {
		var input lt.Camera
		if err := client.Get(source, &input); err != nil {
			log.Fatal(err)
		}
		if input.Video.Signal == "locked" {
			sourceURL = source
			break
		}
	}

	// If no active input video source found, use a default one
	if sourceURL == "" {
		fmt.Println("no active input video source found")
		sourceURL = "cv40:/0/camera/0"
	}

	// Select an absolute path to a capture directory
	wd, err := os.Getwd() // Current working directory
	if err != nil {
		log.Fatal("working directory:", err)
	}

	// Create worker
	err = client.Post(sourceURL+"/file", lt.ImageFileWorker{Media: "image/jpeg", Location: wd}, nil)
	if !errors.Is(err, lt.ErrRedirect) {
		log.Fatal("worker creation failed:", err)
	}
	workerURL := lt.RedirectLocation(err)

	// Fetch worker response
	var worker lt.Worker
	if err := client.Get(workerURL, &worker); err != nil {
		log.Fatal(err)
	}

	// Packet data
	if len(worker.Packets) == 0 {
		log.Fatal("worker packet: not found")
	}
	packet := worker.Packets[0]
	defer packet.Close() // Packet has to be manually released

	// Packet metadata
	var meta lt.ImageMetadata
	if err := json.Unmarshal(packet.Meta, &meta); err != nil {
		log.Fatal("worker packet metadata:", err)
	}

	// Print file information
	fmt.Printf("%s %s %dx%d %d bytes\n", sourceURL, worker.Name, meta.Size[0], meta.Size[1], worker.Length)
}
