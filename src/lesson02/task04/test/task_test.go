package test

import (
	"os"
	"strings"
	"testing"
)

// givenChecks must stay in hello_test.go: they are what grades TestHello.
const givenChecks = `	// Given checks.
	if resp == nil {
		t.Fatal("resp is nil: send the request first")
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d, want 200", resp.StatusCode)
	}
	if string(body) != "Hello, Go!" {
		t.Fatalf("body %q, want %q", body, "Hello, Go!")
	}
	// End of given checks.`

func TestGivenChecksUnchanged(t *testing.T) {
	src, err := os.ReadFile("hello_test.go")
	if err != nil {
		t.Fatalf("hello_test.go is missing: %v", err)
	}
	if !strings.Contains(strings.ReplaceAll(string(src), "\r\n", "\n"), givenChecks) {
		t.Fatal("the given checks at the end of TestHello were changed; restore them (Reset Task restores the file)")
	}
}
