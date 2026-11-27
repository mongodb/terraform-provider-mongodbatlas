package federateddatabaseinstance_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/federateddatabaseinstance"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitDataFederatedInstanceImportID(t *testing.T) {
	tests := map[string]struct {
		importID    string
		projectID   string
		name        string
		s3Bucket    string
		legacy      bool
		errorString bool
	}{
		"slash format": {
			importID:  "project/name",
			projectID: "project",
			name:      "name",
		},
		"slash format with bucket": {
			importID:  "project/name/bucket",
			projectID: "project",
			name:      "name",
			s3Bucket:  "bucket",
		},
		"legacy format": {
			importID:  "project--name",
			projectID: "project",
			name:      "name",
			legacy:    true,
		},
		"legacy format with bucket": {
			importID:  "project--name--bucket",
			projectID: "project",
			name:      "name",
			s3Bucket:  "bucket",
			legacy:    true,
		},
		"invalid format": {
			importID:    "project--name--bucket--extra",
			errorString: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			projectID, instanceName, s3Bucket, legacy, err := federateddatabaseinstance.SplitDataFederatedInstanceImportID(tc.importID)
			if tc.errorString {
				require.Error(t, err)
				assert.False(t, legacy)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.projectID, projectID)
			assert.Equal(t, tc.name, instanceName)
			assert.Equal(t, tc.s3Bucket, s3Bucket)
			assert.Equal(t, tc.legacy, legacy)
		})
	}
}
