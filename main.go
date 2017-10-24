package main

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
		return &spyConnection{conn}, nil
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

		return &spyConnection{tlsConn}, nil
	}

	t := &http.Transport{
		Dial:                  dial,
		DialTLS:               dialTLS,
		DisableCompression:    true,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	timeoutClient := &http.Client{
		Transport: t,
	}

	//timeoutClient.Get("http://example.com/")
	timeoutClient.Get("https://www.google.co.uk/")
	timeoutClient.Get("https://duckduckgo.com/")
}
