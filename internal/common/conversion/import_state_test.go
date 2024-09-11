package conversion_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	projectID           = "123456789012345678901234"
	projectIDInvalid    = "invalid_project_id"
	clusterName         = "clusterName"
	clusterNameWithDash = "cluster-name"
	clusterNameInvalid  = "_invalidClusterName"
)

func TestIDWithProjectIDClusterName(t *testing.T) {
	testCases := map[string]struct {
		projectID           string
		clusterName         string
		expectedID          string
		expectedErrContains string
	}{
		"valid": {
			projectID:   projectID,
			clusterName: clusterName,
			expectedID:  projectID + "-" + clusterName,
		},
		"valid cluster name with dash": {
			projectID:   projectID,
			clusterName: clusterNameWithDash,
			expectedID:  projectID + "-" + clusterNameWithDash,
		},

		"invalid project_id": {
			projectID:           projectIDInvalid,
			clusterName:         clusterName,
			expectedErrContains: projectIDInvalid,
		},
		"invalid cluster_name": {
			projectID:           projectID,
			clusterName:         clusterNameInvalid,
			expectedErrContains: clusterNameInvalid,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			id, err := conversion.IDWithProjectIDClusterName(tc.projectID, tc.clusterName)
			assert.Equal(t, tc.expectedErrContains == "", err == nil)
			if err == nil {
				assert.Equal(t, tc.expectedID, id)
			} else {
				assert.Contains(t, err.Error(), tc.expectedErrContains)
			}
		})
	}
}

func TestValidateProjectID(t *testing.T) {
	require.NoError(t, conversion.ValidateProjectID(projectID))
	require.Error(t, conversion.ValidateProjectID(projectIDInvalid))
}

func TestValidateClusterName(t *testing.T) {
	require.NoError(t, conversion.ValidateClusterName(clusterName))
	require.Error(t, conversion.ValidateClusterName(clusterNameInvalid))
}
