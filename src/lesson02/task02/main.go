// Command stats prints the line and word counts of a file as JSON.
//
//	stats <file>   -> {"lines":3,"words":7}, exit code 0
//	stats          -> usage on stderr, exit code 2
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Stats is the JSON document the program prints.
type Stats struct {
	Lines int `json:"lines"`
	Words int `json:"words"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: stats <file>")
		os.Exit(2)
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var s Stats
	//@@begin
	//@@ph // TODO: count the lines (strings.Split or strings.Count) and the words (strings.Fields) of data.
	//@@ph _, _ = data, strings.Fields
	if text := string(data); text != "" {
		s.Lines = strings.Count(text, "\n")
		if !strings.HasSuffix(text, "\n") {
			s.Lines++ // a last line without a newline still counts
		}
		s.Words = len(strings.Fields(text))
	}
	//@@end
	if err := json.NewEncoder(os.Stdout).Encode(s); err != nil {
		os.Exit(1)
	}
}
