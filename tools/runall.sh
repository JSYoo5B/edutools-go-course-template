#!/bin/bash
# Run every edu task (solution state) serially and write the results to a log.
# Usage: runall.sh <course_dir> <log_file>
COURSE="${1:?course dir}"; LOG="${2:?log file}"
mkdir -p "$(dirname "$LOG")"
LOG="$(cd "$(dirname "$LOG")" && pwd)/$(basename "$LOG")"   # make the path absolute before cd
cd "$COURSE" || exit 1
: > "$LOG"
for m in $(find . -name go.mod | sort); do
  d=$(dirname "$m")
  [ -d "$d/test" ] || continue            # theory/choice tasks have no test/
  start=$(date +%s)
  out=$(cd "$d" && go test ./test/ -count=1 -timeout 20m 2>&1 | grep -E '^(ok|FAIL|--- FAIL|panic)|_test.go:[0-9]+:' | head -8)
  echo "$d | $(($(date +%s)-start))s | $out" >> "$LOG"
done
echo FULLDONE >> "$LOG"
