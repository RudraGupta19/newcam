package lt

// GET /
type Agent struct {
	Version  string `json:"version"`
	Revision string `json:"revision"`
	Time     string `json:"time"`
}

// GET /:board
type Board struct {
	Model  string `json:"model"`
	SN     uint32 `json:"sn"`
	CPU    uint32 `json:"cpu"`
	FPGA   uint32 `json:"fpga"`
	Bridge uint32 `json:"bridge"`
}

// GET /:board/buttons/:id
// GET /:board/camera/:camera/buttons/:id
type Button struct {
	Description  string `json:"description"`
	Pressed      bool   `json:"pressed"`
	PressedCount int    `json:"pressedCount"`
	Timestamp    int64  `json:"timestamp"`
}

// GET /:board/buttons
// GET /:board/camera/:camera/buttons
type Buttons struct {
	Buttons []Button `json:"buttons"`
}

//
type AudioSignal struct {
	Description string `json:"description"`
	Format      string `json:"format"`
	Signal      string `json:"signal"`
	Channels    int    `json:"channels"`
	Samplerate  int    `json:"samplerate"`
	Depth       int    `json:"depth"`
}

//
type VideoSignal struct {
	Description string  `json:"description"`
	Format      string  `json:"format"`
	Signal      string  `json:"signal"`
	Size        [2]int  `json:"size"`
	Framerate   float64 `json:"framerate"`
	Interlaced  bool    `json:"interlaced"`
}

// GET /:board/camera/:id
type Camera struct {
	Model string      `json:"model"`
	SN    uint32      `json:"sn"`
	CPU   uint32      `json:"cpu"`
	FPGA  uint32      `json:"fpga"`
	Audio AudioSignal `json:"audio"`
	Video VideoSignal `json:"video"`
}

// GET,POST /:board/camera/:id/exposure
type CameraExposure struct {
	Framerate float64 `json:"framerate"`
	// Auto/Manual
	IsAuto bool `json:"isAuto"`
	// Manual parameters
	Shutter      float64 `json:"shutter"`
	Gain         float64 `json:"gain"`
	Binning      float64 `json:"binning"`
	LowLightGain float64 `json:"lowLightGain"`
	// Auto parameters
	Level         float64 `json:"level"`
	Speed         float64 `json:"speed"`
	MaxSaturation float64 `json:"maxSaturation"`
	// Limits
	ShutterLimits      [2]float64 `json:"shutterLimits"`
	GainLimits         [2]float64 `json:"gainLimits"`
	BinningLimits      [2]float64 `json:"binningLimits"`
	LowLightGainLimits [2]float64 `json:"lowLightGainLimits"`
	// Window
	Window [4]int `json:"window"`
}

// GET,POST /:board/camera/:id/white
type CameraWhite struct {
	// White balance
	Balance [3]float64 `json:"balance"`
	// Temperature
	Temperature int `json:"temperature"`
}

// GET,POST /:board/camera/:id/colors
type CameraColors struct {
	// Gamma
	Gamma float64 `json:"gamma"`
	// HSBC
	Hue        int `json:"hue"`
	Saturation int `json:"saturation"`
	Brightness int `json:"brightness"`
	Contrast   int `json:"contrast"`
	// Color gain
	ColorGain [3]float64 `json:"colorGain"`
}

// GET,POST /:board/camera/:id/visuals
type CameraVisuals struct {
	// View
	Flip string  `json:"flip"`
	Zoom float64 `json:"zoom"`
	// Details sharpness
	Sharpness      float64 `json:"sharpness"`
	SharpnessFloor int     `json:"sharpnessFloor"`
	// Noise reduction
	Anisotropic int `json:"anisotropic"`
	Bilateral   int `json:"bilateral"`
	// ShadowLighting
	ShadowLightingGain float64 `json:"shadowLightingGain"`
}

// GET /canvas/:id
type Input struct {
	Audio AudioSignal `json:"audio"`
	Video VideoSignal `json:"video"`
}

// GET,POST /:board/:output/:id
type Output struct {
	Source      string `json:"source"`
	Overlay     string `json:"overlay"`
	OverlayMode string `json:"overlayMode"`
	Osd         string `json:"osd"`
	Format      string `json:"format"`
	Link        string `json:"link"`

	Audio AudioSignal `json:"audio"`
	Video VideoSignal `json:"video"`
}

//
// Create DataWorker (POST)
//

type AudioDataWorker struct {
	Media string `json:"media"` // "audio/..."

	// Format
	Channels   int `json:"channels"`
	Samplerate int `json:"samplerate"`
	Depth      int `json:"depth"`
}

type ImageDataWorker struct {
	Media string `json:"media"` // "image/..."

	// Format
	Size [2]int `json:"size"`
}

type VideoDataWorker struct {
	Media string `json:"media"` // "video/..."

	// Format
	Size      [2]int  `json:"size"`
	Framerate float64 `json:"framerate"`
}

//
// Create FileWorker (POST)
//

type AudioFileWorker struct {
	Media string `json:"media"` // "audio/..."

	// File
	Location      string `json:"location"`
	Duration      int64  `json:"duration"`
	SplitSize     int    `json:"splitSize"`
	SplitDuration int64  `json:"splitDuration"`

	// Format
	Channels   int `json:"channels"`
	Samplerate int `json:"samplerate"`
	Depth      int `json:"depth"`
}

type ImageFileWorker struct {
	Media string `json:"media"` // "image/..."

	// File
	Location string `json:"location"`

	// Format
	Size [2]int `json:"size"`
}

type VideoEncoderExtra struct {
	// Encoder
	HW string `json:"hw"`

	//
	Bitrate int `json:"bitrate"`
	Quality int `json:"quality"`
	GOP     int `json:"gop"`

	// Codec
	Codec string `json:"codec"`

	// Preset
	Preset string `json:"preset"`
}

type VideoFileWorker struct {
	Media string `json:"media"` // "video/..."

	// File
	Location      string `json:"location"`
	Duration      int64  `json:"duration"`
	SplitSize     int    `json:"splitSize"`
	SplitDuration int64  `json:"splitDuration"`

	// Format
	Size      [2]int  `json:"size"`
	Framerate float64 `json:"framerate"`

	// Extra parameters
	Extra VideoEncoderExtra `json:"extra"`
}

// Packet metadata
type AudioMetadata struct {
	Channels   int `json:"channels"`
	Samplerate int `json:"samplerate"`
	Depth      int `json:"depth"`
	Samples    int `json:"samples"`
}

type ImageMetadata struct {
	Size [2]int `json:"size"`
}

type VideoMetadata struct {
	Size       [2]int  `json:"size"`
	Framerate  float64 `json:"framerate"`
	Interlaced bool    `json:"interlaced"`
	Keyframe   bool    `json:"keyframe"`
}

//
// Canvas
//

// POST /canvas/:id/ops
type CanvasOps struct {
	Ops []any `json:"ops"`
}

// POST /canvas/:id/init
type CanvasInit struct {
	Op string `json:"op"`

	Source string `json:"source"`

	Color     [4]int  `json:"color"`
	Size      [2]int  `json:"size"`
	Framerate float64 `json:"framerate"`
}

// POST /canvas/:id/clear
type CanvasClear struct {
	Op string `json:"op"`

	Color     [4]int `json:"color"`
	Position  [2]int `json:"position"`
	Size      [2]int `json:"size"`
	Thickness int    `json:"thickness"`
	Source    string `json:"source"`
}

// POST /canvas/:id/text
type CanvasText struct {
	Op string `json:"op"`

	Text  string `json:"text"`
	Align string `json:"align"`

	// Shape
	Font       string `json:"font"`
	FontSize   int    `json:"fontSize"`
	Italic     bool   `json:"italic"`
	Bold       bool   `json:"bold"`
	Color      [4]int `json:"color"`
	Background [4]int `json:"background"`

	// Container
	Angle    float64    `json:"angle"`
	Position [2]int     `json:"position"`
	Size     [2]int     `json:"size"`
	Anchor   [2]float64 `json:"anchor"`
}

// POST /canvas/:id/line
type CanvasLine struct {
	Op string `json:"op"`

	// Shape
	Width   int    `json:"width"`
	Color   [4]int `json:"color"`
	Pattern []int  `json:"pattern"`

	// Container
	Angle    float64    `json:"angle"`
	Position [2]int     `json:"position"`
	Size     [2]int     `json:"size"`
	Anchor   [2]float64 `json:"anchor"`
}

// POST /canvas/:id/ellipse
type CanvasEllipse struct {
	Op string `json:"op"`

	// Shape
	Width   int    `json:"width"`
	Color   [4]int `json:"color"`
	Fill    [4]int `json:"fill"`
	Pattern []int  `json:"pattern"`

	// Container
	Angle    float64    `json:"angle"`
	Position [2]int     `json:"position"`
	Size     [2]int     `json:"size"`
	Anchor   [2]float64 `json:"anchor"`
}

// POST /canvas/:id/rectangle
type CanvasRectangle struct {
	Op string `json:"op"`

	// Shape
	Width   int    `json:"width"`
	Color   [4]int `json:"color"`
	Fill    [4]int `json:"fill"`
	Pattern []int  `json:"pattern"`
	Rounded int    `json:"rounded"`

	// Container
	Angle    float64    `json:"angle"`
	Position [2]int     `json:"position"`
	Size     [2]int     `json:"size"`
	Anchor   [2]float64 `json:"anchor"`
}

// POST /canvas/:id/image
type CanvasImage struct {
	Op string `json:"op"`

	// Source
	Source string `json:"source"`

	// Data
	Format string `json:"format"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Data   []byte `json:"data"`

	// Container
	Angle    float64    `json:"angle"`
	Position [2]int     `json:"position"`
	Size     [2]int     `json:"size"`
	Anchor   [2]float64 `json:"anchor"`
}

// POST /canvas/:id/video
type CanvasVideo struct {
	Op string `json:"op"`

	// Source
	Source string `json:"source"`

	// Container
	Position [2]int     `json:"position"`
	Size     [2]int     `json:"size"`
	Anchor   [2]float64 `json:"anchor"`
}
