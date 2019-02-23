#!/bin/bash

# External environment variables:
# GIT_SEMVER_BINARY: Path to git-semver binary to use in test
# WAIT_AFTER_TEST: Set this variable to leave the container running after executing the test until a key is pressed.
# BIND_WORKDIR: Set this variable to mount the container working directory to ${TEST_DIR}/testworkdir with bindfs (unmount with sudo umount "${TEST_DIR}/testworkdir")

set -e

export WORKDIR=/root/workdir
export TEST_CONTAINER_NAME=git_semver_test

[[ -z "${GIT_SEMVER_BINARY}" ]] && echo "GIT_SEMVER_BINARY not defined" && exit 1

function cleanup {

  if [[ ! -z "$WAIT_AFTER_TEST" ]]; then
    echo -e "\e[0;33mNO_CLEANUP: Leaving ${TEST_CONTAINER_NAME} running!\e[0m"

    if [[ ! -z "$BIND_WORKDIR" ]]; then

      TEST_WORKDIR_MOUNT_POINT="${TEST_DIR}/.testworkdir"

      echo "Mounting container working directory at ${TEST_WORKDIR_MOUNT_POINT}"
      mkdir -p "${TEST_WORKDIR_MOUNT_POINT}"

      CURRENT_USER="$(whoami | tr -d '\n')"
      sudo bindfs --map=root/${CURRENT_USER}:@root/@${CURRENT_USER} "$(docker inspect --format '{{.GraphDriver.Data.MergedDir}}' ${TEST_CONTAINER_NAME} | tr -d '\n')${WORKDIR}" "${TEST_WORKDIR_MOUNT_POINT}"

      read -n1 -r -p "Press any key to stop the container..." key

      sudo umount "${TEST_WORKDIR_MOUNT_POINT}"
    else
      read -n1 -r -p "Press any key to stop the container..." key
    fi
  fi

  docker stop "${TEST_CONTAINER_NAME}"
}

function assertEquals {
  [[ "$1" = "$2" ]] || (echo "'$1' is not equal to '$2'" && echo "FAIL!" && exit 1)
}

function tc_exec {
  docker exec -ti ${TEST_CONTAINER_NAME} "$@"
}

function tc_cp {
  docker cp ${TEST_DIR}/${1} ${TEST_CONTAINER_NAME}:${WORKDIR}/${2}
}

function tc_pretty_git_log {
  tc_exec git --no-pager log --graph --abbrev-commit --decorate --format=format:'%C(bold blue)%h%C(reset) - %C(white)%s %C(bold yellow)%d%C(reset)' --all
  echo
}

function trim {
   echo "$1" | sed -e 's/^[[:space:]]*//' | sed -e 's/[[:space:]]*$//'
}

docker run --rm -tid --name=${TEST_CONTAINER_NAME} --workdir=${WORKDIR} --entrypoint=/bin/bash bitnami/git:2.22.0

trap cleanup EXIT

docker cp ${GIT_SEMVER_BINARY} ${TEST_CONTAINER_NAME}:/bin/git-semver

tc_exec git init .

tc_exec git config user.email test@example.com
tc_exec git config user.name testuser
