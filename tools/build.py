#!/usr/bin/env python3
"""Course generator: src/ (authoring format with //@@ markers) -> EduTools course tree.

Usage: build.py <out_dir> [--student <student_out_dir>]
"""
import json
import os
import re
import shutil
import sys

HERE = os.path.dirname(os.path.abspath(__file__))      # tools/
ROOT = os.path.dirname(HERE)                            # repo root
SRC = os.path.join(ROOT, "src")
STRICT = False
sys.path.insert(0, ROOT)
import structure  # noqa: E402

# texts inserted when a lesson does not provide its own (see structure.COURSE['ui'])
_UI_DEFAULTS = dict(
    quiz_correct="Correct!",
    quiz_incorrect="Not quite. Re-read the explanation for a hint.",
    scratch_comment=[
        "// This is a reading/quiz task. See the 'Task Description' panel.",
        "// This file is a free scratch space and is not used for grading.",
    ],
)
UI = {**_UI_DEFAULTS, **structure.COURSE.get("ui", {})}


def u16len(s):
    return len(s.encode("utf-16-le")) // 2


def yq(s):
    """YAML double-quoted scalar."""
    return json.dumps(s, ensure_ascii=False)


def process_go(text):
    """Strip //@@ markers. Returns (clean_text, placeholders, student_text)."""
    lines = text.split("\n")
    out = []  # list of (line_text)
    phs = []  # (start_line_index_in_out, end_line_index_in_out, indent, placeholder_text)
    i = 0
    while i < len(lines):
        ln = lines[i]
        if ln.strip() == "//@@begin":
            i += 1
            ph_lines = []
            while lines[i].strip().startswith("//@@ph"):
                raw = lines[i].strip()
                ph_lines.append(raw[len("//@@ph "):] if len(raw) > len("//@@ph") else "")
                i += 1
            start = len(out)
            while lines[i].strip() != "//@@end":
                out.append(lines[i])
                i += 1
            end = len(out) - 1
            i += 1  # skip //@@end
            first = out[start]
            indent = first[: len(first) - len(first.lstrip("\t "))]
            phs.append((start, end, indent, "\n".join(ph_lines)))
            continue
        out.append(ln)
        i += 1
    clean = "\n".join(out)
    # offsets
    placeholders = []
    cum = []  # char offset (utf16) of each line start
    pos = 0
    for l in out:
        cum.append(pos)
        pos += u16len(l) + 1
    for (start, end, indent, ph_text) in phs:
        off = cum[start] + u16len(indent)
        last_end = cum[end] + u16len(out[end])
        length = last_end - off
        placeholders.append(
            dict(offset=off, length=length, placeholder_text=indent_cont(ph_text, indent))
        )
    # student text
    stu = []
    for idx, l in enumerate(out):
        stu.append(l)
    student = clean
    for (start, end, indent, ph_text) in reversed(phs):
        before = "\n".join(out[:start])
        first = out[start]
        pre_first = first[: len(indent)]
        after_lines = out[end + 1:]
        mid = pre_first + indent_cont(ph_text, indent)
        parts = []
        if start > 0:
            parts.append(before)
        parts.append(mid)
        parts.extend(after_lines)
        out = (before.split("\n") if start > 0 else []) + mid.split("\n") + after_lines
    student = "\n".join(out)
    return clean, placeholders, student


def indent_cont(ph_text, indent):
    ls = ph_text.split("\n")
    return "\n".join([ls[0]] + [indent + l if l else l for l in ls[1:]])


def write(path, content):
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, "w", encoding="utf-8") as f:
        f.write(content)


def main():
    global STRICT
    STRICT = True
    out_dir = sys.argv[1]
    student_dir = None
    if "--student" in sys.argv:
        student_dir = sys.argv[sys.argv.index("--student") + 1]
    for d in (out_dir, student_dir):
        if d and os.path.exists(d):
            for _ in range(5):
                shutil.rmtree(d, ignore_errors=True)
                if not os.path.exists(d) or os.listdir(d) == [".DS_Store"]:
                    break
    os.makedirs(out_dir, exist_ok=True)
    if student_dir:
        os.makedirs(student_dir, exist_ok=True)

    # course-info.yaml
    c = structure.COURSE
    ci = []
    ci.append(f"title: {yq(c['title'])}")
    ci.append(f"language: {c['language']}")
    ci.append("summary: |-")
    for l in c["summary"].strip().split("\n"):
        ci.append("  " + l if l else "")
    ci.append(f"programming_language: {c['programming_language']}")
    ci.append(f"programming_language_version: {c['programming_language_version']}")
    ci.append("content:")
    for s in structure.SECTIONS:
        ci.append(f"- {s['dir']}")
    if c["additional_files"]:
        ci.append("additional_files:")
        for name, vis in c["additional_files"]:
            ci.append(f"- name: {name}")
            ci.append(f"  visible: {'true' if vis else 'false'}")
    write(os.path.join(out_dir, "course-info.yaml"), "\n".join(ci) + "\n")
    if student_dir:
        write(os.path.join(student_dir, "course-info.yaml"), "\n".join(ci) + "\n")

    for extra in c["additional_files"]:
        name = extra[0]
        p = os.path.join(SRC, "_root", name)
        if os.path.exists(p):
            shutil.copy(p, os.path.join(out_dir, name))
            if student_dir:
                shutil.copy(p, os.path.join(student_dir, name))

    idea_src = os.path.join(SRC, "_root", "idea")
    if os.path.isdir(idea_src):
        # module definition so that GoLand treats the course dir as a Go module (else everything shows as excluded)
        shutil.copytree(idea_src, os.path.join(out_dir, ".idea"), dirs_exist_ok=True)
    rd = os.path.join(SRC, "_root", "README.md")
    if os.path.exists(rd):
        shutil.copy(rd, os.path.join(out_dir, "README.md"))

    for s in structure.SECTIONS:
        si = [f"custom_name: {yq(s['title'])}", "content:"]
        for l in s["lessons"]:
            si.append(f"- {l['dir']}")
        write(os.path.join(out_dir, s["dir"], "section-info.yaml"), "\n".join(si) + "\n")
        if student_dir:
            write(os.path.join(student_dir, s["dir"], "section-info.yaml"), "\n".join(si) + "\n")
        for l in s["lessons"]:
            li = [f"custom_name: {yq(l['title'])}", "content:"]
            for t in l["tasks"]:
                li.append(f"- {t['dir']}")
            ldir = os.path.join(s["dir"], l["dir"])
            write(os.path.join(out_dir, ldir, "lesson-info.yaml"), "\n".join(li) + "\n")
            if student_dir:
                write(os.path.join(student_dir, ldir, "lesson-info.yaml"), "\n".join(li) + "\n")
            for t in l["tasks"]:
                build_task(out_dir, student_dir, s, l, t)


def build_task(out_dir, student_dir, s, l, t):
    tdir = os.path.join(s["dir"], l["dir"], t["dir"])
    src_dir = os.path.join(SRC, l["dir"], t["dir"])
    emit_task(out_dir, student_dir, tdir, src_dir, t, module=f"{l['dir']}_{t['dir']}")


def gen_gomod(out_dir, student_dir, tdir, module):
    """Per-task go.mod/go.sum (EduTools Go layout: each task is its own module)."""
    root = os.path.join(SRC, "_root")
    with open(os.path.join(root, "go.mod"), encoding="utf-8") as f:
        mod = f.read().replace("module coursemod", f"module {module}", 1)
    d = os.path.join(out_dir, tdir)
    write(os.path.join(d, "go.mod"), mod)
    if os.path.exists(os.path.join(root, "go.sum")):
        shutil.copy(os.path.join(root, "go.sum"), os.path.join(d, "go.sum"))
    import subprocess
    r = subprocess.run(["go", "mod", "tidy"], cwd=d, capture_output=True, text=True)
    if r.returncode != 0:
        sys.exit(f"go mod tidy failed in {d}: {r.stdout} {r.stderr}")
    sumf = os.path.join(d, "go.sum")
    if not os.path.exists(sumf):  # without dependencies tidy creates no go.sum -> write an empty file
        open(sumf, "w").close()
    if student_dir:
        for n in ("go.mod", "go.sum"):
            if os.path.exists(os.path.join(d, n)):
                shutil.copy(os.path.join(d, n), os.path.join(student_dir, tdir, n))


def emit_task(out_dir, student_dir, tdir, src_dir, t, module=None):
    if not os.path.isdir(src_dir):
        sys.exit(f"missing source dir {src_dir}")
    ti = []
    typ = t["type"]
    ti.append(f"type: {typ}")
    ti.append(f"custom_name: {yq(t['title'])}")
    if typ == "choice":
        ti.append(f"is_multiple_choice: {'true' if t.get('multiple') else 'false'}")
        ti.append("options:")
        for text, ok in t["options"]:
            ti.append(f"- text: {yq(text)}")
            ti.append(f"  is_correct: {'true' if ok else 'false'}")
        ti.append(f"message_correct: {yq(t.get('msg_ok', UI['quiz_correct']))}")
        ti.append(f"message_incorrect: {yq(t.get('msg_bad', UI['quiz_incorrect']))}")
        if t.get("header"):
            ti.append(f"quiz_header: {yq(t['header'])}")
    files = t.get("files", [])
    if typ in ("theory", "choice") and not files:
        # EduTools shows the task description only when a task file is open in the editor
        # (the reference Tour of Go course gives every task a visible main.go + hidden go.mod).
        mod = f"{module or 'task'}"
        scratch = (
            "\n".join(UI["scratch_comment"]) + "\n"
            "package main\n\nfunc main() {}\n"
        ).replace("\\n", "\n")
        write(os.path.join(out_dir, tdir, "main.go"), scratch)
        write(os.path.join(out_dir, tdir, "go.mod"), f"module {mod}\n\ngo 1.26.0\n")
        if student_dir:
            write(os.path.join(student_dir, tdir, "main.go"), scratch)
            write(os.path.join(student_dir, tdir, "go.mod"), f"module {mod}\n\ngo 1.26.0\n")
        ti.append("files:")
        ti.append("- name: main.go")
        ti.append("  visible: true")
        ti.append("- name: go.mod")
        ti.append("  visible: false")
    elif files:
        ti.append("files:")
    else:
        ti.append("files: []")
    gen_mod = False
    for fname, visible in files:
        src_file = os.path.join(src_dir, fname)
        if fname in ("go.mod", "go.sum") and not os.path.exists(src_file):
            gen_mod = True
            raw = clean = student = None
            phs = []
            ti.append(f"- name: {fname}")
            ti.append(f"  visible: {'true' if visible else 'false'}")
            continue
        with open(src_file, encoding="utf-8") as f:
            raw = f.read()
        if fname.endswith(".go"):
            clean, phs, student = process_go(raw)
        else:
            clean, phs, student = raw, [], raw
        write(os.path.join(out_dir, tdir, fname), clean)
        if student_dir:
            write(os.path.join(student_dir, tdir, fname), student)
        ti.append(f"- name: {fname}")
        ti.append(f"  visible: {'true' if visible else 'false'}")
        if phs:
            ti.append("  placeholders:")
            for p in phs:
                ti.append(f"  - offset: {p['offset']}")
                ti.append(f"    length: {p['length']}")
                ti.append(f"    placeholder_text: {yq(p['placeholder_text'])}")
    if gen_mod:
        gen_gomod(out_dir, student_dir, tdir, module)
    write(os.path.join(out_dir, tdir, "task-info.yaml"), "\n".join(ti) + "\n")
    if student_dir:
        write(os.path.join(student_dir, tdir, "task-info.yaml"), "\n".join(ti) + "\n")
    md = os.path.join(src_dir, "task.md")
    if not os.path.exists(md):
        if STRICT:
            sys.exit(f"missing task.md in {src_dir}")
        print(f"warning: missing task.md in {src_dir}", file=sys.stderr)
        return
    shutil.copy(md, os.path.join(out_dir, tdir, "task.md"))
    if student_dir:
        shutil.copy(md, os.path.join(student_dir, tdir, "task.md"))


if __name__ == "__main__":
    main()
