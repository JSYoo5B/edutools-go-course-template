package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	task "lesson02_task04"
)

// TestHello is the exercise: the server is started for you; send the request and read the body.
// The checks at the end grade your code. Do not change them.
func TestHello(t *testing.T) {
	srv := httptest.NewServer(task.Handler())
	defer srv.Close()

	//@@begin
	//@@ph // TODO: GET srv.URL + "/hello?name=Go" with http.Get, close the body when done,
	//@@ph // and read it with io.ReadAll.
	//@@ph var resp *http.Response
	//@@ph var body []byte
	//@@ph _ = io.ReadAll
	resp, err := http.Get(srv.URL + "/hello?name=Go")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading the body failed: %v", err)
	}
	//@@end

	// Given checks.
	if resp == nil {
		t.Fatal("resp is nil: send the request first")
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d, want 200", resp.StatusCode)
	}
	if string(body) != "Hello, Go!" {
		t.Fatalf("body %q, want %q", body, "Hello, Go!")
	}
	// End of given checks.
}
