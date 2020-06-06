package conventional_commits

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseCommitMessage_ParsesSimpleCommitMessage(t *testing.T) {

	commitMessage, err := ParseCommitMessage(`feat: my description`)

	assert.Nil(t, err)

	assert.Equal(
		t,
		&ConventionalCommitMessage{
			ChangeType:             "feat",
			Description:            "my description",
			ContainsBreakingChange: false,
			Footers:                map[string][]string{},
		},
		commitMessage,
	)

}

func TestParseCommitMessage_ParsesSimpleCommitMessageWithCaseInsensitiveType(t *testing.T) {

	commitMessage, err := ParseCommitMessage(`fEaT: my description`)

	assert.Nil(t, err)

	assert.Equal(
		t,
		&ConventionalCommitMessage{
			ChangeType:             "feat",
			Description:            "my description",
			ContainsBreakingChange: false,
			Footers:                map[string][]string{},
		},
		commitMessage,
	)

}

func TestParseCommitMessage_ParsesCommitMessageWithBreakingChangeIndicator(t *testing.T) {

	commitMessage, err := ParseCommitMessage(`feat!: my description`)

	assert.Nil(t, err)

	assert.Equal(
		t,
		&ConventionalCommitMessage{
			ChangeType:             "feat",
			Description:            "my description",
			ContainsBreakingChange: true,
			Footers:                map[string][]string{},
		},
		commitMessage,
	)

}

func TestParseCommitMessage_ParsesCommitMessageWithBreakingChangeIndicatorAfterScope(t *testing.T) {

	commitMessage, err := ParseCommitMessage(`feat(scope)!: my description`)

	assert.Nil(t, err)

	assert.Equal(
		t,
		&ConventionalCommitMessage{
			ChangeType:             "feat",
			Scope:                  "scope",
			Description:            "my description",
			ContainsBreakingChange: true,
			Footers:                map[string][]string{},
		},
		commitMessage,
	)

}

func TestParseCommitMessage_ParsesSimpleCommitMessageWithLineBreak(t *testing.T) {

	commitMessage, err := ParseCommitMessage("feat: my description\nwith line break")

	assert.Nil(t, err)

	assert.Equal(
		t,
		&ConventionalCommitMessage{
			ChangeType:  "feat",
			Description: "my description\nwith line break",
			Footers:     map[string][]string{},
		},
		commitMessage,
	)

}

func TestParseCommitMessage_ParsesCommitMessageWithBody(t *testing.T) {

	commitMessage, err := ParseCommitMessage(
		`
feat: my description
with line break

and this is a body

This is still the body
`,
	)

	assert.Nil(t, err)

	assert.Equal(
		t,
		&ConventionalCommitMessage{
			ChangeType:  "feat",
			Description: "my description\nwith line break",
			Body:        "and this is a body\n\nThis is still the body",
			Footers:     map[string][]string{},
		},
		commitMessage,
	)

}

func TestParseCommitMessage_ParsesCommitMessageWithFooter(t *testing.T) {

	commitMessage, err := ParseCommitMessage(
		`
feat: my description
with line break

and this is a body
with line break

this is still the body

Fix #123
Fix: http://example.com/123
Custom-Token: Custom-Token-Value
`,
	)

	assert.Nil(t, err)

	assert.Equal(
		t,
		&ConventionalCommitMessage{
			ChangeType:  "feat",
			Description: "my description\nwith line break",
			Body:        "and this is a body\nwith line break\n\nthis is still the body",
			Footers:      map[string][]string {
				"Fix": {"123", "http://example.com/123"},
				"Custom-Token": {"Custom-Token-Value"},
			},
		},
		commitMessage,
	)

}

func TestParseCommitMessage_SetsContainsBreakingChangeToFalseIfBodyContainsBreakingChangeInline(t *testing.T) {

	result, err := ParseCommitMessage(`feat: Some description

Body without BREAKING CHANGE description // BREAKING CHANGE:`)

	assert.Nil(t, err)
	assert.False( t, result.ContainsBreakingChange)

}

func TestParseCommitMessage_SetsContainsBreakingChangeToTrueIfBreakingChangeIndicatorExists(t *testing.T) {

	result, err := ParseCommitMessage(`feat!: Some description`)

	assert.Nil(t, err)
	assert.True( t, result.ContainsBreakingChange)
}

func TestParseCommitMessage_SetsContainsBreakingChangeToFalseIfBreakingChangeDescriptionInBody(t *testing.T) {

	testBodies := []string {
		"BREAKING CHANGE: commit breaks stuff\nsome-text-in-body-to-avoid-parsing-as-footer",
		"BREAKING CHANGES: commit breaks stuff\nsome-text-in-body-to-avoid-parsing-as-footer",
		"BREAKING_CHANGE: commit breaks stuff\nsome-text-in-body-to-avoid-parsing-as-footer",
		"BREAKING_CHANGES: commit breaks stuff\nsome-text-in-body-to-avoid-parsing-as-footer",
		"BREAKING-CHANGE: commit breaks stuff\nsome-text-in-body-to-avoid-parsing-as-footer",
		"BREAKING-CHANGES: commit breaks stuff\nsome-text-in-body-to-avoid-parsing-as-footer",
		`Body with breaking change description in second line:
BREAKING CHANGE: commit breaks stuff`,
	}

	for _, body := range testBodies {
		result, err := ParseCommitMessage("feat: Some description\n\n" + body)

		assert.Nil(t, err)
		assert.NotEmpty(t, result.Body)
		assert.False(t, result.ContainsBreakingChange, "Should not indicate a breaking change: " + body)
	}

}

func TestParseCommitMessage_SetsContainsBreakingChangeToTrueIfBreakingChangeTokenExistsInFooter(t *testing.T) {

	testBodiesAndFooters := []string {
		"BREAKING CHANGE: commit breaks stuff",
		"BREAKING CHANGES: commit breaks stuff",
		"BREAKING_CHANGE: commit breaks stuff",
		"BREAKING_CHANGES: commit breaks stuff",
		"BREAKING-CHANGE: commit breaks stuff",
		"BREAKING-CHANGES: commit breaks stuff",
		`This is the body:

BREAKING CHANGE: commit breaks stuff`,
	}

	for _, bodyAndFooter := range testBodiesAndFooters {
		result, err := ParseCommitMessage("feat: Some description\n\n" + bodyAndFooter)

		assert.Nil(t, err)
		assert.True(t, result.ContainsBreakingChange, "Could not parse breaking change indication from: " + bodyAndFooter)
	}

}

func TestCommitMessage_Compare_should_return_0_if_left_is_breaking_change_and_right_too(t *testing.T) {

	assert.Equal(
		t,
		0,
		(&ConventionalCommitMessage{
			ContainsBreakingChange: true,
		}).Compare(
			&ConventionalCommitMessage{
				ContainsBreakingChange: true,
			},
		),
	)

}

func TestCommitMessage_Compare_should_return_1_if_left_is_breaking_change_and_right_is_not(t *testing.T) {

	assert.Equal(
		t,
		1,
		(&ConventionalCommitMessage{
			ContainsBreakingChange: true,
		}).Compare(
			&ConventionalCommitMessage{
				ChangeType: FEATURE,
			},
		),
	)

	assert.Equal(
		t,
		1,
		(&ConventionalCommitMessage{
			ContainsBreakingChange: true,
		}).Compare(
			&ConventionalCommitMessage{
				ChangeType: FIX,
			},
		),
	)

	assert.Equal(
		t,
		1,
		(&ConventionalCommitMessage{
			ContainsBreakingChange: true,
		}).Compare(
			&ConventionalCommitMessage{
				ChangeType: CHORE,
			},
		),
	)

}

func TestCommitMessage_Compare_should_return_1_if_left_is_not_breaking_change_but_right_is(t *testing.T) {

	assert.Equal(
		t,
		-1,
		(&ConventionalCommitMessage{
			ChangeType: FEATURE,
		}).Compare(
			&ConventionalCommitMessage{
				ContainsBreakingChange: true,
			},
		),
	)

	assert.Equal(
		t,
		-1,
		(&ConventionalCommitMessage{
			ChangeType: FIX,
		}).Compare(
			&ConventionalCommitMessage{
				ContainsBreakingChange: true,
			},
		),
	)

	assert.Equal(
		t,
		-1,
		(&ConventionalCommitMessage{
			ChangeType: CHORE,
		}).Compare(
			&ConventionalCommitMessage{
				ContainsBreakingChange: true,
			},
		),
	)

}

func TestCommitMessage_Compare_should_return_0_if_left_is_feature_and_right_too(t *testing.T) {

	assert.Equal(
		t,
		0,
		(&ConventionalCommitMessage{
			ChangeType: FEATURE,
		}).Compare(
			&ConventionalCommitMessage{
				ChangeType: FEATURE,
			},
		),
	)

}

func TestCommitMessage_Compare_should_return_1_if_left_is_feature_and_right_is_not(t *testing.T) {

	assert.Equal(
		t,
		1,
		(&ConventionalCommitMessage{
			ChangeType: FEATURE,
		}).Compare(
			&ConventionalCommitMessage{
				ChangeType: FIX,
			},
		),
	)

	assert.Equal(
		t,
		1,
		(&ConventionalCommitMessage{
			ChangeType: FEATURE,
		}).Compare(
			&ConventionalCommitMessage{
				ChangeType: CHORE,
			},
		),
	)

}

func TestCommitMessage_Compare_should_return_1_if_left_is_not_feature_but_right_is(t *testing.T) {

	assert.Equal(
		t,
		-1,
		(&ConventionalCommitMessage{
			ChangeType: FIX,
		}).Compare(
			&ConventionalCommitMessage{
				ChangeType: FEATURE,
			},
		),
	)

	assert.Equal(
		t,
		-1,
		(&ConventionalCommitMessage{
			ChangeType: CHORE,
		}).Compare(
			&ConventionalCommitMessage{
				ChangeType: FEATURE,
			},
		),
	)

}

func TestCommitMessage_Compare_should_return_0_if_left_is_fix_and_right_too(t *testing.T) {

	assert.Equal(
		t,
		0,
		(&ConventionalCommitMessage{
			ChangeType: FIX,
		}).Compare(
			&ConventionalCommitMessage{
				ChangeType: FIX,
			},
		),
	)

}

func TestCommitMessage_Compare_should_return_1_if_left_is_fix_and_right_is_chore(t *testing.T) {

	assert.Equal(
		t,
		1,
		(&ConventionalCommitMessage{
			ChangeType: FIX,
		}).Compare(
			&ConventionalCommitMessage{
				ChangeType: CHORE,
			},
		),
	)

}

func TestCommitMessage_Compare_should_return_1_if_left_is_chore_and_right_is_fix(t *testing.T) {

	assert.Equal(
		t,
		-1,
		(&ConventionalCommitMessage{
			ChangeType: CHORE,
		}).Compare(
			&ConventionalCommitMessage{
				ChangeType: FIX,
			},
		),
	)

}

func TestCommitMessage_Compare_should_return_0_if_left_is_chore_and_right_is_doc(t *testing.T) {

	assert.Equal(
		t,
		0,
		(&ConventionalCommitMessage{
			ChangeType: CHORE,
		}).Compare(
			&ConventionalCommitMessage{
				ChangeType: DOCS,
			},
		),
	)

}
