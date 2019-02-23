package semver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersionToStringShouldReturnWellFormedSemanticVersionStringWithoutPreReleaseTag(t *testing.T)  {

	v := &Version{
		Major: 1,
		Minor: 2,
		Patch: 3,
	}

	assert.Equal(t, "1.2.3", v.ToString())

}

func TestVersionToStringShouldReturnWellFormedSemanticVersionStringWithPreReleaseTag(t *testing.T)  {

	v := &Version{
		Major: 1,
		Minor: 2,
		Patch: 3,
		PreReleaseTag: []interface{}{ "alpha", int64(123) },
	}

	assert.Equal(t, "1.2.3-alpha.123", v.ToString())

}
