package semver

import (
	"github.com/psanetra/git-semver/logger"
)

func CompareVersions(v1 *Version, v2 *Version) int {

	if v1 == nil && v2 == nil {
		return 0
	} else if v1 == nil {
		return -1
	} else if v2 == nil {
		return 1
	}

	if v1.Major != v2.Major {
		return v1.Major - v2.Major
	}

	if v1.Minor != v2.Minor {
		return v1.Minor - v2.Minor
	}

	if v1.Patch != v2.Patch {
		return v1.Patch - v2.Patch
	}

	for i := 0; i < len(v1.PreReleaseTag) && i < len(v2.PreReleaseTag); i++ {
		v1TagId := v1.PreReleaseTag[i]
		v2TagId := v2.PreReleaseTag[i]

		comparisonresult := ComparePreReleaseTagIds(v1TagId, v2TagId)

		if comparisonresult != 0 {
			return comparisonresult
		}
	}

	if len(v1.PreReleaseTag) == len(v2.PreReleaseTag) {
		return 0
	} else if len(v1.PreReleaseTag) == 0 {
		return 1
	} else if len(v2.PreReleaseTag) == 0 {
		return -1
	}

	return len(v1.PreReleaseTag) - len(v2.PreReleaseTag)
}

func ComparePreReleaseTagIds(tagId1 interface{}, tagId2 interface{}) int {

	numericV1TagId, v1TagIdIsNumeric := tagId1.(int64)

	numericV2TagId, v2TagIdIsNumeric := tagId2.(int64)

	if v1TagIdIsNumeric != v2TagIdIsNumeric {
		if v1TagIdIsNumeric {
			return -1
		} else {
			return 1
		}
	}

	if v1TagIdIsNumeric {

		if numericV1TagId > numericV2TagId {
			return 1
		} else if numericV2TagId > numericV1TagId {
			return -1
		}

	} else {

		stringV1TagId, ok := tagId1.(string)

		if !ok {
			logger.Logger.Fatalln("Unknown pre-release tag id type")
		}

		stringV2TagId, ok := tagId2.(string)

		if !ok {
			logger.Logger.Fatalln("Unknown pre-release tag id type")
		}

		if stringV1TagId > stringV2TagId {
			return 1
		} else if stringV2TagId > stringV1TagId {
			return -1
		}

	}

	return 0
}
