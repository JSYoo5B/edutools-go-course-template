"""Course structure. For a new course, edit COURSE and SECTIONS only.
Lessons are read from src/lessonNN/lesson.py.

The repository itself (docs, tools, comments) is written in English. The *course content* may be in any
language: set COURSE["language"] to the natural language of the course and write the task.md files,
titles and quiz texts in that language.
"""
import importlib.util
import os

HERE = os.path.dirname(os.path.abspath(__file__))


def load_lesson(name):
    path = os.path.join(HERE, "src", name, "lesson.py")
    spec = importlib.util.spec_from_file_location(f"lesson_{name}", path)
    mod = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(mod)
    return mod.LESSON


COURSE = dict(
    title="Sample Course",
    language="English",  # natural language of the course (EduTools 'language' field), e.g. "Korean", "Japanese"
    summary="""\
A short description of the course: what it teaches, who it is for, what to prepare.""",
    programming_language="Go",
    programming_language_version="1.26",
    additional_files=[],  # every task gets its own generated go.mod, so this is usually empty
    # Texts that the build inserts when a lesson does not provide its own. Translate them for your course language.
    ui=dict(
        quiz_correct="Correct!",
        quiz_incorrect="Not quite. Re-read the explanation for a hint.",
        scratch_comment=[
            "// This is a reading/quiz task. See the 'Task Description' panel.",
            "// This file is a free scratch space and is not used for grading.",
        ],
    ),
)

# Sections -> lesson directory names. Lessons that do not exist yet (no lesson.py) are skipped.
SECTIONS = [
    dict(dir="section01", title="Part 1. Getting started", lessons=["lesson01"]),
]

_resolved = []
for s in SECTIONS:
    ls = [load_lesson(ln) for ln in s["lessons"]
          if os.path.exists(os.path.join(HERE, "src", ln, "lesson.py"))]
    if ls:
        _resolved.append(dict(s, lessons=ls))
SECTIONS = _resolved
