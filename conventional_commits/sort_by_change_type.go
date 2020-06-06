package conventional_commits

type ByChangeTypeDesc []*ConventionalCommitMessage

func (a ByChangeTypeDesc) Len() int { return len(a) }
func (a ByChangeTypeDesc) Less(i, j int) bool {
	iPriority := ChangeTypePriorities[a[i].ChangeType]
	jPriority := ChangeTypePriorities[a[j].ChangeType]

	return (jPriority - iPriority) < 0
}
func (a ByChangeTypeDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
