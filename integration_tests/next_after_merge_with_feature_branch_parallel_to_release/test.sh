#!/bin/bash

export CURRENT_TEST="${BASH_SOURCE[0]}"
export TEST_DIR=$(dirname ${CURRENT_TEST})

source ${TEST_DIR}/../init_test_container.sh

# Produce history:
# *   193c028 - Merge branch 'test_branch'
# |\
# | * 17c1f5a - feat: Add feature in branch (test_branch)
# * | 21c43b3 - feat: More features (tag: v1.0.0)
# |/
# *   1688566 - feat: Add feature

tc_exec mkdir src
tc_exec touch src/a.go
tc_exec git add -A
tc_exec git commit -m "feat: Add feature"
tc_exec git checkout -b test_branch
tc_exec touch src/b.go
tc_exec git add -A
tc_exec git commit -m "feat: Add feature in branch"
tc_exec git checkout master
tc_exec touch src/c.go
tc_exec git add -A
tc_exec git commit -m "feat: More features"
tc_exec git tag "v$(tc_exec git semver next)"
tc_exec git merge test_branch --no-edit
tc_pretty_git_log

set +e
RESULT="$(tc_exec git semver next)"
assertEquals "$?" 0 || exit 1
set -e

assertEquals "1.1.0" "${RESULT}"

echo SUCCESS!
