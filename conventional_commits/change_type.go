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
)

func (c ChangeType) TriggersRelease() bool {
	switch c {
	case FEATURE:
		return true
	case FIX:
		return true
	}

	return false
}

func (c ChangeType) ToChangelogString() string {
	switch c {
	case FEATURE:
		return "Feature"
	case FIX:
		return "Fix"
	case CHORE:
		return "Maintenance"
	case PERF:
		return "Performance"
	case STYLE:
		return "Code style"
	case DOCS:
		return "Documentation"
	case REFACTOR:
		return "Refactoring"
	}

	return string(c)
}
