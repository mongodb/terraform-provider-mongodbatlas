package unit_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApiPathsParsing(t *testing.T) {
	specParts := unit.ReadAPISpecPaths()
	assert.GreaterOrEqual(t, len(specParts), 5)
	assert.Contains(t, specParts, "GET")
	processArgsPath := "/api/atlas/v2/groups/6746ceed6f62fc3c122a3e0e/clusters/test-acc-tf-c-7871793563057636102/processArgs"
	getPaths := specParts["GET"]
	found1, err := unit.FindNormalizedPath(processArgsPath, &getPaths)
	require.NoError(t, err)
	assert.Equal(t, "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/processArgs", found1.Path)
	variables := found1.Variables(processArgsPath)
	assert.Equal(t, "6746ceed6f62fc3c122a3e0e", variables["groupId"])
	assert.Equal(t, "test-acc-tf-c-7871793563057636102", variables["clusterName"])
}
