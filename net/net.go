// net provides tools for eavesdropping on a net.Conn.
package net

import (
	"io"
	"net"
	"os"
)

// SpyConnection wraps a net.Conn, all reads and writes are output to stderr, via WrapConnection().
type SpyConnection struct {
	net.Conn
}

// Read writes all data read from the underlying connection to stderr.
func (sc *SpyConnection) Read(b []byte) (int, error) {
	tr := io.TeeReader(sc.Conn, os.Stderr)
	br, err := tr.Read(b)
	return br, err
}

// Write writes all data written to the underlying connection to stderr.
func (sc *SpyConnection) Write(b []byte) (int, error) {
	mw := io.MultiWriter(sc.Conn, os.Stderr)
	bw, err := mw.Write(b)
	return bw, err
}

// WrapConnection wraps an existing connection, all data read/written is output to stdout.
func WrapConnection(c net.Conn) *SpyConnection {
	return &SpyConnection{c}
}
