package conventional_commits

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseCommitMessageParsesSimpleCommitMessage(t *testing.T) {

	commitMessage, err := ParseCommitMessage(`feat: my description`)

	assert.Nil(t, err)

	assert.Equal(
		t,
		&CommitMessage{
			ChangeType:                 "feat",
			Description:                "my description",
			HasBreakingChangeIndicator: false,
		},
		commitMessage,
	)

}

func TestParseCommitMessageParsesCommitMessageWithBreakingChangeIndicator(t *testing.T) {

	commitMessage, err := ParseCommitMessage(`feat!: my description`)

	assert.Nil(t, err)

	assert.Equal(
		t,
		&CommitMessage{
			ChangeType:                 "feat",
			Description:                "my description",
			HasBreakingChangeIndicator: true,
		},
		commitMessage,
	)

}

func TestParseCommitMessageParsesCommitMessageWithBreakingChangeIndicatorAfterScope(t *testing.T) {

	commitMessage, err := ParseCommitMessage(`feat(scope)!: my description`)

	assert.Nil(t, err)

	assert.Equal(
		t,
		&CommitMessage{
			ChangeType:                 "feat",
			Description:                "my description",
			HasBreakingChangeIndicator: true,
		},
		commitMessage,
	)

}

func TestParseCommitMessageParsesSimpleCommitMessageWithLineBreak(t *testing.T) {

	commitMessage, err := ParseCommitMessage("feat: my description\nwith line break")

	assert.Nil(t, err)

	assert.Equal(
		t,
		&CommitMessage{
			ChangeType:  "feat",
			Description: "my description\nwith line break",
		},
		commitMessage,
	)

}

func TestParseCommitMessageParsesCommitMessageWithBody(t *testing.T) {

	commitMessage, err := ParseCommitMessage(
		`
feat: my description
with line break

and this is a body
`,
	)

	assert.Nil(t, err)

	assert.Equal(
		t,
		&CommitMessage{
			ChangeType:  "feat",
			Description: "my description\nwith line break",
			Body:        "and this is a body",
		},
		commitMessage,
	)

}

func TestParseCommitMessageParsesCommitMessageWithFooter(t *testing.T) {

	commitMessage, err := ParseCommitMessage(
		`
feat: my description
with line break

and this is a body
with line break

this is still the body

this is the footer
with a single line break
`,
	)

	assert.Nil(t, err)

	assert.Equal(
		t,
		&CommitMessage{
			ChangeType:  "feat",
			Description: "my description\nwith line break",
			Body:        "and this is a body\nwith line break\n\nthis is still the body",
			Footer:      "this is the footer\nwith a single line break",
		},
		commitMessage,
	)

}

func TestCommitMessage_IsBreakingChangeReturnsFalseIfNoBreakingChangeDescription(t *testing.T) {

	assert.False(
		t,
		(&CommitMessage{
			Body:   "Body without BREAKING CHANGE description // BREAKING CHANGE:",
			Footer: "Footer without BREAKING CHANGE description // BREAKING CHANGE:",
		}).IsBreakingChange(),
	)

}

func TestCommitMessage_IsBreakingChangeReturnsTrueIfBreakingChangeIndicatorExists(t *testing.T) {

	assert.True(
		t,
		(&CommitMessage{
			HasBreakingChangeIndicator: true,
		}).IsBreakingChange(),
	)
}

func TestCommitMessage_IsBreakingChangeReturnsTrueIfBreakingChangeDescriptionInBody(t *testing.T) {

	assert.True(
		t,
		(&CommitMessage{
			Body: "BREAKING CHANGE: commit breaks stuff",
		}).IsBreakingChange(),
	)

	assert.True(
		t,
		(&CommitMessage{
			Body: "BREAKING CHANGES: commit breaks stuff",
		}).IsBreakingChange(),
	)

	assert.True(
		t,
		(&CommitMessage{
			Body: "BREAKING_CHANGE: commit breaks stuff",
		}).IsBreakingChange(),
	)

	assert.True(
		t,
		(&CommitMessage{
			Body: "BREAKING_CHANGES: commit breaks stuff",
		}).IsBreakingChange(),
	)

	assert.True(
		t,
		(&CommitMessage{
			Body: "BREAKING-CHANGE: commit breaks stuff",
		}).IsBreakingChange(),
	)

	assert.True(
		t,
		(&CommitMessage{
			Body: "BREAKING-CHANGES: commit breaks stuff",
		}).IsBreakingChange(),
	)

	assert.True(
		t,
		(&CommitMessage{
			Body: `
Body with breaking change description in second line:
BREAKING CHANGE: commit breaks stuff`,
		}).IsBreakingChange(),
	)

}

func TestCommitMessage_IsBreakingChangeReturnstrueIfBreakingChangeDescriptionInFooter(t *testing.T) {

	assert.True(
		t,
		(&CommitMessage{
			Footer: "BREAKING CHANGE: commit breaks stuff",
		}).IsBreakingChange(),
	)

	assert.True(
		t,
		(&CommitMessage{
			Footer: "BREAKING CHANGES: commit breaks stuff",
		}).IsBreakingChange(),
	)

	assert.True(
		t,
		(&CommitMessage{
			Footer: "BREAKING_CHANGE: commit breaks stuff",
		}).IsBreakingChange(),
	)

	assert.True(
		t,
		(&CommitMessage{
			Footer: "BREAKING_CHANGES: commit breaks stuff",
		}).IsBreakingChange(),
	)

	assert.True(
		t,
		(&CommitMessage{
			Footer: "BREAKING-CHANGE: commit breaks stuff",
		}).IsBreakingChange(),
	)

	assert.True(
		t,
		(&CommitMessage{
			Footer: "BREAKING-CHANGES: commit breaks stuff",
		}).IsBreakingChange(),
	)

	assert.True(
		t,
		(&CommitMessage{
			Footer: `
Footer with breaking change description in second line:
BREAKING CHANGE: commit breaks stuff`,
		}).IsBreakingChange(),
	)

}

func TestCommitMessage_Compare_should_return_0_if_left_is_breaking_change_and_right_too(t *testing.T) {

	assert.Equal(
		t,
		0,
		(&CommitMessage{
			HasBreakingChangeIndicator: true,
		}).Compare(
			&CommitMessage{
				HasBreakingChangeIndicator: true,
			},
		),
	)

}

func TestCommitMessage_Compare_should_return_1_if_left_is_breaking_change_and_right_is_not(t *testing.T) {

	assert.Equal(
		t,
		1,
		(&CommitMessage{
			HasBreakingChangeIndicator: true,
		}).Compare(
			&CommitMessage{
				ChangeType: FEATURE,
			},
		),
	)

	assert.Equal(
		t,
		1,
		(&CommitMessage{
			HasBreakingChangeIndicator: true,
		}).Compare(
			&CommitMessage{
				ChangeType: FIX,
			},
		),
	)

	assert.Equal(
		t,
		1,
		(&CommitMessage{
			HasBreakingChangeIndicator: true,
		}).Compare(
			&CommitMessage{
				ChangeType: CHORE,
			},
		),
	)

}

func TestCommitMessage_Compare_should_return_1_if_left_is_not_breaking_change_but_right_is(t *testing.T) {

	assert.Equal(
		t,
		-1,
		(&CommitMessage{
			ChangeType: FEATURE,
		}).Compare(
			&CommitMessage{
				HasBreakingChangeIndicator: true,
			},
		),
	)

	assert.Equal(
		t,
		-1,
		(&CommitMessage{
			ChangeType: FIX,
		}).Compare(
			&CommitMessage{
				HasBreakingChangeIndicator: true,
			},
		),
	)

	assert.Equal(
		t,
		-1,
		(&CommitMessage{
			ChangeType: CHORE,
		}).Compare(
			&CommitMessage{
				HasBreakingChangeIndicator: true,
			},
		),
	)

}

func TestCommitMessage_Compare_should_return_0_if_left_is_feature_and_right_too(t *testing.T) {

	assert.Equal(
		t,
		0,
		(&CommitMessage{
			ChangeType: FEATURE,
		}).Compare(
			&CommitMessage{
				ChangeType: FEATURE,
			},
		),
	)

}

func TestCommitMessage_Compare_should_return_1_if_left_is_feature_and_right_is_not(t *testing.T) {

	assert.Equal(
		t,
		1,
		(&CommitMessage{
			ChangeType: FEATURE,
		}).Compare(
			&CommitMessage{
				ChangeType: FIX,
			},
		),
	)

	assert.Equal(
		t,
		1,
		(&CommitMessage{
			ChangeType: FEATURE,
		}).Compare(
			&CommitMessage{
				ChangeType: CHORE,
			},
		),
	)

}

func TestCommitMessage_Compare_should_return_1_if_left_is_not_feature_but_right_is(t *testing.T) {

	assert.Equal(
		t,
		-1,
		(&CommitMessage{
			ChangeType: FIX,
		}).Compare(
			&CommitMessage{
				ChangeType: FEATURE,
			},
		),
	)

	assert.Equal(
		t,
		-1,
		(&CommitMessage{
			ChangeType: CHORE,
		}).Compare(
			&CommitMessage{
				ChangeType: FEATURE,
			},
		),
	)

}

func TestCommitMessage_Compare_should_return_0_if_left_is_fix_and_right_too(t *testing.T) {

	assert.Equal(
		t,
		0,
		(&CommitMessage{
			ChangeType: FIX,
		}).Compare(
			&CommitMessage{
				ChangeType: FIX,
			},
		),
	)

}

func TestCommitMessage_Compare_should_return_1_if_left_is_fix_and_right_is_chore(t *testing.T) {

	assert.Equal(
		t,
		1,
		(&CommitMessage{
			ChangeType: FIX,
		}).Compare(
			&CommitMessage{
				ChangeType: CHORE,
			},
		),
	)

}

func TestCommitMessage_Compare_should_return_1_if_left_is_chore_and_right_is_fix(t *testing.T) {

	assert.Equal(
		t,
		-1,
		(&CommitMessage{
			ChangeType: CHORE,
		}).Compare(
			&CommitMessage{
				ChangeType: FIX,
			},
		),
	)

}

func TestCommitMessage_Compare_should_return_0_if_left_is_chore_and_right_is_doc(t *testing.T) {

	assert.Equal(
		t,
		0,
		(&CommitMessage{
			ChangeType: CHORE,
		}).Compare(
			&CommitMessage{
				ChangeType: DOCS,
			},
		),
	)

}
