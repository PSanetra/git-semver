package git_utils

import (
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func HashListContains(hashList []plumbing.Hash, hash plumbing.Hash) bool {

	for _, h := range hashList {
		if h == hash {
			return true
		}
	}

	return false
}
