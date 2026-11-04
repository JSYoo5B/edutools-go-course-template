"""Exercise the verification CLI with small native-test fixtures."""
import fcntl
from pathlib import Path
import subprocess
import tempfile
import unittest

TOOLS = Path(__file__).resolve().parent


class VerifyTests(unittest.TestCase):
    def run_fixture(self, body, mode, *, output=False, hold_global=False):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            src = root / 'src'
            (src / '_root').mkdir(parents=True)
            (src / '_root/go.mod').write_text('module fixture\n\ngo 1.26.0\n')
            task = src / 'lesson/task/test'
            task.mkdir(parents=True)
            (task.parent / 'task.md').write_text('Fixture')
            (task / 'example_test.go').write_text('package test\nimport "testing"\nfunc TestExample(t *testing.T) { '+body+' }\n')
            definition = dict(dir='task', title='Fixture', type='theory',
                              student_failure='expected operation missing',
                              files=[('go.mod', False), ('test/example_test.go', True)])
            if output:
                definition = dict(dir='task', title='Fixture', type='output',
                                  files=[('go.mod', False), ('main.go', True), ('test/output.txt', False)])
                (task / 'example_test.go').unlink()
                (task.parent / 'main.go').write_text('package main\nimport \"fmt\"\nfunc main() { _ = fmt.Print; '+body+' }\n')
                (task / 'output.txt').write_text('wanted\n')
            (src / 'lesson/lesson.py').write_text('LESSON = '+repr(dict(tasks=[definition])))
            code = f'''import sys
sys.path.insert(0, {str(TOOLS)!r})
import verify
verify.build.SRC = {str(src)!r}
verify.WORK_BASE = {str(root / 'work')!r}
verify.LOCK = {str(root / 'work/lock')!r}
sys.argv = ['verify.py', 'lesson', 'task', {mode!r}]
verify.main()
'''
            (root / 'work').mkdir()
            with (root / 'work/lock').open('w') as lock:
                if hold_global:
                    fcntl.flock(lock, fcntl.LOCK_EX)
                return subprocess.run(['python3', '-c', code], capture_output=True, text=True, timeout=30)

    def test_compile_does_not_wait_for_global_execution_lock(self):
        r = self.run_fixture('t.Log("compile only")', '--student', hold_global=True)
        self.assertEqual(r.returncode, 0, r.stdout+r.stderr)

    def test_output_baseline_without_failure_declaration(self):
        r = self.run_fixture('', '--student-test', output=True)
        self.assertEqual(r.returncode, 0, r.stdout+r.stderr)
        for body in ['fmt.Println("wanted")', 'panic("failed")', 'invalid Go code']:
            with self.subTest(body=body):
                r = self.run_fixture(body, '--student-test', output=True)
                self.assertNotEqual(r.returncode, 0, r.stdout+r.stderr)

    def test_native_example(self):
        r = self.run_fixture('t.Log("native test ran")', '--run')
        self.assertEqual(r.returncode, 0, r.stdout+r.stderr)
        self.assertIn('native test ran', r.stdout)

    def test_expected_assertion(self):
        r = self.run_fixture('t.Fatal("expected operation missing")', '--student-test')
        self.assertEqual(r.returncode, 0, r.stdout+r.stderr)

    def test_other_failures_are_rejected(self):
        for body in ['t.Log("unexpected pass")', 't.Fatal("environment unavailable")',
                     'invalid Go code', 'panic("expected operation missing")']:
            with self.subTest(body=body):
                r = self.run_fixture(body, '--student-test')
                self.assertNotEqual(r.returncode, 0, r.stdout+r.stderr)


if __name__ == '__main__':
    unittest.main()
