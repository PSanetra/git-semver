package git_utils

import (
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func HashListContains(hashList []plumbing.Hash, tag *object.Tag) bool {

	for _, hash := range hashList {
		if hash == tag.Target {
			return true
		}
	}

	return false
}