#!/bin/bash

export CURRENT_TEST="${BASH_SOURCE[0]}"
export TEST_DIR=$(dirname ${CURRENT_TEST})

source ${TEST_DIR}/../init_test_container.sh

tc_exec mkdir src
tc_exec touch src/file.go
tc_exec git add -A
tc_exec git commit -m "feat: Add feature"
tc_exec git tag -a "v$(git semver next)" -m 'release'
tc_exec touch src/file2.go
tc_exec git add -A
tc_exec git commit -m "chore: Some maintenance"

set +e
RESULT="$(tc_exec git semver next)"
assertEquals "$?" 0 || exit 1
set -e

assertEquals "1.0.0" "${RESULT}"

echo SUCCESS!
