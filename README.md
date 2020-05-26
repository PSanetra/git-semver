# git-semver
[![Docker Cloud Automated build](https://img.shields.io/docker/cloud/automated/psanetra/git-semver)](https://hub.docker.com/r/psanetra/git-semver) [![Docker Image Pulls](https://img.shields.io/docker/pulls/psanetra/git-semver)](https://hub.docker.com/r/psanetra/git-semver)

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
$ git-semver compare 1.2.3-alpha 1.2.3+build-2018-12-31
=
```

## License

MIT
