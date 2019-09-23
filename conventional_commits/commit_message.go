package conventional_commits

import (
	"github.com/pkg/errors"
	"github.com/psanetra/git-semver/regex_utils"
	"regexp"
	"strings"
)

var messageRegex = regexp.MustCompile(`(?m)^(?P<ChangeType>[a-z]+)(\([^)]*\))?(?P<BCIndicator>!)?:\s*(?P<Description>[^\n]([^\n]|\n[^\n])*)(\n{2,}(?P<BodyAndFooter>(.|\n)*))?$`)

var breakingChangeRegex = regexp.MustCompile("(?m)^BREAKING[ _\\-]CHANGE(S)?:")

type CommitMessage struct {
	ChangeType                 ChangeType
	HasBreakingChangeIndicator bool
	Description                string
	Body                       string
	Footer                     string
}

// inspired by https://www.conventionalcommits.org
func ParseCommitMessage(message string) (*CommitMessage, error) {

	match := regex_utils.SubmatchMap(messageRegex, message)

	if match == nil {
		return nil, errors.New("Could not parse commit message \"" + message + "\"")
	}

	breakingChangeIndicator := match["BCIndicator"]

	bodyAndFooter := trimWhitespace(match["BodyAndFooter"])

	footerSeparatorIndex := strings.LastIndex(bodyAndFooter, "\n\n")

	body := bodyAndFooter
	footer := ""

	if footerSeparatorIndex >= 0 {
		body = bodyAndFooter[:footerSeparatorIndex]
		footer = trimWhitespace(bodyAndFooter[footerSeparatorIndex+2:])
	}

	return &CommitMessage{
		ChangeType:                 ChangeType(match["ChangeType"]),
		HasBreakingChangeIndicator: breakingChangeIndicator == "!",
		Description:                match["Description"],
		Body:                       body,
		Footer:                     footer,
	}, nil
}

func trimWhitespace(str string) string {
	return strings.Trim(str, " \t\r\n")
}

func (c *CommitMessage) Compare(other *CommitMessage) int {
	if c.IsBreakingChange() {
		if other.IsBreakingChange() {
			return 0
		} else {
			return 1
		}
	} else if other.IsBreakingChange() {
		return -1
	}

	if c.ChangeType == FEATURE {
		if other.ChangeType == FEATURE {
			return 0
		} else {
			return 1
		}
	} else if other.ChangeType == FEATURE {
		return -1
	}

	if c.ChangeType == FIX {
		if other.ChangeType == FIX {
			return 0
		} else {
			return 1
		}
	} else if other.ChangeType == FIX {
		return -1
	}

	return 0
}

func (c *CommitMessage) IsBreakingChange() bool {
	if c.HasBreakingChangeIndicator {
		return true
	}

	containsBreakingChangeDescription := breakingChangeRegex.MatchString(c.Footer)

	if containsBreakingChangeDescription {
		return true
	}

	return breakingChangeRegex.MatchString(c.Body)
}
