# EduTools Go course template

A template repository for building and maintaining **Go learning courses** in the
[JetBrains Academy (EduTools)](https://plugins.jetbrains.com/plugin/10081-jetbrains-academy) format.
It generates the course tree from `src/`, verifies the tasks, converts the result into a tree that opens directly in
learner mode, and packages it. Instructions for agents: [AGENTS.md](AGENTS.md).

The repository is written in English; the course you build can be in any language.

## Quick start

```bash
git clone <this repo> my-course && cd my-course
make check                       # build the sample course + structure check + go vet of every module
make verify L=lesson01 T=task03  # verify one sample exercise (learner state compiles, solution passes)
make run-demo L=lesson01 T=task01
make study                       # learner (Study-mode) tree in build/course-study, validated
make zip                         # course.zip
```

Requirements: Go 1.26+, Python 3. (`tools/yamlcheck` and `tools/to_study` are built with `go build`; they need
`gopkg.in/yaml.v3` in the module cache or network access.)

## Turn it into your course

1. Fill `structure.py` (title, summary, **course language**, default UI texts, sections) and `src/_root/README.md`.
2. Add your subject library to `src/_root/go.mod` (`go get`, `go mod tidy`, keep go.sum).
3. Write lessons under `src/` using the sample `src/lesson01` as the model; rules in
   [docs/authoring-guide.md](docs/authoring-guide.md).
4. Pitfalls: [docs/gotchas.md](docs/gotchas.md). Design: [docs/course-design.md](docs/course-design.md).

## Two trees

| Output | Opens in | Contains |
| --- | --- | --- |
| `build/course` (author tree) | Course Creator mode | solutions, placeholders by offset |
| `build/course-study` (`make study`) | learner mode, jumps to the first task | placeholder text in the files, `status`/`record`, `mode: Study`, no root README |

## Layout

`AGENTS.md` (instructions) | `structure.py` (course metadata) | `src/` (sources) | `tools/` (build, verify, validate, convert)
| `docs/` | `Makefile`

go.mod: the source has one `src/_root/go.mod`; the build generates a `go.mod`/`go.sum` per task.
