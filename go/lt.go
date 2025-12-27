package lt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"sync"
)

//
// Errors
//

var (
	EOF         = io.EOF
	ErrRedirect = errors.New("redirect")
	ErrUpdating = errors.New("updating")
	ErrClosed   = errors.New("use of closed connection")
)

func RedirectLocation(err error) string {
	location, _ := strings.CutPrefix(err.Error(), ErrRedirect.Error()+": ")
	return location
}

//
// Shared buffers local references
//

var sharedBuffers = newBuffers()

type buffers struct {
	m  map[string]*buffer // Handles
	mu sync.RWMutex
}

func newBuffers() *buffers {
	return &buffers{
		m: map[string]*buffer{},
	}
}

func (bs *buffers) Load(handle string, size int) (*buffer, bool) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	// Map shared buffer
	if _, ok := bs.m[handle]; !ok {
		b, err := mapBuffer(handle, size)
		if err != nil {
			return nil, false
		}
		bs.m[handle] = b
	}
	bs.m[handle].ref++
	// Done
	return bs.m[handle], true
}

func (bs *buffers) Delete(handle string) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	// Delete shared buffer
	if b, ok := bs.m[handle]; ok {
		b.ref -= 1
		if b.ref <= 0 {
			//b.Close()
			//delete(bs.m, handle)
		}
	}
}

//
// Packets
//

// Packet provided into the worker response
type Packet struct {
	Track     int
	Media     string
	Signal    string
	Timestamp int64

	Ref  string
	Data []byte
	Meta json.RawMessage

	// Internal
	handle       string
	roundTripper *roundTripper
}

// Handle base64 encoded data or shared memory data
func (p *Packet) UnmarshalJSON(j []byte) error {
	pkt := struct {
		Track     int    `json:"track"`
		Media     string `json:"media"`
		Signal    string `json:"signal"`
		Timestamp int64  `json:"timestamp"`

		// Encoded data (base64)
		Data []byte          `json:"data"`
		Meta json.RawMessage `json:"meta"`

		// Shared memory data
		Ref    string `json:"ref"`
		Handle string `json:"handle"`
		Ptr    int    `json:"ptr"`
		Len    int    `json:"len"`
		Cap    int    `json:"cap"`
	}{}
	if err := json.Unmarshal(j, &pkt); err != nil {
		return err
	}

	// Header
	p.Track = pkt.Track
	p.Media = pkt.Media
	p.Signal = pkt.Signal
	p.Timestamp = pkt.Timestamp
	p.Meta = pkt.Meta

	// Data
	if pkt.Ref == "" {
		// Base64 Encoded
		p.Data = pkt.Data
	} else {
		// Load from local handler
		buffer, ok := sharedBuffers.Load(pkt.Handle, pkt.Cap)
		if !ok {
			return errors.New("packet: shared memory data not found")
		}
		// Data
		p.handle = pkt.Handle
		p.Ref = pkt.Ref
		p.Data = buffer.Data[pkt.Ptr : pkt.Ptr+pkt.Len : pkt.Ptr+pkt.Len]
	}

	// Done
	return nil
}

// Release references
func (p *Packet) Close() error {
	// Local
	if p.handle != "" {
		sharedBuffers.Delete(p.handle)
	}
	// Remote
	if p.Ref != "" && p.roundTripper != nil {
		go func() {
			p.roundTripper.Call("DELETE", p.Ref, nil, nil)
			p.roundTripper.Close()
		}()
	}
	// Done
	return nil
}

//
// Workers
//

type Worker struct {
	Name     string `json:"name"`
	Location string `json:"location"`

	Start    int64  `json:"start"`
	Duration int64  `json:"duration"`
	Length   int    `json:"length"`
	Status   string `json:"status"`

	Packets []Packet `json:"packets"`
}

//
// RoundTripper
//

type roundTripper struct {
	scheme string
	conn   net.Conn
	decode func(any) error
	encode func(any) error
	ref    int
	mu     sync.Mutex
}

func (rt *roundTripper) open(addr string) error {
	conn, err := dial(addr)
	if err != nil {
		return err
	}
	rt.conn = conn
	rt.decode = json.NewDecoder(conn).Decode
	rt.encode = json.NewEncoder(conn).Encode
	return nil
}

func (rt *roundTripper) close() error {
	if rt.ref > 0 {
		rt.ref--
		return nil
	}
	rt.scheme = ""
	if rt.conn == nil {
		return nil
	}
	err := rt.conn.Close()
	rt.conn = nil
	return err
}

func (rt *roundTripper) Clone() *roundTripper {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	rt.ref++
	return rt
}

func (rt *roundTripper) Close() error {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	return rt.close()
}

func (rt *roundTripper) Call(method, location string, body, response any) error {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	// Validate url
	u, err := url.Parse(location)
	if err != nil {
		return err
	}
	if u.Scheme == "" {
		return errors.New("url scheme not found")
	}
	if u.Path == "" {
		return errors.New("url path not found")
	}

	// Validate connection
	if rt.scheme != "" && rt.scheme != u.Scheme {
		return fmt.Errorf("bad url scheme: %s != %s", rt.scheme, u.Scheme)
	}
	if rt.scheme == "" {
		if err := rt.open(location); err != nil {
			return err
		}
		rt.scheme = u.Scheme
	}

	// JSON Request
	request := struct {
		Method string `json:"method"`
		URL    string `json:"url"`
		Body   any    `json:"body"`
	}{
		Method: method,
		URL:    location,
		Body:   body,
	}
	if err := rt.encode(&request); err != nil {
		rt.close()
		return err
	}

	// JSON Response
	var responseMessage json.RawMessage
	if err := rt.decode(&responseMessage); err != nil {
		rt.close()
		return err
	}

	// Check for error
	responseError := struct {
		Location string `json:"location"`
		Error    string `json:"error"`
	}{}
	if err := json.Unmarshal(responseMessage, &responseError); err != nil {
		return err
	}
	switch responseError.Error {
	case "":
		// No error
	case EOF.Error():
		return EOF
	case ErrRedirect.Error():
		return fmt.Errorf("%w: %s", ErrRedirect, responseError.Location)
	case ErrUpdating.Error():
		return ErrUpdating
	case ErrClosed.Error():
		return ErrClosed
	default:
		return errors.New(responseError.Error)
	}

	// Parse response
	if response != nil {
		return json.Unmarshal(responseMessage, response)
	}

	// Done
	return nil
}

//
// Client
//

func Get(url string, response any) error {
	var client Client
	defer client.Close()
	return client.Get(url, response)
}

func Post(url string, body, response any) error {
	var client Client
	defer client.Close()
	return client.Post(url, body, response)
}

func Delete(url string) (err error) {
	var client Client
	defer client.Close()
	return client.Delete(url)
}

type Client struct {
	roundTripper roundTripper
}

func (c *Client) Get(url string, response any) error {
	if err := c.call("GET", url, nil, response); err != nil {
		return err
	}
	if worker, ok := response.(*Worker); ok {
		for i := range worker.Packets {
			if worker.Packets[i].Ref != "" {
				worker.Packets[i].roundTripper = c.roundTripper.Clone()
			}
		}
	}
	return nil
}

func (c *Client) Post(url string, body, response any) error {
	return c.call("POST", url, body, response)
}

func (c *Client) Delete(url string) error {
	return c.call("DELETE", url, nil, nil)
}

func (c *Client) Close() error {
	return c.roundTripper.Close()
}

func (c *Client) call(method, location string, body, response any) error {
	return c.roundTripper.Call(method, location, body, response)
}
