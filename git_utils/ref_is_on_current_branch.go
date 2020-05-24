package git_utils

import (
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/revlist"
)

func AssertRefIsOnHead(repo *git.Repository, ref *plumbing.Reference, message string) error {
	headRef, err := repo.Head()

	if err != nil {
		return err
	}

	headRefList, err := revlist.Objects(
		repo.Storer,
		[]plumbing.Hash{
			headRef.Hash(),
		},
		[]plumbing.Hash{},
	)

	if err != nil {
		return err
	}

	refCommitHash := RefToCommitHash(repo.Storer, ref)

	if !HashListContains(headRefList, refCommitHash) {
		return errors.Errorf(message + " (tag: %s; commit: %s)", ref.Name().String(), refCommitHash.String())
	}

	return nil
}
