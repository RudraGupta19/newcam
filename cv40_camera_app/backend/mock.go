package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	lt "lt/client/go"
)

// MockTransport implements a fake lt.Client for local development
type MockTransport struct {
	recording bool
	paused    bool
	settings  MockSettings
}

type MockSettings struct {
	Visuals  lt.CameraVisuals
	Colors   lt.CameraColors
	White    lt.CameraWhite
	Exposure lt.CameraExposure
}

var mockTransport = &MockTransport{
	settings: MockSettings{
		Visuals: lt.CameraVisuals{
			Zoom:      1.0,
			Sharpness: 0.5,
		},
		Colors: lt.CameraColors{
			Brightness: 0,
			Contrast:   0,
			Saturation: 0,
			Hue:        0,
		},
		White: lt.CameraWhite{
			Temperature: 6500,
		},
		Exposure: lt.CameraExposure{
			LowLightGain: 0.0,
		},
	},
}

// MockClient wraps the real lt.Client but intercepts calls when in mock mode
type MockClient struct {
	real *lt.Client
	mock bool
}

func NewMockClient() *MockClient {
	return &MockClient{
		real: &lt.Client{},
		mock: os.Getenv("MOCK_CAMERA") == "1",
	}
}

func (c *MockClient) Get(url string, response any) error {
	if !c.mock {
		return c.real.Get(url, response)
	}
	return c.mockGet(url, response)
}

func (c *MockClient) Post(url string, body, response any) error {
	if !c.mock {
		return c.real.Post(url, body, response)
	}
	return c.mockPost(url, body, response)
}

func (c *MockClient) Delete(url string) error {
	if !c.mock {
		return c.real.Delete(url)
	}
	return nil // Mock delete always succeeds
}

func (c *MockClient) Close() error {
	if !c.mock {
		return c.real.Close()
	}
	return nil
}

func (c *MockClient) mockGet(url string, response any) error {
	log.Printf("MOCK GET %s", url)
	
	switch {
	case url == "cv40:/0/camera/0":
		if cam, ok := response.(*lt.Camera); ok {
			*cam = lt.Camera{
				Model: "Mock Camera",
				Video: lt.VideoSignal{
					Signal: "locked",
					Size:   [2]int{1920, 1080},
				},
			}
		}
	case url == "cv40:/0/camera/0/visuals":
		if v, ok := response.(*lt.CameraVisuals); ok {
			*v = mockTransport.settings.Visuals
		}
	case url == "cv40:/0/camera/0/colors":
		if v, ok := response.(*lt.CameraColors); ok {
			*v = mockTransport.settings.Colors
		}
	case url == "cv40:/0/camera/0/white":
		if v, ok := response.(*lt.CameraWhite); ok {
			*v = mockTransport.settings.White
		}
	case url == "cv40:/0/camera/0/exposure":
		if v, ok := response.(*lt.CameraExposure); ok {
			*v = mockTransport.settings.Exposure
		}
	}
	return nil
}

func (c *MockClient) mockPost(url string, body, response any) error {
	log.Printf("MOCK POST %s", url)
	
	switch {
	case url == "cv40:/0/camera/0/visuals":
		if v, ok := body.(*lt.CameraVisuals); ok {
			mockTransport.settings.Visuals = *v
		}
	case url == "cv40:/0/camera/0/colors":
		if v, ok := body.(*lt.CameraColors); ok {
			mockTransport.settings.Colors = *v
		}
	case url == "cv40:/0/camera/0/white":
		if v, ok := body.(*lt.CameraWhite); ok {
			mockTransport.settings.White = *v
		}
	case url == "cv40:/0/camera/0/exposure":
		if v, ok := body.(*lt.CameraExposure); ok {
			mockTransport.settings.Exposure = *v
		}
	case url == "cv40:/0/camera/0/file":
		// Mock file worker creation - return redirect error with fake worker URL
		if _, ok := body.(lt.ImageFileWorker); ok {
			return fmt.Errorf("%w: mock://worker/image/%d", lt.ErrRedirect, rand.Int())
		}
		if _, ok := body.(lt.VideoFileWorker); ok {
			return fmt.Errorf("%w: mock://worker/video/%d", lt.ErrRedirect, rand.Int())
		}
	default:
		// Mock worker control endpoints (start/pause/stop)
		if len(url) > 12 && url[:12] == "mock://worker" {
			// Simulate worker operations
			time.Sleep(100 * time.Millisecond) // Simulate processing time
		}
	}
	return nil
}

// Helper to create a mock-aware client
func createClient() *MockClient {
	return NewMockClient()
}
