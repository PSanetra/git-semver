package conventional_commits

import (
	"sort"
	"strings"
)

func ToMarkdown(messages []*ConventionalCommitMessage) string {

	var markDownParts []string

	commitsContainingBreakingChanges := ByChangeTypeDesc(filterBreakingChanges(messages))
	sort.Stable(commitsContainingBreakingChanges)

	if len(commitsContainingBreakingChanges) > 0 {
		markDownParts = append(markDownParts, markdownBreakingChanges(commitsContainingBreakingChanges))
	}

	features := filterByNonBreakingChangeType(FEATURE, messages)

	if len(features) > 0 {
		featuresString := "### Features\n\n"
		featuresString += markdownSimpleChanges(features)
		markDownParts = append(markDownParts, featuresString)
	}

	fixes := filterByNonBreakingChangeType(FIX, messages)

	if len(fixes) > 0 {
		fixesString := "### Bug Fixes\n\n"
		fixesString += markdownSimpleChanges(fixes)
		markDownParts = append(markDownParts, fixesString)
	}

	return strings.Join(markDownParts, "\n")
}

func markdownBreakingChanges(commitsContainingBreakingChanges ByChangeTypeDesc) string {
	ret := "### BREAKING CHANGES\n\n"

	for _, change := range commitsContainingBreakingChanges {

		breakingChangeDescriptions := change.breakingChangeDescriptions()

		if len(breakingChangeDescriptions) == 0 {
			ret += "* "

			if change.Scope != "" {
				ret += "**" + change.Scope + "** "
			}

			ret += change.Description + "\n"

			if change.Body != "" {
				ret += "\n" + change.Body
			}
		} else {
			for _, description := range breakingChangeDescriptions {

				ret += "* "

				if change.Scope != "" {
					ret += "**" + change.Scope + "** "
				}

				ret += description + "\n"

			}
		}

	}

	return ret
}

func markdownSimpleChanges(changes []*ConventionalCommitMessage) string {
	ret := ""

	for _, change := range changes {
		// skip breaking changes without separate description, because they are listed in another section
		if change.ContainsBreakingChange && len(change.breakingChangeDescriptions()) == 0 {
			continue
		}

		ret += "* "

		if change.Scope != "" {
			ret += "**" + change.Scope + "** "
		}

		ret += change.Description + "\n"

		if change.Body != "" {
			ret += change.Body
			ret += "\n"
		}

	}

	return ret
}

func filterBreakingChanges(messages []*ConventionalCommitMessage) []*ConventionalCommitMessage {
	var ret []*ConventionalCommitMessage

	for _, c := range messages {
		if c.ContainsBreakingChange {
			ret = append(ret, c)
		}
	}

	return ret
}

func filterByNonBreakingChangeType(changeType ChangeType, messages []*ConventionalCommitMessage) []*ConventionalCommitMessage {
	var ret []*ConventionalCommitMessage

	for _, c := range messages {
		if c.ChangeType == changeType {
			ret = append(ret, c)
		}
	}

	return ret
}
