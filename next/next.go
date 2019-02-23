package next

import (
	"github.com/pkg/errors"
	"github.com/psanetra/git-semver/conventional_commits"
	"github.com/psanetra/git-semver/git_utils"
	"github.com/psanetra/git-semver/logger"
	"github.com/psanetra/git-semver/semver"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/revlist"
	"io"
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

	latestReleaseVersion, latestReleaseVersionTag, err := findLatestVersion(repo, false)

	if err != nil {
		return nil, errors.WithMessage(err, "Error while trying to find latest release version tag")
	}

	var latestPreReleaseVersion *semver.Version

	if options.PreReleaseOptions.ShouldBePreRelease() {
		latestPreReleaseVersion, _, err = findLatestVersion(repo, true)
	}

	if err != nil {
		return nil, errors.WithMessage(err, "Error while trying to find latest pre-release version tag")
	}

	excludedCommits := make([]plumbing.Hash, 0, 1)

	if latestReleaseVersionTag != nil {
		excludedCommits = append(excludedCommits, latestReleaseVersionTag.Target)
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
		return nil, errors.WithMessage(err, "Could not find commits since latest version")
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


func findLatestVersion(repo *git.Repository, preRelease bool) (*semver.Version, *object.Tag, error) {
	latestVersionTag, err := findLatestVersionTag(repo, preRelease)

	if err != nil {
		return nil, nil, err
	}

	var latestVersion *semver.Version

	if preRelease && latestVersionTag == nil {
		return nil, nil, nil
	}

	if latestVersionTag == nil {
		latestVersion = &semver.Version{
			Major: 0,
			Minor: 0,
			Patch: 0,
		}
	} else {
		latestVersion = tagNameToVersion(latestVersionTag.Name)
	}

	return latestVersion, latestVersionTag, nil
}

func findLatestVersionTag(repo *git.Repository, includePreReleases bool) (*object.Tag, error) {

	tagIter, err := repo.TagObjects()

	if err != nil {
		return nil, err
	}

	defer tagIter.Close()

	var maxVersionTag *object.Tag
	var maxVersion = &semver.Version{
		Major: 0,
		Minor: 0,
		Patch: 0,
	}

	for tag, err := tagIter.Next(); err != io.EOF; tag, err = tagIter.Next() {
		if err != nil {
			return nil, err
		}

		version := tagNameToVersion(tag.Name)

		if version == nil || !includePreReleases && len(version.PreReleaseTag) > 0 {
			continue
		}

		if semver.CompareVersions(version, maxVersion) > 0 {
			maxVersion = version
			maxVersionTag = tag
		}
	}

	if maxVersionTag == nil {
		return nil, nil
	}

	headRef, err := repo.Head()

	if err != nil {
		return nil, err
	}

	headRefList, err := revlist.Objects(
		repo.Storer,
		[]plumbing.Hash{
			headRef.Hash(),
		},
		[]plumbing.Hash{},
	)

	if err != nil {
		return nil, err
	}

	if !git_utils.HashListContains(headRefList, maxVersionTag) {
		return nil, errors.New("latest version tag is not on current branch")
	}

	return maxVersionTag, nil
}

func tagNameToVersion(tagName string) *semver.Version {

	version, err := semver.ParseVersion(tagName)

	if err != nil {
		logger.Logger.Debug(err, ": Tag: ", tagName)
		return nil
	}

	return version
}
