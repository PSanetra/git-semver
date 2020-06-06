package conventional_commits

type ChangeType string

const (
	FEATURE  ChangeType = "feat"
	FIX      ChangeType = "fix"
	CHORE    ChangeType = "chore"
	PERF     ChangeType = "perf"
	STYLE    ChangeType = "style"
	DOCS     ChangeType = "docs"
	REFACTOR ChangeType = "refactor"
	CI       ChangeType = "ci"
)

var ChangeTypePriorities = map[ChangeType]int{
	FEATURE:  10,
	FIX:      9,
	PERF:     8,
	DOCS:     7,
	CHORE:    6,
	CI:       5,
	STYLE:    4,
	REFACTOR: 3,
}
