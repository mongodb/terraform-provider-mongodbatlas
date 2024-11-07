package conversion_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateProjectID(t *testing.T) {
	require.NoError(t, conversion.ValidateProjectID("123456789012345678901234"))
	err := conversion.ValidateProjectID("invalid_project_id")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid_project_id")
	assert.Contains(t, err.Error(), "project_id must be a 24 character hex string")
}

func TestValidateClusterName(t *testing.T) {
	require.NoError(t, conversion.ValidateClusterName("clusterName"))
	err := conversion.ValidateClusterName("_invalidClusterName")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "_invalidClusterName")
	assert.Contains(t, err.Error(), "cluster_name must be a string with length between 1 and 64, starting and ending with an alphanumeric character, and containing only alphanumeric characters and hyphens")
}

func TestGetProjectID(t *testing.T) {
	c := &conversion.ClusterImportAttrNames{
		ProjectID:   "test_project_id",
		ClusterName: "test_cluster_name",
	}
	assert.Equal(t, "test_project_id", c.GetProjectID())
	assert.Equal(t, "test_cluster_name", c.GetClusterName())

	getProjectID := func(names *conversion.ClusterImportAttrNames) string {
		return names.GetProjectID()
	}
	assert.Equal(t, "project_id", getProjectID(nil))
	getClusterName := func(names *conversion.ClusterImportAttrNames) string {
		return names.GetClusterName()
	}
	assert.Equal(t, "cluster_name", getClusterName(nil))
}
