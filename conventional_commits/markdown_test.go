package conventional_commits

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_markdown(t *testing.T) {

	result := ToMarkdown([]*ConventionalCommitMessage{
		{
			ChangeType:             FEATURE,
			Scope:                  "some_component",
			ContainsBreakingChange: true,
			Description:            "Add some feature",
			Body:                   "Lorem ipsum...",
			Footers: map[string][]string{
				"BREAKING CHANGE": {
					`There is a breaking change in some API.
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque facilisis neque nec fermentum placerat. 
Integer placerat leo sed leo ullamcorper, nec fermentum tortor tincidunt. 

Pellentesque blandit justo quis mauris gravida, quis mollis nunc maximus. Nulla a massa vitae urna mollis tincidunt. 
Praesent condimentum pellentesque convallis. 

Mauris vitae risus vel lorem luctus rutrum. 
Phasellus neque nibh, posuere eu nibh nec, feugiat gravida sem. Aliquam posuere sit amet diam ut ultrices. 
Nunc tincidunt odio quis ipsum aliquam, ut posuere enim sollicitudin. Pellentesque eu erat id justo semper laoreet.`,
				},
			},
		},
		{
			ChangeType:             FIX,
			Scope:                  "some_component",
			ContainsBreakingChange: true,
			Description:            "Fix some issue",
			Body:                   "Lorem ipsum...",
			Footers: map[string][]string{
				"BREAKING CHANGE": {
					`There is another breaking change in some API.
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque facilisis neque nec fermentum placerat. 
Integer placerat leo sed leo ullamcorper, nec fermentum tortor tincidunt.`,
				},
			},
		},
		{
			ChangeType:  FIX,
			Scope:       "some_component",
			Description: "Fix another issue",
		},
		{
			ChangeType:  FIX,
			Description: "Fix without scope",
		},
		{
			ChangeType:             FIX,
			Scope:                  "some_component",
			ContainsBreakingChange: true,
			Description:            "Fix with breaking change, but without separate BREAKING CHANGE description.",
			Body:                   "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque facilisis neque nec fermentum placerat.",
		},
		{
			ChangeType:  CHORE,
			Description: "Edit README.md",
		},
		{
			ChangeType:  PERF,
			Description: "Improve performance",
		},
		{
			ChangeType:  STYLE,
			Description: "go fmt",
		},
		{
			ChangeType:  REFACTOR,
			Description: "Refactor something",
		},
		{
			ChangeType:  CI,
			Description: "Fix some pipeline",
		},
		{
			ChangeType:  DOCS,
			Description: "Edit some docs",
		},
	})

	assert.Equal(
		t,
		`### BREAKING CHANGES

* **some_component** There is a breaking change in some API.
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque facilisis neque nec fermentum placerat. 
Integer placerat leo sed leo ullamcorper, nec fermentum tortor tincidunt. 

Pellentesque blandit justo quis mauris gravida, quis mollis nunc maximus. Nulla a massa vitae urna mollis tincidunt. 
Praesent condimentum pellentesque convallis. 

Mauris vitae risus vel lorem luctus rutrum. 
Phasellus neque nibh, posuere eu nibh nec, feugiat gravida sem. Aliquam posuere sit amet diam ut ultrices. 
Nunc tincidunt odio quis ipsum aliquam, ut posuere enim sollicitudin. Pellentesque eu erat id justo semper laoreet.
* **some_component** There is another breaking change in some API.
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque facilisis neque nec fermentum placerat. 
Integer placerat leo sed leo ullamcorper, nec fermentum tortor tincidunt.
* **some_component** Fix with breaking change, but without separate BREAKING CHANGE description.

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque facilisis neque nec fermentum placerat.
### Features

* **some_component** Add some feature
Lorem ipsum...

### Bug Fixes

* **some_component** Fix some issue
Lorem ipsum...
* **some_component** Fix another issue
* Fix without scope
`,
		result,
	)

}
