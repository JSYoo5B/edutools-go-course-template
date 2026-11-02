package test

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// run builds the learner's program (the task root, "..") and runs it. It builds instead of using
// go run, because go run does not pass the program's exit code through.
func run(t *testing.T, args ...string) ([]byte, int) {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "stats")
	if out, err := exec.CommandContext(t.Context(), "go", "build", "-o", bin, "..").CombinedOutput(); err != nil {
		t.Fatalf("the program does not build:\n%s", out)
	}
	out, err := exec.CommandContext(t.Context(), bin, args...).Output()
	var exit *exec.ExitError
	switch {
	case err == nil:
		return out, 0
	case errors.As(err, &exit):
		return out, exit.ExitCode()
	default:
		t.Fatalf("could not run the program: %v", err)
		return nil, 0
	}
}

func TestStats(t *testing.T) {
	cases := []struct {
		name, text   string
		lines, words int
	}{
		{"empty file", "", 0, 0},
		{"last line without newline", "one two\nthree", 2, 3},
		{"blank lines and spaces", "a  b\n\n c\n", 3, 3},
	}
	for _, c := range cases {
		path := filepath.Join(t.TempDir(), "in.txt")
		if err := os.WriteFile(path, []byte(c.text), 0o644); err != nil {
			t.Fatal(err)
		}
		out, code := run(t, path)
		if code != 0 {
			t.Fatalf("%s: exit code %d, want 0", c.name, code)
		}
		// Decode instead of comparing text: spacing and key order do not matter.
		var got struct{ Lines, Words *int }
		if err := json.Unmarshal(out, &got); err != nil || got.Lines == nil || got.Words == nil {
			t.Fatalf("%s: stdout is not {\"lines\":n,\"words\":m}: %q", c.name, out)
		}
		if *got.Lines != c.lines || *got.Words != c.words {
			t.Errorf("%s: got lines=%d words=%d, want lines=%d words=%d", c.name, *got.Lines, *got.Words, c.lines, c.words)
		}
	}
}

func TestUsage(t *testing.T) {
	if _, code := run(t); code != 2 {
		t.Errorf("without an argument the exit code is %d, want 2", code)
	}
}
