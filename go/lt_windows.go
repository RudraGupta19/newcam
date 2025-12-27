package lt

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
	"unsafe"

	"github.com/tarm/serial"
	"golang.org/x/sys/windows"
)

//
// Conn
//

func dial(addr string) (net.Conn, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	// UART
	case "uart":
		port := u.Hostname()
		baurate := 115200
		if u.Port() != "" {
			if baurate, err = strconv.Atoi(u.Port()); err != nil {
				return nil, fmt.Errorf("invalid baudrate: %s %w", u.Port(), err)
			}
		}
		uart, err := serial.OpenPort(&serial.Config{Name: "COM" + port, Baud: baurate})
		if err != nil {
			return nil, fmt.Errorf("unknown uart port: %s %w", u.Scheme, err)
		}
		return &UARTConn{uart}, nil

	// TCP
	case "tcp":
		return net.Dial("tcp", u.Host)

	// Unix
	default:
		switch u.Host {
		// Windows: Unix Domain Socket
		case "unix":
			return net.Dial("unix", path.Join(os.TempDir(), u.Scheme+".sock"))
		// Windows: Named Pipe (default)
		case "pipe", "":
			return pipeDial("\\\\.\\pipe\\" + u.Scheme)
		default:
			return nil, fmt.Errorf("unknown port: %s", u.Port())
		}
	}
}

type UARTConn struct {
	*serial.Port
}

//func (uart *UARTConn) Read(b []byte) (n int, err error)
//func (uart *UARTConn) Write(b []byte) (n int, err error)
//func (uart *UARTConn) Close() error

func (uart *UARTConn) LocalAddr() net.Addr {
	return nil
}

func (uart *UARTConn) RemoteAddr() net.Addr {
	return nil
}

func (uart *UARTConn) SetDeadline(t time.Time) error {
	return nil
}

func (uart *UARTConn) SetReadDeadline(t time.Time) error {
	return nil

}

func (uart *UARTConn) SetWriteDeadline(t time.Time) error {
	return nil
}

//
// Shared buffer
//

type buffer struct {
	Data   []byte
	ref    int
	handle windows.Handle
	ptr    uintptr
}

func mapBuffer(handle string, size int) (*buffer, error) {
	b := &buffer{}
	key, err := windows.UTF16PtrFromString(handle)
	if err != nil {
		return nil, err
	}
	// Open system handle
	b.handle, err = windows.CreateFileMapping(
		windows.InvalidHandle,
		nil,
		windows.PAGE_READONLY,
		0,
		uint32(size),
		key)
	if err != nil && err != windows.ERROR_ALREADY_EXISTS {
		return nil, os.NewSyscallError("CreateFileMapping", err)
	}
	// Map buffer
	b.ptr, err = windows.MapViewOfFile(
		b.handle,
		windows.FILE_MAP_READ,
		0,
		0,
		uintptr(size))
	if err != nil {
		windows.CloseHandle(b.handle)
		return nil, os.NewSyscallError("MapViewOfFile", err)
	}
	b.Data = unsafe.Slice((*byte)(unsafe.Pointer(b.ptr)), size)
	// Done
	return b, nil
}

func (b *buffer) Close() error {
	if b.ptr != uintptr(0) {
		windows.UnmapViewOfFile(b.ptr)
		b.ptr = uintptr(0)
	}
	if b.handle != windows.InvalidHandle {
		windows.CloseHandle(b.handle)
		b.handle = windows.InvalidHandle
	}
	return nil
}

//
// Name pipe
//

type pipeConn struct {
	net.Conn
	handle windows.Handle
}

func pipeDial(addr string) (net.Conn, error) {
	name, err := windows.UTF16PtrFromString(addr)
	if err != nil {
		return nil, err
	}
	// Open named pipe
	c := &pipeConn{}
	c.handle, err = windows.CreateFile(
		name,
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		0,                     // mode
		nil,                   // default security attributes
		windows.OPEN_EXISTING, // opens existing pipe
		0,                     // default attributes
		0,                     // no template file
	)
	if err != nil {
		log.Println("CreateFile", err)
		return nil, ErrClosed
	}
	// Done
	return c, nil
}

func (c *pipeConn) Read(p []byte) (int, error) {
	var n uint32
	err := windows.ReadFile(c.handle, p, &n, nil)
	return int(n), err
}

func (c *pipeConn) Write(p []byte) (int, error) {
	var n uint32
	err := windows.WriteFile(c.handle, p, &n, nil)
	return int(n), err
}

func (c *pipeConn) Close() error {
	return windows.CloseHandle(c.handle)
}
