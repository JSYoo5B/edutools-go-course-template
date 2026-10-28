package main

import "fmt"

func main() {
	counts := map[string]int{"go": 2}
	n := counts["rust"]
	fmt.Println(n)
}
