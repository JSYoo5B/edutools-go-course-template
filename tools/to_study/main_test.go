package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Peek Solution needs possible_answer on each placeholder and the full learner text of the file;
// non-ASCII text checks the UTF-16 offsets.
func TestConvertTaskKeepsPeekSolution(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	prefix := "// é 😀\npackage main\n\nfunc main() {\n\t"
	answer := "// é 😀\nprintln(\"solved\")"
	suffix := "\n}\n"
	placeholder := "// TODO: print"
	for _, dir := range []string{src, dst} {
		if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(prefix+answer+suffix), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	root := mapping(scalar("type"), scalar("edu"), scalar("files"), seq(mapping(
		scalar("name"), scalar("main.go"), scalar("visible"), boolNode(true),
		scalar("placeholders"), seq(mapping(
			scalar("offset"), intNode(len(u16(prefix))),
			scalar("length"), intNode(len(u16(answer))),
			scalar("placeholder_text"), quoted(placeholder))),
	)))
	if err := writeYAML(filepath.Join(src, "task-info.yaml"), root); err != nil {
		t.Fatal(err)
	}
	if problems := convertTask(src, dst); len(problems) > 0 {
		t.Fatal(problems)
	}
	converted, err := readYAML(filepath.Join(dst, "task-info.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	file := get(converted, "files").Content[0]
	if got := get(get(file, "placeholders").Content[0], "possible_answer"); got == nil || got.Value != answer {
		t.Fatalf("possible_answer = %v, want %q", got, answer)
	}
	learner, err := os.ReadFile(filepath.Join(dst, "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	if string(learner) != prefix+placeholder+suffix || strings.Contains(string(learner), "solved") {
		t.Fatalf("learner file is wrong: %q", learner)
	}
	if text := get(file, "text"); text == nil || text.Value != string(learner) {
		t.Fatal("file model text must equal the learner file")
	}
}

// Check restores hidden files from the model, so a hidden file without placeholders keeps its text too.
func TestConvertTaskKeepsHiddenFileText(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	contents := "package test\n\nfunc check() {}\n"
	for _, dir := range []string{src, dst} {
		if err := os.MkdirAll(filepath.Join(dir, "test"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "test/task_test.go"), []byte(contents), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	root := mapping(scalar("type"), scalar("edu"), scalar("files"), seq(mapping(
		scalar("name"), scalar("test/task_test.go"), scalar("visible"), boolNode(false))))
	if err := writeYAML(filepath.Join(src, "task-info.yaml"), root); err != nil {
		t.Fatal(err)
	}
	if problems := convertTask(src, dst); len(problems) > 0 {
		t.Fatal(problems)
	}
	converted, err := readYAML(filepath.Join(dst, "task-info.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if text := get(get(converted, "files").Content[0], "text"); text == nil || text.Value != contents {
		t.Fatal("hidden file text is missing from the model")
	}
}
