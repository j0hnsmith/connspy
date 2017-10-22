package main

import (
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

type spyConnection struct {
	net.Conn
}

func (sc *spyConnection) Read(b []byte) (int, error) {
	tr := io.TeeReader(sc.Conn, os.Stderr)
	br, err := tr.Read(b)
	return br, err
}

func (sc *spyConnection) Write(b []byte) (int, error) {
	mw := io.MultiWriter(sc.Conn, os.Stderr)
	bw, err := mw.Write(b)
	//bw, err := sc.Conn.Write(b)
	return bw, err
}

func main() {
	dialer := (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	})

	dial := func(network, address string) (net.Conn, error) {
		conn, err := dialer.Dial(network, address)
		if err != nil {
			return nil, err
		}
		return &spyConnection{
			conn,
		}, nil
	}

	timeoutClient := &http.Client{
		Transport: &http.Transport{
			Dial:                  dial,
			DisableCompression:    true,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	timeoutClient.Get("http://example.com/")
}
