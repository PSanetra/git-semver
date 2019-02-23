package semver

import (
	"fmt"
	"strings"
)

func (v *Version) ToString() string {
	releaseVersion := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)

	if len(v.PreReleaseTag) > 0 {
		stringTagElements := make([]string, 0, len(v.PreReleaseTag))

		for _, tagElement := range v.PreReleaseTag {
			stringTagElements = append(stringTagElements, fmt.Sprintf("%v", tagElement))
		}

		return fmt.Sprintf("%s-%s", releaseVersion, strings.Join(stringTagElements, "."))
	}

	return releaseVersion
}
