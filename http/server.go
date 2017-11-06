// Also adds a http server wrapper to output data for each connection to a specific io.Writer.
package http

import (
	"io"
	"net"
	"net/http"

	spynet "github.com/j0hnsmith/connspy/net"
)

type SpyServer struct {
	http.Server
	wc func() (io.WriteCloser, error)
}

// ListenAndServe is similar to the net.http Server equivalent, no tcp keep-alive though.
func (ss *SpyServer) ListenAndServe() error {
	addr := ss.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	spyLn := spynet.WrapListener(ln, ss.wc)
	return ss.Serve(spyLn)
}

// WrapServer turns a normal http server into one what writes all connection data to an io.Writer. The function (wf) is
// called each time a new connection is made, it should return a writer that will be used for all the data for that
// connection.
//
//  package main
//
//  import (
//      "fmt"
//      "net/http"
//
//      spyhttp "github.com/j0hnsmith/connspy/http"
//      spynet "github.com/j0hnsmith/connspy/net"
//  )
//
//  type SomeHandler struct{}
//
//  func (sh SomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//      fmt.Fprintf(w, "Hello World!\n")
//  }
//
//  func main() {
//      server := http.Server{
//          Addr:    ":8080",
//          Handler: &SomeHandler{},
//      }
//      spyServer := spyhttp.WrapServer(server, spynet.NewDebugFileWriter)
//      spyServer.ListenAndServe()
//  }
func WrapServer(s http.Server, wf func() (io.WriteCloser, error)) *SpyServer {
	return &SpyServer{s, wf}
}
