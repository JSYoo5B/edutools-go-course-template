#!/usr/bin/env python3
"""Verify one task in isolation (EduTools Go layout: task dir = module, tests in test/).

Usage: verify.py <lesson_dir> <task_dir> [--student | --student-test | --run] [-run REGEX] [-v]
--run: theory/choice example: native Go tests in test/, or go run ., under the global lock
--student-test: require the declared assertion failure, or successful execution with different output
Default: solution state, runs `go test ./test/` (output tasks: go run . vs test/output.txt) under a global lock (serializes heavy tests, e.g. ones that start containers).
--student: learner state, only `go vet ./...` (compile check).
"""
import fcntl
import json
import importlib.util
import os
import shutil
import subprocess
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import build  # noqa: E402

ROOT = build.ROOT
WORK_BASE = os.path.join(ROOT, ".work")
LOCK = os.path.join(WORK_BASE, "heavy-tests.lock")


def main():
    args = sys.argv[1:]
    baseline = "--student-test" in args
    student = "--student" in args
    demo = "--run" in args
    args = [a for a in args if a not in ("--student", "--student-test", "--run")]
    run = None
    verbose = False
    if "-run" in args:
        i = args.index("-run")
        run = args[i + 1]
        del args[i:i + 2]
    if "-v" in args:
        verbose = True
        args.remove("-v")
    lesson, task = args[:2]
    spec = importlib.util.spec_from_file_location("lesson_mod", os.path.join(build.SRC, lesson, "lesson.py"))
    mod = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(mod)
    tdef = next(t for t in mod.LESSON["tasks"] if t["dir"] == task)
    os.makedirs(WORK_BASE, exist_ok=True)
    work = os.path.join(WORK_BASE, f"{lesson}-{task}")
    # Protect this task's generated files without blocking unrelated compile checks.
    task_lock = open(work + ".lock", "w")
    fcntl.flock(task_lock, fcntl.LOCK_EX)
    if os.path.exists(work):
        shutil.rmtree(work)
    os.makedirs(work)
    out, stu = os.path.join(work, "sol"), os.path.join(work, "stu")
    build.emit_task(out, stu, "t", os.path.join(build.SRC, lesson, task), tdef, module=f"{lesson}_{task}")
    target = os.path.join(stu if student or baseline else out, "t")
    if student:
        r = subprocess.run(["go", "vet", "./..."], cwd=target, text=True)
        sys.exit(r.returncode)
    # Only actual execution needs the global resource lock.
    lock_file = open(LOCK, "w")
    fcntl.flock(lock_file, fcntl.LOCK_EX)
    if baseline:
        if tdef["type"] == "output":  # the program must run and print something else
            sys.exit(0 if check_output(target) == 1 else "Expected a successful program with different output")
        expected = tdef.get("student_failure")
        if not expected:
            sys.exit("Set student_failure to the intended assertion message in lesson.py")
        result = subprocess.run(["go", "test", "./test/", "-json", "-count=1", "-timeout", "20m"],
                                cwd=target, capture_output=True, text=True)
        print(result.stdout, end="")
        print(result.stderr, end="", file=sys.stderr)
        events = []
        for line in result.stdout.splitlines():
            try:
                events.append(json.loads(line))
            except ValueError:
                pass
        failed = {e.get("Test") for e in events if e.get("Action") == "fail" and e.get("Test")}
        matched = any(e.get("Test") in failed and expected in e.get("Output", "") for e in events)
        broken = any(x in result.stdout + result.stderr for x in ("[build failed]", "[setup failed]", "panic:"))
        if result.returncode != 1 or not matched or broken:
            sys.exit("Learner state did not fail at the expected assertion")
        print("learner state: expected assertion failure confirmed")
        return
    if demo:
        cmd = (["go", "test", "./test/", "-count=1", "-v"]
               if any(name.startswith("test/") and name.endswith("_test.go") for name, _ in tdef["files"])
               else ["go", "run", "."])
        r = subprocess.run(cmd, cwd=target, text=True, timeout=600)
        sys.exit(r.returncode)
    if tdef["type"] == "output":  # run main, compare stdout with test/output.txt (stdin: test/input.txt)
        sys.exit(check_output(target))
    cmd = ["go", "test", "./test/", "-timeout", "20m", "-count=1"]
    if verbose:
        cmd.append("-v")
    if run:
        cmd += ["-run", run]
    r = subprocess.run(cmd, cwd=target, text=True)
    sys.exit(r.returncode)


def check_output(target):
    """Mimic the plugin's output check: trailing white space is ignored."""
    stdin = None
    inp = os.path.join(target, "test", "input.txt")
    if os.path.exists(inp):
        with open(inp, encoding="utf-8") as f:
            stdin = f.read()
    with open(os.path.join(target, "test", "output.txt"), encoding="utf-8") as f:
        want = f.read()
    r = subprocess.run(["go", "run", "."], cwd=target, input=stdin, capture_output=True, text=True, timeout=600)
    if r.returncode != 0:
        print(r.stderr, file=sys.stderr)
        return 2
    if r.stdout.rstrip() != want.rstrip():
        print(f"output differs\n--- want\n{want}--- got\n{r.stdout}", file=sys.stderr)
        return 1
    print("output OK")
    return 0


if __name__ == "__main__":
    main()
