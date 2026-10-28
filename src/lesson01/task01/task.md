# Maps and counting

A map (`map[K]V`) looks up a value by key. Counting word frequencies in a text is the most common example.

## Operations and functions used

| Operation | Signature / syntax | What it does | Caveat |
| --- | --- | --- | --- |
| Create | `make(map[string]int)` | An empty map | You cannot write to a `nil` map (reading is fine) |
| Write | `m[k]++` | Starts from the zero value when the key is missing | The core of frequency counting |
| Read | `v, ok := m[k]` | The value and whether the key exists | Use `ok` to tell "missing" from "0" |
| `strings.Fields` | `func Fields(s string) []string` | Splits around white space | Runs of spaces count as one separator |
| `strings.ToLower` | `func ToLower(s string) string` | Lower-cases a string | For case-insensitive comparison |
| `sort.Strings` | `func Strings(x []string)` | Sorts a slice | Map iteration order is **random**, so sort the keys before printing |

## Caveats

- The iteration order of a map differs on every run. Sort the keys when you print or compare results.
- Do not write to a map from several goroutines at once (you need a `sync.Mutex`).

Run the example in `main.go` on the left.
