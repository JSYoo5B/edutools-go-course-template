# EduTools Go course template - common commands. See AGENTS.md for the rules.
.PHONY: help build check study vet verify run-demo run-all zip clean

OUT      ?= build/course
STUDENT  ?= build/student
ZIP      ?= course.zip

help:
	@echo "make build                     src/ -> $(OUT) (+ learner state $(STUDENT))"
	@echo "make check                     build + structure check + go vet on every module"
	@echo "make verify L=lesson01 T=task03   one task: learner-state compile + solution test"
	@echo "make run-demo L=lesson01 T=task01 run a theory example (go run .)"
	@echo "make run-all                   run every edu task serially (log: build/run-all.log)"
	@echo "make study                     author tree -> learner (Study-mode) tree (build/course-study) + check"
	@echo "make zip                       package into $(ZIP)"

build:
	python3 tools/build.py $(OUT) --student $(STUDENT)

tools/yamlcheck/yamlcheck: tools/yamlcheck/main.go
	cd tools/yamlcheck && go build -o yamlcheck .

check: build tools/yamlcheck/yamlcheck
	tools/yamlcheck/yamlcheck $(OUT)
	@fail=0; for m in $$(find $(OUT) $(STUDENT) -name go.mod | sort); do \
	  d=$$(dirname $$m); (cd $$d && go vet ./... >/dev/null 2>&1) || { echo "VET FAIL $$d"; fail=1; }; \
	done; test $$fail = 0 && echo "vet OK"

verify:
	python3 tools/verify.py $(L) $(T) --student
	python3 tools/verify.py $(L) $(T) -v

run-demo:
	python3 tools/verify.py $(L) $(T) --run

tools/to_study/to_study: tools/to_study/main.go
	cd tools/to_study && go build -o to_study .

STUDY_OUT ?= build/course-study

study: build tools/to_study/to_study tools/yamlcheck/yamlcheck
	tools/to_study/to_study $(OUT) $(STUDY_OUT)
	tools/yamlcheck/yamlcheck -study $(STUDY_OUT)
	@fail=0; for m in $$(find $(STUDY_OUT) -name go.mod | sort); do \
	  d=$$(dirname $$m); (cd $$d && go vet ./... >/dev/null 2>&1) || { echo "VET FAIL $$d"; fail=1; }; \
	done; test $$fail = 0 && echo "learner tree vet OK"

run-all: build
	tools/runall.sh $(OUT) build/run-all.log
	@cat build/run-all.log

zip: build
	find $(OUT) -name .DS_Store -delete
	cd build && rm -f ../$(ZIP) && zip -rq ../$(ZIP) course -x '*.DS_Store' 'course/.idea/.author_contents_storage_db*' 'course/.idea/workspace.xml'
	@ls -la $(ZIP)

clean:
	rm -rf build .work
