package regex_utils

import (
	"regexp"
)

func SubmatchMap(regexp *regexp.Regexp, str string) map[string]string {

	submatches := regexp.FindStringSubmatch(str)

	if submatches == nil {
		return nil
	}

	return SubmatchMapFromSubmatches(regexp, submatches)
}

func SubmatchMapFromSubmatches(regexp *regexp.Regexp, submatches []string) map[string]string {
	groupNames := regexp.SubexpNames()

	submatchMap := make(map[string]string, len(groupNames))

	for i, groupName := range groupNames {
		submatchMap[groupName] = submatches[i]
	}

	return submatchMap

}
