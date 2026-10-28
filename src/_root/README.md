# <Course title> - JetBrains Academy course

<One paragraph: what the course teaches.>

## Contents
(table of parts / lessons)

## What you need
- Go <version> or newer
- (anything else the course requires)

## Open it (GoLand + JetBrains Academy plugin)
- **Learner mode** (the course screen opens directly): open the `course-student/` directory.
- **Author mode** (edit the course): open the `course/` directory. It contains the solutions.

## Solve or verify without the IDE
Each task is one Go module. The grading tests are in the task's `test/` directory.

```bash
cd section01/lesson01/task03
go test ./test/ -v
```

## Layout
```
course/
  course-info.yaml
  .idea/            module definition (course.iml, modules.xml)
  sectionNN/section-info.yaml
    lessonNN/lesson-info.yaml
      taskNN/
        task-info.yaml, task.md
        go.mod, go.sum        (one task = one module)
        *.go                  visible source
        test/task_test.go     hidden grading test (edu)
```

## Verification status
(State honestly what was verified and what was not.)
