# Style: program behaviour

The learner completes a `main` program, but grading is not a text comparison. The hidden test builds the program
with `go build ..` (the task root, seen from `test/`; `go run` would hide the exit code), runs it and checks what it does: the exit code, files it writes, or
output it decodes (JSON here), so formatting does not matter.

Complete `main.go`: count the lines and the words of the file and fill `s`. A last line without a trailing newline
still counts as a line.

Use this style when the subject is a command or when the program's effects matter more than its exact output.
