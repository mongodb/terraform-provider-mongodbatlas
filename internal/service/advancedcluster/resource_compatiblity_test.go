package advancedcluster_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/stretchr/testify/assert"
)

func TestAdvancedCluster_overrideMongoDBMajorVersion(t *testing.T) {
	testCases := []struct {
		name     string
		before   string
		after    string
		expected bool
	}{
		{
			name:     "keeps previous version when only formatting differs",
			before:   "8",
			after:    "8.0",
			expected: true,
		},
		{
			name:     "uses API value when major version actually changed",
			before:   "7.0",
			after:    "8.0",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, advancedcluster.ShouldUsePreviousMongoDBMajorVersion(tc.before, tc.after))
		})
	}
}
