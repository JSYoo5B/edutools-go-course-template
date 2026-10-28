// Package wordcount is an exercise on counting words with a map. Fill in the TODO part.
package wordcount

import "strings"

// WordCount returns the frequency of each word of text, lower-cased and split around white space.
func WordCount(text string) map[string]int {
	//@@begin
	//@@ph // TODO: count the words with strings.ToLower, strings.Fields and a map, and return it.
	//@@ph _ = strings.Fields
	//@@ph return nil
	counts := make(map[string]int)
	for _, w := range strings.Fields(strings.ToLower(text)) {
		counts[w]++
	}
	return counts
	//@@end
}
