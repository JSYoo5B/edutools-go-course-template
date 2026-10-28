# Course design frame

## Principles
- **Learning goal -> lessons -> (theory piece -> exercise -> ...) -> quiz.** One task.md = one topic, 3-7 KB.
- For a library/tool course, make every lesson show **what builds things (Build) and what verifies them (Inspect)** in
  tables. If the tool deals with failures, add **inject -> observe -> recover** as a third axis.
- End with a **reference card** (cheat sheets for build / inspect / error paths, plus helper code to copy).
- Do not list every option up front: put what the task needs in the text and defer the rest to the reference card.

## Generic lesson layout

| Part | Lessons | Content |
| --- | --- | --- |
| 1 Getting started | intro/setup, core concepts | install, the basic flow, a first exercise |
| 2 Core features | one lesson per feature group | build APIs -> inspect APIs |
| 3 Advanced / services | one lesson per sub-package | their build/inspect APIs |
| 4 Errors and edges | failure injection, recovery, error classification | polling, handle/cleanup contracts |
| 5 Practice and reference | cleanup rules, diagnostics, capstone, cards | tips |

## Kinds of edu exercises
1. **Build**: create the subject with an option/function combination and respect the cleanup rules.
2. **Inspect**: implement a function that reads state; the test compares with what the subject itself reports.
3. **Failure**: inject (keep the handle), observe, restore; the test also checks the recovery.
4. **Rejection path**: provoke an invalid combination and check the error class (fast, no heavy resources).

## Quizzes
- About four plausible options; wrong options come from common misconceptions. Give a correct-answer message and a hint
  for a wrong answer.

## Repeat in every lesson
- Poll asynchronously converging state; handle a non-nil handle returned with an error; cleanup order.

## Checklist for a new version / package
- Defaults, signatures and constraints in the option/method tables.
- Measured claims: durations, delays, error classes.
- Supported environments (releases, images, OS) and the required Go version.
- Reference-card tables and the README (version, task counts).
