package cv40

import (
    "errors"
    "strconv"
    "time"
    lt "lt/client/go"
    "cv40-camera-backend/internal/config"
)

type RealClient struct {
    cfg config.Config
    c   lt.Client
}

func NewRealClient(cfg config.Config) *RealClient {
    return &RealClient{cfg: cfg}
}

func (r *RealClient) Close() { r.c.Close() }

func (r *RealClient) Health() error {
    var cam lt.Camera
    r.c.Timeout = 1500 * time.Millisecond
    return r.c.Get("cv40:/0/camera/0", &cam)
}

func (r *RealClient) CreateVideoWorker(dest string, media string) (string, error) {
    err := r.c.Post("cv40:/0/camera/0/file", lt.VideoFileWorker{Media: media, Location: dest}, nil)
    if !errors.Is(err, lt.ErrRedirect) { return "", err }
    return lt.RedirectLocation(err), nil
}

func (r *RealClient) StartWorker(u string) error { return r.c.Post(u+"/start", nil, nil) }
func (r *RealClient) PauseWorker(u string) error { return r.c.Post(u+"/pause", nil, nil) }
func (r *RealClient) StopWorker(u string) error { return r.c.Post(u+"/stop", nil, nil) }

func (r *RealClient) CaptureStill(dest string) error {
    err := r.c.Post("cv40:/0/camera/0/file", lt.ImageFileWorker{Media: "image/jpeg", Location: dest}, nil)
    if !errors.Is(err, lt.ErrRedirect) { return err }
    return nil
}

func (r *RealClient) SetColors(v lt.CameraColors) error { return r.c.Post("cv40:/0/camera/0/colors", &v, nil) }
func (r *RealClient) SetVisuals(v lt.CameraVisuals) error { return r.c.Post("cv40:/0/camera/0/visuals", &v, nil) }
func (r *RealClient) SetWhite(v lt.CameraWhite) error { return r.c.Post("cv40:/0/camera/0/white", &v, nil) }
func (r *RealClient) SetExposure(v lt.CameraExposure) error { return r.c.Post("cv40:/0/camera/0/exposure", &v, nil) }

func (r *RealClient) GetColors() (lt.CameraColors, error) { var v lt.CameraColors; err := r.c.Get("cv40:/0/camera/0/colors", &v); return v, err }
func (r *RealClient) GetVisuals() (lt.CameraVisuals, error) { var v lt.CameraVisuals; err := r.c.Get("cv40:/0/camera/0/visuals", &v); return v, err }
func (r *RealClient) GetWhite() (lt.CameraWhite, error) { var v lt.CameraWhite; err := r.c.Get("cv40:/0/camera/0/white", &v); return v, err }
func (r *RealClient) GetExposure() (lt.CameraExposure, error) { var v lt.CameraExposure; err := r.c.Get("cv40:/0/camera/0/exposure", &v); return v, err }

func (r *RealClient) ConfigureOutputOverlay(output string, canvasID int) error {
    return r.c.Post("cv40:/0/"+output, lt.JSON{"overlay": "canvas/"+strconv.Itoa(canvasID)}, nil)
}

func (r *RealClient) CanvasInit(id int, size [2]int) error {
    return r.c.Post("cv40:/canvas/"+strconv.Itoa(id)+"/init", lt.CanvasInit{Size: size}, nil)
}

func (r *RealClient) CanvasText(id int, t lt.CanvasText) error {
    return r.c.Post("cv40:/canvas/"+strconv.Itoa(id)+"/text", t, nil)
}

func (r *RealClient) GetOutput(output string) (lt.Output, error) {
    var out lt.Output
    err := r.c.Get("cv40:/0/"+output, &out)
    return out, err
}

func (r *RealClient) GetWorker(u string) (lt.Worker, error) {
    var w lt.Worker
    err := r.c.Get(u, &w)
    return w, err
}
