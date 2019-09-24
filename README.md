# git-semver
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FPSanetra%2Fgit-semver.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2FPSanetra%2Fgit-semver?ref=badge_shield)


git-semver is a command line tool to calculate semantic versions based on the git history and tags of a repository.

git-semver assumes that the commit messages in the git history are wellformed according to the [conventional commit](https://www.conventionalcommits.org/en/v1.0.0-beta.4/) specification.

## Howto

Calculate next semantic version
```bash
$ git-semver next
1.2.3
```

Calculate next unstable semantic version (only if there is no stable version tag yet)
```bash
$ git-semver next --stable=false
0.1.2
```

Calculate next alpha pre-release version with an appended counter
```bash
$ git-semver next --pre-release-tag=alpha --pre-release-counter
1.2.3-alpha.1
```

Print latest semantic version (excluding pre-releases)
```bash
$ git-semver latest
1.2.3
```

Print latest semantic version including pre-releases
```bash
$ git-semver latest --include-pre-releases
1.2.3-beta
```

## License

MIT


[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FPSanetra%2Fgit-semver.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2FPSanetra%2Fgit-semver?ref=badge_large)