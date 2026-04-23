package advancedcluster

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestOverrideAttributesWithPrevStateValueMongoDBMajorVersion(t *testing.T) {
	testCases := []struct {
		name     string
		before   string
		after    string
		expected string
	}{
		{
			name:     "keeps previous version when only formatting differs",
			before:   "8",
			after:    "8.0",
			expected: "8",
		},
		{
			name:     "uses API value when major version actually changed",
			before:   "7.0",
			after:    "8.0",
			expected: "8.0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			modelIn := &TFModel{
				MongoDBMajorVersion: types.StringValue(tc.before),
				Labels:              types.MapNull(types.StringType),
				Tags:                types.MapNull(types.StringType),
			}
			modelOut := &TFModel{
				MongoDBMajorVersion: types.StringValue(tc.after),
				Labels:              types.MapNull(types.StringType),
				Tags:                types.MapNull(types.StringType),
			}

			overrideAttributesWithPrevStateValue(modelIn, modelOut)

			if got := modelOut.MongoDBMajorVersion.ValueString(); got != tc.expected {
				t.Fatalf("unexpected mongo_db_major_version: got %q want %q", got, tc.expected)
			}
		})
	}
}
