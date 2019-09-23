package git_utils

import (
	"github.com/psanetra/git-semver/logger"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

func RefToCommitHash(storer storer.EncodedObjectStorer, tagRef *plumbing.Reference) plumbing.Hash {
	o, err := object.GetObject(storer, tagRef.Hash())

	if err != nil {
		logger.Logger.Fatalln("Error on resolving tag hash (", tagRef.Hash().String(), "): ", err)
	}

	switch o := o.(type) {
	case *object.Commit:
		return o.Hash
	case *object.Tag:
		return o.Target
	default:
		logger.Logger.Fatalln("Error on resolving tag hash (", tagRef.Hash().String(), "): ", err)
		return plumbing.Hash{}
	}
}
