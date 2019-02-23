#!/usr/bin/env bash

export TESTS_DIR=$(dirname ${BASH_SOURCE[0]})

shopt -s globstar

SUCCESS_COUNT=0
FAILED_COUNT=0

for test in ${TESTS_DIR}/**/test.sh; do
    echo "run ${test}"

    TEST_OUTPUT=$(${test})

    TEST_SUCCEEDED="$?"

    if [[ ${TEST_SUCCEEDED} != 0 ]]; then
        echo "$TEST_OUTPUT"
        echo -e "\e[0;31m${test} FAILED\e[0m"
        FAILED_COUNT=$((FAILED_COUNT+1))
    else
        SUCCESS_COUNT=$((SUCCESS_COUNT+1))
    fi

    TEST_OUTPUT=''
done

echo -e "\e[0;32m${SUCCESS_COUNT} tests succeeded\e[0m"
[[ ${FAILED_COUNT} != 0 ]] && echo -e "\e[0;31m${FAILED_COUNT} tests failed!\e[0m"
