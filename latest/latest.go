package latest

import (
	"github.com/pkg/errors"
	"github.com/psanetra/git-semver/git_utils"
	"github.com/psanetra/git-semver/logger"
	"github.com/psanetra/git-semver/semver"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/revlist"
	"io"
)

type LatestOptions struct {
	Workdir            string
	IncludePreReleases bool
}

func Latest(options LatestOptions) (*semver.Version, error) {

	repo, err := git.PlainOpenWithOptions(options.Workdir, &git.PlainOpenOptions{
		DetectDotGit: true,
	})

	if err != nil {
		return nil, errors.WithMessage(err, "Could not open git repository")
	}

	latestReleaseVersion, _, err := FindLatestVersion(repo, options.IncludePreReleases)

	if latestReleaseVersion == nil {
		latestReleaseVersion = &semver.EmptyVersion
	}

	return latestReleaseVersion, err

}

func FindLatestVersion(repo *git.Repository, preRelease bool) (*semver.Version, *plumbing.Reference, error) {
	latestVersionTag, err := findLatestVersionTag(repo, preRelease)

	if err != nil {
		return nil, nil, err
	}

	if latestVersionTag == nil {
		return nil, nil, nil
	}

	return tagNameToVersion(latestVersionTag.Name().Short()), latestVersionTag, nil
}

func findLatestVersionTag(repo *git.Repository, includePreReleases bool) (*plumbing.Reference, error) {

	tagIter, err := repo.Tags()

	if err != nil {
		return nil, err
	}

	defer tagIter.Close()

	var maxVersionTag *plumbing.Reference
	var maxVersion = &semver.EmptyVersion

	for tag, err := tagIter.Next(); err != io.EOF; tag, err = tagIter.Next() {
		if err != nil {
			return nil, err
		}

		version := tagNameToVersion(tag.Name().Short())

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

	maxVersionCommitHash := git_utils.RefToCommitHash(repo.Storer, maxVersionTag)

	if !git_utils.HashListContains(headRefList, maxVersionCommitHash) {
		return nil, errors.Errorf("latest version tag (%s on %s) is not on current branch", maxVersionTag.Name().String(), maxVersionCommitHash.String())
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
