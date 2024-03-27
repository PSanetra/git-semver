package git_utils

import "github.com/go-git/go-git/v5/plumbing/object"

type ByHistoryDesc []*object.Commit

func (a ByHistoryDesc) Len() int { return len(a) }
func (a ByHistoryDesc) Less(i, j int) bool {
	// Swap i and j as we want to sort descanding
	isAncestor, err := a[j].IsAncestor(a[i])

	if err != nil {
		return a[j].Committer.When.Before(a[i].Committer.When)
	}

	return isAncestor
}
func (a ByHistoryDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
