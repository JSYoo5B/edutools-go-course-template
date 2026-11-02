package test

import (
	"testing"

	task "lesson02_task05"
)

// TestClamp is yours to write. It must pass for the correct Clamp and fail for a broken one.
func TestClamp(t *testing.T) {
	//@@begin
	//@@ph // TODO: call task.Clamp with values below, inside and above a range and t.Errorf on a wrong result.
	//@@ph _ = task.Clamp
	cases := []struct{ v, want int }{
		{-5, 0},  // below
		{0, 0},   // on the lower bound
		{7, 7},   // inside
		{10, 10}, // on the upper bound
		{42, 10}, // above
	}
	for _, c := range cases {
		if got := task.Clamp(c.v, 0, 10); got != c.want {
			t.Errorf("Clamp(%d, 0, 10) = %d, want %d", c.v, got, c.want)
		}
	}
	//@@end
}
