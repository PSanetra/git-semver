package semver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompareVersionsReturns0IfV1AndV2AreNil(t *testing.T) {

	result := CompareVersions(nil, nil)

	assert.Equal(t, 0, result)
	assert.True(t, result == 0)

}

func TestCompareVersionsReturnsGreater0IfOnlyV2IsNil(t *testing.T) {

	result := CompareVersions(
		&Version{
			0,
			0,
			0,
			[]interface{}{},
		},
		nil,
	)

	assert.True(t, result > 0)

}

func TestCompareVersionsReturnsLess0IfOnlyV1IsNil(t *testing.T) {

	result := CompareVersions(
		nil,
		&Version{
			0,
			0,
			0,
			[]interface{}{},
		},
	)

	assert.True(t, result < 0)

}

func TestCompareVersionsReturnsGreaterThan0IfV1HasGreaterMajorVersion(t *testing.T) {

	result := CompareVersions(
		&Version{
			2,
			3,
			4,
			[]interface{}{},
		},
		&Version{
			1,
			5,
			6,
			[]interface{}{},
		},
	)

	assert.True(t, result > 0)

}

func TestCompareVersionsReturnsLessThan0IfV2HasGreaterMajorVersion(t *testing.T) {

	result := CompareVersions(
		&Version{
			1,
			5,
			6,
			[]interface{}{},
		},
		&Version{
			2,
			3,
			4,
			[]interface{}{},
		},
	)

	assert.True(t, result < 0)

}

func TestCompareVersionsReturnsGreaterThan0IfV1HasGreaterMinorVersion(t *testing.T) {

	result := CompareVersions(
		&Version{
			1,
			2,
			3,
			[]interface{}{},
		},
		&Version{
			1,
			1,
			4,
			[]interface{}{},
		},
	)

	assert.True(t, result > 0)

}

func TestCompareVersionsReturnsLessThan0IfV2HasGreaterMinorVersion(t *testing.T) {

	result := CompareVersions(
		&Version{
			1,
			1,
			4,
			[]interface{}{},
		},
		&Version{
			1,
			2,
			3,
			[]interface{}{},
		},
	)

	assert.True(t, result < 0)

}

func TestCompareVersionsReturnsGreaterThan0IfV1HasGreaterPatchVersion(t *testing.T) {

	result := CompareVersions(
		&Version{
			1,
			1,
			2,
			[]interface{}{},
		},
		&Version{
			1,
			1,
			1,
			[]interface{}{},
		},
	)

	assert.True(t, result > 0)

}

func TestCompareVersionsReturnsLessThan0IfV2HasGreaterPatchVersion(t *testing.T) {

	result := CompareVersions(
		&Version{
			1,
			1,
			1,
			[]interface{}{},
		},
		&Version{
			1,
			1,
			2,
			[]interface{}{},
		},
	)

	assert.True(t, result < 0)

}

func TestCompareVersionsReturnsGreaterThan0IfV1HasPreReleaseTagWithHigherPrecedence(t *testing.T) {

	result := CompareVersions(
		&Version{
			1,
			1,
			1,
			[]interface{}{
				"beta",
				int64(1),
			},
		},
		&Version{
			1,
			1,
			1,
			[]interface{}{
				"alpha",
				int64(99),
			},
		},
	)

	assert.True(t, result > 0)

}

func TestCompareVersionsReturnsLessThan0IfV2HasPreReleaseTagWithHigherPrecedence(t *testing.T) {

	result := CompareVersions(
		&Version{
			1,
			1,
			1,
			[]interface{}{
				"alpha",
				int64(99),
			},
		},
		&Version{
			1,
			1,
			1,
			[]interface{}{
				"beta",
				int64(1),
			},
		},
	)

	assert.True(t, result < 0)

}

// https://semver.org/#spec-item-11
func TestCompareVersionsLongerPreReleaseTagsHaveHigherPrecedence(t *testing.T) {

	result := CompareVersions(
		&Version{
			1,
			1,
			1,
			[]interface{}{
			},
		},
		&Version{
			1,
			1,
			1,
			[]interface{}{
				"alpha",
			},
		},
	)

	assert.True(t, result < 0)

	result = CompareVersions(
		&Version{
			1,
			1,
			1,
			[]interface{}{
				"alpha",
			},
		},
		&Version{
			1,
			1,
			1,
			[]interface{}{
			},
		},
	)

	assert.True(t, result > 0)

}

func TestCompareVersionsReturns0IfV1AndV2AreEqual(t *testing.T) {

	result := CompareVersions(
		&Version{
			1,
			2,
			3,
			[]interface{}{
				"alpha",
				int64(1),
			},
		},
		&Version{
			1,
			2,
			3,
			[]interface{}{
				"alpha",
				int64(1),
			},
		},
	)

	assert.Equal(t, result, 0)

}

func TestComparePreReleaseTagIdsReturns0IfEqual(t *testing.T) {

	assert.Equal(t, ComparePreReleaseTagIds(int64(1), int64(1)), 0)
	assert.Equal(t, ComparePreReleaseTagIds("abc", "abc"), 0)

}

func TestComparePreReleaseTagIdsReturnsGreaterThan0IfTag1IsGreater(t *testing.T) {

	assert.True(t, ComparePreReleaseTagIds(int64(2), int64(1)) > 0)
	assert.True(t, ComparePreReleaseTagIds("xyz", "abc") > 0)

}

func TestComparePreReleaseTagIdsReturnsLessThan0IfTag2IsGreater(t *testing.T) {

	assert.True(t, ComparePreReleaseTagIds(int64(1), int64(2)) < 0)
	assert.True(t, ComparePreReleaseTagIds("abc", "xyz") < 0)

}

// https://semver.org/#spec-item-11
func TestComparePreReleaseTagIdsReturnsGreaterThan0IfTag1IsStringAndTag2IsInt64(t *testing.T) {

	assert.True(t, ComparePreReleaseTagIds("abc", int64(9999)) > 0)

}

// https://semver.org/#spec-item-11
func TestComparePreReleaseTagIdsReturnsLessThan0IfTag2IsStringAndTag1IsInt64(t *testing.T) {

	assert.True(t, ComparePreReleaseTagIds(int64(9999), "abc") < 0)

}
