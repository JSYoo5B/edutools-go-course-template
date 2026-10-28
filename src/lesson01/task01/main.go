// Example: count word frequencies with a map and print them with sorted keys. Run: go run .
package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	text := "Go is fun and go is fast"

	// With a map made by make, m[k]++ starts from 0 when the key is missing.
	counts := make(map[string]int)
	for _, w := range strings.Fields(strings.ToLower(text)) {
		counts[w]++
	}

	// The iteration order of a map is random, so sort the keys before printing.
	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%-5s %d\n", k, counts[k])
	}

	// comma-ok: tell "missing" from "present with value 0".
	if _, ok := counts["rust"]; !ok {
		fmt.Println("rust is not present")
	}
}
