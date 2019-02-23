package regex_utils

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestSubmatchMapShouldReturnMapWhichMapsGroupNamesToMatchedStrings(t *testing.T) {

	regex := regexp.MustCompile(`abc(?P<group1>def(?P<nestedGroup>ghi))(?P<group2>jkl)`)

	matches := SubmatchMap(regex, "abcdefghijkl")

	assert.Equal(t, "defghi", matches["group1"])
	assert.Equal(t, "ghi", matches["nestedGroup"])
	assert.Equal(t, "jkl", matches["group2"])

}

func TestSubmatchMapShouldReturnNilMapIfRegexDidNotMatch(t *testing.T) {

	regex := regexp.MustCompile("xyz")

	matches := SubmatchMap(regex, "abc")

	assert.Nil(t, matches)

}