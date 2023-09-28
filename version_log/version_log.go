package version_log

import (
	"github.com/pkg/errors"
	"github.com/psanetra/git-semver/git_utils"
	"github.com/psanetra/git-semver/latest"
	"github.com/psanetra/git-semver/logger"
	"github.com/psanetra/git-semver/semver"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/revlist"
	"sort"
)

type VersionLogOptions struct {
	Workdir                  string
	Version                  *semver.Version
	ExcludePreReleaseCommits bool
}

// Returns all commits since the preceding version to options.Version
func VersionLog(options VersionLogOptions) ([]*object.Commit, error) {

	repo, err := git.PlainOpenWithOptions(options.Workdir, &git.PlainOpenOptions{
		DetectDotGit: true,
	})

	if err != nil {
		return nil, errors.WithMessage(err, "Could not open git repository")
	}

	versions, err := git_utils.GetVersions(repo)

	if err != nil {
		return nil, errors.WithMessage(err, "Could not find Tags")
	}

	var targetVersionRef *plumbing.Reference

	if options.Version == nil {
		headRef, err := repo.Head()

		if err != nil {
			return nil, errors.WithMessage(err, "Could not find HEAD")
		}

		targetVersionRef = headRef
	} else {
		targetVersionRef, err = findTagForVersion(repo, options.Version.ToString())

		if err != nil {
			return nil, err
		}
	}

	greatestPreceding := semver.FindGreatestPreceding(options.Version, versions, !options.ExcludePreReleaseCommits)

	var fromVersionTag *plumbing.Reference

	if greatestPreceding == nil {
		if targetVersionRef == nil {
			_, fromVersionTag, err = latest.FindLatestVersion(repo, -1, options.ExcludePreReleaseCommits)

			if err != nil {
				return nil, errors.WithMessage(err, "Could not find latest version")
			}
		}
	} else {
		fromVersionTag, err = findTagForVersion(repo, greatestPreceding.ToString())

		if err != nil {
			return nil, err
		}
	}

	excludedCommits := make([]plumbing.Hash, 0, 1)

	// fromVersionTag may be null if there is no latest version
	if fromVersionTag != nil {
		excludedCommits = append(excludedCommits, fromVersionTag.Hash())
	}

	// historyRange also contains other hashes than commit hashes (e.g. blob or tree hashes)
	historyRange, err := revlist.Objects(
		repo.Storer,
		[]plumbing.Hash{
			targetVersionRef.Hash(),
		},
		excludedCommits,
	)

	if err != nil {
		return nil, errors.WithMessage(err, "Could not find commits.")
	}

	commits := make([]*object.Commit, 0, len(historyRange))

	for _, hash := range historyRange {
		commit, err := repo.CommitObject(hash)

		if err == plumbing.ErrObjectNotFound {
			// hash is not a commit object
			continue
		}

		if err != nil {
			return nil, errors.WithMessage(err, "Could not read commit "+hash.String())
		}

		commits = append(commits, commit)
	}

	sort.Sort(git_utils.ByHistoryDesc(commits))

	// Return most recent commits first
	return commits, nil
}

func findTagForVersion(repo *git.Repository, version string) (*plumbing.Reference, error) {
	tag, err := repo.Tag("v" + version)

	if err != nil {
		logger.Logger.Debugln("Could not find tag v"+version+":", err)

		tag, err = repo.Tag(version)

		if err != nil {
			return nil, errors.WithMessage(err, "Could not find tag "+version+" or v"+version)
		}
	}

	return tag, err
}
