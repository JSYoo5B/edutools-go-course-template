# Style: write the test

Here the learner writes a test. The code under test is given and correct. Grading cannot compare output, so a
hidden test runs the learner's `TestClamp` against broken copies of `clamp.go` (mutants, swapped in with
`go test -overlay`) and expects it to fail on each one. It must also pass on the real code, because Check runs it
too.

Write `TestClamp` in `test/clamp_test.go`. Cover values below, inside and above the range.

Use this style to teach testing itself: assertions, table tests, edge cases. Pick mutants that a reasonable test
catches, and name the bug in the failure message without giving away the test.
