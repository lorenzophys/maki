package main

import (
	"reflect"
	"testing"
)

func TestGetTargetsFromMakeDb(t *testing.T) {
	isopsMakeDb := []byte(`# GNU Make 3.81
# Copyright (C) 2006  Free Software Foundation, Inc.
# This is free software; see the source for copying conditions.
# There is NO warranty; not even for MERCHANTABILITY or FITNESS FOR A
# PARTICULAR PURPOSE.

# This program built for i386-apple-darwin11.3.0

# Make data base, printed on Sun Feb 12 10:59:48 2023

# Variables
# LOTS OF STUFF HERE
# makefile (from 'Makefile', line 1)
MAKEFILE_LIST :=  YOUcanNaMeYourMakeFILEweirdStuff
# default
MAKEFILES := xxx
# automatic
MAIN_PATH = isops
# variable set hash-table stats:
# Load=91/1024=9%, Rehash=0, Collisions=4/114=4%

# Pattern-specific Variable Values

# No pattern-specific variable values.

# Directories


# No files, no impossibilities in 0 directories.

# Implicit Rules

# No implicit rules.

# Files

# Not a target:
::
#  Command-line target.
#  Implicit rule search has been done.
#  File does not exist.
#  File has not been updated.
# variable set hash-table stats:
# Load=0/32=0%, Rehash=0, Collisions=0/0=0%

.PHONY: clean format lint test coverage tox
#  Implicit rule search has not been done.
#  Modification time never checked.
#  File has not been updated.

# Not a target:
.SUFFIXES:
#  Implicit rule search has not been done.
#  Modification time never checked.
#  File has not been updated.

format:
#  Phony target (prerequisite of .PHONY).
#  Implicit rule search has not been done.
#  File does not exist.
#  File has not been updated.
	isort ${MAIN_PATH} ${TESTS_PATH}
	black ${MAIN_PATH} ${TESTS_PATH}
	

test:
#  Phony target (prerequisite of .PHONY).
#  Implicit rule search has not been done.
#  File does not exist.
#  File has not been updated.
	pytest -vvv
	

# Not a target:
.DEFAULT:
#  Implicit rule search has not been done.
#  Modification time never checked.
#  File has not been updated.

clean:
#  Phony target (prerequisite of .PHONY).
#  Implicit rule search has not been done.
#  File does not exist.
#  File has not been updated.
	find . -name '__pycache__' -exec rm -rf {} +
	find . -name '.DS_Store' -exec rm -rf {} +
	find . -name '.pytest_cache' -exec rm -rf {} +
	find . -name '.mypy_cache' -exec rm -rf {} +
	find . -name '.tox' -exec rm -rf {} +
	rm -f .coverage
	

coverage:
#  Phony target (prerequisite of .PHONY).
#  Implicit rule search has not been done.
#  File does not exist.
#  File has not been updated.
	pytest --no-cov-on-fail --cov-report term-missing --cov=${MAIN_PATH} tests/
	

# Not a target:
YOUcanNaMeYourMakeFILEweirdStuff:
#  Implicit rule search has been done.
#  Last modified 2023-02-12 10:49:18
#  File has been updated.
#  Successfully updated.
# variable set hash-table stats:
# Load=0/32=0%, Rehash=0, Collisions=0/0=0%

lint:
#  Phony target (prerequisite of .PHONY).
#  Implicit rule search has not been done.
#  File does not exist.
#  File has not been updated.
	flake8 ${MAIN_PATH} ${TESTS_PATH}
	mypy --no-incremental ${MAIN_PATH} # https://github.com/python/mypy/issues/7276
	pydocstyle ${MAIN_PATH} --add-ignore=${PYDOCSTYLE_IGNORE}
	

tox:
#  Phony target (prerequisite of .PHONY).
#  Implicit rule search has not been done.
#  File does not exist.
#  File has not been updated.
	tox --recreate --parallel
	

# files hash-table stats:
# Load=11/1024=1%, Rehash=0, Collisions=0/34=0%
# VPATH Search Paths

# No 'vpath' search paths.

# No general ('VPATH' variable) search path.

# # of strings in strcache: 1
# # of strcache buffers: 1
# strcache size: total = 4096 / max = 4096 / min = 4096 / avg = 4096
# strcache free: total = 4080 / max = 4080 / min = 4080 / avg = 4080

# Finished Make data base on Sun Feb 12 10:59:48 2023
`)
	got, _ := getTargetsFromMakeDb(isopsMakeDb)
	expected := []string{"clean", "coverage", "format", "lint", "test", "tox"}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %q but got %q", expected, got)
	}

}
