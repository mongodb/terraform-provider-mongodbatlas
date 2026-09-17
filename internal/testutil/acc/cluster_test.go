package acc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestClusterTier(t *testing.T) {
	testCases := map[string]struct {
		envValue string
		expected string
	}{
		"unset uses M10": {
			envValue: "",
			expected: constant.M10,
		},
		"plain value": {
			envValue: "M50",
			expected: "M50",
		},
		"whitespace padded value": {
			envValue: " M50 ",
			expected: "M50",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv(acc.TestClusterTierEnvName, tc.envValue)
			assert.Equal(t, tc.expected, acc.TestClusterTier())
		})
	}
}
