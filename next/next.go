package next

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/revlist"
	"github.com/pkg/errors"
	"github.com/psanetra/git-semver/conventional_commits"
	"github.com/psanetra/git-semver/git_utils"
	"github.com/psanetra/git-semver/latest"
	"github.com/psanetra/git-semver/logger"
	"github.com/psanetra/git-semver/semver"
)

type NextOptions struct {
	Workdir            string
	Stable             bool
	MajorVersionFilter int
	PreReleaseOptions  semver.PreReleaseOptions
}

func Next(options NextOptions) (*semver.Version, error) {

	repo, err := git.PlainOpenWithOptions(options.Workdir, &git.PlainOpenOptions{
		DetectDotGit: true,
	})

	if err != nil {
		return nil, errors.WithMessage(err, "Could not open git repository")
	}

	headRef, err := repo.Head()

	if err != nil {
		return nil, errors.WithMessage(err, "Could not find HEAD")
	}

  defaultBranchRef, err := repo.Reference(plumbing.ReferenceName("refs/remotes/origin/HEAD"), true)
  if err != nil {
      return nil, errors.WithMessage(err, "Could not determine default branch")
  }

  latestReleaseVersion, latestReleaseVersionTag, err := latest.FindLatestVersionOnBranch(repo, options.MajorVersionFilter, defaultBranchRef.Name().Short(), false)

	if err != nil {
		return nil, errors.WithMessage(err, "Error while trying to find latest release version tag")
	}

	if latestReleaseVersionTag != nil {
		if err = git_utils.AssertRefIsReachable(repo, latestReleaseVersionTag, headRef, "Latest tag is not on HEAD. This is necessary as the next version is calculated based on the commits since the latest version tag."); err != nil {
			return nil, err
		}
	}

	if latestReleaseVersion == nil {
		latestReleaseVersion = &semver.EmptyVersion
	}

	var latestPreReleaseVersion *semver.Version
	var latestPreReleaseVersionTag *plumbing.Reference

	if options.PreReleaseOptions.ShouldBePreRelease() {
		defaultBranchRef, err := repo.Reference(plumbing.ReferenceName("refs/remotes/origin/HEAD"), true)
    if err != nil {
        return nil, errors.WithMessage(err, "Could not determine default branch")
    }

    latestPreReleaseVersion, latestPreReleaseVersionTag, err = latest.FindLatestVersionOnBranch(repo, options.MajorVersionFilter, defaultBranchRef.Name().Short(), true)
  }

	if err != nil {
		return nil, errors.WithMessage(err, "Error while trying to find latest pre-release version tag")
	}

	if latestPreReleaseVersionTag != nil {
		if err = git_utils.AssertRefIsReachable(repo, latestPreReleaseVersionTag, headRef, "Latest tag is not on HEAD. This is necessary as the next version is calculated based on the commits since the latest version tag."); err != nil {
			return nil, err
		}
	}

	excludedCommits := make([]plumbing.Hash, 0, 1)

	if latestReleaseVersionTag != nil {
		excludedCommits = append(excludedCommits, latestReleaseVersionTag.Hash())
	}

	// historyDiff also contains other hashes than commit hashes (e.g. blob or tree hashes)
	historyDiff, err := revlist.Objects(
		repo.Storer,
		[]plumbing.Hash{
			headRef.Hash(),
		},
		excludedCommits,
	)

	if err != nil {
		objectInfo := " (HEAD: " + headRef.Hash().String()

		if latestReleaseVersionTag != nil {
			objectInfo = ", Latest Version: {" + latestReleaseVersionTag.Name().Short() + ", " + latestReleaseVersionTag.Hash().String() + "}"
		}

		objectInfo += ")"

		return nil, errors.WithMessage(err, "Could not find commits since latest version"+objectInfo)
	}

	var nextVersion semver.Version

	maxPrioCommitMessage := &conventional_commits.ConventionalCommitMessage{}

	for _, hash := range historyDiff {
		commit, err := repo.CommitObject(hash)

		if err == plumbing.ErrObjectNotFound {
			// hash is not a commit object
			continue
		}

		if err != nil {
			return nil, errors.WithMessage(err, "Could not read commit "+hash.String())
		}

		message, err := conventional_commits.ParseCommitMessage(commit.Message)

		if err != nil {
			logger.Logger.Debug(err)
			continue
		}

		if message.Compare(maxPrioCommitMessage) <= 0 {
			continue
		}

		maxPrioCommitMessage = message

		if message.ContainsBreakingChange {
			break
		}
	}

	nextVersion, err = semver.Increment(
		*latestReleaseVersion,
		latestPreReleaseVersion,
		options.Stable,
		commitMessageToSemverChange(maxPrioCommitMessage),
		&options.PreReleaseOptions,
	)

	if err != nil {
		return nil, errors.WithMessage(err, "Could not increment version")
	}

	return &nextVersion, nil

}

func commitMessageToSemverChange(msg *conventional_commits.ConventionalCommitMessage) semver.Change {

	var semverChange semver.Change

	if msg == nil {
		return semverChange
	} else if msg.ContainsBreakingChange {
		semverChange = semver.BREAKING
	} else if msg.ChangeType == conventional_commits.FEATURE {
		semverChange = semver.NEW_FEATURE
	} else if msg.ChangeType == conventional_commits.FIX {
		semverChange = semver.FIX
	}

	return semverChange
}
