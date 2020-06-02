package semver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindPrecedingShouldReturnNilArgsAreNil(t *testing.T) {

	result := FindGreatestPreceding(nil, nil, false)

	assert.Nil(t, result)

}

func TestFindPrecedingShouldReturnGreatestIfVersionIsNil(t *testing.T) {

	result := FindGreatestPreceding(
		nil,
		[]*Version{
			{1, 2, 3, nil},
			{3, 2, 1, nil},
			{2, 1, 3, nil},
		},
		false,
	)

	assert.Equal(t, 3, result.Major)
}

func TestFindPrecedingShouldReturnNilIfListIsNil(t *testing.T) {

	result := FindGreatestPreceding(&Version{1, 2, 3, nil}, nil, false)

	assert.Nil(t, result)
}

func TestFindPrecedingShouldReturnGreatestPrecedingOfVersion(t *testing.T) {

	result := FindGreatestPreceding(
		&Version{2, 1, 3, nil},
		[]*Version{
			{1, 2, 3, nil},
			{3, 2, 1, nil},
			{2, 1, 3, nil},
		},
		false,
	)

	assert.Equal(t, 1, result.Major)
}

func TestFindPrecedingWithoutIgnoringPreReleasesShouldReturnGreatestPrecedingPreReleaseOfVersion(t *testing.T) {

	result := FindGreatestPreceding(
		&Version{3, 2, 1, nil},
		[]*Version{
			{1, 2, 3, nil},
			{3, 2, 1, nil},
			{2, 1, 3, nil},
			{3, 2, 1, []interface{}{"alpha"}},
		},
		false,
	)

	assert.Equal(t, 3, result.Major)
	assert.Equal(t, "alpha", result.PreReleaseTag[0])
}

func TestFindPrecedingWithIgnoringPreReleasesShouldReturnGreatestPrecedingNonPreReleaseOfVersion(t *testing.T) {

	result := FindGreatestPreceding(
		&Version{3, 2, 1, nil},
		[]*Version{
			{1, 2, 3, nil},
			{3, 2, 1, nil},
			{2, 1, 3, nil},
			{3, 2, 1, []interface{}{"alpha"}},
		},
		true,
	)

	assert.Equal(t, 2, result.Major)
}
