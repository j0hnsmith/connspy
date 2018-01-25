# connspy

Tools for spying on connections, all output read/written to stderr 

### `http` package 

A `http.Client` suitable for debugging, writes all http data to stdout.

```go
client := connspy.NewClient(nil, nil)

resp, _ := client.Get("http://example.com/")
// ensure all of the body is read
ioutil.ReadAll(resp.Body)
resp.Body.Close()

resp, _ = client.Get("https://example.com/")
ioutil.ReadAll(resp.Body)
resp.Body.Close()
```

![http output to stderr](https://dl.dropboxusercontent.com/s/ved2xxrp3rbzome/Screen%20Shot%202017-11-02%20at%2022.48.35.png?dl=0)

### `net` package

Provides a `net.Conn` wrapper that writes all reads/writes to stderr.

## Docs

[![GoDoc](https://godoc.org/github.com/j0hnsmith/connspy?status.svg)](https://godoc.org/github.com/j0hnsmith/connspy) 

## Background info

[https://medium.com/@j0hnsmith/eavesdrop-on-a-golang-http-client-c4dc49af9d5e](https://medium.com/@j0hnsmith/eavesdrop-on-a-golang-http-client-c4dc49af9d5e)
