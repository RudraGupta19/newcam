package main

import (
	"log"

	lt "lt/client/go"
)

func main() {
	// Initialize the client to communicate with the server
	var client lt.Client
	defer client.Close()

	// Flag to track if the visuals have changed
	var zoomChanged bool

	// Get the current zoom level
	var visuals lt.CameraVisuals
	if err := client.Get("cv40:/0/camera/0/visuals", &visuals); err != nil {
		log.Fatal(err)
	}

	for {
		// URL to get the state of the buttons
		sourceURL := "cv40:/0/camera/0/buttons" // "cv40:/0/buttons" / "cv40:/0/camera/0/buttons"

		// Get buttons state
		var buttons lt.Buttons
		if err := client.Get(sourceURL, &buttons); err != nil {
			log.Fatal(err)
		}

		if len(buttons.Buttons) == 0 {
			log.Fatal("Not enough buttons available")
		}

		// Zoom in if the first button is pressed
		if !buttons.Buttons[0].Pressed {
			visuals.Zoom += 0.1
			if visuals.Zoom >= 4.0 {
				visuals.Zoom = 4.0
			}
			zoomChanged = true
		}

		// Zoom out if the second button is pressed
		if !buttons.Buttons[1].Pressed {
			visuals.Zoom -= 0.1
			if visuals.Zoom <= 1.0 {
				visuals.Zoom = 1.0
			}
			zoomChanged = true
		}

		// Set the new zoom level if it has changed
		if zoomChanged {
			if err := client.Post("cv40:/0/camera/0/visuals", &visuals, nil); err != nil {
				log.Fatal(err)
			}
			zoomChanged = false
		}
	}
}
