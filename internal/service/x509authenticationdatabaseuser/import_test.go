package x509authenticationdatabaseuser_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/x509authenticationdatabaseuser"
	"github.com/stretchr/testify/assert"
)

func TestSplitX509ImportID(t *testing.T) {
	tests := map[string]struct {
		importID  string
		projectID string
		username  string
		legacy    bool
	}{
		"slash format": {
			importID:  "project/username",
			projectID: "project",
			username:  "username",
		},
		"project only is valid in both formats": {
			importID:  "project",
			projectID: "project",
		},
		"legacy format": {
			importID:  "project-username",
			projectID: "project",
			username:  "username",
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
