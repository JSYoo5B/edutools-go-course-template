package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStudyAnswerModel(t *testing.T) {
	for _, tc := range []struct{ name, model, answer, want string }{
		{"valid", "text: TODO", "possible_answer: solved", ""},
		{"missing text", "", "possible_answer: solved", "model text"},
		{"wrong text", "text: other", "possible_answer: solved", "model text"},
		{"missing answer", "text: TODO", "", "missing possible_answer"},
		{"encrypted answer", "text: TODO", "encrypted_possible_answer: ciphertext", "encrypted answers"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			for name, data := range map[string]string{
				"course-info.yaml": "mode: Study\n",
				"task.go":          "TODO",
				"task-info.yaml":   "status: Unchecked\nrecord: -1\nfiles:\n- name: task.go\n  " + tc.model + "\n  placeholders:\n  - offset: 0\n    length: 4\n    placeholder_text: TODO\n    initial_state: {offset: 0, length: 4}\n    " + tc.answer + "\n",
			} {
				if err := os.WriteFile(filepath.Join(root, name), []byte(data), 0600); err != nil {
					t.Fatal(err)
				}
			}
			errs = nil
			t.Cleanup(func() { errs = nil })
			checkStudyTree(root)
			got := strings.Join(errs, "\n")
			if tc.want == "" && got != "" || tc.want != "" && !strings.Contains(got, tc.want) {
				t.Fatalf("errors=%q, want %q", got, tc.want)
			}
		})
	}
}
