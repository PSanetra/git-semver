#!/bin/bash

export CURRENT_TEST="${BASH_SOURCE[0]}"
export TEST_DIR=$(dirname ${CURRENT_TEST})

source ${TEST_DIR}/../init_test_container.sh

set +e
RESULT="$(tc_exec git semver next)"
assertEquals "$?" 1 || exit 1
set -e

RESULT=$(trim "${RESULT}")

assertEquals "$(echo -e '\e[31mFATA\e[0m[0000] Could not find HEAD: reference not found')" "${RESULT}"

echo SUCCESS!
