// yamlcheck validates an EduTools course tree produced by build.py.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode/utf16"

	"gopkg.in/yaml.v3"
)

type CourseInfo struct {
	Title       string   `yaml:"title"`
	Language    string   `yaml:"language"`
	Summary     string   `yaml:"summary"`
	ProgLang    string   `yaml:"programming_language"`
	ProgLangVer string   `yaml:"programming_language_version"`
	Content     []string `yaml:"content"`
	Additional  []struct {
		Name    string `yaml:"name"`
		Visible bool   `yaml:"visible"`
	} `yaml:"additional_files"`
}
type ItemInfo struct {
	CustomName string   `yaml:"custom_name"`
	Content    []string `yaml:"content"`
}
type Placeholder struct {
	Offset int    `yaml:"offset"`
	Length int    `yaml:"length"`
	Text   string `yaml:"placeholder_text"`
}
type TaskFile struct {
	Name         string        `yaml:"name"`
	Visible      bool          `yaml:"visible"`
	Editable     *bool         `yaml:"editable"`
	Placeholders []Placeholder `yaml:"placeholders"`
}
type Option struct {
	Text    string `yaml:"text"`
	Correct bool   `yaml:"is_correct"`
}
type TaskInfo struct {
	Type       string     `yaml:"type"`
	CustomName string     `yaml:"custom_name"`
	Files      []TaskFile `yaml:"files"`
	Multiple   bool       `yaml:"is_multiple_choice"`
	Options    []Option   `yaml:"options"`
}

var errs []string

func fail(format string, a ...any) { errs = append(errs, fmt.Sprintf(format, a...)) }

func readYAML(path string, v any) bool {
	b, err := os.ReadFile(path)
	if err != nil {
		fail("%s: %v", path, err)
		return false
	}
	if err := yaml.Unmarshal(b, v); err != nil {
		fail("%s: yaml: %v", path, err)
		return false
	}
	return true
}

func subdirs(dir string) []string {
	es, _ := os.ReadDir(dir)
	var out []string
	for _, e := range es {
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out
}

func checkContent(dir string, content []string) {
	listed := map[string]bool{}
	for _, c := range content {
		listed[c] = true
		if st, err := os.Stat(filepath.Join(dir, c)); err != nil || !st.IsDir() {
			fail("%s: content entry %q is not a directory", dir, c)
		}
	}
	for _, d := range subdirs(dir) {
		if !listed[d] {
			fail("%s: directory %q is not listed in content", dir, d)
		}
	}
}

var pkgRe = regexp.MustCompile(`(?m)^package\s+(\w+)`)

var study bool

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "-study" {
		study = true
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: yamlcheck [-study] <course_dir>")
		os.Exit(2)
	}
	root := args[0]
	var ci CourseInfo
	if !readYAML(filepath.Join(root, "course-info.yaml"), &ci) {
		report()
	}
	if ci.Title == "" || ci.Language == "" || ci.ProgLang != "Go" {
		fail("course-info: title, language and programming_language must be set (and programming_language must be Go)")
	}
	for _, a := range ci.Additional {
		if _, err := os.Stat(filepath.Join(root, a.Name)); err != nil {
			fail("additional file %s missing", a.Name)
		}
	}
	checkContent(root, ci.Content)
	pkgs := map[string]string{}
	nTasks := map[string]int{}
	for _, s := range ci.Content {
		var si ItemInfo
		sdir := filepath.Join(root, s)
		if !readYAML(filepath.Join(sdir, "section-info.yaml"), &si) {
			continue
		}
		if si.CustomName == "" {
			fail("%s: no custom_name", sdir)
		}
		checkContent(sdir, si.Content)
		for _, l := range si.Content {
			var li ItemInfo
			ldir := filepath.Join(sdir, l)
			if !readYAML(filepath.Join(ldir, "lesson-info.yaml"), &li) {
				continue
			}
			if li.CustomName == "" {
				fail("%s: no custom_name", ldir)
			}
			checkContent(ldir, li.Content)
			for _, t := range li.Content {
				tdir := filepath.Join(ldir, t)
				var ti TaskInfo
				if !readYAML(filepath.Join(tdir, "task-info.yaml"), &ti) {
					continue
				}
				nTasks[ti.Type]++
				if ti.CustomName == "" {
					fail("%s: no custom_name", tdir)
				}
				md, err := os.ReadFile(filepath.Join(tdir, "task.md"))
				if err != nil || len(strings.TrimSpace(string(md))) < 20 {
					fail("%s: task.md missing/empty", tdir)
				}
				switch ti.Type {
				case "theory":
					if !hasVisibleFile(ti) {
						fail("%s: theory without a visible file (description panel would not show)", tdir)
					}
				case "choice":
					if !hasVisibleFile(ti) {
						fail("%s: choice without a visible file (description panel would not show)", tdir)
					}
					ok := 0
					for _, o := range ti.Options {
						if o.Correct {
							ok++
						}
					}
					if len(ti.Options) < 3 || ok == 0 || (!ti.Multiple && ok != 1) {
						fail("%s: bad choice options (n=%d correct=%d multiple=%v)", tdir, len(ti.Options), ok, ti.Multiple)
					}
				case "edu":
					if len(ti.Files) < 2 {
						fail("%s: edu needs >=2 files", tdir)
					}
					hasPH, hasHiddenTest := false, false
					for _, f := range ti.Files {
						b, err := os.ReadFile(filepath.Join(tdir, f.Name))
						if err != nil {
							fail("%s: file %s: %v", tdir, f.Name, err)
							continue
						}
						if strings.HasSuffix(f.Name, "_test.go") {
							if f.Visible {
								fail("%s: test file %s must be hidden", tdir, f.Name)
							}
							hasHiddenTest = true
						}
						if strings.Contains(string(b), "//@@") {
							fail("%s: %s still has //@@ markers", tdir, f.Name)
						}
						if strings.HasSuffix(f.Name, "_test.go") && !strings.HasPrefix(f.Name, "test/") {
							fail("%s: test file %s must live in test/", tdir, f.Name)
						}
						if m := pkgRe.FindStringSubmatch(string(b)); m != nil && !strings.HasPrefix(f.Name, "test/") {
							if prev, dup := pkgs[m[1]]; dup && prev != tdir {
								fail("%s: package %q also used in %s", tdir, m[1], prev)
							}
							pkgs[m[1]] = tdir
						}
						u := utf16.Encode([]rune(string(b)))
						last := 0
						phs := append([]Placeholder(nil), f.Placeholders...)
						sort.Slice(phs, func(i, j int) bool { return phs[i].Offset < phs[j].Offset })
						for _, p := range phs {
							hasPH = true
							if p.Offset < last || p.Length <= 0 || p.Offset+p.Length > len(u) {
								fail("%s: %s: bad placeholder %+v (file len %d)", tdir, f.Name, p, len(u))
								continue
							}
							last = p.Offset + p.Length
							sol := string(utf16.Decode(u[p.Offset : p.Offset+p.Length]))
							if strings.TrimSpace(sol) == "" || strings.TrimSpace(p.Text) == "" {
								fail("%s: %s: empty placeholder region/text", tdir, f.Name)
							}
							if !f.Visible {
								fail("%s: hidden file %s has placeholders", tdir, f.Name)
							}
						}
					}
					gm, err := os.ReadFile(filepath.Join(tdir, "go.mod"))
					if err != nil || !strings.HasPrefix(string(gm), "module ") {
						fail("%s: edu task needs its own go.mod", tdir)
					}
					if !hasPH {
						fail("%s: edu task without placeholders", tdir)
					}
					if !hasHiddenTest {
						fail("%s: edu task without hidden *_test.go", tdir)
					}
				default:
					fail("%s: unknown task type %q", tdir, ti.Type)
				}
			}
		}
	}
	if study {
		checkStudyTree(root)
	} else {
		checkAuthorTree(root)
	}
	fmt.Printf("tasks: %v\n", nTasks)
	report()
}

func report() {
	if len(errs) == 0 {
		fmt.Println("OK: course tree is consistent")
		return
	}
	for _, e := range errs {
		fmt.Println("ERROR:", e)
	}
	os.Exit(1)
}

func hasVisibleFile(ti TaskInfo) bool {
	for _, f := range ti.Files {
		if f.Visible {
			return true
		}
	}
	return false
}

// checkAuthorTree: rules that apply to the author (Course Creator) tree.
func checkAuthorTree(root string) {
	filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || d.Name() != "task-info.yaml" {
			return nil
		}
		b, _ := os.ReadFile(p)
		if strings.Contains(string(b), "editable:") {
			// "Failed to launch checking": the plugin writes the hidden test file at Check time
			fail("%s: contains `editable:` (read-only hidden files break Check)", p)
		}
		return nil
	})
	if b, err := os.ReadFile(filepath.Join(root, "course-info.yaml")); err == nil {
		if regexp.MustCompile(`(?m)^mode:`).Match(b) {
			fail("course-info.yaml has `mode:` - the author tree must open in Course Creator mode (use `make study` for the learner tree)")
		}
	}
	if _, err := os.Stat(filepath.Join(root, ".idea", "modules.xml")); err != nil {
		fail(".idea/modules.xml is missing - the project would have no module (everything shows as excluded)")
	}
}

// checkStudyTree: rules that apply to the learner (Study-mode) tree.
func checkStudyTree(root string) {
	b, err := os.ReadFile(filepath.Join(root, "course-info.yaml"))
	if err != nil || !regexp.MustCompile(`(?m)^mode:\s*Study`).Match(b) {
		fail("course-info.yaml must contain `mode: Study`")
	}
	if _, err := os.Stat(filepath.Join(root, "README.md")); err == nil {
		fail("root README.md must not exist in the learner tree (it is opened first and the plugin then skips the jump to the first task)")
	}
	if regexp.MustCompile(`(?m)^additional_files:`).Match(b) && strings.Contains(string(b), "README.md") {
		fail("course-info.yaml lists README.md in additional_files")
	}
	for _, leaked := range []string{".idea/workspace.xml", ".idea/frameworkLessonHistoryCC"} {
		if _, err := os.Stat(filepath.Join(root, leaked)); err == nil {
			fail("%s must not be in the learner tree (author-only IDE state)", leaked)
		}
	}
	filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || d.Name() != "task-info.yaml" {
			return nil
		}
		dir := filepath.Dir(p)
		var raw map[string]any
		data, _ := os.ReadFile(p)
		if err := yaml.Unmarshal(data, &raw); err != nil {
			fail("%s: %v", p, err)
			return nil
		}
		if raw["status"] == nil || raw["record"] == nil {
			fail("%s: learner tasks need `status` and `record`", p)
		}
		files, _ := raw["files"].([]any)
		for _, fe := range files {
			fm, _ := fe.(map[string]any)
			phs, _ := fm["placeholders"].([]any)
			if len(phs) == 0 {
				continue
			}
			name, _ := fm["name"].(string)
			text, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				fail("%s: %v", p, err)
				continue
			}
			u := utf16.Encode([]rune(string(text)))
			for _, pe := range phs {
				pm, _ := pe.(map[string]any)
				off, _ := pm["offset"].(int)
				ln, _ := pm["length"].(int)
				want, _ := pm["placeholder_text"].(string)
				if pm["initial_state"] == nil {
					fail("%s/%s: placeholder without initial_state", dir, name)
				}
				if off < 0 || off+ln > len(u) || string(utf16.Decode(u[off:off+ln])) != want {
					fail("%s/%s: file text at offset %d length %d is not the placeholder text", dir, name, off, ln)
				}
			}
		}
		return nil
	})
}
