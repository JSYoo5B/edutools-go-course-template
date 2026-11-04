# EduTools Go course template - common commands. See AGENTS.md for the rules.
.PHONY: test-tools help build check study vet verify run-demo run-all zip zip-author clean

OUT      ?= build/course
STUDENT  ?= build/student
ZIP      ?= course.zip
AUTHOR_ZIP ?= course-author.zip

help:
	@echo "make build                     src/ -> $(OUT) (+ learner state $(STUDENT))"
	@echo "make check                     build + structure check + go vet on every module"
	@echo "make verify L=lesson01 T=task03   one task: learner compile + expected failure + solution test"
	@echo "make run-demo L=lesson01 T=task01 run a theory example (native tests or go run .)"
	@echo "make run-all                   run every edu task serially (log: build/run-all.log)"
	@echo "make study                     author tree -> learner (Study-mode) tree (build/course-study) + check"
	@echo "make zip                       Study tree into $(ZIP)"
	@echo "make zip-author                Course Creator tree into $(AUTHOR_ZIP)"

build:
	python3 tools/build.py $(OUT) --student $(STUDENT)

tools/yamlcheck/yamlcheck: tools/yamlcheck/main.go
	cd tools/yamlcheck && go build -o yamlcheck .

check: build test-tools tools/yamlcheck/yamlcheck
	tools/yamlcheck/yamlcheck $(OUT)
	@fail=0; for m in $$(find $(OUT) $(STUDENT) -name go.mod | sort); do \
	  d=$$(dirname $$m); (cd $$d && go vet ./... >/dev/null 2>&1) || { echo "VET FAIL $$d"; fail=1; }; \
	done; test $$fail = 0 && echo "vet OK"

verify:
	python3 tools/verify.py $(L) $(T) --student
	python3 tools/verify.py $(L) $(T) --student-test
	python3 tools/verify.py $(L) $(T) -v

run-demo:
	python3 tools/verify.py $(L) $(T) --run

tools/to_study/to_study: $(wildcard tools/to_study/*.go)
	cd tools/to_study && go build -o to_study .

STUDY_OUT ?= build/course-study

study: build test-tools tools/to_study/to_study tools/yamlcheck/yamlcheck
	tools/to_study/to_study $(OUT) $(STUDY_OUT)
	tools/yamlcheck/yamlcheck -study $(STUDY_OUT)
	@fail=0; for m in $$(find $(STUDY_OUT) -name go.mod | sort); do \
	  d=$$(dirname $$m); (cd $$d && go vet ./... >/dev/null 2>&1) || { echo "VET FAIL $$d"; fail=1; }; \
	done; test $$fail = 0 && echo "learner tree vet OK"

run-all: build
	tools/runall.sh $(OUT) build/run-all.log
	@cat build/run-all.log

test-tools:
	python3 -m unittest discover -s tools -p 'test_*.py'
	cd tools/to_study && go test ./... -count=1
	cd tools/yamlcheck && go test ./... -count=1

zip: study
	python3 tools/package.py "$(STUDY_OUT)" "$(ZIP)"

zip-author: check
	python3 tools/package.py --author "$(OUT)" "$(AUTHOR_ZIP)"

clean:
	rm -rf build .work
