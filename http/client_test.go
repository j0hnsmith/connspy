package http_test

import (
	"io/ioutil"

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
