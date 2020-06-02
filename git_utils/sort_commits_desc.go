package git_utils

import "gopkg.in/src-d/go-git.v4/plumbing/object"

type ByCommitTimeDesc []*object.Commit

func (a ByCommitTimeDesc) Len() int           { return len(a) }
func (a ByCommitTimeDesc) Less(i, j int) bool { return a[i].Committer.When.After(a[j].Committer.When) }
func (a ByCommitTimeDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
