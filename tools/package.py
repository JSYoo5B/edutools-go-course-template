#!/usr/bin/env python3
"""Package a validated course tree, then validate a fresh extraction."""
import argparse
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import zipfile


def excluded(relative):
    parts = relative.parts
    if '.DS_Store' in parts:
        return True
    return any(parts[i] == '.idea' and
               (parts[i + 1] == 'workspace.xml' or parts[i + 1].startswith('.author_contents_storage_db'))
               for i in range(len(parts) - 1))


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--author', action='store_true', help='package a Course Creator tree instead of Study')
    parser.add_argument('source', type=Path)
    parser.add_argument('archive', type=Path)
    args = parser.parse_args()
    source, archive = args.source.resolve(), args.archive.resolve()
    flags = [] if args.author else ['-study']
    checker = Path(__file__).resolve().parent / 'yamlcheck' / 'yamlcheck'
    if source == archive or source in archive.parents:
        sys.exit('The archive must be outside the source tree')
    archive.parent.mkdir(parents=True, exist_ok=True)
    with tempfile.TemporaryDirectory() as tmp:
        # Validate the filtered tree: an IDE-opened Study source may contain workspace.xml.
        staged = Path(tmp) / 'staged' / source.name
        staged.mkdir(parents=True)
        for file in sorted(source.rglob('*')):
            relative = file.relative_to(source)
            if file.is_file() and not excluded(relative):
                target = staged / relative
                target.parent.mkdir(parents=True, exist_ok=True)
                shutil.copyfile(file, target)
        subprocess.run([str(checker), *flags, str(staged)], check=True)
        candidate = Path(tmp) / 'course.zip'
        with zipfile.ZipFile(candidate, 'w', zipfile.ZIP_DEFLATED) as z:
            for file in sorted(staged.rglob('*')):
                if file.is_file():
                    z.write(file, Path(source.name) / file.relative_to(staged))
        extracted = Path(tmp) / 'extracted'
        with zipfile.ZipFile(candidate) as z:
            z.extractall(extracted)
        subprocess.run([str(checker), *flags, str(extracted / source.name)], check=True)
        archive.write_bytes(candidate.read_bytes())
    print(f'Packaged and checked fresh extraction: {archive}')


if __name__ == '__main__':
    main()
