# AGENTS.md - EduTools Go course repository

This repository is a working template for building and maintaining **Go learning courses in the JetBrains Academy
(EduTools) format**. The course subject is empty. Clone this repo to start a new course and fill `src/`.
Anything that depends on the subject library lives only in the clone's `src/` and `docs/`.

## Language policy

- **This repository is written in English**: AGENTS.md, README.md, docs/, tool code, comments, commit messages.
- **Course content may be in any language.** Write `task.md`, titles, quiz texts and example comments in the course
  language and set `COURSE["language"]` (and `COURSE["ui"]`, the build's default texts) in `structure.py`.
- The English docs describe *how to build* courses; the course language only affects *what the learner reads*.

## 0. How to read this file

- New here: sections 1-3. While working: section 4 (procedures) and section 6 (definition of done).
- When something breaks in the IDE: `docs/gotchas.md` (symptom -> cause -> fix). Design: `docs/course-design.md`.
  Rules for writing tasks: `docs/authoring-guide.md`.
- Claude Code (v2.1.277+) reads this file directly when there is no `CLAUDE.md`. Do not add a `CLAUDE.md` unless
  it imports this file (`@AGENTS.md`), because a `CLAUDE.md` makes Claude Code ignore `AGENTS.md`.

## 1. Repository map

```
AGENTS.md               instructions for agents (this file)
README.md               quick start for humans
Makefile                build / check / verify / run-demo / run-all / study / zip
structure.py            course metadata (title, summary, language, ui texts) and the section -> lesson map; edit first
src/_root/              go.mod (+go.sum), README.md (copied to the course root), idea/ (IDE module definition)
src/lessonNN/           lessons: lesson.py + taskMM/ (sample lesson01: one theory, one quiz, one exercise; stdlib only)
tools/build.py          src/ -> author course tree (+ learner-state tree); markers, per-task go.mod generation
tools/verify.py         verify one task in isolation (global lock): solution go test / --student go vet / --run
tools/yamlcheck/        structure validator (Go): author tree, or `-study` for the learner tree
tools/to_study/         author tree -> learner (Study-mode) tree converter (Go)
tools/runall.sh         run every edu task serially
docs/                   authoring-guide, gotchas, course-design
build/  .work/          generated output (never committed)
```

**go.mod layout**: the source has exactly one `src/_root/go.mod` (the dependency seed). The build generates a
`go.mod`/`go.sum` **per task** in the output (one task = one Go module). Do not put `go.mod` in task or lesson
source directories.

## 2. Hard rules (break one and the course does not work in the IDE - all of these were hit for real)

1. **The Go checker runs only the tests in the task's `test/` directory.** One task = one Go module: `go.mod` and the
   source `*.go` in the task root, `test/task_test.go` (`package test`, importing the source as `task "<module>"`).
   Module name = `lessonNN_taskMM`. (Violation: Check says "No tests have run".)
2. **Never write `editable: false`** on hidden files (`test/task_test.go`, `go.sum`). The plugin writes the hidden test
   file at Check time. (Violation: "Failed to launch checking".)
3. **Every task, theory and choice included, needs a visible file**; otherwise the Task Description panel does not show.
   For theory/choice tasks that is an **example `main.go` showing the concept**. No empty `main()`. No answer hints in
   a quiz example.
4. The course root needs **`.idea/modules.xml` + `.iml`** (copied from `src/_root/idea/`). Without them the project has
   no module and the top directory shows as excluded.
5. Create placeholders **only with the markers** (`//@@begin` / `//@@ph ` / `//@@end`). Never compute offset/length by
   hand (UTF-16 units, wrong for non-ASCII text).
6. The **learner state** (placeholders substituted) must also compile, and an unimplemented task must **fail fast with a
   clear message** in the course language.
7. Run heavy tests (processes, containers) **only through `tools/verify.py`** (global lock, serial). No parallel
   `go test`. No global delete commands.
8. **One task description = one topic, 3-7 KB.** Interleave: theory -> related exercise -> theory -> exercise.
9. Display order is the **order of `tasks` in `lesson.py`**, not the directory numbers. If you rename a directory, change
   the module name, the import in the test and `lesson.py` together.
10. If the docs and the real behaviour differ, **follow the real run**. State measured claims (timings, error classes,
    defaults) only when you checked them.
11. **A directory with `course-info.yaml` opens in author (Course Creator) mode unless it has a `mode:` key.** To hand
    learners a tree that opens straight in learner mode, generate it with `make study` (section 4). Do not just add
    `mode: Study` by hand: the files and placeholder offsets must be in learner form.
12. **The learner tree must not contain a root `README.md`**, `.idea/workspace.xml` or `.idea/frameworkLessonHistoryCC`.
    If a file is already open when the project opens, the plugin skips its jump to the first task.
13. **Validate freshly generated output, not a tree the IDE has opened.** The IDE rewrites the YAML (indentation, line
    folding) and may re-add `additional_files: README.md`.

## 3. Source at a glance (`src/lessonNN/`)

```
lesson.py                LESSON = dict(dir, title, tasks=[dict(dir, type, title, files, ...)])
taskMM/task.md           description in the course language (all types)
taskMM/main.go           theory/choice: example that shows the concept (visible file)
taskMM/<name>.go         edu: visible source = solution + markers
taskMM/test/task_test.go edu: hidden grading test (package test)
```

- Types: `theory`, `choice` (options live in lesson.py), `edu`.
- List `go.mod`/`go.sum` in `files`; the **build generates them** (never create them in the source). Copy the sample
  `src/lesson01/lesson.py`.
- Marker example:
  ```go
  //@@begin
  //@@ph // TODO: do ...
  //@@ph return nil
  <solution lines>
  //@@end
  ```

## 4. Procedures

### Start a new course
1. Clone this repo. Edit `structure.py` (title, summary, language, ui texts, sections) and `src/_root/README.md`.
2. Research the subject (official docs, public API via `go doc`, examples, pitfalls). Verify facts by running code.
3. Add the subject library to `src/_root/go.mod` (`go get <lib>@<ver>`, `go mod tidy`) and keep go.sum in `src/_root/`
   (copied into every task, then tidied).
4. Design the lessons (`docs/course-design.md`), write `src/lessonNN/` (copy `lesson01` and edit). Replace or delete the sample.
5. `make verify L=lessonNN T=taskMM` per task, then `make check`, then `make run-all`.
6. `make study` (learner tree) -> `make zip` -> verify in the IDE (section 5). Record in the README what was and was
   not verified.

### Update for a new version of the subject library
1. Collect the change log, upgrade guide and the public API difference between the two versions into a table
   (added / changed / removed / behaviour change).
2. Raise the version in `src/_root/go.mod`, regenerate go.sum, `grep` the version strings in titles/README.
3. `grep -rn` the changed or removed symbols in `src/`; fix theory, examples and edu (source + test). New features:
   a theory piece (+ an exercise if valuable).
4. Re-run every measured claim in the text on the new version.
5. Re-verify everything (`make check && make run-all`), also re-read explanations whose output changed even if the tests
   still pass. Regenerate (`make study`) and repackage.

### Add a new package / topic
Write an API table (function | signature | what it does | caveat), add a lesson (theory pieces -> 2-4 exercises -> quiz),
update `structure.py`, the README and the reference card.

### Working with parallel author agents
- Split by lesson (writing is parallel; heavy tests are serial because of the lock). Do not give one agent the heaviest lessons.
- Brief them with section 2, `docs/authoring-guide.md`, the sample lesson and the planned tasks. No edits outside their lessons.
- Treat their reports as hints. After assembling, run `make check` / `make run-all` yourself.

## 5. IDE verification (a human step - it cannot be automated)

- In GoLand with the JetBrains Academy plugin, **close and re-open** `build/course` (author) and `build/course-study`
  (learner), then check: (1) the theory description shows, (2) the exercise Check works, (3) the top directory is not
  excluded, (4) the learner tree jumps to the first task.
- On failure read `~/Library/Logs/JetBrains/GoLand*/idea.log` (`Local check of task` exceptions; the line before
  `Course opened:` shows which file was opened first). The state of an open project can be queried through the
  goland MCP (`get_project_modules`, with `projectPath`).
- Reference for comparison: the official JetBrains *Tour of Go* course (its `task-info.yaml` and `.idea`).
- Building into a directory the IDE has open can delete `.idea/workspace.xml`: close and re-open the project.

## 6. Definition of done

- [ ] `make check` passes (author tree rules + `go vet` of every module, solution and learner state)
- [ ] `make study` passes (`yamlcheck -study` + `go vet` of the learner tree)
- [ ] every edu solution passes `make verify` **at least twice in a row** (no flakiness); the learner state fails fast
- [ ] theory/choice examples compile; representative ones were executed (`make run-demo`)
- [ ] facts in the text (signatures, defaults, measured values) were confirmed by docs / `go doc` / running
- [ ] README version, task counts and verification status match reality; unverified items are stated
- [ ] `make zip` produced the package

## 7. Style (for the course content)

Pick one register in the course language and keep it. No emoji. Introduce foreign terms with a short gloss on first use.
Theory: Build/Inspect tables. Exercise: `# Exercise:` / API table / what to do / how it is graded / caveats.
Code fences always carry a language tag. Details: `docs/authoring-guide.md`.
