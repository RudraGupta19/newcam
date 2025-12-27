package main

import (
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
		"camera/0",
	}
	for _, source := range sources {
		var input lt.Camera
		if err := client.Get("cv40:/0/"+source, &input); err != nil {
			log.Fatal(err)
		}
		if input.Video.Signal == "locked" {
			sourceURL = source
			break
		}
	}

	// No active input video source found
	if sourceURL == "" {
		log.Fatal("no active input video source found")
	}

	// Set hdmi-out source
	if err := client.Post("cv40:/0/hdmi-out/0", lt.JSON{"source": sourceURL}, nil); err != nil {
		log.Fatal(err)
	}

	// Print hdmi-out source
	fmt.Printf("hdmi-out source: %s\n", sourceURL)

	// Get hdmi-out source
	var output lt.Output
	if err := client.Get("cv40:/0/hdmi-out/0", &output); err != nil {
		log.Fatal(err)
	}

	// Initialize canvas
	if err := client.Post("cv40:/canvas/0/init", lt.CanvasInit{Size: output.Video.Size}, nil); err != nil {
		log.Fatal(err)
	}

	// Draw text
	if err := client.Post("cv40:/canvas/0/text",
		lt.CanvasText{
			Text:     "Hello, World!",
			FontSize: 200, Color: [4]int{255, 255, 255, 255},
			Size: output.Video.Size,
		},
		nil); err != nil {
		log.Fatal(err)
	}

	/* Add your own drawing code here
	 * Images (PNG, JPEG, ...)
	 * Rectangles
	 * Lines
	 * Ellipses
	 * ...
	 */

	// Set hdmi-out overlay
	if err := client.Post("cv40:/0/hdmi-out/0", lt.JSON{"overlay": "canvas/0"}, nil); err != nil {
		log.Fatal(err)
	}

}
