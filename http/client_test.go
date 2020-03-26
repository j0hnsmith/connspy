package http_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/j0hnsmith/connspy/http"
)

func ExampleNewClient() {
	client := http.NewClient(nil, nil)

	resp, _ := client.Get("http://example.com/")
	// ensure all of the body is read
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	resp, _ = client.Get("https://example.com/")
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// Output:
}

func ExampleNewClientWriter_toStderr() {
	client := http.NewClientWriter(nil, nil, os.Stderr)

	resp, _ := client.Get("http://example.com/")
	// ensure all of the body is read
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	resp, _ = client.Get("https://example.com/")
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// Output:
}

func ExampleNewClientWriter_toCustomBuffer() {
	buf := new(bytes.Buffer)

	client := http.NewClientWriter(nil, nil, buf)

	resp, _ := client.Get("http://example.com/")
	// ensure all of the body is read
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	resp, _ = client.Get("https://example.com/")
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	httpLog := buf.String()

	s := strings.Split(httpLog, "\n")

	fmt.Println(s[4])

	// Output:
	// Host: example.com
}

func ExampleNewClientWriter_toCustomBufferWithRedaction() {
	buf := new(bytes.Buffer)

	client := http.NewClientWriter(nil, nil, buf)

	resp, _ := client.Get("http://example.com/")
	// ensure all of the body is read
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	resp, _ = client.Get("https://example.com/")
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	httpLog := buf.String()

	s := strings.Split(httpLog, "\n")

	for count, line := range s {
		rgx := regexp.MustCompile(`^(Host: )(.+)$`)
		line = rgx.ReplaceAllString(line, `$1[REDACTED]`)
		s[count] = line
	}

	fmt.Println(s[4])
	// Output:
	// Host: [REDACTED]
}
