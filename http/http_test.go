package http

import (
	"io/ioutil"
)

func ExampleNewClient() {
	client := connspy.NewClient(nil, nil)

	resp, _ := client.Get("http://example.com/")
	// ensure all of the body is read
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	resp, _ = client.Get("https://example.com/")
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// Output:
}
