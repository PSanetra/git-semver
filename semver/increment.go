package semver

import "github.com/pkg/errors"

var VersionAlreadyStableError = errors.New("There is already a stable version of this project.")

type Change int

const (
	FIX         Change = 1
	NEW_FEATURE Change = 2
	BREAKING    Change = 3
)

type PreReleaseOptions struct {
	Label string
	AppendCounter bool
}

func (o *PreReleaseOptions) ShouldBePreRelease() bool {
	return o != nil && (o.Label != "" || o.AppendCounter)
}

// Increments a semantic version
func Increment(
	latestRelease Version,
	latestPreRelease *Version,
	shouldBeStable bool,
	highestPriorityChange Change,
	preReleaseOpts *PreReleaseOptions) (Version, error) {

	if latestRelease.IsStable() && !shouldBeStable {
		return Version{}, VersionAlreadyStableError
	}

	newVersion := latestRelease
	newVersion.PreReleaseTag = []interface{}{}

	if shouldBeStable && newVersion.Major < 1 {
		highestPriorityChange = BREAKING
	}

	switch highestPriorityChange {
	case BREAKING:
		newVersion.incrementOnBreakingChange(shouldBeStable)
		break
	case NEW_FEATURE:
		newVersion.incrementMinor()
		break
	case FIX:
		newVersion.incrementPatch()
		break
	}

	if !preReleaseOpts.ShouldBePreRelease() {
		return newVersion, nil
	}

	newPreReleaseTag, err := parsePreReleaseTag(preReleaseOpts.Label)

	if err != nil {
		return Version{}, errors.WithMessage(err, "Could not parse pre-release tag")
	}

	if preReleaseOpts.AppendCounter {
		newPreReleaseTag = append(newPreReleaseTag, int64(1))
	}

	if latestPreRelease == nil ||
		newVersion.Major != latestPreRelease.Major ||
		newVersion.Minor != latestPreRelease.Minor ||
		newVersion.Patch != latestPreRelease.Patch {

		newVersion.PreReleaseTag = newPreReleaseTag
		return newVersion, nil
	}

	if preReleaseOpts.AppendCounter && preReleaseTagsWithCounterAreSimilar(newPreReleaseTag, latestPreRelease.PreReleaseTag) {
		newPreReleaseTag = incrementPreReleaseTagCounter(latestPreRelease.PreReleaseTag)
	}

	newVersion.PreReleaseTag = newPreReleaseTag
	return newVersion, nil
}

func (v *Version) incrementOnBreakingChange(shouldBeStable bool) {
	if shouldBeStable {
		v.Major += 1
		v.Minor = 0
		v.Patch = 0
	} else {
		v.incrementMinor()
	}
}

func (v *Version) incrementMinor() {
	v.Minor += 1
	v.Patch = 0
}

func (v *Version) incrementPatch() {
	v.Patch += 1
}

// Checks if two PreRelease Tags are equal. Will not compare the counter.
func preReleaseTagsWithCounterAreSimilar(tag1 []interface{}, tag2 []interface{}) bool {
	if len(tag1) != len(tag2) {
		return false
	}

	for i := 0; i < (len(tag1) - 1); i++ {
		if tag1[i] != tag2[i] {
			return false
		}
	}

	_, t1eIsInt64 := tag1[len(tag1) - 1].(int64)
	_, t2eIsInt64 := tag2[len(tag1) - 1].(int64)
	return t1eIsInt64 && t2eIsInt64
}

func incrementPreReleaseTagCounter(tag []interface{}) ([]interface{}) {
	newTag := make([]interface{}, 0, len(tag))
	newTag = append(newTag, tag...)

	counter, _ := newTag[len(newTag)-1].(int64)

	newTag[len(newTag)-1] = counter + 1

	return newTag
}
