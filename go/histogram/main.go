package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

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

	// Create worker
	err := client.Post(sourceURL+"/data", lt.ImageDataWorker{Media: "image/yuv422"}, nil)
	if !errors.Is(err, lt.ErrRedirect) {
		log.Fatal("worker creation failed:", err)
	}
	workerURL := lt.RedirectLocation(err)

	// Fetch worker response
	var worker lt.Worker
	if err := client.Get(workerURL, &worker); err != nil {
		log.Fatal(err)
	}

	// Check packet
	if len(worker.Packets) == 0 {
		log.Fatal("worker packet: not found")
	}
	packet := worker.Packets[0]
	defer packet.Close()

	// Packet metadata
	var meta lt.ImageMetadata
	if err := json.Unmarshal(packet.Meta, &meta); err != nil {
		log.Fatal("worker packet metadata:", err)
	}

	// Print infos
	fmt.Printf("%s %dx%d %d bytes\n", sourceURL, meta.Size[0], meta.Size[1], len(packet.Data))

	// Create histogram
	n := meta.Size[0] * meta.Size[1]
	h := [8]int{}
	for i := 0; i < n; i++ {
		h[packet.Data[i]/32]++
	}

	// Print histogram
	fmt.Println("")
	fmt.Printf("  0.. 31: %.1f%%\n", 100*float64(h[0])/float64(n))
	fmt.Printf(" 32.. 63: %.1f%%\n", 100*float64(h[1])/float64(n))
	fmt.Printf(" 64.. 95: %.1f%%\n", 100*float64(h[2])/float64(n))
	fmt.Printf(" 96..127: %.1f%%\n", 100*float64(h[3])/float64(n))
	fmt.Printf("128..159: %.1f%%\n", 100*float64(h[4])/float64(n))
	fmt.Printf("160..191: %.1f%%\n", 100*float64(h[5])/float64(n))
	fmt.Printf("192..223: %.1f%%\n", 100*float64(h[6])/float64(n))
	fmt.Printf("224..256: %.1f%%\n", 100*float64(h[7])/float64(n))
}
