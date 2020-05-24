#!/bin/bash

export CURRENT_TEST="${BASH_SOURCE[0]}"
export TEST_DIR=$(dirname ${CURRENT_TEST})

source ${TEST_DIR}/../init_test_container.sh

tc_exec git checkout -b master
tc_exec mkdir src
tc_exec touch src/file.go
tc_exec git add -A
tc_exec git commit -m "feat: Master commit"
tc_exec git checkout -b v1
tc_exec touch src/file2.go
tc_exec git add -A
tc_exec git commit -m "feat: v1 commit"
tc_exec git tag "v1.0.0"
tc_exec git checkout -b v2
tc_exec touch src/file3.go
tc_exec git add -A
tc_exec git commit -m "feat: v2 commit"
tc_exec git tag "v2.0.0"

set +e
RESULT="$(tc_exec git semver latest --major-version=1)"
assertEquals "$?" 0 || ( echo "${RESULT}" && exit 1 )
set -e

assertEquals "1.0.0" "${RESULT}"

echo SUCCESS!
