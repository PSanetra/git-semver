package git_utils

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/revlist"
	"github.com/pkg/errors"
)

func AssertRefIsReachable(repo *git.Repository, precedingRef *plumbing.Reference, headRef *plumbing.Reference, message string) error {
	toRefList, err := revlist.Objects(
		repo.Storer,
		[]plumbing.Hash{
			headRef.Hash(),
		},
		[]plumbing.Hash{},
	)

	if err != nil {
		return err
	}

	refCommitHash := RefToCommitHash(repo.Storer, precedingRef)

	if !HashListContains(toRefList, refCommitHash) {
		return errors.Errorf(message+" (tag: %s; commit: %s)", precedingRef.Name().String(), refCommitHash.String())
	}

	return nil
}
