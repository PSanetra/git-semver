package next

import (
	"github.com/psanetra/git-semver/semver"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_tagNameToVersion_should_return_version(t *testing.T) {
	version := tagNameToVersion("1.2.3")

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

func Test_tagNameToVersion_should_return_version_if_tag_has_v_prefix(t *testing.T) {
	version := tagNameToVersion("v1.2.3")

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
