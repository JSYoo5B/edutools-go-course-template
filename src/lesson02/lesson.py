# Sample lesson: one task per exercise style (see docs/exercise-styles.md). Delete it in a real course.
LESSON = dict(
    dir="lesson02",
    title="Exercise styles",
    tasks=[
        # output: run main, compare stdout with the hidden test/output.txt (optional hidden test/input.txt = stdin)
        dict(dir="task01", type="output", title="Style: program output",
             files=[("go.mod", False), ("go.sum", False), ("main.go", True), ("test/output.txt", False)]),
        # edu: a main program graded by a hidden test that runs it and checks exit code and decoded JSON
        dict(dir="task02", type="edu", title="Style: program behaviour",
             files=[("go.mod", False), ("go.sum", False), ("main.go", True), ("test/task_test.go", False)]),
        # edu: a function fragment, a visible starter test the learner can run, a hidden behaviour test
        dict(dir="task03", type="edu", title="Style: fragment with a starter test",
             files=[("go.mod", False), ("go.sum", False), ("transfer.go", True),
                    ("test/starter_test.go", True), ("test/task_test.go", False)]),
        # edu: the exercise is the body of a visible test; given assertions grade it, a hidden test guards them
        dict(dir="task04", type="edu", title="Style: code inside a test",
             files=[("go.mod", False), ("go.sum", False), ("hello.go", True),
                    ("test/hello_test.go", True), ("test/task_test.go", False)]),
        # edu: the learner writes a test; a hidden test runs it against broken versions of the code
        dict(dir="task05", type="edu", title="Style: write the test",
             files=[("go.mod", False), ("go.sum", False), ("clamp.go", True),
                    ("test/clamp_test.go", True), ("test/task_test.go", False)]),
    ],
)
