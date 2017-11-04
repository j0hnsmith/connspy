// httpspy provides a http client that outputs raw http to stdout. Also makes the underlying net.Client implementation
// available.
package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// NewClient returns a http.Client that will output all http data to stdout. The client has various default timeouts,
// call with nil values to use them, otherwise pass arguments to customise.
func NewClient(dialer *net.Dialer, transport *http.Transport) *http.Client {
	if dialer == nil {
		dialer = &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
	}

	if transport == nil {
		transport = &http.Transport{
			DisableCompression:    true, // humans can't read a compressed response
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
	}

	dial := func(network, address string) (net.Conn, error) {
		conn, err := dialer.Dial(network, address)
		if err != nil {
			return nil, err
		}

		fmt.Fprint(os.Stderr, fmt.Sprintf("\n%s\n\n", strings.Repeat("-", 80)))
		return &SpyConnection{conn}, nil // return a wrapped net.Conn
	}

	dialTLS := func(network, address string) (net.Conn, error) {
		plainConn, err := dialer.Dial(network, address)
		if err != nil {
			return nil, err
		}

		//Initiate TLS and check remote host name against certificate.
		cfg := new(tls.Config)

		// add https:// to satisfy url.Parse(), we won't use it
		u, err := url.Parse(fmt.Sprintf("https://%s", address))
		if err != nil {
			return nil, err
		}

		serverName := u.Host[:strings.LastIndex(u.Host, ":")]
		cfg.ServerName = serverName

		tlsConn := tls.Client(plainConn, cfg)

		errc := make(chan error, 2)
		timer := time.AfterFunc(time.Second, func() {
			errc <- errors.New("TLS handshake timeout")
		})
		go func() {
			err := tlsConn.Handshake()
			timer.Stop()
			errc <- err
		}()
		if err := <-errc; err != nil {
			plainConn.Close()
			return nil, err
		}
		if !cfg.InsecureSkipVerify {
			if err := tlsConn.VerifyHostname(cfg.ServerName); err != nil {
				plainConn.Close()
				return nil, err
			}
		}

		fmt.Fprint(os.Stderr, fmt.Sprintf("\n%s\n\n", strings.Repeat("-", 80)))
		return &SpyConnection{tlsConn}, nil // return a wrapped net.Conn
	}

	transport.Dial = dial
	transport.DialTLS = dialTLS

	timeoutClient := &http.Client{
		Transport: transport,
	}

	return timeoutClient
}

// SpyConnection wraps a net.Conn, all reads and writes are output to stderr
type SpyConnection struct {
	net.Conn
}

// Read writes all data read from the underlying connection to stderr
func (sc *SpyConnection) Read(b []byte) (int, error) {
	tr := io.TeeReader(sc.Conn, os.Stderr)
	br, err := tr.Read(b)
	return br, err
}

// Write writes all data written to the underlying connection to stderr
func (sc *SpyConnection) Write(b []byte) (int, error) {
	mw := io.MultiWriter(sc.Conn, os.Stderr)
	bw, err := mw.Write(b)
	return bw, err
}
