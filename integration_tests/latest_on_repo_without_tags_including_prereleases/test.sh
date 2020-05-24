#!/bin/bash

export CURRENT_TEST="${BASH_SOURCE[0]}"
export TEST_DIR=$(dirname ${CURRENT_TEST})

source ${TEST_DIR}/../init_test_container.sh

set +e
RESULT="$(tc_exec git semver latest --include-pre-releases)"
assertEquals "$?" 0 || exit 1
set -e

assertEquals "0.0.0" "${RESULT}"

echo SUCCESS!
