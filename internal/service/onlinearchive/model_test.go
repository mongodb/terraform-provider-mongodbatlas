package onlinearchive_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/onlinearchive"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitOnlineArchiveImportID(t *testing.T) {
	const errorString = "import format error to import a MongoDB Atlas Online Archive, use the format {project_id}-{cluster_name}-{archive_id}"
	tests := map[string]struct {
		importID    string
		projectID   string
		clusterName string
		archiveID   string
		errorString string
	}{
		"cluster name without hyphens": {
			importID:    "664619d870c247237f4b86a6-Cluster0-664619d870c247237f4b86a7",
			projectID:   "664619d870c247237f4b86a6",
			clusterName: "Cluster0",
			archiveID:   "664619d870c247237f4b86a7",
		},
		"cluster name with one hyphen": {
			importID:    "664619d870c247237f4b86a6-my-cluster-664619d870c247237f4b86a7",
			projectID:   "664619d870c247237f4b86a6",
			clusterName: "my-cluster",
			archiveID:   "664619d870c247237f4b86a7",
		},
		"cluster name with several hyphens": {
			importID:    "664619d870c247237f4b86a6-my-prod-cluster-01-664619d870c247237f4b86a7",
			projectID:   "664619d870c247237f4b86a6",
			clusterName: "my-prod-cluster-01",
			archiveID:   "664619d870c247237f4b86a7",
		},
		"missing archive ID": {
			importID:    "664619d870c247237f4b86a6-Cluster0",
			errorString: errorString,
		},
		"empty input": {
			importID:    "",
			errorString: errorString,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			projectID, clusterName, archiveID, err := onlinearchive.SplitOnlineArchiveImportID(tc.importID)
			if tc.errorString != "" {
				require.EqualError(t, err, tc.errorString)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.projectID, projectID)
			assert.Equal(t, tc.clusterName, clusterName)
			assert.Equal(t, tc.archiveID, archiveID)
		})
	}
}
