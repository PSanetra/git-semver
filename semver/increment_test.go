package semver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIncrementShouldSetMajorVersionTo1IfVersionShouldBeStableEvenIfThereWasNoBreakingChange(t *testing.T) {

	newVersion, err := Increment(
		Version{
			Major: 0,
			Minor: 0,
			Patch: 0,
		},
		nil,
		true,
		FIX,
		nil,
	)

	assert.Nil(t, err)
	assert.Equal(
		t,
		Version{
			Major:         1,
			Minor:         0,
			Patch:         0,
			PreReleaseTag: []interface{}{},
		},
		newVersion,
	)

}

func TestIncrementShouldIncrementMajorVersionOnBreakingChange(t *testing.T) {

	newVersion, err := Increment(
		Version{
			Major: 1,
			Minor: 1,
			Patch: 1,
		},
		nil,
		true,
		BREAKING,
		nil,
	)

	assert.Nil(t, err)
	assert.Equal(
		t,
		Version{
			Major:         2,
			Minor:         0,
			Patch:         0,
			PreReleaseTag: []interface{}{},
		},
		newVersion,
	)

}

func TestIncrementShouldIncrementMinorVersionOnBreakingChangeIfVersionShouldBeUnstable(t *testing.T) {

	newVersion, err := Increment(
		Version{
			Major: 0,
			Minor: 1,
			Patch: 1,
		},
		nil,
		false,
		BREAKING,
		nil,
	)

	assert.Nil(t, err)
	assert.Equal(
		t,
		Version{
			Major:         0,
			Minor:         2,
			Patch:         0,
			PreReleaseTag: []interface{}{},
		},
		newVersion,
	)

}

func TestIncrementShouldReturnErrorIfVersionShouldBeUnstableButLatestVersionIsAlreadyStable(t *testing.T) {

	_, err := Increment(
		Version{
			Major: 1,
		},
		nil,
		false,
		BREAKING,
		nil,
	)

	assert.Equal(t, VersionAlreadyStableError, err)

}

func TestIncrementShouldIncrementMinorVersionOnNewFeature(t *testing.T) {

	newVersion, err := Increment(
		Version{
			Major: 1,
			Minor: 1,
			Patch: 1,
		},
		nil,
		true,
		NEW_FEATURE,
		nil,
	)

	assert.Nil(t, err)
	assert.Equal(
		t,
		Version{
			Major:         1,
			Minor:         2,
			Patch:         0,
			PreReleaseTag: []interface{}{},
		},
		newVersion,
	)

}

func TestIncrementShouldIncrementPatchVersionOnFix(t *testing.T) {

	newVersion, err := Increment(
		Version{
			Major: 1,
			Minor: 1,
			Patch: 1,
		},
		nil,
		true,
		FIX,
		nil,
	)

	assert.Nil(t, err)
	assert.Equal(
		t,
		Version{
			Major:         1,
			Minor:         1,
			Patch:         2,
			PreReleaseTag: []interface{}{},
		},
		newVersion,
	)

}

func TestIncrementShouldIncrementVersionAndApplyLabelOnPreRelease(t *testing.T) {

	newVersion, err := Increment(
		Version{},
		nil,
		true,
		BREAKING,
		&PreReleaseOptions{
			Label: "alpha.2018-12-31",
		},
	)

	assert.Nil(t, err)
	assert.Equal(
		t,
		Version{
			Major:         1,
			PreReleaseTag: []interface{}{"alpha", "2018-12-31"},
		},
		newVersion,
	)

}

func TestIncrementShouldIncrementVersionAndApplyLabelAndAppendCounterOnPreRelease(t *testing.T) {

	newVersion, err := Increment(
		Version{},
		nil,
		true,
		BREAKING,
		&PreReleaseOptions{
			Label:         "alpha.2018-12-31",
			AppendCounter: true,
		},
	)

	assert.Nil(t, err)
	assert.Equal(
		t,
		Version{
			Major:         1,
			PreReleaseTag: []interface{}{"alpha", "2018-12-31", int64(1)},
		},
		newVersion,
	)

}

func TestIncrementShouldIncrementVersionAndApplyLabelAndIncrementCounterOnExistingPreRelease(t *testing.T) {

	newVersion, err := Increment(
		Version{},
		&Version{
			Major:         1,
			PreReleaseTag: []interface{}{"alpha", "2018-12-31", int64(99)},
		},
		true,
		BREAKING,
		&PreReleaseOptions{
			Label:         "alpha.2018-12-31",
			AppendCounter: true,
		},
	)

	assert.Nil(t, err)
	assert.Equal(
		t,
		Version{
			Major:         1,
			PreReleaseTag: []interface{}{"alpha", "2018-12-31", int64(100)},
		},
		newVersion,
	)

}

func TestIncrementShouldCompareLengthsOfPreReleaseTagsIfExistingPreReleaseTagIsLongerThanNewPreReleaseTagAndCounterShouldBeAppended(t *testing.T) {

	newVersion, err := Increment(
		Version{},
		&Version{
			Major:         1,
			PreReleaseTag: []interface{}{"alpha", int64(1), int64(1)},
		},
		true,
		BREAKING,
		&PreReleaseOptions{
			Label:         "alpha",
			AppendCounter: true,
		},
	)

	assert.Nil(t, err)
	assert.Equal(
		t,
		Version{
			Major:         1,
			PreReleaseTag: []interface{}{"alpha", int64(1)},
		},
		newVersion,
	)

}

func TestIncrementShouldCompareLengthsOfPreReleaseTagsIfExistingPreReleaseTagIsShorterThanNewPreReleaseTagAndCounterShouldBeAppended(t *testing.T) {

	newVersion, err := Increment(
		Version{},
		&Version{
			Major:         1,
			PreReleaseTag: []interface{}{"alpha", int64(1)},
		},
		true,
		BREAKING,
		&PreReleaseOptions{
			Label:         "alpha.1",
			AppendCounter: true,
		},
	)

	assert.Nil(t, err)
	assert.Equal(
		t,
		Version{
			Major:         1,
			PreReleaseTag: []interface{}{"alpha", int64(1), int64(1)},
		},
		newVersion,
	)

}

func TestIncrementShouldIncrementVersionAndIgnoreExistingPreReleaseIfPreReleaseTagsDoNotMatch(t *testing.T) {

	newVersion, err := Increment(
		Version{},
		&Version{
			Major:         1,
			PreReleaseTag: []interface{}{"alpha", "2018-12-31", int64(99)},
		},
		true,
		BREAKING,
		&PreReleaseOptions{
			Label:         "alpha.2019-01-01",
			AppendCounter: true,
		},
	)

	assert.Nil(t, err)
	assert.Equal(
		t,
		Version{
			Major:         1,
			PreReleaseTag: []interface{}{"alpha", "2019-01-01", int64(1)},
		},
		newVersion,
	)

}

func TestIncrementShouldIncrementVersionAndIgnoreExistingPreReleaseIfExistingPreReleaseHasDifferentMajorVersion(t *testing.T) {

	newVersion, err := Increment(
		Version{},
		&Version{
			Major:         0,
			Minor:         1,
			Patch:         0,
			PreReleaseTag: []interface{}{"alpha", "2018-12-31", int64(99)},
		},
		true,
		BREAKING,
		&PreReleaseOptions{
			Label:         "alpha.2018-12-31",
			AppendCounter: true,
		},
	)

	assert.Nil(t, err)
	assert.Equal(
		t,
		Version{
			Major:         1,
			Minor:         0,
			Patch:         0,
			PreReleaseTag: []interface{}{"alpha", "2018-12-31", int64(1)},
		},
		newVersion,
	)

}

func TestIncrementShouldIncrementVersionAndIgnoreExistingPreReleaseIfExistingPreReleaseHasDifferentMinorVersion(t *testing.T) {

	newVersion, err := Increment(
		Version{},
		&Version{
			Major:         0,
			Minor:         0,
			Patch:         1,
			PreReleaseTag: []interface{}{"alpha", "2018-12-31", int64(99)},
		},
		false,
		NEW_FEATURE,
		&PreReleaseOptions{
			Label:         "alpha.2018-12-31",
			AppendCounter: true,
		},
	)

	assert.Nil(t, err)
	assert.Equal(
		t,
		Version{
			Major:         0,
			Minor:         1,
			Patch:         0,
			PreReleaseTag: []interface{}{"alpha", "2018-12-31", int64(1)},
		},
		newVersion,
	)

}

func TestIncrementShouldIncrementVersionAndIgnoreExistingPreReleaseIfExistingPreReleaseHasDifferentPatchVersion(t *testing.T) {

	newVersion, err := Increment(
		Version{},
		&Version{
			Major:         0,
			Minor:         0,
			Patch:         0,
			PreReleaseTag: []interface{}{"alpha", "2018-12-31", int64(99)},
		},
		false,
		FIX,
		&PreReleaseOptions{
			Label:         "alpha.2018-12-31",
			AppendCounter: true,
		},
	)

	assert.Nil(t, err)
	assert.Equal(
		t,
		Version{
			Major:         0,
			Minor:         0,
			Patch:         1,
			PreReleaseTag: []interface{}{"alpha", "2018-12-31", int64(1)},
		},
		newVersion,
	)

}
