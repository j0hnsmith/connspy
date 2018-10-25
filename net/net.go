// net provides tools for eavesdropping on a net.Conn.
package net

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

// WrapConnection wraps an existing connection, all data read/written is written to w (os.Stderr if w == nil).
func WrapConnection(c net.Conn, output io.Writer) net.Conn {
	return &spyConnection{
		Conn:   c,
		Reader: io.TeeReader(c, output),
		Writer: io.MultiWriter(output, c),
	}
}

// spyConnection wraps a net.Conn, all reads and writes are output to stderr, via WrapConnection().
type spyConnection struct {
	net.Conn
	io.Reader
	io.Writer
}

// Read writes all data read from the underlying connection to sc.Writer.
func (sc *spyConnection) Read(b []byte) (int, error) {
	return sc.Reader.Read(b)
}

// Write writes all data written to the underlying connection to sc.Writer.
func (sc *spyConnection) Write(b []byte) (int, error) {
	return sc.Writer.Write(b)
}

// spyListener wraps a net.Listener, all reads and writes to/from accepted conns are output to stderr
type spyListener struct {
	net.Listener
	getConnWriter func() (io.WriteCloser, error)
}

// Accept on the underlying listener and wrap conn with spyConnection. Calls ConnWriter for each new connection to get
// somewhere to write data.
func (sl *spyListener) Accept() (net.Conn, error) {
	conn, err := sl.Listener.Accept()
	if err != nil {
		return nil, err
	}
	w, err := sl.getConnWriter()
	if err != nil {
		return nil, err
	}
	return WrapConnection(conn, w), nil
}

// WrapListener wraps a net.Listener, all reads and writes to/from accepted conns are output to a writer from f.
func WrapListener(l net.Listener, f func() (io.WriteCloser, error)) net.Listener {
	return &spyListener{l, f}
}

// NewDebugFileWriter opens a file /tmp/{timestamp}.txt, used to write all output for a single connection.
func NewDebugFileWriter() (io.WriteCloser, error) {
	var esLayout = "20060102T150405.000Z"
	tStr := time.Now().UTC().Format(esLayout)

	return os.Create(fmt.Sprintf("/tmp/%s.txt", tStr))
}
