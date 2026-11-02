// Package hello serves a greeting over HTTP.
package hello

import (
	"fmt"
	"net/http"
)

// Handler answers GET /hello?name=X with "Hello, X!".
func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Query().Get("name"))
	})
	return mux
}
