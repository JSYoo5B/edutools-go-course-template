# Exercise styles

Pick the style from how the subject is used, not from habit. Each style has a runnable sample in `src/lesson01`
(function fragment) and `src/lesson02` (the others); copy the one you need.

| Style | Learner edits | Graded by | Sample |
| --- | --- | --- | --- |
| Function fragment | a function body | hidden test calls the function | `lesson01/task03` |
| Program output | `main` | stdout vs hidden `test/output.txt` (`type: output`) | `lesson02/task01` |
| Program behaviour | `main` | hidden test builds and runs it; exit code, files, decoded output | `lesson02/task02` |
| Fragment + starter test | a function body | visible starter test (learner can run it) + hidden behaviour test | `lesson02/task03` |
| Code inside a test | the body of a visible `TestXxx` | given assertions in that test + hidden guard | `lesson02/task04` |
| Write the test | a whole visible `TestXxx` | hidden test runs it against mutants | `lesson02/task05` |

Rough guide: a library -> fragment (+ starter test); a command-line tool -> program behaviour; a testing tool or an
API used from tests -> code inside a test; testing skills -> write the test. Use program output only when the
printed text is the point.

## How Check works (from the plugin source)

- `edu`: Check runs every Go test in the task's `test/` directory; all must pass. Visible test files in `test/` run
  too.
- `output`: Check runs the file that has `main()`, feeds hidden `test/input.txt` (if any) to stdin and compares stdout
  with hidden `test/output.txt`, ignoring trailing white space.
- In a learner tree, Check first writes every **hidden** file of `test/` from the task model (`CheckAction.createTests`)
  and may empty them again afterwards. Visible files are left alone. So:
  - never let a visible file depend on a hidden one (it may be empty on disk between checks);
  - the model must carry the hidden files' text (`tools/to_study` does it).

## Grading without comparing output

- **Decode, then compare values.** Parse JSON or another structured result instead of matching text.
- **Exit codes and side effects.** Build the program (`go build -o <tmp> ..` from `test/`) and run the binary;
  `go run` turns any non-zero exit of the program into 1. Check files written, not messages printed.
- **Recording fakes.** Make the code take an interface and pass a fake that records calls; assert on what happened
  (no writes after an error, call order, arguments).
- **Error classes.** Assert with `errors.Is` / `errors.As` on exported sentinel or typed errors, not on messages.
- **Invariants and edge cases.** Several inputs, boundaries, empty and invalid values; hide them, keep one obvious
  case visible in a starter test.
- **Mutants for learner-written tests.** Swap the code under test for broken copies with `go test -overlay` and
  require the learner's test to fail on each; it must pass on the real code too. Choose bugs a reasonable test
  catches and name the bug, not the missing case.
- **Guard given code.** When visible assertions grade the exercise, a hidden test checks they are unchanged (read the
  file, look for the exact block).

## Rules that apply to every style

- Visible and hidden test files share `package test`: give hidden helpers distinct names so a learner's own helper
  does not collide with them.
- The placeholder holds the concept being taught. Setup, cleanup and routine error handling stay outside it.
- The learner state compiles (`_ = x` for unused names) and fails with a message that says what is missing.
- A hidden test that runs `go` itself (build, `-overlay`) relies on the Go toolchain the IDE already uses; keep it to a
  few short runs.
