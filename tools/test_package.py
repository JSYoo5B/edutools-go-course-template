"""Run the package CLI on course trees containing saved IDE state."""
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import unittest
import zipfile

TOOLS = Path(__file__).resolve().parent


class PackageTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.temp = tempfile.TemporaryDirectory()
        cls.addClassCleanup(cls.temp.cleanup)
        cls.root = Path(cls.temp.name)
        cls.tools = cls.root / 'tools'
        (cls.tools / 'yamlcheck').mkdir(parents=True)
        shutil.copyfile(TOOLS / 'package.py', cls.tools / 'package.py')
        subprocess.run(['go', 'build', '-o', str(cls.tools / 'yamlcheck/yamlcheck'), '.'],
                       cwd=TOOLS / 'yamlcheck', check=True, capture_output=True)
        converter = cls.tools / 'to_study'
        subprocess.run(['go', 'build', '-o', str(converter), '.'],
                       cwd=TOOLS / 'to_study', check=True, capture_output=True)
        cls.author = cls.root / 'author'
        cls.study = cls.root / 'study'
        subprocess.run([sys.executable, str(TOOLS / 'build.py'), str(cls.author)],
                       check=True, capture_output=True)
        subprocess.run([str(converter), str(cls.author), str(cls.study)], check=True, capture_output=True)

    def test_direct_packaging_excludes_ide_state_without_changing_source(self):
        for author in (False, True):
            with self.subTest(author=author), tempfile.TemporaryDirectory() as tmp:
                source = Path(tmp) / 'opened-course'
                shutil.copytree(self.author if author else self.study, source)
                task = next(source.rglob('task-info.yaml')).parent.relative_to(source)
                unwanted = [Path('.DS_Store'), task / '.DS_Store',
                            Path('.idea/workspace.xml'), Path('.idea/.author_contents_storage_db'),
                            Path('.idea/.author_contents_storage_db-wal'),
                            Path('.idea/.author_contents_storage_db-shm'),
                            Path('.idea/.author_contents_storage_db.files/data')]
                # Similar names outside .idea and required module metadata must survive.
                retained = [task / 'workspace.xml', task / '.author_contents_storage_db-note']
                for relative in unwanted + retained:
                    path = source / relative
                    path.parent.mkdir(parents=True, exist_ok=True)
                    path.write_text('saved IDE state or course content')
                before = {p.relative_to(source): p.read_bytes() for p in source.rglob('*') if p.is_file()}
                archive = Path(tmp) / 'course.zip'
                result = subprocess.run([sys.executable, str(self.tools / 'package.py'),
                                         *(['--author'] if author else []), str(source), str(archive)],
                                        capture_output=True, text=True)
                self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
                with zipfile.ZipFile(archive) as z:
                    actual = {Path(name).relative_to(source.name): z.read(name) for name in z.namelist()}
                self.assertEqual(actual, {p: data for p, data in before.items() if p not in unwanted})
                self.assertIn(Path('.idea/modules.xml'), actual)
                self.assertEqual(before, {p.relative_to(source): p.read_bytes()
                                          for p in source.rglob('*') if p.is_file()})

    def test_invalid_course_is_still_rejected(self):
        with tempfile.TemporaryDirectory() as tmp:
            source = Path(tmp) / 'invalid'
            shutil.copytree(self.study, source)
            (source / 'course-info.yaml').write_text('mode: Study\n')
            archive = Path(tmp) / 'course.zip'
            result = subprocess.run([sys.executable, str(self.tools / 'package.py'), str(source), str(archive)],
                                    capture_output=True, text=True)
            self.assertNotEqual(result.returncode, 0)
            self.assertFalse(archive.exists())


if __name__ == '__main__':
    unittest.main()
