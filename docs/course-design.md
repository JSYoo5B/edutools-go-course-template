# Course design frame

## Teaching principles
- **Small steps, in the order the learner needs them.** One task introduces one concept. More lessons are fine if each
  one is understood; never front-load a list of every option.
- **An example shows the behaviour it explains.** The learner runs it and sees the concept happen (state changes,
  output, an error). Code that only prints the values from the question teaches nothing.
- **An exercise combines what was just taught** into working code and lets the learner observe the result. Change a
  condition, compare with a baseline, check against what the subject itself reports.
- **Exercises use the subject the way its users do.** For a library, the learner writes code that calls its API; for a
  testing tool, the learner writes a Go test (`TestXxx(t *testing.T)`) that uses it. Supply the boilerplate (setup,
  cleanup, routine error handling) and leave the meaningful calls to the learner.
- **Do not force one shape on every lesson.** Pick the number and kind of tasks from the topic, not from a template.
- **Respect the audience.** Assume what the brief says they know; explain the concepts, not the obvious.

Reference course: JetBrains Academy's adaptation of *A Tour of Go* -
[source](https://github.com/jetbrains-academy/tour-of-go) (author tree, same layout as this template's output).
Look at how a topic is built up from small concept tasks before an exercise ties them together, and at how every
theory task has a runnable example. Link to it; do not copy its content.

## A possible lesson layout (adapt freely)

| Part | Lessons | Content |
| --- | --- | --- |
| 1 Getting started | setup, core concepts | the basic flow, a first exercise |
| 2 Core features | one lesson per feature group | how to build things -> how to inspect them |
| 3 Advanced | one lesson per sub-package or topic | |
| 4 Errors and edges | failures, recovery, error classes | |
| 5 Reference | cheat sheets, helper code to copy | |

## Kinds of edu exercises (examples, not a quota)
1. **Build**: create something with a combination of options/functions and clean it up.
2. **Inspect**: read state; the test compares with what the subject itself reports.
3. **Failure**: inject, observe, restore; the test also checks the recovery.
4. **Rejection path**: provoke an invalid combination and check the error class.

## Quizzes
- About four plausible options; wrong options come from common misconceptions. Give a correct-answer message and a hint
  for a wrong answer.

## Checklist for a new version / package
- Signatures, defaults and constraints in the text.
- Measured claims: durations, delays, error classes.
- Supported environments and the required Go version.
- README (version, task counts).
