package x509authenticationdatabaseuser_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/x509authenticationdatabaseuser"
	"github.com/stretchr/testify/assert"
)

func TestSplitX509ImportID(t *testing.T) {
	const projectID = "664619d870c247237f4b86a6"
	tests := map[string]struct {
		importID  string
		projectID string
		username  string
		legacy    bool
	}{
		"slash format": {
			importID:  projectID + "/username",
			projectID: projectID,
			username:  "username",
		},
		"slash format with slash in username": {
			importID:  projectID + "/CN=team/service",
			projectID: projectID,
			username:  "CN=team/service",
		},
		"project only is valid in both formats": {
			importID:  projectID,
			projectID: projectID,
		},
		"legacy format": {
			importID:  projectID + "-username",
			projectID: projectID,
			username:  "username",
			legacy:    true,
		},
		"legacy format with slash in username": {
			importID:  projectID + "-CN=team/service",
			projectID: projectID,
			username:  "CN=team/service",
			legacy:    true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			projectID, username, legacy := x509authenticationdatabaseuser.SplitX509ImportID(tc.importID)
			assert.Equal(t, tc.projectID, projectID)
			assert.Equal(t, tc.username, username)
			assert.Equal(t, tc.legacy, legacy)
		})
	}
}
