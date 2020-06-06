package conventional_commits

import (
	"github.com/pkg/errors"
	"github.com/psanetra/git-semver/regex_utils"
	"regexp"
	"strings"
)

var messageRegex = regexp.MustCompile(`(?m)^(?P<ChangeType>[a-zA-Z]+)(\((?P<Scope>[^)]*)\))?(?P<BCIndicator>!)?:\s*(?P<Description>[^\n]([^\n]|\n[^\n])*)(?P<BodyAndFooter>\n{2,}(.|\n)*)?$`)
var footerTokenRegex = regexp.MustCompile(`(?P<Token>[a-zA-Z_\-]+|BREAKING[ _\-]CHANGE(S)?)( \#|: )(?P<Value>[^\n]+)(\n|(\n)?$)`)
var footerRegex = regexp.MustCompile("\n\n(" + footerTokenRegex.String() + ")+$")

var breakingChangeRegex = regexp.MustCompile("^BREAKING[ _\\-]CHANGE(S)?$")

type ConventionalCommitMessage struct {
	ChangeType             ChangeType          `json:"type"`
	Scope                  string              `json:"scope,omitempty"`
	ContainsBreakingChange bool                `json:"breaking_change,omitempty"`
	Description            string              `json:"description"`
	Body                   string              `json:"body,omitempty"`
	Footers                map[string][]string `json:"footers,omitempty"`
}

// inspired by https://www.conventionalcommits.org
func ParseCommitMessage(message string) (*ConventionalCommitMessage, error) {

	match := regex_utils.SubmatchMap(messageRegex, message)

	if match == nil {
		return nil, errors.New("Could not parse commit message \"" + message + "\"")
	}

	breakingChangeIndicator := match["BCIndicator"]

	bodyAndFooter := match["BodyAndFooter"]

	body := bodyAndFooter
	footer := make(map[string][]string)

	var footerIndicies = footerRegex.FindStringIndex(bodyAndFooter)

	if len(footerIndicies) > 0 {
		body = bodyAndFooter[:footerIndicies[0]]

		for _, submatches := range footerTokenRegex.FindAllStringSubmatch(bodyAndFooter[footerIndicies[0]:], 100) {
			submatchMap := regex_utils.SubmatchMapFromSubmatches(footerTokenRegex, submatches)
			token := submatchMap["Token"]
			tokenValueList := footer[token]
			footer[token] = append(tokenValueList, submatchMap["Value"])
		}
	} else {
		body = bodyAndFooter
	}

	body = trimWhitespace(body)

	commitMessage := &ConventionalCommitMessage{
		ChangeType:             ChangeType(strings.ToLower(match["ChangeType"])),
		Scope:                  match["Scope"],
		ContainsBreakingChange: breakingChangeIndicator == "!",
		Description:            match["Description"],
		Body:                   body,
		Footers:                footer,
	}

	commitMessage.ContainsBreakingChange = commitMessage.ContainsBreakingChange || commitMessage.footerHasBreakingChange()

	return commitMessage, nil
}

func trimWhitespace(str string) string {
	return strings.Trim(str, " \t\r\n")
}

func (c *ConventionalCommitMessage) Compare(other *ConventionalCommitMessage) int {
	if c.ContainsBreakingChange {
		if other.ContainsBreakingChange {
			return 0
		} else {
			return 1
		}
	} else if other.ContainsBreakingChange {
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

func (c *ConventionalCommitMessage) footerHasBreakingChange() bool {
	for key, _ := range c.Footers {
		if breakingChangeRegex.MatchString(key) {
			return true
		}
	}

	return false
}

func (c *ConventionalCommitMessage) breakingChangeDescriptions() []string {
	var ret []string

	for key, value := range c.Footers {
		if breakingChangeRegex.MatchString(key) {
			ret = append(ret, value...)
		}
	}

	return ret
}
