package cloudbackupsnapshotrestorejob_test

import (
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudbackupsnapshotrestorejob"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312016/admin"
)

func TestFlattenDesiredTimestamp(t *testing.T) {
	sampleDate := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	sampleIncrement := 42

	testCases := []struct {
		name     string
		input    *admin.ApiBSONTimestamp
		expected []map[string]any
	}{
		{
			name:     "nil input returns nil",
			input:    nil,
			expected: nil,
		},
		{
			name: "populated timestamp",
			input: &admin.ApiBSONTimestamp{
				Date:      &sampleDate,
				Increment: &sampleIncrement,
			},
			expected: []map[string]any{
				{
					"date":      conversion.StringPtr(sampleDate.Format(time.RFC3339)),
					"increment": sampleIncrement,
				},
			},
		},
		{
			name: "nil date and increment",
			input: &admin.ApiBSONTimestamp{
				Date:      nil,
				Increment: nil,
			},
			expected: []map[string]any{
				{
					"date":      (*string)(nil),
					"increment": 0,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := cloudbackupsnapshotrestorejob.FlattenDesiredTimestamp(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
