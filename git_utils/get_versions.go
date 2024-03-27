package git_utils

import (
	"github.com/go-git/go-git/v5"
	"github.com/psanetra/git-semver/semver"
	"io"
)

func GetVersions(repo *git.Repository) ([]*semver.Version, error) {

	tagIter, err := repo.Tags()

	if err != nil {
		return nil, err
	}

	defer tagIter.Close()

	var ret []*semver.Version

	for tag, err := tagIter.Next(); err != io.EOF; tag, err = tagIter.Next() {
		if err != nil {
			return nil, err
		}

		tagName := tag.Name().Short()

		version, err := semver.ParseVersion(tagName)

		if err != nil {
			continue
		}

		ret = append(ret, version)
	}

	return ret, nil
}
