package lt

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
	"strconv"
	"syscall"
	"time"

	"github.com/tarm/serial"
	"golang.org/x/sys/unix"
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
		baudrate := 115200
		if u.Port() != "" {
			if baudrate, err = strconv.Atoi(u.Port()); err != nil {
				return nil, fmt.Errorf("invalid baudrate: %s %w", u.Port(), err)
			}
		}
		uart, err := serial.OpenPort(&serial.Config{Name: "TTY" + port, Baud: baudrate})
		if err != nil {
			return nil, fmt.Errorf("uart error: %s %s %d %w", u.Scheme, port, baudrate, err)
		}
		return &UARTConn{uart}, nil

	// TCP
	case "tcp":
		return net.Dial("tcp", u.Host)

	// Unix
	default:
		return net.Dial("unix", path.Join(os.TempDir(), u.Scheme+".sock"))
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
	name string
	fd   int
	Data []byte
	ref  int
}

const shmPath = "/dev/shm/"

func mapBuffer(handle string, size int) (*buffer, error) {
	b := &buffer{}

	// Open the shared memory
	perm := uint32(unix.S_IRUSR | unix.S_IWUSR | unix.S_IRGRP | unix.S_IWGRP)
	fd, err := syscall.Open(shmPath+handle, unix.O_RDWR, perm)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: handle, Err: err}
	}

	// Setup the buffer now so we can clean it up easy
	b.name = handle
	b.fd = fd
	b.Data = nil
	b.ref = 0

	// Resize to the required size
	err = syscall.Ftruncate(b.fd, int64(size))
	if err != nil {
		b.Close()
		return nil, &os.PathError{Op: "ftruncate", Path: handle, Err: err}
	}

	// Map the file in memory
	b.Data, err = syscall.Mmap(b.fd, 0, size, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		b.Close()
		return nil, &os.PathError{Op: "mmap", Path: handle, Err: err}
	}

	// Done
	return b, nil
}

func (b *buffer) Close() error {
	if b == nil {
		return nil
	}
	// Unmap memory
	if b.Data != nil {
		syscall.Munmap(b.Data)
	}
	// Close file descriptor
	if b.fd >= 0 {
		syscall.Close(b.fd)
	}

	return nil
}
