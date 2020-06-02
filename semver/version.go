package semver

import (
	"github.com/pkg/errors"
	"github.com/psanetra/git-semver/regex_utils"
	"regexp"
	"strconv"
	"strings"
)

// source: https://github.com/semver/semver/issues/232#issuecomment-430840155
var VersionRegex = regexp.MustCompile("^v?(?P<Major>0|[1-9]\\d*)\\.(?P<Minor>0|[1-9]\\d*)\\.(?P<Patch>0|[1-9]\\d*)(?P<PreReleaseTagWithSeparator>-(?P<PreReleaseTag>(0|[1-9]\\d*|\\d*[A-Za-z-][\\dA-Za-z-]*)(\\.(0|[1-9]\\d*|\\d*[A-Za-z-][\\dA-Za-z-]*))*))?(?P<BuildMetadataTagWithSeparator>\\+(?P<BuildMetadataTag>[\\dA-Za-z-]+(\\.[\\dA-Za-z-]*)*))?$")

var VersionParsingError = errors.New("Could not parse version")

var EmptyVersion = Version{}

type Version struct {
	Major int
	Minor int
	Patch int
	// PreReleaseTag array can contain strings and int64s
	PreReleaseTag []interface{}
}

func ParseVersion(str string) (*Version, error) {

	submatches := regex_utils.SubmatchMap(VersionRegex, str)

	if submatches == nil {
		return nil, VersionParsingError
	}

	major, err := strconv.Atoi(submatches["Major"])

	if err != nil {
		return nil, err
	}

	minor, err := strconv.Atoi(submatches["Minor"])

	if err != nil {
		return nil, err
	}

	patch, err := strconv.Atoi(submatches["Patch"])

	if err != nil {
		return nil, err
	}

	preReleaseTagStr := submatches["PreReleaseTag"]

	preReleaseTag, err := parsePreReleaseTag(preReleaseTagStr)

	if err != nil {
		return nil, errors.WithMessage(err, "Could not parse pre-release tag")
	}

	return &Version{
		Major:         major,
		Minor:         minor,
		Patch:         patch,
		PreReleaseTag: preReleaseTag,
	}, nil

}

func parsePreReleaseTag(str string) ([]interface{}, error) {

	if len(str) == 0 {
		return []interface{}{}, nil
	}

	parts := strings.Split(str, ".")

	preReleaseTag := make([]interface{}, 0, len(parts)+1)

	for _, part := range parts {

		parsedPart, err := strconv.ParseInt(part, 10, 64)

		if err != nil {
			numErr, ok := err.(*strconv.NumError)

			if !ok || numErr.Err != strconv.ErrSyntax {
				return nil, errors.WithMessage(err, "Could not parse part '"+part+"' as int64")
			}

			preReleaseTag = append(preReleaseTag, part)
		} else {
			preReleaseTag = append(preReleaseTag, parsedPart)
		}

	}

	return preReleaseTag, nil
}

func (v *Version) IsStable() bool {
	return v.Major != 0
}

func (v *Version) IsPreRelease() bool {
	return len(v.PreReleaseTag) > 0
}
