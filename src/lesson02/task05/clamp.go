// Package clamp limits a value to a range.
package clamp

// Clamp returns lo if v < lo, hi if v > hi, and v otherwise.
func Clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
