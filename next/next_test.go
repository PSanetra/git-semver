package next

import (
	"github.com/psanetra/git-semver/semver"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagNameToVersionShouldReturnVersionIfProjectNameMatches(t *testing.T) {
	version := tagNameToVersion("myProject@v1.2.3", "myProject", false)

	assert.Equal(
		t,
		&semver.Version{
			Major:         1,
			Minor:         2,
			Patch:         3,
			PreReleaseTag: []interface{}{},
		},
		version,
	)
}

func TestTagNameToVersionShouldReturnNilIfVersionIsForOtherProject(t *testing.T) {
	version := tagNameToVersion("myProject@v1.2.3", "myOtherProject", false)

	assert.Nil(t, version)
}

func TestTagNameToVersionShouldReturnNilIfTagContainsProjectNameAndProjectNameIsEmpty(t *testing.T) {
	version := tagNameToVersion("myProject@v1.2.3", "", false)

	assert.Nil(t, version)
}

func TestTagNameToVersionShouldReturnNilIfVersionIsForMainProjectAndProjectIsNotMainProject(t *testing.T) {
	version := tagNameToVersion("v1.2.3", "myOtherProject", false)

	assert.Nil(t, version)
}

func TestTagNameToVersionShouldReturnVersionIfTagContainsNoProjectNameAndProjectIsMainProject(t *testing.T) {
	version := tagNameToVersion("v1.2.3", "myMainProject", true)

	assert.Equal(
		t,
		&semver.Version{
			Major:         1,
			Minor:         2,
			Patch:         3,
			PreReleaseTag: []interface{}{},
		},
		version,
	)
}
