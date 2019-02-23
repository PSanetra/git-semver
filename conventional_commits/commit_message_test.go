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
			ChangeType:  "feat",
			Description: "my description",
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
			ChangeType:  "feat",
			Description: "my description",
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
			ChangeType:  "feat",
			Description: "my description",
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
				ChangeType:       "feat",
				Description:      "my description\nwith line break",
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
				ChangeType:       "feat",
				Description:      "my description\nwith line break",
				Body:             "and this is a body",
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
				ChangeType:       "feat",
				Description:      "my description\nwith line break",
				Body:             "and this is a body\nwith line break\n\nthis is still the body",
				Footer:           "this is the footer\nwith a single line break",
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
