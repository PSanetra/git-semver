# git-semver
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org) [![Docker Cloud Automated build](https://img.shields.io/docker/cloud/automated/psanetra/git-semver)](https://hub.docker.com/r/psanetra/git-semver) [![Docker Image Pulls](https://img.shields.io/docker/pulls/psanetra/git-semver)](https://hub.docker.com/r/psanetra/git-semver)

git-semver is a command line tool to calculate [semantic versions](https://semver.org/spec/v2.0.0.html) based on the git history and tags of a repository.

git-semver assumes that the commit messages in the git history are well-formed according to the [conventional commit](https://www.conventionalcommits.org/en/v1.0.0-beta.4/) specification.

## Pull Docker Image

```bash
$ docker pull psanetra/git-semver
```

## Commands

### latest

The `latest` command prints the latest semantic version in the current repository by comparing all git tags. Tag names may have a "v" prefix, but this command always prints the version without that prefix.

#### Examples

Print the latest semantic version (ignoring pre-releases).
```bash
$ git-semver latest
1.2.3
```

Print the latest semantic version including pre-releases.
```bash
$ git-semver latest --include-pre-releases
1.2.3-beta
```

### next

The `next` command calculates the next semantic version based on the history of the **default branch**. It ensures that the latest version tag is reachable on the default branch before computing the next version.

#### **Changes in `next` Command**
- The `next` command now only considers tags that exist on the **default branch**.
- It dynamically determines the default branch (e.g., `main` or `master`).
- It ensures pre-release versions are also based only on the default branch.

#### Examples

Calculate the next semantic version. (Will print the latest version if there were no relevant changes.)
```bash
$ git-semver next
1.2.3
```

Calculate the next unstable semantic version. (Only if there is no stable version tag yet.)
```bash
$ git-semver next --stable=false
0.1.2
```

Calculate the next alpha pre-release version with an appended counter.
```bash
$ git-semver next --pre-release-tag=alpha --pre-release-counter
1.2.3-alpha.1
```

### log

The `log` command prints the commit log of all commits contained in a specified version or all commits since the latest version if no version is specified.

#### Examples

Print the commits added in version 1.0.0.
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
  }
]
```

### compare

The `compare` command is a utility command to compare two semantic versions.

- Prints `=` if both provided versions are equal.
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

## License

MIT
