package git_utils

import (
	"github.com/go-git/go-git/v5/plumbing"
)

func HashListContains(hashList []plumbing.Hash, hash plumbing.Hash) bool {

	for _, h := range hashList {
		if h == hash {
			return true
		}
	}

	return false
}
