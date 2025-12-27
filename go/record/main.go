package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/term"

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
	err = client.Post(sourceURL+"/file", lt.VideoFileWorker{Media: "video/mp4", Location: wd}, nil)
	if !errors.Is(err, lt.ErrRedirect) {
		log.Fatal("worker creation failed:", err)
	}
	workerURL := lt.RedirectLocation(err)

	// Allow terminal keyboard events
	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), state)

	// Capture keyboard events
	events := make(chan byte, 1)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			c, err := reader.ReadByte()
			if err != nil {
				log.Fatal(err)
			}
			// Ctrl+C
			if c == 0x03 {
				fmt.Println("")
				os.Exit(0)
			}
			// Forward event
			events <- c
		}
	}()

	var paused bool
	for {
		select {
		case c := <-events:
			switch c {
			// Start/Pause (space key)
			case ' ':
				if paused {
					// Start
					if err := client.Post(workerURL+"/start", nil, nil); err != nil {
						log.Fatal(err)
					}
					paused = false
				} else {
					// Pause
					if err := client.Post(workerURL+"/pause", nil, nil); err != nil {
						log.Fatal(err)
					}
					paused = true
				}

			// Stop (any other key)
			default:
				fmt.Println("stop")
				if err := client.Post(workerURL+"/stop", nil, nil); err != nil {
					log.Fatal(err)
				}
			}

		default:
			// Fetch worker response
			var worker lt.Worker
			if err := client.Get(workerURL, &worker); err != nil {
				log.Fatal(err)
			}

			// Loop over written packets
			for _, packet := range worker.Packets {
				switch strings.Split(packet.Media, "/")[0] {
				// Parse audio data
				case "audio":
					var meta lt.AudioMetadata
					if err := json.Unmarshal(packet.Meta, &meta); err != nil {
						log.Fatal("worker packet metadata:", err)
					}
					fmt.Printf("%s audio #%d  chan %d rate %d - %d bytes\n", worker.Name, packet.Track, meta.Channels, meta.Samplerate, worker.Length)

				// Parse video data
				case "video":
					var meta lt.VideoMetadata
					if err := json.Unmarshal(packet.Meta, &meta); err != nil {
						log.Fatal("worker packet metadata:", err)
					}
					fmt.Printf("\r%s %s %dx%d %d bytes", sourceURL, worker.Name, meta.Size[0], meta.Size[1], worker.Length)

				default:
					log.Fatal("unknown type: ", packet.Media)
				}

				// Release packet
				if err := packet.Close(); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
