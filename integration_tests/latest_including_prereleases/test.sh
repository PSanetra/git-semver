#!/bin/bash

export CURRENT_TEST="${BASH_SOURCE[0]}"
export TEST_DIR=$(dirname ${CURRENT_TEST})

source ${TEST_DIR}/../init_test_container.sh

tc_exec mkdir src
tc_exec touch src/file.go
tc_exec git add -A
tc_exec git commit -m "feat: Add feature"
tc_exec git tag "v1.2.3"
tc_exec touch src/file2.go
tc_exec git add -A
tc_exec git commit -m "feat: Add feature 2"
tc_exec git tag "v1.3.0-beta"

set +e
RESULT="$(tc_exec git semver latest --include-pre-releases)"
assertEquals "$?" 0 || exit 1
set -e

assertEquals "1.3.0-beta" "${RESULT}"

echo SUCCESS!
