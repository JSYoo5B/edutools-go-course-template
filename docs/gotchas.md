# EduTools / GoLand pitfalls (in the order they were hit)

Each item: symptom -> cause -> fix / how to check.

## 1. Check says "No tests have run"
- Cause: the Go checker (`GoEduTaskChecker`) runs only the tests in the **`test/` directory of the task**. A test in the
  task root, or a single `go.mod` at the course root, does not work.
- Fix: one task = one Go module: `go.mod` (module e.g. `lesson03_task02`) + source in the root + `test/task_test.go`
  (`package test`, importing the source as `task "<module>"`). Plugin source:
  `educational-plugin/intellij-plugin/Edu-Go/src/com/jetbrains/edu/go/` (GoConfigurator: test file `task_test.go`, test directory `test`).

## 2. Check says "Failed to launch checking"
- Cause: a hidden file in `task-info.yaml` had `editable: false`. Check fills the test file through
  `createTests -> replaceFileText` and hit a `ReadOnlyModificationException`.
- Fix: never use `editable: false`. Diagnose in `~/Library/Logs/JetBrains/GoLand<ver>/idea.log`: look for the
  `Local check of task` exception and its stack (`CheckAction$StudyCheckTask`).

## 3. The description of a theory/quiz task does not show
- Cause: with no visible file (`files: []`) nothing opens in the editor, so the Task Description panel does not appear.
  The official *Tour of Go* gives every task, theory included, a visible `main.go` and a hidden `go.mod`.
- Fix: every task gets a visible file. We use an example `main.go` that shows the concept.

## 4. The top directory shows as excluded / the project has no module
- Cause: the course `.idea/` has no `modules.xml` and no `<name>.iml` (`get_project_modules` -> `[]`).
- Fix: `.idea/course.iml` (`<module type="WEB_MODULE">`, `<component name="Go" enabled="true"/>`,
  `<content url="file://$MODULE_DIR$"/>`; for an iml inside `.idea`, `$MODULE_DIR$` is the project root),
  `modules.xml`, `externalDependencies.xml` (Go plugin). Templates: `src/_root/idea/`. The IDE does not reload
  automatically: close and re-open the project.

## 5. Placeholder positions are off
- offset/length are in **UTF-16 code units** (BMP characters count 1, emoji 2). Use the marker generator.
  `yamlcheck` checks ranges and empty regions.

## 6. Heavy tests run in parallel and time out
- Starting processes or containers in several tests at once exhausts memory. `tools/verify.py` serialises them with
  an `fcntl.flock` global lock. Tell agents not to bypass it with a direct `go test`.

## 7. Preview versus opening the folder
- Preview copies the course to `/private/var/folders/.../T/course_previewN/`: after editing the source, start a new
  preview. Building over an open project can delete `.idea/workspace.xml`.

## 8. Miscellaneous
- BSD/macOS vs GNU tools: if `sed -i ''` fails, replace with Python.
- In zsh a glob such as `*_test.go` with no match is an error: use `find`.
- `go mod tidy` can be slow (GOPROXY=direct, ...): create a seed go.sum once and copy it into every task.
- A dependency's `toolchain` line does not change your module's toolchain (only keep the `go` line right).
- The official *Tour of Go* uses `output` tasks (output comparison). The canonical `edu` + Go test layout is in the
  plugin's file templates (`task_test.go.ft`, `edu_task.go.ft`, `go.mod.ft`).

## 9. Author-provided multilingual content is not displayed (checked by experiment)
- `task_en.md` / `task_ko.md` next to `task.md` did not show in GoLand + JetBrains Academy. The plugin model has one
  description per task; `task_<lang>.md` looks like a cache of the machine-translation feature.
- Title translations are stored under the item IDs of marketplace courses, so they are unusable in a local course.
- For several languages, **build and publish one course per language**.

## 10. Opening a course shows the author (Course Creator) screen, not the learner screen
- Cause (from the plugin source, `YamlDeserializer.deserializeCourse`): `mode:` present in `course-info.yaml` -> learner
  mode, absent -> author mode. Our build writes no `mode`, so it always opens as author. A learner view exists only in
  Preview or by importing a course archive.
- Adding `mode: Study` by hand is not enough. The author tree files contain the **solution** and the placeholder
  offset/length refer to the solution. A learner-mode tree has the **placeholder text** in the files and offsets that refer
  to that text (`initial_state`, `status` included).
- Fix: `tools/to_study` (`make study`) replaces the solution regions with the placeholder text, recomputes
  offset/length/`initial_state`, adds `status`/`record` (and `selected_options` for quizzes) and `mode: Study`. It reads and
  writes YAML with `yaml.v3`, because the IDE rewrites an author tree's YAML in its own style, so line-based parsing fails.
  Validate with `yamlcheck -study`. Known limit: no encrypted solution, so "peek solution" is unavailable.
- Reference: the official *Tour of Go* is a learner tree (`mode: Study`, `status`, `record`).

## 11. The learner tree opens but does not jump to the first task
- Cause: `EduProjectActivity` jumps to the first task **only when no file is open**
  (`!fileEditorManager.hasOpenFiles() -> NavigationUtils.openFirstTask`). With a root `README.md` the IDE opens it first
  (`idea.log`: `fileOpened README.md`, then `Course opened`), so the jump is skipped.
- Fix: **no root `README.md` in the learner tree** (also removed from `additional_files`), no saved
  `.idea/workspace.xml`, no `frameworkLessonHistoryCC`. `tools/to_study` does all three. Diagnose: the file opened just
  before `Course opened: <name>` in `idea.log`.
- Note: when the IDE opens a tree it may add `additional_files: README.md` back into `course-info.yaml` (the file stays
  absent; harmless). Run the validators on freshly generated output, not on a tree the IDE has touched.
