package acc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestDefaultRegions(t *testing.T) {
	testCases := map[string]struct {
		envValue string
		expected []string
	}{
		"empty uses default list": {
			envValue: "",
			expected: []string{constant.UsWest2, "US_EAST_1", "US_EAST_2"},
		},
		"single region": {
			envValue: "EU_WEST_1",
			expected: []string{"EU_WEST_1"},
		},
		"multiple regions": {
			envValue: "EU_WEST_1,US_EAST_1,US_WEST_2",
			expected: []string{"EU_WEST_1", "US_EAST_1", "US_WEST_2"},
		},
		"spaces and empty entries are dropped": {
			envValue: " EU_WEST_1 , , US_EAST_1,",
			expected: []string{"EU_WEST_1", "US_EAST_1"},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv(acc.TestRegionsEnvName, tc.envValue)
			assert.Equal(t, tc.expected, acc.DefaultRegions())
			assert.Equal(t, tc.expected[0], acc.DefaultRegion())
		})
	}
}
