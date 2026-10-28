# Authoring guide (for authors and agents)

Audience: learners who know Go and testing basics but are new to the subject. Goal: they learn **which APIs build what,
how to verify the result, and (where it applies) how to handle errors and edge cases** - by writing code.

## Source layout (`src/`)

```
src/_root/go.mod, go.sum, README.md, idea/{course.iml,modules.xml,externalDependencies.xml}
src/lessonNN/lesson.py                 # LESSON dict
src/lessonNN/taskMM/task.md            # description in the course language
src/lessonNN/taskMM/<file>.go          # edu: visible source = solution + markers / theory,choice: example main.go
src/lessonNN/taskMM/test/task_test.go  # edu: hidden grading test (package test)
```
The task directory number is **not** the display order (that is the order in `lesson.py`). The module is named
`lessonNN_taskMM`, so renaming a directory means changing the import in the test and `lesson.py` too.

### lesson.py

```python
LESSON = dict(
  dir="lesson03", title="Lesson title",
  tasks=[
    dict(dir="task01", type="theory", title="...",
         files=[("go.mod", False), ("go.sum", False), ("main.go", True)]),      # the example code is the visible file
    dict(dir="task02", type="edu", title="Exercise: ...",
         files=[("go.mod", True), ("go.sum", False), ("exercise.go", True), ("test/task_test.go", False)]),
    dict(dir="task06", type="choice", title="Quiz: ...", multiple=False,
         options=[("option", True), ("option", False), ...],                    # about four, plausible wrong ones
         msg_ok="Correct! why", msg_bad="a hint",                               # optional; defaults from structure.COURSE["ui"]
         files=[("go.mod", False), ("go.sum", False), ("main.go", True)]),
  ])
```
- Do **not** create go.mod / go.sum in the source: the build generates them from `src/_root/go.mod` and tidies them.
- Order: theory piece -> related edu -> next theory piece -> edu ... -> choice quiz at the end of the lesson.
- If `files` is omitted for theory/choice the build inserts an empty scratch `main.go`. **Do not rely on it** (write a
  real example).

## edu task files

### Visible source (solution + markers)

```go
func WordCount(text string) map[string]int {
	//@@begin
	//@@ph // TODO: 1) do ...
	//@@ph _ = strings.Fields
	//@@ph return nil
	counts := make(map[string]int)        // <- the solution (between the markers)
	...
	return counts
	//@@end
}
```
- Lines starting with `//@@ph ` right after `//@@begin` are the placeholder text the learner sees; everything up to
  `//@@end` is the solution. The markers are removed from the output.
- The learner state must compile -> silence unused imports/variables with `_ = x`. Check with `make verify L=.. T=..`
  (= `--student` vet + the solution test).
- In the learner state the grading test must fail **fast** with a clear message in the course language (for example when
  a function returns the zero value, fail before any heavy setup). A package-level `ErrNotImplemented` sentinel and a
  probe call (nil input / cancelled context) is one proven pattern.
- Package names must be unique per task, and start with a package doc comment.
- Keep the boilerplate error handling (`t.Fatal`) outside the placeholder so the learner does not touch the signature.

### Hidden test `test/task_test.go`

- `package test`; the source is imported as `task "<module>"` (module = `lessonNN_taskMM`), so only **exported**
  identifiers are visible. Export whatever the test needs.
- Use `t.Context()`, but in cleanup use a **fresh context** (`context.WithTimeout(context.Background(), ...)`):
  `t.Context()` is already cancelled when cleanup runs.
- Safety net: terminate what the learner code may forget to clean up, in the test's own `t.Cleanup`.
- Poll anything that converges asynchronously (background work, caches, statistics) with a deadline; never assert at once.
- Failure messages: what was expected, what was seen, which API to look at - in the course language.
- Keep a test light (one heavy resource), at most 3-4 minutes. No `t.Parallel` when tests are heavy.
- **No flakiness**: pass at least twice in a row and record the duration.
- Verify "did the learner clean up" by checking that the resource is gone after the subtest ends.
- If the documentation and the real behaviour differ, follow the real run and reflect it in the text.

## theory / choice tasks: the example `main.go`

- `package main`, must compile (`--student` vet). Prefer a demo that **runs with `go run .`**: prepare -> run the concept ->
  print -> `defer` cleanup.
- Do not paste the snippets of task.md: weave them into one coherent story of 40-120 lines with comments that explain
  *why* and point out pitfalls.
- A risky or slow example may be compile-only with a note at the top. **No empty main.**
- A quiz example contains the code of the question and **no hint about the answer**.
- Do not expose the solution of an earlier exercise.

## task.md structure

- One register in the course language, no emoji, a short gloss for foreign terms on first use.
- **theory**: concept + **Build/Inspect tables** (`function | signature | what it does | caveat`). Every signature and
  behaviour confirmed by docs, `go doc` or running. 3-7 KB per piece.
- **edu**: `# Exercise: title` / APIs used (table) / what to do (numbered) / how it is graded / caveats. Show the key calls
  but not the whole answer.
- **choice**: the question only (options live in lesson.py).
- Code fences carry a language tag (```go, ```bash).

## Parallel author agents

- Split by lesson. Writing is parallel; heavy tests are serial (global lock in `tools/verify.py`).
  Do not hand the heaviest lessons to one agent.
- Brief: this document, `AGENTS.md` section 2, the sample lesson, the planned tasks per lesson (APIs, test ideas),
  "follow the real run when the docs differ and report it".
- Forbidden: editing outside their lessons, global delete commands, editing the subject library.
- Report: APIs taught per task, measured durations, docs-vs-reality differences or flakiness, extra modules, TODOs.
- A report is a hint. After assembling, run the verification yourself.
