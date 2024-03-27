# git-semver
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org) [![Docker Cloud Automated build](https://img.shields.io/docker/cloud/automated/psanetra/git-semver)](https://hub.docker.com/r/psanetra/git-semver) [![Docker Image Pulls](https://img.shields.io/docker/pulls/psanetra/git-semver)](https://hub.docker.com/r/psanetra/git-semver)

git-semver is a command line tool to calculate [semantic versions](https://semver.org/spec/v2.0.0.html) based on the git history and tags of a repository.

git-semver assumes that the commit messages in the git history are wellformed according to the [conventional commit](https://www.conventionalcommits.org/en/v1.0.0-beta.4/) specification.

## Pull Docker Image

```bash
$ docker pull psanetra/git-semver
```

## Commands

### latest

The `latest` command prints the latest semantic version in the current repository by comparing all git tags. Tag names may have a "v" prefix, but this commands prints the version always without that prefix. 

#### Examples

Print latest semantic version (ignoring pre-releases).
```bash
$ git-semver latest
1.2.3
```

Print latest semantic version including pre-releases.
```bash
$ git-semver latest --include-pre-releases
1.2.3-beta
```

### next

The `next` command can be used to calculate the next semantic version based on the history of the current branch. It fails if the git tag of the latest semantic version is not reachable on the current branch or if the tagged commit is not reachable because the repository is shallow.

#### Examples

Calculate next semantic version. (Will print the latest version if there were no relevant changes.)
```bash
$ git-semver next
1.2.3
```

Calculate next unstable semantic version. (Only if there is no stable version tag yet.)
```bash
$ git-semver next --stable=false
0.1.2
```

Calculate next alpha pre-release version with an appended counter.
```bash
$ git-semver next --pre-release-tag=alpha --pre-release-counter
1.2.3-alpha.1
```

### log

The `log` command prints the commit log of all commits, which were contained in a specified version or all commits since the latest version if no version is specified.

#### Examples

Print the commits, added in version 1.0.0.
```bash
$ git-semver log v1.0.0
commit 478bb9dfdca43216cda6cedcab27faf5c8fd68c0
Author: John Doe <john.doe@example.com>
Date:   Wed Jun 03 20:17:23 2020 +0000

    fix(some_component): Add fix

    Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nunc bibendum vulputate sapien vel mattis.

    Vivamus faucibus leo id libero suscipit, varius tincidunt neque interdum. Mauris rutrum at velit vitae semper.

    Fixes: http://issues.example.com/123
    BREAKING CHANGE: This commit is breaking some API.

commit f716712a4a26491533ba3b6d95e29f9beed85f47
Author: John Doe <john.doe@example.com>
Date:   Wed Jun 03 20:17:23 2020 +0000

    Some non-conventional-commit

commit d44f505f677d52ca23fb9a69de1f5bb6e6085a74
Author: John Doe <john.doe@example.com>
Date:   Wed Jun 03 20:17:22 2020 +0000

    feat: Add feature
```

Print only conventional commits, formatted as JSON.
```bash
$ git-semver log --conventional-commits v1.0.0
[
  {
    "type": "fix",
    "scope": "some_component",
    "breaking_change": true,
    "description": "Add fix",
    "body": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nunc bibendum vulputate sapien vel mattis.\n\nVivamus faucibus leo id libero suscipit, varius tincidunt neque interdum. Mauris rutrum at velit vitae semper.",
    "footers": {
      "BREAKING CHANGE": [
        "This commit is breaking some API."
      ],
      "Fixes": [
        "http://issues.example.com/123"
      ]
    }
  },
  {
    "type": "feat",
    "description": "Add feature"
  }
]
```

Print changelog formatted as markdown.
```bash
$ git-semver log --markdown v1.0.0
### BREAKING CHANGES

* **some_component** This commit is breaking some API.

### Features

* Add feature

### Bug Fixes

* **some_component** Add fix
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nunc bibendum vulputate sapien vel mattis.

Vivamus faucibus leo id libero suscipit, varius tincidunt neque interdum. Mauris rutrum at velit vitae semper.
```

### compare

The `compare` command is an utility command to compare two semantic versions.

- Prints `=` if both provided versions are equals.
- Prints `<` if the first provided version is lower than the second version.
- Prints `>` if the first provided version is greater than the second version.

#### Examples

Compare the versions `1.2.3` and `1.2.3-beta`
```bash
$ git-semver compare 1.2.3 1.2.3-beta
>
```

Compare the versions `1.2.3-alpha` and `1.2.3-beta`
```bash
$ git-semver compare 1.2.3-alpha 1.2.3-beta
<
```

Compare the versions `1.2.3` and `1.2.3+build-2018-12-31`
```bash
$ git-semver compare 1.2.3 1.2.3+build-2018-12-31
=
```

## Example GitLab Job Template

```yaml
stages:
  - tag

tag:
  image:
    name: psanetra/git-semver:latest
    entrypoint:
      - "/usr/bin/env"
      - "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
  stage: tag
  variables:
    GIT_DEPTH: 0
    GIT_FETCH_EXTRA_FLAGS: "--prune --prune-tags --tags"
  before_script:
    - apk add --upgrade --no-cache curl
  script: |
    set -ex
    LATEST_VERSION="$(git semver latest)"
    NEXT_VERSION="$(git semver next)"
    if [ "${LATEST_VERSION}" != "${NEXT_VERSION}" ]; then
      NEXT_TAG="v${NEXT_VERSION}"
      git tag "${NEXT_TAG}"
      CHANGELOG="$(git semver log --markdown ${NEXT_TAG})"
      curl -X POST \
        --header "JOB-TOKEN: ${CI_JOB_TOKEN}" \
        --form "tag_name=v${NEXT_VERSION}" \
        --form "ref=${CI_COMMIT_SHA}" \
        --form "description=${CHANGELOG}" \
        "${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/releases"
    fi
  only:
    - main
    - master
  except:
    - tags
    - schedules
```

## License

MIT
