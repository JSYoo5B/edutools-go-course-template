# Style: program output

The learner completes a `main` program. Check runs it and compares standard output with the hidden
`test/output.txt` (trailing white space aside). A hidden `test/input.txt`, if present, is fed to standard input.

Print the numbers 1 to 15, one per line. Print `Fizz` instead of multiples of 3, `Buzz` instead of multiples
of 5 and `FizzBuzz` instead of multiples of both.

Use this style when the output itself is what is taught. It cannot grade anything that is not printed, and an exact
text comparison fails on harmless formatting differences.
