package acc_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/assert"
)

const (
	projectID          = "123456789012345678901234"
	projectIDInvalid   = "invalid_project_id"
	clusterName        = "clusterName"
	clusterNameInvalid = "_invalidClusterName"
)

func TestIDWithProjectIDClusterName(t *testing.T) {
	testCases := map[string]struct {
		projectID         string
		clusterName       string
		expectedID        string
		expectedHasErrors bool
	}{
		"valid": {
			projectID:   projectID,
			clusterName: clusterName,
			expectedID:  projectID + "-" + clusterName,
		},
		"invalid project_id": {
			projectID:         projectIDInvalid,
			clusterName:       clusterName,
			expectedHasErrors: true,
		},
		"invalid cluster_name": {
			projectID:         projectID,
			clusterName:       clusterNameInvalid,
			expectedHasErrors: true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			id, err := acc.IDWithProjectIDClusterName(tc.projectID, tc.clusterName)
			assert.Equal(t, tc.expectedHasErrors, err != nil)
			if err == nil {
				assert.Equal(t, tc.expectedID, id)
			}
		})
	}
}
