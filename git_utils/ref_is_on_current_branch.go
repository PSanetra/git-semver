package git_utils

import (
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/revlist"
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
		return errors.Errorf(message + " (tag: %s; commit: %s)", precedingRef.Name().String(), refCommitHash.String())
	}

	return nil
}
