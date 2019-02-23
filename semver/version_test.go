package semver

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseVersion(t *testing.T) {

	version, err := ParseVersion("1.2.3")

	assert.NoError(t, err)
	assert.Equal(t, version, &Version{
		Major:         1,
		Minor:         2,
		Patch:         3,
		PreReleaseTag: []interface{}{},
	})

}

func TestParseVersionWithVPrefix(t *testing.T) {

	version, err := ParseVersion("v1.2.3")

	assert.NoError(t, err)
	assert.Equal(t, version, &Version{
		Major:         1,
		Minor:         2,
		Patch:         3,
		PreReleaseTag: []interface{}{},
	})

}

func TestParseVersionWithPreReleaseTag(t *testing.T) {

	version, err := ParseVersion("1.2.3-alpha")

	assert.NoError(t, err)
	assert.Equal(t, version, &Version{
		Major: 1,
		Minor: 2,
		Patch: 3,
		PreReleaseTag: []interface{}{
			"alpha",
		},
	})

}

func TestParseVersionWithNumericIdInPreReleaseTag(t *testing.T) {

	version, err := ParseVersion("1.2.3-alpha.4")

	assert.NoError(t, err)
	assert.Equal(t, version, &Version{
		Major: 1,
		Minor: 2,
		Patch: 3,
		PreReleaseTag: []interface{}{
			"alpha",
			int64(4),
		},
	})

}

func TestParseVersionIgnoresBuildMetadata(t *testing.T) {

	version, err := ParseVersion("1.2.3+mymetadata")

	assert.NoError(t, err)
	assert.Equal(t, version, &Version{
		Major:         1,
		Minor:         2,
		Patch:         3,
		PreReleaseTag: []interface{}{},
	})

}

func TestParseVersionReturnsErrorOnInvalidSyntax(t *testing.T) {

	version, err := ParseVersion("1.2.3Invalid")

	assert.Error(t, err)
	assert.Nil(t, version)

}

func TestVersionRegexOnValidStrings(t *testing.T) {

	// source: https://github.com/semver/semver/issues/232#issuecomment-430813095
	validStrings := []string{
		"0.0.4",
		"1.2.3",
		"10.20.30",
		"1.1.2-prerelease+meta",
		"1.1.2+meta",
		"1.1.2+meta-valid",
		"1.0.0-alpha",
		"1.0.0-beta",
		"1.0.0-alpha.beta",
		"1.0.0-alpha.beta.1",
		"1.0.0-alpha.1",
		"1.0.0-alpha0.valid",
		"1.0.0-alpha.0valid",
		"1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay",
		"1.0.0-rc.1+build.1",
		"2.0.0-rc.1+build.123",
		"1.2.3-beta",
		"10.2.3-DEV-SNAPSHOT",
		"1.2.3-SNAPSHOT-123",
		"1.0.0",
		"2.0.0",
		"1.1.7",
		"2.0.0+build.1848",
		"2.0.1-alpha.1227",
		"1.0.0-alpha+beta",
		"1.2.3----RC-SNAPSHOT.12.9.1--.12+788",
		"1.2.3----R-S.12.9.1--.12+meta",
		"1.2.3----RC-SNAPSHOT.12.9.1--.12",
		"1.0.0+0.build.1-rc.10000aaa-kk-0.1",
		"99999999999999999999999.999999999999999999.99999999999999999",
		"1.0.0-0A.is.legal",
	}

	for _, str := range validStrings {

		assert.True(t, VersionRegex.Match([]byte(str)), fmt.Sprintf("\"%s\" did not match", str))

	}

}

func TestVersionRegexOnInvalidStrings(t *testing.T) {

	// source: https://github.com/semver/semver/issues/232#issuecomment-430813095
	invalidStrings := []string{
		"1",
		"1.2",
		"1.2.3-0123",
		"1.2.3-0123.0123",
		"1.1.2+.123",
		"+invalid",
		"-invalid",
		"-invalid+invalid",
		"-invalid.01",
		"alpha",
		"alpha.beta",
		"alpha.beta.1",
		"alpha.1",
		"alpha+beta",
		"alpha_beta",
		"alpha.",
		"alpha..",
		"beta",
		"1.0.0-alpha_beta",
		"-alpha.",
		"1.0.0-alpha..",
		"1.0.0-alpha..1",
		"1.0.0-alpha...1",
		"1.0.0-alpha....1",
		"1.0.0-alpha.....1",
		"1.0.0-alpha......1",
		"1.0.0-alpha.......1",
		"01.1.1",
		"1.01.1",
		"1.1.01",
		"1.2",
		"1.2.3.DEV",
		"1.2-SNAPSHOT",
		"1.2.31.2.3----RC-SNAPSHOT.12.09.1--..12+788",
		"1.2-RC-SNAPSHOT",
		"-1.0.3-gamma+b7718",
		"+justmeta",
		"9.8.7+meta+meta",
		"9.8.7-whatever+meta+meta",
		"99999999999999999999999.999999999999999999.99999999999999999----RC-SNAPSHOT.12.09.1--------------------------------..12",
	}

	for _, str := range invalidStrings {

		assert.False(t, VersionRegex.Match([]byte(str)), fmt.Sprintf("\"%s\" did match", str))

	}

}

func TestVersionIsStableIfMajorVersionIsNot0(t *testing.T) {

	v := &Version{
		Major: 1,
	}

	assert.True(
		t,
		v.IsStable(),
	)

}

func TestVersionIsUnstableIfMajorVersionIs0(t *testing.T) {

	v := &Version{
		Major: 0,
	}

	assert.False(
		t,
		v.IsStable(),
	)

}
