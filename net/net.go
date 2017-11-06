// net provides tools for eavesdropping on a net.Conn.
package net

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

// SpyConnection wraps a net.Conn, all reads and writes are output to stderr, via WrapConnection().
type SpyConnection struct {
	net.Conn
	DebugWriter io.Writer
}

// Read writes all data read from the underlying connection to stderr.
func (sc *SpyConnection) Read(b []byte) (int, error) {
	tr := io.TeeReader(sc.Conn, sc.DebugWriter)
	br, err := tr.Read(b)
	return br, err
}

// Write writes all data written to the underlying connection to stderr.
func (sc *SpyConnection) Write(b []byte) (int, error) {
	mw := io.MultiWriter(sc.Conn, sc.DebugWriter)
	bw, err := mw.Write(b)
	return bw, err
}

// WrapConnection wraps an existing connection, all data read/written is written to w (os.Stderr if w == nil).
func WrapConnection(c net.Conn, w io.Writer) *SpyConnection {
	if w == nil {
		w = os.Stderr
	}
	return &SpyConnection{c, w}
}

// SpyListener wraps a net.Listener, all reads and writes to/from accepted conns are output to stderr
type SpyListener struct {
	net.Listener
	ConnWriter func() (io.WriteCloser, error)
}

// Accept on the underlying listener and wrap conn with SpyConnection. Calls ConnWriter for each new connection to get
// somewhere to write data.
func (sl *SpyListener) Accept() (net.Conn, error) {
	conn, err := sl.Listener.Accept()
	if err != nil {
		return nil, err
	}
	w, err := sl.ConnWriter()
	if err != nil {
		return nil, err
	}
	return WrapConnection(conn, w), nil
}

// WrapListener wraps a net.Listener, , all reads and writes to/from accepted conns are output to stderr
func WrapListener(l net.Listener, f func() (io.WriteCloser, error)) *SpyListener {
	return &SpyListener{l, f}
}

// DebugFile is an os.File based io.WriteCloser for use with WrapListener (via NewDebugFileWriter).
type DebugFile struct {
	f *os.File
}

// Write bytes to file.
func (df DebugFile) Write(b []byte) (int, error) {
	return df.f.Write(b)
}

// Close file.
func (df DebugFile) Close() error {
	return df.f.Close()
}

// NewDebugFileWriter opens a file /tmp/{timestamp}.txt, used to write all output for a single connection.
func NewDebugFileWriter() (io.WriteCloser, error) {
	var esLayout = "20060102T150405.000Z"
	tStr := time.Now().UTC().Format(esLayout)

	f, err := os.Create(fmt.Sprintf("/tmp/%s.txt", tStr))
	if err != nil {
		return DebugFile{}, err
	}
	return DebugFile{f}, nil
}

// NewStderrWriter is used to write all output to stderr. If the server handles multiple connetions at the same time
// data will appear all jumbled up, try NewDebugFileWriter instead.
func NewStderrWriter() (io.WriteCloser, error) {
	return os.Stderr, nil
}
