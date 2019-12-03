package next

import (
	"github.com/pkg/errors"
	"github.com/psanetra/git-semver/conventional_commits"
	"github.com/psanetra/git-semver/latest"
	"github.com/psanetra/git-semver/logger"
	"github.com/psanetra/git-semver/semver"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/revlist"
)

type NextOptions struct {
	Workdir           string
	Stable            bool
	PreReleaseOptions semver.PreReleaseOptions
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

	latestReleaseVersion, latestReleaseVersionTag, err := latest.FindLatestVersion(repo, false)

	if err != nil {
		return nil, errors.WithMessage(err, "Error while trying to find latest release version tag")
	}

	var latestPreReleaseVersion *semver.Version

	if options.PreReleaseOptions.ShouldBePreRelease() {
		latestPreReleaseVersion, _, err = latest.FindLatestVersion(repo, true)
	}

	if err != nil {
		return nil, errors.WithMessage(err, "Error while trying to find latest pre-release version tag")
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

	maxPrioCommitMessage := &conventional_commits.CommitMessage{}

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

		if message.IsBreakingChange() {
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

func commitMessageToSemverChange(msg *conventional_commits.CommitMessage) semver.Change {

	var semverChange semver.Change

	if msg == nil {
		return semverChange
	} else if msg.IsBreakingChange() {
		semverChange = semver.BREAKING
	} else if msg.ChangeType == conventional_commits.FEATURE {
		semverChange = semver.NEW_FEATURE
	} else if msg.ChangeType == conventional_commits.FIX {
		semverChange = semver.FIX
	}

	return semverChange
}
