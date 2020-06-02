package semver

func FindGreatestPreceding(version *Version, versions []*Version, ignorePreReleases bool) *Version {
	var ret *Version

	for _, v := range versions {
		if v == nil ||
			ignorePreReleases && v.IsPreRelease() ||
			version != nil && CompareVersions(v, version) >= 0 {
			continue
		}

		if CompareVersions(v, ret) > 0 {
			ret = v
		}
	}

	return ret
}
