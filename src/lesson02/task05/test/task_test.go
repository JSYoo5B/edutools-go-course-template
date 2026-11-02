package test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// mutants are broken versions of clamp.go. The learner's TestClamp must fail on each of them.
var mutants = []struct {
	bug      string
	old, new string
}{
	{"values below the range are not raised", "return lo", "return v"},
	{"values above the range are not lowered", "return hi", "return v"},
	{"values above the range become the lower bound", "return hi", "return lo"},
}

// TestYourTestCatchesBugs runs TestClamp against each mutant with go test -overlay, which swaps
// clamp.go for the mutant without touching the learner's files.
func TestYourTestCatchesBugs(t *testing.T) {
	orig, err := filepath.Abs(filepath.Join("..", "clamp.go"))
	if err != nil {
		t.Fatal(err)
	}
	src, err := os.ReadFile(orig)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range mutants {
		if !strings.Contains(string(src), m.old) {
			t.Fatalf("course error: mutant %q does not apply to clamp.go", m.bug)
		}
		dir := t.TempDir()
		mutant := filepath.Join(dir, "clamp.go")
		if err := os.WriteFile(mutant, []byte(strings.Replace(string(src), m.old, m.new, 1)), 0o644); err != nil {
			t.Fatal(err)
		}
		overlay, _ := json.Marshal(map[string]map[string]string{"Replace": {orig: mutant}})
		overlayFile := filepath.Join(dir, "overlay.json")
		if err := os.WriteFile(overlayFile, overlay, 0o644); err != nil {
			t.Fatal(err)
		}
		cmd := exec.CommandContext(t.Context(), "go", "test", "-count=1", "-overlay", overlayFile, "-run", "^TestClamp$", ".")
		out, err := cmd.CombinedOutput()
		if err == nil {
			t.Errorf("TestClamp passes even when %s. Add a case that catches it.", m.bug)
		} else if strings.Contains(string(out), "build failed") || strings.Contains(string(out), "[setup failed]") {
			t.Fatalf("could not build against the mutant:\n%s", out)
		}
	}
}
