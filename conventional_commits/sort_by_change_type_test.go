package conventional_commits

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestSort_ByChangeTypeDesc(t *testing.T) {

	messages := ByChangeTypeDesc{
		{ChangeType: FIX},
		{ChangeType: PERF},
		{ChangeType: FEATURE},
	}

	sort.Sort(messages)

	assert.Equal(t, FEATURE, messages[0].ChangeType)
	assert.Equal(t, FIX, messages[1].ChangeType)
	assert.Equal(t, PERF, messages[2].ChangeType)

}
