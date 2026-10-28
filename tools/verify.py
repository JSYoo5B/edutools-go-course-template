#!/usr/bin/env python3
"""Verify one task in isolation (EduTools Go layout: task dir = module, tests in test/).

Usage: verify.py <lesson_dir> <task_dir> [--student | --run] [-run REGEX] [-v]
--run: theory/choice example: 'go run .' (solution state) under the global lock
Default: solution state, runs `go test ./test/` under a global lock (serializes heavy tests, e.g. ones that start containers).
--student: learner state, only `go vet ./...` (compile check).
"""
import fcntl
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
    student = "--student" in args
    demo = "--run" in args
    args = [a for a in args if a not in ("--student", "--run")]
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
    if os.path.exists(work):
        shutil.rmtree(work)
    os.makedirs(work)
    out, stu = os.path.join(work, "sol"), os.path.join(work, "stu")
    build.emit_task(out, stu, "t", os.path.join(build.SRC, lesson, task), tdef, module=f"{lesson}_{task}")
    target = os.path.join(stu if student else out, "t")
    if student:
        r = subprocess.run(["go", "vet", "./..."], cwd=target, text=True)
        sys.exit(r.returncode)
    if demo:  # run a theory example (go run .) under the global lock, with a time limit
        with open(LOCK, "w") as lf:
            print("waiting for test lock...", flush=True)
            fcntl.flock(lf, fcntl.LOCK_EX)
            print("lock acquired, go run .", flush=True)
            r = subprocess.run(["go", "run", "."], cwd=target, text=True, timeout=600)
            sys.exit(r.returncode)
    with open(LOCK, "w") as lf:
        print("waiting for test lock...", flush=True)
        fcntl.flock(lf, fcntl.LOCK_EX)
        print("lock acquired, running tests", flush=True)
        cmd = ["go", "test", "./test/", "-timeout", "20m", "-count=1"]
        if verbose:
            cmd.append("-v")
        if run:
            cmd += ["-run", run]
        r = subprocess.run(cmd, cwd=target, text=True)
        sys.exit(r.returncode)


if __name__ == "__main__":
    main()
