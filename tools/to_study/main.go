// to_study converts an author (Course Creator) course tree into a learner (Study-mode) tree.
//
// Usage: to_study <author_course_dir> <out_dir>
//
// The plugin opens a course in learner mode if course-info.yaml has a `mode:` key, and in author mode otherwise.
// A learner-mode tree (1) has the placeholder text in the files instead of the solution,
// (2) has placeholder offset/length in task-info.yaml relative to that text (including initial_state), and
// (3) has status/record on every task. YAML is read and written with yaml.v3, so the format the IDE rewrites is handled too.
package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf16"

	"gopkg.in/yaml.v3"
)

func u16(s string) []uint16 { return utf16.Encode([]rune(s)) }
func fromU16(u []uint16) string { return string(utf16.Decode(u)) }

func scalar(v string) *yaml.Node { return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: v} }

// quoted forces the double-quoted style (\n, \t escapes). The block style may be rejected by the IDE parser when the text contains tabs.
func quoted(v string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: v, Style: yaml.DoubleQuotedStyle}
}
func intNode(n int) *yaml.Node  { return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: fmt.Sprint(n)} }
func boolNode(b bool) *yaml.Node {
	v := "false"
	if b {
		v = "true"
	}
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: v}
}
func mapping(kv ...*yaml.Node) *yaml.Node { return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map", Content: kv} }
func seq(items ...*yaml.Node) *yaml.Node  { return &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Content: items} }

func get(m *yaml.Node, key string) *yaml.Node {
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			return m.Content[i+1]
		}
	}
	return nil
}
func has(m *yaml.Node, key string) bool { return get(m, key) != nil }
func add(m *yaml.Node, key string, v *yaml.Node) { m.Content = append(m.Content, scalar(key), v) }

func readYAML(path string) (*yaml.Node, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(b, &doc); err != nil {
		return nil, err
	}
	return doc.Content[0], nil
}

func writeYAML(path string, n *yaml.Node) error {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(n); err != nil {
		return err
	}
	enc.Close()
	return os.WriteFile(path, buf.Bytes(), 0o644)
}

func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, p)
		base := d.Name()
		if base == ".DS_Store" || base == "workspace.xml" || strings.HasPrefix(base, ".author_contents_storage_db") || base == "frameworkLessonHistoryCC" {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(target, b, 0o644)
	})
}

type ph struct {
	offset, length int
	text           string
	node           *yaml.Node
}

func convertTask(srcDir, dstDir string) []string {
	var problems []string
	root, err := readYAML(filepath.Join(srcDir, "task-info.yaml"))
	if err != nil {
		return []string{fmt.Sprintf("%s: %v", srcDir, err)}
	}
	ttype := ""
	if t := get(root, "type"); t != nil {
		ttype = t.Value
	}
	if files := get(root, "files"); files != nil && files.Kind == yaml.SequenceNode {
		for _, f := range files.Content {
			name := get(f, "name").Value
			phsNode := get(f, "placeholders")
			if phsNode != nil && len(phsNode.Content) > 0 {
				var phs []ph
				for _, p := range phsNode.Content {
					var o, l int
					fmt.Sscan(get(p, "offset").Value, &o)
					fmt.Sscan(get(p, "length").Value, &l)
					phs = append(phs, ph{offset: o, length: l, text: get(p, "placeholder_text").Value})
				}
				b, err := os.ReadFile(filepath.Join(srcDir, name))
				if err != nil {
					problems = append(problems, fmt.Sprintf("%s/%s: %v", srcDir, name, err))
					continue
				}
				u := u16(string(b))
				var out []uint16
				pos, delta := 0, 0
				newPhs := make([]*yaml.Node, 0, len(phs))
				for _, p := range phs { // offsets must be ascending within a file
					if p.offset < pos || p.offset+p.length > len(u) {
						problems = append(problems, fmt.Sprintf("%s/%s: bad placeholder %d+%d", srcDir, name, p.offset, p.length))
						continue
					}
					out = append(out, u[pos:p.offset]...)
					pu := u16(p.text)
					newOff := p.offset + delta
					out = append(out, pu...)
					pos = p.offset + p.length
					delta += len(pu) - p.length
					newPhs = append(newPhs, mapping(
						scalar("offset"), intNode(newOff),
						scalar("length"), intNode(len(pu)),
						scalar("placeholder_text"), quoted(p.text),
						scalar("initial_state"), mapping(scalar("length"), intNode(len(pu)), scalar("offset"), intNode(newOff)),
						scalar("initialized_from_dependency"), boolNode(false),
						scalar("selected"), boolNode(false),
						scalar("status"), scalar("Unchecked"),
					))
				}
				out = append(out, u[pos:]...)
				newText := fromU16(out)
				// self-check: the converted file region must equal the placeholder text
				nu := u16(newText)
				for i, p := range newPhs {
					var o, l int
					fmt.Sscan(get(p, "offset").Value, &o)
					fmt.Sscan(get(p, "length").Value, &l)
					if fromU16(nu[o:o+l]) != get(p, "placeholder_text").Value {
						problems = append(problems, fmt.Sprintf("%s/%s: placeholder %d mismatch", srcDir, name, i))
					}
				}
				if err := os.WriteFile(filepath.Join(dstDir, name), []byte(newText), 0o644); err != nil {
					problems = append(problems, err.Error())
				}
				phsNode.Content = newPhs
			}
			if !has(f, "learner_created") {
				add(f, "learner_created", boolNode(false))
			}
		}
	}
	if ttype == "choice" && !has(root, "selected_options") {
		add(root, "selected_options", seq())
	}
	if !has(root, "status") {
		add(root, "status", scalar("Unchecked"))
	}
	if !has(root, "record") {
		add(root, "record", intNode(-1))
	}
	if err := writeYAML(filepath.Join(dstDir, "task-info.yaml"), root); err != nil {
		problems = append(problems, err.Error())
	}
	return problems
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: to_study <author_course_dir> <out_dir>")
		os.Exit(2)
	}
	src, dst := os.Args[1], os.Args[2]
	os.RemoveAll(dst)
	if err := copyTree(src, dst); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var problems []string
	n := 0
	filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && d.Name() == "task-info.yaml" {
			rel, _ := filepath.Rel(src, filepath.Dir(p))
			problems = append(problems, convertTask(filepath.Dir(p), filepath.Join(dst, rel))...)
			n++
		}
		return nil
	})
	// course-info.yaml: a `mode` key means learner mode
	cp := filepath.Join(dst, "course-info.yaml")
	root, err := readYAML(cp)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if !has(root, "mode") {
		// insert before `content` (key order is irrelevant but easier to read)
		idx := len(root.Content)
		for i := 0; i+1 < len(root.Content); i += 2 {
			if root.Content[i].Value == "content" {
				idx = i
				break
			}
		}
		kv := []*yaml.Node{scalar("mode"), scalar("Study")}
		root.Content = append(root.Content[:idx], append(kv, root.Content[idx:]...)...)
	}
	// Remove the root README.md from the learner tree. If README.md is opened first when the project opens,
	// the plugin skips its "only when no file is open" jump to the first task (EduProjectActivity).
	os.Remove(filepath.Join(dst, "README.md"))
	if af := get(root, "additional_files"); af != nil && af.Kind == yaml.SequenceNode {
		kept := af.Content[:0]
		for _, e := range af.Content {
			if n := get(e, "name"); n != nil && n.Value == "README.md" {
				continue
			}
			kept = append(kept, e)
		}
		af.Content = kept
		if len(kept) == 0 {
			for i := 0; i+1 < len(root.Content); i += 2 {
				if root.Content[i].Value == "additional_files" {
					root.Content = append(root.Content[:i], root.Content[i+2:]...)
					break
				}
			}
		}
	}
	if err := writeYAML(cp, root); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("converted %d tasks -> %s\n", n, dst)
	if len(problems) > 0 {
		fmt.Println(strings.Join(problems, "\n"))
		os.Exit(1)
	}
}
