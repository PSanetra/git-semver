package regex_utils

import (
	"regexp"
)

func SubmatchMap(regexp *regexp.Regexp, str string) map[string]string {

	submatches := regexp.FindStringSubmatch(str)

	if submatches == nil {
		return nil
	}

	groupNames := regexp.SubexpNames()

	submatchMap := make(map[string]string, len(groupNames))

	for i, groupName := range groupNames {
		submatchMap[groupName] = submatches[i]
	}

	return submatchMap

}
